# PowerShell install script for Gitti

$ErrorActionPreference = "Stop"

# Version to install
$Version = "v0.1.3"

Write-Host "Installing Gitti version: $Version" -ForegroundColor Cyan

# Detect Architecture
if ($env:PROCESSOR_ARCHITECTURE -eq "AMD64") {
    $Arch = "amd64"
} elseif ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") {
    $Arch = "arm64"
} else {
    Write-Error "Unsupported architecture: $env:PROCESSOR_ARCHITECTURE"
    exit 1
}

Write-Host "Detected Architecture: $Arch" -ForegroundColor Gray

# Construct URL
# Naming convention: gitti-v{version}-windows-{arch}.zip
$VersionNum = $Version -replace "^v", ""
$FileName = "gitti-v$VersionNum-windows-$Arch.zip"
$DownloadUrl = "https://github.com/gohyuhan/gitti/releases/download/$Version/$FileName"

Write-Host "Download URL: $DownloadUrl" -ForegroundColor Gray

# Temp paths
$TempDir = [System.IO.Path]::GetTempPath()
$ZipPath = Join-Path $TempDir $FileName

# Download
Write-Host "Downloading..." -ForegroundColor Cyan
try {
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $ZipPath
} catch {
    Write-Error "Failed to download: $_"
    exit 1
}

# Install Directory
$InstallDir = Join-Path $env:LOCALAPPDATA "gitti"
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
}

# Extract
Write-Host "Extracting..." -ForegroundColor Cyan
# Using Expand-Archive. Force to overwrite existing files.
Expand-Archive -Path $ZipPath -DestinationPath $InstallDir -Force

# Cleanup
Remove-Item $ZipPath -ErrorAction SilentlyContinue

# Verify Binary
$BinaryPath = Join-Path $InstallDir "gitti.exe"
if (-not (Test-Path $BinaryPath)) {
    # Check if it's in a subfolder?
    $Found = Get-ChildItem -Path $InstallDir -Filter "gitti.exe" -Recurse | Select-Object -First 1
    if ($Found) {
        Move-Item $Found.FullName $InstallDir -Force
        $BinaryPath = Join-Path $InstallDir "gitti.exe"
    } else {
        Write-Error "Binary 'gitti.exe' not found in extracted files."
        exit 1
    }
}

Write-Host "Installed to: $BinaryPath" -ForegroundColor Green

# Add to PATH
$UserPath = [Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::User)
if ($UserPath -notlike "*$InstallDir*") {
    Write-Host "Adding to PATH..." -ForegroundColor Cyan
    $NewPath = "$UserPath;$InstallDir"
    [Environment]::SetEnvironmentVariable("Path", $NewPath, [EnvironmentVariableTarget]::User)
    $env:Path = "$env:Path;$InstallDir" # Update current session
    Write-Host "Added to PATH. You may need to restart your terminal." -ForegroundColor Yellow
} else {
    Write-Host "Already in PATH." -ForegroundColor Gray
}

Write-Host "Installation complete! Run 'gitti --version' to verify." -ForegroundColor Green
