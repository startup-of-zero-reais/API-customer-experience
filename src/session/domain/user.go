package domain

type (
	// User struct represents a user entity
	User struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		Password string `json:"-"`
	}

	UserRepository interface {
		Find(email string) (*User, error)
	}
)
