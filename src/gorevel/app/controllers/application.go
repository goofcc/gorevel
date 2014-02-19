package controllers

import (
	"image"
	"io"
	"mime/multipart"
	"net/smtp"
	"os"
	"path"
	"strings"

	"code.google.com/p/go-uuid/uuid"
	"github.com/disintegration/imaging"
	"github.com/lunny/xorm"
	"github.com/robfig/revel"

	"gorevel/app/models"
	"gorevel/app/routes"
)

type Vars map[string]interface{}

var (
	engine      *xorm.Engine
	UPLOAD_PATH string
	imageExts   string = ".jpg.jpeg.png"
)

type Application struct {
	*revel.Controller
	userId int64
}

func (c *Application) injector() revel.Result {
	c.RenderArgs["active"] = c.Name
	user := c.connected()
	if user != nil {
		c.RenderArgs["user"] = user
		c.userId = user.Id
	} else {
		c.userId = 0
	}

	// 检查是否需要授权
	value, ok := Permissions[strings.TrimSuffix(c.Action, "Post")]
	if ok {
		if user == nil {
			c.Flash.Error("请先登录")
			c.Session["preUrl"] = c.Request.Request.URL.String()
			return c.Redirect(routes.User.Signin())
		} else {
			perm := user.GetPermissions()
			_, ok := perm[value]
			if !ok {
				return c.Forbidden("抱歉，您没有得到授权！")
			}
		}
	}

	return nil
}

func (c *Application) connected() *models.User {
	if c.RenderArgs["user"] != nil {
		return c.RenderArgs["user"].(*models.User)
	}

	if username, ok := c.Session["user"]; ok {
		return c.getUser(username)
	}

	return nil
}

func (c *Application) getUser(username string) *models.User {
	var user models.User
	has, _ := engine.Where("name = ?", username).Get(&user)

	if !has {
		return nil
	}

	return &user
}

func (c *Application) Vars(args Vars) {
	for k, v := range args {
		c.RenderArgs[k] = v
	}
}

func saveFile(file *multipart.File, filePath string) error {
	os.MkdirAll(path.Dir(filePath), 0777)

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()

	if err != nil {
		revel.ERROR.Println(err)
	} else {
		_, err = io.Copy(f, *file)
	}

	return err
}

func thumbFile(filePath string) {
	src, _ := imaging.Open(filePath)
	var dst *image.NRGBA

	dst = imaging.Thumbnail(src, 48, 48, imaging.CatmullRom)
	imaging.Save(dst, filePath)
}

func deleteFile(filepath string) error {
	err := os.Remove(filepath)

	if err != nil {
		revel.ERROR.Println(err)
	}

	return err
}

func checkFileExt(v *revel.Validation, header *multipart.FileHeader, fileExts, formField, message string) bool {
	if !strings.Contains(fileExts, strings.ToLower(path.Ext(header.Filename))) {
		err := &revel.ValidationError{
			Message: message,
			Key:     formField,
		}
		v.Errors = append(v.Errors, err)
		return false
	}
	return true
}

func uuidFileName(fileName string) string {
	return strings.Replace(uuid.NewUUID().String(), "-", "", -1) + path.Ext(fileName)
}

func uuidName() string {
	return strings.Replace(uuid.NewUUID().String(), "-", "", -1)
}

func sendMail(subject, content string, tos []string) error {
	message := `From: Revel中文社区
To: ` + strings.Join(tos, ",") + `
Subject: ` + subject + `
Content-Type: text/html;charset=UTF-8

` + content

	Smtp := models.Smtp
	auth := smtp.PlainAuth("", Smtp.Username, Smtp.Password, Smtp.Host)
	err := smtp.SendMail(Smtp.Address, auth, Smtp.From, tos, []byte(message))
	if err != nil {
		revel.ERROR.Println(err)
		return err
	}

	return nil
}
