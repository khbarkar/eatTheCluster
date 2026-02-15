// internal/game/pacman.go
package game

// ChaosMode determines what happens when Pacman eats a dot.
type ChaosMode int

const (
	ChompMode  ChaosMode = iota // Kill the resource
	PoisonMode                  // Degrade the resource
)

// Direction represents movement direction.
type Direction int

const (
	DirNone Direction = iota
	DirUp
	DirDown
	DirLeft
	DirRight
)

// Pacman is the player character.
type Pacman struct {
	X, Y      int
	Lives     int
	Mode      ChaosMode
	Score     int
	direction Direction
}

// NewPacman creates a pacman at the given position with 3 lives.
func NewPacman(x, y int) *Pacman {
	return &Pacman{
		X:     x,
		Y:     y,
		Lives: 3,
		Mode:  ChompMode,
	}
}

// ToggleMode switches between Chomp and Poison mode.
func (p *Pacman) ToggleMode() {
	if p.Mode == ChompMode {
		p.Mode = PoisonMode
	} else {
		p.Mode = ChompMode
	}
}

// SetDirection sets the intended movement direction.
func (p *Pacman) SetDirection(d Direction) {
	p.direction = d
}

// Move attempts to move Pacman in the current direction.
func (p *Pacman) Move(m *Maze) {
	nx, ny := p.X, p.Y
	switch p.direction {
	case DirUp:
		ny--
	case DirDown:
		ny++
	case DirLeft:
		nx--
	case DirRight:
		nx++
	}
	if !m.IsWall(nx, ny) {
		p.X = nx
		p.Y = ny
	}
}

func (p *Pacman) ModeString() string {
	if p.Mode == ChompMode {
		return "CHOMP"
	}
	return "POISON"
}
