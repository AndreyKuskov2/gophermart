package jwt

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

const testSecretKey = "test-secret-key"

func TestCreateJwtToken(t *testing.T) {
	userID := 123

	token, err := CreateJwtToken(testSecretKey, userID)
	if err != nil {
		t.Fatalf("CreateJwtToken failed: %v", err)
	}

	if token == "" {
		t.Fatal("CreateJwtToken returned empty token")
	}

	claims, err := VerifyToken(token, testSecretKey)
	if err != nil {
		t.Fatalf("Failed to verify created token: %v", err)
	}

	if claims.Subject != "123" {
		t.Errorf("Expected subject '123', got '%s'", claims.Subject)
	}

	if claims.ExpiresAt == nil {
		t.Fatal("Token has no expiration time")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		t.Fatal("Token has already expired")
	}
}

func TestCreateJwtTokenWithDifferentUserIDs(t *testing.T) {
	testCases := []int{1, 100, 999999}

	for _, userID := range testCases {
		t.Run(fmt.Sprintf("UserID_%d", userID), func(t *testing.T) {
			token, err := CreateJwtToken(testSecretKey, userID)
			if err != nil {
				t.Fatalf("CreateJwtToken failed for userID %d: %v", userID, err)
			}

			claims, err := VerifyToken(token, testSecretKey)
			if err != nil {
				t.Fatalf("Failed to verify token for userID %d: %v", userID, err)
			}

			expectedSubject := strconv.Itoa(userID)
			if claims.Subject != expectedSubject {
				t.Errorf("Expected subject '%s', got '%s'", expectedSubject, claims.Subject)
			}
		})
	}
}

func TestVerifyToken(t *testing.T) {
	userID := 456
	token, err := CreateJwtToken(testSecretKey, userID)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	claims, err := VerifyToken(token, testSecretKey)
	if err != nil {
		t.Fatalf("VerifyToken failed for valid token: %v", err)
	}

	if claims.Subject != "456" {
		t.Errorf("Expected subject '456', got '%s'", claims.Subject)
	}

	_, err = VerifyToken(token, "wrong-secret")
	if err == nil {
		t.Fatal("VerifyToken should fail with wrong secret key")
	}

	_, err = VerifyToken("invalid.token.here", testSecretKey)
	if err == nil {
		t.Fatal("VerifyToken should fail with malformed token")
	}

	_, err = VerifyToken("", testSecretKey)
	if err == nil {
		t.Fatal("VerifyToken should fail with empty token")
	}
}

func TestVerifyTokenWithExpiredToken(t *testing.T) {
	claims := JWTClaims{
		jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(-time.Second)),
			Subject:   "123",
		},
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(testSecretKey))
	if err != nil {
		t.Fatalf("Failed to create expired token: %v", err)
	}

	_, err = VerifyToken(tokenString, testSecretKey)
	if err == nil {
		t.Fatal("VerifyToken should fail with expired token")
	}
}

func TestGetJwtClaims(t *testing.T) {
	testClaims := &JWTClaims{
		jwtlib.RegisteredClaims{
			Subject: "789",
		},
	}

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("Failed to create test request: %v", err)
	}

	ctx := context.WithValue(req.Context(), "claims", testClaims)
	req = req.WithContext(ctx)

	claims, err := GetJwtClaims(req)
	if err != nil {
		t.Fatalf("GetJwtClaims failed: %v", err)
	}

	if claims.Subject != "789" {
		t.Errorf("Expected subject '789', got '%s'", claims.Subject)
	}

	reqWithoutClaims, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("Failed to create test request: %v", err)
	}

	_, err = GetJwtClaims(reqWithoutClaims)
	if err == nil {
		t.Fatal("GetJwtClaims should fail when no claims in context")
	}

	reqWrongType, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("Failed to create test request: %v", err)
	}

	ctxWrongType := context.WithValue(reqWrongType.Context(), "claims", "not-a-claim")
	reqWrongType = reqWrongType.WithContext(ctxWrongType)

	_, err = GetJwtClaims(reqWrongType)
	if err == nil {
		t.Fatal("GetJwtClaims should fail when context contains wrong type")
	}
}

func TestJWTClaimsStruct(t *testing.T) {
	claims := &JWTClaims{
		jwtlib.RegisteredClaims{
			Subject:   "test-subject",
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	if claims.Subject != "test-subject" {
		t.Errorf("Expected subject 'test-subject', got '%s'", claims.Subject)
	}

	if claims.ExpiresAt == nil {
		t.Fatal("ExpiresAt should not be nil")
	}
}
