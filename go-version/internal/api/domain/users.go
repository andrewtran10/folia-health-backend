package domain

type UserCreateDomain struct {
	Name     *string
	Email    *string
	Password *string
}

type UserUpdateDomain struct {
	UserID   string
	Name     *string
	Email    *string
	Password *string
}

type UserGetDomain struct {
	UserID string
}
