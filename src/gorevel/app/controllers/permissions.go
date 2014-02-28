package controllers

const (
	_ = iota
	ADMIN_GROUP
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

	// Topic
	"Topic.New":   MEMBER_GROUP,
	"Topic.Edit":  MEMBER_GROUP,
	"Topic.Reply": MEMBER_GROUP,

	// User
	"User.Edit": MEMBER_GROUP,
}
