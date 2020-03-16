package accounts

import (
	"context"
	"testing"

	"github.com/go-redis/redis"
	"github.com/golang/mock/gomock"
	"github.com/sleekservices/ServiceRenderer/common"
	"github.com/sleekservices/ServiceRenderer/common/authority"
	"github.com/sleekservices/ServiceRenderer/common/badactor"
	"github.com/sleekservices/ServiceRenderer/common/errors"
	"github.com/sleekservices/ServiceRenderer/common/password"
	r "github.com/sleekservices/ServiceRenderer/common/redis"
	"github.com/sleekservices/ServiceRenderer/mocks"
	"github.com/stretchr/testify/assert"
	"time"
)

var (
	mockProviderService *mocks.MockServiceProvider
	mockBadactorService *mocks.MockBadactorService
	ctx                 context.Context
)

func sessionStore() *r.SessionStore {
	c := redis.NewClient(&redis.Options{
		Addr:        "localhost:6379",
		Password:    "",
		DialTimeout: 5 * time.Minute,
		MaxRetries:  10,
	})
	redisClient := r.Client{
		Client: c,
	}
	return r.NewSessionStore(&redisClient)
}

func TestRegisterServiceProvider__Should_register_provider_successfully(t *testing.T) {
	ctx := context.Background()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockProviderService := mocks.NewMockServiceProvider(mockCtrl)

	mockProviderService.
		EXPECT().
		CreateServiceProvider(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	h := NewHandler(mockProviderService, nil, nil)

	r := &authority.ServiceRenderer{
		FirstName: "Mark",
		LastName:  "Hill",
	}

	err := h.RegisterServiceProvider(ctx, r)
	assert.NoError(t, err)
	assert.Equal(t, "Mark", r.FirstName)
}

func TestServiceProviderLogin__Should_login_user_with_valid_email_and_password(t *testing.T) {
	ctx := context.Background()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockProviderService := mocks.NewMockServiceProvider(mockCtrl)

	userPin := "123456"

	hashPin, err := password.HashPin(userPin)

	assert.NoError(t, err)

	providerDetails := &authority.ServiceRenderer{
		FirstName: "Mark",
		LastName:  "Hill",
		Pin:       hashPin,
		Email:     "test@rendize.com",
	}

	mockProviderService.
		EXPECT().
		FindServiceProviderByEmail(gomock.Any(), gomock.Any()).
		Return(providerDetails, nil).AnyTimes()

	mockBadactorService := mocks.NewMockBadactorService(mockCtrl)
	mockBadactorService.
		EXPECT().
		FindJail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errors.ErrorLog(errors.ErrNoDocumentFound)).
		AnyTimes()

	badactorStudio := badactor.NewStudio(ctx, true, mockBadactorService)

	h := NewHandler(mockProviderService, badactorStudio, sessionStore())

	session := new(common.Session)
	session.ServiceRenderer = providerDetails

	loginUser, err := h.ServiceProviderLogin(ctx, session, providerDetails.Email, userPin)

	assert.NoError(t, err)
	assert.NotEmpty(t, loginUser)
	assert.Equal(t, providerDetails.FirstName, loginUser.FirstName)
}

func TestServiceProviderLogin__Should_err_with_invalid_email_and_password(t *testing.T) {
	ctx := context.Background()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockProviderService := mocks.NewMockServiceProvider(mockCtrl)

	userPin := "123456"

	hashPin, err := password.HashPin(userPin)

	assert.NoError(t, err)

	providerDetails := &authority.ServiceRenderer{
		FirstName: "Mark",
		LastName:  "Hill",
		Pin:       hashPin,
		Email:     "test@rendize.com",
	}

	session := new(common.Session)
	session.ServiceRenderer = providerDetails

	mockProviderService.
		EXPECT().
		FindServiceProviderByEmail(gomock.Any(), gomock.Any()).
		Return(providerDetails, nil).AnyTimes()

	mockBadactorService := mocks.NewMockBadactorService(mockCtrl)

	mockBadactorService.
		EXPECT().
		FindJail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errors.ErrorLog(errors.ErrNoDocumentFound)).
		AnyTimes()

	mockBadactorService.
		EXPECT().
		CreateInfraction(gomock.Any(), gomock.Any()).
		Return(nil, nil).AnyTimes()

	mockBadactorService.
		EXPECT().
		CountInfraction(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(int64(1), nil)

	badactorStudio := badactor.NewStudio(ctx, true, mockBadactorService)

	h := NewHandler(mockProviderService, badactorStudio, sessionStore())

	_, err = h.ServiceProviderLogin(ctx, session, providerDetails.Email, "1234")
	assert.Error(t, err)
	assert.EqualError(t, err, errors.ErrorLog(errors.ErrInvalidCredentials).Error())
}
