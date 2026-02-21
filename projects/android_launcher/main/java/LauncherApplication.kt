package com.alwaldend.src.projects.android_launcher

import android.app.Application
import android.content.Context
import android.content.pm.LauncherApps
import androidx.activity.ComponentActivity
import androidx.datastore.dataStore
import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.ViewModelProvider.AndroidViewModelFactory.Companion.APPLICATION_KEY
import androidx.lifecycle.viewmodel.initializer
import androidx.lifecycle.viewmodel.viewModelFactory

private val Context.dataStore by dataStore(
    fileName = "launcherApps.pb",
    serializer = LauncherStateSerializer,
)

class LauncherApplication : Application() {
    lateinit var manager: LauncherManager
    lateinit var stateRepo: LauncherStateRepository
    val navigator: LauncherNavigator = LauncherNavigatorImpl()

    fun init(activity: ComponentActivity): LauncherApplication {
        stateRepo = LauncherStateRepository(dataStore)
        manager = LauncherManager(
            context = activity,
            packageManager = activity.packageManager,
            launcherAppsService = getSystemService(LauncherApps::class.java),
            lifecycle = activity.lifecycle
        )
        return this
    }
}

inline fun <reified T : ViewModel> launcherViewModelFactory(
    crossinline init: LauncherApplication.() -> T
): ViewModelProvider.Factory {
    return viewModelFactory {
        initializer {
            (this[APPLICATION_KEY] as LauncherApplication).run(init)
        }
    }
}
