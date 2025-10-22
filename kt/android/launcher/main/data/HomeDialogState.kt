package com.alwaldend.src.kt.android.launcher.data

import androidx.annotation.StringRes
import androidx.compose.runtime.Immutable

@Immutable
sealed interface HomeDialogState {

    enum class HeaderActions {
        GoToInfo, HideOrShow, Uninstall
    }

    enum class LauncherDialogAction {
        GoToSettings, HideHiddenApps, ShowHiddenApps,
    }

    data class Buttons<T>(val buttons: List<Button<T>>) {
        constructor(vararg buttons: Button<T>) : this(buttons.toList())
    }

    data class Button<T>(val label: Label, val data: T) {
        constructor(label: String, data: T) : this(Label.Str(label), data)

        constructor(@StringRes label: Int, data: T) : this(
            Label.Res(label),
            data
        )

        @Immutable
        sealed interface Label {
            data class Str(val string: String) : Label

            data class Res(@StringRes val res: Int) : Label
        }
    }
}
