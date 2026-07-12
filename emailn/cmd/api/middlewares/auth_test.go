package middlewares_test

import (
	"emailn/cmd/api/middlewares"
	"emailn/internal/domain/auth"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
	args := s.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*oidc.IDToken), args.Error(1)
}

type testCase struct {
	name       string
	header     string
	setupMock  func(*ServiceMock)
	wantStatus int
	wantBody   func(*testing.T, *httptest.ResponseRecorder)
}

func TestAuth(t *testing.T) {
	for _, tc := range []testCase{
		{
			name:       "missing Authorization header",
			header:     "",
			setupMock:  func(m *ServiceMock) {},
			wantStatus: http.StatusUnauthorized,
			wantBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var body map[string]string
				json.Unmarshal(w.Body.Bytes(), &body)
				assert.Equal(t, "Unauthorized", body["error"])
			},
		},
		{
			name:       "invalid Authorization format - no Bearer prefix",
			header:     "Token abc123",
			setupMock:  func(m *ServiceMock) {},
			wantStatus: http.StatusUnauthorized,
			wantBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var body map[string]string
				json.Unmarshal(w.Body.Bytes(), &body)
				assert.Equal(t, "Unauthorized", body["error"])
			},
		},
		{
			name:       "invalid Authorization format - only Bearer",
			header:     "Bearer",
			setupMock:  func(m *ServiceMock) {},
			wantStatus: http.StatusUnauthorized,
			wantBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var body map[string]string
				json.Unmarshal(w.Body.Bytes(), &body)
				assert.Equal(t, "Unauthorized", body["error"])
			},
		},
		{
			name:   "invalid token",
			header: "Bearer invalid-token",
			setupMock: func(m *ServiceMock) {
				m.On("VerifyToken", "invalid-token").
					Return(nil, errors.New("oidc: token format invalid"))
			},
			wantStatus: http.StatusUnauthorized,
			wantBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var body map[string]string
				json.Unmarshal(w.Body.Bytes(), &body)
				assert.Equal(t, "Unauthorized", body["error"])
			},
		},
		{
			name:   "valid token",
			header: "Bearer valid-token",
			setupMock: func(m *ServiceMock) {
				idToken := &oidc.IDToken{
					Subject:  "user-123",
					Issuer:   "http://localhost:8080",
					Audience: []string{"emailn"},
					Expiry:   time.Now().Add(1 * time.Hour),
				}
				m.On("VerifyToken", "valid-token").Return(idToken, nil)
			},
			wantStatus: http.StatusOK,
			wantBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name:   "expired token",
			header: "Bearer expired-token",
			setupMock: func(m *ServiceMock) {
				m.On("VerifyToken", "expired-token").
					Return(nil, errors.New("oidc: token is expired"))
			},
			wantStatus: http.StatusUnauthorized,
			wantBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var body map[string]string
				json.Unmarshal(w.Body.Bytes(), &body)
				assert.Equal(t, "Unauthorized", body["error"])
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			serviceMock := new(ServiceMock)
			tc.setupMock(serviceMock)

			w := httptest.NewRecorder()
			router := gin.New()

			router.GET("/test", middlewares.Auth(serviceMock), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "ok"})
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
			tc.wantBody(t, w)
			serviceMock.AssertExpectations(t)
		})
	}
}
