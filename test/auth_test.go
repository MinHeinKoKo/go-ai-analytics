package test

import (
	"ai-analytics/internal/utils"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestPasswordHashing(t *testing.T) {
	password := "testpassword123"

	// Test password hashing
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Test password verification
	err = utils.CheckPassword(hashedPassword, password)
	if err != nil {
		t.Fatalf("Password verification failed: %v", err)
	}

	// Test wrong password
	err = utils.CheckPassword(hashedPassword, "wrongpassword")
	if err == nil {
		t.Fatal("Expected password verification to fail with wrong password")
	}
}

func TestJWTToken(t *testing.T) {
	userID := primitive.NewObjectID()
	email := "test@example.com"
	secret := "test-secret"
	expiryHours := 24

	// Generate token
	token, err := utils.GenerateToken(userID, email, secret, expiryHours)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Validate token
	claims, err := utils.ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	// Check claims
	if claims.UserID != userID {
		t.Fatalf("Expected user ID %v, got %v", userID, claims.UserID)
	}

	if claims.Email != email {
		t.Fatalf("Expected email %s, got %s", email, claims.Email)
	}
}
