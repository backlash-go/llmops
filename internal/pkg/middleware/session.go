package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	goredis "github.com/go-redis/redis/v7"
	"github.com/marmotedu/component-base/pkg/core"
	"github.com/marmotedu/errors"

	"llmops/internal/apiserver/store/redis"
	"llmops/internal/pkg/code"
	"llmops/internal/pkg/session"
)

// CookieSession authenticates API requests with the server-side browser session.
func CookieSession(store redis.RStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		if store == nil {
			core.WriteResponse(c, errors.WithCode(code.ErrDatabase, "redis store is not initialized"), nil)
			c.Abort()

			return
		}

		sessionID, err := c.Cookie(session.CookieName)
		if err != nil || strings.TrimSpace(sessionID) == "" {
			core.WriteResponse(c, errors.WithCode(code.ErrMissingHeader, "session cookie cannot be empty."), nil)
			c.Abort()

			return
		}

		payload, err := store.Rdb().Get(session.Key(sessionID)).Result()
		if err != nil {
			if err == goredis.Nil {
				core.WriteResponse(c, errors.WithCode(code.ErrExpired, "session expired please to login"), nil)
			} else {
				core.WriteResponse(c, errors.WithCode(code.ErrDatabase, err.Error()), nil)
			}
			c.Abort()

			return
		}

		var data session.Data
		if err := json.Unmarshal([]byte(payload), &data); err != nil {
			core.WriteResponse(c, errors.WithCode(code.ErrDecodingJSON, err.Error()), nil)
			c.Abort()

			return
		}

		if data.ExpiresAt > 0 && time.Now().Unix() > data.ExpiresAt {
			_ = store.Rdb().Del(session.Key(sessionID)).Err()
			clearSessionCookie(c)
			core.WriteResponse(c, errors.WithCode(code.ErrExpired, "session expired"), nil)
			c.Abort()

			return
		}

		c.Set(SessionIDKey, sessionID)
		c.Set(UserIDKey, data.UserID)
		c.Set(UsernameKey, data.Username)
		c.Next()
	}
}

func clearSessionCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(session.CookieName, "", -1, session.CookiePath, "", false, true)
}
