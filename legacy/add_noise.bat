@echo off
setlocal EnableDelayedExpansion

set "DIR_START=start"
set "DIR_END=end"
set "VOL_NOISE=0.3"

:loop
if "%~1"=="" goto end

call :GetRandomFile "%DIR_START%" FILE_START
call :GetRandomFile "%DIR_END%" FILE_END

if not defined FILE_START (
    echo Error: No wav files in %DIR_START%
    pause
    goto end
)
if not defined FILE_END (
    echo Error: No wav files in %DIR_END%
    pause
    goto end
)

ffmpeg -y -i "!FILE_START!" -i "%~1" -i "!FILE_END!" -filter_complex "[0:a]volume=%VOL_NOISE%,aformat=sample_rates=44100:channel_layouts=mono[a0];[1:a]highpass=f=300,lowpass=f=3000,acrusher=bits=8:mode=log:aa=1,aformat=sample_rates=44100:channel_layouts=mono[a1];[2:a]volume=%VOL_NOISE%,aformat=sample_rates=44100:channel_layouts=mono[a2];[a0][a1][a2]concat=n=3:v=0:a=1[out]" -map "[out]" "%~dpn1_w_noise.wav"

if !errorlevel! neq 0 (
    echo Error processing: "%~1"
    pause
)
shift
goto loop

:end
exit /b

:GetRandomFile
set "cnt=0"
for %%f in ("%~1\*.wav") do set /a cnt+=1
if !cnt!==0 (
    set "%~2="
    exit /b
)
set /a r=(!RANDOM! %% cnt) + 1
set "i=0"
for %%f in ("%~1\*.wav") do (
    set /a i+=1
    if !i!==!r! set "%~2=%%f"
)
exit /b
