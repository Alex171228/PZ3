param(
  [string]$Port = "8080"
)

Write-Host "Running on PORT=$Port"
$env:PORT = $Port
go run ./cmd/server
