package sessionmanager

import (
	"context"

	"github.com/sleekservices/ServiceRenderer/common"
	"github.com/sleekservices/ServiceRenderer/common/authority"
)

type Store interface {
	// GetSession retrieves a saved session from storage. If the
	// session was not found or expired, nil should be returned with
	// an appropriate error.
	GetSession(context.Context, string) (*common.Session, error)

	// UpdateSession updates a session by saving it to storage. If a
	// session already exists in storage it should be replaced. If the
	// session is authenticated, it should be stored under that user.
	UpdateSession(context.Context, *common.Session) error

	// DeleteSession deletes a session from storage.
	DeleteSession(context.Context, *common.Session) error

	// DeleteSessionsForAuthority deletes all sessions attached to a authority.
	DeleteSessionsForAuthority(context.Context, *authority.ServiceRenderer) error
}
