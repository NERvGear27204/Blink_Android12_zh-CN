package com.cflat.blink.shizuku

import androidx.annotation.Keep
import java.io.BufferedReader
import java.io.InputStreamReader

/** Runs the fixed device-idle operations as Shizuku's shell/root identity. */
@Keep
class ShizukuShellUserService : IShizukuShellService.Stub() {

    override fun whitelistPackage(packageName: String): String =
        runCommand("cmd", "deviceidle", "whitelist", "+$packageName")

    override fun listDeviceIdleWhitelist(): String =
        runCommand("cmd", "deviceidle", "whitelist")

    override fun destroy() {
        System.exit(0)
    }

    private fun runCommand(vararg command: String): String {
        val process = ProcessBuilder(*command)
            .redirectErrorStream(true)
            .start()
        val output = BufferedReader(InputStreamReader(process.inputStream)).use { it.readText() }
        val exitCode = process.waitFor()
        check(exitCode == 0) { "Command failed with exit code $exitCode: $output" }
        return output
    }
}
