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

rootProject.name = "activitywatch_ingester_android"

include(":app")
