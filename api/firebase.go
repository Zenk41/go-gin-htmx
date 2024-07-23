package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type FirebaseApi interface {
	SignInWithPassword(email, password string) (map[string]interface{}, error)
	SignUpWithPassword(name, email, password string) (map[string]interface{}, error)
	ExchangeRefreshTokenForIDToken(refreshToken string) (map[string]interface{}, error)
}

type firebaseApi struct {
	apiKey string
}

func NewFirebaseApi(apiKey string) FirebaseApi {
	return &firebaseApi{
		apiKey: apiKey,
	}
}

func (fi *firebaseApi) SignInWithPassword(email, password string) (map[string]interface{}, error) {
	loginPayload := map[string]interface{}{
		"email":             email,
		"password":          password,
		"returnSecureToken": true,
	}

	return fi.callFirebaseAuthAPI("signInWithPassword", loginPayload)
}

func (fi *firebaseApi) SignUpWithPassword(name, email, password string) (map[string]interface{}, error) {
	signUpPayload := map[string]interface{}{
		"email":             email,
		"password":          password,
		"returnSecureToken": true,
		"displayName":       name,
	}

	return fi.callFirebaseAuthAPI("signUp", signUpPayload)
}

func (fi *firebaseApi) ExchangeRefreshTokenForIDToken(refreshToken string) (map[string]interface{}, error) {
	tokenPayload := map[string]interface{}{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
	}

	return fi.callFirebaseAuthAPI("token", tokenPayload)
}

func (fi *firebaseApi) callFirebaseAuthAPI(endpoint string, payload map[string]interface{}) (map[string]interface{}, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.New("failed to marshal request payload: " + err.Error())
	}

	url := "https://identitytoolkit.googleapis.com/v1/accounts:" + endpoint + "?key=" + fi.apiKey
	if endpoint == "token" {
		url = "https://securetoken.googleapis.com/v1/" + endpoint + "?key=" + fi.apiKey
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, errors.New("failed to call Firebase Auth API: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, errors.New("failed to decode Firebase Auth API error response: " + err.Error())
		}
		return nil, errors.New("Firebase Auth API error: " + errorResponse["error"].(map[string]interface{})["message"].(string))
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.New("failed to decode Firebase Auth API response: " + err.Error())
	}

	return response, nil
}
