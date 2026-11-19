package com.alwaldend.src.projects.activitywatch_ingester_android

import com.sun.net.httpserver.HttpServer
import java.io.ByteArrayOutputStream
import java.io.InputStream
import java.net.InetSocketAddress
import java.nio.charset.StandardCharsets
import java.time.Duration
import java.time.Instant
import java.util.concurrent.Executor
import org.json.JSONObject
import org.junit.Assert.assertEquals
import org.junit.Test
import org.junit.runner.RunWith
import org.robolectric.RobolectricTestRunner

@RunWith(RobolectricTestRunner::class)
class ActivityWatchIngesterTest {
  private data class ReceivedRequest(
      val method: String,
      val path: String,
      val authorization: String? = null,
      val body: String,
  )

  @Test
  fun `sends bucket create and heartbeat to the configured server`() {
    val server = HttpServer.create(InetSocketAddress("127.0.0.1", 0), 0)
    val received = ArrayList<ReceivedRequest>()

    server.createContext("/") { exchange ->
      val body = readAll(exchange.requestBody)
      synchronized(received) {
        received.add(
            ReceivedRequest(
                method = exchange.requestMethod,
                path = exchange.requestURI.path,
                authorization = exchange.requestHeaders.getFirst("Authorization"),
                body = body,
            ))
      }
      exchange.sendResponseHeaders(200, -1)
      exchange.close()
    }
    server.start()
    try {
      val ingester =
          ActivityWatchIngester(
              baseUrl = "http://127.0.0.1:${server.address.port}",
              apiKey = "test-secret",
              now = { Instant.parse("2026-09-08T00:00:00Z") },
              // Run synchronously so the test observes completion deterministically.
              executor = Executor { task -> task.run() },
          )
      ingester.send("chrome", "https://example.com/page", Duration.ofSeconds(5))
    } finally {
      server.stop(0)
    }

    val bucketRequest = received.firstOrNull { it.path == "/api/0/buckets/web.tab.current" }
    val heartbeatRequest =
        received.firstOrNull { it.path == "/api/0/buckets/web.tab.current/heartbeat" }
    assertEquals("bucket create must be a POST", "POST", bucketRequest?.method)
    assertEquals("heartbeat must be a POST", "POST", heartbeatRequest?.method)
    assertEquals("Bearer test-secret", bucketRequest?.authorization)
    assertEquals(
        """{"client":"aw-watcher-android-web","type":"currenttab","hostname":"android"}""",
        bucketRequest?.body)
    val heartbeatJson = JSONObject(heartbeatRequest?.body ?: "")
    assertEquals("example.com/page", heartbeatJson.getJSONObject("data").getString("url"))
    assertEquals(5.0, heartbeatJson.getDouble("duration"), 0.0)
  }

  private fun readAll(input: InputStream): String {
    val buffer = ByteArrayOutputStream()
    input.copyTo(buffer)
    return buffer.toString(StandardCharsets.UTF_8.name())
  }
}
