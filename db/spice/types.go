package spice

type ObjectType string
type RelationType string
type PermissionType string
type CheckPermissionType string

type Relation struct {
	ResourceType ObjectType
	ResourceID   string
	SubjectType  ObjectType
	SubjectID    string
	Relation     RelationType
}

type Permission struct {
	ResourceType ObjectType
	ResourceID   string
	SubjectType  ObjectType
	SubjectID    string
	Permission   PermissionType
}

const (
	HasPermission         CheckPermissionType = "HasPermission"
	UnspecifiedPermission CheckPermissionType = "UnspecifiedPermission"
	NoPermission          CheckPermissionType = "NoPermission"
)
