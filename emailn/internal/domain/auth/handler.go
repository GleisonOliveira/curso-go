package auth

import (
	"crypto/rand"
	"emailn/cmd/api/handlers"
	"emailn/internal/internalerrors"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	config  *Config
	service ServiceInterface
}

func NewHandler(config *Config, service ServiceInterface) *Handler {
	return &Handler{config: config, service: service}
}

func (h *Handler) HandleLogin(c *gin.Context) {
	state := generateState()
	authURL := h.config.BaseURL + "/" + h.config.AuthURI + "?" +
		"client_id=" + h.config.ClientID +
		"&redirect_uri=" + h.config.CallbackURL +
		"&response_type=code" +
		"&scope=openid email" +
		"&state=" + state

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// TODO: Verificar se o state recebido corresponde ao state armazenado no HandleLogin.
// Atualmente so valida se esta vazio, mas nao compara com o original.
func (h *Handler) HandleCallback(c *gin.Context) *handlers.ResponseConfig[*TokenResponse] {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		return &handlers.ResponseConfig[*TokenResponse]{
			ErrorStatus: http.StatusBadRequest,
			Error:       internalerrors.AuthInternal,
		}
	}

	token, err := h.service.ExchangeCode(code)

	if err != nil {
		return &handlers.ResponseConfig[*TokenResponse]{
			ErrorStatus: http.StatusBadRequest,
			Error:       err,
		}
	}

	return &handlers.ResponseConfig[*TokenResponse]{
		SuccessStatus: http.StatusOK,
		Data:          token,
	}
}

func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
