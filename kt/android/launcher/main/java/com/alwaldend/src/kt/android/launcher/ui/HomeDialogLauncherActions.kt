package com.alwaldend.src.kt.android.launcher.ui

import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import com.alwaldend.src.kt.android.launcher.R
import com.alwaldend.src.kt.android.launcher.data.HomeDialogState
import com.alwaldend.src.kt.android.launcher.data.HomeDialogState.Button
import com.alwaldend.src.kt.android.launcher.data.HomeDialogState.Buttons
import com.alwaldend.src.kt.android.launcher.data.HomeDialogState.LauncherDialogAction.GoToSettings
import com.alwaldend.src.kt.android.launcher.data.HomeDialogState.LauncherDialogAction.HideHiddenApps
import com.alwaldend.src.kt.android.launcher.data.HomeDialogState.LauncherDialogAction.ShowHiddenApps


private val buttons = Buttons(
    Button(
        label = R.string.launcher_dialog_settings,
        data = GoToSettings
    ),
    Button(
        label = R.string.launcher_dialog_hide_hidden_apps,
        data = HideHiddenApps
    ),
    Button(
        label = R.string.launcher_dialog_show_hidden_apps,
        data = ShowHiddenApps
    )
)

@Composable
fun HomeDialogLauncherActions(
    showHiddenApps: Boolean,
    onLauncherDialogAction: (HomeDialogState.LauncherDialogAction) -> Unit,
    onDismissRequest: () -> Unit,
    modifier: Modifier = Modifier
) {
    HomeDialog(
        onDismissRequest = onDismissRequest,
        modifier = modifier,
        actionButtons = buttons,
        showButton = { button ->
            when (button) {
                HideHiddenApps -> showHiddenApps
                ShowHiddenApps -> !showHiddenApps
                else -> true
            }
        },
        onButtonClick = onLauncherDialogAction
    )
}
