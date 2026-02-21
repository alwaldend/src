package com.alwaldend.src.projects.android_launcher.ui

import androidx.annotation.StringRes
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.ColumnScope
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.RowScope
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.heightIn
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.LazyItemScope
import androidx.compose.material3.Checkbox
import androidx.compose.material3.ElevatedCard
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.RadioButton
import androidx.compose.material3.Surface
import androidx.compose.material3.Switch
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableIntStateOf
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalConfiguration
import androidx.compose.ui.res.dimensionResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.window.Dialog
import com.alwaldend.src.projects.android_launcher.R

private fun <T : Enum<*>> Class<T>.names() = enumConstants!!.mapNotNull {
    if (it.name == "UNRECOGNIZED") null else it.name
}

private fun <T : Enum<*>> Class<T>.find(name: String) = enumConstants!!.first {
    it.name == name
}

@Composable
fun SettingsGroup(
    @StringRes name: Int,
    content: @Composable ColumnScope.() -> Unit
) {
    Column(modifier = Modifier.padding(dimensionResource(R.dimen.padding_medium))) {
        Text(
            text = stringResource(name),
            modifier = Modifier.padding(dimensionResource(R.dimen.padding_medium))
        )
        Surface(modifier = Modifier.fillMaxWidth()) {
            Column(content = content)
        }
    }
}

@Composable
fun SettingsField(
    @StringRes name: Int,
    onClick: () -> Unit,
    modifier: Modifier = Modifier,
    label: String? = null,
    content: @Composable RowScope.() -> Unit = {}
) {
    Surface(modifier = modifier, onClick = onClick) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(dimensionResource(R.dimen.padding_medium)),
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.SpaceBetween
        ) {
            Column {
                Text(
                    text = stringResource(id = name),
                    textAlign = TextAlign.Start,
                    style = MaterialTheme.typography.headlineSmall,
                )
                if (label != null) {
                    Text(
                        text = label,
                        textAlign = TextAlign.Start,
                        style = MaterialTheme.typography.labelLarge,
                    )
                }
            }
            content()
        }
    }
}

@Composable
fun SettingsFieldSwitch(
    @StringRes name: Int,
    @StringRes label: Int,
    isToggled: Boolean,
    modifier: Modifier = Modifier,
    onClick: () -> Unit
) {
    SettingsField(
        name = name, label = stringResource(label),
        modifier = modifier, onClick = onClick
    ) {
        Switch(checked = isToggled, onCheckedChange = null)
    }
}


@Composable
fun SettingsFieldSingleChoice(
    @StringRes name: Int,
    options: List<String>,
    selectedOption: Int,
    modifier: Modifier = Modifier,
    onConfirm: (Int) -> Unit,
) {
    var showDialog by remember { mutableStateOf(false) }
    SettingsField(
        name = name,
        modifier = modifier,
        label = options[selectedOption],
        onClick = { showDialog = true }
    )
    if (!showDialog) {
        return
    }
    var curChoice by remember { mutableIntStateOf(selectedOption) }
    SettingsDialog(
        modifier = modifier,
        name = name,
        onConfirm = { onConfirm(curChoice); showDialog = false },
        onDismissRequest = { showDialog = false },
        itemsCount = options.size,
        itemsKey = { i -> options[i] }
    ) { i ->
        SettingsButton(
            text = options[i], isSelected = i == curChoice, isRadio = true
        ) {
            curChoice = i
        }
    }
}


@Composable
fun <T : Enum<*>> SettingsFieldSingleChoiceEnum(
    @StringRes name: Int,
    selectedOption: T,
    modifier: Modifier = Modifier,
    onConfirm: (enum: T) -> Unit,
) {
    val options by remember { mutableStateOf(selectedOption::class.java.names()) }
    SettingsFieldSingleChoice(
        modifier = modifier,
        name = name,
        options = options,
        selectedOption = options.indexOf(selectedOption.name),
        onConfirm = { i -> onConfirm(selectedOption::class.java.find(options[i])) }
    )
}

@Composable
fun SettingsDialog(
    @StringRes name: Int,
    onDismissRequest: () -> Unit,
    onConfirm: () -> Unit,
    itemsCount: Int,
    itemsKey: (Int) -> Any,
    modifier: Modifier = Modifier,
    showCancel: Boolean = true,
    itemContent: @Composable LazyItemScope.(Int) -> Unit
) {
    Dialog(onDismissRequest = onDismissRequest) {
        ElevatedCard(modifier = modifier) {
            Column {
                Text(
                    text = stringResource(name),
                    style = MaterialTheme.typography.headlineSmall,
                    modifier = Modifier.padding(dimensionResource(R.dimen.padding_medium))
                )
                LazyColumn(
                    modifier = Modifier
                        .heightIn(
                            0.dp,
                            (LocalConfiguration.current.screenHeightDp * 0.6).dp
                        )
                ) {
                    items(
                        count = itemsCount,
                        key = itemsKey,
                        itemContent = itemContent
                    )
                }
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.End
                ) {
                    if (showCancel) {
                        TextButton(
                            modifier = Modifier.padding(dimensionResource(R.dimen.padding_small)),
                            onClick = onDismissRequest
                        ) {
                            Text(text = stringResource(R.string.settings_choice_dialog_cancel))
                        }
                    }
                    TextButton(
                        modifier = Modifier.padding(dimensionResource(R.dimen.padding_small)),
                        onClick = onConfirm
                    ) {
                        Text(text = stringResource(R.string.settings_choice_dialog_ok))
                    }
                }
            }

        }
    }
}

@Composable
fun SettingsButton(
    text: String,
    modifier: Modifier = Modifier,
    isSelected: Boolean = false, isRadio: Boolean = false,
    addButton: Boolean = true, onClick: () -> Unit = { }
) {
    Surface(modifier = modifier, onClick = onClick) {
        Row(
            modifier = Modifier
                .padding(
                    dimensionResource(R.dimen.padding_medium),
                    dimensionResource(R.dimen.padding_small)
                )
                .fillMaxWidth(),
            verticalAlignment = Alignment.CenterVertically
        ) {
            if (addButton && isRadio) {
                RadioButton(selected = isSelected, onClick = null)
            } else if (addButton) {
                Checkbox(checked = isSelected, onCheckedChange = null)
            }
            Text(
                text = text,
                modifier = Modifier.padding(dimensionResource(R.dimen.padding_small))
            )
        }
    }
}
