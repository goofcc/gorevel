package controllers

import (
	"fmt"

	"github.com/robfig/revel"

	"gorevel/app/models"
	"gorevel/app/routes"
)

type User struct {
	Application
}

func (c User) Signup() revel.Result {
	return c.Render()
}

func (c User) SignupPost(user models.User) revel.Result {
	user.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.User.Signup())
	}

	salt := uuidName()
	user.Type = MEMBER_GROUP
	user.Avatar = models.DefaultAvatar
	user.ValidateCode = uuidName()
	user.Salt = salt
	user.HashedPassword = models.EncryptPassword(user.Password, salt)

	aff, _ := engine.Insert(&user)
	if aff == 0 {
		c.Flash.Error("注册用户失败")
		return c.Redirect(routes.User.Signup())
	}

	subject := "激活账号 —— Revel中文社区"
	content := `<h2><a href="http://gorevel.cn/user/validate/` + user.ValidateCode + `">激活账号</a></h2>`
	go sendMail(subject, content, []string{user.Email})

	c.Flash.Success(fmt.Sprintf("%s 注册成功，请到您的邮箱 %s 激活账号！", user.Name, user.Email))

	engine.Insert(&models.Permissions{
		UserId: user.Id,
		Perm:   MEMBER_GROUP,
	})

	return c.Redirect(routes.User.Signin())
}

func (c User) Signin() revel.Result {
	return c.Render()
}

func (c User) SigninPost(name, password string) revel.Result {
	c.Validation.Required(name).Message("请输入用户名")
	c.Validation.Required(password).Message("请输入密码")
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.User.Signin())
	}

	var user models.User
	var count int64
	var hashedPassword string
	has, _ := engine.Where("name = ?", name).Get(&user)
	if has {
		// 密码加盐兼容旧密码
		if user.Salt == "" {
			hashedPassword = models.EncryptPassword(password, "")
		} else {
			hashedPassword = models.EncryptPassword(password, user.Salt)
		}

		count, _ = engine.Where("name = ? AND hashed_password = ?", name, hashedPassword).Count(&models.User{})

		// 密码加盐兼容旧密码
		if count > 0 && user.Salt == "" {
			salt := uuidName()
			hashedPassword = models.EncryptPassword(password, salt)
			engine.Id(user.Id).Update(&models.User{
				Salt:           salt,
				HashedPassword: hashedPassword,
			})
		}
	}

	if !has || count == 0 {
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Out["user"] = name
		c.Flash.Error("用户名或密码错误")
		return c.Redirect(routes.User.Signin())
	}

	if !user.IsActive() {
		c.Flash.Error(fmt.Sprintf("您的账号 %s 尚未激活，请到您的邮箱 %s 激活账号！", user.Name, user.Email))
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.User.Signin())
	}

	c.Session["user"] = name

	if preUrl, ok := c.Session["preUrl"]; ok {
		return c.Redirect(preUrl)
	}

	return c.Redirect(routes.App.Index())
}

func (c User) Signout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}

	return c.Redirect(routes.App.Index())
}

func (c User) Edit() revel.Result {
	avatars := models.Avatars

	return c.Render(avatars)
}

func (c User) EditPost(avatar string) revel.Result {
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.User.Edit())
	}

	var user models.User
	has, _ := engine.Id(c.user().Id).Get(&user)
	if !has {
		return c.NotFound("用户不存在")
	}

	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()
		if ok := checkImageExt(c.Validation, &file, header, IMAGE_EXTS); ok {
			fileName := uuidFileName(header.Filename)
			err, ret := qiniuUploadImage(&file, fileName)
			if err != nil {
				c.Flash.Error("上传头像到七牛出错，请检查七牛配置。")
				return c.Redirect(routes.User.Edit())
			} else {
				if user.IsCustomAvatar() {
					qiniuDeleteImage(user.Avatar)
				}
				user.Avatar = ret.Key
			}
		}
	} else if avatar != "" {
		if user.IsCustomAvatar() {
			qiniuDeleteImage(user.Avatar)
		}
		user.Avatar = avatar
	}

	aff, _ := engine.Id(c.user().Id).Cols("avatar").Update(&user)
	if aff > 0 {
		c.Flash.Success("保存成功")
	} else {
		c.Flash.Error("保存失败")
	}

	return c.Redirect(routes.User.Edit())
}

func (c User) Validate(code string) revel.Result {
	var user models.User
	has, _ := engine.Where("validate_code = ?", code).Get(&user)
	if !has {
		return c.NotFound("用户不存在或校验码错误")
	}

	user.Status = models.USER_STATUS_ACTIVATED
	user.ValidateCode = ""
	aff, _ := engine.Cols("status", "validate_code").Update(&user)
	if aff > 0 {
		c.Flash.Success("您的账号成功激活，请登录！")
	} else {
		c.Flash.Error("抱歉，您的账号未能激活，请与管理员联系！")
	}

	return c.Redirect(routes.User.Signin())
}

func (c User) ForgotPassword() revel.Result {
	return c.Render()
}

func (c User) ForgotPasswordPost(email string) revel.Result {
	c.Validation.Required(email).Message("请填写Email")
	c.Validation.Email(email).Message("Email格式不正确")
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.User.ForgotPassword())
	}

	var user models.User
	has, _ := engine.Where("email = ?", email).Get(&user)
	if !has {
		return c.NotFound("用户不存在")
	}

	user.ValidateCode = uuidName()
	engine.Where("email = ?", email).Cols("validate_code").Update(&user)

	subject := "重设密码 —— Revel社区"
	content := `<h2><a href="http://gorevel.cn/reset_password/` + user.ValidateCode + `">重设密码</a></h2>`

	go sendMail(subject, content, []string{user.Email})

	c.Flash.Success(fmt.Sprintf("你好，重设密码的链接已经发送到您的邮箱，请到您的邮箱 %s 重设密码！", user.Email))

	return c.Redirect(routes.User.Signin())
}

func (c User) ResetPassword(code string) revel.Result {
	return c.Render(code)
}

func (c User) ResetPasswordPost(code, password, confirmPassword string) revel.Result {
	var user models.User
	has, _ := engine.Where("validate_code = ?", code).Get(&user)
	if code == "" || !has {
		return c.NotFound("用户不存在或验证码错误")
	}

	c.Validation.Required(password).Message("请填写新密码")
	c.Validation.Required(confirmPassword == password).Message("新密码不一致")
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.User.ResetPassword(code))
	}

	salt := uuidName()

	aff, _ := engine.Id(user.Id).Update(&models.User{
		HashedPassword: models.EncryptPassword(password, salt),
		Salt:           salt,
		ValidateCode:   "",
	})

	if aff > 0 {
		c.Flash.Success(fmt.Sprintf("%s，你好！重设密码成功，请登录！", user.Name))
	} else {
		c.Flash.Error("出现未知错误，请与管理员联系！")
	}

	return c.Redirect(routes.User.Signin())
}
