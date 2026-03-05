@echo off
SETLOCAL EnableDelayedExpansion

echo ===========================================
echo   DevAle v2 - Easy Setup & Build Script
echo ===========================================
echo.

:: Check for Go
where go >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo [ERROR] Go is not installed. Please install Go from https://go.dev/
    pause
    exit /b 1
)
echo [OK] Go found.

:: Check for Node.js
where npm >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo [ERROR] Node.js is not installed. Please install Node.js from https://nodejs.org/
    pause
    exit /b 1
)
echo [OK] Node.js found.

:: Check for Wails
where wails >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo [INFO] Wails CLI not found. Attempting to install...
    go install github.com/wailsapp/wails/v2/cmd/wails@v2.11.0
    if %ERRORLEVEL% neq 0 (
        echo [ERROR] Failed to install Wails. Please install it manually: go install github.com/wailsapp/wails/v2/cmd/wails@v2.11.0
        pause
        exit /b 1
    )
    for /f "delims=" %%G in ('go env GOPATH') do set "GOPATH_DIR=%%G"
    if exist "%GOPATH_DIR%\bin\wails.exe" set "PATH=%PATH%;%GOPATH_DIR%\bin"
    echo [OK] Wails installed successfully.
) else (
    echo [OK] Wails found.
)

echo.
echo [BUILD] Starting production build for Windows...
wails build -platform windows/amd64

if %ERRORLEVEL% equ 0 (
    echo.
    echo ===========================================
    echo   BUILD SUCCESSFUL!
    echo   Your app is ready at: build\bin\devale-v2.exe
    echo ===========================================
) else (
    echo.
    echo [ERROR] Build failed. Please check the logs above.
)

pause
exit /b %ERRORLEVEL%
