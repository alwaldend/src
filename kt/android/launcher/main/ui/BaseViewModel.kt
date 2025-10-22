package com.alwaldend.src.kt.android.launcher.ui

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.alwaldend.src.kt.android.launcher.Defaults
import com.alwaldend.src.kt.android.launcher.LauncherNavigator
import com.alwaldend.src.kt.android.launcher.LauncherStateRepository
import com.alwaldend.src.kt.android.launcher.data.UiState
import kotlinx.coroutines.CoroutineExceptionHandler
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.SharingStarted
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.map
import kotlinx.coroutines.flow.stateIn
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch


abstract class BaseViewModel(
    private val stateRepo: LauncherStateRepository,
    private val navigator: LauncherNavigator
) : ViewModel(), LauncherNavigator by navigator {

    private val _errorFlow = MutableStateFlow<Throwable?>(null)
    val errorFlow = _errorFlow.asStateFlow()

    private val exceptionHandler = CoroutineExceptionHandler { _, e ->
        _errorFlow.update { e }
    }

    protected var prevState = Defaults.State

    val uiStateFlow = stateRepo.observeState()
        .map { stateResult ->
            UiState.Success(stateResult.fold(
                onSuccess = {
                    prevState = it
                    it
                },
                onFailure = {
                    _errorFlow.update { it }
                    prevState
                }
            ))
        }.stateIn(
            scope = viewModelScope,
            initialValue = UiState.Loading(),
            started = SharingStarted.Lazily
        )

    protected fun stateAction(action: suspend LauncherStateRepository.() -> Unit) {
        launch { action.invoke(stateRepo) }
    }

    protected fun launch(action: suspend CoroutineScope.() -> Unit) {
        viewModelScope.launch(exceptionHandler) { action() }
    }
}
