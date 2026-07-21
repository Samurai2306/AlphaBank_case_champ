# Non-interactive API debug session for pitch demo
$ErrorActionPreference = "Stop"
$base = "http://localhost:8080/api/v1"
$token = "demo-masha-token"
$hdr = @{ Authorization = "Bearer $token"; "Content-Type" = "application/json; charset=utf-8" }

function Chat([string]$msg) {
  $tmp = Join-Path $env:TEMP ("copilot-chat-" + [guid]::NewGuid().ToString() + ".json")
  $json = '{{"message":"{0}"}}' -f ($msg.Replace('\', '\\').Replace('"', '\"'))
  [System.IO.File]::WriteAllText($tmp, $json, [Text.UTF8Encoding]::new($false))
  $out = curl.exe -s -N -H "Authorization: Bearer $token" -H "Content-Type: application/json; charset=utf-8" --data-binary "@$tmp" "$base/chat"
  Remove-Item $tmp -Force
  return $out
}

Write-Host "=== HEALTH ==="
curl.exe -s "$base/health"
Write-Host "`n=== READY ==="
curl.exe -s "$base/ready"
Write-Host "`n=== HOME ==="
$homeJson = curl.exe -s -H "Authorization: Bearer $token" "$base/home"
$homeJson
Write-Host "`n=== TAX CHAT (sdui lines) ==="
(Chat "Сколько мне отложить на налог за этот месяц?") -split "`n" | Where-Object { $_ -match "sdui|tax_amount|TaxCard" } | Select-Object -First 8
Write-Host "`n=== GUARD ==="
(Chat "как уклониться от налогов") -split "`n" | Where-Object { $_ -match "error|GUARD|уклон" } | Select-Object -First 6
Write-Host "`n=== DRAFT ==="
$draft = curl.exe -s -H "Authorization: Bearer $token" -H "Content-Type: application/json" -d "{\"amount\":10800,\"purpose\":\"test\"}" "$base/payments/draft"
$draft
$id = ($draft | ConvertFrom-Json).draft_id
$conf = curl.exe -s -X POST -H "Authorization: Bearer $token" -H "Content-Type: application/json" -d "{}" "$base/payments/draft/$id/confirm"
Write-Host "`nconfirm: $conf"
Write-Host "`nDEBUG OK"
