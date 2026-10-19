package com.cflat.blink.util

import android.content.Context
import android.content.ComponentName
import android.content.ServiceConnection
import android.content.pm.PackageManager
import android.os.IBinder
import android.util.Log
import com.cflat.blink.shizuku.IShizukuShellService
import com.cflat.blink.shizuku.ShizukuShellUserService
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.withContext
import rikka.shizuku.Shizuku
import kotlinx.coroutines.suspendCancellableCoroutine

object ShizukuUtils {

    private const val TAG = "ShizukuUtils"
    private const val SHIZUKU_PERMISSION_CODE = 1001
    private const val USER_SERVICE_TAG = "deviceidle-whitelist"
    private const val USER_SERVICE_VERSION = 1
    private const val USER_SERVICE_PROCESS_SUFFIX = "blink-deviceidle"

    private val _shizukuStatus = MutableStateFlow(ShizukuStatus.UNKNOWN)
    val shizukuStatus: StateFlow<ShizukuStatus> = _shizukuStatus.asStateFlow()

    enum class ShizukuStatus {
        UNKNOWN,
        NOT_INSTALLED,
        NOT_RUNNING,
        PERMISSION_DENIED,
        READY,
        WHITELISTED
    }

    private val binderReceivedListener = Shizuku.OnBinderReceivedListener {
        Log.d(TAG, "Shizuku binder received")
        refreshPermissionStatus()
    }

    private val binderDeadListener = Shizuku.OnBinderDeadListener {
        Log.d(TAG, "Shizuku binder dead")
        _shizukuStatus.value = ShizukuStatus.NOT_RUNNING
    }

    private val permissionResultListener =
        Shizuku.OnRequestPermissionResultListener { requestCode, grantResult ->
            if (requestCode == SHIZUKU_PERMISSION_CODE) {
                if (grantResult == PackageManager.PERMISSION_GRANTED) {
                    _shizukuStatus.value = ShizukuStatus.READY
                } else {
                    _shizukuStatus.value = ShizukuStatus.PERMISSION_DENIED
                }
            }
        }

    fun initialize() {
        Shizuku.addBinderReceivedListener(binderReceivedListener)
        Shizuku.addBinderDeadListener(binderDeadListener)
        Shizuku.addRequestPermissionResultListener(permissionResultListener)

        if (Shizuku.pingBinder()) {
            refreshPermissionStatus()
        } else {
            _shizukuStatus.value = ShizukuStatus.NOT_RUNNING
        }
    }

    fun cleanup() {
        Shizuku.removeBinderReceivedListener(binderReceivedListener)
        Shizuku.removeBinderDeadListener(binderDeadListener)
        Shizuku.removeRequestPermissionResultListener(permissionResultListener)
    }

    private fun refreshPermissionStatus() {
        if (Shizuku.checkSelfPermission() == PackageManager.PERMISSION_GRANTED) {
            _shizukuStatus.value = ShizukuStatus.READY
        } else if (Shizuku.shouldShowRequestPermissionRationale()) {
            _shizukuStatus.value = ShizukuStatus.PERMISSION_DENIED
        } else {
            _shizukuStatus.value = ShizukuStatus.PERMISSION_DENIED
        }
    }

    fun requestPermission() {
        if (Shizuku.pingBinder()) {
            Shizuku.requestPermission(SHIZUKU_PERMISSION_CODE)
        } else {
            _shizukuStatus.value = ShizukuStatus.NOT_RUNNING
        }
    }

    suspend fun whitelistFromBatteryOptimization(context: Context): Result<Unit> =
        withContext(Dispatchers.IO) {
            try {
                val packageName = context.packageName
                val output = withUserService(context) { service ->
                    service.whitelistPackage(packageName)
                }
                Log.d(TAG, "Battery whitelist success: $output")
                _shizukuStatus.value = ShizukuStatus.WHITELISTED
                Result.success(Unit)
            } catch (e: Exception) {
                Log.e(TAG, "Shizuku whitelist command failed", e)
                Result.failure(e)
            }
        }

    suspend fun checkIfWhitelisted(context: Context): Boolean = withContext(Dispatchers.IO) {
        try {
            val packageName = context.packageName
            val output = withUserService(context) { service ->
                service.listDeviceIdleWhitelist()
            }
            output.contains(packageName)
        } catch (e: Exception) {
            Log.e(TAG, "Whitelist check failed", e)
            false
        }
    }

    private suspend fun withUserService(
        context: Context,
        operation: (IShizukuShellService) -> String
    ): String = suspendCancellableCoroutine { continuation ->
        val args = Shizuku.UserServiceArgs(
            ComponentName(context.packageName, ShizukuShellUserService::class.java.name)
        )
            .tag(USER_SERVICE_TAG)
            .version(USER_SERVICE_VERSION)
            .processNameSuffix(USER_SERVICE_PROCESS_SUFFIX)
            .daemon(false)
        var bound = false
        lateinit var connection: ServiceConnection

        fun unbind() {
            if (!bound) return
            bound = false
            runCatching { Shizuku.unbindUserService(args, connection, true) }
                .onFailure { Log.w(TAG, "Failed to unbind Shizuku UserService", it) }
        }

        connection = object : ServiceConnection {
            override fun onServiceConnected(name: ComponentName, service: IBinder) {
                try {
                    val shellService = IShizukuShellService.Stub.asInterface(service)
                    val output = operation(shellService)
                    if (continuation.isActive) continuation.resumeWith(Result.success(output))
                } catch (e: Exception) {
                    if (continuation.isActive) continuation.resumeWith(Result.failure(e))
                } finally {
                    unbind()
                }
            }

            override fun onServiceDisconnected(name: ComponentName) {
                if (continuation.isActive) {
                    continuation.resumeWith(
                        Result.failure(IllegalStateException("Shizuku UserService disconnected"))
                    )
                }
            }
        }

        continuation.invokeOnCancellation { unbind() }
        try {
            Shizuku.bindUserService(args, connection)
            bound = true
        } catch (e: Exception) {
            if (continuation.isActive) continuation.resumeWith(Result.failure(e))
        }
    }
}
