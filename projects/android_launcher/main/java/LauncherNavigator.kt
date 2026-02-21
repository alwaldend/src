package com.alwaldend.src.projects.android_launcher

import androidx.compose.runtime.Immutable
import com.alwaldend.src.projects.android_launcher.ui.LauncherRoute
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update

@Immutable
sealed interface LauncherNavigationState {

    data object Idle : LauncherNavigationState

    data class GoToRoute(val route: LauncherRoute<*>) : LauncherNavigationState

    data object PopBackStack : LauncherNavigationState
}

interface LauncherNavigator {
    val navigationStateFlow: StateFlow<LauncherNavigationState>

    fun onNavigation(state: LauncherNavigationState)

    fun popBackStack()

    fun goToRoute(route: LauncherRoute<*>)
}

class LauncherNavigatorImpl : LauncherNavigator {
    private val navigationState = MutableStateFlow<LauncherNavigationState>(
        LauncherNavigationState.Idle
    )

    override val navigationStateFlow = navigationState.asStateFlow()

    override fun onNavigation(state: LauncherNavigationState) {
        navigationState.compareAndSet(state, LauncherNavigationState.Idle)
    }

    override fun popBackStack() {
        navigate(LauncherNavigationState.PopBackStack)
    }

    override fun goToRoute(route: LauncherRoute<*>) {
        navigate(LauncherNavigationState.GoToRoute(route))
    }

    private fun navigate(state: LauncherNavigationState) {
        navigationState.update { state }
    }
}
