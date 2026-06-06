package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("63")).
			Bold(true)

	accentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	dangerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	mutedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	phaseOK    = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("✓")
	phaseActive = lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Render("◉")
	phasePending = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("□")
)

type model struct {
	activeTab int
	tabs      []string
	width     int
	height    int
	logs      []string
}

func initialModel() model {
	return model{
		activeTab: 0,
		tabs:      []string{"Dashboard", "Campaigns", "Agents", "AI", "Logs", "Lab"},
		logs: []string{
			"[14:32:01] WS1: PrivEsc SUID python → root ✓",
			"[14:31:45] Recon: Port 445 open on .10",
			"[14:31:30] AI: GTFOBins vector suggested (conf=0.89)",
			"[14:31:00] Agent WS1 checked in (session=abc123)",
			"[14:30:45] Horizon-Intel: 23 hosts discovered",
		},
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab", "right", "l":
			m.activeTab = (m.activeTab + 1) % len(m.tabs)
		case "shift+tab", "left", "h":
			m.activeTab = (m.activeTab - 1 + len(m.tabs)) % len(m.tabs)
		}
	}
	return m, nil
}

func (m model) View() string {
	header := m.renderHeader()
	tabs := m.renderTabs()
	content := m.renderContent()

	return lipgloss.JoinVertical(lipgloss.Left, header, tabs, content)
}

func (m model) renderHeader() string {
	banner := titleStyle.Render(`
██╗  ██╗ ██╗  ██╗  ██████╗  ██╗  ██╗ ██╗  ██╗
╚██╗██╔╝ ██║  ██║ ██╔═══██╗ ██║  ██║ ╚██╗██╔╝
 ╚███╔╝  ███████║ ██║   ██║ ███████║  ╚███╔╝
 ██╔██╗  ╚════██║ ██║   ██║ ╚════██║  ██╔██╗
██╔╝ ██╗     ██╔╝ ╚██████╔╝     ██╔╝ ██╔╝ ██╗
╚═╝  ╚═╝     ╚═╝   ╚═════╝      ╚═╝  ╚═╝  ╚═╝`)

	info := mutedStyle.Render(fmt.Sprintf("v1.0 | %dx%d | [CTRL+C] Quit | [Tab] Switch", m.width, m.height))
	return lipgloss.JoinHorizontal(lipgloss.Top, banner, "   ", info)
}

func (m model) renderTabs() string {
	var rendered []string
	for i, tab := range m.tabs {
		if i == m.activeTab {
			rendered = append(rendered, accentStyle.Render(fmt.Sprintf("[%s]", tab)))
		} else {
			rendered = append(rendered, mutedStyle.Render(fmt.Sprintf(" %s ", tab)))
		}
	}
	return strings.Join(rendered, " │ ")
}

func (m model) renderContent() string {
	switch m.activeTab {
	case 0:
		return m.dashboardTab()
	case 1:
		return m.campaignsTab()
	case 2:
		return m.agentsTab()
	case 3:
		return m.aiTab()
	case 4:
		return m.logsTab()
	case 5:
		return m.labTab()
	default:
		return "No content"
	}
}

func (m model) dashboardTab() string {
	title := titleStyle.Render("╭─ DASHBOARD ─ Campaign: \"TFG-Demo\" · Phase: Exploitation ──────╮")

	killChain := m.renderKillChain()

	networkMap := borderStyle.Copy().Width(34).Height(10).Render(
		" NETWORK MAP\n" +
			" ● DC (10.0.0.10) [compromised]\n" +
			" │  ├─ ● WS1 (10.0.0.50) [user]\n" +
			" │  └─ ● WS2 (10.0.0.51) [user]\n" +
			" ● DB (10.0.0.20) [scanned]")

	aiPanel := borderStyle.Copy().Width(34).Height(10).Render(
		" AI CONSOLE\n" +
			accentStyle.Render(" > analyze current target\n") +
			"\n" +
			" [Specter] Analysis:\n" +
			"  ● SMB 445 → EternalBlue\n" +
			"  ● RDP 3389 → BlueKeep\n" +
			"  ● OS: Win2019\n" +
			"\n" +
			accentStyle.Render(" Recommendations:") + "\n" +
			"  [1] EternalBlue [0.95]\n" +
			"  [2] Kerberoast  [0.82]")

	events := borderStyle.Copy().Width(70).Height(6).Render(
		" LIVE FEED\n" + strings.Repeat("─", 68) + "\n" +
			strings.Join(m.logs, "\n"))

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, networkMap, "  ", aiPanel)
	bottomRow := events

	return lipgloss.JoinVertical(lipgloss.Left, title, killChain, " ", topRow, " ", bottomRow)
}

func (m model) renderKillChain() string {
	phases := []struct {
		name string
		done bool
		act  bool
	}{
		{"Recon", true, false},
		{"Weaponize", true, false},
		{"Deliver", true, false},
		{"Exploit", false, true},
		{"Install", false, false},
		{"C2", false, false},
	}

	var items []string
	for _, p := range phases {
		switch {
		case p.done:
			items = append(items, fmt.Sprintf(" %s %s", p.name, phaseOK))
		case p.act:
			items = append(items, fmt.Sprintf(" %s %s", p.name, phaseActive))
		default:
			items = append(items, fmt.Sprintf(" %s %s", p.name, phasePending))
		}
	}

	bar := "██████████████░░░░░░░"
	pc := " 67%"

	return borderStyle.Copy().Width(70).Render(
		" KILL CHAIN " + strings.Join(items, " │ ") + "\n" +
			" " + accentStyle.Render(bar) + mutedStyle.Render(pc))
}

func (m model) campaignsTab() string {
	return borderStyle.Copy().Width(70).Height(10).Render(
		titleStyle.Render(" CAMPAIGNS") + "\n\n" +
			"  Active:\n" +
			"  ▸ TFG-Demo  | running  | exploitation | 5 agents | 67%\n" +
			"\n" +
			"  Completed:\n" +
			"  (none)")
}

func (m model) agentsTab() string {
	return borderStyle.Copy().Width(70).Height(10).Render(
		titleStyle.Render(" AGENTS") + "\n\n" +
			"  ID       Hostname    OS          User     Status\n" +
			"  ──────── ─────────── ─────────── ──────── ──────\n" +
			"  abc123   WS1         Windows2019 SYSTEM   " + accentStyle.Render("online") + "\n" +
			"  def456   DB          Ubuntu24.04 root     " + accentStyle.Render("online") + "\n" +
			"  ghi789   WS2         Windows11   user     " + mutedStyle.Render("idle"))
}

func (m model) aiTab() string {
	return borderStyle.Copy().Width(70).Height(10).Render(
		titleStyle.Render(" AI CONSOLE (Specter + Apex + Ollama)") + "\n\n" +
			"  Model: llama3.2 | Mode: manual | Ollama: localhost:11434\n" +
			"  ─────────────────────────────────────────────────────\n" +
			"  > analyze target 10.0.0.10\n" +
			"\n" +
			"  [Specter] Host Analysis:\n" +
			"    OS: Windows Server 2019\n" +
			"    Services: SMB(445), RDP(3389), DNS(53)\n" +
			"    Vulnerabilities: MS17-010, CVE-2020-1472\n" +
			"\n" +
			"  [Apex] Recommendation: Use EternalBlue → Kerberoast chain\n" +
			"  [1] Accept  [2] Reject  [3] More options")
}

func (m model) logsTab() string {
	return borderStyle.Copy().Width(70).Height(12).Render(
		titleStyle.Render(" EVENT LOG") + "\n" +
			strings.Repeat("─", 68) + "\n" +
			strings.Join(m.logs, "\n") + "\n" +
			strings.Repeat("─", 68) + "\n" +
			mutedStyle.Render(" [F] Filter [S] Search [/] Find [R] Refresh"))
}

func (m model) labTab() string {
	return borderStyle.Copy().Width(70).Height(10).Render(
		titleStyle.Render(" LAB ENVIRONMENT") + "\n\n" +
			"  Container          Status    IP\n" +
			"  ─────────────────  ────────  ───────────\n" +
			"  x404x-attacker     " + accentStyle.Render("running") + "  172.20.0.10\n" +
			"  x404x-target1      " + accentStyle.Render("running") + "  172.20.0.20\n" +
			"  x404x-target2      " + accentStyle.Render("running") + "  172.20.0.21\n" +
			"  x404x-dashboard    " + accentStyle.Render("running") + "  172.20.0.30\n" +
			"  x404x-ollama       " + accentStyle.Render("running") + "  172.20.0.40\n" +
			"\n" +
			"  Dashboard: " + accentStyle.Render("http://localhost:3000") + "\n" +
			mutedStyle.Render("  [U] Start  [D] Stop  [S] Status"))
}

func StartTUI() error {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
