package main

import (
	"github.com/gin-gonic/gin"
	"github.com/somil342/uber/myobjects/cabs"
	"github.com/somil342/uber/myobjects/maps"
	"github.com/somil342/uber/routers"
)

func main() {
	//load cabs from DB in memory in a map
	cabs.LoadCabs()

	//load map from a file "map.txt" in memory in a 2D slice
	maps.LoadMap("")
	startServer()
}

func startServer() {
	router := gin.Default()
	//http://localhost:8080/pastbooking
	//header-> username
	router.GET("/pastbooking", routers.PastBooking)
	//http://localhost:8080/cabsnearbyme?loc=1,2&distance=10
	router.GET("/cabsnearbyme", routers.CabsNearByMe)

	//http://localhost:8080/bookcab
	//header-> username
	//json->{"fromrow":0,"fromcol":1,"torow":0,"tocol":4}
	router.POST("/bookcab", routers.BookCab)
	router.Run()
}
