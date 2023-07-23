package tmpl

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/goo"
	"github.com/qf0129/goo/crud"
)

const TokenKey = "tk"

func RequireTokenFromCookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		tk, err := c.Cookie(TokenKey)
		if err != nil {
			goo.RespFailWith401(c, "InvalidToken: "+err.Error())
			c.Abort()
			return
		}

		userIdStr, err := goo.ParseToken(tk)
		if err != nil {
			goo.RespFailWith401(c, "InvalidToken: "+err.Error())
			c.Abort()
			return
		}

		existsUser := &User{}
		err = crud.QueryOneTarget[User](userIdStr, existsUser)
		if err != nil {
			goo.RespFailWith401(c, "InvalidUser: "+err.Error())
			c.Abort()
			return
		}
		c.Set("user", existsUser)
		c.Next()
	}
}

// 跨域请求
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
