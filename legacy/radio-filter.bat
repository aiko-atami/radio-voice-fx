@echo off
:loop
if "%~1"=="" goto end
ffmpeg -i "%~1" -af "highpass=f=300,lowpass=f=3000,aformat=sample_rates=11025,compand,volume=1.2" "%~dpn1_filtered.wav"
if %errorlevel% neq 0 (
    echo Error processing: %~1
    pause
)
shift
goto loop
:end
