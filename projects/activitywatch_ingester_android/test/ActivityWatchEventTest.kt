package com.alwaldend.src.projects.activitywatch_ingester_android

import java.time.Duration
import java.time.Instant
import org.junit.Assert.assertEquals
import org.junit.Test
import org.junit.runner.RunWith
import org.robolectric.RobolectricTestRunner

@RunWith(RobolectricTestRunner::class)
class ActivityWatchEventTest {
  @Test
  fun `strips browser protocol`() {
    val event =
        ActivityWatchEvent(
            Instant.parse("2026-09-08T00:00:00Z"), "chrome", "https://example.com/page")

    assertEquals("example.com/page", event.toJson().getJSONObject("data").getString("url"))
  }

  @Test
  fun `serializes duration in seconds`() {
    val event =
        ActivityWatchEvent(
            Instant.parse("2026-09-08T00:00:00Z"),
            "chrome",
            "https://example.com/page",
            Duration.ofSeconds(5),
        )

    assertEquals(5.0, event.toJson().getDouble("duration"), 0.0)
  }

  @Test
  fun `tracks duration only while the same url stays active`() {
    val tracker = SessionTracker()
    val start = Instant.parse("2026-09-08T00:00:00Z")

    assertEquals(Duration.ZERO, tracker.durationFor("https://a.example", start))
    assertEquals(
        Duration.ofSeconds(3),
        tracker.durationFor("https://a.example", start.plusSeconds(3)),
    )
    assertEquals(Duration.ZERO, tracker.durationFor("https://b.example", start.plusSeconds(6)))
  }
}
