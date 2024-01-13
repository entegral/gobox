package dynamo

import "context"

func PutUser(ctx context.Context, email string) (user *User) {
	if email == "" {
		email = "test@gmail.com"
	}
	name, age := "Test User Name", 30
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
