import * as fs from "node:fs/promises";
import { constants as fsConstants } from "node:fs";
import * as path from "node:path";
import process from "node:process";

export const description = "Bounded repository context, reads, and searches.";

const plugin = {
  name: "repo_context",
  description,

  apply(ctx) {
    const root = path.resolve(ctx.resolveWorkspace("."));
    const canonicalRoot = fs.realpath(root);
    const MAX_FILE_BYTES = 2 * 1024 * 1024;
    const MAX_OUTPUT_BYTES = 1024 * 1024;

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
      if (value.includes("\0")) {
        throw new TypeError(`${label} must not contain a NUL byte`);
      }
      if ((!allowEmpty && value.length === 0) || value.length > maximum) {
        throw new RangeError(`${label} has an invalid length`);
      }
      return value;
    }

    function workspacePath(value, label = "path") {
      const requested = string(value ?? ".", label, 4096);
      if (path.isAbsolute(requested)) {
        throw new TypeError(`${label} must be relative to the workspace`);
      }
      const absolute = path.resolve(ctx.resolveWorkspace(requested));
      const relative = path.relative(root, absolute);
      if (escapes(root, absolute)) {
        throw new RangeError(`${label} resolves outside the workspace`);
      }
      return {
        absolute,
        relative: relative.split(path.sep).join("/") || ".",
      };
    }

    function clipUtf8(value, maximum) {
      const source = Buffer.from(String(value), "utf8");
      if (source.length <= maximum) {
        return { text: source.toString("utf8"), bytes: source.length,
          truncated: false };
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

    function resultText(result, field) {
      const value = result[field];
      if (Buffer.isBuffer(value)) {
        return value.toString("utf8");
      }
      return value === undefined ? "" : String(value);
    }

    function textLines(content) {
      if (content.length === 0) return [];
      const lines = content.split(/\r?\n/);
      if (content.endsWith("\n")) lines.pop();
      return lines;
    }

    function executionState(result, label) {
      if (!result || typeof result !== "object") {
        throw new Error(`${label} returned an invalid execution result`);
      }
      const outputLimitExceeded = result.outputLimitExceeded === true;
      if (result.truncated && !outputLimitExceeded) {
        throw new Error(
          `${label} returned incomplete output without an output-limit marker`,
        );
      }
      if (outputLimitExceeded) {
        return { code: null, outputLimitExceeded: true };
      }
      if (result.signal !== null && result.signal !== undefined) {
        throw new Error(`${label} was terminated by ${result.signal}`);
      }
      const code = result.code ?? result.exitCode;
      if (!Number.isInteger(code)) {
        throw new Error(`${label} ended without an exit code`);
      }
      return { code, outputLimitExceeded: false };
    }

    function escapes(base, candidate) {
      const relative = path.relative(base, candidate);
      return relative === ".." || relative.startsWith(`..${path.sep}`) ||
        path.isAbsolute(relative);
    }

    async function canonicalPath(absolute, label) {
      const [workspace, selected] = await Promise.all([
        canonicalRoot,
        fs.realpath(absolute),
      ]);
      if (escapes(workspace, selected)) {
        throw new RangeError(`${label} resolves outside the workspace`);
      }
      return selected;
    }

    async function bindPath(absolute, label) {
      const canonical = await canonicalPath(absolute, label);
      const handle = await fs.open(
        canonical,
        fsConstants.O_RDONLY | fsConstants.O_NOFOLLOW,
      );
      try {
        const [workspace, opened, info] = await Promise.all([
          canonicalRoot,
          fs.realpath(`/proc/self/fd/${handle.fd}`),
          handle.stat(),
        ]);
        if (escapes(workspace, opened)) {
          throw new RangeError(`${label} resolves outside the workspace`);
        }
        return {
          absolute: opened,
          commandPath: `/proc/${process.pid}/fd/${handle.fd}`,
          handle,
          info,
        };
      } catch (error) {
        await handle.close();
        throw error;
      }
    }

    async function run(file, args, options = {}) {
      return ctx.exec(file, args, {
        cwd: options.cwd ?? root,
        ...(file === "git" ? {
          env: {
            GIT_ALTERNATE_OBJECT_DIRECTORIES: null,
            GIT_COMMON_DIR: null,
            GIT_CONFIG: null,
            GIT_CONFIG_COUNT: null,
            GIT_CONFIG_PARAMETERS: null,
            GIT_DIR: null,
            GIT_GRAFT_FILE: null,
            GIT_IMPLICIT_WORK_TREE: null,
            GIT_INDEX_FILE: null,
            GIT_INTERNAL_SUPER_PREFIX: null,
            GIT_NAMESPACE: null,
            GIT_NO_LAZY_FETCH: "1",
            GIT_NO_REPLACE_OBJECTS: null,
            GIT_OBJECT_DIRECTORY: null,
            GIT_OPTIONAL_LOCKS: "0",
            GIT_PREFIX: null,
            GIT_REPLACE_REF_BASE: null,
            GIT_SHALLOW_FILE: null,
            GIT_TERMINAL_PROMPT: "0",
            GIT_WORK_TREE: null,
          },
        } : {}),
        allowTruncatedOutput: true,
        timeoutMs: options.timeoutMs ?? 10_000,
        maxBytes: options.maxBytes ?? 512 * 1024,
      });
    }

    async function readFile(location, maximum = MAX_FILE_BYTES) {
      const binding = await bindPath(
        location.absolute,
        location.relative,
      );
      try {
        if (!binding.info.isFile()) {
          throw new TypeError(`${location.relative} is not a regular file`);
        }
        if (binding.info.size > maximum) {
          throw new RangeError(
            `${location.relative} is ${binding.info.size} bytes; ` +
              `limit is ${maximum}`,
          );
        }
        const buffer = Buffer.allocUnsafe(maximum + 1);
        let bytes = 0;
        while (bytes < buffer.length) {
          const result = await binding.handle.read(
            buffer,
            bytes,
            buffer.length - bytes,
            bytes,
          );
          if (result.bytesRead === 0) break;
          bytes += result.bytesRead;
        }
        if (bytes > maximum) {
          throw new RangeError(`${location.relative} exceeded its byte limit`);
        }
        const raw = buffer.subarray(0, bytes);
        return {
          content: raw.toString("utf8"),
          bytes,
          raw,
        };
      } finally {
        await binding.handle.close();
      }
    }

    async function fallbackSearch({
      query,
      paths,
      fixed,
      caseSensitive,
      contextLines,
      maxMatches,
    }) {
      const MAX_SCAN_FILES = 10_000;
      const MAX_SCAN_BYTES = 64 * 1024 * 1024;
      const MAX_RESULT_BYTES = 4 * 1024 * 1024;
      const matches = [];
      let matchCount = 0;
      let scannedBytes = 0;
      let scannedFiles = 0;
      let resultBytes = 2;
      let truncated = false;
      let detailsTruncated = false;

      if (!fixed) {
        throw new Error(
          "ripgrep is required for regular-expression searches",
        );
      }
      const fixedExpression = !caseSensitive
        ? new RegExp(
          query.replace(/[.*+?^${}()|[\]\\]/g, "\\$&"),
          "gidu",
        )
        : undefined;

      function emit(entry) {
        const bytes = Buffer.byteLength(JSON.stringify(entry), "utf8") +
          (matches.length === 0 ? 0 : 1);
        if (resultBytes + bytes > MAX_RESULT_BYTES) {
          truncated = true;
          return false;
        }
        matches.push(entry);
        resultBytes += bytes;
        return true;
      }

      function submatches(line) {
        const found = [];
        const occurrences = [];
        if (caseSensitive) {
          let offset = 0;
          while (true) {
            const start = line.indexOf(query, offset);
            if (start < 0) break;
            occurrences.push({ start, text: query });
            offset = start + query.length;
          }
        } else {
          fixedExpression.lastIndex = 0;
          for (const match of line.matchAll(fixedExpression)) {
            occurrences.push({ start: match.index, text: match[0] });
          }
        }
        for (const occurrence of occurrences) {
          if (found.length >= 32) {
            detailsTruncated = true;
            break;
          }
          const text = clipUtf8(occurrence.text, 512);
          const byteStart = Buffer.byteLength(
            line.slice(0, occurrence.start),
            "utf8",
          );
          const byteEnd = byteStart + Buffer.byteLength(
            occurrence.text,
            "utf8",
          );
          detailsTruncated ||= text.truncated;
          found.push({
            start: byteStart,
            end: byteEnd,
            text: text.text,
          });
        }
        return found;
      }

      async function searchFile(relativePath) {
        if (truncated) return;
        const normalized = relativePath.split(path.sep).join("/");
        const location = workspacePath(normalized);
        const info = await fs.stat(location.absolute);
        if (!info.isFile() || info.size > MAX_FILE_BYTES) return;
        if (scannedFiles >= MAX_SCAN_FILES ||
            scannedBytes + info.size > MAX_SCAN_BYTES) {
          truncated = true;
          return;
        }
        scannedFiles += 1;
        scannedBytes += info.size;
        const file = await readFile(location);
        if (!Buffer.from(file.content, "utf8").equals(file.raw)) {
          throw new Error(
            "ripgrep is required to search files containing invalid UTF-8",
          );
        }
        if (file.content.includes("\0")) return;
        const lines = textLines(file.content);
        const foundByLine = lines.map((line) => submatches(line));
        const selectedMatches = [];
        let omittedMatch = false;
        for (let index = 0; index < lines.length; index += 1) {
          if (foundByLine[index].length === 0) continue;
          if (matchCount + selectedMatches.length >= maxMatches) {
            omittedMatch = true;
            continue;
          }
          selectedMatches.push(index);
        }

        const selectedSet = new Set(selectedMatches);
        const outputLines = new Set();
        for (const index of selectedMatches) {
          const start = Math.max(0, index - contextLines);
          const end = Math.min(lines.length - 1, index + contextLines);
          for (let outputIndex = start; outputIndex <= end; outputIndex += 1) {
            if (foundByLine[outputIndex].length === 0 ||
                selectedSet.has(outputIndex)) {
              outputLines.add(outputIndex);
            }
          }
        }

        const orderedLines = [...outputLines].sort((left, right) => left - right);
        for (const index of orderedLines) {
          const isMatch = selectedSet.has(index);
          const clipped = clipUtf8(lines[index], 4096);
          detailsTruncated ||= clipped.truncated;
          const emitted = emit({
            kind: isMatch ? "match" : "context",
            path: normalized,
            line: index + 1,
            text: clipped.text,
            textTruncated: clipped.truncated,
            ...(isMatch ? { submatches: foundByLine[index] } : {}),
          });
          if (!emitted) break;
          if (isMatch) matchCount += 1;
        }
        truncated ||= omittedMatch;
      }

      async function visit(relativePath, selected = false) {
        if (truncated) return;
        const location = workspacePath(relativePath);
        let info = await fs.lstat(location.absolute);
        if (info.isSymbolicLink()) {
          if (!selected) return;
          await canonicalPath(location.absolute, location.relative);
          info = await fs.stat(location.absolute);
        }
        if (info.isFile()) {
          await searchFile(location.relative);
          return;
        }
        if (!info.isDirectory()) return;
        if (selected) {
          throw new Error(
            "ripgrep is required to search directories; the JavaScript " +
              "fallback accepts only explicitly selected files",
          );
        }
        const entries = await fs.readdir(location.absolute, {
          withFileTypes: true,
        });
        entries.sort((left, right) => left.name.localeCompare(right.name));
        for (const entry of entries) {
          if (entry.name === ".git" || entry.name === "node_modules") {
            continue;
          }
          await visit(path.join(location.relative, entry.name));
          if (truncated) break;
        }
      }

      for (const selectedPath of paths) {
        await visit(selectedPath, true);
        if (truncated) break;
      }
      return {
        query,
        matches,
        matchCount,
        truncated: truncated || detailsTruncated,
        engine: "javascript",
        scannedFiles,
        scannedBytes,
        resultBytes,
      };
    }

    async function gitContext(directory, maximum) {
      const workspace = await canonicalRoot;
      const selected = await bindPath(directory, "selected Git directory");
      let rootResult;
      try {
        rootResult = await run(
          "git",
          ["-C", selected.commandPath, "rev-parse", "--show-toplevel"],
          { maxBytes: 16 * 1024 },
        );
      } finally {
        await selected.handle.close();
      }
      const rootState = executionState(rootResult, "git rev-parse");
      if (rootState.outputLimitExceeded) {
        return {
          error: "Git repository root discovery output was truncated",
          rootTruncated: true,
        };
      }
      if (rootState.code !== 0) {
        return null;
      }

      const rootText = resultText(rootResult, "stdout").trim();
      if (!path.isAbsolute(rootText)) {
        throw new Error("git rev-parse returned an invalid repository root");
      }
      let repository;
      try {
        repository = await bindPath(path.resolve(rootText), "Git repository");
      } catch (error) {
        if (error instanceof RangeError) {
          return {
            error: "the selected path belongs to a repository outside the " +
              "workspace",
          };
        }
        throw error;
      }
      const relative = path.relative(workspace, repository.absolute);

      let headResult;
      let statusResult;
      try {
        [headResult, statusResult] = await Promise.all([
          run(
            "git",
            [
              "-C",
              repository.commandPath,
              "rev-parse",
              "--verify",
              "HEAD",
            ],
            { maxBytes: 16 * 1024 },
          ),
          run(
            "git",
            [
              "-C",
              repository.commandPath,
              "status",
              "--short",
              "--branch",
              "--untracked-files=no",
            ],
            { maxBytes: maximum },
          ),
        ]);
      } finally {
        await repository.handle.close();
      }
      const headState = executionState(headResult, "git rev-parse HEAD");
      const statusState = executionState(statusResult, "git status");
      const output = {
        root: relative.split(path.sep).join("/") || ".",
        rootTruncated: false,
      };
      if (headState.outputLimitExceeded) {
        output.head = null;
        output.headAvailable = false;
        output.headTruncated = true;
      } else if (headState.code === 0) {
        const head = resultText(headResult, "stdout").trim();
        if (!/^[0-9a-f]{40,64}$/i.test(head)) {
          throw new Error("git rev-parse HEAD returned an invalid hash");
        }
        output.head = head;
        output.headAvailable = true;
        output.headTruncated = false;
      } else {
        output.head = null;
        output.headAvailable = false;
        output.headTruncated = false;
        output.headError = clipUtf8(
          resultText(headResult, "stderr"),
          4096,
        ).text;
      }

      if (!statusState.outputLimitExceeded && statusState.code !== 0) {
        output.status = null;
        output.statusAvailable = false;
        output.statusTruncated = false;
        output.statusError = clipUtf8(
          resultText(statusResult, "stderr"),
          4096,
        ).text;
      } else {
        const clipped = clipUtf8(
          resultText(statusResult, "stdout"),
          maximum,
        );
        output.status = clipped.text;
        output.statusAvailable = true;
        output.statusTruncated =
          clipped.truncated || statusState.outputLimitExceeded;
      }
      return output;
    }

    async function instructionContext(directory, maximum) {
      const workspace = await canonicalRoot;
      const directories = [];
      let current = directory;
      while (true) {
        directories.push(current);
        if (current === workspace) {
          break;
        }
        const parent = path.dirname(current);
        if (parent === current || escapes(workspace, parent)) {
          break;
        }
        current = parent;
      }
      directories.reverse();

      const documents = [];
      let remaining = maximum;
      let truncated = false;
      for (const candidateDirectory of directories) {
        const candidate = path.join(candidateDirectory, "AGENTS.md");
        try {
          const info = await fs.stat(candidate);
          if (!info.isFile()) {
            continue;
          }
          if (remaining === 0 && info.size > 0) {
            truncated = true;
            continue;
          }
          const relative = path.relative(workspace, candidate)
            .split(path.sep).join("/");
          const file = await readFile(
            workspacePath(relative),
            Math.min(MAX_FILE_BYTES, Math.max(info.size, 1)),
          );
          const clipped = clipUtf8(file.content, remaining);
          documents.push({
            path: relative,
            content: clipped.text,
            truncated: clipped.truncated,
          });
          remaining -= clipped.bytes;
          truncated ||= clipped.truncated;
        } catch (error) {
          if (error?.code !== "ENOENT") {
            throw error;
          }
        }
      }
      return {
        documents,
        usedBytes: maximum - remaining,
        truncated,
      };
    }

    ctx.tool(
      {
        name: "repo_context_get",
        description:
          "Describe a workspace path, its Git state, and applicable " +
          "AGENTS.md files.",
        inputSchema: {
          type: "object",
          additionalProperties: false,
          properties: {
            path: { type: "string", default: ".", maxLength: 4096 },
            include_git: { type: "boolean", default: true },
            include_instructions: { type: "boolean", default: true },
            max_bytes: {
              type: "integer",
              minimum: 1024,
              maximum: MAX_OUTPUT_BYTES,
              default: 128 * 1024,
            },
          },
        },
      },
      async (input = {}) => {
        const location = workspacePath(input.path ?? ".");
        const maximum = integer(
          input.max_bytes,
          128 * 1024,
          1024,
          MAX_OUTPUT_BYTES,
          "max_bytes",
        );
        const selected = await bindPath(
          location.absolute,
          location.relative,
        );
        try {
          const absolute = selected.absolute;
          const info = selected.info;
          const directory = info.isDirectory()
            ? absolute
            : path.dirname(absolute);
          const output = {
            path: location.relative,
            kind: info.isDirectory()
              ? "directory"
              : info.isFile() ? "file" : "other",
            size: info.size,
          };

          let remaining = maximum;
          if (input.include_git !== false) {
            const git = await gitContext(
              directory,
              Math.min(remaining, 64 * 1024),
            );
            output.git = git;
            if (git?.status) {
              remaining -= Buffer.byteLength(git.status, "utf8");
            }
          }
          if (input.include_instructions !== false) {
            const instructions = await instructionContext(directory, remaining);
            output.instructions = instructions.documents;
            output.instructionsTruncated = instructions.truncated;
            remaining -= instructions.usedBytes;
          }

          if (info.isDirectory()) {
            const allEntries = (await fs.readdir(selected.commandPath, {
              withFileTypes: true,
            }))
              .sort((left, right) => left.name.localeCompare(right.name));
            const entries = allEntries.slice(0, 128)
              .map((entry) => ({
                name: entry.name,
                kind: entry.isDirectory()
                  ? "directory"
                  : entry.isFile() ? "file" : "other",
              }));
            output.entries = entries;
            output.entriesTruncated = allEntries.length > 128;
          }
          output.maxBytes = maximum;
          output.usedTextBytes = maximum - remaining;
          return output;
        } finally {
          await selected.handle.close();
        }
      },
    );

    ctx.tool(
      {
        name: "repo_context_search",
        description:
          "Search workspace files with ripgrep and return bounded " +
          "structured matches.",
        inputSchema: {
          type: "object",
          additionalProperties: false,
          required: ["query"],
          properties: {
            query: { type: "string", minLength: 1, maxLength: 4096 },
            paths: {
              type: "array",
              maxItems: 32,
              items: { type: "string", minLength: 1, maxLength: 4096 },
              default: ["."],
            },
            globs: {
              type: "array",
              maxItems: 32,
              items: { type: "string", minLength: 1, maxLength: 256 },
              default: [],
            },
            fixed: { type: "boolean", default: true },
            case_sensitive: { type: "boolean", default: true },
            context_lines: {
              type: "integer",
              minimum: 0,
              maximum: 20,
              default: 0,
            },
            max_matches: {
              type: "integer",
              minimum: 1,
              maximum: 500,
              default: 100,
            },
          },
        },
      },
      async (input) => {
        const query = string(input.query, "query", 4096);
        const contextLines = integer(
          input.context_lines,
          0,
          0,
          20,
          "context_lines",
        );
        const maxMatches = integer(
          input.max_matches,
          100,
          1,
          500,
          "max_matches",
        );
        const selectedPaths = input.paths ?? ["."];
        if (!Array.isArray(selectedPaths) || selectedPaths.length > 32) {
          throw new RangeError("paths must contain at most 32 entries");
        }
        const locations = selectedPaths.map((value, index) =>
          workspacePath(value, `paths[${index}]`),
        );
        const paths = locations.map((location) => location.relative);
        const bindings = [];
        try {
          for (let index = 0; index < locations.length; index += 1) {
            const binding = await bindPath(
              locations[index].absolute,
              `paths[${index}]`,
            );
            bindings.push({
              ...binding,
              relative: locations[index].relative,
            });
          }
        } catch (error) {
          await Promise.all(bindings.map((binding) => binding.handle.close()));
          throw error;
        }
        const globs = input.globs ?? [];
        if (!Array.isArray(globs) || globs.length > 32) {
          throw new RangeError("globs must contain at most 32 entries");
        }
        const validatedGlobs = globs.map((value, index) =>
          string(value, `globs[${index}]`, 256),
        );

        const args = [
          "--json",
          "--no-config",
          "--max-columns=4096",
          "--max-columns-preview",
        ];
        if (input.fixed !== false) {
          args.push("--fixed-strings");
        }
        args.push(input.case_sensitive === false
          ? "--ignore-case"
          : "--case-sensitive");
        if (contextLines > 0) {
          args.push("--context", String(contextLines));
        }
        for (const glob of validatedGlobs) {
          args.push("--glob", glob);
        }
        args.push(
          "--",
          query,
          ...bindings.map((binding) => binding.commandPath),
        );

        let execution;
        try {
          try {
            execution = await run("rg", args, {
              maxBytes: 4 * 1024 * 1024,
              timeoutMs: 20_000,
            });
          } catch (error) {
            if (error?.code !== "ENOENT") throw error;
            return await fallbackSearch({
              query,
              paths,
              fixed: input.fixed !== false,
              caseSensitive: input.case_sensitive !== false,
              contextLines,
              maxMatches,
            });
          }
        } finally {
          await Promise.all(bindings.map((binding) => binding.handle.close()));
        }

        function displayPath(reported) {
          for (const binding of bindings) {
            if (reported === binding.commandPath) return binding.relative;
            const prefix = `${binding.commandPath}/`;
            if (!reported.startsWith(prefix)) continue;
            const suffix = reported.slice(prefix.length);
            return binding.relative === "."
              ? suffix
              : `${binding.relative}/${suffix}`;
          }
          return reported;
        }
        const state = executionState(execution, "ripgrep");
        if (!state.outputLimitExceeded &&
            state.code !== 0 && state.code !== 1) {
          throw new Error(
            `ripgrep exited with ${state.code}: ` +
              clipUtf8(resultText(execution, "stderr"), 4096).text,
          );
        }

        const matches = [];
        let matchCount = 0;
        let parseTruncated = state.outputLimitExceeded;
        const stdout = resultText(execution, "stdout");
        const lines = stdout.split("\n");
        if (state.outputLimitExceeded && !stdout.endsWith("\n")) {
          lines.pop();
        }
        for (const line of lines) {
          if (!line) {
            continue;
          }
          let event;
          try {
            event = JSON.parse(line);
          } catch {
            parseTruncated = true;
            continue;
          }
          if (event.type !== "match" && event.type !== "context") {
            continue;
          }
          if (event.type === "match" && matchCount >= maxMatches) {
            parseTruncated = true;
            break;
          }
          const data = event.data;
          const lineBytes = typeof data.lines?.bytes === "string";
          const pathBytes = typeof data.path?.bytes === "string";
          const lineText = clipUtf8(data.lines?.text ?? "", 4096);
          parseTruncated ||= lineText.truncated || lineBytes || pathBytes;
          const entry = {
            kind: event.type,
            path: displayPath(data.path?.text ?? ""),
            pathTruncated: pathBytes,
            line: data.line_number ?? null,
            text: lineText.text.replace(/\r?\n$/, ""),
            textTruncated: lineText.truncated || lineBytes,
          };
          if (event.type === "match") {
            const sourceSubmatches = data.submatches ?? [];
            parseTruncated ||= sourceSubmatches.length > 32;
            entry.submatches = sourceSubmatches.slice(0, 32).map(
              (submatch) => {
                const matchBytes =
                  typeof submatch.match?.bytes === "string";
                const text = clipUtf8(
                  submatch.match?.text ?? "",
                  512,
                );
                parseTruncated ||= text.truncated || matchBytes;
                return {
                  start: submatch.start,
                  end: submatch.end,
                  text: text.text,
                  textTruncated: text.truncated || matchBytes,
                };
              },
            );
            matchCount += 1;
          }
          matches.push(entry);
        }
        return {
          query,
          matches,
          matchCount,
          truncated: parseTruncated,
          engine: "ripgrep",
        };
      },
    );

    ctx.tool(
      {
        name: "repo_context_read",
        description:
          "Read bounded line ranges from regular files inside the workspace.",
        inputSchema: {
          type: "object",
          additionalProperties: false,
          required: ["files"],
          properties: {
            files: {
              type: "array",
              minItems: 1,
              maxItems: 32,
              items: {
                type: "object",
                additionalProperties: false,
                required: ["path"],
                properties: {
                  path: { type: "string", minLength: 1, maxLength: 4096 },
                  start_line: { type: "integer", minimum: 1 },
                  end_line: { type: "integer", minimum: 1 },
                },
              },
            },
            max_total_bytes: {
              type: "integer",
              minimum: 1,
              maximum: MAX_OUTPUT_BYTES,
              default: 128 * 1024,
            },
          },
        },
      },
      async (input) => {
        if (!Array.isArray(input.files) || input.files.length === 0 ||
            input.files.length > 32) {
          throw new RangeError("files must contain between 1 and 32 entries");
        }
        const maximum = integer(
          input.max_total_bytes,
          128 * 1024,
          1,
          MAX_OUTPUT_BYTES,
          "max_total_bytes",
        );
        let remaining = maximum;
        const files = [];
        for (let index = 0; index < input.files.length; index += 1) {
          const request = input.files[index];
          const location = workspacePath(
            request.path,
            `files[${index}].path`,
          );
          const start = integer(
            request.start_line,
            1,
            1,
            Number.MAX_SAFE_INTEGER,
            `files[${index}].start_line`,
          );
          const end = integer(
            request.end_line,
            Number.MAX_SAFE_INTEGER,
            1,
            Number.MAX_SAFE_INTEGER,
            `files[${index}].end_line`,
          );
          if (end < start) {
            throw new RangeError("end_line must not be before start_line");
          }
          const file = await readFile(location);
          const lines = textLines(file.content);
          const rangeExists = start <= lines.length;
          const selected = !rangeExists
            ? ""
            : lines.slice(start - 1, Math.min(end, lines.length)).join("\n");
          const clipped = clipUtf8(selected, remaining);
          const retainedNewlines = clipped.text.match(/\n/g)?.length ?? 0;
          const retainedEnd = clipped.text === ""
            ? null
            : start + retainedNewlines -
              (clipped.truncated && clipped.text.endsWith("\n") ? 1 : 0);
          files.push({
            path: location.relative,
            startLine: start,
            endLine: clipped.truncated
              ? retainedEnd
              : rangeExists ? Math.min(end, lines.length) : null,
            totalLines: lines.length,
            content: clipped.text,
            bytes: clipped.bytes,
            truncated: clipped.truncated,
          });
          remaining -= clipped.bytes;
          if (remaining === 0) {
            break;
          }
        }
        return {
          files,
          maxTotalBytes: maximum,
          usedBytes: maximum - remaining,
          truncated: files.length < input.files.length ||
            files.some((file) => file.truncated),
        };
      },
    );
  },
};

plugin.apply.description = description;

export default plugin;
