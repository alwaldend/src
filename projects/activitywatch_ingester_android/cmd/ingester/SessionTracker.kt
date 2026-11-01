package com.alwaldend.src.projects.activitywatch_ingester_android

import java.time.Duration
import java.time.Instant

class SessionTracker {
  private var lastUrl: String? = null
  private var lastTimestamp: Instant? = null

  fun durationFor(url: String, now: Instant): Duration {
    val duration =
        if (lastUrl != null && lastTimestamp != null && url == lastUrl) {
          Duration.between(lastTimestamp, now)
        } else {
          Duration.ZERO
        }
    lastUrl = url
    lastTimestamp = now
    return duration
  }
}
