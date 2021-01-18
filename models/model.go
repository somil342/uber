package models

import (
	"fmt"
	"time"

	pg "github.com/go-pg/pg/v10"
	orm "github.com/go-pg/pg/v10/orm"
	"github.com/somil342/uber/database"
)

//Cab ...
type Cab struct {
	tableName  string    `pg:"cabs"`
	CabID      int       `pg:"cab_id,pk,unique"`
	DriverName string    `json:"driver_name" pg:"driver_name"`
	CarName    string    `json:"car_name" pg:"car_name"`
	Email      string    `json:"email" pg:"email,notnull,unique"`
	PhoneNum   string    `json:"phone_num" pg:"phone_num"`
	CreatedOn  time.Time `json:"created_on" pg:"created_on,default:now()"`
}

//Rider ...
type Rider struct {
	tableName string `pg:"riders"`
	RiderID   int    `json:"rider_id" pg:"rider_id,pk,unique"`
	UserName  string `json:"user_name" pg:"user_name,unique"`
}

//Trip ...
type Trip struct {
	tableName     string    `pg:"trips"`
	TripID        int       `pg:"trip_id,pk,unique"`
	TripStartTime time.Time `json:"trip_start_time" pg:"trip_start_time,default:now()"`
	From          string    `json:"from" pg:"from"`
	To            string    `json:"to" pg:"to"`
	Price         int       `json:"price" pg:"price"`
	TripStatus    string    `json:"trip_status" pg:"trip_status"`
	Path          []string  `json:"path" pg:"path,array"`

	Cab     *Cab   `json:"cab" pg:"cab,rel:has-one"`
	CabID   int    `json:"cab_id" pg:"cab_cab_id,fk,notnull,on_delete:CASCADE"`
	Rider   *Rider `json:"rider" pg:"rel:has-one"`
	RiderID int    `json:"rider_id" pg:"rider_rider_id,fk,notnull,on_delete:CASCADE"`
}

//SelectAllTripsOfARider ...
func SelectAllTripsOfARider(riderID int) ([]Trip, error) {

	var trip []Trip

	conn, err := database.Connect()
	if err != nil {
		return trip, err
	}
	defer conn.Close()

	q := conn.Model(&trip).Relation("Cab").Where("rider_id=?0", riderID).Relation("Rider")

	err = q.Select()
	if err != nil {
		return trip, err
	}

	return trip, nil
}

//insert dummy cabs
func insertDummyCabs(conn *pg.DB) {

	dummyCabs := []Cab{
		{DriverName: "mohan lal", CarName: "santro", Email: "mohan@gmail.com,", PhoneNum: "11111"},
		{DriverName: "sohan lal", CarName: "swift", Email: "sohan@gmail.com,", PhoneNum: "22222"},
		{DriverName: "riya yadav", CarName: "accent", Email: "riya@gmail.com,", PhoneNum: "33333"},
		{DriverName: "rakesh kumar", CarName: "suzuki", Email: "rakesh@gmail.com,", PhoneNum: "44444"},
	}
	_, err := conn.Model(&dummyCabs).Insert()
	if err != nil {
		panic(err)
	}
}

//insert dummy riders
func insertDummyRiders(conn *pg.DB) {

	dummyRiders := []Rider{
		{UserName: "abhay"},
		{UserName: "jatin"},
		{UserName: "sahil"},
	}
	_, err := conn.Model(&dummyRiders).Insert()
	if err != nil {
		panic(err)
	}
}

//creating needed tables and inserting dummy data for demo purpose
func init() {
	//fmt.Println("creating tables")
	conn, err := database.Connect()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for _, model := range []interface{}{&Cab{}, &Rider{}, &Trip{}} {

		err = conn.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
			//Temp:        true,
			FKConstraints: true,
		})
		if err != nil {
			panic(err)
		}
	}

	//insertDummyRiders(conn)
	//insertDummyCabs(conn)

}

//SelectAllCabs ...
func SelectAllCabs() ([]Cab, error) {

	var cabs []Cab
	conn, err := database.Connect()
	if err != nil {
		return cabs, err
	}
	defer conn.Close()

	err = conn.Model(&cabs).Select()
	if err != nil {
		return cabs, err
	}

	return cabs, nil
}

//InsertTrip ...
func InsertTrip(conn *pg.DB, from string, to string, price int, tripStatus string, path []string, CabCabID int, RiderRiderID int) (int, error) {

	trip := &Trip{
		From:       from,
		To:         to,
		Price:      price,
		TripStatus: tripStatus,
		Path:       path,
		CabID:      CabCabID,
		RiderID:    RiderRiderID,
	}

	_, err := conn.Model(trip).Returning("trip_id").Insert()
	if err != nil {
		return -1, err
	}

	return trip.TripID, nil
}

//EditTrip ...
func EditTrip(conn *pg.DB, tripid int, tripStatus string) error {

	fmt.Println(tripid, tripStatus)
	trip := &Trip{
		TripID: tripid,
	}
	err := conn.Model(trip).WherePK().Select()
	if err != nil {
		return err
	}

	//fmt.Println(trip.TripID)

	trip.TripStatus = tripStatus
	_, err = conn.Model(trip).WherePK().Update()
	if err != nil {
		return err
	}

	return nil
}

//SelectRider ...
func SelectRider(userName string) (Rider, error) {

	var rider Rider
	conn, err := database.Connect()
	if err != nil {
		return rider, err
	}
	defer conn.Close()

	err = conn.Model(&rider).Where("user_name = ?", userName).Select()
	if err != nil {
		return rider, err
	}

	return rider, nil
}
