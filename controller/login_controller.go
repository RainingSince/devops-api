package controller

import (
	"cicd/db"
	"cicd/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type UserInfo struct {
	Account  string `json:"account" bason:"account"`
	Password string `json:"password" bason:"password"`
	UserId   string `json:"userId" bson:"userId"`
	Token    string `json:"token" bson:"token"`
}

func Login(c *gin.Context) {

	req := &UserInfo{}
	user := &UserInfo{}

	err := c.BindJSON(req)

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, utils.ErrorWithMessage("参数校验失败！"))
		return
	}

	err = db.DbClient.C("user").
		Find(bson.M{
			"account":  req.Account,
			"password": req.Password,
		}).
		One(&user)

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, utils.ErrorWithMessage("登录失败！"))
		return
	}

	token, err := utils.CreateToken(&utils.TokenClaims{ID: user.UserId})

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, utils.ErrorWithMessage("Token生成失败！"))
		return
	}

	c.JSON(http.StatusOK, utils.Ok(token))
}
