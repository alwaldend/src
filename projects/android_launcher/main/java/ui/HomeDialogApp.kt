package com.alwaldend.src.projects.android_launcher.ui

import androidx.compose.runtime.Composable
import androidx.compose.runtime.remember
import com.alwaldend.src.projects.android_launcher.Model
import com.alwaldend.src.projects.android_launcher.data.HomeDialogState

@Composable
fun HomeDialogApp(
    app: Model.App,
    showShortcuts: Boolean,
    onAppShortcutClick: (Model.AppShortcut) -> Unit,
    onHeaderAction: (Model.App, HomeDialogState.HeaderActions) -> Unit,
    onDismissRequest: () -> Unit
) {
    val buttons = remember(app) {
        HomeDialogState.Buttons(app.shortcutsList.map { shortcut ->
            HomeDialogState.Button(shortcut.label, shortcut)
        })
    }
    HomeDialog(
        onDismissRequest = onDismissRequest,
        actionButtons = buttons,
        app = app,
        onHeaderAction = onHeaderAction,
        showButtons = showShortcuts,
        onButtonClick = { onAppShortcutClick(it) }
    )
}
