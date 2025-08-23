@echo off
REM Windows batch script for building go-strutil
REM Usage: build_windows.bat [test|bench|build|all]

setlocal enabledelayedexpansion

echo Starting Windows build for go-strutil...

REM Check if Go is installed
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo Error: Go is not installed or not in PATH
    exit /b 1
)

REM Get Go version
for /f "tokens=3" %%v in ('go version') do set GO_VERSION=%%v
set GO_VERSION=%GO_VERSION:go=%
echo Detected Go version: %GO_VERSION%

REM Create directories
if not exist "build" mkdir build
if not exist "reports" mkdir reports
if not exist "reports\benchmarks" mkdir reports\benchmarks
if not exist "reports\coverage" mkdir reports\coverage

REM Set build tags based on Go version
set BUILD_TAGS=
for /f "tokens=1,2 delims=." %%a in ("%GO_VERSION%") do (
    set MAJOR=%%a
    set MINOR=%%b
)

REM Simple version comparison for Windows
if %MAJOR% gtr 1 (
    set BUILD_TAGS=generics
) else (
    if %MINOR% geq 18 (
        set BUILD_TAGS=generics
    ) else (
        if %MINOR% lss 7 (
            set BUILD_TAGS=legacy
        )
    )
)

echo Using build tags: %BUILD_TAGS%

REM Parse command line argument
set COMMAND=%1
if "%COMMAND%"=="" set COMMAND=all

if "%COMMAND%"=="test" goto test
if "%COMMAND%"=="bench" goto bench
if "%COMMAND%"=="build" goto build
if "%COMMAND%"=="fmt" goto fmt
if "%COMMAND%"=="vet" goto vet
if "%COMMAND%"=="all" goto all

:fmt
echo Formatting Go code...
go fmt ./...
if %errorlevel% neq 0 (
    echo Error during formatting
    exit /b 1
)
echo Code formatted successfully
if "%COMMAND%"=="fmt" goto end
goto :eof

:vet
echo Running go vet...
go vet ./...
if %errorlevel% neq 0 (
    echo Error during vet check
    exit /b 1
)
echo Vet check passed
if "%COMMAND%"=="vet" goto end
goto :eof

:build
echo Building project...
if "%BUILD_TAGS%"=="" (
    go build -v ./...
) else (
    go build -tags="%BUILD_TAGS%" -v ./...
)
if %errorlevel% neq 0 (
    echo Error during build
    exit /b 1
)
echo Build completed successfully
if "%COMMAND%"=="build" goto end
goto :eof

:test
echo Running tests...
if "%BUILD_TAGS%"=="" (
    go test -v ./...
) else (
    go test -tags="%BUILD_TAGS%" -v ./...
)
if %errorlevel% neq 0 (
    echo Error during testing
    exit /b 1
)
echo All tests passed
if "%COMMAND%"=="test" goto end
goto :eof

:bench
echo Running benchmarks...
if "%BUILD_TAGS%"=="" (
    go test -bench=. -benchmem ./... > reports\benchmarks\bench-windows.txt
) else (
    go test -tags="%BUILD_TAGS%" -bench=. -benchmem ./... > reports\benchmarks\bench-windows.txt
)
if %errorlevel% neq 0 (
    echo Error during benchmarking
    exit /b 1
)
echo Benchmarks completed: reports\benchmarks\bench-windows.txt
if "%COMMAND%"=="bench" goto end
goto :eof

:all
echo Running all build tasks...
call :fmt
call :vet  
call :build
call :test
echo All build tasks completed successfully
goto :eof

:end
echo Done.