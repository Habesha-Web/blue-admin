package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// // Mock configuration for testing
// func init() {
// 	// Set a fixed salt for consistent testing
// 	configs.AppConfig = &mockConfig{
// 		salt: "test_salt",
// 	}
// }

// Mock configuration struct to override AppConfig.Get method
type mockConfig struct {
	salt string
}

func (m *mockConfig) Get(key string) string {
	if key == "SECRETE_SALT" {
		return m.salt
	}
	return ""
}

func TestHashFunc(t *testing.T) {
	password := "password123"
	expectedHash := "d81a4bcb5018ed48fd54d83f7f3d164a1fd2e2c4039fc99b297e7fbe23af373d7309cebdc3a928c9dd01ec51f70d50d1a4ec98ef398fd13034b58760c60e79f6"

	hashedPassword := HashFunc(password)
	assert.Equal(t, expectedHash, hashedPassword, "Hashed password should match the expected hash")
}

func TestPasswordsMatch(t *testing.T) {
	password := "password123"
	hashedPassword := HashFunc(password)

	match := PasswordsMatch(hashedPassword, password)
	assert.True(t, match, "Passwords should match")

	noMatch := PasswordsMatch(hashedPassword, "wrongpassword")
	assert.False(t, noMatch, "Passwords should not match")
}

func TestCreateJWTToken(t *testing.T) {
	email := "test@example.com"
	uuid := "test-uuid"
	roles := []string{"user"}
	duration := 10
	user_id := 4

	token, err := CreateJWTToken(email, uuid, user_id, roles, duration)
	require.NoError(t, err, "Token creation should not return an error")

	// Check if the token can be parsed
	parsedClaim, err := ParseJWTToken(token)
	require.NoError(t, err, "Token parsing should not return an error")

	assert.Equal(t, email, parsedClaim.Email, "Email in token should match")
	assert.Equal(t, uuid, parsedClaim.UUID, "UUID in token should match")
	assert.ElementsMatch(t, roles, parsedClaim.Roles, "Roles in token should match")
}

func TestParseJWTToken(t *testing.T) {
	email := "test@example.com"
	uuid := "test-uuid"
	roles := []string{"user"}
	duration := 10
	user_id := 4

	token, err := CreateJWTToken(email, uuid, user_id, roles, duration)
	require.NoError(t, err, "Token creation should not return an error")

	parsedClaim, err := ParseJWTToken(token)
	require.NoError(t, err, "Token parsing should not return an error")

	assert.Equal(t, email, parsedClaim.Email, "Email in token should match")
	assert.Equal(t, uuid, parsedClaim.UUID, "UUID in token should match")
	assert.ElementsMatch(t, roles, parsedClaim.Roles, "Roles in token should match")

	// Test invalid token
	invalidToken := "invalid.token.here"
	_, err = ParseJWTToken(invalidToken)
	assert.Error(t, err, "Parsing an invalid token should return an error")
}

func TestUniqueSlice(t *testing.T) {
	input := []string{"a", "b", "a", "c", "b"}
	expected := []string{"a", "b", "c"}

	result := UniqueSlice(input)
	assert.ElementsMatch(t, expected, result, "Unique slice should match the expected unique values")
}

func TestCheckValueExistsInSlice(t *testing.T) {
	slice_super := []string{"admin", "user", "guest", "superuser"}
	slice := []string{"admin", "user", "guest"}

	assert.True(t, CheckValueExistsInSlice(slice, "admin"), "Slice should contain 'admin'")
	assert.True(t, CheckValueExistsInSlice(slice_super, "superuser"), "Slice should contain 'superuser'")
	assert.False(t, CheckValueExistsInSlice(slice, "unknown"), "Slice should not contain 'unknown'")
}
