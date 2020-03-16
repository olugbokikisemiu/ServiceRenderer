package password

import (
	"crypto/rand"
	"crypto/subtle"
	"regexp"
	"strings"

	"github.com/sleekservices/ServiceRenderer/common/must"
	"golang.org/x/crypto/scrypt"
)

const (
	PasswordLen = 8
	HashLen     = 64
	SaltLen     = 32
)

var (
	bigLetters   = regexp.MustCompile(`^.*[A-Z].*$`)
	smallLetters = regexp.MustCompile(`^.*[a-z].*$`)
	numbers      = regexp.MustCompile(`^.*\d.*$`)
)

type Hash struct {
	Hash []byte `bson:"hash"`
	Salt []byte `bosn:"salt"`
}

func ParsePassword(password string) bool {
	return len(password) >= PasswordLen &&
		bigLetters.MatchString(password) &&
		smallLetters.MatchString(password) &&
		numbers.MatchString(password)
}

func generateSalt() []byte {
	salt := make([]byte, SaltLen)
	must.DoF(func() error {
		_, err := rand.Read(salt)
		return err
	})
	return salt
}

func createPasswordHash(password string, salt []byte) ([]byte, error) {
	password = strings.TrimSpace(password)
	hash, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, HashLen)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func HashPassword(password string) (*Hash, error) {
	salt := generateSalt()
	hash, err := createPasswordHash(password, salt)
	if err != nil {
		return nil, err
	}
	return &Hash{Hash: hash, Salt: salt}, nil
}

func verifyPassword(password string, hash []byte, salt []byte) bool {
	verifyPass, err := createPasswordHash(password, salt)
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare(verifyPass, hash) == 1
}

func (h *Hash) IsEqualTo(password string) bool {
	return verifyPassword(password, h.Hash, h.Salt)
}
