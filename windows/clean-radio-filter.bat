@echo off
setlocal
set "SCRIPT_DIR=%~dp0"
set "BIN=%SCRIPT_DIR%radiofx.exe"

if not exist "%BIN%" set "BIN=%SCRIPT_DIR%..\radiofx.exe"

if not exist "%BIN%" (
    echo Error: radiofx.exe not found next to this wrapper or in its parent folder.
    pause
    exit /b 1
)

"%BIN%" apply -preset clean_radio -mode filter -suffix clean_radio %*
if errorlevel 1 (
    pause
    exit /b 1
)
