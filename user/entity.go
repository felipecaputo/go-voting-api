package user

type User struct {
	ID       string `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name,omitempty"`
	IsAdmin  bool   `json:"is_admin,omitempty" db:"is_admin"`
}
