package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/coreos/go-oidc"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Service struct {
	config   *Config
	client   HTTPClient
	verifier TokenVerifier
}

var _ ServiceInterface = (*Service)(nil)

func NewService(config *Config, client HTTPClient, verifier TokenVerifier) *Service {
	return &Service{config: config, client: client, verifier: verifier}
}

func (s *Service) ExchangeCode(code string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", s.config.ClientID)
	data.Set("client_secret", s.config.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", s.config.CallbackURL)

	tokenURL := s.config.BaseURL + "/" + s.config.TokenURI

	req, err := http.NewRequest(http.MethodPost, tokenURL, strings.NewReader(data.Encode()))

	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to call token endpoint: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token endpoint returned status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResponse TokenResponse

	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokenResponse, nil
}

func (s *Service) VerifyToken(token string) (*oidc.IDToken, error) {
	idToken, err := s.verifier.Verify(context.Background(), token)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	return idToken, nil
}
