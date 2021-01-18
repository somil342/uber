package trips

import (
	"encoding/json"
	"time"

	"github.com/somil342/uber/database"
	"github.com/somil342/uber/models"
	"github.com/somil342/uber/myobjects/cabs"
	"github.com/somil342/uber/myobjects/locations"
	"github.com/somil342/uber/myobjects/maps"
	"github.com/somil342/uber/myobjects/riders"
)

//Trip ...
type Trip struct {
	TripID        int                `json:"tripid,omitempty"`
	Rider         riders.Rider       `json:"rider,omitempty"`
	Cab           cabs.Cab           `json:"cab,omitempty"`
	TripStatus    string             `json:"tripstatus,omitempty"`
	FromPt        locations.Location `json:"frompt,omitempty"`
	ToPt          locations.Location `json:"topt,omitempty"`
	Path          []string           `json:"path,omitempty"`
	Price         int                `json:"price,omitempty"`
	Eta           int                `json:"eta_in_minutes,omitempty"`
	TripStartTime time.Time          `json:"trip_start_time,omitempty"`
}

//Trips ...
var Trips map[int]*Trip = make(map[int]*Trip)

//GetTripData ...
func GetTripData(tripID int) (string, error) {
	if tripID == -1 {
		b, err := json.Marshal(&Trips)
		return string(b), err
	}
	trip := Trips[tripID]
	b, err := json.Marshal(&trip)
	return string(b), err
}

//creating intance for trip
func new(rider riders.Rider, cab cabs.Cab, from locations.Location, to locations.Location, path []string, price int, eta int) *Trip {

	trip := &Trip{
		Rider:      rider,
		Cab:        cab,
		TripStatus: "CREATED",
		FromPt:     from,
		ToPt:       to,
		Path:       path,
		Price:      price,
		Eta:        eta,
	}

	return trip
}

//CreateTrip ...
func CreateTrip(rider riders.Rider, from *locations.Location, to *locations.Location, defaultCab ...int) (*Trip, error) {

	cabID := -1

	if len(defaultCab) > 0 && defaultCab[0] >= 0 {
		cb, ok := cabs.AllCabs[defaultCab[0]]
		if ok && cb.IsAvailable {
			cabID = defaultCab[0]
		}
	}

	if cabID == -1 {
		nearByCabs, err := cabs.FindCabsNearBy(from, 1000000)
		if err != nil {
			return nil, err
		}
		cabID = nearByCabs[0]
	}
	cb := cabs.AllCabs[cabID]

	cb.IsAvailable = false

	path, price, etaMin, err := maps.MapObj.Bfs(from, to)
	if err != nil {
		return nil, err
	}

	trip := new(rider, *cb, *from, *to, path, price, etaMin)

	return trip, nil
}

//StartTrip ...
func (t *Trip) StartTrip() error {
	t.TripStartTime = time.Now().Local()
	t.TripStatus = "STARTED"

	//insert in db
	conn, err := database.Connect()
	if err != nil {
		return err
	}

	id, err := models.InsertTrip(conn, t.FromPt.String(), t.ToPt.String(), t.Price, "STARTED", t.Path, t.Cab.CabID, t.Rider.RiderID)

	if err == nil {
		t.TripID = id
		Trips[id] = t
	}
	return err
}

//EndTrip ...
func (t *Trip) EndTrip(status string) error {
	t.TripStatus = status

	//update in db
	Trips[t.TripID] = nil

	conn, err := database.Connect()
	if err != nil {
		return err
	}
	err = models.EditTrip(conn, t.TripID, status)
	if err != nil {
		return nil
	}
	cabs.AllCabs[t.Cab.CabID].IsAvailable = true
	return nil
}
