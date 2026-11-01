package com.alwaldend.src.projects.activitywatch_ingester_android

import java.time.Duration
import java.time.Instant
import org.json.JSONObject

data class ActivityWatchEvent(
    val timestamp: Instant,
    val browser: String,
    val url: String,
    val duration: Duration = Duration.ZERO,
) {
  fun toJson(): JSONObject =
      JSONObject()
          .put("timestamp", timestamp.toString())
          .put("duration", duration.toMillis() / 1000.0)
          .put("data", dataJson())

  private fun dataJson(): JSONObject =
      JSONObject()
          .put("audible", false)
          .put("browser", browser)
          .put("incognito", false)
          .put("title", "")
          .put("url", stripProtocol(url))
}

internal fun stripProtocol(url: String): String =
    url.removePrefix("http://").removePrefix("https://")
