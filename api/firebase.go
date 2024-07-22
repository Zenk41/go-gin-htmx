package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type FirebaseApi interface {
	SignInWithPassword(email, password string) (*http.Response, error)
	SignUpWithPassword(name, email, password string) (*http.Response, error)
	ExchangeRefreshTokenforIDToken(refreshToken string) (*http.Response, error)
}

type firebaseApi struct {
	apiKey string
}

func NewFirebaseApi(apiKey string) FirebaseApi { // Return FirebaseApi interface
	return &firebaseApi{
		apiKey: apiKey,
	}
}

func (fi *firebaseApi) SignInWithPassword(email, password string) (*http.Response, error) {
	// Call Firebase Auth API to verify the user credentials
	loginPayload := map[string]interface{}{
		"email":             email,
		"password":          password,
		"returnSecureToken": true,
	}

	payloadBytes, err := json.Marshal(loginPayload)
	if err != nil {
		return nil, errors.New("failed to marshal request payload")
	}

	url := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=" + fi.apiKey
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, errors.New("failed to call Firebase Auth API")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, errors.New("failed to decode Firebase Auth API response")
		}
		// Optionally, you can extract the specific error message from the errorResponse
		return nil, errors.New("Firebase Auth API error: " + resp.Status)
	}

	return resp, nil
}

func (fi *firebaseApi) SignUpWithPassword(name, email, password string) (*http.Response, error) {
	// // Call Firebase Auth API to create the user
	firebasePayload := map[string]interface{}{
		"name":              name,
		"email":             email,
		"password":          password,
		"returnSecureToken": true,
	}

	payloadBytes, err := json.Marshal(firebasePayload)
	if err != nil {
		return &http.Response{}, errors.New("failed to marshal request payload")
	}

	resp, err := http.Post("https://identitytoolkit.googleapis.com/v1/accounts:signUp?key="+fi.apiKey, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return &http.Response{}, errors.New("failed to call Firebase Auth API")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return &http.Response{}, errors.New("failed to decode firebase Auth api response")
		}
		// Optionally, you can extract the specific error message from the errorResponse
		return nil, errors.New("firebase auth api error: " + resp.Status)
	}

	var registerResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&registerResponse); err != nil {
		return &http.Response{}, errors.New("failed to decode firebase auth api response")
	}
	return resp, nil
}

func (fi *firebaseApi) ExchangeRefreshTokenforIDToken(refreshToken string) (*http.Response, error) {
	firebasePayload := map[string]interface{}{
		"grand_type":              "refresh_token",
		"refresh_token":            refreshToken,
	}

	payloadBytes, err := json.Marshal(firebasePayload)
	if err != nil {
		return &http.Response{}, errors.New("failed to marshal request payload")
	}
	resp, err := http.Post("https://securetoken.googleapis.com/v1/token?key="+fi.apiKey,"application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return &http.Response{}, errors.New("failed to call Firebase Auth API")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return &http.Response{}, errors.New("failed to decode firebase Auth api response")
		}
		// Optionally, you can extract the specific error message from the errorResponse
		return nil, errors.New("firebase auth api error: " + resp.Status)
	}

	var tokenResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return &http.Response{}, errors.New("failed to decode firebase auth api response")
	}
	return resp, nil
}
