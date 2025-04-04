package requests

type RolePermissionRequest struct {
	Role           int64
	PermissionCode string
	Action         string
}

type RolesRequest struct {
	Role        string
	Description string
}

type PermissionRequest struct {
	Permission  string
	Description string
}
