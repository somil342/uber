package cabs

import (
	"errors"
	"sync"

	"github.com/somil342/uber/models"
	"github.com/somil342/uber/myobjects/locations"
	"github.com/somil342/uber/myobjects/maps"
)

//Cab ...
type Cab struct {
	CabID       int                `json:"cabid,omitempty"`
	DriverName  string             `json:"drivername,omitempty"`
	CurLocation locations.Location `json:"curlocation,omitempty"`
	IsAvailable bool               `json:"isavailable"`
}

//AllCabs ...
var AllCabs map[int]*Cab

//GetCabData ...
func GetCabData(cabID int) Cab {

	cab, ok := AllCabs[cabID]
	if !ok {
		return Cab{}
	}
	return *cab
}

//cab map
type cabMaps struct {
	grid  [][]([]int)
	r     int
	c     int
	mutex *sync.Mutex
}

//cabmap
var cabMap cabMaps

func init() {
	AllCabs = make(map[int]*Cab)

	cabMap = cabMaps{}
	cabMap.r = maps.MaxRow
	cabMap.c = maps.MaxCol
	cabMap.mutex = &sync.Mutex{}

	cabMap.grid = make([][]([]int), cabMap.r)
	for i := 0; i < cabMap.r; i++ {
		cabMap.grid[i] = make([]([]int), cabMap.c)
		for j := 0; j < cabMap.c; j++ {
			cabMap.grid[i][j] = make([]int, 0)
		}
	}

}

//LoadCabs ...load all cabs in AllCabs(map of all cabs) and cabMap(2d map)
func LoadCabs() error {

	cabs, err := models.SelectAllCabs()
	if err != nil {
		panic(err)
	}

	r := 0
	c := 0
	i := 0
	for i < len(cabs) {
		v := cabs[i]

		_, err := New(v.CabID, v.DriverName, locations.Location{Row: r, Col: c})
		if err == nil {
			i++
		}

		c++
		if c >= maps.MaxCol {
			r++
			c = 0
		}
		if r >= maps.MaxRow {
			break
		}
	}

	return nil
}

//UpdateLocation ...
func (c *Cab) UpdateLocation(loc locations.Location) error {

	if !cabMap.isThereAWay(&loc) {
		return errors.New("invalid location")
	}

	cabMap.mutex.Lock()

	//old cab
	if c.CurLocation != (locations.Location{}) {
		for i, v := range cabMap.grid[c.CurLocation.Row][c.CurLocation.Col] {
			if v == c.CabID {
				cabMap.grid[c.CurLocation.Row][c.CurLocation.Col][i] = -1
			}
		}
	}

	for i, v := range cabMap.grid[loc.Row][loc.Col] {
		if v == -1 {
			cabMap.grid[loc.Row][loc.Col][i] = c.CabID
			cabMap.mutex.Unlock()
			c.CurLocation = loc
			return nil
		}
	}

	cabMap.grid[loc.Row][loc.Col] = append(cabMap.grid[loc.Row][loc.Col], c.CabID)
	cabMap.mutex.Unlock()
	c.CurLocation = loc

	return nil
}

//New ...
func New(cabID int, driverName string, loc locations.Location) (*Cab, error) {

	cab := &Cab{
		CabID:       cabID,
		DriverName:  driverName,
		IsAvailable: true,
	}
	err := cab.UpdateLocation(loc)

	if err != nil {
		return cab, nil
	}

	AllCabs[cab.CabID] = cab
	return cab, nil
}

//FindCabsNearBy ... give me all nearest cabs which are available
func FindCabsNearBy(l *locations.Location, distance int) ([]int, error) {

	nearByCabs := make([]int, 0)

	if distance < 0 {
		return nearByCabs, errors.New("distance can'nt be negative")
	}

	if !cabMap.isThereAWay(l) {
		return nearByCabs, errors.New("invalid location")
	}

	type qItem struct {
		l *locations.Location
		d int
	}

	//directions
	dirs := [][]int{
		{1, 0},  //up
		{0, 1},  //right
		{-1, 0}, //down
		{0, -1}, //left
	}

	//visited location
	vis := make([][]bool, cabMap.r)

	for i := range vis {
		vis[i] = make([]bool, cabMap.c)
	}

	//queue (for bfs)
	q := make(chan qItem, cabMap.r*cabMap.c)

	s := qItem{
		l: &locations.Location{
			Row: l.Row,
			Col: l.Col,
		},
		d: 0,
	}

	q <- s
	vis[s.l.Row][s.l.Col] = true

L:
	for {
		select {
		case node := <-q:
			if node.d > distance {
				continue
			}

			nearByCabs = append(nearByCabs, cabMap.allAvailableCabsAtLoc(node.l)...)

			for _, v := range dirs {

				neigh := qItem{
					l: &locations.Location{
						Row: node.l.Row + v[0],
						Col: node.l.Col + v[1],
					},
					d: node.d + 1,
				}
				if cabMap.isThereAWay(neigh.l) && !vis[neigh.l.Row][neigh.l.Col] && neigh.d <= distance {
					q <- neigh
					vis[neigh.l.Row][neigh.l.Col] = true
				}
			}

		default:
			break L
		}

	}

	if len(nearByCabs) == 0 {
		return nearByCabs, errors.New("No cabs available nearby you")
	}

	return nearByCabs, nil
}

func (cm *cabMaps) isThereAWay(l *locations.Location) bool {
	return l.Row >= 0 && l.Row < cm.r && l.Col >= 0 && l.Col < cm.c && maps.MapObj.Grid[l.Row][l.Col] == 0
}

// give me all cabs at location loc
func (cm *cabMaps) allCabsAtLoc(loc *locations.Location) []int {

	cabs := make([]int, 0)

	cabMap.mutex.Lock()
	for _, v := range cm.grid[loc.Row][loc.Col] {
		if v == -1 {
			continue
		}
		cabs = append(cabs, v)
	}
	cabMap.mutex.Unlock()

	return cabs
}

// give me all cabs at location loc which are available
func (cm *cabMaps) allAvailableCabsAtLoc(loc *locations.Location) []int {

	cabs := make([]int, 0)

	cabMap.mutex.Lock()
	for _, v := range cm.grid[loc.Row][loc.Col] {
		if v == -1 || !AllCabs[v].IsAvailable {
			continue
		}
		cabs = append(cabs, v)
	}
	cabMap.mutex.Unlock()

	return cabs
}
