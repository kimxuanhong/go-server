package jwt

import (
	"github.com/kimxuanhong/go-server/core"
	"strings"
)

const UserInfoKey = "userInfo"

// AuthMiddleware
// Example
// user := c.Get(UserInfoKey).(*UserInfo)
func AuthMiddleware(jwtComp *Jwt) core.Handler {
	return func(c core.Context) {
		authHeader := c.Header("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(core.StatusUnauthorized, map[string]interface{}{"error": "missing token"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		user, err := jwtComp.Validate(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(core.StatusUnauthorized, map[string]interface{}{"error": "invalid token"})
			return
		}

		c.Set(UserInfoKey, user)
		c.Next()
	}
}
