package com.alwaldend.src.projects.activitywatch_ingester_android

import android.view.accessibility.AccessibilityNodeInfo

private val URL_BAR_VIEW_IDS =
    listOf("url_bar", "url_bar_editor", "addressbar_edit", "address_bar_edit_text", "edit_text_url")

internal fun extractUrl(source: AccessibilityNodeInfo, packageName: String): String? =
    URL_BAR_VIEW_IDS.asSequence()
        .flatMap { viewId ->
          source
              .findAccessibilityNodeInfosByViewId("$packageName:id/$viewId")
              .orEmpty()
              .asSequence()
        }
        .mapNotNull { it.text?.toString() }
        .map(String::trim)
        .firstOrNull(String::isNotEmpty)
