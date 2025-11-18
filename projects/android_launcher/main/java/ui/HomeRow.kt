package com.alwaldend.src.kt.android.launcher.ui

import androidx.compose.foundation.ExperimentalFoundationApi
import androidx.compose.foundation.combinedClickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.ExperimentalLayoutApi
import androidx.compose.foundation.layout.FlowRow
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.defaultMinSize
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.key
import androidx.compose.runtime.remember
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.semantics.Role
import androidx.compose.ui.semantics.role
import androidx.compose.ui.semantics.semantics
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import com.alwaldend.src.kt.android.launcher.Defaults
import com.alwaldend.src.kt.android.launcher.Model
import com.alwaldend.src.kt.android.launcher.Model.Settings.HorizontalArrangement


@OptIn(ExperimentalLayoutApi::class)
@Composable
fun HomeRow(
    apps: Model.Apps,
    showHiddenApps: Boolean,
    onClick: (Model.App) -> Unit,
    onLongClick: (Model.App) -> Unit,
    sortItems: (Model.Apps, Model.Settings.Layout) -> List<Model.App>,
    getItemLabel: (Model.App, Model.Settings.AppCard) -> String,
    settings: Model.Settings,
    modifier: Modifier = Modifier,
) {
    val items = remember(apps) { sortItems(apps, settings.layout) }
    FlowRow(
        modifier = modifier
            .verticalScroll(rememberScrollState())
            .padding(
                settings.layout.horizontalPadding.dp,
                settings.layout.verticalPadding.dp
            ),
        horizontalArrangement = getArrangement(settings.layout.horizontalArrangement),
        verticalArrangement = getArrangement(settings.layout.verticalArrangement),
    ) {
        items.forEach {
            key(it.hashCode()) {
                val label = remember { getItemLabel(it, settings.appCard) }
                if (!it.isHidden || showHiddenApps) {
                    HomeRowItemCard(
                        label = label,
                        onClick = { onClick(it) },
                        onLongClick = { onLongClick(it) },
                        settings = settings.appCard
                    )
                }
            }
        }
    }
}


@OptIn(ExperimentalFoundationApi::class)
@Composable
private fun HomeRowItemCard(
    label: String,
    onClick: () -> Unit,
    onLongClick: () -> Unit,
    settings: Model.Settings.AppCard,
    modifier: Modifier = Modifier,
) {
    Box(
        modifier = modifier
            .semantics { role = Role.Button }
            .clip(ButtonDefaults.textShape)
    ) {
        Row(
            modifier = Modifier
                .combinedClickable(
                    onLongClick = onLongClick,
                    onClick = onClick
                )
                .defaultMinSize(
                    minWidth = ButtonDefaults.MinWidth,
                    minHeight = ButtonDefaults.MinHeight
                )
                .padding(ButtonDefaults.TextButtonContentPadding),
            horizontalArrangement = Arrangement.Center,
            verticalAlignment = Alignment.CenterVertically,
        ) {
            Text(
                modifier = Modifier.padding(settings.padding.dp),
                style = getTextStyle(
                    settings.textStyle,
                    MaterialTheme.typography
                ),
                fontFamily = getFontFamily(settings.fontFamily),
                color = getTextColor(settings.textColor),
                text = label
            )
        }
    }
}

private fun getTextStyle(
    textStyle: Model.Settings.TextStyle,
    typography: androidx.compose.material3.Typography
): TextStyle {
    return typography.let {
        when (textStyle) {
            Model.Settings.TextStyle.DisplayLarge -> it.displayLarge
            Model.Settings.TextStyle.DisplayMedium -> it.displayMedium
            Model.Settings.TextStyle.DisplaySmall -> it.displaySmall
            Model.Settings.TextStyle.HeadlineLarge -> it.headlineLarge
            Model.Settings.TextStyle.HeadlineMedium -> it.headlineMedium
            Model.Settings.TextStyle.HeadlineSmall -> it.headlineSmall
            Model.Settings.TextStyle.TitleLarge -> it.titleLarge
            Model.Settings.TextStyle.TitleMedium -> it.titleMedium
            Model.Settings.TextStyle.TitleSmall -> it.titleSmall
            Model.Settings.TextStyle.BodyLarge -> it.bodyLarge
            Model.Settings.TextStyle.BodyMedium -> it.bodyMedium
            Model.Settings.TextStyle.BodySmall -> it.bodySmall
            Model.Settings.TextStyle.LabelLarge -> it.labelLarge
            Model.Settings.TextStyle.LabelMedium -> it.labelMedium
            Model.Settings.TextStyle.LabelSmall -> it.labelSmall
            else -> it.bodyMedium
        }
    }
}

private fun getTextColor(textColor: Model.Settings.TextColor): Color {
    return when (textColor) {
        Model.Settings.TextColor.Unspecified -> Color.Unspecified
        Model.Settings.TextColor.Black -> Color.Black
        Model.Settings.TextColor.DarkGray -> Color.DarkGray
        Model.Settings.TextColor.Gray -> Color.Gray
        Model.Settings.TextColor.LightGray -> Color.LightGray
        Model.Settings.TextColor.White -> Color.White
        Model.Settings.TextColor.Red -> Color.Red
        Model.Settings.TextColor.Green -> Color.Green
        Model.Settings.TextColor.Blue -> Color.Blue
        Model.Settings.TextColor.Yellow -> Color.Yellow
        Model.Settings.TextColor.Cyan -> Color.Cyan
        Model.Settings.TextColor.Magenta -> Color.Magenta
        Model.Settings.TextColor.Transparent -> Color.Transparent
        else -> Color.Unspecified
    }
}

private fun getFontFamily(fontFamily: Model.Settings.FontFamily): FontFamily {
    return when (fontFamily) {
        Model.Settings.FontFamily.Default -> FontFamily.Default
        Model.Settings.FontFamily.Cursive -> FontFamily.Cursive
        Model.Settings.FontFamily.Monospace -> FontFamily.Monospace
        Model.Settings.FontFamily.SansSerif -> FontFamily.SansSerif
        Model.Settings.FontFamily.Serif -> FontFamily.Serif
        else -> FontFamily.Default
    }
}

@Preview
@Composable
private fun HomeItemCardPreview() {
    LauncherTheme(darkTheme = true) {
        HomeRowItemCard(
            modifier = Modifier.padding(30.dp),
            label = "app",
            onClick = { /*TODO*/ },
            onLongClick = { /*TODO*/ },
            settings = Defaults.AppCardSettings
        )
    }
}

private fun getArrangement(arrangement: HorizontalArrangement): Arrangement.Horizontal {
    return when (arrangement) {
        HorizontalArrangement.HorizontalStart -> Arrangement.Start
        HorizontalArrangement.HorizontalEnd -> Arrangement.End
        HorizontalArrangement.HorizontalCenter -> Arrangement.Center
        HorizontalArrangement.HorizontalSpaceEvenly -> Arrangement.SpaceEvenly
        HorizontalArrangement.HorizontalSpaceBetween -> Arrangement.SpaceBetween
        HorizontalArrangement.HorizontalSpaceAround -> Arrangement.SpaceAround
        else -> Arrangement.Center
    }
}

private fun getArrangement(arrangement: Model.Settings.VerticalArrangement): Arrangement.Vertical {
    return when (arrangement) {
        Model.Settings.VerticalArrangement.VerticalTop -> Arrangement.Top
        Model.Settings.VerticalArrangement.VerticalBottom -> Arrangement.Bottom
        Model.Settings.VerticalArrangement.VerticalCenter -> Arrangement.Center
        Model.Settings.VerticalArrangement.VerticalSpaceEvenly -> Arrangement.SpaceEvenly
        Model.Settings.VerticalArrangement.VerticalSpaceBetween -> Arrangement.SpaceBetween
        Model.Settings.VerticalArrangement.VerticalSpaceAround -> Arrangement.SpaceAround
        else -> Arrangement.Center
    }
}
