package request

type (
	Login struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	ChangePassword struct {
		PasswordBaru string `json:"password_baru" validate:"required"`
	}
)
