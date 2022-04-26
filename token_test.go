package main

import (
	"encoding/json"
	"testing"
)

func TestGenerateUserToken(t *testing.T) {
	token, err := GenerateUserToken(User{0, "test", "test"})
	if err != nil {
		t.Error(err)
	}

	err = ValidateUserToken(string(token))
	if err != nil {
		t.Error(err)
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	token, err := GenerateRefreshToken(User{0, "test", "test"})
	if err != nil {
		t.Error(err)
	}

	_, err = ValidateRefreshToken(token)
	if err != nil {
		t.Error(err)
	}
}

func TestGenerateTokenPair(t *testing.T) {
	tokens, err := GenerateTokenPair(User{0, "Test", "Test"})
	if err != nil {
		t.Error(err)
	}

	accessTokenError := ValidateUserToken(string(tokens.AccessToken))
	_, refreshTokenError := ValidateRefreshToken(tokens.RefreshToken)
	if accessTokenError != nil || refreshTokenError != nil {
		t.Error(accessTokenError.Error() + "\n" + refreshTokenError.Error())
	}
}

func TestTokenPairJSON(t *testing.T) {
	tokens, err := GenerateTokenPair(User{0, "Test", "$2a$10$4H.g0FZ5Lcuy5NxaSr5fLOvunWZaTnplGcBVl7igsccuqftXZdDJu}"})
	if err != nil {
		t.Error(err)
	}

	tokensJSON, err := json.Marshal(tokens)
	if err != nil {
		t.Error(err)
	}

	var tokensDecoded TokenPair
	err = json.Unmarshal(tokensJSON, &tokensDecoded)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", string(tokensJSON))
}
