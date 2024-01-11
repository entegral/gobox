package exampleLib

import "github.com/entegral/gobox/dynamo"

type User struct {
	dynamo.Row
	Email string
	Name  string
	Age   int
}

// Keys returns the partition key and sort key for the row
func (u *User) Keys(gsi int) (string, string, error) {
	// For this example, assuming GUID is the partition key and Email is the sort key.
	// Additional logic can be added to handle different GSIs if necessary.
	switch gsi {
	default:
		// Handle other GSIs or return an error
		return u.Email, "info", nil
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
	*dynamo.MonoLink[*User]
	Phone string
	Addr  string
}

func (c *ContactInfo) Type() string {
	return "contact"
}
