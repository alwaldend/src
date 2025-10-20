package com.alwaldend.src.kt.android.launcher

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.lifecycle.lifecycleScope
import com.alwaldend.src.kt.android.launcher.ui.LauncherUi
import kotlinx.coroutines.launch

private val tag = MainActivity::class.simpleName


class MainActivity : ComponentActivity() {

    private lateinit var app: LauncherApplication
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        app = (application as LauncherApplication).init(this)

        lifecycleScope.launch {
            app.stateRepo.reloadState(
                newApps = app.manager.fetchAllApps(),
                isHomeApp = app.manager.isHomeApp()
            )
        }
        app.manager.addCallback(
            onChanged = { packageName ->
                lifecycleScope.launch {
                    app.stateRepo.reloadApp(app.manager.getApp(packageName))
                }
            },
            onRemoved = { packageName ->
                lifecycleScope.launch {
                    app.stateRepo.removeApp(packageName)
                }
            }
        )

        setContent {
            LauncherUi()
        }
    }

    override fun onResume() {
        super.onResume()
        lifecycleScope.launch {
            app.stateRepo.updateIsHomeApp(app.manager.isHomeApp())
        }
    }

    override fun onStop() {
        super.onStop()
        app.manager.removeCallbacks()
    }
}
