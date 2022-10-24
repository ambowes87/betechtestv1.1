package data

type UserData struct {
	ID        string `json="id"`
	FirstName string `json="first_name"`
	LastName  string `json="last_name"`
	Nickname  string `json="nickname"`
	Password  string `json="password"`
	Email     string `json="email"`
	Country   string `json="country"`
	Created   string `json="created_at"`
	Updated   string `json="updated_at"`
}
