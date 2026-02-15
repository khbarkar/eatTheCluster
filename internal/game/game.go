// internal/game/game.go
package game

type GameState int

const (
	StateWarning GameState = iota
	StatePlaying
	StatePaused
	StateGameOver
	StateWin
)

type Resource struct {
	Name      string
	Namespace string
	Kind      string
}

type ChaosEvent struct {
	Resource Resource
	Action   string
	Status   string
}

type Game struct {
	Maze       *Maze
	Pacman     *Pacman
	Ghosts     []*Ghost
	State      GameState
	DryRun     bool
	Resources  map[Position]Resource
	ChaosLog   []ChaosEvent
	TickCount  int
	TotalDots  int
	DotsEaten  int
	GhostSpeed int
}

func NewGame(resources []Resource, dryRun bool) *Game {
	maze := NewMaze()
	dotPositions := maze.DotPositions()

	resMap := make(map[Position]Resource)
	for i, r := range resources {
		if i >= len(dotPositions) {
			break
		}
		resMap[dotPositions[i]] = r
	}

	return &Game{
		Maze:       maze,
		Pacman:     NewPacman(14, 23),
		Ghosts:     CreateGhosts(),
		State:      StateWarning,
		DryRun:     dryRun,
		Resources:  resMap,
		ChaosLog:   make([]ChaosEvent, 0),
		TotalDots:  maze.RemainingDots(),
		GhostSpeed: 2,
	}
}

func (g *Game) AcceptWarning() {
	g.State = StatePlaying
}

func (g *Game) TogglePause() {
	if g.State == StatePlaying {
		g.State = StatePaused
	} else if g.State == StatePaused {
		g.State = StatePlaying
	}
}

func (g *Game) Tick() []ChaosEvent {
	if g.State != StatePlaying {
		return nil
	}
	g.TickCount++

	var events []ChaosEvent

	g.Pacman.Move(g.Maze)

	px, py := g.Pacman.X, g.Pacman.Y
	if g.Maze.EatDot(px, py) {
		g.DotsEaten++
		g.Pacman.Score += 10
		if res, ok := g.Resources[Position{X: px, Y: py}]; ok {
			action := "killed"
			if g.Pacman.Mode == PoisonMode {
				action = "degraded"
			}
			event := ChaosEvent{Resource: res, Action: action, Status: "pending"}
			events = append(events, event)
			g.ChaosLog = append(g.ChaosLog, event)
			if len(g.ChaosLog) > 3 {
				g.ChaosLog = g.ChaosLog[len(g.ChaosLog)-3:]
			}
		}
	}

	if g.Maze.EatPellet(px, py) {
		g.Pacman.Score += 50
		for _, ghost := range g.Ghosts {
			ghost.Frighten(50)
		}
	}

	if g.TickCount%g.GhostSpeed == 0 {
		target := Position{X: g.Pacman.X, Y: g.Pacman.Y}
		for _, ghost := range g.Ghosts {
			ghost.Move(g.Maze, target)
		}
	}

	totalDots := g.TotalDots
	if totalDots < 1 {
		totalDots = 1
	}
	progress := float64(g.DotsEaten) / float64(totalDots)
	if progress > 0.75 {
		g.GhostSpeed = 1
	} else if progress > 0.5 {
		g.GhostSpeed = 1
	}

	for _, ghost := range g.Ghosts {
		if ghost.X == g.Pacman.X && ghost.Y == g.Pacman.Y {
			if ghost.Frightened {
				ghost.Respawn()
				g.Pacman.Score += 200
			} else {
				g.Pacman.Lives--
				if g.Pacman.Lives <= 0 {
					g.State = StateGameOver
				} else {
					g.Pacman.X = 14
					g.Pacman.Y = 23
					g.Pacman.SetDirection(DirNone)
				}
				if len(g.ChaosLog) > 0 {
					g.ChaosLog[len(g.ChaosLog)-1].Status = "blocked"
				}
				break
			}
		}
	}

	g.checkWinCondition()
	return events
}

func (g *Game) checkWinCondition() {
	if g.Maze.RemainingDots() == 0 && len(g.Maze.PowerPelletPositions()) == 0 {
		g.State = StateWin
	}
}
