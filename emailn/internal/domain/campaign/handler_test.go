package campaign_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"emailn/cmd/api/container"
	"emailn/cmd/api/routes"
	"emailn/internal/domain/auth"
	"emailn/internal/domain/campaign"
	"emailn/internal/types"

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

type AuthServiceMock struct {
	mock.Mock
}

func (s *AuthServiceMock) ExchangeCode(code string) (*auth.TokenResponse, error) {
	args := s.Called(code)
	return nil, args.Error(1)
}

func (s *AuthServiceMock) VerifyToken(token string) (*oidc.IDToken, error) {
	args := s.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*oidc.IDToken), args.Error(1)
}

func (s *ServiceMock) Create(newCampaign *campaign.CreateCampaignRequest) (*campaign.CampaignResponse, error) {
	args := s.Called(newCampaign)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*campaign.CampaignResponse), args.Error(1)
}

func (s *ServiceMock) Get() (*[]campaign.CampaignResponse, error) {
	args := s.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]campaign.CampaignResponse), args.Error(1)
}

func (s *ServiceMock) Show(id types.UUID) (*campaign.CampaignResponse, error) {
	args := s.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*campaign.CampaignResponse), args.Error(1)
}

func (s *ServiceMock) Cancel(id types.UUID) (*campaign.CampaignResponse, error) {
	args := s.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*campaign.CampaignResponse), args.Error(1)
}

func (s *ServiceMock) Delete(id types.UUID) error {
	args := s.Called(id)
	return args.Error(0)
}

type testCase struct {
	name       string
	method     string
	url        string
	body       string
	setupMock  func(*ServiceMock)
	wantStatus int
	wantBody   func(*testing.T, *httptest.ResponseRecorder, *ServiceMock)
}

func TestHandler(t *testing.T) {
	for _, tc := range []testCase{
		{
			name:   "create campaign success",
			method: http.MethodPost,
			url:    "/campaigns",
			body:   `{"name":"Campaign Name","content":"Body content here","emails":["user@test.com"]}`,
			setupMock: func(m *ServiceMock) {
				expected := &campaign.CampaignResponse{Id: uuid.New(), Name: "Campaign Name"}
				m.On("Create", mock.MatchedBy(func(c *campaign.CreateCampaignRequest) bool {
					return c.Name == "Campaign Name"
				})).Return(expected, nil)
			},
			wantStatus: http.StatusCreated,
			wantBody: func(t *testing.T, w *httptest.ResponseRecorder, m *ServiceMock) {
				var result campaign.CampaignResponse
				json.Unmarshal(w.Body.Bytes(), &result)
				assert.NotEqual(t, uuid.Nil, result.Id)
				assert.Equal(t, "Campaign Name", result.Name)
				m.AssertExpectations(t)
			},
		},
		{
			name:   "create campaign error",
			method: http.MethodPost,
			url:    "/campaigns",
			body:   `{"name":"Campaign Name","content":"Body content here","emails":["user@test.com"]}`,
			setupMock: func(m *ServiceMock) {
				m.On("Create", mock.Anything).Return(nil, errors.New("service error"))
			},
			wantStatus: http.StatusBadRequest,
			wantBody: func(t *testing.T, w *httptest.ResponseRecorder, m *ServiceMock) {
				var result map[string]string
				json.Unmarshal(w.Body.Bytes(), &result)
				assert.Equal(t, "service error", result["error"])
				m.AssertExpectations(t)
			},
		},
		{
			name:   "list campaigns success",
			method: http.MethodGet,
			url:    "/campaigns",
			body:   "",
			setupMock: func(m *ServiceMock) {
				expected := &[]campaign.CampaignResponse{
					{Id: uuid.New(), Name: "Campaign 1"},
					{Id: uuid.New(), Name: "Campaign 2"},
				}
				m.On("Get").Return(expected, nil)
			},
			wantStatus: http.StatusOK,
			wantBody: func(t *testing.T, w *httptest.ResponseRecorder, m *ServiceMock) {
				var result []campaign.CampaignResponse
				json.Unmarshal(w.Body.Bytes(), &result)
				assert.Len(t, result, 2)
				assert.Equal(t, "Campaign 1", result[0].Name)
				assert.Equal(t, "Campaign 2", result[1].Name)
				m.AssertExpectations(t)
			},
		},
		{
			name:   "list campaigns error",
			method: http.MethodGet,
			url:    "/campaigns",
			body:   "",
			setupMock: func(m *ServiceMock) {
				m.On("Get").Return(nil, errors.New("service error"))
			},
			wantStatus: http.StatusBadRequest,
			wantBody: func(t *testing.T, w *httptest.ResponseRecorder, m *ServiceMock) {
				var result map[string]string
				json.Unmarshal(w.Body.Bytes(), &result)
				assert.Equal(t, "service error", result["error"])
				m.AssertExpectations(t)
			},
		},
		{
			name:   "show campaign success",
			method: http.MethodGet,
			url:    "/campaigns/" + uuid.New().String(),
			body:   `{"id":"` + uuid.New().String() + `"}`,
			setupMock: func(m *ServiceMock) {
				m.On("Show", mock.MatchedBy(func(types.UUID) bool {
					return true
				})).Return(&campaign.CampaignResponse{Id: uuid.New(), Name: "Campaign Name"}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody: func(t *testing.T, w *httptest.ResponseRecorder, m *ServiceMock) {
				var result campaign.CampaignResponse
				json.Unmarshal(w.Body.Bytes(), &result)
				assert.NotEqual(t, uuid.Nil, result.Id)
				assert.Equal(t, "Campaign Name", result.Name)
				m.AssertExpectations(t)
			},
		},
		{
			name:   "show campaign error",
			method: http.MethodGet,
			url:    "/campaigns/" + uuid.New().String(),
			body:   `{"id":"` + uuid.New().String() + `"}`,
			setupMock: func(m *ServiceMock) {
				m.On("Show", mock.Anything).Return(nil, errors.New("service error"))
			},
			wantStatus: http.StatusBadRequest,
			wantBody: func(t *testing.T, w *httptest.ResponseRecorder, m *ServiceMock) {
				var result map[string]string
				json.Unmarshal(w.Body.Bytes(), &result)
				assert.Equal(t, "service error", result["error"])
				m.AssertExpectations(t)
			},
		},
		{
			name:   "cancel campaign success",
			method: http.MethodPatch,
			url:    "/campaigns/cancel/" + uuid.New().String(),
			body:   `{"id":"` + uuid.New().String() + `"}`,
			setupMock: func(m *ServiceMock) {
				m.On("Cancel", mock.MatchedBy(func(types.UUID) bool {
					return true
				})).Return(&campaign.CampaignResponse{Id: uuid.New(), Name: "Campaign Name"}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody: func(t *testing.T, w *httptest.ResponseRecorder, m *ServiceMock) {
				var result campaign.CampaignResponse
				json.Unmarshal(w.Body.Bytes(), &result)
				assert.NotEqual(t, uuid.Nil, result.Id)
				assert.Equal(t, "Campaign Name", result.Name)
				m.AssertExpectations(t)
			},
		},
		{
			name:   "cancel campaign error",
			method: http.MethodPatch,
			url:    "/campaigns/cancel/" + uuid.New().String(),
			body:   `{"id":"` + uuid.New().String() + `"}`,
			setupMock: func(m *ServiceMock) {
				m.On("Cancel", mock.Anything).Return(nil, errors.New("service error"))
			},
			wantStatus: http.StatusBadRequest,
			wantBody: func(t *testing.T, w *httptest.ResponseRecorder, m *ServiceMock) {
				var result map[string]string
				json.Unmarshal(w.Body.Bytes(), &result)
				assert.Equal(t, "service error", result["error"])
				m.AssertExpectations(t)
			},
		},
		{
			name:   "delete campaign success",
			method: http.MethodDelete,
			url:    "/campaigns/" + uuid.New().String(),
			body:   "",
			setupMock: func(m *ServiceMock) {
				m.On("Delete", mock.MatchedBy(func(types.UUID) bool {
					return true
				})).Return(nil)
			},
			wantStatus: http.StatusNoContent,
			wantBody: func(t *testing.T, w *httptest.ResponseRecorder, m *ServiceMock) {
				m.AssertExpectations(t)
			},
		},
		{
			name:   "delete campaign error",
			method: http.MethodDelete,
			url:    "/campaigns/" + uuid.New().String(),
			body:   "",
			setupMock: func(m *ServiceMock) {
				m.On("Delete", mock.Anything).Return(errors.New("service error"))
			},
			wantStatus: http.StatusBadRequest,
			wantBody: func(t *testing.T, w *httptest.ResponseRecorder, m *ServiceMock) {
				var result map[string]string
				json.Unmarshal(w.Body.Bytes(), &result)
				assert.Equal(t, "service error", result["error"])
				m.AssertExpectations(t)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			serviceMock := new(ServiceMock)
			tc.setupMock(serviceMock)

			authServiceMock := new(AuthServiceMock)
			authServiceMock.On("VerifyToken", "test-token").Return((*oidc.IDToken)(nil), nil)

			handler := campaign.NewCampaignHandler(serviceMock)
			ctn := &container.Container{CampaignHandler: handler, AuthService: authServiceMock}

			router := gin.New()
			routes.RegisterRoutes(router, ctn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tc.method, tc.url, strings.NewReader(tc.body))
			req.Header.Set("Authorization", "Bearer test-token")
			if tc.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
			tc.wantBody(t, w, serviceMock)
		})
	}
}
