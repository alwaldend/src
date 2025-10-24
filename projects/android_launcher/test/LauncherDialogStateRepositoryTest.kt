package com.alwaldend.src.kt.android.launcher

import androidx.datastore.core.DataStore
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.updateAndGet
import kotlinx.coroutines.runBlocking
import java.io.ByteArrayInputStream
import org.robolectric.RobolectricTestRunner
import org.junit.Test
import org.junit.Before
import org.junit.runner.RunWith

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
        repoDefault = LauncherStateRepository(
            dataStore = DataStoreMock(Defaults.State)
        )
        repoEmpty = LauncherStateRepository(
            dataStore = DataStoreMock(Model.State.getDefaultInstance())
        )
        repoWithApps = LauncherStateRepository(
            dataStore = DataStoreMock(
                Defaults.State
            )
        )
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
        val newSettings = repoEmpty.updateSettings {
            appCard = Defaults.AppCardSettings
        }
        assert(newSettings.appCard == Defaults.AppCardSettings) {
            "updated settings are not equal to default settings: $newSettings, ${Defaults.AppCardSettings}"
        }
    }

    @Test
    fun updateIsHomeApp() = runBlocking {
        assert(!repoEmpty.stateOrThrow().isHomeApp)
        repoEmpty.updateIsHomeApp(true)
        assert(repoEmpty.stateOrThrow().isHomeApp) {
            "updated isHomeApp should be true"
        }
    }

    @Test
    fun setIsHidden() {
    }

    @Test
    fun removeApp() = runBlocking {
    }

    @Test
    fun setShowHiddenApps() = runBlocking {
        assert(!repoEmpty.stateOrThrow().showHiddenApps)
        repoEmpty.setShowHiddenApps(true)
        assert(repoEmpty.stateOrThrow().showHiddenApps) {
            "updated showHiddenApps should be true"
        }
    }

    @Test
    fun reloadApp() {

    }

    @Test
    fun resetSettings() = runBlocking {
        val settings = repoEmpty.resetSettings()
        assert(settings == Defaults.State.settings) {
            "reset settings are not equal to default settings: " +
                    "$settings, ${Defaults.State.settings}"
        }
    }

    @Test
    fun updateSettingsFromInputStream() = runBlocking {
        val stream = ByteArrayInputStream(Defaults.Settings.toByteArray())
        repoEmpty.updateSettingsFromInputStream(stream)
        val settings = repoEmpty.stateOrThrow().settings
        assert(settings == Defaults.Settings) {
            "updated settings are not equal to default settings: " +
                    "$settings, ${Defaults.Settings}"
        }
    }

    @Test
    fun reloadState() {
    }
}
