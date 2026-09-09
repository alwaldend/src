pluginManagement {
  repositories {
    google()
    mavenCentral()
    gradlePluginPortal()
  }
}

dependencyResolutionManagement {
  repositoriesMode = RepositoriesMode.FAIL_ON_PROJECT_REPOS
  repositories {
    google()
    mavenCentral()
  }
  versionCatalogs {
    create("androidPlugins") { from(files("../../tools/android/gradle/libs.versions.toml")) }
  }
}

rootProject.name = "android_launcher"

include(":app")
