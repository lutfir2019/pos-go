package types

type User struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type UpdateUser struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Username        string `json:"username"`
	Role            string `json:"role"`
	CurrentPassword string `json:"currentPassword"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"passwordConfirm"`
}
