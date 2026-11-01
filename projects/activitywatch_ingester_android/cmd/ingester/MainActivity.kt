package com.alwaldend.src.projects.activitywatch_ingester_android

import android.content.Context
import android.content.Intent
import android.content.res.Configuration
import android.os.Bundle
import android.provider.Settings
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.core.view.WindowCompat

private const val PREFS_NAME = "activitywatch_ingester"
private const val KEY_BASE_URL = "base_url"
private const val KEY_API_KEY = "api_key"
private const val KEY_ACCESSIBILITY_ENABLED = "accessibility_enabled"
private const val KEY_URL_EXTRACTION_ENABLED = "url_extraction_enabled"
private const val KEY_KEY_LOGGING_ENABLED = "key_logging_enabled"

class MainActivity : ComponentActivity() {
  override fun onCreate(savedInstanceState: Bundle?) {
    super.onCreate(savedInstanceState)
    WindowCompat.setDecorFitsSystemWindows(window, false)
    updateStatusBarContrast()
    val prefs = getSharedPreferences(PREFS_NAME, Context.MODE_PRIVATE)
    val initialBaseUrl = prefs.getString(KEY_BASE_URL, DEFAULT_BASE_URL) ?: DEFAULT_BASE_URL
    val initialApiKey = prefs.getString(KEY_API_KEY, "") ?: ""
    val initialAccessibilityEnabled = prefs.getBoolean(KEY_ACCESSIBILITY_ENABLED, false)
    val initialUrlExtractionEnabled = prefs.getBoolean(KEY_URL_EXTRACTION_ENABLED, false)
    val initialKeyLoggingEnabled = prefs.getBoolean(KEY_KEY_LOGGING_ENABLED, false)
    setContent {
      var baseUrlState by remember { mutableStateOf(initialBaseUrl) }
      var apiKeyState by remember { mutableStateOf(initialApiKey) }
      var accessibilityEnabled by remember { mutableStateOf(initialAccessibilityEnabled) }
      var urlExtractionEnabled by remember { mutableStateOf(initialUrlExtractionEnabled) }
      var keyLoggingEnabled by remember { mutableStateOf(initialKeyLoggingEnabled) }
      MainContent(
          baseUrl = baseUrlState,
          apiKey = apiKeyState,
          onBaseUrlChange = { baseUrl ->
            baseUrlState = baseUrl
            prefs.edit().putString(KEY_BASE_URL, baseUrl).apply()
          },
          onApiKeyChange = { apiKey ->
            apiKeyState = apiKey
            prefs.edit().putString(KEY_API_KEY, apiKey).apply()
          },
          accessibilityEnabled = accessibilityEnabled,
          urlExtractionEnabled = urlExtractionEnabled,
          keyLoggingEnabled = keyLoggingEnabled,
          onAccessibilityChange = { enabled ->
            accessibilityEnabled = enabled
            prefs.edit().putBoolean(KEY_ACCESSIBILITY_ENABLED, enabled).apply()
          },
          onUrlExtractionChange = { enabled ->
            urlExtractionEnabled = enabled
            prefs.edit().putBoolean(KEY_URL_EXTRACTION_ENABLED, enabled).apply()
          },
          onKeyLoggingChange = { enabled ->
            keyLoggingEnabled = enabled
            prefs.edit().putBoolean(KEY_KEY_LOGGING_ENABLED, enabled).apply()
          },
          onOpenAccessibilitySettings = {
            startActivity(Intent(Settings.ACTION_ACCESSIBILITY_SETTINGS))
          },
      )
    }
  }

  override fun onResume() {
    super.onResume()
    updateStatusBarContrast()
  }

  /**
   * Keeps the status-bar region inside the themed surface and flips the clock and battery icons to
   * match the active light or dark appearance.
   */
  private fun updateStatusBarContrast() {
    val darkAppearance =
        (resources.configuration.uiMode and Configuration.UI_MODE_NIGHT_MASK) ==
            Configuration.UI_MODE_NIGHT_YES
    WindowCompat.getInsetsController(window, window.decorView).apply {
      isAppearanceLightStatusBars = !darkAppearance
      isAppearanceLightNavigationBars = !darkAppearance
    }
  }
}

private const val DEFAULT_BASE_URL = "http://127.0.0.1:5600"

@Composable
private fun MainContent(
    baseUrl: String,
    apiKey: String,
    accessibilityEnabled: Boolean,
    urlExtractionEnabled: Boolean,
    keyLoggingEnabled: Boolean,
    onBaseUrlChange: (String) -> Unit,
    onApiKeyChange: (String) -> Unit,
    onAccessibilityChange: (Boolean) -> Unit,
    onUrlExtractionChange: (Boolean) -> Unit,
    onKeyLoggingChange: (Boolean) -> Unit,
    onOpenAccessibilitySettings: () -> Unit,
) {
  IngesterTheme {
    Surface(
        modifier = Modifier.fillMaxSize(),
        color = MaterialTheme.colorScheme.background,
    ) {
      IngesterSettingsScreen(
          baseUrl = baseUrl,
          apiKey = apiKey,
          accessibilityEnabled = accessibilityEnabled,
          urlExtractionEnabled = urlExtractionEnabled,
          keyLoggingEnabled = keyLoggingEnabled,
          onBaseUrlChange = onBaseUrlChange,
          onApiKeyChange = onApiKeyChange,
          onAccessibilityChange = onAccessibilityChange,
          onUrlExtractionChange = onUrlExtractionChange,
          onKeyLoggingChange = onKeyLoggingChange,
          onOpenAccessibilitySettings = onOpenAccessibilitySettings)
    }
  }
}
