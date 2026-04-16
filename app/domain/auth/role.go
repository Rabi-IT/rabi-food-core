package auth

type Role string

const (
	User        Role = "user"
	TenantOwner Role = "tenant_owner"
	Backoffice  Role = "backoffice"
	System      Role = "system"
)

func (r Role) IsUser() bool {
	return r == User
}

func (r Role) IsTenantOwner() bool {
	return r == TenantOwner
}

func (r Role) IsBackoffice() bool {
	return r == Backoffice
}

func (r Role) IsSystem() bool {
	return r == System
}

func (r Role) String() string {
	return string(r)
}
