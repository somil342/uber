package auth

import (
	"errors"

	"github.com/somil342/uber/models"
	"github.com/somil342/uber/myobjects/riders"
)

//GetRider ...for authenticating user
func GetRider(userName string) (riders.Rider, error) {
	rider, err := models.SelectRider(userName)
	if err != nil {
		return riders.Rider{}, err
	}
	if rider.UserName != userName {
		return riders.Rider{}, errors.New("user invalid")
	}
	r := riders.New(rider.RiderID, rider.UserName)
	return *r, err
}
