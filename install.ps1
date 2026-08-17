# speedtest_cli - Windows PowerShell Installer
# Usage: iwr -useb https://raw.githubusercontent.com/Dolyyyy/speedtest_cli/main/install.ps1 | iex

$ErrorActionPreference = "Stop"

Write-Host "🚀 Installing speedtest_cli for Windows..." -ForegroundColor Cyan

$Url = "https://raw.githubusercontent.com/Dolyyyy/speedtest_cli/main/binaries/speedtest-windows-amd64.exe"
$InstallDir = "$env:LOCALAPPDATA\Programs\speedtest"
$TargetExe = "$InstallDir\speedtest.exe"

if (-not (Test-Path -Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
}

Write-Host "📥 Downloading speedtest-windows-amd64.exe..." -ForegroundColor Yellow
Invoke-WebRequest -Uri $Url -OutFile $TargetExe

# Check if Path contains InstallDir
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    Write-Host "🔧 Adding $InstallDir to User PATH..." -ForegroundColor Magenta
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
    $env:Path += ";$InstallDir"
}

Write-Host "`n✅ speedtest_cli installed successfully!" -ForegroundColor Green
Write-Host "Please restart your PowerShell terminal and type 'speedtest' to run." -ForegroundColor White
