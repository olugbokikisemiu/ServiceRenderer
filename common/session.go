package common

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sleekservices/ServiceRenderer/common/authority"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Claim
type Claims struct {
	SessionID string `bson:"session_id"`
	jwt.StandardClaims
}

// Session contains all claims and extra data which form a stateful session.
type Session struct {
	Claims          `json:"claims"`
	ServiceRenderer *authority.ServiceRenderer `json:"serviceRender"`
	Context         *gin.Context               `json:"-"`
}

// Store is responsible for storage of sessions.
type Store interface {
	// GetSession retrieves a saved session from storage. If the
	// session was not found or expired, nil should be returned with
	// an appropriate error.
	GetSession(context.Context, string) (*Session, error)

	// UpdateSession updates a session by saving it to storage. If a
	// session already exists in storage it should be replaced. If the
	// session is authenticated, it should be stored under that user.
	UpdateSession(context.Context, *Session) error

	// DeleteSession deletes a session from storage.
	DeleteSession(context.Context, *Session) error
}

func (s *Session) IsAuthenticated() bool {
	return !s.ServiceRenderer.ID.IsZero()
}

func (s *Session) UserID() primitive.ObjectID {
	return s.ServiceRenderer.ID
}

func (s *Session) Expired() bool {
	return s.ExpiresAt < time.Now().Unix()
}

func (s *Session) AuthorityID() string {
	return s.ServiceRenderer.ID.Hex()
}

func (s *Session) UserAgent() string {
	return s.Context.Request.UserAgent()
}

func (s *Session) IsAuthenticatedAuthority(authorityID primitive.ObjectID) bool {
	return s.IsAuthenticated() && s.ServiceRenderer.ID == authorityID
}

func (s *Session) SetClaims(expiration time.Time) {
	s.Claims.ExpiresAt = expiration.Unix()
}
