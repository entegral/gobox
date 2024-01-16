package dynamo

type User struct {
	Row
	Email string
	Name  string
	Age   int
}

type ErrMissingEmail struct{}

func (e ErrMissingEmail) Error() string {
	return "missing email"
}

// Keys returns the partition key and sort key for the row
func (u *User) Keys(gsi int) (string, string, error) {
	if u.Email == "" {
		return "", "", ErrMissingEmail{}
	}
	// For this example, assuming GUID is the partition key and Email is the sort key.
	// Additional logic can be added to handle different GSIs if necessary.
	u.Pk = u.Email
	u.Sk = "info"
	switch gsi {
	default:
		// Handle other GSIs or return an error
		return u.Pk, u.Sk, nil
	}
}

func (u *User) Type() string {
	return "user"
}

func CreateUser(email string) *User {
	return &User{
		Email: email,
	}
}

type ContactInfo struct {
	*MonoLink[*User]
	Phone string
	Addr  string
}

func (c *ContactInfo) Type() string {
	return "contact"
}
