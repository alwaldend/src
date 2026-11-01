package com.alwaldend.src.projects.activitywatch_ingester_android

import java.io.IOException
import java.net.HttpURLConnection
import java.net.URL
import java.time.Duration
import java.time.Instant
import java.util.concurrent.Executor
import java.util.concurrent.Executors
import java.util.concurrent.atomic.AtomicReference
import org.json.JSONObject

/**
 * Coalesces queued events to the latest browser URL and sends one bucket create plus one heartbeat
 * per settled value, so a slow or unreachable server cannot accumulate an unbounded backlog of
 * stale heartbeats.
 */
private class HeartbeatWorker : Executor {
  private val latest = AtomicReference<Runnable?>(null)
  private val worker = Executors.newSingleThreadExecutor()

  override fun execute(task: Runnable) {
    latest.set(task)
    worker.execute {
      val pending = latest.getAndSet(null) ?: return@execute
      pending.run()
    }
  }
}

class ActivityWatchIngester(
    private val baseUrl: String,
    private val apiKey: String? = null,
    private val now: () -> Instant = Instant::now,
    private val executor: Executor = HeartbeatWorker(),
) {
  fun send(browser: String, url: String, duration: Duration = Duration.ZERO) {
    val event = ActivityWatchEvent(now(), browser, url, duration)
    executor.execute {
      ensureBucket()
      post(event.toJson())
    }
  }

  private fun ensureBucket() {
    try {
      val connection =
          URL("$baseUrl/api/0/buckets/web.tab.current").openConnection() as HttpURLConnection
      try {
        connection.requestMethod = "POST"
        connection.connectTimeout = 1_000
        connection.readTimeout = 2_000
        connection.setRequestProperty("Content-Type", "application/json")
        apiKey?.let { connection.setRequestProperty("Authorization", "Bearer $it") }
        connection.doOutput = true
        val bucket =
            JSONObject()
                .put("client", "aw-watcher-android-web")
                .put("type", "currenttab")
                .put("hostname", "android")
        connection.outputStream.use { output -> output.write(bucket.toString().toByteArray()) }
        connection.responseCode
      } finally {
        connection.disconnect()
      }
    } catch (_: IOException) {}
  }

  private fun post(event: JSONObject) {
    try {
      val connection =
          URL("$baseUrl/api/0/buckets/web.tab.current/heartbeat?pulsetime=1.0").openConnection()
              as HttpURLConnection
      try {
        connection.requestMethod = "POST"
        connection.connectTimeout = 1_000
        connection.readTimeout = 2_000
        connection.setRequestProperty("Content-Type", "application/json")
        apiKey?.let { connection.setRequestProperty("Authorization", "Bearer $it") }
        connection.doOutput = true
        connection.outputStream.use { output -> output.write(event.toString().toByteArray()) }
        connection.responseCode
      } finally {
        connection.disconnect()
      }
    } catch (_: IOException) {}
  }
}
