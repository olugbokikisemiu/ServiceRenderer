package password

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidatePassword__Should_validate_password(t *testing.T) {
	assert.True(t, ParsePassword("Password123"))
	assert.True(t, ParsePassword("Pa55word123"))
	assert.True(t, ParsePassword("pAssWord1"))

	assert.False(t, ParsePassword("pass"))
	assert.False(t, ParsePassword("password"))
	assert.False(t, ParsePassword("PASSWORD123"))
	assert.False(t, ParsePassword("password123"))

}

func TestPasswordHash__Should_hash_password_successfully(t *testing.T) {
	password := []string{
		"Pa55word123",
		"Password123",
	}
	hash, err := HashPassword(password[0])
	assert.Nil(t, err)
	assert.True(t, hash.IsEqualTo(password[0]))
	assert.False(t, hash.IsEqualTo(password[1]))
}
