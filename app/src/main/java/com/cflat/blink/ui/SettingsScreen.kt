package com.cflat.blink.ui

import android.content.Context
import android.content.Intent
import android.content.pm.PackageManager
import android.os.Build
import android.os.PowerManager
import android.provider.Settings as AndroidSettings
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.rounded.ArrowBack
import androidx.compose.material.icons.rounded.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import com.cflat.blink.util.PrefsManager
import com.cflat.blink.util.ShizukuUtils
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch
import kotlin.math.roundToInt

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SettingsScreen(
    context: Context,
    onBack: () -> Unit
) {
    // ── State ──────────────────────────────────────────────────────────
    var workMinutes by remember {
        mutableFloatStateOf(PrefsManager.getWorkDurationMs(context) / 60000f)
    }
    var restSeconds by remember {
        mutableFloatStateOf(PrefsManager.getRestDurationMs(context) / 1000f)
    }
    var alertMode by remember { mutableStateOf(PrefsManager.getAlertMode(context)) }
    var overlayDismissable by remember { mutableStateOf(PrefsManager.isOverlayDismissable(context)) }
    var overlayOpacity by remember { mutableFloatStateOf(PrefsManager.getOverlayOpacity(context)) }
    var overlayStyle by remember { mutableStateOf(PrefsManager.getOverlayStyle(context)) }
    var overlayThemeMode by remember { mutableStateOf(PrefsManager.getOverlayThemeMode(context)) }
    var overlayWhitelist by remember { mutableStateOf(PrefsManager.getOverlayAppWhitelist(context)) }
    var dismissButtonStyle by remember { mutableStateOf(PrefsManager.getOverlayDismissButtonStyle(context)) }
    var showAppPicker by remember { mutableStateOf(false) }
    var availableApps by remember { mutableStateOf<List<AppEntry>>(emptyList()) }
    var appSearch by remember { mutableStateOf("") }
    var shizukuMode by remember { mutableStateOf(PrefsManager.isShizukuModeEnabled(context)) }
    var shizukuBusy by remember { mutableStateOf(false) }
    val shizukuStatus by ShizukuUtils.shizukuStatus.collectAsStateWithLifecycle()
    var autoStart by remember { mutableStateOf(PrefsManager.isAutoStartEnabled(context)) }
    val scope = rememberCoroutineScope()

    LaunchedEffect(showAppPicker) {
        if (showAppPicker && availableApps.isEmpty()) {
            availableApps = loadLaunchableApps(context)
        }
    }

    DisposableEffect(shizukuMode) {
        if (shizukuMode) {
            ShizukuUtils.initialize()
            onDispose { ShizukuUtils.cleanup() }
        } else {
            onDispose { }
        }
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Settings", fontWeight = FontWeight.Bold) },
                navigationIcon = {
                    IconButton(onClick = onBack) {
                        Icon(Icons.AutoMirrored.Rounded.ArrowBack, contentDescription = "Back")
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = MaterialTheme.colorScheme.surface
                )
            )
        }
    ) { padding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .verticalScroll(rememberScrollState())
                .padding(padding)
                .padding(horizontal = 16.dp, vertical = 8.dp),
            verticalArrangement = Arrangement.spacedBy(4.dp)
        ) {

            // ── Timers ───────────────────────────────────────────────
            SettingsSectionHeader("Timers", Icons.Rounded.Timer)
            SettingsCard {
                SliderRow(
                    label = "Work Duration",
                    value = workMinutes,
                    displayText = "${workMinutes.roundToInt()} min",
                    range = 1f..60f,
                    steps = 58,
                    onValueChange = { workMinutes = it },
                    onValueChangeFinished = {
                        PrefsManager.setWorkDurationMs(
                            context, workMinutes.roundToInt() * 60 * 1000L
                        )
                    }
                )
                HorizontalDivider(modifier = Modifier.padding(horizontal = 16.dp))
                SliderRow(
                    label = "Rest Duration",
                    value = restSeconds,
                    displayText = "${restSeconds.roundToInt()} sec",
                    range = 10f..60f,
                    steps = 49,
                    onValueChange = { restSeconds = it },
                    onValueChangeFinished = {
                        PrefsManager.setRestDurationMs(
                            context, restSeconds.roundToInt() * 1000L
                        )
                    }
                )
            }

            Spacer(Modifier.height(12.dp))

            // ── Alert Mode ───────────────────────────────────────────
            SettingsSectionHeader("Alert", Icons.Rounded.Notifications)
            SettingsCard {
                Column(
                    modifier = Modifier.fillMaxWidth().padding(16.dp),
                    verticalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    Text(
                        "Alert Mode",
                        style = MaterialTheme.typography.titleSmall,
                        fontWeight = FontWeight.SemiBold
                    )
                    Text(
                        "How you're reminded to take a break",
                        style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.onSurfaceVariant
                    )
                    Spacer(Modifier.height(4.dp))
                    @OptIn(ExperimentalMaterial3Api::class)
                    SingleChoiceSegmentedButtonRow(modifier = Modifier.fillMaxWidth()) {
                        val options = listOf(
                            "both" to "Both",
                            "overlay" to "Overlay",
                            "overlay_media" to "Overlay + Media",
                            "notification" to "Notif"
                        )
                        options.forEachIndexed { index, (key, label) ->
                            UniformSegmentedButton(
                                modifier = Modifier.weight(1f),
                                shape = SegmentedButtonDefaults.itemShape(index, options.size),
                                onClick = {
                                    alertMode = key
                                    PrefsManager.setAlertMode(context, key)
                                },
                                selected = alertMode == key,
                                label = label
                            )
                        }
                    }
                }
            }

            Spacer(Modifier.height(12.dp))

            // ── Overlay Appearance ───────────────────────────────────
            SettingsSectionHeader("Overlay", Icons.Rounded.Layers)
            SettingsCard {
                // Dismissable toggle
                ListItem(
                    headlineContent = { Text("Dismissable", fontWeight = FontWeight.SemiBold) },
                    supportingContent = { Text("Show a close button during rest") },
                    leadingContent = {
                        Icon(
                            Icons.Rounded.Close,
                            contentDescription = null,
                            tint = MaterialTheme.colorScheme.primary
                        )
                    },
                    trailingContent = {
                        Switch(
                            checked = overlayDismissable,
                            onCheckedChange = {
                                overlayDismissable = it
                                PrefsManager.setOverlayDismissable(context, it)
                            }
                        )
                    },
                    colors = ListItemDefaults.colors(
                        containerColor = MaterialTheme.colorScheme.surfaceContainerLow
                    )
                )

                HorizontalDivider(modifier = Modifier.padding(horizontal = 16.dp))

                // Opacity slider
                SliderRow(
                    label = "Opacity",
                    value = overlayOpacity,
                    displayText = "${(overlayOpacity * 100).roundToInt()}%",
                    range = 0.3f..1.0f,
                    steps = 13,
                    onValueChange = { overlayOpacity = it },
                    onValueChangeFinished = {
                        PrefsManager.setOverlayOpacity(context, overlayOpacity)
                    }
                )

                HorizontalDivider(modifier = Modifier.padding(horizontal = 16.dp))

                // Style selector
                Column(
                    modifier = Modifier.fillMaxWidth().padding(16.dp),
                    verticalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    Text(
                        "Overlay Style",
                        style = MaterialTheme.typography.titleSmall,
                        fontWeight = FontWeight.SemiBold
                    )
                    Text(
                        "Visual theme of the rest screen",
                        style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.onSurfaceVariant
                    )
                    Spacer(Modifier.height(4.dp))
                    @OptIn(ExperimentalMaterial3Api::class)
                    SingleChoiceSegmentedButtonRow(modifier = Modifier.fillMaxWidth()) {
                        val options = listOf(
                            "expressive" to "Expressive",
                            "calm" to "Calm",
                            "minimal" to "Minimal"
                        )
                        options.forEachIndexed { index, (key, label) ->
                            UniformSegmentedButton(
                                modifier = Modifier.weight(1f),
                                shape = SegmentedButtonDefaults.itemShape(index, options.size),
                                onClick = {
                                    overlayStyle = key
                                    PrefsManager.setOverlayStyle(context, key)
                                },
                                selected = overlayStyle == key,
                                label = label
                            )
                        }
                    }
                }

                HorizontalDivider(modifier = Modifier.padding(horizontal = 16.dp))

                Column(
                    modifier = Modifier.fillMaxWidth().padding(16.dp),
                    verticalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    Text(
                        "Overlay Theme",
                        style = MaterialTheme.typography.titleSmall,
                        fontWeight = FontWeight.SemiBold
                    )
                    Text(
                        "Choose light, dark, or follow system theme",
                        style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.onSurfaceVariant
                    )
                    Spacer(Modifier.height(4.dp))
                    @OptIn(ExperimentalMaterial3Api::class)
                    SingleChoiceSegmentedButtonRow(modifier = Modifier.fillMaxWidth()) {
                        val options = listOf(
                            "system" to "System",
                            "dark" to "Dark",
                            "light" to "Light"
                        )
                        options.forEachIndexed { index, (key, label) ->
                            UniformSegmentedButton(
                                modifier = Modifier.weight(1f),
                                shape = SegmentedButtonDefaults.itemShape(index, options.size),
                                onClick = {
                                    overlayThemeMode = key
                                    PrefsManager.setOverlayThemeMode(context, key)
                                },
                                selected = overlayThemeMode == key,
                                label = label
                            )
                        }
                    }
                }

                HorizontalDivider(modifier = Modifier.padding(horizontal = 16.dp))

                Column(
                    modifier = Modifier.fillMaxWidth().padding(16.dp),
                    verticalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    Text(
                        "Skip Overlay in Apps",
                        style = MaterialTheme.typography.titleSmall,
                        fontWeight = FontWeight.SemiBold
                    )
                    Text(
                        "Select apps where Blink should skip full-screen overlay.",
                        style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.onSurfaceVariant
                    )
                    FilledTonalButton(
                        onClick = { showAppPicker = true },
                        shape = RoundedCornerShape(12.dp)
                    ) {
                        Text("Select Apps (${overlayWhitelist.size})")
                    }
                    TextButton(
                        onClick = {
                            overlayWhitelist = emptySet()
                            PrefsManager.setOverlayAppWhitelist(context, emptySet())
                        },
                        enabled = overlayWhitelist.isNotEmpty()
                    ) { Text("Clear Selection") }
                    TextButton(
                        onClick = {
                            val intent = Intent(AndroidSettings.ACTION_USAGE_ACCESS_SETTINGS).apply {
                                addFlags(Intent.FLAG_ACTIVITY_NEW_TASK)
                            }
                            context.startActivity(intent)
                        }
                    ) {
                        Text("Grant usage access")
                    }
                }

                HorizontalDivider(modifier = Modifier.padding(horizontal = 16.dp))

                Column(
                    modifier = Modifier.fillMaxWidth().padding(16.dp),
                    verticalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    Text(
                        "Dismiss Button Style",
                        style = MaterialTheme.typography.titleSmall,
                        fontWeight = FontWeight.SemiBold
                    )
                    Text(
                        "Customize dismiss control style on the rest overlay.",
                        style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.onSurfaceVariant
                    )
                    @OptIn(ExperimentalMaterial3Api::class)
                    SingleChoiceSegmentedButtonRow(modifier = Modifier.fillMaxWidth()) {
                        val options = listOf(
                            "icon" to "Icon",
                            "icon_text" to "Icon + Text",
                            "text" to "Text"
                        )
                        options.forEachIndexed { index, (key, label) ->
                            UniformSegmentedButton(
                                modifier = Modifier.weight(1f),
                                shape = SegmentedButtonDefaults.itemShape(index, options.size),
                                onClick = {
                                    dismissButtonStyle = key
                                    PrefsManager.setOverlayDismissButtonStyle(context, key)
                                },
                                selected = dismissButtonStyle == key,
                                label = label
                            )
                        }
                    }
                }
            }

            Spacer(Modifier.height(12.dp))

            // ── Shizuku ───────────────────────────────────────────────
            SettingsSectionHeader("Shizuku", Icons.Rounded.Security)
            SettingsCard {
                ListItem(
                    headlineContent = { Text("Enable Shizuku Mode", fontWeight = FontWeight.SemiBold) },
                    supportingContent = { Text("Request and use Shizuku only when enabled") },
                    leadingContent = {
                        Icon(
                            Icons.Rounded.Security,
                            contentDescription = null,
                            tint = if (shizukuMode) MaterialTheme.colorScheme.primary
                            else MaterialTheme.colorScheme.outline
                        )
                    },
                    trailingContent = {
                        Switch(
                            checked = shizukuMode,
                            onCheckedChange = {
                                shizukuMode = it
                                PrefsManager.setShizukuModeEnabled(context, it)
                            }
                        )
                    },
                    colors = ListItemDefaults.colors(
                        containerColor = MaterialTheme.colorScheme.surfaceContainerLow
                    )
                )

                if (shizukuMode) {
                    HorizontalDivider(modifier = Modifier.padding(horizontal = 16.dp))
                    Column(
                        modifier = Modifier.fillMaxWidth().padding(16.dp),
                        verticalArrangement = Arrangement.spacedBy(8.dp)
                    ) {
                        Text(
                            "Status: ${shizukuStatus.toReadableText()}",
                            style = MaterialTheme.typography.bodyMedium
                        )
                        Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                            FilledTonalButton(
                                onClick = { ShizukuUtils.requestPermission() },
                                shape = RoundedCornerShape(12.dp)
                            ) { Text("Request Permission") }
                            FilledTonalButton(
                                onClick = {
                                    shizukuBusy = true
                                    scope.launch {
                                        ShizukuUtils.whitelistFromBatteryOptimization(context)
                                        delay(400)
                                        shizukuBusy = false
                                    }
                                },
                                enabled = !shizukuBusy && shizukuStatus == ShizukuUtils.ShizukuStatus.READY,
                                shape = RoundedCornerShape(12.dp)
                            ) {
                                Text(if (shizukuBusy) "Applying..." else "Whitelist App")
                            }
                        }
                    }
                }
            }

            Spacer(Modifier.height(12.dp))

            // ── System ───────────────────────────────────────────────
            SettingsSectionHeader("System", Icons.Rounded.PowerSettingsNew)
            SettingsCard {
                // Battery Optimization
                val pm = context.getSystemService(Context.POWER_SERVICE) as PowerManager
                var isIgnoring by remember {
                    mutableStateOf(pm.isIgnoringBatteryOptimizations(context.packageName))
                }

                ListItem(
                    headlineContent = { Text("Battery Optimization", fontWeight = FontWeight.SemiBold) },
                    supportingContent = {
                        Text(
                            if (isIgnoring) "Unrestricted — timers are reliable"
                            else "Restricted — tap to fix"
                        )
                    },
                    leadingContent = {
                        Icon(
                            Icons.Rounded.Shield,
                            contentDescription = null,
                            tint = if (isIgnoring) MaterialTheme.colorScheme.primary
                            else MaterialTheme.colorScheme.error
                        )
                    },
                    trailingContent = {
                        if (!isIgnoring) {
                            FilledTonalButton(
                                onClick = {
                                    @Suppress("BatteryLife")
                                    val intent = Intent(
                                        AndroidSettings.ACTION_REQUEST_IGNORE_BATTERY_OPTIMIZATIONS,
                                        android.net.Uri.parse("package:${context.packageName}")
                                    )
                                    context.startActivity(intent)
                                    scope.launch {
                                        delay(1500)
                                        isIgnoring = pm.isIgnoringBatteryOptimizations(context.packageName)
                                    }
                                },
                                shape = RoundedCornerShape(12.dp)
                            ) { Text("Disable") }
                        } else {
                            Icon(
                                Icons.Rounded.CheckCircle,
                                contentDescription = null,
                                tint = MaterialTheme.colorScheme.primary,
                                modifier = Modifier.size(24.dp)
                            )
                        }
                    },
                    colors = ListItemDefaults.colors(
                        containerColor = MaterialTheme.colorScheme.surfaceContainerLow
                    )
                )

                HorizontalDivider(modifier = Modifier.padding(horizontal = 16.dp))

                ListItem(
                    headlineContent = { Text("Start on Boot", fontWeight = FontWeight.SemiBold) },
                    supportingContent = { Text("Resume automatically after reboot") },
                    leadingContent = {
                        Icon(
                            Icons.Rounded.PowerSettingsNew,
                            contentDescription = null,
                            tint = if (autoStart) MaterialTheme.colorScheme.primary
                            else MaterialTheme.colorScheme.outline
                        )
                    },
                    trailingContent = {
                        Switch(
                            checked = autoStart,
                            onCheckedChange = {
                                autoStart = it
                                PrefsManager.setAutoStartEnabled(context, it)
                            }
                        )
                    },
                    colors = ListItemDefaults.colors(
                        containerColor = MaterialTheme.colorScheme.surfaceContainerLow
                    )
                )
            }

            Spacer(Modifier.height(12.dp))

            // ── OEM-specific guidance ────────────────────────────────
            val oemGuide = getOemBatteryGuide()
            if (oemGuide != null) {
                SettingsSectionHeader("Device-Specific", Icons.Rounded.PhoneAndroid)
                SettingsCard {
                    Column(
                        modifier = Modifier.fillMaxWidth().padding(16.dp),
                        verticalArrangement = Arrangement.spacedBy(8.dp)
                    ) {
                        Text(
                            oemGuide.title,
                            style = MaterialTheme.typography.titleSmall,
                            fontWeight = FontWeight.SemiBold
                        )
                        Text(
                            oemGuide.description,
                            style = MaterialTheme.typography.bodySmall,
                            color = MaterialTheme.colorScheme.onSurfaceVariant
                        )
                        Spacer(Modifier.height(4.dp))
                        oemGuide.steps.forEach { step ->
                            Text(
                                step,
                                style = MaterialTheme.typography.bodySmall,
                                color = MaterialTheme.colorScheme.onSurface
                            )
                        }
                    }
                }
            }

            if (showAppPicker) {
                AppWhitelistDialog(
                    apps = availableApps,
                    selected = overlayWhitelist,
                    query = appSearch,
                    onQueryChange = { appSearch = it },
                    onToggle = { packageName ->
                        overlayWhitelist = overlayWhitelist.toMutableSet().apply {
                            if (contains(packageName)) remove(packageName) else add(packageName)
                        }
                        PrefsManager.setOverlayAppWhitelist(context, overlayWhitelist)
                    },
                    onDismiss = { showAppPicker = false }
                )
            }

            Spacer(Modifier.height(32.dp))
        }
    }
}

// ── Shared helpers ─────────────────────────────────────────────────────────

@Composable
private fun SettingsSectionHeader(title: String, icon: ImageVector) {
    Row(
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.spacedBy(8.dp),
        modifier = Modifier.padding(horizontal = 4.dp, vertical = 6.dp)
    ) {
        Icon(
            imageVector = icon,
            contentDescription = null,
            tint = MaterialTheme.colorScheme.primary,
            modifier = Modifier.size(18.dp)
        )
        Text(
            text = title,
            style = MaterialTheme.typography.labelLarge,
            color = MaterialTheme.colorScheme.primary,
            fontWeight = FontWeight.Bold
        )
    }
}

@Composable
private fun SettingsCard(content: @Composable ColumnScope.() -> Unit) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        shape = RoundedCornerShape(20.dp),
        colors = CardDefaults.cardColors(
            containerColor = MaterialTheme.colorScheme.surfaceContainerLow
        )
    ) {
        Column(content = content)
    }
}

@Composable
private fun SliderRow(
    label: String,
    value: Float,
    displayText: String,
    range: ClosedFloatingPointRange<Float>,
    steps: Int,
    onValueChange: (Float) -> Unit,
    onValueChangeFinished: () -> Unit
) {
    Column(
        modifier = Modifier.fillMaxWidth().padding(horizontal = 16.dp, vertical = 12.dp)
    ) {
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Text(
                label,
                style = MaterialTheme.typography.titleSmall,
                fontWeight = FontWeight.SemiBold
            )
            Text(
                displayText,
                style = MaterialTheme.typography.bodyMedium,
                color = MaterialTheme.colorScheme.primary,
                fontWeight = FontWeight.Medium
            )
        }
        Slider(
            value = value,
            onValueChange = onValueChange,
            onValueChangeFinished = onValueChangeFinished,
            valueRange = range,
            steps = steps,
            modifier = Modifier.fillMaxWidth()
        )
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun UniformSegmentedButton(
    modifier: Modifier = Modifier,
    shape: androidx.compose.ui.graphics.Shape,
    onClick: () -> Unit,
    selected: Boolean,
    label: String
) {
    SegmentedButton(
        modifier = modifier.heightIn(min = 44.dp),
        shape = shape,
        onClick = onClick,
        selected = selected,
        label = {
            Text(
                text = label,
                style = MaterialTheme.typography.labelMedium,
                maxLines = 1,
                overflow = TextOverflow.Ellipsis
            )
        }
    )
}

private data class AppEntry(
    val packageName: String,
    val label: String
)

@Composable
private fun AppWhitelistDialog(
    apps: List<AppEntry>,
    selected: Set<String>,
    query: String,
    onQueryChange: (String) -> Unit,
    onToggle: (String) -> Unit,
    onDismiss: () -> Unit
) {
    val filtered = remember(apps, query) {
        val q = query.trim().lowercase()
        if (q.isBlank()) apps
        else apps.filter { it.label.lowercase().contains(q) || it.packageName.lowercase().contains(q) }
    }

    AlertDialog(
        onDismissRequest = onDismiss,
        title = { Text("Select Apps") },
        text = {
            Column(verticalArrangement = Arrangement.spacedBy(8.dp)) {
                OutlinedTextField(
                    value = query,
                    onValueChange = onQueryChange,
                    singleLine = true,
                    placeholder = { Text("Search app") },
                    modifier = Modifier.fillMaxWidth()
                )
                LazyColumn(
                    modifier = Modifier
                        .fillMaxWidth()
                        .heightIn(min = 180.dp, max = 420.dp),
                    verticalArrangement = Arrangement.spacedBy(2.dp)
                ) {
                    items(filtered, key = { it.packageName }) { app ->
                        val checked = selected.contains(app.packageName)
                        ListItem(
                            modifier = Modifier.clickable { onToggle(app.packageName) },
                            headlineContent = {
                                Text(
                                    app.label,
                                    maxLines = 1,
                                    overflow = TextOverflow.Ellipsis
                                )
                            },
                            supportingContent = {
                                Text(
                                    app.packageName,
                                    maxLines = 1,
                                    overflow = TextOverflow.Ellipsis
                                )
                            },
                            trailingContent = {
                                Checkbox(
                                    checked = checked,
                                    onCheckedChange = { onToggle(app.packageName) }
                                )
                            },
                            colors = ListItemDefaults.colors(
                                containerColor = MaterialTheme.colorScheme.surfaceContainerLow
                            )
                        )
                    }
                }
            }
        },
        confirmButton = {
            TextButton(onClick = onDismiss) { Text("Done") }
        }
    )
}

private fun loadLaunchableApps(context: Context): List<AppEntry> {
    val pm = context.packageManager
    val intent = Intent(Intent.ACTION_MAIN).apply {
        addCategory(Intent.CATEGORY_LAUNCHER)
    }
    return pm.queryIntentActivities(intent, PackageManager.MATCH_ALL)
        .map {
            val packageName = it.activityInfo.packageName
            val label = it.loadLabel(pm)?.toString().orEmpty().ifBlank { packageName }
            AppEntry(packageName = packageName, label = label)
        }
        .distinctBy { it.packageName }
        .sortedBy { it.label.lowercase() }
}

private fun ShizukuUtils.ShizukuStatus.toReadableText(): String =
    when (this) {
        ShizukuUtils.ShizukuStatus.UNKNOWN -> "Unknown"
        ShizukuUtils.ShizukuStatus.NOT_INSTALLED -> "Shizuku not installed"
        ShizukuUtils.ShizukuStatus.NOT_RUNNING -> "Shizuku not running"
        ShizukuUtils.ShizukuStatus.PERMISSION_DENIED -> "Permission required"
        ShizukuUtils.ShizukuStatus.READY -> "Ready"
        ShizukuUtils.ShizukuStatus.WHITELISTED -> "Whitelisted"
    }

// ── OEM battery guidance ──────────────────────────────────────────────────────

private data class OemBatteryGuide(
    val title: String,
    val description: String,
    val steps: List<String>
)

private fun getOemBatteryGuide(): OemBatteryGuide? {
    val manufacturer = Build.MANUFACTURER.lowercase()
    return when {
        manufacturer.contains("xiaomi") || manufacturer.contains("redmi") || manufacturer.contains("poco") ->
            OemBatteryGuide(
                title = "Xiaomi/MIUI Battery Settings",
                description = "MIUI aggressively kills background apps. Follow these steps:",
                steps = listOf(
                    "1. Settings > Apps > Manage Apps > Blink > Autostart > enable",
                    "2. Settings > Apps > Manage Apps > Blink > Battery Saver > No restrictions",
                    "3. Security app > Battery > App Battery Saver > Blink > No restrictions",
                    "4. Lock Blink in Recents (swipe down on the app card)",
                    "",
                    "Note: If the system force-stops Blink, autostart on boot won't work until you open the app once."
                )
            )
        manufacturer.contains("samsung") ->
            OemBatteryGuide(
                title = "Samsung Battery Settings",
                description = "Samsung's battery optimization can stop background apps. Follow these steps:",
                steps = listOf(
                    "1. Settings > Battery > Background usage limits > Never sleeping apps > Add Blink",
                    "2. Settings > Apps > Blink > Battery > Unrestricted",
                    "3. Lock Blink in Recents (tap the app icon > Lock)",
                    "",
                    "Note: If the system force-stops Blink, autostart on boot won't work until you open the app once."
                )
            )
        manufacturer.contains("huawei") || manufacturer.contains("honor") ->
            OemBatteryGuide(
                title = "Huawei/Honor Battery Settings",
                description = "EMUI has aggressive power management. Follow these steps:",
                steps = listOf(
                    "1. Settings > Battery > App launch > Blink > Manage manually > enable all toggles",
                    "2. Settings > Apps > Apps > Blink > Battery > enable Run in background",
                    "3. Lock Blink in Recents (swipe down on the app card)"
                )
            )
        manufacturer.contains("oppo") || manufacturer.contains("realme") || manufacturer.contains("oneplus") ->
            OemBatteryGuide(
                title = "OPPO/Realme/OnePlus Battery Settings",
                description = "ColorOS aggressively manages background apps. Follow these steps:",
                steps = listOf(
                    "1. Settings > App management > App list > Blink > Auto-launch > Allow",
                    "2. Settings > Battery > More settings > Optimize battery use > Blink > Don't optimize",
                    "3. Settings > Apps > App management > Blink > Battery usage > Allow background activity",
                    "4. Lock Blink in Recents (tap the lock icon on the app card)",
                    "",
                    "Note: If the system force-stops Blink, autostart on boot won't work until you open the app once."
                )
            )
        manufacturer.contains("vivo") || manufacturer.contains("iqoo") ->
            OemBatteryGuide(
                title = "Vivo/iQOO Battery Settings",
                description = "Funtouch OS restricts background apps. Follow these steps:",
                steps = listOf(
                    "1. Settings > Apps > Autostart > enable Blink",
                    "2. Settings > Battery > Background power consumption > Blink > Don't restrict",
                    "3. Settings > Apps > Blink > Battery > Allow high background power consumption",
                    "4. Lock Blink in Recents (swipe down on the app card)",
                    "",
                    "Note: If the system force-stops Blink, autostart on boot won't work until you open the app once."
                )
            )
        manufacturer.contains("motorola") || manufacturer.contains("lenovo") ->
            OemBatteryGuide(
                title = "Motorola Battery Settings",
                description = "Moto's system can still kill foreground services. Follow these steps:",
                steps = listOf(
                    "1. Settings > Battery > Battery optimization > All apps > Blink > Don't optimize",
                    "2. Settings > Apps > Blink > Battery > Allow background activity",
                    "3. Settings > Battery > Adaptive battery > turn off (optional, affects all apps)",
                    "",
                    "Note: If the system force-stops Blink, autostart on boot won't work until you open the app once."
                )
            )
        manufacturer.contains("asus") ->
            OemBatteryGuide(
                title = "ASUS Battery Settings",
                description = "ASUS power management can restrict background apps. Follow these steps:",
                steps = listOf(
                    "1. Settings > Battery > PowerMaster > Battery optimization > Blink > Don't optimize",
                    "2. Settings > Battery > PowerMaster > Auto-start manager > enable Blink"
                )
            )
        manufacturer.contains("nokia") ->
            OemBatteryGuide(
                title = "Nokia Battery Settings",
                description = "Nokia uses DuraSpeed which limits background processes. Follow these steps:",
                steps = listOf(
                    "1. Settings > Battery > Background activity > enable Blink",
                    "2. Settings > Apps > Blink > Battery > Unrestricted"
                )
            )
        else -> null
    }
}
