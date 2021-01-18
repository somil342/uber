package maps

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/somil342/uber/myobjects/locations"
)

//MaxRow ...
var MaxRow int = 1000

//MaxCol ...
var MaxCol int = 1000

//Map ...
type Map struct {
	Grid [][]int `json:"grid"`
	Row  int     `json:"row"`
	Col  int     `json:"col"`
}

//MapObj ...
var MapObj *Map = &Map{}

func init() {

	MapObj.Row = MaxRow
	MapObj.Col = MaxCol

	MapObj.Grid = make([][]int, MaxRow)
	for r := range MapObj.Grid {
		MapObj.Grid[r] = make([]int, MaxCol)
	}
}

//LoadMap ...
func LoadMap(path string) error {

	if path == "" {
		path = "map.txt"
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	rowNum := 0
	colNum := 0
	mapFromFile := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		rowNum++
		line := scanner.Text()
		line = strings.Trim(line, " ")
		res := strings.Split(line, ",")
		line = strings.Join(res, "")
		if colNum < len(res) {
			colNum = len(res)
		}
		mapFromFile = append(mapFromFile, line)
	}

	for i := 0; i < rowNum; i++ {
		for j := 0; j < colNum; j++ {
			MapObj.Grid[i][j] = 1
		}
	}

	for i := 0; i < len(mapFromFile); i++ {
		for j := 0; j < len(mapFromFile[i]); j++ {
			val := string(mapFromFile[i][j])
			if val == "0" {
				MapObj.Grid[i][j] = 0
			}
		}
	}

	return nil
}

//Bfs ... returns Shortest path,price,eta between 2 location using bfs traversal
func (mp *Map) Bfs(l1 *locations.Location, l2 *locations.Location) ([]string, int, int, error) {

	type qItem struct {
		l *locations.Location
	}

	fmt.Println("src", *l1)
	fmt.Println("dest", *l2)
	if !mp.checkBoundary(l1) || !mp.checkBoundary(l2) || *l1 == *l2 {
		return []string{}, 0, 0, errors.New("invalid location")
	}

	if *l1 == *l2 {
		return []string{}, 0, 0, errors.New("src and dest location are same")
	}

	//directions
	dirs := [][]int{
		{1, 0},  //up
		{0, 1},  //right
		{-1, 0}, //down
		{0, -1}, //left
	}

	vis := make([][]bool, mp.Row)
	parent := make([][]*locations.Location, mp.Row)

	for i := range vis {
		vis[i] = make([]bool, mp.Col)
		parent[i] = make([]*locations.Location, mp.Col)
	}

	q := make(chan qItem, MaxRow*MaxCol)

	s := qItem{
		l: &locations.Location{
			Row: l1.Row,
			Col: l1.Col,
		},
	}

	q <- s
	vis[s.l.Row][s.l.Col] = true

	possible := false
L:
	for {
		select {
		case node := <-q:
			if node.l.Row == l2.Row && node.l.Col == l2.Col {
				possible = true
				break L
			}

			for _, v := range dirs {
				neigh := qItem{
					l: &locations.Location{
						Row: node.l.Row + v[0],
						Col: node.l.Col + v[1],
					},
				}
				if mp.checkBoundary(neigh.l) && !vis[neigh.l.Row][neigh.l.Col] {
					q <- neigh
					vis[neigh.l.Row][neigh.l.Col] = true
					parent[neigh.l.Row][neigh.l.Col] = &locations.Location{
						Row: node.l.Row,
						Col: node.l.Col,
					}
				}
			}

		default:
			break L
		}

	}

	if !possible {
		return []string{}, 0, 0, errors.New("Not possible")
	}

	temp := l2

	path := []string{}
	path = append(path, fmt.Sprint("(", temp.Row, ",", temp.Col, ")"))
	for {
		temp = parent[temp.Row][temp.Col]
		path = append(path, fmt.Sprint("(", temp.Row, ",", temp.Col, ")"))
		if *temp == *l1 {
			break
		}
	}

	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	//path len * 100 price , path len * 10 eta(time in minutes)
	return path, len(path) * 100, len(path) * 10, nil
}

// check for boundary condtion and blocked area
func (mp *Map) checkBoundary(l *locations.Location) bool {
	return l.Row >= 0 && l.Row < mp.Row && l.Col >= 0 && l.Col < mp.Col && mp.Grid[l.Row][l.Col] == 0
}

//PrintPath ...
func (mp *Map) PrintPath(path []string) string {

	if len(path) == 0 {
		return ""
	}

	str := "[" + path[0]
	for i := 1; i < len(path); i++ {
		str += " -> " + path[i]
	}
	return str + "]"
}
