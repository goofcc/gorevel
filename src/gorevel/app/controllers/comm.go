package controllers

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/robfig/config"
	"github.com/robfig/revel"
	"io"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
)

type SmtpType struct {
	username string
	password string
	host     string
	address  string
	from     string
}

var (
	avatars = []string{
		"gopher_teal.jpg",
		"gopher_aqua.jpg",
		"gopher_brown.jpg",
		"gopher_strawberry_bg.jpg",
		"gopher_strawberry.jpg",
	}
	defaultAvatar = avatars[0]

	basePath   string = ""
	uploadPath string = ""
	imageExts  string = ".jpg.jpeg.png"
	Smtp       SmtpType
)

func saveFile(r *revel.Request, formField string) string {
	file, header, err := r.FormFile(formField)
	if err != nil {
		return ""
	}
	defer file.Close()

	uuid := strings.Replace(uuid.NewUUID().String(), "-", "", -1)
	ext := filepath.Ext(header.Filename)
	fileName := uuid + ext

	os.MkdirAll(uploadPath, 0777)

	f, err := os.OpenFile(uploadPath+fileName, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()

	if err != nil {
		fmt.Println(err)
	} else {
		io.Copy(f, file)
	}

	return fileName
}

func deleteFile(fileName string) error {
	err := os.Remove(uploadPath + fileName)

	if err != nil {
		fmt.Println(err)
	}

	return err
}

func getFileExt(r *revel.Request, formField string) (bool, string) {
	file, header, err := r.FormFile(formField)
	if err != nil {
		return false, ""
	}
	defer file.Close()

	return true, strings.ToLower(filepath.Ext(header.Filename))
}

func checkFileExt(c *revel.Controller, fileExts, formField, message string) {
	if ok, ext := getFileExt(c.Request, formField); ok && !strings.Contains(fileExts, ext) {
		err := &revel.ValidationError{
			Message: message,
			Key:     formField,
		}
		c.Validation.Errors = append(c.Validation.Errors, err)
	}
}

func sendMail(subject, content string, tos []string) error {
	message := `From: Revel中文社区
To: ` + strings.Join(tos, ",") + `
Subject: ` + subject + `
Content-Type: text/html;charset=UTF-8

` + content
	if Smtp.username == "" {
		path, _ := filepath.Abs("")
		c, _ := config.ReadDefault(fmt.Sprintf("%s/src/gorevel/conf/my.conf", path))

		Smtp.username, _ = c.String("smtp", "smtp.username")
		Smtp.password, _ = c.String("smtp", "smtp.password")
		Smtp.address, _ = c.String("smtp", "smtp.address")
		Smtp.from, _ = c.String("smtp", "smtp.from")
		Smtp.host, _ = c.String("smtp", "smtp.host")
	}

	auth := smtp.PlainAuth("", Smtp.username, Smtp.password, Smtp.host)
	err := smtp.SendMail(Smtp.address, auth, Smtp.from, tos, []byte(message))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
