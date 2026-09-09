package com.alwaldend.src.projects.activitywatch_ingester_android

import androidx.test.ext.junit.runners.AndroidJUnit4
import androidx.test.platform.app.InstrumentationRegistry
import org.junit.Assert.assertEquals
import org.junit.Test
import org.junit.runner.RunWith

@RunWith(AndroidJUnit4::class)
class ActivityWatchInstrumentationTest {
  @Test
  fun appContextMatchesPackage() {
    val appContext = InstrumentationRegistry.getInstrumentation().targetContext
    assertEquals(
        "com.alwaldend.src.projects.activitywatch_ingester_android", appContext.packageName)
  }

  @Test
  fun accessibilityServiceIsDeclared() {
    val appContext = InstrumentationRegistry.getInstrumentation().targetContext
    val serviceInfo =
        appContext.packageManager.getServiceInfo(
            android.content.ComponentName(
                appContext, ActivityWatchAccessibilityService::class.java),
            0,
        )
    assertEquals(
        "android.permission.BIND_ACCESSIBILITY_SERVICE",
        serviceInfo.permission,
    )
  }
}
