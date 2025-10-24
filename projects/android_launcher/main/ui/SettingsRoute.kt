package com.alwaldend.src.kt.android.launcher.ui

import android.net.Uri
import android.os.ParcelFileDescriptor
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableIntStateOf
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.stringResource
import androidx.lifecycle.viewmodel.compose.viewModel
import com.alwaldend.src.kt.android.launcher.Model
import com.alwaldend.src.kt.android.launcher.R
import com.alwaldend.src.kt.android.launcher.data.UiState
import java.io.InputStream

private val layoutChoices = object {
    val layoutHorizontalPadding = (0..150 step 10).map { it.toString() }
    val layoutVerticalPadding = (0..150 step 10).map { it.toString() }
}


object SettingsRoute : LauncherRoute<SettingsViewModel> {
    override val url = "settings_screen"
    override val label = R.string.navigation_settings
    override val showScaffold = true

    @Composable
    override fun getViewModel(): SettingsViewModel {
        return viewModel(factory = SettingsViewModel.Factory)
    }

    @Composable
    override fun Content(
        vm: SettingsViewModel,
        uiState: UiState.Success<Model.State>
    ) {
        SettingsScreen(
            state = uiState.state,
            writeSettingsToFile = vm::writeSettingsToFile,
            updateSettingsFromStream = vm::updateSettingsFromStream,
            setPadding = vm::setAppCardPadding,
            setTextColor = vm::setAppCardTextColor,
            setTextStyle = vm::setAppCardTextStyle,
            setFontFamily = vm::setAppCardFontFamily,
            setSortBy = vm::setAppLayoutSortBy,
            setLabelLowercase = vm::setAppCardLabelLowercase,
            setVerticalArrangement = vm::setAppLayoutVerticalArrangement,
            setHorizontalArrangement = vm::setAppLayoutHorizontalArrangement,
            setVerticalPadding = vm::setAppLayoutVerticalPadding,
            setHorizontalPadding = vm::setAppLayoutHorizontalPadding,
            backupReset = vm::backupReset,
            setLabelRemoveSpaces = vm::setAppCardLabelRemoveSpaces,
            setReverseOrder = vm::setAppLayoutReverseOrder
        )
    }
}

@Composable
private fun SettingsScreen(
    modifier: Modifier = Modifier,
    state: Model.State,
    writeSettingsToFile: (() -> ParcelFileDescriptor?) -> Unit,
    backupReset: () -> Unit,
    updateSettingsFromStream: (() -> InputStream?) -> Unit,
    setReverseOrder: (Boolean) -> Unit,
    setHorizontalPadding: (Int) -> Unit,
    setVerticalPadding: (Int) -> Unit,
    setHorizontalArrangement: (Model.Settings.HorizontalArrangement) -> Unit,
    setVerticalArrangement: (Model.Settings.VerticalArrangement) -> Unit,
    setSortBy: (Model.Settings.SortBy) -> Unit,
    setLabelLowercase: (Boolean) -> Unit,
    setLabelRemoveSpaces: (Boolean) -> Unit,
    setFontFamily: (Model.Settings.FontFamily) -> Unit,
    setTextStyle: (Model.Settings.TextStyle) -> Unit,
    setTextColor: (Model.Settings.TextColor) -> Unit,
    setPadding: (Int) -> Unit
) {
    Column(modifier = modifier.verticalScroll(rememberScrollState())) {
        SettingsGroup(R.string.settings_group_general) {
            BackupExport(writeSettingsToFile = writeSettingsToFile)
            BackupImport(updateSettingsFromStream = updateSettingsFromStream)
            BackupReset(resetSettings = backupReset)
        }
        SettingsGroup(R.string.settings_group_layout) {
            val settings = state.settings.layout
            LayoutReverseOrder(
                reverse = settings.reverseOrder,
                setReverse = setReverseOrder
            )
            LayoutHorizontalPadding(
                padding = settings.horizontalPadding,
                setPadding = setHorizontalPadding
            )
            LayoutVerticalPadding(
                padding = settings.verticalPadding,
                setPadding = setVerticalPadding
            )
            LayoutHorizontalArrangement(
                arrangement = settings.horizontalArrangement,
                setArrangement = setHorizontalArrangement
            )
            LayoutVerticalArrangement(
                arrangement = settings.verticalArrangement,
                setArrangement = setVerticalArrangement
            )
            LayoutSortBy(
                sortBy = settings.sortBy,
                setSortBy = setSortBy
            )
        }
        SettingsGroup(R.string.settings_group_app_card) {
            val settings = state.settings.appCard
            AppCardRemoveSpaces(
                removeSpaces = settings.labelRemoveSpaces,
                setRemoveSpaces = setLabelRemoveSpaces
            )
            AppCardLowercase(
                lowercase = settings.labelLowercase,
                setLowercase = setLabelLowercase
            )
            AppCardFontFamily(
                fontFamily = settings.fontFamily,
                setFontFamily = setFontFamily
            )
            AppCardTextStyle(
                textStyle = settings.textStyle,
                setTextStyle = setTextStyle
            )
            AppCardTextColor(
                color = settings.textColor,
                setColor = setTextColor
            )
            AppCardPadding(
                padding = settings.padding,
                setPadding = setPadding
            )
        }
    }
}

@Composable
private fun BackupReset(resetSettings: () -> Unit) {
    var showDialog by remember { mutableStateOf(false) }
    SettingsField(
        name = R.string.settings_backup_reset,
        label = stringResource(R.string.settings_backup_reset_label),
        onClick = { showDialog = true }
    )
    if (showDialog) {
        SettingsDialog(
            name = R.string.settings_backup_reset,
            onConfirm = { resetSettings(); showDialog = false },
            onDismissRequest = { showDialog = false },
            itemsCount = 1,
            itemsKey = { 0 }
        ) {
            SettingsButton(
                text = stringResource(R.string.settings_backup_reset_confirmation),
                addButton = false
            )
        }
    }
}

@Composable
private fun BackupImport(
    updateSettingsFromStream: (() -> InputStream?) -> Unit
) {
    var result by remember { mutableStateOf<Uri?>(null) }
    val launcher = rememberLauncherForActivityResult(
        contract = ActivityResultContracts.OpenDocument()
    ) { result = it }
    SettingsField(
        name = R.string.settings_backup_import,
        label = stringResource(R.string.settings_backup_import_label),
        onClick = { launcher.launch(arrayOf("application/*")) }
    )
    val contentResolver = LocalContext.current.contentResolver
    LaunchedEffect(result) {
        result?.let { uri ->
            updateSettingsFromStream { contentResolver.openInputStream(uri) }
            result = null
        }
    }
}

@Composable
private fun BackupExport(
    writeSettingsToFile: (() -> ParcelFileDescriptor?) -> Unit
) {
    var result by remember { mutableStateOf<Uri?>(null) }
    val launcher = rememberLauncherForActivityResult(
        ActivityResultContracts.CreateDocument("application/binpb")
    ) {
        result = it
    }
    val filename = stringResource(R.string.settings_backup_export_filename)
    SettingsField(
        name = R.string.settings_backup_export,
        label = stringResource(R.string.settings_backup_export_label),
        onClick = { launcher.launch(filename) }
    )
    val contentResolver = LocalContext.current.contentResolver
    LaunchedEffect(result) {
        result?.let { uri ->
            writeSettingsToFile {
                contentResolver.openFileDescriptor(uri, "w")
            }
            result = null
        }
    }
}

@Composable
private fun LayoutVerticalPadding(
    padding: Int,
    setPadding: (Int) -> Unit
) {
    var curOption by remember {
        mutableIntStateOf(layoutChoices.layoutVerticalPadding.indexOf(
            padding.toString()
        ).let { if (it != -1) it else 0 })
    }
    SettingsFieldSingleChoice(
        name = R.string.settings_layout_vertical_padding,
        selectedOption = curOption,
        options = layoutChoices.layoutVerticalPadding,
        onConfirm = { choice ->
            curOption = choice
            setPadding(layoutChoices.layoutVerticalPadding[choice].toInt())
        }
    )
}

@Composable
private fun LayoutHorizontalPadding(
    padding: Int,
    setPadding: (Int) -> Unit
) {
    var curOption by remember {
        mutableIntStateOf(layoutChoices.layoutHorizontalPadding.indexOf(
            padding.toString()
        ).let { if (it != -1) it else 0 })
    }
    SettingsFieldSingleChoice(
        name = R.string.settings_layout_horizontal_padding,
        selectedOption = curOption,
        options = layoutChoices.layoutHorizontalPadding,
        onConfirm = { choice ->
            curOption = choice
            setPadding(layoutChoices.layoutHorizontalPadding[choice].toInt())
        }
    )
}

@Composable
private fun LayoutSortBy(
    sortBy: Model.Settings.SortBy,
    setSortBy: (Model.Settings.SortBy) -> Unit
) {
    SettingsFieldSingleChoiceEnum(
        name = R.string.settings_layout_sort_by,
        selectedOption = sortBy,
        onConfirm = { setSortBy(it) }
    )
}

@Composable
private fun LayoutReverseOrder(
    reverse: Boolean,
    setReverse: (Boolean) -> Unit
) {
    SettingsFieldSwitch(
        name = R.string.settings_layout_reverse_order,
        label = R.string.settings_layout_reverse_order_label,
        isToggled = reverse,
        onClick = { setReverse(reverse.not()) }
    )
}

@Composable
private fun LayoutVerticalArrangement(
    arrangement: Model.Settings.VerticalArrangement,
    setArrangement: (Model.Settings.VerticalArrangement) -> Unit
) {
    SettingsFieldSingleChoiceEnum(
        name = R.string.settings_layout_vertical_arrangement,
        selectedOption = arrangement,
        onConfirm = setArrangement
    )
}


@Composable
private fun LayoutHorizontalArrangement(
    arrangement: Model.Settings.HorizontalArrangement,
    setArrangement: (Model.Settings.HorizontalArrangement) -> Unit
) {
    SettingsFieldSingleChoiceEnum(
        name = R.string.settings_layout_horizontal_arrangement,
        selectedOption = arrangement,
        onConfirm = setArrangement
    )
}

@Composable
private fun AppCardFontFamily(
    fontFamily: Model.Settings.FontFamily,
    setFontFamily: (Model.Settings.FontFamily) -> Unit
) {
    SettingsFieldSingleChoiceEnum(
        name = R.string.settings_app_card_font_family,
        selectedOption = fontFamily,
        onConfirm = setFontFamily
    )
}

@Composable
private fun AppCardTextStyle(
    textStyle: Model.Settings.TextStyle,
    setTextStyle: (Model.Settings.TextStyle) -> Unit
) {
    SettingsFieldSingleChoiceEnum(
        name = R.string.settings_app_card_text_style,
        selectedOption = textStyle,
        onConfirm = setTextStyle
    )
}

@Composable
private fun AppCardRemoveSpaces(
    removeSpaces: Boolean,
    setRemoveSpaces: (Boolean) -> Unit
) {
    SettingsFieldSwitch(
        name = R.string.settings_app_label_remove_spaces,
        label = R.string.settings_app_label_remove_spaces_label,
        isToggled = removeSpaces,
        onClick = { setRemoveSpaces(removeSpaces.not()) }
    )
}

@Composable
private fun AppCardLowercase(
    lowercase: Boolean,
    setLowercase: (Boolean) -> Unit
) {
    SettingsFieldSwitch(
        name = R.string.settings_app_label_lowercase,
        label = R.string.settings_app_label_lowercase_label,
        isToggled = lowercase,
        onClick = { setLowercase(lowercase.not()) }
    )
}

@Composable
private fun AppCardTextColor(
    color: Model.Settings.TextColor,
    setColor: (Model.Settings.TextColor) -> Unit
) {
    SettingsFieldSingleChoiceEnum(
        name = R.string.settings_app_card_text_color,
        selectedOption = color,
        onConfirm = setColor
    )
}

@Composable
private fun AppCardPadding(
    padding: Int,
    setPadding: (Int) -> Unit,
) {
    val options by remember { mutableStateOf((0..30).map { it.toString() }) }
    var curOption by remember { mutableIntStateOf(options.indexOf(padding.toString())) }
    SettingsFieldSingleChoice(
        name = R.string.settings_app_card_padding,
        selectedOption = if (curOption == -1) 0 else curOption,
        options = options,
        onConfirm = { curOption = it; setPadding(options[it].toInt()) }
    )
}
