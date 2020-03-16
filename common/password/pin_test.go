package password
import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidatePin__Should_validate_pin(t *testing.T) {
	assert.True(t, ParsePin("123456"))

	assert.False(t, ParsePin("1234"))
}

func TestHashPin__Should_hash_pin_successful(t *testing.T) {
	pin := []string{
		"123456",
		"987654",
	}

	hash, err := HashPin(pin[0])
	assert.Nil(t, err)

	assert.True(t, hash.IsEqualTo(pin[0]))
	assert.False(t, hash.IsEqualTo(pin[1]))
}
