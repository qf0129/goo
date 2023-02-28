package simple

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
			crud.RespFailWith401(c, err.Error())
			c.Abort()
			return
		}

		userIdStr, err := goo.ParseToken(tk)
		if err != nil {
			crud.RespFailWith401(c, err.Error())
			c.Abort()
			return
		}

		existsUser := &User{}
		err = crud.QueryOneTarget[User](userIdStr, existsUser)
		if err != nil {
			crud.RespFailWith401(c, err.Error())
			c.Abort()
			return
		}
		c.Set("user", existsUser)
		c.Next()
	}
}
