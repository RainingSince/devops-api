package middleware

import (
	"cicd/config"
	"cicd/db"
	"cicd/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strings"
)

func Author() gin.HandlerFunc {
	return authPath
}

type AuthUser struct {
	UserID string `json:"userId"`
}

/*
* 验证 token 并添加 userId
*/
func authPath(c *gin.Context) {

	methods := c.Request.Method
	if methods == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	path := c.FullPath()
	ingore := config.GlobabConfig.AuthConfig.IngorePath

	if strings.Contains(ingore, path) {
		c.Next()
		return
	}

	token := c.Request.Header.Get("Authorization")

	if token == "" {
		c.JSON(http.StatusUnauthorized, utils.ErrorWithMessage("接口认证失败"))
		return
	}

	tokeDetail, err := utils.ParseToken(token)

	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorWithMessage(err.Error()))
		return
	}

	user := db.DbClient.C("user")

	result := AuthUser{}

	err = user.Find(bson.M{"userId": tokeDetail.ID}).One(&result)

	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorWithMessage(err.Error()))
		return
	}

	c.Set("userId", tokeDetail.ID)

	c.Next()
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		var headerKeys []string
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
			c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json")
		}


		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}