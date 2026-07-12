package auth_test

import (
	"emailn/cmd/api/container"
	"emailn/cmd/api/routes"
	"emailn/internal/domain/auth"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) ExchangeCode(code string) (*auth.TokenResponse, error) {
	args := s.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.TokenResponse), args.Error(1)
}
func (s *ServiceMock) VerifyToken(token string) (*oidc.IDToken, error) {
	return nil, nil
}

type testCase struct {
	name       string
	method     string
	url        string
	setupMock  func(*ServiceMock)
	wantStatus int
	wantBody   func(*testing.T, *httptest.ResponseRecorder, *ServiceMock)
}

func TestHandler(t *testing.T) {
	for _, tc := range []testCase{
		{
			name:       "login redirect to keycloak",
			method:     http.MethodGet,
			url:        "/auth",
			setupMock:  func(m *ServiceMock) {},
			wantStatus: http.StatusTemporaryRedirect,
			wantBody: func(t *testing.T, w *httptest.ResponseRecorder, m *ServiceMock) {
				location := w.Header().Get("Location")
				assert.Contains(t, location, "client_id=test-client")
				assert.Contains(t, location, "response_type=code")
				assert.Contains(t, location, "scope=openid")
				assert.Contains(t, location, "state=")
				m.AssertNotCalled(t, "ExchangeCode")
			},
		},
		{
			name:   "callback success",
			method: http.MethodGet,
			url:    "/auth/callback?code=valid-code&state=valid-state",
			setupMock: func(m *ServiceMock) {
				m.On("ExchangeCode", "valid-code").Return(&auth.TokenResponse{
					AccessToken:  "access-123",
					RefreshToken: "refresh-456",
					ExpiresIn:    300,
					TokenType:    "Bearer",
				}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody: func(t *testing.T, w *httptest.ResponseRecorder, m *ServiceMock) {
				var result auth.TokenResponse
				json.Unmarshal(w.Body.Bytes(), &result)
				assert.Equal(t, "access-123", result.AccessToken)
				assert.Equal(t, "refresh-456", result.RefreshToken)
				assert.Equal(t, 300, result.ExpiresIn)
				assert.Equal(t, "Bearer", result.TokenType)
				m.AssertExpectations(t)
			},
		},
		{
			name:       "callback missing code and state",
			method:     http.MethodGet,
			url:        "/auth/callback",
			setupMock:  func(m *ServiceMock) {},
			wantStatus: http.StatusUnauthorized,
			wantBody: func(t *testing.T, w *httptest.ResponseRecorder, m *ServiceMock) {
				var result map[string]string
				json.Unmarshal(w.Body.Bytes(), &result)
				assert.Equal(t, "Incorrect auth flow credentials", result["error"])
				m.AssertNotCalled(t, "ExchangeCode")
			},
		},
		{
			name:       "callback missing only code",
			method:     http.MethodGet,
			url:        "/auth/callback?state=valid-state",
			setupMock:  func(m *ServiceMock) {},
			wantStatus: http.StatusUnauthorized,
			wantBody: func(t *testing.T, w *httptest.ResponseRecorder, m *ServiceMock) {
				var result map[string]string
				json.Unmarshal(w.Body.Bytes(), &result)
				assert.Equal(t, "Incorrect auth flow credentials", result["error"])
				m.AssertNotCalled(t, "ExchangeCode")
			},
		},
		{
			name:   "callback service error",
			method: http.MethodGet,
			url:    "/auth/callback?code=bad-code&state=valid-state",
			setupMock: func(m *ServiceMock) {
				m.On("ExchangeCode", "bad-code").Return(nil, errors.New("token endpoint returned status 401"))
			},
			wantStatus: http.StatusBadRequest,
			wantBody: func(t *testing.T, w *httptest.ResponseRecorder, m *ServiceMock) {
				var result map[string]string
				json.Unmarshal(w.Body.Bytes(), &result)
				assert.Equal(t, "token endpoint returned status 401", result["error"])
				m.AssertExpectations(t)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			serviceMock := new(ServiceMock)
			tc.setupMock(serviceMock)

			config := &auth.Config{
				BaseURL:     "http://localhost:8080",
				AuthURI:     "realms/emailn/protocol/openid-connect/auth",
				CallbackURL: "http://localhost:8081/auth/callback",
				ClientID:    "test-client",
			}

			handler := auth.NewHandler(config, serviceMock)
			ctn := &container.Container{AuthHandler: handler}

			router := gin.New()
			routes.RegisterRoutes(router, ctn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tc.method, tc.url, nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
			tc.wantBody(t, w, serviceMock)
		})
	}
}
