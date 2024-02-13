package row

import "context"

var tstCtx = context.Background()

type User struct {
	Email string `dynamodbav:"email" json:"email"`
	Name  string `dynamodbav:"name" json:"name"`
}

func (u *User) Keys(gsi int) (string, string, error) {
	switch gsi {
	case 0:
		return u.Email, "details", nil
	case 1:
		return u.Name, "details", nil
	}
	return "", "", nil
}

func (u *User) Type() string {
	return "User"
}
