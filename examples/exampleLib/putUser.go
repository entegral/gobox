package exampleLib

import "context"

func PutUser(ctx context.Context) (user *User) {
	email, name, age := "test@gmail.com", "Test User Name", 30
	user = &User{
		Email: email,
		Name:  name,
		Age:   age,
	}
	err := user.Put(ctx, user)
	if err != nil {
		panic(err)
	}
	return user
}
