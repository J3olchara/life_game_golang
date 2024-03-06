package life

import (
	"bufio"
	"errors"
	"math/rand"
	"os"
	"time"
)

type World struct {
	Height int
	Width  int
	Cells  [][]bool
}

func NewWorld(height, width int) (*World, error) {
	if height <= 0 || width <= 0 {
		return nil, errors.New("height and width have to be positive")
	}
	cells := make([][]bool, height)
	for i := range cells {
		cells[i] = make([]bool, width)
	}
	return &World{
		Height: height,
		Width:  width,
		Cells:  cells,
	}, nil
}

func (w *World) RandInit(percentage int) {
	numAlive := percentage * w.Height * w.Width / 100
	w.fillAlive(numAlive)
	r := rand.New(rand.NewSource(time.Now().Unix()))

	for i := 0; i < w.Height*w.Width; i++ {
		randRowLeft := r.Intn(w.Width)
		randColLeft := r.Intn(w.Height)
		randRowRight := r.Intn(w.Width)
		randColRight := r.Intn(w.Height)

		w.Cells[randRowLeft][randColLeft] = w.Cells[randRowRight][randColRight]
	}
}

func (w *World) fillAlive(num int) {
	aliveCount := 0
	for j, row := range w.Cells {
		for k := range row {
			w.Cells[j][k] = true
			aliveCount++
			if aliveCount == num {

				return
			}
		}
	}
}

func (w *World) Neighbors(x, y int) int {
	var cnt int
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			cordX := x + i
			cordY := y + j
			if cordY == -1 {
				cordY = len(w.Cells) - 1
			} else if cordY == len(w.Cells) {
				cordY = 0
			}
			if cordX == -1 {
				cordX = len(w.Cells[0]) - 1
			} else if cordX == len(w.Cells) {
				cordX = 0
			}
			if w.Cells[cordY][cordX] && (cordY != y || cordX != x) {
				cnt++
			}
		}
	}
	return cnt
}

func (w *World) Next(x, y int) bool {
	n := w.Neighbors(x, y)
	alive := w.Cells[y][x]
	if n < 4 && n > 1 && alive {
		return true
	}
	if n == 3 && !alive {
		return true
	}

	return false
}

func NextState(oldWorld, newWorld *World) {
	for i := 0; i < oldWorld.Height; i++ {
		for j := 0; j < oldWorld.Width; j++ {
			newWorld.Cells[i][j] = oldWorld.Next(j, i)
		}
	}
}

func (w *World) Seed() {
	for _, row := range w.Cells {
		for i := range row {
			if rand.Intn(10) == 1 {
				row[i] = true
			}
		}
	}
}

func (w *World) String(tCell, fCell string) string {
	var world string
	for i := range w.Cells {
		for j := range w.Cells[i] {
			l := fCell
			if w.Cells[i][j] {
				l = tCell
			}
			world += l
		}
		if i < len(w.Cells)-1 {
			world += "\n"
		}
	}
	return world
}

func (w *World) SaveState(filename string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0600)
	defer f.Close()
	if err != nil {
		return err
	}
	for i := range w.Cells {
		for j := range w.Cells[i] {
			l := "0"
			if w.Cells[i][j] {
				l = "1"
			}
			_, err := f.WriteString(l)
			if err != nil {
				return err
			}
		}
		if i != len(w.Cells)-1 {
			_, err := f.WriteString("\n")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (w *World) LoadState(filename string) error {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0600)
	defer f.Close()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(f)
	w.Cells = make([][]bool, 0)
	widthSet := false
	for scanner.Scan() {
		line := scanner.Text()
		if widthSet && w.Width != len(line) {
			return errors.New("file corrupted")
		}
		if !widthSet {
			widthSet = true
			w.Width = len(line)
		}

		worldLine := make([]bool, w.Width)
		for i := range line {
			worldLine[i] = line[i] == '1'
		}

		w.Cells = append(w.Cells, worldLine)
		w.Height++
	}
	return nil
}
