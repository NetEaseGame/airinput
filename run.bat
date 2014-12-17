@echo off

echo Android AirInput -- Initializing...

echo Waiting for device to be connected...
adb.exe wait-for-device

echo - Installing native service...
adb.exe push air-native "/data/local/tmp/air-native"
adb.exe shell chmod 755 "/data/local/tmp/air-native"

echo Starting...
adb.exe shell kill -9 "/data/local/tmp/air-native"
start /B adb.exe shell "/data/local/tmp/air-native -daemon"

echo Service started successfully.
