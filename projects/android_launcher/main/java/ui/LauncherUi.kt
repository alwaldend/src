package com.alwaldend.src.projects.android_launcher.ui

import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.Surface
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.rememberNavController

private val startDestination: LauncherRoute<*> = HomeRoute

object Routes {
    private val home = HomeRoute
    val settings = SettingsRoute
    val allRoutes = listOf(home, settings)
}

@Composable
fun LauncherUi(modifier: Modifier = Modifier) {
    LauncherTheme {
        Surface(modifier = modifier.fillMaxSize()) {
            val navController = rememberNavController()
            NavHost(
                modifier = modifier,
                navController = navController,
                startDestination = startDestination.url
            ) {
                Routes.allRoutes.forEach { route ->
                    launcherComposable(
                        route = route,
                        navController = navController
                    )
                }
            }
        }
    }
}
