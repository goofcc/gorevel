package controllers

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/coocood/qbs"
	"github.com/disintegration/imaging"
	"github.com/robfig/revel"
	"gorevel/app/models"
	"gorevel/app/routes"
	"image"
	"strings"
)

type User struct {
	Application
}

func (c *User) Signup() revel.Result {
	return c.Render()
}

func (c *User) SignupPost(user models.User) revel.Result {
	user.Validate(c.q, c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.User.Signup())
	}

	user.Type = MemberGroup
	user.Avatar = defaultAvatar
	user.ValidateCode = strings.Replace(uuid.NewUUID().String(), "-", "", -1)

	if !user.Save(c.q) {
		c.Flash.Error("注册用户失败")
		return c.Redirect(routes.User.Signup())
	}

	subject := "激活账号 —— Revel社区"
	content := `<h2><a href="http://gorevel.cn/user/validate/` + user.ValidateCode + `">激活账号</a></h2>`
	go sendMail(subject, content, []string{user.Email})

	c.Flash.Success(fmt.Sprintf("%s 注册成功，请到您的邮箱 %s 激活账号！", user.Name, user.Email))

	perm := new(models.Permissions)
	perm.UserId = user.Id
	perm.Perm = MemberGroup
	perm.Save(c.q)

	return c.Redirect(routes.User.Signin())
}

func (c *User) Signin() revel.Result {
	return c.Render()
}

func (c *User) SigninPost(name, password string) revel.Result {
	c.Validation.Required(name).Message("请输入用户名")
	c.Validation.Required(password).Message("请输入密码")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.User.Signin())
	}

	user := new(models.User)
	condition := qbs.NewCondition("name = ?", name).
		And("hashed_password = ?", models.EncryptPassword(password))
	c.q.Condition(condition).Find(user)

	if user.Id == 0 {
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Out["user"] = name
		c.Flash.Error("用户名或密码错误")
		return c.Redirect(routes.User.Signin())
	}

	if !user.IsActive {
		c.Flash.Error(fmt.Sprintf("您的账号 %s 尚未激活，请到您的邮箱 %s 激活账号！", user.Name, user.Email))
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.User.Signin())
	}

	c.Session["user"] = name

	preUrl, ok := c.Session["preUrl"]
	if ok {
		return c.Redirect(preUrl)
	}

	return c.Redirect(routes.App.Index())
}

func (c *User) Signout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}

	return c.Redirect(routes.App.Index())
}

func (c *User) Edit() revel.Result {
	id := c.RenderArgs["user"].(*models.User).Id
	user := findUserById(c.q, id)
	if user.Id == 0 {
		return c.NotFound("用户不存在")
	}

	return c.Render(user, avatars)
}

func (c *User) EditPost(avatar string) revel.Result {
	id := c.RenderArgs["user"].(*models.User).Id
	checkFileExt(c.Controller, imageExts, "picture", "Only image")
	user := findUserById(c.q, id)
	if user.Id == 0 {
		return c.NotFound("用户不存在")
	}

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.User.Edit())
	}

	if ok, _ := getFileExt(c.Request, "picture"); ok {
		picture := saveFile(c.Request, "picture")
		src, _ := imaging.Open(uploadPath + picture)
		var dst *image.NRGBA

		dst = imaging.Thumbnail(src, 48, 48, imaging.CatmullRom)
		avatar = "thumb" + picture
		imaging.Save(dst, uploadPath+avatar)
		deleteFile(picture)
	}

	if avatar != "" {
		if strings.HasPrefix(user.Avatar, "thumb") {
			deleteFile(user.Avatar)
		}
		user.Avatar = avatar
	}

	if user.Save(c.q) {
		c.Flash.Success("保存成功")
	} else {
		c.Flash.Error("保存信息失败")
	}

	return c.Redirect(routes.User.Edit())
}

func (c *User) Validate(code string) revel.Result {
	user := findUserByCode(c.q, code)
	if user.Id == 0 {
		return c.NotFound("用户不存在或校验码错误")
	}

	user.IsActive = true
	user.Save(c.q)

	c.Flash.Success("您的账号成功激活，请登录！")

	return c.Redirect(routes.User.Signin())
}

func findUserById(q *qbs.Qbs, id int64) *models.User {
	user := new(models.User)
	q.WhereEqual("id", id).Find(user)

	return user
}

func findUserByCode(q *qbs.Qbs, code string) *models.User {
	user := new(models.User)
	q.WhereEqual("validate_code", code).Find(user)

	return user
}
