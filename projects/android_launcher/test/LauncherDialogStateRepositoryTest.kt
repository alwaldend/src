package com.alwaldend.src.projects.android_launcher

import androidx.datastore.core.DataStore
import java.io.ByteArrayInputStream
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.updateAndGet
import kotlinx.coroutines.runBlocking
import org.junit.Before
import org.junit.Test
import org.junit.runner.RunWith
import org.robolectric.RobolectricTestRunner

private class DataStoreMock(state: Model.State) : DataStore<Model.State> {

  private val _data = MutableStateFlow(state)
  override val data: Flow<Model.State> = _data.asStateFlow()

  override suspend fun updateData(transform: suspend (t: Model.State) -> Model.State): Model.State {
    return _data.updateAndGet { transform(it) }
  }
}

@RunWith(RobolectricTestRunner::class)
class LauncherDialogStateRepositoryTest {

  private lateinit var repoDefault: LauncherStateRepository
  private lateinit var repoEmpty: LauncherStateRepository
  private lateinit var repoWithApps: LauncherStateRepository

  @Before
  fun setUp() {
    repoDefault = LauncherStateRepository(dataStore = DataStoreMock(Defaults.State))
    repoEmpty = LauncherStateRepository(dataStore = DataStoreMock(Model.State.getDefaultInstance()))
    repoWithApps = LauncherStateRepository(dataStore = DataStoreMock(Defaults.State))
  }

  @Test
  fun stateOrThrow() = runBlocking {
    val state = repoDefault.stateOrThrow()
    assert(state == Defaults.State) {
      "state is not equal to the default state: $state, ${Defaults.State}"
    }
  }

  @Test
  fun observeState() = runBlocking {
    val state = repoDefault.observeState().first().getOrThrow()
    assert(state == Defaults.State) {
      "state is not equal to the default state: $state, ${Defaults.State}"
    }
  }

  @Test
  fun updateSettings() = runBlocking {
    val newSettings = repoEmpty.updateSettings { appCard = Defaults.AppCardSettings }
    assert(newSettings.appCard == Defaults.AppCardSettings) {
      "updated settings are not equal to default settings: $newSettings, ${Defaults.AppCardSettings}"
    }
  }

  @Test
  fun updateIsHomeApp() = runBlocking {
    assert(!repoEmpty.stateOrThrow().isHomeApp)
    repoEmpty.updateIsHomeApp(true)
    assert(repoEmpty.stateOrThrow().isHomeApp) { "updated isHomeApp should be true" }
  }

  @Test fun setIsHidden() {}

  @Test
  fun removeApp() = runBlocking {
    val app =
        Model.App.newBuilder().setPackageName("com.example.removed").setLabel("Removed").build()
    repoWithApps.reloadApp(app)
    repoWithApps.removeApp(app.packageName)
    assert(!repoWithApps.stateOrThrow().apps.appsMap.containsKey(app.packageName)) {
      "removed app should not remain in state: ${app.packageName}"
    }
  }

  @Test
  fun updateMissingApp() = runBlocking {
    val packageName = "com.example.newly.added"
    val app = Model.App.newBuilder().setPackageName(packageName).setLabel("Newly added").build()
    repoEmpty.reloadApp(app)
    assert(repoEmpty.stateOrThrow().apps.appsMap[packageName] == app) {
      "missing package-added app should be inserted: ${repoEmpty.stateOrThrow()}"
    }
  }

  @Test
  fun setShowHiddenApps() = runBlocking {
    assert(!repoEmpty.stateOrThrow().showHiddenApps)
    repoEmpty.setShowHiddenApps(true)
    assert(repoEmpty.stateOrThrow().showHiddenApps) { "updated showHiddenApps should be true" }
  }

  @Test
  fun reloadApp() = runBlocking {
    val app =
        Model.App.newBuilder().setPackageName("com.example.reload").setLabel("Reloaded").build()
    repoEmpty.reloadApp(app)
    assert(repoEmpty.stateOrThrow().apps.appsMap[app.packageName] == app) {
      "reloaded app should be stored: ${repoEmpty.stateOrThrow()}"
    }
  }

  @Test
  fun `reloadState preserves app metadata and hidden state`() = runBlocking {
    val installedApp =
        Model.App.newBuilder().setPackageName("com.example.installed").setLabel("Installed").build()
    val existingHiddenApp =
        Model.App.newBuilder()
            .setPackageName("com.example.hidden")
            .setLabel("Hidden")
            .setIsHidden(true)
            .build()
    val refreshedHiddenApp =
        Model.App.newBuilder()
            .setPackageName("com.example.hidden")
            .setLabel("Hidden refreshed")
            .setIsHidden(true)
            .build()
    repoEmpty.reloadApp(existingHiddenApp)
    repoEmpty.reloadState(
        newApps =
            Model.Apps.newBuilder()
                .putApps(installedApp.packageName, installedApp)
                .putApps(refreshedHiddenApp.packageName, refreshedHiddenApp)
                .build(),
        isHomeApp = true)
    val state = repoEmpty.stateOrThrow()
    assert(state.apps.appsMap[installedApp.packageName] == installedApp) {
      "installed app should be stored: $state"
    }
    assert(state.apps.appsMap[refreshedHiddenApp.packageName] == refreshedHiddenApp) {
      "existing hidden state should be preserved: $state"
    }
    assert(state.isHomeApp) { "home state should be updated: $state" }
  }

  @Test
  fun `reloadApp permits empty icon and shortcut metadata`() = runBlocking {
    val packageName = "com.example.label.only"
    val app = Model.App.newBuilder().setPackageName(packageName).setLabel("Label only").build()
    repoEmpty.reloadApp(app)
    assert(repoEmpty.stateOrThrow().apps.appsMap[packageName] == app) {
      "label-only app should be stored without icon or shortcuts: ${repoEmpty.stateOrThrow()}"
    }
  }

  @Test
  fun `reloadApp merges icon into an existing label-only entry`() = runBlocking {
    val packageName = "com.example.icon"
    val labelOnly = Model.App.newBuilder().setPackageName(packageName).setLabel("Icon").build()
    repoEmpty.reloadApp(labelOnly)
    val withIcon =
        labelOnly.toBuilder().setIcon(com.google.protobuf.ByteString.copyFrom(ICON_BYTES)).build()
    repoEmpty.reloadApp(withIcon)
    val stored = repoEmpty.stateOrThrow().apps.appsMap[packageName]
    assert(stored?.icon?.toByteArray()?.contentEquals(ICON_BYTES) == true) {
      "merged entry should carry the decoded icon bytes: $stored"
    }
    assert(stored?.label == "Icon") { "merged entry should keep its label: $stored" }
  }

  @Test
  fun resetSettings() = runBlocking {
    val settings = repoEmpty.resetSettings()
    assert(settings == Defaults.State.settings) {
      "reset settings are not equal to default settings: " + "$settings, ${Defaults.State.settings}"
    }
  }

  @Test
  fun updateSettingsFromInputStream() = runBlocking {
    val stream = ByteArrayInputStream(Defaults.Settings.toByteArray())
    repoEmpty.updateSettingsFromInputStream(stream)
    val settings = repoEmpty.stateOrThrow().settings
    assert(settings == Defaults.Settings) {
      "updated settings are not equal to default settings: " + "$settings, ${Defaults.Settings}"
    }
  }

  @Test fun reloadState() {}

  private companion object {
    // Smallest valid PNG payload; the dialog only needs decodable bytes.
    val ICON_BYTES =
        byteArrayOf(
            0x89.toByte(),
            0x50,
            0x4E,
            0x47,
            0x0D,
            0x0A,
            0x1A,
            0x0A,
            0x00,
            0x00,
            0x00,
            0x0D,
            0x49,
            0x48,
            0x44,
            0x52,
            0x00,
            0x00,
            0x00,
            0x01,
            0x00,
            0x00,
            0x00,
            0x01,
            0x08,
            0x06,
            0x00,
            0x00,
            0x00,
            0x1F,
            0x15,
            0xC4.toByte(),
            0x89.toByte(),
            0x00,
            0x00,
            0x00,
            0x0A,
            0x49,
            0x44,
            0x41,
            0x54,
            0x78,
            0x9C.toByte(),
            0x63,
            0x00,
            0x01,
            0x00,
            0x00,
            0x05,
            0x00,
            0x01,
            0x0D,
            0x0A,
            0x2D,
            0xB4.toByte(),
            0x00,
            0x00,
            0x00,
            0x00,
            0x49,
            0x45,
            0x4E,
            0x44,
            0xAE.toByte(),
            0x42,
            0x60,
            0x82.toByte(),
        )
  }
}
