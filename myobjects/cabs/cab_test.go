package cabs

import (
	"testing"

	"github.com/somil342/uber/myobjects/maps"

	"github.com/somil342/uber/models"
	"github.com/somil342/uber/myobjects/locations"
)

func TestUpdateLocation(t *testing.T) {

	cab := &Cab{}

	r := -1
	c := 1
	err := cab.UpdateLocation(locations.Location{Row: r, Col: c})

	if cab == nil {
		t.Errorf("TestUpdateLocation() = nil; want invalid location")
	}

	r = 1
	c = 1
	err = cab.UpdateLocation(locations.Location{Row: r, Col: c})

	if err != nil {
		t.Errorf("TestUpdateLocation() = %s; want nil", err.Error())
	}

	if cab.CurLocation.Row != r || cab.CurLocation.Col != c {
		t.Errorf("TestUpdateLocation() = %d,%d; want %d,%d", cab.CurLocation.Row, cab.CurLocation.Col, r, c)
	}

	ok := false
	for _, v := range cabMap.grid[r][c] {
		if v == cab.CabID {
			ok = true
			break
		}
	}

	if !ok {
		t.Errorf("TestUpdateLocation() = cab loc not updated on map; want %d", cab.CabID)
	}

}

func TestNew(t *testing.T) {

	if len(AllCabs) != 0 {
		t.Errorf("TestNew() = %d; want 0", len(AllCabs))
	}

	c, _ := New(1, "abhay", locations.Location{Row: 1, Col: 2})
	if c == nil {
		t.Errorf("TestNew() = %s; want cab{1,'abhay',{1,2}}", "nil")
	}

	if len(AllCabs) != 1 {
		t.Errorf("TestNew() = %d; want 1", len(AllCabs))
	}

	New(2, "jatin", locations.Location{Row: 3, Col: 4})

	if len(AllCabs) != 2 {
		t.Errorf("TestNew() = %d; want 2", len(AllCabs))
	}

}

func TestFindCabsNearBy(t *testing.T) {

	var err error
	err = LoadCabs()

	if err != nil {
		t.Errorf("TestFindCabsNearBy() = non-nil; want nil")
	}

	_, err = FindCabsNearBy(&locations.Location{Row: 0, Col: 3}, -1)
	if err == nil {
		t.Errorf("TestFindCabsNearBy() = nil; want distance can'nt be negative")
	}

	r := 0
	c := 1

	_, err = FindCabsNearBy(&locations.Location{Row: r, Col: c}, 0)

	if err != nil {
		t.Errorf("TestFindCabsNearBy() = %s; want nil", err.Error())
	}

	want := 0
	for _, v := range AllCabs {
		if v.CurLocation.Row == r && v.CurLocation.Col == c {
			want++
		}
	}
	cabs, _ := FindCabsNearBy(&locations.Location{Row: r, Col: c}, 0)
	got := len(cabs)

	if want != got {
		t.Errorf("TestFindCabsNearBy() = %d; want %d", got, want)
	}

	TestUnLoadCabs(t)

}

func TestLoadCabs(t *testing.T) {

	cabs, err := models.SelectAllCabs()
	if err != nil {
		panic(err)
	}

	LoadCabs()

	if len(AllCabs) != len(cabs) {
		t.Errorf("TestLoadCabs = %d; want %d", len(AllCabs), len(cabs))
	}

	TestUnLoadCabs(t)
}

func TestIsThereAWay(t *testing.T) {

	if cabMap.isThereAWay(&locations.Location{Row: 0, Col: -1}) || cabMap.isThereAWay(&locations.Location{Row: -1, Col: 0}) {
		t.Errorf("TestIsThereAWay = true; want false")
	}

	if cabMap.isThereAWay(&locations.Location{Row: maps.MapObj.Row, Col: 0}) || cabMap.isThereAWay(&locations.Location{Row: 0, Col: maps.MapObj.Col}) {
		t.Errorf("TestIsThereAWay = true; want false")
	}

	maps.MapObj.Grid[0][0] = 1

	if cabMap.isThereAWay(&locations.Location{Row: 0, Col: 0}) {
		t.Errorf("TestIsThereAWay = true; want false")
	}
	maps.MapObj.Grid[0][0] = 0
}

func TestAllCabsAtLoc(t *testing.T) {
	loc := locations.Location{Row: 0, Col: 0}
	c1, err := New(1, "a", loc)
	if err != nil {
		t.Errorf("TestAllCabsAtLoc = %s; want nil", err.Error())
	}
	c1.IsAvailable = false
	_, err = New(2, "b", loc)
	if err != nil {
		t.Errorf("TestAllCabsAtLoc = %s; want nil", err.Error())
	}

	cabs := cabMap.allCabsAtLoc(&loc)
	got := len(cabs)
	want := 2

	if got != want {
		t.Errorf("TestAllCabsAtLoc = %d; want %d", got, want)
	}

	TestUnLoadCabs(t)

}

func TestAllAvailableCabsAtLoc(t *testing.T) {
	loc := locations.Location{Row: 0, Col: 0}
	c1, err := New(1, "a", loc)
	if err != nil {
		t.Errorf("TestAllAvailableCabsAtLoc = %s; want nil", err.Error())
	}
	c1.IsAvailable = false
	_, err = New(2, "b", loc)
	if err != nil {
		t.Errorf("TestAllAvailableCabsAtLoc = %s; want nil", err.Error())
	}

	cabs := cabMap.allAvailableCabsAtLoc(&loc)
	got := len(cabs)
	want := 1

	if got != want {
		t.Errorf("TestAllAvailableCabsAtLoc = %d; want %d", got, want)
	}

}

func TestUnLoadCabs(t *testing.T) {
	for i := 0; i < cabMap.r; i++ {
		for j := 0; j < cabMap.c; j++ {
			cabMap.grid[i][j] = make([]int, 0)
		}
	}
	AllCabs = make(map[int]*Cab)
}
