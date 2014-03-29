package models

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lunny/xorm"
	. "github.com/qiniu/api/conf"
	"github.com/robfig/config"
	"github.com/robfig/revel"
)

var (
	Engine      *xorm.Engine
	Smtp        SmtpType
	QiniuScope  string
	QiniuDomain string
)

type SmtpType struct {
	Username string
	Password string
	Host     string
	Address  string
	From     string
}

func init() {
	revel.OnAppStart(Init)
}

func Init() {
	c, err := config.ReadDefault(revel.BasePath + "/conf/my.conf")

	driver, _ := c.String("database", "db.driver")
	dbname, _ := c.String("database", "db.dbname")
	user, _ := c.String("database", "db.user")
	password, _ := c.String("database", "db.password")
	host, _ := c.String("database", "db.host")

	params := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", user, password, host, dbname)
	Engine, err = xorm.NewEngine(driver, params)
	// defer Engine.Close()

	if err != nil {
		panic(err)
	}

	err = Engine.Sync(
		new(User),
		new(Category),
		new(Topic),
		new(Reply),
		new(Permissions),
	)

	if err != nil {
		panic(err)
	}

	// 如果是空数据库，自动添加管理员账号 admin/123
	count, _ := Engine.Count(new(User))
	if count == 0 {
		Engine.Insert(&User{
			Name:           "admin",
			Email:          "admin@admin.com",
			Avatar:         DefaultAvatar,
			Type:           1,
			IsActive:       true,
			HashedPassword: EncryptPassword("123"),
		})

		Engine.Insert(
			&Permissions{UserId: 1, Perm: 1},
			&Permissions{UserId: 1, Perm: 2},
		)
	}

	// Engine.ShowSQL = revel.DevMode
	// Engine.ShowDebug = revel.DevMode
	// Engine.ShowWarn = revel.DevMode
	Engine.ShowErr = revel.DevMode

	Smtp.Username, _ = c.String("smtp", "smtp.username")
	Smtp.Password, _ = c.String("smtp", "smtp.password")
	Smtp.Address, _ = c.String("smtp", "smtp.address")
	Smtp.From, _ = c.String("smtp", "smtp.from")
	Smtp.Host, _ = c.String("smtp", "smtp.host")

	ACCESS_KEY, _ = c.String("qiniu", "access_key")
	SECRET_KEY, _ = c.String("qiniu", "secret_key")
	QiniuScope, _ = c.String("qiniu", "scope")
	QiniuDomain, _ = c.String("qiniu", "qiniuDomain")
}
