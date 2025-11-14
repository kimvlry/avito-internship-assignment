package entity

type User struct {
	ID       string
	Username string
	TeamName string
	IsActive bool
}

func (u *User) CanReview() bool {
	return u.IsActive
}
