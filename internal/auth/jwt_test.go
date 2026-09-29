package auth

import (
	"testing"
	"time"

	"github.com/miriyaladhanwinn/on-duty/internal/db"
)

func TestTokenGenerationAndValidation(t *testing.T) {
	user := &db.User{
		ID:         "test-user-uuid",
		Email:      "student@cb.students.amrita.edu",
		Name:       "Test Student",
		RollNo:     "CB.SC.U4CSE25001",
		Department: "CSE",
		Role:       db.RoleStudent,
	}

	token, err := GenerateToken(user, 1*time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.UserID != user.ID || claims.Email != user.Email || claims.Role != user.Role {
		t.Fatalf("Claims mismatch: %+v vs %+v", claims, user)
	}
}

func TestPasswordHashingAndVerification(t *testing.T) {
	rawPassword := "SecureAmritaPass@2026"
	hash, err := HashPassword(rawPassword)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if !CheckPassword(rawPassword, hash) {
		t.Fatal("CheckPassword failed with correct password")
	}

	if CheckPassword("WrongPassword123", hash) {
		t.Fatal("CheckPassword succeeded with incorrect password")
	}
}
