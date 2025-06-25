@echo off
chcp 65001
setlocal enabledelayedexpansion

REM 현재 날짜를 가져옵니다.
for /f "tokens=2 delims==" %%I in ('"wmic os get localdatetime /value"') do set datetime=%%I
set today=%datetime:~0,8%

REM 오늘이 월요일인지 확인합니다.
for /f "tokens=*" %%A in ('powershell -command "(Get-Date).DayOfWeek"') do set dayofweek=%%A

if /i "%dayofweek%"=="Monday" (
    REM 오늘이 월요일이면 3일 전 날짜를 계산합니다.
    for /f "tokens=*" %%A in ('powershell -command "(Get-Date).AddDays(-3).ToString(\"yyyyMMdd\")"') do set referenceDate=%%A
) else (
    REM 오늘이 월요일이 아니면 어제 날짜를 참조합니다.
    for /f "tokens=*" %%A in ('powershell -command "(Get-Date).AddDays(-1).ToString(\"yyyyMMdd\")"') do set referenceDate=%%A
)

echo 참조 날짜: %referenceDate%
echo 오늘 날짜: %today%

for %%F in (*.xlsx *.docx) do (
    set filename=%%~nF

    if not exist "2025" mkdir "2025"
    copy "%%F" "2025\"
    echo 파일을 2025 폴더로 복사: "%%~nF%%~xF"
)


echo 파일 이름 변경 작업을 시작합니다...
for %%F in (*.xlsx *.docx) do (
    set filename=%%~nF
    set extension=%%~xF

    REM 파일 이름에 참조 날짜가 포함되어 있는지 확인
    echo !filename! | findstr /c:"%referenceDate%" >nul
    if not errorlevel 1 (
        REM 파일 이름에서 참조 날짜를 오늘 날짜로 변경
        set newfilename=!filename:%referenceDate%=%today%!
        ren "%%F" "!newfilename!!extension!"
        echo 파일 이름 변경: "%%~nF%%~xF" → "!newfilename!!extension!"
    )
)

echo 파일 이름 변경 작업이 완료되었습니다.
pause
