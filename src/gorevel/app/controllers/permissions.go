package controllers

const (
	_ = iota
	AdminGroup
	MemberGroup
)

var Permissions = map[string]int{
	// Admin
	"Admin.Index":          AdminGroup,
	"Admin.ListUser":       AdminGroup,
	"Admin.DeleteUser":     AdminGroup,
	"Admin.ListCategory":   AdminGroup,
	"Admin.DeleteCategory": AdminGroup,
	"Admin.NewCategory":    AdminGroup,
	"Admin.EditCategory":   AdminGroup,

	// User
	"User.Edit": MemberGroup,

	// Topic
	"Topic.New":   MemberGroup,
	"Topic.Edit":  MemberGroup,
	"Topic.Reply": MemberGroup,
}
