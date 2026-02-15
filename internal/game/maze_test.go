// internal/game/maze_test.go
package game

import "testing"

func TestNewMaze(t *testing.T) {
	m := NewMaze()
	if m.Width() == 0 || m.Height() == 0 {
		t.Fatal("maze dimensions must be non-zero")
	}
}

func TestMazeWalls(t *testing.T) {
	m := NewMaze()
	if !m.IsWall(0, 0) {
		t.Fatal("expected (0,0) to be a wall")
	}
}

func TestMazeDotPositions(t *testing.T) {
	m := NewMaze()
	positions := m.DotPositions()
	if len(positions) == 0 {
		t.Fatal("maze must have dot positions")
	}
	for _, p := range positions {
		if m.IsWall(p.X, p.Y) {
			t.Fatalf("dot position (%d,%d) is inside a wall", p.X, p.Y)
		}
	}
}

func TestMazePowerPelletPositions(t *testing.T) {
	m := NewMaze()
	pellets := m.PowerPelletPositions()
	if len(pellets) != 4 {
		t.Fatalf("expected 4 power pellets, got %d", len(pellets))
	}
}

func TestMazeRender(t *testing.T) {
	m := NewMaze()
	rendered := m.Render()
	if len(rendered) == 0 {
		t.Fatal("render must produce output")
	}
}
