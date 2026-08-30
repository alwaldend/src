import { promises as dns } from "node:dns";
import * as net from "node:net";
import { performance } from "node:perf_hooks";
import * as tls from "node:tls";

export const description = "Bounded DNS, TCP/TLS, and HTTP diagnostics.";

const plugin = {
  name: "network_probe",
  description,

  apply(ctx) {
    const MAX_BODY_BYTES = 1024 * 1024;
    const RECORD_TYPES = new Set([
      "A",
      "AAAA",
      "CAA",
      "CNAME",
      "MX",
      "NAPTR",
      "NS",
      "PTR",
      "SOA",
      "SRV",
      "TXT",
    ]);

    function integer(value, fallback, minimum, maximum, label) {
      const selected = value === undefined ? fallback : value;
      if (!Number.isInteger(selected)) {
        throw new TypeError(`${label} must be an integer`);
      }
      if (selected < minimum || selected > maximum) {
        throw new RangeError(
          `${label} must be between ${minimum} and ${maximum}`,
        );
      }
      return selected;
    }

    function string(value, label, maximum, allowEmpty = false) {
      if (typeof value !== "string") {
        throw new TypeError(`${label} must be a string`);
      }
      if (/[\0\r\n]/.test(value)) {
        throw new TypeError(`${label} contains an invalid control character`);
      }
      if ((!allowEmpty && value.length === 0) || value.length > maximum) {
        throw new RangeError(`${label} has an invalid length`);
      }
      return value;
    }

    function host(value) {
      const selected = string(value, "host", 253);
      if (/\s/.test(selected)) {
        throw new TypeError("host must not contain whitespace");
      }
      return selected;
    }

    function clipUtf8(value, maximum) {
      const source = Buffer.from(String(value), "utf8");
      if (source.length <= maximum) {
        return {
          text: source.toString("utf8"),
          bytes: source.length,
          truncated: false,
        };
      }
      let end = maximum;
      while (end > 0 && (source[end] & 0xc0) === 0x80) {
        end -= 1;
      }
      return {
        text: source.subarray(0, end).toString("utf8"),
        bytes: end,
        truncated: true,
      };
    }

    function decodeUtf8Prefix(source) {
      let start = source.length - 1;
      while (start >= 0 && (source[start] & 0xc0) === 0x80) {
        start -= 1;
      }
      if (start < 0) return "";
      const leading = source[start];
      const expected = leading < 0x80
        ? 1
        : leading >= 0xc2 && leading <= 0xdf
          ? 2
          : leading >= 0xe0 && leading <= 0xef
            ? 3
            : leading >= 0xf0 && leading <= 0xf4 ? 4 : 1;
      const end = source.length - start < expected ? start : source.length;
      return source.subarray(0, end).toString("utf8");
    }

    function errorDetails(error) {
      return {
        code: typeof error?.code === "string" ? error.code : null,
        message: clipUtf8(error?.message ?? String(error), 2048).text,
      };
    }

    function elapsed(started) {
      return Math.round((performance.now() - started) * 1000) / 1000;
    }

    function boundedRecords(records, maximumRecords, maximumBytes) {
      const selected = [];
      let truncated = records.length > maximumRecords;
      for (const record of records.slice(0, maximumRecords)) {
        selected.push(record);
        if (Buffer.byteLength(JSON.stringify(selected), "utf8") > maximumBytes) {
          selected.pop();
          truncated = true;
          break;
        }
      }
      return { records: selected, truncated };
    }

    ctx.tool(
      {
        name: "dns_lookup",
        description:
          "Resolve one DNS record type with an optional resolver and hard bounds.",
        inputSchema: {
          type: "object",
          additionalProperties: false,
          required: ["host"],
          properties: {
            host: { type: "string", minLength: 1, maxLength: 253 },
            rrtype: {
              type: "string",
              enum: [...RECORD_TYPES],
              default: "A",
            },
            resolver: { type: "string", minLength: 1, maxLength: 300 },
            timeout_ms: {
              type: "integer",
              minimum: 100,
              maximum: 30_000,
              default: 5_000,
            },
            max_records: {
              type: "integer",
              minimum: 1,
              maximum: 100,
              default: 25,
            },
            max_bytes: {
              type: "integer",
              minimum: 1024,
              maximum: 256 * 1024,
              default: 64 * 1024,
            },
          },
        },
      },
      async (input) => {
        const hostname = host(input.host);
        const rrtype = String(input.rrtype ?? "A").toUpperCase();
        if (!RECORD_TYPES.has(rrtype)) {
          throw new RangeError(`unsupported rrtype: ${rrtype}`);
        }
        const timeoutMs = integer(
          input.timeout_ms,
          5_000,
          100,
          30_000,
          "timeout_ms",
        );
        const maxRecords = integer(
          input.max_records,
          25,
          1,
          100,
          "max_records",
        );
        const maxBytes = integer(
          input.max_bytes,
          64 * 1024,
          1024,
          256 * 1024,
          "max_bytes",
        );
        const resolver = new dns.Resolver({ timeout: timeoutMs, tries: 1 });
        if (input.resolver !== undefined) {
          resolver.setServers([string(input.resolver, "resolver", 300)]);
        }

        const started = performance.now();
        let timer;
        try {
          const timeout = new Promise((_, reject) => {
            timer = setTimeout(() => {
              const error = new Error(`DNS lookup timed out after ${timeoutMs}ms`);
              error.code = "ETIMEDOUT";
              resolver.cancel();
              reject(error);
            }, timeoutMs);
          });
          const resolved = await Promise.race([
            resolver.resolve(hostname, rrtype),
            timeout,
          ]);
          const records = Array.isArray(resolved) ? resolved : [resolved];
          const bounded = boundedRecords(records, maxRecords, maxBytes);
          return {
            ok: true,
            host: hostname,
            rrtype,
            resolver: input.resolver ?? null,
            records: bounded.records,
            truncated: bounded.truncated,
            elapsedMs: elapsed(started),
          };
        } catch (error) {
          return {
            ok: false,
            host: hostname,
            rrtype,
            resolver: input.resolver ?? null,
            error: errorDetails(error),
            elapsedMs: elapsed(started),
          };
        } finally {
          clearTimeout(timer);
        }
      },
    );

    ctx.tool(
      {
        name: "tcp_probe",
        description:
          "Open a bounded TCP or TLS connection without sending application data.",
        inputSchema: {
          type: "object",
          additionalProperties: false,
          required: ["host", "port"],
          properties: {
            host: { type: "string", minLength: 1, maxLength: 253 },
            port: { type: "integer", minimum: 1, maximum: 65535 },
            timeout_ms: {
              type: "integer",
              minimum: 100,
              maximum: 30_000,
              default: 5_000,
            },
            tls: { type: "boolean", default: false },
            server_name: { type: "string", minLength: 1, maxLength: 253 },
            verify_certificate: { type: "boolean", default: true },
          },
        },
      },
      async (input) => {
        const hostname = host(input.host);
        const port = integer(input.port, undefined, 1, 65535, "port");
        const timeoutMs = integer(
          input.timeout_ms,
          5_000,
          100,
          30_000,
          "timeout_ms",
        );
        const useTls = input.tls === true;
        const serverName = input.server_name === undefined
          ? (net.isIP(hostname) ? undefined : hostname)
          : host(input.server_name);
        const started = performance.now();

        return new Promise((resolve) => {
          let settled = false;
          let timer;
          const options = { host: hostname, port };
          const socket = useTls
            ? tls.connect({
                ...options,
                servername: serverName,
                rejectUnauthorized: input.verify_certificate !== false,
              })
            : net.connect(options);

          const finish = (output) => {
            if (settled) {
              return;
            }
            settled = true;
            clearTimeout(timer);
            socket.removeAllListeners();
            socket.destroy();
            resolve({
              host: hostname,
              port,
              tls: useTls,
              elapsedMs: elapsed(started),
              ...output,
            });
          };

          timer = setTimeout(() => {
            finish({
              ok: false,
              error: {
                code: "ETIMEDOUT",
                message: `connection timed out after ${timeoutMs}ms`,
              },
            });
          }, timeoutMs);
          socket.once("error", (error) => {
            finish({ ok: false, error: errorDetails(error) });
          });
          socket.once(useTls ? "secureConnect" : "connect", () => {
            const output = {
              ok: true,
              localAddress: socket.localAddress ?? null,
              localPort: socket.localPort ?? null,
              remoteAddress: socket.remoteAddress ?? null,
              remoteFamily: socket.remoteFamily ?? null,
              remotePort: socket.remotePort ?? null,
            };
            if (useTls) {
              const certificate = socket.getPeerCertificate(false);
              output.authorized = socket.authorized;
              output.authorizationError = socket.authorizationError ?? null;
              output.protocol = socket.getProtocol();
              output.cipher = socket.getCipher()?.standardName ??
                socket.getCipher()?.name ?? null;
              output.alpnProtocol = socket.alpnProtocol || null;
              output.certificate = certificate &&
                  Object.keys(certificate).length > 0
                ? {
                    subject: certificate.subject ?? null,
                    issuer: certificate.issuer ?? null,
                    validFrom: certificate.valid_from ?? null,
                    validTo: certificate.valid_to ?? null,
                    fingerprint256: certificate.fingerprint256 ?? null,
                  }
                : null;
            }
            finish(output);
          });
        });
      },
    );

    ctx.tool(
      {
        name: "http_probe",
        description:
          "Issue a bounded HTTP GET or HEAD request and return metadata plus a body preview.",
        inputSchema: {
          type: "object",
          additionalProperties: false,
          required: ["url"],
          properties: {
            url: { type: "string", minLength: 1, maxLength: 8192 },
            method: { type: "string", enum: ["GET", "HEAD"], default: "GET" },
            headers: {
              type: "object",
              maxProperties: 32,
              additionalProperties: {
                type: "string",
                maxLength: 8192,
              },
            },
            timeout_ms: {
              type: "integer",
              minimum: 100,
              maximum: 60_000,
              default: 10_000,
            },
            max_body_bytes: {
              type: "integer",
              minimum: 0,
              maximum: MAX_BODY_BYTES,
              default: 64 * 1024,
            },
            follow_redirects: { type: "boolean", default: false },
            include_sensitive_headers: { type: "boolean", default: false },
          },
        },
      },
      async (input) => {
        const selectedUrl = string(input.url, "url", 8192);
        const url = new URL(selectedUrl);
        if (url.protocol !== "http:" && url.protocol !== "https:") {
          throw new TypeError("url must use http or https");
        }
        if (url.username || url.password) {
          throw new TypeError("credentials must not be embedded in url");
        }
        const method = String(input.method ?? "GET").toUpperCase();
        if (method !== "GET" && method !== "HEAD") {
          throw new RangeError("method must be GET or HEAD");
        }
        const headers = input.headers ?? {};
        if (headers === null || Array.isArray(headers) ||
            typeof headers !== "object") {
          throw new TypeError("headers must be an object");
        }
        const entries = Object.entries(headers);
        if (entries.length > 32) {
          throw new RangeError("headers must contain at most 32 entries");
        }
        for (const [name, value] of entries) {
          string(name, "header name", 256);
          string(value, `header ${name}`, 8192, true);
        }
        const timeoutMs = integer(
          input.timeout_ms,
          10_000,
          100,
          60_000,
          "timeout_ms",
        );
        const maximum = integer(
          input.max_body_bytes,
          64 * 1024,
          0,
          MAX_BODY_BYTES,
          "max_body_bytes",
        );
        const controller = new AbortController();
        const timer = setTimeout(() => {
          controller.abort(new Error(`HTTP probe timed out after ${timeoutMs}ms`));
        }, timeoutMs);
        const started = performance.now();

        try {
          const response = await fetch(url, {
            method,
            headers,
            redirect: input.follow_redirects === true ? "follow" : "manual",
            signal: controller.signal,
          });
          const responseHeaders = {};
          const sensitive = new Set([
            "proxy-authenticate",
            "set-cookie",
            "www-authenticate",
          ]);
          let headerBytes = 0;
          let headersTruncated = false;
          for (const [name, value] of response.headers) {
            if (sensitive.has(name) && input.include_sensitive_headers !== true) {
              responseHeaders[name] = "[redacted]";
              continue;
            }
            const clipped = clipUtf8(value, 8192);
            const bytes = Buffer.byteLength(name, "utf8") + clipped.bytes;
            if (headerBytes + bytes > 64 * 1024) {
              headersTruncated = true;
              break;
            }
            responseHeaders[name] = clipped.text;
            headerBytes += bytes;
            headersTruncated ||= clipped.truncated;
          }

          const chunks = [];
          let keptBytes = 0;
          let bodyTruncated = false;
          if (response.body && method !== "HEAD") {
            const reader = response.body.getReader();
            while (true) {
              const { done, value } = await reader.read();
              if (done) {
                break;
              }
              const chunk = Buffer.from(value);
              const remaining = maximum - keptBytes;
              if (remaining > 0) {
                chunks.push(chunk.subarray(0, remaining));
                keptBytes += Math.min(chunk.length, remaining);
              }
              if (chunk.length > remaining) {
                bodyTruncated = true;
                await reader.cancel("body preview limit reached");
                break;
              }
            }
          }
          const body = Buffer.concat(chunks, keptBytes);
          const contentType = response.headers.get("content-type") ?? "";
          const textual = /^(text\/)|json|javascript|xml|x-www-form-urlencoded/i
            .test(contentType);
          return {
            ok: true,
            requestedUrl: url.toString(),
            finalUrl: response.url,
            redirected: response.redirected,
            status: response.status,
            statusText: response.statusText,
            headers: responseHeaders,
            headersTruncated,
            body: textual
              ? bodyTruncated
                ? decodeUtf8Prefix(body)
                : body.toString("utf8")
              : body.toString("base64"),
            bodyEncoding: textual ? "utf8" : "base64",
            bodyBytes: keptBytes,
            bodyTruncated,
            elapsedMs: elapsed(started),
          };
        } catch (error) {
          return {
            ok: false,
            requestedUrl: url.toString(),
            error: errorDetails(error),
            elapsedMs: elapsed(started),
          };
        } finally {
          clearTimeout(timer);
        }
      },
    );
  },
};

plugin.apply.description = description;

export default plugin;
