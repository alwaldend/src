package com.alwaldend.src.projects.activitywatch_ingester_android

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.Button
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Switch
import androidx.compose.material3.Text
import androidx.compose.material3.TopAppBar
import androidx.compose.material3.TopAppBarDefaults
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.semantics.heading
import androidx.compose.ui.semantics.semantics
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp

@OptIn(ExperimentalMaterial3Api::class)
@Composable
internal fun IngesterSettingsScreen(
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
    modifier: Modifier = Modifier,
) {
  Scaffold(
      modifier = modifier,
      containerColor = MaterialTheme.colorScheme.background,
      topBar = {
        TopAppBar(
            title = { Text(stringResource(R.string.settings_title)) },
            colors =
                TopAppBarDefaults.topAppBarColors(
                    containerColor = MaterialTheme.colorScheme.surfaceContainer))
      },
      bottomBar = {
        Button(
            onClick = onOpenAccessibilitySettings,
            modifier = Modifier.fillMaxWidth().padding(16.dp),
        ) {
          Text(stringResource(R.string.settings_open_accessibility))
        }
      },
  ) { contentPadding ->
    Column(
        modifier =
            Modifier.padding(contentPadding)
                .verticalScroll(rememberScrollState())
                .padding(horizontal = 16.dp, vertical = 8.dp),
        verticalArrangement = Arrangement.spacedBy(12.dp),
    ) {
      Text(
          text = stringResource(R.string.settings_subtitle),
          modifier = Modifier.fillMaxWidth(),
          style = MaterialTheme.typography.bodyMedium,
          color = MaterialTheme.colorScheme.onSurfaceVariant,
          textAlign = TextAlign.Start)
      SettingsGroup(title = stringResource(R.string.settings_group_server)) {
        OutlinedTextField(
            value = baseUrl,
            onValueChange = onBaseUrlChange,
            label = { Text(stringResource(R.string.settings_base_url)) },
            supportingText = { Text(stringResource(R.string.settings_base_url_support)) },
            singleLine = true,
            modifier = Modifier.fillMaxWidth())
        OutlinedTextField(
            value = apiKey,
            onValueChange = onApiKeyChange,
            label = { Text(stringResource(R.string.settings_api_key)) },
            singleLine = true,
            modifier = Modifier.fillMaxWidth())
      }
      SettingsGroup(title = stringResource(R.string.settings_group_capture)) {
        SettingsSwitchRow(
            label = stringResource(R.string.settings_accessibility),
            supportingText = stringResource(R.string.settings_accessibility_support),
            checked = accessibilityEnabled,
            onCheckedChange = onAccessibilityChange)
        SettingsSwitchRow(
            label = stringResource(R.string.settings_url_extraction),
            supportingText = stringResource(R.string.settings_url_extraction_support),
            checked = urlExtractionEnabled,
            enabled = accessibilityEnabled,
            onCheckedChange = onUrlExtractionChange)
        SettingsSwitchRow(
            label = stringResource(R.string.settings_key_logging),
            supportingText = stringResource(R.string.settings_key_logging_support),
            checked = keyLoggingEnabled,
            enabled = accessibilityEnabled,
            onCheckedChange = onKeyLoggingChange)
      }
    }
  }
}

@Composable
private fun SettingsGroup(title: String, content: @Composable () -> Unit) {
  Card(
      modifier = Modifier.fillMaxWidth(),
      colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surfaceVariant),
  ) {
    Column(modifier = Modifier.padding(16.dp), verticalArrangement = Arrangement.spacedBy(8.dp)) {
      Text(
          text = title,
          modifier = Modifier.semantics { heading() },
          style = MaterialTheme.typography.titleMedium,
          color = MaterialTheme.colorScheme.onSurfaceVariant)
      content()
    }
  }
}

@Composable
private fun SettingsSwitchRow(
    label: String,
    supportingText: String,
    checked: Boolean,
    onCheckedChange: (Boolean) -> Unit,
    enabled: Boolean = true,
) {
  Row(
      modifier = Modifier.fillMaxWidth(),
      verticalAlignment = Alignment.CenterVertically,
      horizontalArrangement = Arrangement.spacedBy(12.dp),
  ) {
    Column(modifier = Modifier.weight(1f)) {
      Text(text = label, style = MaterialTheme.typography.bodyLarge)
      Text(
          text = supportingText,
          style = MaterialTheme.typography.bodySmall,
          color = MaterialTheme.colorScheme.onSurfaceVariant)
    }
    Switch(checked = checked, enabled = enabled, onCheckedChange = onCheckedChange)
  }
}
