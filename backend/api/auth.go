package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

const sessionCookieName = "whishper_session"

// authEnabled reports whether API authentication is turned on. Auth is enabled
// only when an API key is configured, keeping older deployments backwards
// compatible (no key -> open API, same as before).
func authEnabled() bool {
	return os.Getenv("WHISHPER_API_KEY") != ""
}

// signToken creates an HMAC-signed session token for the given subject. The
// token format is: base64url(payload).base64url(hmac_sha256(payload)).
func signToken(subject string) string {
	secret := os.Getenv("WHISHPER_API_KEY")
	exp := time.Now().Add(7 * 24 * time.Hour).Unix()
	payload := fmt.Sprintf(`{"sub":%q,"exp":%d}`, subject, exp)
	b64p := base64.RawURLEncoding.EncodeToString([]byte(payload))

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(b64p))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return b64p + "." + sig
}

// validateToken verifies the signature and expiry of a session token.
func validateToken(token string) bool {
	secret := os.Getenv("WHISHPER_API_KEY")
	if secret == "" || token == "" {
		return false
	}

	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(parts[0]))
	expected := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(parts[1])) {
		return false
	}

	raw, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}
	var p struct {
		Exp int64 `json:"exp"`
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return false
	}
	return time.Now().Unix() <= p.Exp
}

// RequireAuth is a Fiber middleware that protects the API. A request is allowed
// when it carries either a valid X-API-Key header (programmatic access) or a
// valid session cookie (web UI access). The /api/auth/* routes are always open.
func (s *Server) RequireAuth(c *fiber.Ctx) error {
	path := c.Path()

	// Login/logout endpoints must stay reachable.
	if strings.HasPrefix(path, "/api/auth") {
		return c.Next()
	}

	// Auth disabled -> behave exactly as before.
	if !authEnabled() {
		return c.Next()
	}

	// Programmatic access via API key.
	if key := c.Get("X-API-Key"); key != "" {
		if hmac.Equal([]byte(key), []byte(os.Getenv("WHISHPER_API_KEY"))) {
			return c.Next()
		}
	}

	// Web UI access via session cookie.
	if validateToken(c.Cookies(sessionCookieName)) {
		return c.Next()
	}

	return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
}

// handleLogin validates credentials against the configured env vars and, on
// success, issues a session cookie.
func (s *Server) handleLogin(c *fiber.Ctx) error {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(c.Body(), &body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Bad request")
	}

	wantUser := os.Getenv("WHISHPER_USERNAME")
	wantPass := os.Getenv("WHISHPER_PASSWORD")
	if wantUser == "" || wantPass == "" {
		log.Warn().Msg("Login attempted but WHISHPER_USERNAME/WHISHPER_PASSWORD are not configured")
		return fiber.NewError(fiber.StatusInternalServerError, "Login is not configured")
	}

	userOk := hmac.Equal([]byte(body.Username), []byte(wantUser))
	passOk := hmac.Equal([]byte(body.Password), []byte(wantPass))
	if !userOk || !passOk {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
	}

	token := signToken(wantUser)
	c.Cookie(&fiber.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   7 * 24 * 3600,
		HTTPOnly: true,
		SameSite: "Lax",
	})
	return c.JSON(fiber.Map{"token": token})
}

// handleLogout clears the session cookie.
func (s *Server) handleLogout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HTTPOnly: true,
		SameSite: "Lax",
	})
	return c.SendStatus(fiber.StatusOK)
}
