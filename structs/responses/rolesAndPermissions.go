package responses

import "system_service/models"

type RolePermissionResponseDTO struct {
	StatusCode     int
	RolePermission *models.Role_permissions
	StatusDesc     string
}

type RolePermissionsResponseDTO struct {
	StatusCode      int
	RolePermissions *[]models.Role_permissions
	StatusDesc      string
}

type RolePermissionsAllResponseDTO struct {
	StatusCode      int
	RolePermissions *[]interface{}
	StatusDesc      string
}

type RoleResponseDTO struct {
	StatusCode int
	Role       *models.Roles
	StatusDesc string
}

type RolesResponseDTO struct {
	StatusCode int
	Roles      *[]models.Roles
	StatusDesc string
}

type RolesAllResponseDTO struct {
	StatusCode int
	Roles      *[]interface{}
	StatusDesc string
}

type PermissionResponseDTO struct {
	StatusCode int
	Permission *models.Permissions
	StatusDesc string
}

type PermissionsResponseDTO struct {
	StatusCode  int
	Permissions *[]models.Permissions
	StatusDesc  string
}

type PermissionsAllResponseDTO struct {
	StatusCode  int
	Permissions *[]interface{}
	StatusDesc  string
}
