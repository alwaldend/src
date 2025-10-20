package com.alwaldend.src.kt.android.launcher.ui

import android.graphics.BitmapFactory
import androidx.compose.foundation.Image
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.ColumnScope
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.ElevatedCard
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.key
import androidx.compose.runtime.remember
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.asImageBitmap
import androidx.compose.ui.res.dimensionResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.window.Dialog
import com.alwaldend.src.kt.android.launcher.Model
import com.alwaldend.src.kt.android.launcher.R
import com.alwaldend.src.kt.android.launcher.data.HomeDialogState
import com.alwaldend.src.kt.android.launcher.data.HomeDialogState.Button.Label

@Composable
fun <T> HomeDialog(
    onDismissRequest: () -> Unit,
    onButtonClick: (T) -> Unit,
    actionButtons: HomeDialogState.Buttons<T>,
    modifier: Modifier = Modifier,
    showButton: (T) -> Boolean = { true },
) {
    HomeDialogBase(
        modifier = modifier,
        onDismissRequest = onDismissRequest
    ) {
        HomeDialogButtons(
            onDismissRequest = onDismissRequest,
            onButtonClick = onButtonClick,
            showButton = showButton,
            actionButtons = actionButtons,
        )
    }
}

@Composable
fun <T> HomeDialog(
    onDismissRequest: () -> Unit,
    onButtonClick: (T) -> Unit,
    actionButtons: HomeDialogState.Buttons<T>,
    app: Model.App,
    onHeaderAction: (Model.App, HomeDialogState.HeaderActions) -> Unit,
    modifier: Modifier = Modifier,
    showButton: (T) -> Boolean = { true },
    showButtons: Boolean = true,
) {
    HomeDialogBase(
        modifier = modifier,
        onDismissRequest = onDismissRequest
    ) {
        HomeDialogHeader(
            app = app,
            onHeaderAction = onHeaderAction,
            onDismissRequest = onDismissRequest
        )
        Spacer(modifier = Modifier.height(dimensionResource(id = R.dimen.padding_medium)))
        if (actionButtons.buttons.isNotEmpty() && showButtons) {
            HomeDialogButtons(
                onDismissRequest = onDismissRequest,
                onButtonClick = onButtonClick,
                showButton = showButton,
                actionButtons = actionButtons,
            )
        }
    }
}

@Composable
fun HomeDialogButton(
    text: String,
    onClick: () -> Unit,
    modifier: Modifier = Modifier,
    textAlign: TextAlign = TextAlign.Center
) {
    Surface(
        modifier = modifier.fillMaxWidth(),
        onClick = onClick,
    ) {
        Text(
            text = text,
            modifier = Modifier.padding(dimensionResource(R.dimen.padding_medium)),
            style = MaterialTheme.typography.bodyLarge,
            textAlign = textAlign
        )
    }
}

@Composable
private fun <T> HomeDialogButtons(
    onDismissRequest: () -> Unit,
    onButtonClick: (T) -> Unit,
    showButton: (T) -> Boolean,
    actionButtons: HomeDialogState.Buttons<T>,
) {
    ElevatedCard {
        Column {
            actionButtons.buttons.forEach { button ->
                key(button.hashCode()) {
                    if (showButton(button.data)) {
                        HomeDialogButton(
                            text = when (button.label) {
                                is Label.Res -> stringResource(button.label.res)
                                is Label.Str -> button.label.string
                            },
                            textAlign = TextAlign.Start,
                            onClick = {
                                onButtonClick(button.data)
                                onDismissRequest()
                            }
                        )
                    }
                }
            }
        }
    }
}

@Composable
private fun HomeDialogBase(
    onDismissRequest: () -> Unit,
    modifier: Modifier = Modifier,
    content: @Composable ColumnScope.() -> Unit
) {
    Dialog(onDismissRequest = onDismissRequest) {
        Column(
            modifier = modifier.verticalScroll(rememberScrollState()),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.Center,
            content = content
        )
    }
}

@Composable
private fun HomeDialogHeader(
    app: Model.App,
    onHeaderAction: (Model.App, HomeDialogState.HeaderActions) -> Unit,
    onDismissRequest: () -> Unit,
    modifier: Modifier = Modifier
) {
    val icon = remember {
        val array = app.icon.toByteArray()
        BitmapFactory.decodeByteArray(array, 0, array.size).asImageBitmap()
    }
    Row(
        modifier = modifier.fillMaxWidth(),
        horizontalArrangement = Arrangement.SpaceBetween,
        verticalAlignment = Alignment.CenterVertically
    ) {
        ElevatedCard(modifier = Modifier.weight(1f)) {
            Column(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(0.dp, dimensionResource(R.dimen.padding_medium)),
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.Center
            ) {
                Image(
                    modifier = Modifier
                        .padding(dimensionResource(R.dimen.padding_small))
                        .size(dimensionResource(R.dimen.image_size)),
                    bitmap = icon,
                    contentDescription = "App icon",
                )
                Text(
                    text = app.label,
                    modifier = Modifier.padding(dimensionResource(R.dimen.padding_small)),
                    style = MaterialTheme.typography.bodyLarge,
                    textAlign = TextAlign.Center
                )
            }
        }
        Spacer(modifier = Modifier.width(dimensionResource(R.dimen.padding_medium)))
        ElevatedCard(modifier = Modifier.weight(1f)) {
            Column(
                modifier = Modifier.fillMaxWidth(),
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.Center
            ) {
                HomeDialogButton(
                    text = stringResource(R.string.drawer_app_info),
                    onClick = {
                        onHeaderAction(
                            app,
                            HomeDialogState.HeaderActions.GoToInfo
                        )
                        onDismissRequest()
                    }
                )
                HomeDialogButton(
                    text = stringResource(
                        if (app.isHidden) R.string.drawer_app_show else R.string.drawer_app_hide
                    ),
                    onClick = {
                        onHeaderAction(
                            app,
                            HomeDialogState.HeaderActions.HideOrShow
                        )
                        onDismissRequest()
                    }
                )
                HomeDialogButton(
                    text = stringResource(R.string.drawer_app_uninstall),
                    onClick = {
                        onHeaderAction(
                            app,
                            HomeDialogState.HeaderActions.Uninstall
                        )
                        onDismissRequest()
                    }
                )
            }
        }
    }
}

@Preview
@Composable
private fun HomeDialogButtonPreview() {
    LauncherTheme {
        HomeDialogButton(text = "button", onClick = { })
    }
}
