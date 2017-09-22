package modelUser

//User model
type User struct {
	FirstName     string
	MiddleInitial string
	LastName      string
	Username      string
	Password      string
	Status        int
	TimeLocked    int32
	Email         string
	PhoneNumber   string
}
