package types

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Register struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"passwordConfirm"`
	Role            string `json:"role"`
}
