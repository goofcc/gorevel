package models

import (
	"crypto/md5"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/robfig/revel"
)

const (
	USER_STATUS_INACTIVE = iota
	USER_STATUS_ACTIVATED
	USER_STATUS_DISABLE
)

var (
	Avatars = []string{
		"gopher_teal.jpg",
		"gopher_aqua.jpg",
		"gopher_brown.jpg",
		"gopher_strawberry_bg.jpg",
		"gopher_strawberry.jpg",
	}
	DefaultAvatar = Avatars[0]
)

type User struct {
	Id              int64
	Name            string
	Email           string
	HashedPassword  string
	Salt            string
	Type            int // 1管理员，2普通用户
	Avatar          string
	ValidateCode    string
	Status          int
	Created         time.Time   `xorm:"created"`
	Updated         time.Time   `xorm:"updated"`
	Password        string      `xorm:"-"`
	ConfirmPassword string      `xorm:"-"`
	Permissions     map[int]int `xorm:"-"`
}

var (
	nameRegex = regexp.MustCompile("^\\w*$")
)

func (user User) Validate(v *revel.Validation) {
	v.Required(user.Name).Message("请输入用户名")
	v.Match(user.Name, nameRegex).Message("只能使用字母、数字和下划线")

	if user.HasName() {
		err := &revel.ValidationError{
			Message: "用户名已经注册过",
			Key:     "user.Name",
		}
		v.Errors = append(v.Errors, err)
	}

	v.Required(user.Email).Message("请输入Email")
	v.Email(user.Email).Message("无效的电子邮件")

	if user.HasEmail() {
		err := &revel.ValidationError{
			Message: "邮件已经注册过",
			Key:     "user.Email",
		}
		v.Errors = append(v.Errors, err)
	}

	v.Required(user.Password).Message("请输入密码")
	v.MinSize(user.Password, 3).Message("密码最少三位")
	v.Required(user.ConfirmPassword == user.Password).Message("密码不一致")
}

func (u User) HasName() bool {
	var user User
	has, _ := Engine.Where("name = ?", u.Name).Get(&user)

	return has
}

func (u User) HasEmail() bool {
	var user User
	has, _ := Engine.Where("email = ?", u.Email).Get(&user)

	return has
}

// 加密密码,转成md5
func EncryptPassword(password, salt string) string {
	h := md5.New()
	io.WriteString(h, password+salt)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func (u User) GetPermissions() map[int]int {
	if u.Permissions == nil {
		u.Permissions = make(map[int]int)
		var permissions []Permissions
		Engine.Where("user_id = ?", u.Id).Find(&permissions)

		for _, perm := range permissions {
			u.Permissions[perm.Perm] = perm.Perm
		}
	}

	return u.Permissions
}

func (u User) IsAdmin() bool {
	return u.Type == 1
}

func (u User) IsActive() bool {
	return u.Status == USER_STATUS_ACTIVATED
}

// 是否是默认头像
func (u User) IsDefaultAvatar(avatar string) bool {
	return avatar == u.Avatar
}

// 头像的图片地址
func (u User) AvatarImgSrc() string {
	if u.IsCustomAvatar() {
		return fmt.Sprintf("/%s?imageMogr/v2/thumbnail/48x48!", u.Avatar)
	}

	return fmt.Sprintf("/public/img/%s?imageMogr/v2/thumbnail/48x48!", u.Avatar)
}

func (u User) IsCustomAvatar() bool {
	for _, avatar := range Avatars {
		if avatar == u.Avatar {
			return false
		}
	}

	return true
}
