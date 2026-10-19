package com.cflat.blink.shizuku;

interface IShizukuShellService {
    // Transaction reserved by Shizuku to terminate the UserService process.
    void destroy() = 16777114;

    String whitelistPackage(String packageName) = 1;
    String listDeviceIdleWhitelist() = 2;
}
