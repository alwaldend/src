package com.alwaldend.src.projects.android_launcher.ui

import androidx.compose.foundation.ExperimentalFoundationApi
import androidx.compose.foundation.combinedClickable
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.safeDrawingPadding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Warning
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.Card
import androidx.compose.material3.Icon
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.dimensionResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.tooling.preview.Preview
import androidx.lifecycle.viewmodel.compose.viewModel
import com.alwaldend.src.projects.android_launcher.Defaults
import com.alwaldend.src.projects.android_launcher.Model
import com.alwaldend.src.projects.android_launcher.R
import com.alwaldend.src.projects.android_launcher.data.HomeDialogState
import com.alwaldend.src.projects.android_launcher.data.HomeDialogState.HeaderActions
import com.alwaldend.src.projects.android_launcher.data.UiState

object HomeRoute : LauncherRoute<HomeViewModel> {
    override val url = "home_screen"
    override val label = R.string.navigation_home
    override val showScaffold = false

    @Composable
    override fun getViewModel(): HomeViewModel {
        return viewModel(factory = HomeViewModel.Factory)
    }

    @Composable
    override fun Content(
        vm: HomeViewModel,
        uiState: UiState.Success<Model.State>
    ) {
        HomeScreen(
            state = uiState.state,
            onHeaderAction = vm::onHeaderAction,
            onAppDialogClick = vm::onAppDialogClick,
            onRowItemClick = vm::onRowItemClick,
            onLauncherDialogAction = vm::onLauncherDialogAction,
            getRowItemLabel = vm::getRowItemLabel,
            sortRowItems = vm::sortRowItems
        )
    }
}


@OptIn(ExperimentalFoundationApi::class)
@Composable
private fun HomeScreen(
    state: Model.State,
    onAppDialogClick: (Model.AppShortcut) -> Unit,
    onLauncherDialogAction: (HomeDialogState.LauncherDialogAction) -> Unit,
    onHeaderAction: (Model.App, HeaderActions) -> Unit,
    onRowItemClick: (Model.App) -> Unit,
    sortRowItems: (Model.Apps, Model.Settings.Layout) -> List<Model.App>,
    getRowItemLabel: (Model.App, Model.Settings.AppCard) -> String,
    modifier: Modifier = Modifier
) {
    var showLauncherDialog by remember { mutableStateOf(false) }
    var showAppDialog by remember { mutableStateOf<Model.App?>(null) }
    val interactionSource = remember { MutableInteractionSource() }

    Column(modifier = modifier.safeDrawingPadding()) {
        if (!state.isHomeApp) {
            NotAHomeAppBanner()
        }
        HomeRow(
            modifier = Modifier
                .fillMaxSize()
                .combinedClickable(
                    interactionSource = interactionSource,
                    indication = null,
                    onLongClick = { showLauncherDialog = true },
                    onClick = { }
                ),
            settings = state.settings,
            apps = state.apps,
            showHiddenApps = state.showHiddenApps,
            onClick = onRowItemClick,
            onLongClick = { showAppDialog = it },
            getItemLabel = getRowItemLabel,
            sortItems = sortRowItems
        )
    }

    if (showLauncherDialog) {
        HomeDialogLauncherActions(
            showHiddenApps = state.showHiddenApps,
            onLauncherDialogAction = onLauncherDialogAction,
            onDismissRequest = { showLauncherDialog = false }
        )
    }
    showAppDialog?.let {
        HomeDialogApp(
            app = it,
            showShortcuts = state.isHomeApp,
            onAppShortcutClick = onAppDialogClick,
            onHeaderAction = onHeaderAction,
            onDismissRequest = { showAppDialog = null }
        )
    }
}


@Composable
private fun NotAHomeAppBanner() {
    Card {
        Row(
            modifier = Modifier
                .padding(dimensionResource(R.dimen.padding_medium))
                .fillMaxWidth(),
            horizontalArrangement = Arrangement.Center,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Icon(
                imageVector = Icons.Filled.Warning,
                contentDescription = "Warning icon",
                modifier = Modifier.size(ButtonDefaults.IconSize * 2),
            )
            Spacer(modifier = Modifier.width(dimensionResource(R.dimen.padding_medium)))
            Text(stringResource(R.string.not_a_home_app_banner))
        }
    }
}


@Preview(
    name = "phone",
    device = "spec:shape=Normal,width=360,height=640,unit=dp,dpi=480"
)
@Preview(
    name = "landscape",
    device = "spec:shape=Normal,width=640,height=360,unit=dp,dpi=480"
)
@Preview(
    name = "foldable",
    device = "spec:shape=Normal,width=673,height=841,unit=dp,dpi=480"
)
@Preview(
    name = "tablet",
    device = "spec:shape=Normal,width=1280,height=800,unit=dp,dpi=480"
)
@Composable
private fun HomePagePreview() {
    val state = Defaults.State
    LauncherTheme(darkTheme = true) {
        Box {
            HomeScreen(
                modifier = Modifier,
                state = state,
                onLauncherDialogAction = { },
                onAppDialogClick = {},
                onHeaderAction = { _, _ -> },
                onRowItemClick = { _ -> },
                sortRowItems = { _, _ -> listOf() },
                getRowItemLabel = { app, _ -> app.label }
            )
        }
    }
}
