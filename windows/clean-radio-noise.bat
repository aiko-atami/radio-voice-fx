@echo off
setlocal
set "ROOT=%~dp0.."
set "BIN=%ROOT%\radiofx.exe"

if not exist "%BIN%" (
    echo Error: "%BIN%" not found.
    echo Build or copy radiofx.exe into this folder first.
    pause
    exit /b 1
)

"%BIN%" apply -preset clean_radio -mode noise -suffix clean_noise %*
if errorlevel 1 (
    pause
    exit /b 1
)
