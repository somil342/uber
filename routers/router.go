package routers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/somil342/uber/auth"
	"github.com/somil342/uber/models"
	"github.com/somil342/uber/myobjects/cabs"
	"github.com/somil342/uber/myobjects/locations"
	"github.com/somil342/uber/myobjects/trips"
)

//for authenticating rider
type header struct {
	UserName string `header:"username"`
}

type cabbook struct {
	FromRow int `json:"fromrow"`
	FromCol int `json:"fromcol"`
	ToRow   int `json:"torow"`
	ToCol   int `json:"tocol"`
}

//BookCab ...
func BookCab(c *gin.Context) {
	/*{
		"fromrow":0,
		"fromcol":1,
		"torow":0,
		"tocol":4
	}*/

	h := header{}
	if err := c.ShouldBindHeader(&h); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	rider, err := auth.GetRider(h.UserName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	data := cabbook{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	src := &locations.Location{Row: data.FromRow, Col: data.FromCol}
	dest := &locations.Location{Row: data.ToRow, Col: data.ToCol}

	trip, err := trips.CreateTrip(rider, src, dest)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	err = trip.StartTrip()

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"trip": trip})

	go func() {
		time.Sleep(time.Duration(trip.Eta) * time.Millisecond)
		trip.EndTrip("COMPLETED")
	}()

}

//CabsNearByMe ...
func CabsNearByMe(c *gin.Context) {

	mp := c.Request.URL.Query()

	locStr, ok := mp["loc"]

	if !ok || len(locStr) == 0 {
		c.JSON(200, gin.H{"err": "invalid location"})
		return
	}

	res := strings.Split(locStr[0], ",")
	if len(res) != 2 {
		c.JSON(200, gin.H{"err": "invalid location"})
		return
	}

	row, err1 := strconv.Atoi(res[0])
	col, err2 := strconv.Atoi(res[1])

	if err1 != nil || err2 != nil {
		c.JSON(200, gin.H{"err": "invalid location"})
		return
	}

	loc := locations.Location{Row: row, Col: col}

	disStr, ok := mp["distance"]
	dis := 10
	if ok && len(disStr) > 0 {
		val, _ := strconv.Atoi(disStr[0])
		if val > 0 {
			dis = val
		}
	}
	nearByCabsIDs, err := cabs.FindCabsNearBy(&loc, dis)

	if err != nil {
		c.JSON(200, gin.H{"err": err.Error()})
		return
	}

	nearByCabs := make([]cabs.Cab, 0)
	for _, v := range nearByCabsIDs {
		nearByCabs = append(nearByCabs, cabs.GetCabData(v))
	}

	c.JSON(200, gin.H{"nearByCabs": nearByCabs})
}

//PastBooking ...
func PastBooking(c *gin.Context) {

	h := header{}
	if err := c.ShouldBindHeader(&h); err != nil {
		c.JSON(200, gin.H{"err": err.Error()})
		return
	}

	r, err := auth.GetRider(h.UserName)
	if err != nil {
		c.JSON(200, gin.H{"err": err.Error()})
		return
	}

	trips, err := models.SelectAllTripsOfARider(r.RiderID)
	if err != nil {
		c.JSON(200, err)
		return
	}

	c.JSON(200, gin.H{"trips": trips})
}
