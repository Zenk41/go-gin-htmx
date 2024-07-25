package models

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the application
type User struct {
	UserID    string    `firestore:"user_id"`
	Email     string    `firestore:"email"`
	Password  string    `firestore:"password"` // Ensure this is hashed before storing
	Name      string    `firestore:"name"`
	CreatedAt time.Time `firestore:"created_at"`
	UpdatedAt time.Time `firestore:"updated_at"`
}

func (u *User) EncryptPassword (password string) error {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(pass)
	return nil
}

func (u *User) CheckPassword(encryptedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(encryptedPassword))
	if err != nil {
		return err
	}
	return nil
}

type RegisterPayload struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Password string `json:"password"`
}

// LogInResponseWithEmailPassword represents the response payload for a login request
type LogInResponseWithEmailPassword struct {
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
}

// RegisterResponseWithEmailPassword represents the response payload for a register request
type RegisterResponseWithEmailPassword struct {
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	LocalID      string `json:"localId"`
}

type userRepository struct {
	client *firestore.Client
}

type UserRepository interface {
	CreateUser(ctx context.Context, user User) error
	GetUser(ctx context.Context, userID string) (*User, error)
}

func NewUserRepository(client *firestore.Client) UserRepository {
	return &userRepository{
		client: client,
	}
}

// CreateUser creates a new user in Firestore
func (ur *userRepository) CreateUser(ctx context.Context, user User) error {
	_, err := ur.client.Collection("users").Doc(user.UserID).Set(ctx, user)
	return err
}

// GetUser retrieves a user by ID from Firestore
func (ur *userRepository) GetUser(ctx context.Context, userID string) (*User, error) {
	doc, err := ur.client.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		return nil, err
	}
	var user User
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
