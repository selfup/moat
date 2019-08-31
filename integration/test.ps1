go install

if (Test-Path archive) {
    Remove-Item -Recurse -Force archive
}

if (Test-Path fixtures) {
    Remove-Item -Recurse -Force fixtures
}

$WOW = "wow this is going to be encrypted and saved to a cloud service directory"

moat -home="archive" -service="fixtures"

$WOW | Set-Content archive/Moat/wow.txt

moat -home="archive" -service="fixtures" -cmd=push

Remove-Item archive\Moat\wow.txt

moat -home="archive" -service="fixtures" -cmd=pull

$contents = Get-Content archive/Moat/wow.txt

if ($contents -eq $WOW) {
    Write-Host "moat passed"
} 
else {
    Write-Host "moat failed"
    exit 1
}

Remove-Item -Recurse -Force archive
Remove-Item -Recurse -Force fixtures
