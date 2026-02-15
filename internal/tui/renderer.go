package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kristinb/eatthecluster/internal/game"
)

var (
	wallStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("27"))
	dotPodStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
	dotDeployStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("226"))
	dotStateful    = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	pacmanStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true)
	ghostStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	ghostScared    = lipgloss.NewStyle().Foreground(lipgloss.Color("21"))
	pelletStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true)
	hudStyle       = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1)
	warningStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	titleStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true)
)

type TickMsg time.Time

type Model struct {
	Game      *game.Game
	Context   string
	ChaosFunc func(game.ChaosEvent)
	quitting  bool
}

func NewModel(g *game.Game, context string, chaosFunc func(game.ChaosEvent)) Model {
	return Model{
		Game:      g,
		Context:   context,
		ChaosFunc: chaosFunc,
	}
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.Game.State {
		case game.StateWarning:
			if msg.Type == tea.KeyEnter {
				m.Game.AcceptWarning()
				return m, tickCmd()
			}
			if msg.String() == "q" {
				m.quitting = true
				return m, tea.Quit
			}
		case game.StatePlaying:
			switch msg.String() {
			case "q":
				m.quitting = true
				return m, tea.Quit
			case "up", "k":
				m.Game.Pacman.SetDirection(game.DirUp)
			case "down", "j":
				m.Game.Pacman.SetDirection(game.DirDown)
			case "left", "h":
				m.Game.Pacman.SetDirection(game.DirLeft)
			case "right", "l":
				m.Game.Pacman.SetDirection(game.DirRight)
			case "tab":
				m.Game.Pacman.ToggleMode()
			case " ":
				m.Game.TogglePause()
			}
		case game.StatePaused:
			switch msg.String() {
			case "q":
				m.quitting = true
				return m, tea.Quit
			case " ":
				m.Game.TogglePause()
			}
		case game.StateGameOver, game.StateWin:
			if msg.String() == "q" {
				m.quitting = true
				return m, tea.Quit
			}
		}

	case TickMsg:
		if m.Game.State == game.StatePlaying {
			events := m.Game.Tick()
			if m.ChaosFunc != nil {
				for _, e := range events {
					go m.ChaosFunc(e)
				}
			}
		}
		return m, tickCmd()
	}

	return m, nil
}

func (m Model) View() string {
	switch m.Game.State {
	case game.StateWarning:
		return m.viewWarning()
	case game.StateGameOver:
		return m.viewGameOver()
	case game.StateWin:
		return m.viewWin()
	default:
		return m.viewGame()
	}
}

func (m Model) viewWarning() string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(warningStyle.Render("  *** WARNING ***"))
	b.WriteString("\n\n")
	b.WriteString("  This tool performs REAL chaos actions on your cluster.\n")
	b.WriteString(fmt.Sprintf("  Context: %s\n", m.Context))
	if m.Game.DryRun {
		b.WriteString("  Mode: DRY RUN (no real actions)\n")
	} else {
		b.WriteString(warningStyle.Render("  Mode: LIVE (resources WILL be affected)"))
		b.WriteString("\n")
	}
	b.WriteString("\n  Press ENTER to accept, q to quit.\n")
	return b.String()
}

func (m Model) viewGame() string {
	maze := m.renderMaze()
	hud := m.renderHUD()
	return lipgloss.JoinHorizontal(lipgloss.Top, maze, "  ", hud)
}

func (m Model) renderMaze() string {
	g := m.Game
	var b strings.Builder

	ghostPos := make(map[game.Position]*game.Ghost)
	for _, gh := range g.Ghosts {
		ghostPos[game.Position{X: gh.X, Y: gh.Y}] = gh
	}

	for y := 0; y < g.Maze.Height(); y++ {
		for x := 0; x < g.Maze.Width(); x++ {
			if g.Pacman.X == x && g.Pacman.Y == y {
				b.WriteString(pacmanStyle.Render("C"))
				continue
			}
			if gh, ok := ghostPos[game.Position{X: x, Y: y}]; ok {
				if gh.Frightened {
					b.WriteString(ghostScared.Render("M"))
				} else {
					b.WriteString(ghostStyle.Render("M"))
				}
				continue
			}
			cell := g.Maze.Cell(x, y)
			switch cell {
			case game.Wall:
				b.WriteString(wallStyle.Render("#"))
			case game.Dot:
				if res, ok := g.Resources[game.Position{X: x, Y: y}]; ok {
					switch res.Kind {
					case "StatefulSet":
						b.WriteString(dotStateful.Render("."))
					case "Deployment":
						b.WriteString(dotDeployStyle.Render("."))
					default:
						b.WriteString(dotPodStyle.Render("."))
					}
				} else {
					b.WriteString(dotPodStyle.Render("."))
				}
			case game.PowerPellet:
				b.WriteString(pelletStyle.Render("O"))
			default:
				b.WriteString(" ")
			}
		}
		b.WriteString("\n")
	}

	if g.State == game.StatePaused {
		b.WriteString("\n")
		b.WriteString(titleStyle.Render("  PAUSED - Press SPACE to resume"))
	}

	return b.String()
}

func (m Model) renderHUD() string {
	g := m.Game
	var lines []string

	lines = append(lines, titleStyle.Render("EAT THE CLUSTER"))
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("Mode:  %s", g.Pacman.ModeString()))
	lines = append(lines, fmt.Sprintf("Lives: %s", strings.Repeat("C ", g.Pacman.Lives)))
	lines = append(lines, fmt.Sprintf("Score: %d", g.Pacman.Score))
	lines = append(lines, fmt.Sprintf("Eaten: %d / %d", g.DotsEaten, g.TotalDots))
	lines = append(lines, "")
	lines = append(lines, "Chaos Log:")
	for _, e := range g.ChaosLog {
		lines = append(lines, fmt.Sprintf("  %s %s/%s [%s]", e.Action, e.Resource.Namespace, e.Resource.Name, e.Status))
	}
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("Cluster: %s", m.Context))
	if g.DryRun {
		lines = append(lines, warningStyle.Render("DRY RUN"))
	}

	return hudStyle.Render(strings.Join(lines, "\n"))
}

func (m Model) viewGameOver() string {
	return fmt.Sprintf("\n%s\n\n  Monitoring caught you! Score: %d\n  Press q to quit.\n",
		warningStyle.Render("  GAME OVER"), m.Game.Pacman.Score)
}

func (m Model) viewWin() string {
	return fmt.Sprintf("\n%s\n\n  You ate the entire cluster! Score: %d\n  Press q to quit.\n",
		titleStyle.Render("  YOU WIN!"), m.Game.Pacman.Score)
}
