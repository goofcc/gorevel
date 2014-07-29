package controllers

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/revel/revel/cache"
	"image"
	"io"
	"mime/multipart"
	"net/smtp"
	"os"
	"path"
	"strings"
	// "github.com/go-xorm/xorm"
	"github.com/nashtsai/xormrevelmodule"
	qio "github.com/qiniu/api/io"
	"github.com/qiniu/api/rs"
	"github.com/revel/revel"

	"gorevel/app/models"
	"gorevel/app/routes"
)

var (
	//engine           *xorm.Engine
	UPLOAD_PATH      string
	IMAGE_EXTS       string = ".jpg.jpeg.png"
	IMAGE_LIMIT_SIZE int64  = 500 * 1024
)

type Vars map[string]interface{}

type Sizer interface {
	Size() int64
}

type Application struct {
	*revel.Controller
	xormmodule.XormController
}

func (c *Application) checkUser() revel.Result {
	user := c.user()
	if user != nil {
		c.RenderArgs["user"] = user
	}

	// 检查是否需要授权
	action := strings.TrimSuffix(c.Action, "Post")
	if value, needCheck := Permissions[action]; needCheck {
		if user == nil {
			c.Flash.Error("请先登录")
			c.Session["preUrl"] = c.Request.Request.URL.String()
			return c.Redirect(routes.User.Signin())
		} else {
			perm := user.GetPermissions()
			if _, ok := perm[value]; !ok {
				return c.Forbidden("抱歉，您没有得到授权！")
			}
		}
	}

	c.addExtraArgs()

	return nil
}

func (c *Application) user() *models.User {
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
	has, _ := c.Engine.Where("name = ?", username).Get(&user)
	if !has {
		return nil
	}

	return &user
}

func (c *Application) bindVars(vars Vars) {
	for k, v := range vars {
		c.RenderArgs[k] = v
	}
}

func (c *Application) addExtraArgs() {
	c.bindVars(Vars{
		"active":      c.Name,
		"action":      c.Action,
		"qiniuDomain": models.QiniuDomain,
		"categories":  getCategories(),
	})
}

func saveFile(file *multipart.File, filePath string) error {
	os.MkdirAll(path.Dir(filePath), 0777)

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err == nil {
		defer f.Close()
		_, err = io.Copy(f, *file)
	} else {
		revel.ERROR.Println(err)
	}

	return err
}

func thumbFile(filePath string) {
	src, _ := imaging.Open(filePath)
	var dst *image.NRGBA

	dst = imaging.Thumbnail(src, 48, 48, imaging.CatmullRom)
	imaging.Save(dst, filePath)
}

func deleteFile(filePath string) error {
	if fileExist(filePath) {
		if err := os.Remove(filePath); err != nil {
			revel.ERROR.Println(err)
		}
	}

	return nil
}

func fileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || os.IsExist(err)
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

func checkImageExt(v *revel.Validation, file *multipart.File, header *multipart.FileHeader, formField string) bool {
	if !strings.Contains(IMAGE_EXTS, strings.ToLower(path.Ext(header.Filename))) {
		err := &revel.ValidationError{
			Message: "只能上传图片",
			Key:     formField,
		}
		v.Errors = append(v.Errors, err)
		return false
	}

	if (*file).(Sizer).Size() > IMAGE_LIMIT_SIZE {
		err := &revel.ValidationError{
			Message: fmt.Sprintf("图片尺寸不超过%dKB", IMAGE_LIMIT_SIZE),
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

func qiniuUploadImage(file *multipart.File, fileName string) (error, qio.PutRet) {
	var ret qio.PutRet
	var policy = rs.PutPolicy{
		Scope: models.QiniuScope,
	}
	err := qio.Put(nil, &ret, policy.Token(nil), fileName, *file, nil)
	if err != nil {
		revel.ERROR.Println("io.Put failed:", err)
	}

	return err, ret
}

func qiniuDeleteImage(fileName string) error {
	var client rs.Client
	client = rs.New(nil)

	err := client.Delete(nil, models.QiniuScope, fileName)
	if err != nil {
		revel.ERROR.Println("rs.Delete failed:", err)
	}

	return err
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

func getCategories() []models.Category {
	var categories []models.Category
	if err := cache.Get("categories", &categories); err != nil {
		xormmodule.Engine.Find(&categories)
		go cache.Set("categories", categories, cache.FOREVER)
	}

	return categories
}
