package controller

import (
	"cicd/db"
	"cicd/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type UserDetail struct {
	Repos   string `json:"repos" bson:"repos"`
	Account string `json:"account" bson:"account"`
}

func GetUserDetail(c *gin.Context) {
	userId, exit := c.Get("userId")

	if !exit {
		c.JSON(http.StatusUnauthorized, utils.Error())
		return
	}

	resp := &UserDetail{}

	detail := &UserDetail{}

	err := db.DbClient.C("user").Find(bson.M{"userId": userId}).One(&detail)

	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.Error())
		return
	}

	resp.Account = detail.Account

	err = db.DbClient.C("user_repos_config").Find(bson.M{"user": userId}).One(&detail)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.Error())
		return
	}

	resp.Repos = detail.Repos

	c.JSON(http.StatusOK, utils.Ok(resp))
}
