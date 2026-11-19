package com.alwaldend.src.projects.activitywatch_ingester_android

import android.accessibilityservice.AccessibilityService
import android.content.Intent
import android.content.SharedPreferences
import android.view.accessibility.AccessibilityEvent
import java.time.Duration

private const val PREFS_NAME = "activitywatch_ingester"
private const val KEY_BASE_URL = "base_url"
private const val KEY_API_KEY = "api_key"
private const val KEY_ACCESSIBILITY_ENABLED = "accessibility_enabled"
private const val KEY_URL_EXTRACTION_ENABLED = "url_extraction_enabled"
private const val KEY_KEY_LOGGING_ENABLED = "key_logging_enabled"
private const val DEFAULT_BASE_URL = "http://127.0.0.1:5600"

class ActivityWatchAccessibilityService : AccessibilityService() {
  private var ingester = ActivityWatchIngester(DEFAULT_BASE_URL)
  private var lastUrl: String? = null
  private var lastSentAt: java.time.Instant? = null
  private var preferences: SharedPreferences? = null
  private var accessibilityEnabled = false
  private var urlExtractionEnabled = false
  private var keyLoggingEnabled = false
  private val preferenceListener =
      SharedPreferences.OnSharedPreferenceChangeListener { prefs, _ ->
        ingester = createIngester(prefs)
        readAccessibilityPreferences(prefs)
      }

  override fun onServiceConnected() {
    super.onServiceConnected()
    preferences = getSharedPreferences(PREFS_NAME, MODE_PRIVATE)
    ingester = createIngester(preferences!!)
    readAccessibilityPreferences(preferences!!)
    preferences!!.registerOnSharedPreferenceChangeListener(preferenceListener)
  }

  override fun onUnbind(intent: Intent?): Boolean {
    preferences?.unregisterOnSharedPreferenceChangeListener(preferenceListener)
    preferences = null
    return super.onUnbind(intent)
  }

  private fun createIngester(preferences: SharedPreferences): ActivityWatchIngester {
    val baseUrl = preferences.getString(KEY_BASE_URL, DEFAULT_BASE_URL) ?: DEFAULT_BASE_URL
    val apiKey = preferences.getString(KEY_API_KEY, "").orEmpty().takeIf(String::isNotBlank)
    return ActivityWatchIngester(baseUrl, apiKey)
  }

  private fun readAccessibilityPreferences(preferences: SharedPreferences) {
    accessibilityEnabled =
        preferences.getBoolean(KEY_ACCESSIBILITY_ENABLED, false) ||
            preferences.getBoolean(KEY_URL_EXTRACTION_ENABLED, false) ||
            preferences.getBoolean(KEY_KEY_LOGGING_ENABLED, false)
    urlExtractionEnabled = preferences.getBoolean(KEY_URL_EXTRACTION_ENABLED, false)
    keyLoggingEnabled = preferences.getBoolean(KEY_KEY_LOGGING_ENABLED, false)
  }

  override fun onAccessibilityEvent(event: AccessibilityEvent) {
    if (!accessibilityEnabled) {
      return
    }
    val packageName = event.packageName?.toString() ?: return
    val source = event.source ?: return
    try {
      val url = extractUrl(source, packageName) ?: return
      val now = java.time.Instant.now()
      if (url != lastUrl) {
        lastUrl?.let { previous ->
          val elapsed = lastSentAt?.let { Duration.between(it, now) } ?: Duration.ZERO
          if (elapsed > Duration.ZERO) {
            ingester.send(packageName, previous, elapsed)
          }
        }
        lastUrl = url
        lastSentAt = now
      }
      ingester.send(packageName, url, Duration.ZERO)
    } finally {
      source.recycle()
    }
  }

  override fun onInterrupt() {}
}
