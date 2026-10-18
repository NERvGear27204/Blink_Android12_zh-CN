package com.cflat.blink.util

import android.app.usage.UsageEvents
import android.app.usage.UsageStatsManager
import android.content.Context

object ForegroundAppUtils {

    fun getForegroundPackage(context: Context, lookbackMs: Long = 15_000L): String? {
        val usageStatsManager =
            context.getSystemService(Context.USAGE_STATS_SERVICE) as? UsageStatsManager ?: return null

        val endTime = System.currentTimeMillis()
        val events = usageStatsManager.queryEvents(endTime - lookbackMs, endTime)
        val event = UsageEvents.Event()
        var foregroundPackage: String? = null

        while (events.hasNextEvent()) {
            events.getNextEvent(event)
            if (event.eventType == UsageEvents.Event.MOVE_TO_FOREGROUND ||
                event.eventType == UsageEvents.Event.ACTIVITY_RESUMED
            ) {
                val packageName = event.packageName
                if (!packageName.isNullOrBlank()) {
                    foregroundPackage = packageName
                }
            }
        }

        return foregroundPackage
    }
}
