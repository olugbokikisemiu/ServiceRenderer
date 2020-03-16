package password

import "regexp"

import "github.com/sleekservices/ServiceRenderer/common/errors"

const (
	PinLen = 6
)

var pinRegex = regexp.MustCompile(`^\d+$`)

func ParsePin(p string) bool {
	return len(p)>= PinLen && pinRegex.MatchString(p)
}

func HashPin(pin string) (*Hash, error) {
	if ParsePin(pin) {
		return HashPassword(pin)
	}
	return nil, errors.ErrorLog(errors.ErrInvalidPinFormat)
}
