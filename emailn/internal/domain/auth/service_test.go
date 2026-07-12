package auth

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type HTTPClientMock struct {
	mock.Mock
}

func (m *HTTPClientMock) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*http.Response), args.Error(1)
}

type TokenVerifierMock struct {
	mock.Mock
}

func (m *TokenVerifierMock) Verify(ctx context.Context, rawIDToken string) (*oidc.IDToken, error) {
	args := m.Called(ctx, rawIDToken)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*oidc.IDToken), args.Error(1)
}

var testConfig = &Config{
	BaseURL:      "http://localhost:8080",
	TokenURI:     "realms/emailn/protocol/openid-connect/token",
	ClientID:     "emailn",
	ClientSecret: "secret",
	CallbackURL:  "http://localhost:8081/auth/callback",
}

func Test_ExchangeCode_Success(t *testing.T) {
	assert := assert.New(t)
	clientMock := new(HTTPClientMock)
	service := NewService(testConfig, clientMock, nil)

	expectedToken := &TokenResponse{
		AccessToken:  "access-token-123",
		RefreshToken: "refresh-token-456",
		ExpiresIn:    300,
		TokenType:    "Bearer",
	}

	clientMock.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.Method == http.MethodPost &&
			req.Header.Get("Content-Type") == "application/x-www-form-urlencoded"
	})).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{"access_token":"access-token-123","refresh_token":"refresh-token-456","expires_in":300,"token_type":"Bearer"}`)),
	}, nil)

	token, err := service.ExchangeCode("valid-code")

	assert.Nil(err)
	assert.Equal(expectedToken, token)
	clientMock.AssertExpectations(t)
}

func Test_ExchangeCode_Error(t *testing.T) {
	assert := assert.New(t)
	clientMock := new(HTTPClientMock)
	service := NewService(testConfig, clientMock, nil)

	clientMock.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusUnauthorized,
		Body:       io.NopCloser(bytes.NewBufferString("unauthorized")),
	}, nil)

	token, err := service.ExchangeCode("bad-code")

	assert.Nil(token)
	assert.NotNil(err)
	assert.Contains(err.Error(), "token endpoint returned status 401")
	clientMock.AssertExpectations(t)
}

func Test_ExchangeCode_PropagatesConfig(t *testing.T) {
	assert := assert.New(t)
	clientMock := new(HTTPClientMock)
	service := NewService(testConfig, clientMock, nil)

	clientMock.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		body, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(body))
		return req.Method == http.MethodPost &&
			req.URL.String() == "http://localhost:8080/realms/emailn/protocol/openid-connect/token"
	})).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{"access_token":"token"}`)),
	}, nil)

	token, err := service.ExchangeCode("code-123")

	assert.Nil(err)
	assert.Equal("token", token.AccessToken)
	clientMock.AssertExpectations(t)
}

func Test_ExchangeCode_HttpError(t *testing.T) {
	assert := assert.New(t)
	clientMock := new(HTTPClientMock)
	service := NewService(testConfig, clientMock, nil)

	clientMock.On("Do", mock.Anything).Return(nil, errors.New("connection refused"))

	token, err := service.ExchangeCode("code-123")

	assert.Nil(token)
	assert.NotNil(err)
	assert.Contains(err.Error(), "failed to call token endpoint")
	clientMock.AssertExpectations(t)
}

func Test_VerifyToken_Success(t *testing.T) {
	assert := assert.New(t)
	verifierMock := new(TokenVerifierMock)
	service := NewService(testConfig, nil, verifierMock)

	expectedIDToken := &oidc.IDToken{
		Subject:  "user-123",
		Issuer:   "http://localhost:8080",
		Audience: []string{"emailn"},
		Expiry:   time.Now().Add(1 * time.Hour),
	}

	verifierMock.On("Verify", context.Background(), "valid-token").Return(expectedIDToken, nil)

	idToken, err := service.VerifyToken("valid-token")

	assert.Nil(err)
	assert.Equal(expectedIDToken, idToken)
	verifierMock.AssertExpectations(t)
}

func Test_VerifyToken_InvalidToken(t *testing.T) {
	assert := assert.New(t)
	verifierMock := new(TokenVerifierMock)
	service := NewService(testConfig, nil, verifierMock)

	verifierMock.On("Verify", mock.Anything, "invalid-token").
		Return(nil, errors.New("oidc: token format invalid"))

	_, err := service.VerifyToken("invalid-token")

	assert.NotNil(err)
	assert.Contains(err.Error(), "failed to verify token")
	verifierMock.AssertExpectations(t)
}

func Test_VerifyToken_ExpiredToken(t *testing.T) {
	assert := assert.New(t)
	verifierMock := new(TokenVerifierMock)
	service := NewService(testConfig, nil, verifierMock)

	verifierMock.On("Verify", mock.Anything, "expired-token").
		Return(nil, errors.New("oidc: token is expired"))

	_, err := service.VerifyToken("expired-token")

	assert.NotNil(err)
	assert.Contains(err.Error(), "failed to verify token")
	verifierMock.AssertExpectations(t)
}

func Test_VerifyToken_WrongAudience(t *testing.T) {
	assert := assert.New(t)
	verifierMock := new(TokenVerifierMock)
	service := NewService(testConfig, nil, verifierMock)

	verifierMock.On("Verify", mock.Anything, "wrong-aud-token").
		Return(nil, errors.New("oidc: expected audience \"emailn\" got \"other-client\""))

	_, err := service.VerifyToken("wrong-aud-token")

	assert.NotNil(err)
	assert.Contains(err.Error(), "failed to verify token")
	verifierMock.AssertExpectations(t)
}
