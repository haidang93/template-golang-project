package user

type UserRequestDto struct {
	ArrayEmail  *[]string `json:"arrayEmail"`
	ArrayID     *[]string `json:"arrayId"`
	ArrayStatus *[]string `json:"arrayStatus"`
	Page        *int      `json:"page"`
	Limit       *int      `json:"limit"`
}

type UserUpdateDto struct {
	DisplayName       *string `json:"displayName"`
	FirstName         *string `json:"firstName"`
	LastName          *string `json:"lastName"`
	PhoneNumber       *string `json:"phoneNumber"`
	Bio               *string `json:"bio"`
	PreferredLanguage *string `json:"preferredLanguage"`
	ReceiveEmail      *bool   `json:"receiveEmail"`
}
