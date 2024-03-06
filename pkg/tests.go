package life_test

import (
	"github.com/J3olchara/game/pkg/life"
	"testing"
)

func TestNewWorld(t *testing.T) {
	height := 10
	width := 4

	world, _ := life.NewWorld(height, width)

	if world.Height != height {
		t.Errorf("expected height: %d, actual height: %d", height, world.Height)
	}

	if world.Width != width {
		t.Errorf("expected width: %d, actual width: %d", width, world.Width)
	}

	if len(world.Cells) != height {
		t.Errorf("expected height: %d, actual number of rows: %d", height, len(world.Cells))
	}
	// Проверяем, что у каждого элемента — заданная длина
	for i, row := range world.Cells {
		if len(row) != width {
			t.Errorf("expected width: %d, actual row's %d len: %d", width, i, world.Width)
		}
	}
}
