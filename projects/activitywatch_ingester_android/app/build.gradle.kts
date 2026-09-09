plugins {
  id("com.android.application")
  id("org.jetbrains.kotlin.android")
  id("org.jetbrains.kotlin.plugin.compose")
}

android {
  namespace = "com.alwaldend.src.projects.activitywatch_ingester_android"
  compileSdk = 36

  defaultConfig {
    applicationId = namespace
    minSdk = 25
    targetSdk = 36
    versionCode = 1
    versionName = "0.0.0-dev"
  }

  compileOptions {
    sourceCompatibility = JavaVersion.VERSION_11
    targetCompatibility = JavaVersion.VERSION_11
  }

  kotlinOptions { jvmTarget = "11" }

  buildFeatures { compose = true }

  lint { checkReleaseBuilds = false }
}

// Configure source directories to point to the Bazel-layout sources
android.sourceSets.getByName("main") {
  java.srcDirs("../main/java")
  res.srcDirs("../res")
  manifest.srcFile("../AndroidManifest.xml")
}

android.sourceSets.getByName("androidTest") { java.srcDirs("../androidTest") }

tasks.register<Copy>("copyReleaseApkToFdroidOutput") {
  from(layout.buildDirectory.dir("outputs/apk/release")) { include("*.apk") }
  into(rootProject.layout.projectDirectory.dir("build/outputs/apk/release"))
}

tasks
    .matching { it.name == "packageRelease" }
    .configureEach { finalizedBy("copyReleaseApkToFdroidOutput") }

dependencies {
  implementation(libs.androidx.activity.activity.compose)
  implementation(libs.androidx.activity.activity)
  implementation(libs.androidx.compose.foundation.foundation)
  implementation(libs.androidx.compose.foundation.foundation.layout)
  implementation(libs.androidx.compose.material3.material3)
  implementation(libs.androidx.compose.runtime.runtime)
  implementation(libs.androidx.compose.ui.ui)
  implementation(libs.androidx.compose.ui.ui.graphics)
  implementation(libs.androidx.compose.ui.ui.tooling)
  implementation(libs.androidx.core.core)
  implementation(libs.androidx.core.core.ktx)
  implementation(libs.org.json.json)
  implementation(libs.org.jetbrains.kotlin.kotlin.stdlib)
  implementation(libs.com.google.guava.guava)

  testImplementation(libs.junit.junit)
  testImplementation(libs.org.robolectric.robolectric)
  androidTestImplementation(libs.androidx.test.ext.junit.ktx)
  androidTestImplementation(libs.androidx.test.runner)
  androidTestImplementation(libs.androidx.test.espresso.espresso.core)
}
