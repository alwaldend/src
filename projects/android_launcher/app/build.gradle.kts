import com.google.protobuf.gradle.id
import com.google.protobuf.gradle.proto

plugins {
  id("com.android.application")
  id("org.jetbrains.kotlin.android")
  id("org.jetbrains.kotlin.plugin.compose")
  id("com.google.protobuf")
}

android {
  namespace = "com.alwaldend.src.projects.android_launcher"
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
  proto { srcDir("../main/proto") }
  manifest.srcFile("../AndroidManifest.xml")
}

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
  implementation(libs.androidx.compose.foundation.foundation.layout)
  implementation(libs.androidx.compose.material3.material3)
  implementation(libs.androidx.compose.runtime.runtime)
  implementation(libs.androidx.compose.ui.ui)
  implementation(libs.androidx.compose.ui.ui.graphics)
  implementation(libs.androidx.compose.ui.ui.tooling)
  implementation(libs.androidx.core.core.ktx)
  implementation(libs.androidx.datastore.datastore)
  implementation(libs.androidx.lifecycle.lifecycle.viewmodel.compose)
  implementation(libs.androidx.lifecycle.lifecycle.runtime.compose)
  implementation(libs.androidx.navigation.navigation.compose)
  implementation(libs.com.google.guava.guava)
  implementation(libs.org.jetbrains.kotlin.kotlin.stdlib)

  implementation("com.google.protobuf:protobuf-javalite:4.33.0")

  testImplementation(libs.junit.junit)
  testImplementation(libs.org.robolectric.robolectric)
  testImplementation(libs.androidx.test.ext.junit.ktx)
  testImplementation(libs.androidx.test.espresso.espresso.core)
}

protobuf {
  protoc { artifact = "com.google.protobuf:protoc:4.33.0" }
  generateProtoTasks {
    all().forEach { task -> task.builtins { create("java") { option("lite") } } }
  }
}
