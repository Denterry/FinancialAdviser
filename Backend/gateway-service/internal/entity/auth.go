package entity

// Claims is what our middleware pulls out of a validated JWT
type Claims struct {
	UserID   string
	Email    string
	Username string
	IsAdmin  bool
}
