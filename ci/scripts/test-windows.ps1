# Windows용 크로스 플랫폼 테스트 스크립트
param(
    [string]$OutputDir = "test-results",
    [switch]$Verbose = $false,
    [switch]$Benchmark = $false
)

# 색상 및 로깅 함수
function Write-ColorOutput($ForegroundColor) {
    $fc = $host.UI.RawUI.ForegroundColor
    $host.UI.RawUI.ForegroundColor = $ForegroundColor
    if ($args) {
        Write-Output $args
    } else {
        $input | Write-Output
    }
    $host.UI.RawUI.ForegroundColor = $fc
}

function Write-Info($message) {
    Write-ColorOutput Blue "[INFO] $message"
}

function Write-Success($message) {
    Write-ColorOutput Green "[SUCCESS] $message"
}

function Write-Error($message) {
    Write-ColorOutput Red "[ERROR] $message"
}

function Write-Warning($message) {
    Write-ColorOutput Yellow "[WARNING] $message"
}

# 결과 디렉터리 생성
$timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
$resultsPath = Join-Path $OutputDir $timestamp
New-Item -ItemType Directory -Path $resultsPath -Force | Out-Null
Write-Info "Test results will be saved to: $resultsPath"

# Go 설치 확인
function Test-GoInstallation {
    try {
        $goVersion = & go version 2>$null
        Write-Info "Go is installed: $goVersion"
        return $true
    }
    catch {
        Write-Error "Go is not installed or not in PATH"
        return $false
    }
}

# Docker 설치 확인 (Windows)
function Test-DockerInstallation {
    try {
        $dockerVersion = & docker --version 2>$null
        Write-Info "Docker is installed: $dockerVersion"
        return $true
    }
    catch {
        Write-Warning "Docker is not available"
        return $false
    }
}

# 네이티브 Windows 테스트
function Test-WindowsNative {
    Write-Info "Running native Windows tests..."
    
    $logFile = Join-Path $resultsPath "windows-native.log"
    
    try {
        # 기본 테스트
        $testOutput = & go test -v ./... 2>&1
        $testOutput | Out-File -FilePath $logFile
        
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Windows native tests passed"
        } else {
            Write-Error "Windows native tests failed"
            Get-Content $logFile | Select-Object -Last 20
        }
        
        # 벤치마크 (옵션)
        if ($Benchmark) {
            $benchFile = Join-Path $resultsPath "windows-native-benchmark.log"
            $benchOutput = & go test -bench=. -benchmem ./... 2>&1
            $benchOutput | Out-File -FilePath $benchFile
            Write-Info "Benchmark results saved to: $benchFile"
        }
    }
    catch {
        Write-Error "Failed to run Windows native tests: $_"
    }
}

# PowerShell Core 테스트
function Test-PowerShellCore {
    if (Get-Command pwsh -ErrorAction SilentlyContinue) {
        Write-Info "Testing with PowerShell Core..."
        
        $logFile = Join-Path $resultsPath "windows-pwsh.log"
        
        try {
            $testOutput = & pwsh -Command "go test -v ./..." 2>&1
            $testOutput | Out-File -FilePath $logFile
            
            if ($LASTEXITCODE -eq 0) {
                Write-Success "PowerShell Core tests passed"
            } else {
                Write-Error "PowerShell Core tests failed"
            }
        }
        catch {
            Write-Error "Failed to run PowerShell Core tests: $_"
        }
    } else {
        Write-Warning "PowerShell Core (pwsh) not found, skipping"
    }
}

# WSL 테스트 (Windows 10/11에서 사용 가능)
function Test-WSL {
    if (Get-Command wsl -ErrorAction SilentlyContinue) {
        Write-Info "Testing with WSL..."
        
        # 사용 가능한 WSL 배포판 확인
        $wslDistros = & wsl -l -q 2>$null | Where-Object { $_ -ne "" }
        
        foreach ($distro in $wslDistros) {
            Write-Info "Testing WSL distro: $distro"
            $logFile = Join-Path $resultsPath "wsl-$distro.log"
            
            try {
                # WSL에서 Go 설치 확인 및 테스트 실행
                $wslCommand = @"
cd /mnt/c/$(($pwd.Path -replace '\\', '/') -replace 'C:', '')
if command -v go >/dev/null 2>&1; then
    go version
    go mod download
    go test -v ./...
else
    echo "Go not installed in WSL distro: $distro"
    exit 1
fi
"@
                
                $testOutput = & wsl -d $distro bash -c $wslCommand 2>&1
                $testOutput | Out-File -FilePath $logFile
                
                if ($LASTEXITCODE -eq 0) {
                    Write-Success "WSL $distro tests passed"
                } else {
                    Write-Warning "WSL $distro tests failed (Go may not be installed)"
                }
            }
            catch {
                Write-Error "Failed to run WSL tests for $distro : $_"
            }
        }
    } else {
        Write-Warning "WSL not available"
    }
}

# Docker Windows Container 테스트
function Test-DockerWindows {
    if (Test-DockerInstallation) {
        Write-Info "Testing with Docker Windows containers..."
        
        # Windows Server Core 기반 테스트
        $dockerFile = @"
FROM mcr.microsoft.com/windows/servercore:ltsc2022
SHELL ["powershell", "-Command"]

# Go 설치
RUN Invoke-WebRequest -Uri 'https://golang.org/dl/go1.21.0.windows-amd64.zip' -OutFile 'go.zip'; \
    Expand-Archive -Path 'go.zip' -DestinationPath 'C:\'; \
    Remove-Item 'go.zip'

ENV PATH="C:\go\bin;${PATH}"

WORKDIR C:\app
COPY . .

CMD ["powershell", "-Command", "go version; go mod download; go test -v ./..."]
"@
        
        $dockerFile | Out-File -FilePath "Dockerfile.windows" -Encoding utf8
        
        try {
            # Docker 이미지 빌드
            Write-Info "Building Windows Docker image..."
            & docker build -f Dockerfile.windows -t strutil-test-windows .
            
            if ($LASTEXITCODE -eq 0) {
                Write-Info "Running Windows Docker container tests..."
                $logFile = Join-Path $resultsPath "docker-windows.log"
                $testOutput = & docker run --rm strutil-test-windows 2>&1
                $testOutput | Out-File -FilePath $logFile
                
                if ($LASTEXITCODE -eq 0) {
                    Write-Success "Docker Windows tests passed"
                } else {
                    Write-Error "Docker Windows tests failed"
                }
            } else {
                Write-Error "Failed to build Windows Docker image"
            }
        }
        catch {
            Write-Error "Docker Windows tests failed: $_"
        }
        finally {
            Remove-Item "Dockerfile.windows" -ErrorAction SilentlyContinue
        }
    }
}

# 다중 아키텍처 테스트 (Windows)
function Test-MultiArch {
    Write-Info "Testing multiple architectures on Windows..."
    
    $architectures = @("amd64", "386")
    
    foreach ($arch in $architectures) {
        Write-Info "Testing architecture: $arch"
        $logFile = Join-Path $resultsPath "windows-$arch.log"
        
        try {
            $env:GOARCH = $arch
            $testOutput = & go test -v ./... 2>&1
            $testOutput | Out-File -FilePath $logFile
            
            if ($LASTEXITCODE -eq 0) {
                Write-Success "Windows $arch tests passed"
            } else {
                Write-Error "Windows $arch tests failed"
            }
        }
        catch {
            Write-Error "Failed to test architecture $arch : $_"
        }
        finally {
            Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
        }
    }
}

# 성능 리포트 생성
function Generate-PerformanceReport {
    Write-Info "Generating performance report..."
    
    $reportFile = Join-Path $resultsPath "performance-summary.md"
    $report = @"
# Windows Performance Test Summary

Generated on: $(Get-Date)

## Test Results

"@
    
    Get-ChildItem $resultsPath -Filter "*.log" | ForEach-Object {
        $platform = $_.BaseName
        $report += @"

### $platform
``````
"@
        
        # 벤치마크 결과 추출
        $content = Get-Content $_.FullName | Where-Object { $_ -match "Benchmark|PASS|FAIL" } | Select-Object -First 20
        if ($content) {
            $report += ($content -join "`n")
        } else {
            $report += "No benchmark data found"
        }
        
        $report += @"
``````

"@
    }
    
    $report | Out-File -FilePath $reportFile -Encoding utf8
    Write-Success "Performance report saved to: $reportFile"
}

# 시스템 정보 수집
function Collect-SystemInfo {
    Write-Info "Collecting system information..."
    
    $sysInfoFile = Join-Path $resultsPath "system-info.txt"
    
    $systemInfo = @"
System Information
==================

OS Version: $([Environment]::OSVersion.VersionString)
OS Platform: $([Environment]::OSVersion.Platform)
Machine Name: $([Environment]::MachineName)
User Name: $([Environment]::UserName)
Processor Count: $([Environment]::ProcessorCount)
CLR Version: $([Environment]::Version)
Working Set: $([Math]::Round([Environment]::WorkingSet / 1MB, 2)) MB

PowerShell Version: $($PSVersionTable.PSVersion)
Go Version: $(try { & go version } catch { "Not installed" })
Docker Version: $(try { & docker --version } catch { "Not installed" })

Current Directory: $(Get-Location)
"@
    
    $systemInfo | Out-File -FilePath $sysInfoFile -Encoding utf8
    Write-Info "System information saved to: $sysInfoFile"
}

# 메인 실행 함수
function Main {
    Write-Info "Starting Windows cross-platform tests..."
    
    # 전제 조건 확인
    if (-not (Test-GoInstallation)) {
        Write-Error "Go is required to run tests"
        exit 1
    }
    
    # 시스템 정보 수집
    Collect-SystemInfo
    
    # 각종 테스트 실행
    Test-WindowsNative
    Test-PowerShellCore
    Test-WSL
    Test-MultiArch
    
    # Docker 테스트 (옵션)
    if (Test-DockerInstallation) {
        Test-DockerWindows
    }
    
    # 성능 리포트 생성
    Generate-PerformanceReport
    
    Write-Success "All Windows tests completed!"
    Write-Info "Results available in: $resultsPath"
    
    # 결과 폴더 열기 (옵션)
    if ((Read-Host "Open results folder? (y/N)") -eq "y") {
        explorer $resultsPath
    }
}

# 스크립트 실행
Main