package com.alwaldend.src.kt.android.launcher

import android.app.role.RoleManager
import android.content.ComponentName
import android.content.Context
import android.content.Intent
import android.content.Intent.ACTION_MAIN
import android.content.Intent.CATEGORY_HOME
import android.content.IntentFilter
import android.content.pm.ActivityInfo
import android.content.pm.LauncherApps
import android.content.pm.PackageManager
import android.content.pm.ResolveInfo
import android.content.pm.ShortcutInfo
import android.graphics.Bitmap
import android.graphics.drawable.Drawable
import android.net.Uri
import android.os.Build
import android.os.Process
import android.os.UserHandle
import android.provider.Settings
import android.util.Log
import androidx.core.graphics.drawable.toBitmap
import androidx.lifecycle.DefaultLifecycleObserver
import androidx.lifecycle.Lifecycle
import androidx.lifecycle.LifecycleOwner
import com.google.protobuf.ByteString
import java.io.ByteArrayOutputStream


private val tag = LauncherManager::class.simpleName

class LauncherManager(
    private val context: Context,
    private val lifecycle: Lifecycle,
    private val packageManager: PackageManager,
    private val launcherAppsService: LauncherApps
) {

    companion object {
        private val callbacks = mutableListOf<LauncherApps.Callback>()
    }

    fun isHomeApp(): Boolean {
        return when {
            Build.VERSION.SDK_INT >= 29 -> context
                .getSystemService(RoleManager::class.java)
                .isRoleHeld(RoleManager.ROLE_HOME)

            else -> listOf<ComponentName>().apply {
                packageManager
                    .getPreferredActivities(
                        listOf(
                            IntentFilter(ACTION_MAIN).apply {
                                addCategory(CATEGORY_HOME)
                            }
                        ),
                        this,
                        context.packageName
                    )
            }.isNotEmpty()
        }
    }

    fun launchApp(packageName: String, flags: Int? = null) {
        val intent = packageManager
            .getLaunchIntentForPackage(packageName)
            .let { if (flags == null) it else it?.setFlags(flags) }
        context.startActivity(intent)
    }

    fun launchAppInfo(packageName: String) {
        val intent = Intent(
            Settings.ACTION_APPLICATION_DETAILS_SETTINGS,
            Uri.parse("package:${packageName}")
        )
        context.startActivity(intent)
    }

    fun launchAppUninstall(packageName: String) {
        val intent = Intent(
            Intent.ACTION_DELETE,
            Uri.parse("package:${packageName}")
        )
        context.startActivity(intent)
    }

    fun launchAppShortcut(shortcut: Model.AppShortcut) {
        launcherAppsService.startShortcut(
            shortcut.packageName,
            shortcut.shortcutId,
            null,
            null,
            Process.myUserHandle()
        )
    }

    fun removeCallbacks() {
        for (callback in callbacks) {
            launcherAppsService.unregisterCallback(callback)
        }
        callbacks.clear()
    }

    private fun queryPackageManager(intent: Intent): List<ResolveInfo> {
        return when {
            Build.VERSION.SDK_INT >= 33 -> packageManager.queryIntentActivities(
                intent,
                PackageManager.ResolveInfoFlags.of(PackageManager.MATCH_ALL.toLong())
            )

            else -> packageManager.queryIntentActivities(intent, 0)
        }
    }

    private fun getAppShortcuts(packageName: String): List<Model.AppShortcut> {
        val userHandle = Process.myUserHandle()
        val query = LauncherApps.ShortcutQuery()
            .setPackage(packageName)
            .setQueryFlags(
                LauncherApps.ShortcutQuery.FLAG_MATCH_DYNAMIC or
                        LauncherApps.ShortcutQuery.FLAG_MATCH_PINNED or
                        LauncherApps.ShortcutQuery.FLAG_MATCH_MANIFEST
            )
        val shortcuts = try {
            launcherAppsService.getShortcuts(query, userHandle) ?: listOf()
        } catch (e: SecurityException) {
            Log.d(tag, "got a security exception: $e")
            listOf()
        }
        return shortcuts.map { info ->
            val label = info.longLabel ?: info.shortLabel ?: info.`package`
            Model.AppShortcut.newBuilder()
                .setShortcutId(info.id)
                .setPackageName(info.`package`)
                .setLabel(label.toString())
                .build()
        }
    }

    fun getApp(packageName: String): Model.App? {
        val queryResults = queryPackageManager(
            Intent(ACTION_MAIN, null)
                .addCategory(Intent.CATEGORY_LAUNCHER)
                .setPackage(packageName)
        )
        if (queryResults.isEmpty()) {
            return null
        }
        return infoToApp(queryResults[0].activityInfo)
    }

    private fun infoToApp(info: ActivityInfo): Model.App {
        return Model.App.newBuilder()
            .setLabel(info.loadLabel(packageManager).toString())
            .setPackageName(info.packageName)
            .addAllShortcuts(getAppShortcuts(info.packageName))
            .setIcon(compressIcon(info.loadIcon(packageManager)))
            .build()
    }

    private fun compressIcon(icon: Drawable): ByteString {
        val stream = ByteArrayOutputStream()
        icon.toBitmap(config = Bitmap.Config.ARGB_8888)
            .compress(Bitmap.CompressFormat.PNG, 100, stream)
        return ByteString.copyFrom(stream.toByteArray())
    }

    fun fetchAllApps(): Model.Apps {
        val queryResults = queryPackageManager(
            Intent(ACTION_MAIN, null)
                .addCategory(Intent.CATEGORY_LAUNCHER)
        )
        val apps = queryResults.associate { info ->
            val activity = info.activityInfo
            activity.packageName to infoToApp(activity)
        }
        return Model.Apps.newBuilder().putAllApps(apps).build()
    }

    fun addCallback(
        onRemoved: (packageName: String) -> Unit,
        onChanged: (packageName: String) -> Unit
    ) = launcherAppsService.registerCallback(object : LauncherApps.Callback() {
        override fun onPackageRemoved(
            packageName: String, user: UserHandle
        ) = onRemoved(packageName)

        override fun onPackageAdded(
            packageName: String, user: UserHandle
        ) = onChanged(packageName)

        override fun onPackageChanged(
            packageName: String, p2: UserHandle
        ) = onChanged(packageName)

        override fun onShortcutsChanged(
            packageName: String,
            shortcuts: MutableList<ShortcutInfo>,
            user: UserHandle
        ) = onChanged(packageName)

        override fun onPackagesAvailable(
            packageNames: Array<out String>?,
            user: UserHandle?,
            replacing: Boolean
        ) = Unit

        override fun onPackagesUnavailable(
            packageNames: Array<out String>?,
            user: UserHandle?,
            replacing: Boolean
        ) = Unit
    })

    fun launchSplitScreen(appFirst: String, appSecond: String) {
        lifecycle.addObserver(object : DefaultLifecycleObserver {
            override fun onStop(owner: LifecycleOwner) {
                lifecycle.removeObserver(this)
                launchApp(
                    packageName = appSecond,
                    flags = Intent.FLAG_ACTIVITY_NEW_TASK
                            or Intent.FLAG_ACTIVITY_LAUNCH_ADJACENT
                            or Intent.FLAG_ACTIVITY_CLEAR_TOP
                )
            }
        })
        launchApp(
            packageName = appFirst,
            flags = Intent.FLAG_ACTIVITY_NEW_TASK or Intent.FLAG_ACTIVITY_CLEAR_TOP
        )
    }
}
