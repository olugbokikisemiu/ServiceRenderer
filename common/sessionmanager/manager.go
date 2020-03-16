package sessionmanager

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	geoip "github.com/oschwald/geoip2-golang"
	"github.com/sleekservices/ServiceRenderer/common"
	"github.com/sleekservices/ServiceRenderer/common/authority"
	"github.com/sleekservices/ServiceRenderer/common/errors"
	"github.com/sleekservices/ServiceRenderer/common/generator"
)

// Cookie string constants
const (
	AuthHeaderName   = "Authorization"
	authHeaderBearer = "bearer"
	authHeaderToken  = "token"

	DefaultCookieName = "sessionid"
	ContextKey        = "session"
)

// ErrNoSession is returned when session cookie is not found in the request header.
var ErrNoSession = fmt.Errorf("session: not found")
var ErrNoAuthHeader = fmt.Errorf("auth header: invalid or not found")

// Manager manages all sessions.
type Manager struct {
	secret       string
	secureCookie bool
	geoip        *geoip.Reader
	store        Store
}

func NewManager(
	secret string,
	secureCookie bool,
	store Store,
) *Manager {
	db, err := geoip.Open("./_static/GeoLite2-Country.mmdb")
	if err != nil {
		log.Panic("Unable to open geoip database %v", err)
	}
	return &Manager{
		secret:       secret,
		secureCookie: secureCookie,
		geoip:        db,
		store:        store,
	}
}

func getCookieExpirationDate(ttl time.Duration) int64 {
	return time.Now().Add(ttl).Unix()
}

func (m *Manager) applyCookie(c *gin.Context, name string, value string, expires time.Time) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   m.secureCookie,
		SameSite: http.SameSiteLaxMode,
	}

	// Add a Set-Cookie header to the response
	http.SetCookie(c.Writer, &cookie)

	// Also append the cookie to the currently processed request, so that subsequent tries to access the cookie will work.
	// For this to work, we need to get the last session cookie in request.
	c.Request.AddCookie(&cookie)
}

func (m *Manager) EncodeSession(session *common.Session) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, session.Claims)
	tokenValue, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", err
	}

	return tokenValue, nil
}

func (m *Manager) applySessionCookie(c *gin.Context, session *common.Session) error {
	signedToken, err := m.EncodeSession(session)
	if err != nil {
		return err
	}

	m.applyCookie(c, DefaultCookieName, signedToken, time.Unix(session.ExpiresAt, 0))

	return nil
}

// latestCookie returns the last occurrence of the named cookie in the cookie header,
// or ErrNoCookie if such a cookie was not found.
// Method was introduced because Gin provides API only for getting the first occurrence.
func (m *Manager) latestCookie(c *gin.Context, name string) (*http.Cookie, error) {
	cookies := c.Request.Cookies()
	for i := len(cookies) - 1; i >= 0; i-- {
		if cookies[i].Name == name {
			return cookies[i], nil
		}
	}
	return nil, http.ErrNoCookie
}

// Extracts Authorization header value
// and returns the Auth header value
func (m *Manager) extractAuthHeader(c *gin.Context) (string, error) {
	header := c.Request.Header.Get(AuthHeaderName)
	if len(header) == 0 {
		return "", ErrNoAuthHeader
	}
	tokens := strings.Fields(header)
	if len(tokens) < 2 {
		return "", ErrNoAuthHeader
	}
	// Due to spec ambiguity and http clients
	// implementation,
	// search token-type in case-insensitive
	switch strings.ToLower(tokens[0]) {
	case authHeaderBearer, authHeaderToken:
		return tokens[1], nil
	default:
		return "", ErrNoAuthHeader
	}
}

func (m *Manager) getEncodedToken(c *gin.Context) string {
	authValue, err := m.extractAuthHeader(c)
	if err == nil {
		log.Debug("got session from auth header: %s", authValue)
		return authValue
	}
	cookie, err := m.latestCookie(c, DefaultCookieName)
	if err == nil {
		log.Debug("got session from cookie: %s", cookie.Value)
		return cookie.Value
	}

	return ""
}

func (m *Manager) UpdateSession(c *gin.Context, s *common.Session) error {
	return m.store.UpdateSession(c, s)
}

func NewAnonymousUserSession(sID string, ttl time.Duration) *common.Session {
	return &common.Session{
		ServiceRenderer: &authority.Anonymous,
		Claims: common.Claims{
			sID,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(ttl).Unix(),
			},
		},
	}
}

func (m *Manager) newAnonSession(c *gin.Context) (*common.Session, error) {
	log.Debug("creating new anon session")

	sID := generator.GenerateRandomToken()

	s := NewAnonymousUserSession(sID, time.Duration(86400)*time.Minute)
	if err := m.store.UpdateSession(c, s); err != nil {
		return nil, err
	}
	return s, nil
}

func (m *Manager) GetSession(c *gin.Context, t string) (*common.Session, error) {
	token, err := jwt.ParseWithClaims(
		t, &common.Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(m.secret), nil
		})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*common.Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid session claims")
	}

	if claims.Valid() != nil {
		return nil, fmt.Errorf("invalid session claims")
	}

	return m.store.GetSession(c, claims.SessionID)
}

func (m *Manager) GetOrCreateSession(c *gin.Context) (*common.Session, error) {
	t := m.getEncodedToken(c)

	if t != "" {
		s, err := m.GetSession(c, t)
		if err == nil {
			return s, err
		}
	}

	return m.newAnonSession(c)
}

// ForContext finds the session from the context
func ForContext(ctx context.Context) (*common.Session, error) {
	raw, ok := ctx.Value(ContextKey).(*common.Session)
	if !ok {
		return nil, errors.ErrorLog(errors.ErrUnableToGetSession)
	}
	return raw, nil
}
