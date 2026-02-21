package com.alwaldend.src.projects.android_launcher.data

import androidx.compose.runtime.Immutable

@Immutable
sealed interface UiState<T> {
    data class Loading<T>(val placeholder: Unit = Unit) : UiState<T>
    data class Success<T>(val state: T) : UiState<T>
}
