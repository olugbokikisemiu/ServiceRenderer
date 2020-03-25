package sessionmanager

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"github.com/sleekservices/ServiceRenderer/common"
	"github.com/sleekservices/ServiceRenderer/common/errors"
)

type SuccessHandler func(c *gin.Context, s *common.Session)
type FailHandler func(c *gin.Context, err error)

func (m *Manager) FindOrCreateSession(c *gin.Context) {
	s, err := m.GetOrCreateSession(c)
	if err != nil {
		log.Debug("issue finding/creating session err=%v", err)
		m.handleFailure(c, err)
		return
	}

	err = m.applySessionCookie(c, s)
	if err != nil {
		log.Debug("unable to apply session err=%v", err)
		m.handleFailure(c, err)
		return
	}
	c.Set(ContextKey, s)
}

func (m *Manager) Authorized(success SuccessHandler) SuccessHandler {
	return func(c *gin.Context, s *common.Session) {
		if !s.IsAuthenticated() {
			m.handleFailure(c, errors.ErrorLog(errors.ErrUnauthenticated))
			return
		}
		c.Set(ContextKey, s)
		success(c, s)
	}
}

func (m *Manager) handleFailure(c *gin.Context, err error) {
	switch err.(type) {
	case *errors.Error:
		c.AbortWithStatusJSON(http.StatusForbidden, err)
	default:
		c.AbortWithStatusJSON(http.StatusForbidden, errors.ErrorLog(errors.ErrUnknown))
	}
}
