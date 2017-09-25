package modelStatus

//Status Model
type Status struct {
	ID   int    `bson:"id"`
	Name string `bson:"name"`
}

//UNVERIFIED status
var UNVERIFIED = Status{ID: 1, Name: "UNVERIFIED"}

//ACTIVE status
var ACTIVE = Status{ID: 2, Name: "ACTIVE"}

//PASSWORD_LOCK status
var PASSWORD_LOCK = Status{ID: 3, Name: "PASSWORD_LOCK"}

//ACCOUNT_LOCK status
var ACCOUNT_LOCK = Status{ID: 4, Name: "ACCOUNT_LOCK"}

//FROZEN status
var FROZEN = Status{ID: 5, Name: "FROZEN"}

//PASSWORD_RESET status
var PASSWORD_RESET = Status{ID: 6, Name: "PASSWORD_RESET"}

//New - Creates a new Status
func New(id int, name string) Status {
	return Status{ID: id, Name: name}
}

//GetID - Get ID of Status
func (s *Status) GetID() int {
	return s.ID
}

//SetID - Set ID of Status
func (s *Status) SetID(id int) {
	s.ID = id
}

//GetName - Get Name of Status
func (s *Status) GetName() string {
	return s.Name
}

//SetName - Set Name of Status
func (s *Status) SetName(name string) {
	s.Name = name
}
