// internal/game/game_test.go
package game

import "testing"

func TestNewGame(t *testing.T) {
	resources := []Resource{
		{Name: "nginx-abc", Namespace: "default", Kind: "Pod"},
		{Name: "redis-xyz", Namespace: "default", Kind: "Pod"},
	}
	g := NewGame(resources, true)
	if g.State != StateWarning {
		t.Fatal("game should start in warning state")
	}
	if !g.DryRun {
		t.Fatal("expected dry run mode")
	}
}

func TestGameAcceptWarning(t *testing.T) {
	g := NewGame(nil, false)
	g.AcceptWarning()
	if g.State != StatePlaying {
		t.Fatal("expected StatePlaying after accepting warning")
	}
}

func TestGameTick(t *testing.T) {
	resources := []Resource{
		{Name: "pod1", Namespace: "default", Kind: "Pod"},
	}
	g := NewGame(resources, true)
	g.AcceptWarning()
	g.Pacman.SetDirection(DirRight)
	g.Tick() // should not panic
}

func TestGameWin(t *testing.T) {
	g := NewGame(nil, true)
	g.AcceptWarning()
	if g.Maze.RemainingDots() == 0 && len(g.Maze.PowerPelletPositions()) == 0 {
		g.checkWinCondition()
	}
}
