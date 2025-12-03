package user

type UserRepositoryInterface interface {
	Create(dto *User) (*User, error)
	GetOne(dto *UserRequestDto) (*User, error)
	Update(userID *string, data *User) (*User, error)
}
