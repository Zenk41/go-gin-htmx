package firebase

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

func Auth(serviceKeyFile string) (*auth.Client, error) {
    ctx := context.Background()
    app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(serviceKeyFile))
    if err != nil {
        return nil, err
    }

    authClient, err := app.Auth(ctx)
    if err != nil {
        return nil, err
    }

    return authClient, nil
}