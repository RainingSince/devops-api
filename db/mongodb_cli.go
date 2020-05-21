package db

import (
	"cicd/config"
	"gopkg.in/mgo.v2"
)

var DbClient *mgo.Database

func DBInit() (err error) {
	session, err := mgo.Dial(config.GlobabConfig.DbConfig.Hosts)

	if err != nil {
		return err
	}

	session.SetMode(mgo.Monotonic, true)

	DbClient = session.DB(config.GlobabConfig.DbConfig.DataBase)
	return err
}
