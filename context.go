package poeia

import "context"

func NewContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, "poeia_user", user)
}

func GetUserFromContext(ctx context.Context) *User {
	user, _ := ctx.Value("poeia_user").(*User)
	return user
}
