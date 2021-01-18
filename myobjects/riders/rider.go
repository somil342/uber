package riders

//Rider ...
type Rider struct {
	RiderID  int    `json:"riderid"`
	UserName string `json:"username"`
}

//New ...
func New(id int, name string) *Rider {
	return &Rider{
		RiderID:  id,
		UserName: name,
	}
}
