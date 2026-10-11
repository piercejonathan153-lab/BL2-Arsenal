# BL2 Arsenal - installer, updater and launcher for The Pre-Sequel
# Runs on Windows PowerShell 5.1 (built into Windows 10/11); nothing else to install.
$ErrorActionPreference = 'Stop'
Add-Type -AssemblyName System.Windows.Forms, System.Drawing, System.IO.Compression.FileSystem
[System.Windows.Forms.Application]::EnableVisualStyles()
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$LauncherVersion = '1.0.0'
$Owner = 'piercejonathan153-lab'
$Repo = 'BL2-Arsenal'
$TpsAppId = 261640
$ModFolder = 'bl2_arsenal'
$CfgDir = Join-Path $env:APPDATA 'BL2Arsenal'
$CfgFile = Join-Path $CfgDir 'launcher.json'

# ------------------------------------------------------------------ helpers
function Load-Config {
    if (Test-Path $CfgFile) { try { return Get-Content $CfgFile -Raw | ConvertFrom-Json } catch {} }
    return [pscustomobject]@{ TpsPath = $null }
}
function Save-Config($cfg) {
    New-Item -ItemType Directory -Force -Path $CfgDir | Out-Null
    $cfg | ConvertTo-Json | Set-Content -Path $CfgFile -Encoding UTF8
}
function Get-SteamLibraries {
    $libs = @()
    foreach ($key in 'HKCU:\Software\Valve\Steam', 'HKLM:\SOFTWARE\WOW6432Node\Valve\Steam', 'HKLM:\SOFTWARE\Valve\Steam') {
        try {
            $p = (Get-ItemProperty $key -ErrorAction Stop)
            foreach ($v in $p.SteamPath, $p.InstallPath) { if ($v) { $libs += ($v -replace '/', '\') } }
        } catch {}
    }
    foreach ($root in @($libs)) {
        $vdf = Join-Path $root 'steamapps\libraryfolders.vdf'
        if (Test-Path $vdf) {
            foreach ($m in [regex]::Matches((Get-Content $vdf -Raw), '"path"\s+"([^"]+)"')) {
                $libs += ($m.Groups[1].Value -replace '\\\\', '\')
            }
        }
    }
    return $libs | Where-Object { $_ } | Select-Object -Unique
}
function Find-Game($folder, $exe) {
    foreach ($lib in Get-SteamLibraries) {
        $p = Join-Path $lib "steamapps\common\$folder"
        if (Test-Path (Join-Path $p "Binaries\Win32\$exe")) { return $p }
    }
    return $null
}
function Is-Tps($p) { return $p -and (Test-Path (Join-Path $p 'Binaries\Win32\BorderlandsPreSequel.exe')) }
function Parse-Version($s) {
    if (-not $s) { return $null }
    $m = [regex]::Match($s, '(\d+)\.(\d+)\.(\d+)')
    if (-not $m.Success) { return $null }
    return [version]("{0}.{1}.{2}" -f $m.Groups[1].Value, $m.Groups[2].Value, $m.Groups[3].Value)
}
function Get-InstalledVersion($tps) {
    $toml = Join-Path $tps "sdk_mods\$ModFolder\pyproject.toml"
    if (Test-Path $toml) {
        $m = [regex]::Match((Get-Content $toml -Raw), 'version\s*=\s*"([^"]+)"')
        if ($m.Success) { return $m.Groups[1].Value }
    }
    $sd = Get-ChildItem (Join-Path $tps 'sdk_mods') -Filter "$ModFolder*.sdkmod" -ErrorAction SilentlyContinue | Select-Object -First 1
    if ($sd) { return (Parse-Version $sd.Name) }
    return $null
}
function Game-Running { return [bool](Get-Process -Name 'BorderlandsPreSequel' -ErrorAction SilentlyContinue) }

# ------------------------------------------------------------------ UI
$cfg = Load-Config
if (-not (Is-Tps $cfg.TpsPath)) {
    $here = Split-Path -Parent $MyInvocation.MyCommand.Path
    foreach ($c in @($here, (Split-Path $here -Parent), (Find-Game 'BorderlandsPreSequel' 'BorderlandsPreSequel.exe'))) {
        if (Is-Tps $c) { $cfg.TpsPath = $c; break }
    }
}
$bl2Path = Find-Game 'Borderlands 2' 'Borderlands2.exe'
$script:latest = $null

$bg = [Drawing.Color]::FromArgb(24, 26, 31); $panel = [Drawing.Color]::FromArgb(34, 37, 44)
$fg = [Drawing.Color]::FromArgb(230, 232, 236); $muted = [Drawing.Color]::FromArgb(150, 156, 168)
$accent = [Drawing.Color]::FromArgb(242, 169, 0); $good = [Drawing.Color]::FromArgb(110, 200, 120); $bad = [Drawing.Color]::FromArgb(235, 100, 90)

$form = New-Object Windows.Forms.Form
$form.Text = "BL2 Arsenal Launcher"
$form.Size = New-Object Drawing.Size(640, 560)
$form.StartPosition = 'CenterScreen'; $form.FormBorderStyle = 'FixedSingle'; $form.MaximizeBox = $false
$form.BackColor = $bg; $form.ForeColor = $fg
$form.Font = New-Object Drawing.Font('Segoe UI', 9.5)

$title = New-Object Windows.Forms.Label
$title.Text = 'BL2 ARSENAL'; $title.Font = New-Object Drawing.Font('Segoe UI', 20, [Drawing.FontStyle]::Bold)
$title.ForeColor = $accent; $title.AutoSize = $true; $title.Location = New-Object Drawing.Point(20, 14)
$form.Controls.Add($title)
$sub = New-Object Windows.Forms.Label
$sub.Text = "Borderlands 2 gear for The Pre-Sequel  -  launcher v$LauncherVersion"
$sub.ForeColor = $muted; $sub.AutoSize = $true; $sub.Location = New-Object Drawing.Point(23, 54)
$form.Controls.Add($sub)

function New-Row($y, $name) {
    $l = New-Object Windows.Forms.Label; $l.Text = $name; $l.ForeColor = $muted
    $l.Location = New-Object Drawing.Point(24, $y); $l.Size = New-Object Drawing.Size(150, 20); $form.Controls.Add($l)
    $v = New-Object Windows.Forms.Label; $v.Location = New-Object Drawing.Point(175, $y)
    $v.Size = New-Object Drawing.Size(430, 20); $v.AutoEllipsis = $true; $form.Controls.Add($v)
    return $v
}
$vTps = New-Row 90 'Pre-Sequel'
$vSdk = New-Row 114 'Mod manager'
$vBl2 = New-Row 138 'Borderlands 2'
$vInst = New-Row 162 'Installed version'
$vLatest = New-Row 186 'Latest version'

function New-Btn($text, $x, $y, $w, $primary) {
    $b = New-Object Windows.Forms.Button; $b.Text = $text
    $b.Location = New-Object Drawing.Point($x, $y); $b.Size = New-Object Drawing.Size($w, 36)
    $b.FlatStyle = 'Flat'; $b.FlatAppearance.BorderSize = 0; $b.Cursor = 'Hand'
    if ($primary) { $b.BackColor = $accent; $b.ForeColor = [Drawing.Color]::Black; $b.Font = New-Object Drawing.Font('Segoe UI', 10, [Drawing.FontStyle]::Bold) }
    else { $b.BackColor = $panel; $b.ForeColor = $fg }
    $form.Controls.Add($b); return $b
}
$bPlay = New-Btn 'PLAY' 24 222 140 $true
$bCheck = New-Btn 'Check for updates' 172 222 140 $false
$bInstall = New-Btn 'Install / Update' 320 222 140 $false
$bFile = New-Btn 'Install from file...' 468 222 140 $false
$bFolder = New-Btn 'Choose game folder...' 24 264 186 $false
$bOpen = New-Btn 'Open mod folder' 218 264 186 $false
$bUninstall = New-Btn 'Uninstall' 412 264 196 $false

$log = New-Object Windows.Forms.TextBox
$log.Multiline = $true; $log.ReadOnly = $true; $log.ScrollBars = 'Vertical'
$log.Location = New-Object Drawing.Point(24, 312); $log.Size = New-Object Drawing.Size(584, 190)
$log.BackColor = $panel; $log.ForeColor = $fg; $log.BorderStyle = 'None'; $log.Font = New-Object Drawing.Font('Consolas', 9)
$form.Controls.Add($log)

function Say($msg) { $log.AppendText("$msg`r`n"); [Windows.Forms.Application]::DoEvents() }
function Set-Status($label, $text, $ok) {
    $label.Text = $text
    if ($ok -eq $true) { $label.ForeColor = $good } elseif ($ok -eq $false) { $label.ForeColor = $bad } else { $label.ForeColor = $fg }
}

function Refresh-Status {
    $tps = $cfg.TpsPath
    if (Is-Tps $tps) { Set-Status $vTps $tps $true } else { Set-Status $vTps 'Not found - use "Choose game folder..."' $false }
    $hasSdk = (Is-Tps $tps) -and (Test-Path (Join-Path $tps 'Binaries\Win32\Plugins\pyunrealsdk.dll')) -and (Test-Path (Join-Path $tps 'sdk_mods\mods_base.sdkmod'))
    if ($hasSdk) { Set-Status $vSdk 'willow2-mod-manager installed' $true }
    else { Set-Status $vSdk 'Missing - get it from github.com/bl-sdk/willow2-mod-manager/releases' $false }
    if ($bl2Path) { Set-Status $vBl2 "$bl2Path (DLC gear unlocks from this install)" $true }
    else { Set-Status $vBl2 'Not found - the mod needs Borderlands 2 installed through Steam' $false }
    $iv = if (Is-Tps $tps) { Get-InstalledVersion $tps } else { $null }
    if ($iv) { Set-Status $vInst "v$iv" $null } else { Set-Status $vInst 'Not installed' $null }
    $bPlay.Enabled = (Is-Tps $tps); $bInstall.Enabled = (Is-Tps $tps); $bFile.Enabled = (Is-Tps $tps)
    $bUninstall.Enabled = [bool]$iv; $bOpen.Enabled = (Is-Tps $tps)
    if ($script:latest) {
        $lv = Parse-Version $script:latest.tag_name; $cur = Parse-Version $iv
        if ($cur -and $lv -le $cur) { Set-Status $vLatest "v$lv - you're up to date" $true; $bInstall.Text = 'Repair' }
        else { Set-Status $vLatest "v$lv available" $null; $bInstall.Text = $(if ($iv) { "Update to v$lv" } else { "Install v$lv" }) }
    } else { Set-Status $vLatest 'Not checked yet' $null; $bInstall.Text = 'Install / Update' }
}

function Get-Latest {
    $uri = "https://api.github.com/repos/$Owner/$Repo/releases/latest"
    $r = Invoke-RestMethod -Uri $uri -Headers @{ 'User-Agent' = 'BL2ArsenalLauncher'; 'Accept' = 'application/vnd.github+json' }
    return $r
}

function Check-Updates($quiet) {
    try {
        Say 'Checking GitHub for the latest version...'
        $script:latest = Get-Latest
        $asset = $script:latest.assets | Where-Object { $_.name -like '*.sdkmod' } | Select-Object -First 1
        if (-not $asset) { Say 'The latest release has no .sdkmod file attached.'; return }
        Say ("Latest: {0}  ({1})" -f $script:latest.tag_name, $script:latest.name)
        if ($script:latest.body) { Say ''; Say ($script:latest.body -replace "`r?`n", "`r`n"); Say '' }
    } catch {
        if (-not $quiet) { Say "Couldn't reach GitHub: $($_.Exception.Message)" } else { Say 'Update check skipped (offline?).' }
    }
    Refresh-Status
}

function Install-Sdkmod($file) {
    $tps = $cfg.TpsPath
    if (Game-Running) { [Windows.Forms.MessageBox]::Show('Close The Pre-Sequel first, then try again.', 'BL2 Arsenal') | Out-Null; return $false }
    $sdk = Join-Path $tps 'sdk_mods'
    if (-not (Test-Path $sdk)) { Say 'No sdk_mods folder - install willow2-mod-manager first.'; return $false }
    $tmp = Join-Path $env:TEMP ("bl2arsenal_" + [guid]::NewGuid().ToString('N'))
    try {
        Say "Unpacking $(Split-Path $file -Leaf)..."
        [IO.Compression.ZipFile]::ExtractToDirectory($file, $tmp)
        $src = Join-Path $tmp $ModFolder
        if (-not (Test-Path (Join-Path $src '__init__.py'))) { Say 'That file is not a BL2 Arsenal package.'; return $false }
        $dst = Join-Path $sdk $ModFolder
        $old = "$dst.old"
        if (Test-Path $old) { Remove-Item $old -Recurse -Force }
        if (Test-Path $dst) { Rename-Item $dst "$ModFolder.old" }
        try { Move-Item $src $dst }
        catch { if (Test-Path $old) { Rename-Item $old $ModFolder }; throw }
        if (Test-Path $old) { Remove-Item $old -Recurse -Force }
        # a loose .sdkmod copy would load the mod twice
        Get-ChildItem $sdk -Filter "$ModFolder*.sdkmod" -ErrorAction SilentlyContinue | Remove-Item -Force
        Say "Installed v$(Get-InstalledVersion $tps). Your settings, save backups and built packages were kept."
        return $true
    } catch { Say "Install failed: $($_.Exception.Message)"; return $false }
    finally { if (Test-Path $tmp) { Remove-Item $tmp -Recurse -Force -ErrorAction SilentlyContinue } }
}

$bCheck.Add_Click({ Check-Updates $false })
$bInstall.Add_Click({
    if (-not $script:latest) { Check-Updates $false }
    if (-not $script:latest) { return }
    $asset = $script:latest.assets | Where-Object { $_.name -like '*.sdkmod' } | Select-Object -First 1
    if (-not $asset) { Say 'No .sdkmod in the latest release.'; return }
    $dl = Join-Path $env:TEMP $asset.name
    try {
        Say "Downloading $($asset.name) ($([math]::Round($asset.size / 1MB, 2)) MB)..."
        Invoke-WebRequest -Uri $asset.browser_download_url -OutFile $dl -UseBasicParsing -Headers @{ 'User-Agent' = 'BL2ArsenalLauncher' }
        Install-Sdkmod $dl | Out-Null
    } catch { Say "Download failed: $($_.Exception.Message)" }
    finally { Remove-Item $dl -Force -ErrorAction SilentlyContinue }
    Refresh-Status
})
$bFile.Add_Click({
    $d = New-Object Windows.Forms.OpenFileDialog
    $d.Filter = 'BL2 Arsenal package (*.sdkmod)|*.sdkmod'; $d.Title = 'Choose a BL2 Arsenal .sdkmod'
    if ($d.ShowDialog() -eq 'OK') { Install-Sdkmod $d.FileName | Out-Null; Refresh-Status }
})
$bFolder.Add_Click({
    $d = New-Object Windows.Forms.FolderBrowserDialog
    $d.Description = 'Choose your Borderlands The Pre-Sequel folder (the one with Binaries and sdk_mods)'
    if ($d.ShowDialog() -eq 'OK') {
        if (Is-Tps $d.SelectedPath) { $cfg.TpsPath = $d.SelectedPath; Save-Config $cfg; Say "Game folder set: $($d.SelectedPath)" }
        else { Say "That folder doesn't contain Binaries\Win32\BorderlandsPreSequel.exe." }
        Refresh-Status
    }
})
$bOpen.Add_Click({ Start-Process explorer.exe (Join-Path $cfg.TpsPath 'sdk_mods') })
$bUninstall.Add_Click({
    $ans = [Windows.Forms.MessageBox]::Show("Remove BL2 Arsenal?`n`nPorted items disappear from your saves while the mod is removed (they come back if you reinstall). Your settings folder is kept.", 'Uninstall BL2 Arsenal', 'YesNo', 'Warning')
    if ($ans -ne 'Yes') { return }
    if (Game-Running) { [Windows.Forms.MessageBox]::Show('Close The Pre-Sequel first.', 'BL2 Arsenal') | Out-Null; return }
    $dst = Join-Path $cfg.TpsPath "sdk_mods\$ModFolder"
    try {
        if (Test-Path $dst) { Remove-Item $dst -Recurse -Force }
        Get-ChildItem (Join-Path $cfg.TpsPath 'sdk_mods') -Filter "$ModFolder*.sdkmod" -ErrorAction SilentlyContinue | Remove-Item -Force
        Say 'BL2 Arsenal removed.'
    } catch { Say "Uninstall failed: $($_.Exception.Message)" }
    Refresh-Status
})
$bPlay.Add_Click({ Say 'Starting The Pre-Sequel through Steam...'; Start-Process "steam://rungameid/$TpsAppId" })

$form.Add_Shown({
    Refresh-Status
    if (Is-Tps $cfg.TpsPath) { Save-Config $cfg }
    Check-Updates $true
})
[void]$form.ShowDialog()
