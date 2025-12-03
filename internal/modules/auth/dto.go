package auth

type SigninDto struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignupDto struct {
	FirstName    *string `json:"firstName" validate:"required"`
	LastName     *string `json:"lastName" validate:"required"`
	Email        *string `json:"email" validate:"required"`
	Username     *string `json:"username"`
	Password     *string `json:"password" validate:"required"`
	AcceptPolicy *bool   `json:"acceptPolicy" validate:"required"`
	ReceiveEmail *bool   `json:"receiveEmail" validate:"required"`
}
