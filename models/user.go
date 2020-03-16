package models

// User provides a regular struct to get JWT, and it will save in db.
type User struct {
	UserName string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password,omitempty"`
}

// JwtToken provides authority to access service
type JwtToken struct {
	Token string `json:"token"`
}
