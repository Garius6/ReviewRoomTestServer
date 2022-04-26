package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTokenPair(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/auth/token?username=rest&password=rest", nil)
	w := httptest.NewRecorder()
	getTokenPair(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Error(res.StatusCode)
	}

	defer res.Body.Close()
	var tokenPair TokenPair
	err := json.NewDecoder(res.Body).Decode(&tokenPair)
	if err != nil {
		t.Fatal(err)
	}

	accessTokenError := ValidateUserToken(string(tokenPair.AccessToken))
	_, refreshTokenError := ValidateRefreshToken(tokenPair.RefreshToken)
	if accessTokenError != nil {
		t.Error("access token error" + accessTokenError.Error())
	}

	if refreshTokenError != nil {
		t.Error("refresh token error" + refreshTokenError.Error())
	}
}
