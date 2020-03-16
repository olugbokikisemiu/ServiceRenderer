package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sleekservices/ServiceRenderer/common"
	"github.com/sleekservices/ServiceRenderer/common/authority"
	"github.com/sleekservices/ServiceRenderer/common/errors"
)

const (
	redisNamespace = "session"
)

type SessionStore struct {
	c *Client
}

func NewSessionStore(c *Client) *SessionStore {
	return &SessionStore{
		c: c,
	}
}

func (ss *SessionStore) cleanupSessionsForAuthority(uID string) error {
	set := fmt.Sprintf("%s-%s", redisNamespace, uID)
	mems, err := ss.c.SMembers(set)
	if err != nil {
		return err
	}

	// Clean up expired session
	for _, m := range mems {
		exists, _ := ss.c.Exists(m)
		if !exists {
			ss.c.SRemove(set, m)
		}
	}
	return nil
}

func (ss *SessionStore) updateSessionsForAuthority(uID, sID string) error {
	if err := ss.cleanupSessionsForAuthority(uID); err != nil {
		return err
	}
	set := fmt.Sprintf("%s-%s", redisNamespace, uID)
	return ss.c.SAdd(set, sID)
}

func (ss *SessionStore) GetSession(c context.Context, sID string) (*common.Session, error) {
	var s *common.Session

	key := fmt.Sprintf("%s-%s", redisNamespace, sID)
	exists, _ := ss.c.Exists(key)
	if !exists {
		return nil, errors.ErrorLog(errors.ErrNotFound)
	}

	val, err := ss.c.Get(key)
	if err != nil {
		return nil, err
	}

	return s, json.Unmarshal([]byte(val.(string)), &s)
}

func (ss *SessionStore) UpdateSession(c context.Context, s *common.Session) error {
	val, err := json.Marshal(s)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s-%s", redisNamespace, s.SessionID)
	if err := ss.c.Set(key, val); err != nil {
		return err
	}

	if s.IsAuthenticated() {
		if err := ss.updateSessionsForAuthority(s.ServiceRenderer.ID.Hex(), s.SessionID); err != nil {
			return err
		}
	}

	return nil
}

func (ss *SessionStore) DeleteSession(c context.Context, s *common.Session) error {
	key := fmt.Sprintf("%s-%s", redisNamespace, s.SessionID)
	if _, err := ss.c.Delete(key); err != nil {
		return err
	}

	if s.IsAuthenticated() {
		if err := ss.cleanupSessionsForAuthority(s.ServiceRenderer.ID.Hex()); err != nil {
			return err
		}
	}
	return nil
}

func (ss *SessionStore) DeleteSessionsForAuthority(c context.Context, a *authority.ServiceRenderer) error {
	set := fmt.Sprintf("%s-%s", redisNamespace, a.ID.String())
	mems, err := ss.c.SMembers(set)
	if err != nil {
		return err
	}
	// Cleanup all sessions
	for _, m := range mems {
		ss.c.SRemove(set, m)
		ss.c.Delete(m)
	}
	return nil
}
