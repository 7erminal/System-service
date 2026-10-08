package requests

type RolePermissionRequest struct {
	Role           int64
	PermissionCode string
	Action         string
	AddedBy        string
}

type RolesRequest struct {
	Role        string
	Description string
	AddedBy     string
}

type PermissionRequest struct {
	Permission  string
	Description string
	AddedBy     string
}
