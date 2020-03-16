package accounts

import (
	"context"
	"fmt"

	"github.com/sleekservices/ServiceRenderer/common"
	"github.com/sleekservices/ServiceRenderer/common/authority"
	"github.com/sleekservices/ServiceRenderer/common/badactor"
	"github.com/sleekservices/ServiceRenderer/common/errors"
)

const (
	loginRule = "Login"
	waitTine  = "5"
)

type Handler struct {
	serviceProvider authority.ServiceProvider
	studio          *badactor.Studio
	store           common.Store
}

func NewHandler(sp authority.ServiceProvider, studio *badactor.Studio, store common.Store) *Handler {
	return &Handler{serviceProvider: sp, studio: studio}
}

func (h *Handler) RegisterServiceProvider(ctx context.Context, r *authority.ServiceRenderer) error {
	return h.serviceProvider.CreateServiceProvider(ctx, r)
}

func (h *Handler) ServiceProviderLogin(ctx context.Context, s *common.Session, email, pin string) (*authority.ServiceRenderer, error) {
	if s.IsAuthenticated() {
		return nil, errors.ErrorLog(errors.ErrAuthenticated)
	}

	if h.studio.IsJailedFor(ctx, email, loginRule) {
		return nil, errors.ErrorLog(errors.ErrJailedLogin, fmt.Sprintf("kindly try again after %s mins", waitTine))
	}

	serviceProvider, err := h.serviceProvider.FindServiceProviderByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !serviceProvider.Pin.IsEqualTo(pin) {
		if err := h.studio.Infraction(ctx, email, loginRule); err != nil {
			return nil, errors.ErrorLog(errors.ErrBadActor, err.Error())
		}

		if h.studio.IsJailedFor(ctx, email, loginRule) {
			return nil, errors.ErrorLog(errors.ErrJailedLogin, fmt.Sprintf("kindly try again after %s mins", waitTine))
		}

		return nil, errors.ErrorLog(errors.ErrInvalidCredentials)
	}
	h.setAuthority(ctx, s, serviceProvider)
	return serviceProvider, nil
}

func (h *Handler) setAuthority(ctx context.Context, s *common.Session, provider *authority.ServiceRenderer) {
	s.ServiceRenderer = provider
	h.store.UpdateSession(ctx, s)
}
