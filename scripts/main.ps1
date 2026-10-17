# The same extension in PowerShell (Windows twin of scripts/main.sh).
# Manifest for this variant: runtime = powershell, run = scripts/main.ps1, args = -NoProfile -File
# Protocol: https://rhun.app/extensions/v1/
$utf8 = New-Object System.Text.UTF8Encoding($false)
$in = [Console]::OpenStandardInput()
$out = [Console]::OpenStandardOutput()
$reader = New-Object System.IO.StreamReader($in, $utf8)
function Send($s) { $b = $utf8.GetBytes($s); $out.Write($b, 0, $b.Length); $out.Flush() }
$id = 0
function Request($line, $body) {
    $script:id++
    Send "$script:id $line`n"
    if ($body) { Send $body }
}
Send "hello 1`n"
while ($null -ne ($line = $reader.ReadLine())) {
    $f = $line -split ' '
    if ($f[0] -eq '0' -and $f[1] -eq 'run') {
        if ($f[2] -eq 'hello.say') { Request 'toast 22' 'Hello from PowerShell' }
    } elseif ($f[0] -eq '0' -and $f[1] -eq 'bye') { break }
}
