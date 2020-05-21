package utils

import (
	"github.com/bndr/gojenkins"
	"net/http"
)

type JenkinsUtil struct {
	Instance *gojenkins.Jenkins
}

func GetInstances(client *http.Client, url string, account string, pwd string) (instance *JenkinsUtil, err error) {
	util := &JenkinsUtil{}
	util.Instance = gojenkins.CreateJenkins(client, url, account, pwd)
	_, err = util.Instance.Init()


	if err != nil {
		return
	}
	return util, err
}
