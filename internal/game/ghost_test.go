// internal/game/ghost_test.go
package game

import "testing"

func TestNewGhost(t *testing.T) {
	g := NewGhost(14, 14, GhostChaser)
	if g.X != 14 || g.Y != 14 {
		t.Fatalf("expected (14,14), got (%d,%d)", g.X, g.Y)
	}
	if g.Behavior != GhostChaser {
		t.Fatal("expected GhostChaser behavior")
	}
	if g.Frightened {
		t.Fatal("ghost should not start frightened")
	}
}

func TestGhostChaserMovesTowardTarget(t *testing.T) {
	m := NewMaze()
	g := NewGhost(1, 5, GhostChaser)
	g.Move(m, Position{X: 10, Y: 5})
	if g.X <= 1 {
		t.Fatal("chaser should move toward target")
	}
}

func TestGhostScatter(t *testing.T) {
	g := NewGhost(14, 14, GhostChaser)
	g.Frighten(5)
	if !g.Frightened {
		t.Fatal("ghost should be frightened after Frighten()")
	}
}

func TestCreateAllGhosts(t *testing.T) {
	ghosts := CreateGhosts()
	if len(ghosts) != 4 {
		t.Fatalf("expected 4 ghosts, got %d", len(ghosts))
	}
}
