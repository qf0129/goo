package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/goo"
	"github.com/qf0129/goo/crud"
	"github.com/qf0129/goo/model"
)

type getTokenBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetTokenHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req getTokenBody
		var existUser model.BaseUser

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
		token, err := goo.CreateToken(req.Username)
		if err != nil {
			crud.RespFail(c, "CreateTokenErr, "+err.Error())
			return
		}
		crud.RespOk(c, gin.H{"token": token})
	}
}
