# Simple Go auto-reload script for Windows
$lastWrite = (Get-Date).AddDays(-1)

Write-Host "Watching Go files for changes..." -ForegroundColor Green
Write-Host "Press Ctrl+C to stop" -ForegroundColor Yellow

# Start initial build
go run .

while ($true) {
    $files = Get-ChildItem -Path . -Include *.go -Recurse | Where-Object { $_.LastWriteTime -gt $lastWrite }
    
    if ($files) {
        Write-Host "Changes detected, rebuilding..." -ForegroundColor Cyan
        $lastWrite = Get-Date
        
        # Kill existing process (if any)
        Get-Process -Name "main" -ErrorAction SilentlyContinue | Stop-Process -Force
        
        # Rebuild and run
        go run .
        
        Start-Sleep -Seconds 1
    }
    
    Start-Sleep -Milliseconds 500
}