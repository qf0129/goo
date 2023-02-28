package simple

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qf0129/goo"
	"github.com/qf0129/goo/crud"
)

type UserReqBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func UserLoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UserReqBody
		var existUser User

		if err := c.ShouldBindJSON(&req); err != nil {
			crud.RespFail(c, "InvalidJsonData, "+err.Error())
			return
		}

		err := goo.DB.Where(map[string]any{"username": req.Username}).First(&existUser).Error
		if err != nil {
			crud.RespFail(c, "FindNotUser, "+err.Error())
			return
		}

		result := goo.VerifyPassword(req.Password, existUser.PasswordHash)
		if !result {
			crud.RespFail(c, "InvalidPassword")
			return
		}
		token, err := goo.CreateToken(strconv.Itoa(int(existUser.Id)))
		if err != nil {
			crud.RespFail(c, "CreateTokenErr, "+err.Error())
			return
		}
		c.SetCookie("tk", token, int(goo.Config.TokenExpiredTime), "/", "*", true, true)
		c.SetCookie("username", existUser.Username, int(goo.Config.TokenExpiredTime), "/", "*", true, false)
		crud.RespOk(c, gin.H{"token": token})
	}
}

func UserRegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UserReqBody
		var existUser User

		if err := c.ShouldBindJSON(&req); err != nil {
			crud.RespFail(c, "InvalidJsonData, "+err.Error())
			return
		}

		goo.DB.Where(map[string]any{"username": req.Username}).First(&existUser)
		if existUser.Id > 0 {
			crud.RespFail(c, "UserAleardyExists")
			return
		}

		psdHash, err := goo.HashPassword(req.Password)
		if err != nil {
			crud.RespFail(c, "InvalidPassword, "+err.Error())
			return
		}

		u := &User{
			Username:     req.Username,
			PasswordHash: psdHash,
		}

		err = crud.CreateOne[User](u)
		if err != nil {
			crud.RespFail(c, "CreateUserErr, "+err.Error())
			return
		}

		crud.RespOk(c, gin.H{"id": u.Id})
	}
}
