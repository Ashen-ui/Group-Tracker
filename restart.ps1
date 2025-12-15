# Script pour arrêter le processus sur le port 8081 et relancer le serveur

Write-Host "Recherche du processus sur le port 8081..."

$port = 8081
$process = Get-NetTCPConnection -LocalPort $port -ErrorAction SilentlyContinue | Select-Object -ExpandProperty OwningProcess -Unique

if ($process) {
    Write-Host "Arrêt du processus PID: $process"
    Stop-Process -Id $process -Force -ErrorAction SilentlyContinue
    Start-Sleep -Seconds 1
}

Write-Host "Démarrage du serveur..."
go run main.go

