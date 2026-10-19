package com.cflat.blink.util

import android.content.Context

object PrefsManager {

    private const val PREFS_NAME = "eye_care_prefs"
    private const val KEY_AUTO_START = "auto_start_on_boot"
    private const val KEY_WORK_DURATION_MS = "work_duration_ms"
    private const val KEY_REST_DURATION_MS = "rest_duration_ms"
    private const val KEY_ALERT_MODE = "alert_mode"            // "both" | "overlay" | "notification"
    private const val KEY_PAUSE_MEDIA_DURING_REST = "pause_media_during_rest"
    private const val KEY_OVERLAY_DISMISSABLE = "overlay_dismissable"
    private const val KEY_OVERLAY_OPACITY = "overlay_opacity"  // 0.4..1.0
    private const val KEY_OVERLAY_STYLE = "overlay_style"      // "expressive" | "calm" | "minimal"
    private const val KEY_OVERLAY_THEME_MODE = "overlay_theme_mode" // "system" | "dark" | "light"
    private const val KEY_OVERLAY_APP_WHITELIST = "overlay_app_whitelist"
    private const val KEY_OVERLAY_DISMISS_BUTTON_STYLE = "overlay_dismiss_button_style" // "icon" | "icon_text" | "text"
    private const val KEY_SHIZUKU_MODE = "shizuku_mode"

    private const val DEFAULT_WORK_DURATION_MS = 20 * 60 * 1000L
    private const val DEFAULT_REST_DURATION_MS = 20 * 1000L

    private fun prefs(context: Context) =
        context.getSharedPreferences(PREFS_NAME, Context.MODE_PRIVATE)

    fun isAutoStartEnabled(context: Context): Boolean =
        prefs(context).getBoolean(KEY_AUTO_START, false)

    fun setAutoStartEnabled(context: Context, enabled: Boolean) =
        prefs(context).edit().putBoolean(KEY_AUTO_START, enabled).apply()

    fun getWorkDurationMs(context: Context): Long =
        prefs(context).getLong(KEY_WORK_DURATION_MS, DEFAULT_WORK_DURATION_MS)

    fun setWorkDurationMs(context: Context, value: Long) =
        prefs(context).edit().putLong(KEY_WORK_DURATION_MS, value).apply()

    fun getRestDurationMs(context: Context): Long =
        prefs(context).getLong(KEY_REST_DURATION_MS, DEFAULT_REST_DURATION_MS)

    fun setRestDurationMs(context: Context, value: Long) =
        prefs(context).edit().putLong(KEY_REST_DURATION_MS, value).apply()

    fun getAlertMode(context: Context): String {
        migrateLegacyAlertMode(context)
        return prefs(context).getString(KEY_ALERT_MODE, "both") ?: "both"
    }

    fun setAlertMode(context: Context, mode: String) =
        prefs(context).edit().putString(KEY_ALERT_MODE, mode).apply()

    fun isPauseMediaDuringRestEnabled(context: Context): Boolean {
        migrateLegacyAlertMode(context)
        return prefs(context).getBoolean(KEY_PAUSE_MEDIA_DURING_REST, false)
    }

    fun setPauseMediaDuringRestEnabled(context: Context, enabled: Boolean) =
        prefs(context).edit().putBoolean(KEY_PAUSE_MEDIA_DURING_REST, enabled).apply()

    private fun migrateLegacyAlertMode(context: Context) {
        val preferences = prefs(context)
        if (preferences.getString(KEY_ALERT_MODE, "both") == "overlay_media") {
            preferences.edit()
                .putString(KEY_ALERT_MODE, "overlay")
                .putBoolean(KEY_PAUSE_MEDIA_DURING_REST, true)
                .apply()
        }
    }

    fun isOverlayDismissable(context: Context): Boolean =
        prefs(context).getBoolean(KEY_OVERLAY_DISMISSABLE, false)

    fun setOverlayDismissable(context: Context, value: Boolean) =
        prefs(context).edit().putBoolean(KEY_OVERLAY_DISMISSABLE, value).apply()

    fun getOverlayOpacity(context: Context): Float =
        prefs(context).getFloat(KEY_OVERLAY_OPACITY, 1.0f)

    fun setOverlayOpacity(context: Context, value: Float) =
        prefs(context).edit().putFloat(KEY_OVERLAY_OPACITY, value).apply()

    fun getOverlayStyle(context: Context): String =
        prefs(context).getString(KEY_OVERLAY_STYLE, "expressive") ?: "expressive"

    fun setOverlayStyle(context: Context, style: String) =
        prefs(context).edit().putString(KEY_OVERLAY_STYLE, style).apply()

    fun getOverlayThemeMode(context: Context): String =
        prefs(context).getString(KEY_OVERLAY_THEME_MODE, "system") ?: "system"

    fun setOverlayThemeMode(context: Context, mode: String) =
        prefs(context).edit().putString(KEY_OVERLAY_THEME_MODE, mode).apply()

    fun getOverlayAppWhitelist(context: Context): Set<String> =
        prefs(context).getStringSet(KEY_OVERLAY_APP_WHITELIST, emptySet())?.toSet() ?: emptySet()

    fun setOverlayAppWhitelist(context: Context, packages: Set<String>) =
        prefs(context).edit().putStringSet(KEY_OVERLAY_APP_WHITELIST, packages).apply()

    fun getOverlayDismissButtonStyle(context: Context): String =
        prefs(context).getString(KEY_OVERLAY_DISMISS_BUTTON_STYLE, "icon") ?: "icon"

    fun setOverlayDismissButtonStyle(context: Context, style: String) =
        prefs(context).edit().putString(KEY_OVERLAY_DISMISS_BUTTON_STYLE, style).apply()

    fun isShizukuModeEnabled(context: Context): Boolean =
        prefs(context).getBoolean(KEY_SHIZUKU_MODE, false)

    fun setShizukuModeEnabled(context: Context, enabled: Boolean) =
        prefs(context).edit().putBoolean(KEY_SHIZUKU_MODE, enabled).apply()
}
