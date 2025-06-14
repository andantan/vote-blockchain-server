while ($true) {
    Clear-Host
    Write-Host "현재 시간: $(Get-Date -Format 'HH:mm:ss')"
    Write-Host "==================================="
    Write-Host "프로세스 이름         메모리 사용량 (MB)"
    Write-Host "==================================="

    tasklist | Select-String -Pattern "^blockchain-node" -CaseSensitive:$false | ForEach-Object {
        $line = $_.ToString()
        if ($line -match '(\d{1,3}(,\d{3})*) K') {
            $kbMemory = $matches[1] -replace ',', ''
            $mbMemory = [math]::Round([double]$kbMemory / 1024, 2)
            $processName = ($line -split '\s+')[0]

            Write-Host ("{0,-20} {1,10} MB" -f $processName, $mbMemory)
        } else {
            Write-Host $line
        }
    }
    Write-Host "==================================="
    Start-Sleep -Seconds 3
}