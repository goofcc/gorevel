package controllers

import (
	"fmt"
	_ "github.com/coocood/mysql"
	"github.com/coocood/qbs"
	"github.com/robfig/config"
	"github.com/robfig/revel"
	"gorevel/app/models"
)

type Qbs struct {
	*revel.Controller
	q *qbs.Qbs
}

func (c *Qbs) Begin() revel.Result {
	q, err := qbs.GetQbs()
	if err != nil {
		fmt.Println(err)
	}
	c.q = q

	return nil
}

func (c *Qbs) End() revel.Result {
	c.q.Close()

	return nil
}

func Init() {
	basePath = revel.BasePath
	uploadPath = basePath + "/public/upload/"

	c, _ := config.ReadDefault(basePath + "/conf/my.conf")
	driver, _ := c.String("database", "db.driver")
	dbname, _ := c.String("database", "db.dbname")
	user, _ := c.String("database", "db.user")
	password, _ := c.String("database", "db.password")
	host, _ := c.String("database", "db.host")

	registerDb(driver, dbname, user, password, host)
}

func registerDb(driver, dbname, user, password, host string) {
	params := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", user, password, host, dbname)
	qbs.Register(driver, params, dbname, qbs.NewMysql())
	err := createTabel()
	if err != nil {
		fmt.Println(err)
	}
}

func createTabel() error {
	migration, err := qbs.GetMigration()
	if err != nil {
		return err
	}
	defer migration.Close()

	err = migration.CreateTableIfNotExists(new(models.User))
	err = migration.CreateTableIfNotExists(new(models.Category))
	err = migration.CreateTableIfNotExists(new(models.Topic))
	err = migration.CreateTableIfNotExists(new(models.Reply))
	err = migration.CreateTableIfNotExists(new(models.Permissions))

	return err
}
