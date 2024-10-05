package request

type CreateUserRequest struct {
	FirstName       string `validate:"required" json:"first_name"`
	LastName        string `validate:"required" json:"last_name"`
	Email           string `validate:"required,email" json:"email"`
	Password        string `validate:"required,min=8" json:"password"`
	ConfirmPassword string `validate:"required,min=8" json:"confirm_password"`
}
