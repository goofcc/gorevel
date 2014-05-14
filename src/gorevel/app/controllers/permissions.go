package controllers

const (
	ADMIN_GROUP = iota + 1
	MEMBER_GROUP
)

var Permissions = map[string]int{
	// Admin
	"Admin.Index":          ADMIN_GROUP,
	"Admin.ListUser":       ADMIN_GROUP,
	"Admin.DeleteUser":     ADMIN_GROUP,
	"Admin.NewCategory":    ADMIN_GROUP,
	"Admin.EditCategory":   ADMIN_GROUP,
	"Admin.ListCategory":   ADMIN_GROUP,
	"Admin.DeleteCategory": ADMIN_GROUP,
	"Admin.DeleteProduct":  ADMIN_GROUP,

	// Topic
	"Topic.New":     MEMBER_GROUP,
	"Topic.Edit":    MEMBER_GROUP,
	"Topic.Reply":   MEMBER_GROUP,
	"Topic.SetGood": ADMIN_GROUP,

	// User
	"User.Edit": MEMBER_GROUP,

	// Product
	"Product.New":  MEMBER_GROUP,
	"Product.Edit": MEMBER_GROUP,
}
