package maps

import (
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestLoadMap(t *testing.T) {
	path := "map.txt"
	grid := []string{"1,0,0,0", "0,0,0,1", "0,1,0,0", "0,0,0,0", "0,0,1,0"}
	gridStr := strings.Join(grid, "\n")

	f, _ := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	f.WriteString(gridStr)

	err := LoadMap(path)
	if err != nil {
		t.Errorf("TestLoadMap = %s; want nil", err.Error())
	}

	for i := 0; i < len(grid); i++ {
		grid[i] = strings.ReplaceAll(grid[i], ",", "")
		want := grid[i]

		st := 0
		end := len(grid[i])
		got := ""
		for _, v := range MapObj.Grid[i][st:end] {
			got += strconv.Itoa(v)
		}

		if got != want {
			t.Errorf("TestLoadMap = %s; want %s", got, want)
		}
	}

}
