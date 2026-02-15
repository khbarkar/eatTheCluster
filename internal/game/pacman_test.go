// internal/game/pacman_test.go
package game

import "testing"

func TestNewPacman(t *testing.T) {
	p := NewPacman(14, 23)
	if p.X != 14 || p.Y != 23 {
		t.Fatalf("expected position (14,23), got (%d,%d)", p.X, p.Y)
	}
	if p.Lives != 3 {
		t.Fatalf("expected 3 lives, got %d", p.Lives)
	}
	if p.Mode != ChompMode {
		t.Fatal("expected default mode to be ChompMode")
	}
}

func TestPacmanToggleMode(t *testing.T) {
	p := NewPacman(14, 23)
	p.ToggleMode()
	if p.Mode != PoisonMode {
		t.Fatal("expected PoisonMode after toggle")
	}
	p.ToggleMode()
	if p.Mode != ChompMode {
		t.Fatal("expected ChompMode after second toggle")
	}
}

func TestPacmanMove(t *testing.T) {
	m := NewMaze()
	p := NewPacman(1, 1)
	p.SetDirection(DirRight)
	p.Move(m)
	if p.X != 2 || p.Y != 1 {
		t.Fatalf("expected (2,1) after moving right, got (%d,%d)", p.X, p.Y)
	}
}

func TestPacmanBlockedByWall(t *testing.T) {
	m := NewMaze()
	p := NewPacman(1, 1)
	p.SetDirection(DirUp)
	p.Move(m)
	if p.X != 1 || p.Y != 1 {
		t.Fatalf("pacman should not move into wall, got (%d,%d)", p.X, p.Y)
	}
}
