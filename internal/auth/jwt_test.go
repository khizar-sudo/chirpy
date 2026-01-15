package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	t.Run("successfully creates a valid JWT token", func(t *testing.T) {
		userID := uuid.New()
		tokenSecret := "test-secret"
		expiresIn := time.Hour

		token, err := MakeJWT(userID, tokenSecret, expiresIn)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if token == "" {
			t.Fatal("expected token to be non-empty")
		}
	})

	t.Run("creates token with correct claims", func(t *testing.T) {
		userID := uuid.New()
		tokenSecret := "test-secret"
		expiresIn := time.Hour
		beforeCreation := time.Now().UTC()

		token, err := MakeJWT(userID, tokenSecret, expiresIn)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Parse the token to verify claims
		parsed, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(tokenSecret), nil
		})
		if err != nil {
			t.Fatalf("failed to parse token: %v", err)
		}

		claims, ok := parsed.Claims.(*jwt.RegisteredClaims)
		if !ok {
			t.Fatal("failed to get claims")
		}

		// Verify issuer
		if claims.Issuer != "chirpy" {
			t.Errorf("expected issuer to be 'chirpy', got %s", claims.Issuer)
		}

		// Verify subject (user ID)
		if claims.Subject != userID.String() {
			t.Errorf("expected subject to be %s, got %s", userID.String(), claims.Subject)
		}

		// Verify issued at time is recent (within 1 second tolerance)
		timeDiffIssued := claims.IssuedAt.Time.Sub(beforeCreation)
		if timeDiffIssued < 0 {
			timeDiffIssued = -timeDiffIssued
		}
		if timeDiffIssued > time.Second {
			t.Errorf("issued at time differs by more than 1 second: %v", timeDiffIssued)
		}

		// Verify expiration time is roughly correct
		expectedExpiry := beforeCreation.Add(expiresIn)
		timeDiff := claims.ExpiresAt.Time.Sub(expectedExpiry)
		if timeDiff < 0 {
			timeDiff = -timeDiff
		}
		if timeDiff > time.Second {
			t.Errorf("expiration time differs by more than 1 second: %v", timeDiff)
		}
	})

	t.Run("creates different tokens for different users", func(t *testing.T) {
		userID1 := uuid.New()
		userID2 := uuid.New()
		tokenSecret := "test-secret"
		expiresIn := time.Hour

		token1, err := MakeJWT(userID1, tokenSecret, expiresIn)
		if err != nil {
			t.Fatalf("failed to create token1: %v", err)
		}

		token2, err := MakeJWT(userID2, tokenSecret, expiresIn)
		if err != nil {
			t.Fatalf("failed to create token2: %v", err)
		}

		if token1 == token2 {
			t.Error("expected different tokens for different users")
		}
	})

	t.Run("creates token with custom expiration", func(t *testing.T) {
		userID := uuid.New()
		tokenSecret := "test-secret"
		expiresIn := 30 * time.Minute

		token, err := MakeJWT(userID, tokenSecret, expiresIn)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		parsed, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(tokenSecret), nil
		})
		if err != nil {
			t.Fatalf("failed to parse token: %v", err)
		}

		claims, ok := parsed.Claims.(*jwt.RegisteredClaims)
		if !ok {
			t.Fatal("failed to get claims")
		}

		actualDuration := claims.ExpiresAt.Time.Sub(claims.IssuedAt.Time)
		timeDiff := actualDuration - expiresIn
		if timeDiff < 0 {
			timeDiff = -timeDiff
		}
		if timeDiff > time.Second {
			t.Errorf("expiration duration differs from expected: expected %v, got %v", expiresIn, actualDuration)
		}
	})
}

func TestValidateJWT(t *testing.T) {
	t.Run("successfully validates a valid token", func(t *testing.T) {
		userID := uuid.New()
		tokenSecret := "test-secret"
		expiresIn := time.Hour

		token, err := MakeJWT(userID, tokenSecret, expiresIn)
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}

		validatedUserID, err := ValidateJWT(token, tokenSecret)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if validatedUserID != userID {
			t.Errorf("expected user ID %s, got %s", userID, validatedUserID)
		}
	})

	t.Run("fails validation with wrong secret", func(t *testing.T) {
		userID := uuid.New()
		tokenSecret := "test-secret"
		wrongSecret := "wrong-secret"
		expiresIn := time.Hour

		token, err := MakeJWT(userID, tokenSecret, expiresIn)
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}

		_, err = ValidateJWT(token, wrongSecret)
		if err == nil {
			t.Fatal("expected error when validating with wrong secret, got nil")
		}
	})

	t.Run("fails validation with expired token", func(t *testing.T) {
		userID := uuid.New()
		tokenSecret := "test-secret"
		expiresIn := -time.Hour // Token already expired

		token, err := MakeJWT(userID, tokenSecret, expiresIn)
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}

		_, err = ValidateJWT(token, tokenSecret)
		if err == nil {
			t.Fatal("expected error when validating expired token, got nil")
		}
	})

	t.Run("fails validation with malformed token", func(t *testing.T) {
		tokenSecret := "test-secret"
		malformedToken := "not.a.valid.jwt.token"

		_, err := ValidateJWT(malformedToken, tokenSecret)
		if err == nil {
			t.Fatal("expected error when validating malformed token, got nil")
		}
	})

	t.Run("fails validation with empty token", func(t *testing.T) {
		tokenSecret := "test-secret"

		_, err := ValidateJWT("", tokenSecret)
		if err == nil {
			t.Fatal("expected error when validating empty token, got nil")
		}
	})

	t.Run("fails validation with token containing invalid UUID", func(t *testing.T) {
		tokenSecret := "test-secret"
		expiresIn := time.Hour

		// Create a token with invalid UUID as subject
		timeNow := time.Now().UTC()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  jwt.NewNumericDate(timeNow),
			ExpiresAt: jwt.NewNumericDate(timeNow.Add(expiresIn)),
			Subject:   "not-a-valid-uuid",
		})
		signed, err := token.SignedString([]byte(tokenSecret))
		if err != nil {
			t.Fatalf("failed to create test token: %v", err)
		}

		_, err = ValidateJWT(signed, tokenSecret)
		if err == nil {
			t.Fatal("expected error when validating token with invalid UUID, got nil")
		}
	})

	t.Run("validates multiple tokens with same secret", func(t *testing.T) {
		userID1 := uuid.New()
		userID2 := uuid.New()
		tokenSecret := "test-secret"
		expiresIn := time.Hour

		token1, err := MakeJWT(userID1, tokenSecret, expiresIn)
		if err != nil {
			t.Fatalf("failed to create token1: %v", err)
		}

		token2, err := MakeJWT(userID2, tokenSecret, expiresIn)
		if err != nil {
			t.Fatalf("failed to create token2: %v", err)
		}

		validatedUserID1, err := ValidateJWT(token1, tokenSecret)
		if err != nil {
			t.Fatalf("failed to validate token1: %v", err)
		}

		validatedUserID2, err := ValidateJWT(token2, tokenSecret)
		if err != nil {
			t.Fatalf("failed to validate token2: %v", err)
		}

		if validatedUserID1 != userID1 {
			t.Errorf("expected user ID %s, got %s", userID1, validatedUserID1)
		}

		if validatedUserID2 != userID2 {
			t.Errorf("expected user ID %s, got %s", userID2, validatedUserID2)
		}
	})

	t.Run("returns zero UUID on validation failure", func(t *testing.T) {
		tokenSecret := "test-secret"
		invalidToken := "invalid.token.here"

		validatedUserID, err := ValidateJWT(invalidToken, tokenSecret)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		zeroUUID := uuid.UUID{}
		if validatedUserID != zeroUUID {
			t.Errorf("expected zero UUID on error, got %s", validatedUserID)
		}
	})
}
