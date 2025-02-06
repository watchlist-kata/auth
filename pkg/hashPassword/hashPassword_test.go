package hashPassword

import "testing"

func TestHashPassword(t *testing.T) {
	password := "testpassword"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Errorf("Error hashing password: %v", err)
	}

	if len(hashedPassword) == 0 {
		t.Errorf("Expected hashed password to be non-empty")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testpassword"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Errorf("Error hashing password: %v", err)
	}

	if !CheckPasswordHash(password, hashedPassword) {
		t.Errorf("Expected password hash to match")
	}

	if CheckPasswordHash("wrongpassword", hashedPassword) {
		t.Errorf("Expected password hash not to match for wrong password")
	}
}
