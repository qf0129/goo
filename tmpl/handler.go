package tmpl

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qf0129/goo"
	"github.com/qf0129/goo/crud"
)

type AuthRequestBody struct {
	Username string
	Password string
}

func AuthLoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AuthRequestBody
		var existUser User

		if err := c.ShouldBindJSON(&req); err != nil {
			goo.RespFail(c, "InvalidJsonData, "+err.Error())
			return
		}

		if req.Username == "" {
			goo.RespFail(c, "InvalidUsername")
			return
		}

		if req.Password == "" {
			goo.RespFail(c, "InvalidPassword")
			return
		}

		err := goo.DB.Where(map[string]any{"username": req.Username}).First(&existUser).Error
		if err != nil {
			goo.RespFail(c, "FindNotUser, "+err.Error())
			return
		}

		result := goo.VerifyPassword(req.Password, existUser.PasswordHash)
		if !result {
			goo.RespFail(c, "InvalidPassword")
			return
		}
		token, err := goo.CreateToken(strconv.Itoa(int(existUser.Id)))
		if err != nil {
			goo.RespFail(c, "CreateTokenErr, "+err.Error())
			return
		}
		c.SetCookie("tk", token, int(goo.Config.TokenExpiredTime), "/", "*", true, true)
		c.SetCookie("username", existUser.Username, int(goo.Config.TokenExpiredTime), "/", "*", true, false)
		goo.RespOk(c, gin.H{"token": token})
	}
}

func AuthRegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AuthRequestBody
		var existUser User

		if err := c.ShouldBindJSON(&req); err != nil {
			goo.RespFail(c, "InvalidJsonData, "+err.Error())
			return
		}

		if req.Username == "" {
			goo.RespFail(c, "InvalidUsername")
			return
		}

		if req.Password == "" {
			goo.RespFail(c, "InvalidPassword")
			return
		}

		goo.DB.Where(map[string]any{"username": req.Username}).First(&existUser)
		if existUser.Id > 0 {
			goo.RespFail(c, "UserAleardyExists")
			return
		}

		psdHash, err := goo.HashPassword(req.Password)
		if err != nil {
			goo.RespFail(c, "InvalidPassword, "+err.Error())
			return
		}

		u := &User{
			Username:     req.Username,
			PasswordHash: psdHash,
		}

		err = crud.CreateOne[User](u)
		if err != nil {
			goo.RespFail(c, "CreateUserErr, "+err.Error())
			return
		}

		goo.RespOk(c, gin.H{"id": u.Id})
	}
}
