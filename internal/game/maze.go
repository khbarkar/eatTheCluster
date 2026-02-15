// internal/game/maze.go
package game

type Position struct {
	X, Y int
}

type CellType int

const (
	Wall CellType = iota
	Empty
	Dot
	PowerPellet
)

const (
	mazeWidth  = 28
	mazeHeight = 31
)

var mazeTemplate = []string{
	"############################",
	"#............##............#",
	"#.####.#####.##.#####.####.#",
	"#o####.#####.##.#####.####o#",
	"#.####.#####.##.#####.####.#",
	"#..........................#",
	"#.####.##.########.##.####.#",
	"#.####.##.########.##.####.#",
	"#......##....##....##......#",
	"######.##### ## #####.######",
	"     #.##### ## #####.#     ",
	"     #.##          ##.#     ",
	"     #.## ###--### ##.#     ",
	"######.## #      # ##.######",
	"      .   #      #   .      ",
	"######.## #      # ##.######",
	"     #.## ######## ##.#     ",
	"     #.##          ##.#     ",
	"     #.## ######## ##.#     ",
	"######.## ######## ##.######",
	"#............##............#",
	"#.####.#####.##.#####.####.#",
	"#.####.#####.##.#####.####.#",
	"#o..##................##..o#",
	"###.##.##.########.##.##.###",
	"###.##.##.########.##.##.###",
	"#......##....##....##......#",
	"#.##########.##.##########.#",
	"#.##########.##.##########.#",
	"#..........................#",
	"############################",
}

type Maze struct {
	cells   [][]CellType
	dots    map[Position]bool
	pellets map[Position]bool
}

func NewMaze() *Maze {
	m := &Maze{
		cells:   make([][]CellType, mazeHeight),
		dots:    make(map[Position]bool),
		pellets: make(map[Position]bool),
	}
	for y, row := range mazeTemplate {
		m.cells[y] = make([]CellType, mazeWidth)
		for x, ch := range row {
			if x >= mazeWidth {
				break
			}
			pos := Position{X: x, Y: y}
			switch ch {
			case '#':
				m.cells[y][x] = Wall
			case '.':
				m.cells[y][x] = Dot
				m.dots[pos] = true
			case 'o':
				m.cells[y][x] = PowerPellet
				m.pellets[pos] = true
			default:
				m.cells[y][x] = Empty
			}
		}
	}
	return m
}

func (m *Maze) Width() int  { return mazeWidth }
func (m *Maze) Height() int { return mazeHeight }

func (m *Maze) IsWall(x, y int) bool {
	if x < 0 || x >= mazeWidth || y < 0 || y >= mazeHeight {
		return true
	}
	return m.cells[y][x] == Wall
}

func (m *Maze) Cell(x, y int) CellType {
	if x < 0 || x >= mazeWidth || y < 0 || y >= mazeHeight {
		return Wall
	}
	return m.cells[y][x]
}

func (m *Maze) DotPositions() []Position {
	var positions []Position
	for p := range m.dots {
		positions = append(positions, p)
	}
	return positions
}

func (m *Maze) PowerPelletPositions() []Position {
	var positions []Position
	for p := range m.pellets {
		positions = append(positions, p)
	}
	return positions
}

func (m *Maze) HasDot(x, y int) bool {
	return m.dots[Position{X: x, Y: y}]
}

func (m *Maze) HasPellet(x, y int) bool {
	return m.pellets[Position{X: x, Y: y}]
}

func (m *Maze) EatDot(x, y int) bool {
	p := Position{X: x, Y: y}
	if m.dots[p] {
		delete(m.dots, p)
		m.cells[y][x] = Empty
		return true
	}
	return false
}

func (m *Maze) EatPellet(x, y int) bool {
	p := Position{X: x, Y: y}
	if m.pellets[p] {
		delete(m.pellets, p)
		m.cells[y][x] = Empty
		return true
	}
	return false
}

func (m *Maze) RemainingDots() int {
	return len(m.dots)
}

func (m *Maze) Render() string {
	var result []byte
	for y := 0; y < mazeHeight; y++ {
		for x := 0; x < mazeWidth; x++ {
			switch m.cells[y][x] {
			case Wall:
				result = append(result, '#')
			case Dot:
				result = append(result, '.')
			case PowerPellet:
				result = append(result, 'o')
			default:
				result = append(result, ' ')
			}
		}
		result = append(result, '\n')
	}
	return string(result)
}
