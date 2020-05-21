package controller

import (
	"cicd/db"
	"cicd/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"strconv"
)

type ReposConfig struct {
	Url     string `json:"url" bson:"url"`
	Token   string `json:"token" bson:"token"`
	ID      string `json:"id" bson:"id"`
	Version string `json:"version" bson:"version"`
}

type UserRepos struct {
	User  string `json:"user" bson:"user"`
	Repos string `json:"repos" bson:"repos"`
}

func GetReposConfig(c *gin.Context) {
	userId, exit := c.Get("userId")
	if !exit {
		c.JSON(http.StatusOK, utils.Ok(""))
		return
	}
	repos := &UserRepos{}
	err := db.DbClient.C("user_repos_config").Find(bson.M{"user": userId}).One(repos)
	if err != nil {
		c.JSON(http.StatusOK, utils.Ok(""))
		return
	}

	req := &ReposConfig{}

	err = db.DbClient.C("repos_config").Find(bson.M{"id": repos.Repos}).One(req)

	if err != nil {
		c.JSON(http.StatusOK, utils.Ok(""))
		return
	}

	c.JSON(http.StatusOK, utils.Ok(req))
}

func SaveReposConfig(c *gin.Context) {

	req := &ReposConfig{}

	err := c.BindJSON(req)

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, utils.ErrorWithMessage("参数校验失败！"))
		return
	}

	req.ID = strconv.FormatUint(utils.GetIntId(), 10)
	userId, exit := c.Get("userId")

	if !exit {
		c.JSON(http.StatusServiceUnavailable, utils.ErrorWithMessage("参数校验失败！:userId"))
		return
	}

	err = db.DbClient.C("user_repos_config").
		Insert(bson.M{
			"user":  userId,
			"repos": req.ID,
		})

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, utils.ErrorWithMessage("保存失败！"))
		return
	}

	err = db.DbClient.C("repos_config").
		Insert(req)

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, utils.ErrorWithMessage("保存失败！"))
		return
	}

	c.JSON(http.StatusOK, utils.Ok(req))
}

func UpdateReposConfig(c *gin.Context) {
	req := &ReposConfig{}

	err := c.BindJSON(req)

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, utils.ErrorWithMessage("参数校验失败！"))
		return
	}

	err = db.DbClient.C("repos_config").Update(bson.M{"id": req.ID}, req)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, utils.ErrorWithMessage("更新失败！"))
		return
	}

	c.JSON(http.StatusOK, utils.Ok(req.ID))
}

func GetReposProjects(c *gin.Context) {
	repos := c.Query("repos")
	if repos == "" {
		c.JSON(http.StatusServiceUnavailable, utils.ErrorWithMessage("参数校验失败！"))
		return
	}

	repoConfig := &ReposConfig{}

	err := db.DbClient.C("repos_config").Find(bson.M{
		"id": repos,
	}).One(repoConfig)

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, utils.ErrorWithMessage("仓库配置加载失败！"))
		return
	}

	url := repoConfig.Url + "/api/v" + repoConfig.Version + "/projects?private_token=" + repoConfig.Token

	resp, err := http.Get(url)

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, utils.ErrorWithMessage("数据加载失败！"))
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 200 {
		c.JSON(http.StatusOK, utils.Ok(string(body)))
	}
}
