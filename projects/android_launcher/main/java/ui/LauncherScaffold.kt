package com.alwaldend.src.projects.android_launcher.ui

import androidx.annotation.StringRes
import androidx.compose.foundation.layout.Column
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.material3.TopAppBar
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.stringResource


@Composable
fun LauncherScaffold(
    @StringRes label: Int,
    goBack: () -> Unit,
    showTopAppBar: Boolean,
    modifier: Modifier = Modifier,
    content: @Composable () -> Unit
) {
    Column(modifier = modifier) {
        if (showTopAppBar) {
            LauncherTopBar(
                label = stringResource(id = label),
                goBack = goBack
            )
        }
        content()
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun LauncherTopBar(label: String, goBack: () -> Unit) {
    TopAppBar(
        title = {
            Text(
                text = label,
                style = MaterialTheme.typography.headlineLarge
            )
        },
        navigationIcon = {
            IconButton(onClick = goBack) {
                Icon(
                    imageVector = Icons.AutoMirrored.Filled.ArrowBack,
                    contentDescription = "Arrow back"
                )
            }
        }
    )
}
