param(

    [Alias("Target")]
    [ValidateSet("run","build","test","help")]
    [string]$Action = "run",


    [string]$Port = "8080"
)

$ErrorActionPreference = "Stop"

function Show-Help {
    Write-Host "Usage:"
    Write-Host "  .\build.ps1 -Action run   [-Port 8080]"
    Write-Host "  .\build.ps1 -Action build"
    Write-Host "  .\build.ps1 -Action test"
    Write-Host "  .\build.ps1 -Action help"
    Write-Host ""
    Write-Host "Aliases:"
    Write-Host "  -Target можно использовать вместо -Action (напр. -Target test)"
}

function Ensure-BinDir {
    if (-not (Test-Path ".\bin")) {
        New-Item -ItemType Directory -Path ".\bin" | Out-Null
    }
}


try {
    $goVersion = & go version
} catch {
    Write-Error "Go не найден в PATH. Установи Go и перезапусти терминал."
    exit 1
}

switch ($Action) {

    "help" {
        Show-Help
        exit 0
    }

    "run" {
        Write-Host " Running server on PORT=$Port"
        $env:PORT = $Port
        # Важно: запуск только сервера
        & go run ./cmd/server
        exit $LASTEXITCODE
    }

    "build" {
        Write-Host " Building server..."
        Ensure-BinDir
        & go build -o .\bin\server.exe .\cmd\server
        if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
        Write-Host " Build complete: .\bin\server.exe"
        exit 0
    }

    "test" {
        Write-Host " Running tests..."
        # Вербозный прогон всех пакетов
        & go test -v ./...
        exit $LASTEXITCODE
    }

    default {
        Show-Help
        exit 1
    }
}
