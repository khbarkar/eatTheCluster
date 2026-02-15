// internal/game/ghost.go
package game

import "math/rand"

type GhostBehavior int

const (
	GhostChaser   GhostBehavior = iota
	GhostPatrol
	GhostRandom
	GhostAmbusher
)

type Ghost struct {
	X, Y           int
	Behavior       GhostBehavior
	Frightened     bool
	FrightenTimer  int
	SpawnX, SpawnY int
}

func NewGhost(x, y int, behavior GhostBehavior) *Ghost {
	return &Ghost{
		X:        x,
		Y:        y,
		SpawnX:   x,
		SpawnY:   y,
		Behavior: behavior,
	}
}

func CreateGhosts() []*Ghost {
	return []*Ghost{
		NewGhost(12, 14, GhostChaser),
		NewGhost(13, 14, GhostPatrol),
		NewGhost(14, 14, GhostRandom),
		NewGhost(15, 14, GhostAmbusher),
	}
}

func (g *Ghost) Frighten(ticks int) {
	g.Frightened = true
	g.FrightenTimer = ticks
}

func (g *Ghost) Move(m *Maze, target Position) {
	if g.Frightened {
		g.FrightenTimer--
		if g.FrightenTimer <= 0 {
			g.Frightened = false
		}
		g.moveRandom(m)
		return
	}

	switch g.Behavior {
	case GhostChaser:
		g.moveToward(m, target)
	case GhostPatrol:
		g.moveToward(m, Position{X: g.SpawnX, Y: g.SpawnY})
	case GhostRandom:
		g.moveRandom(m)
	case GhostAmbusher:
		g.moveToward(m, Position{X: target.X + 4, Y: target.Y})
	}
}

func (g *Ghost) moveToward(m *Maze, target Position) {
	type candidate struct {
		x, y int
		dist int
	}
	moves := []struct{ dx, dy int }{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	var best candidate
	best.dist = 999999
	first := true

	for _, mv := range moves {
		nx, ny := g.X+mv.dx, g.Y+mv.dy
		if m.IsWall(nx, ny) {
			continue
		}
		dx := nx - target.X
		dy := ny - target.Y
		dist := dx*dx + dy*dy
		if first || dist < best.dist {
			best = candidate{nx, ny, dist}
			first = false
		}
	}
	if !first {
		g.X = best.x
		g.Y = best.y
	}
}

func (g *Ghost) moveRandom(m *Maze) {
	moves := []struct{ dx, dy int }{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	var valid []struct{ dx, dy int }
	for _, mv := range moves {
		if !m.IsWall(g.X+mv.dx, g.Y+mv.dy) {
			valid = append(valid, mv)
		}
	}
	if len(valid) > 0 {
		mv := valid[rand.Intn(len(valid))]
		g.X += mv.dx
		g.Y += mv.dy
	}
}

func (g *Ghost) Respawn() {
	g.X = g.SpawnX
	g.Y = g.SpawnY
	g.Frightened = false
	g.FrightenTimer = 0
}
