package types

type LoginInput struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}
type UserData struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
