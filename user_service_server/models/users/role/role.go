package modelRole

//Role model
type Role struct {
	ID   int    `bson:"id"`
	Name string `bson:"name"`
}

//USER role
var USER = Role{ID: 1, Name: "User"}

//ADMIN role
var ADMIN = Role{ID: 99, Name: "Enterprise Administrator"}

//New - Creates new role
func New(id int, name string) Role {
	return Role{ID: id, Name: name}
}

//GetID - Get ID of Role
func (r *Role) GetID() int {
	return r.ID
}

//SetID - Set ID of Role
func (r *Role) SetID(id int) {
	r.ID = id
}

//GetName - Get Name of Role
func (r *Role) GetName() string {
	return r.Name
}

//SetName - Set Name of Role
func (r *Role) SetName(name string) {
	r.Name = name
}
