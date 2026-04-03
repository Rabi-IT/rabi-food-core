package auth

type Role string

const (
	User       Role = "user"
	Staff      Role = "staff"
	Backoffice Role = "backoffice"
	System     Role = "system"
)

func (r Role) IsBackoffice() bool {
	return r == Backoffice
}

func (r Role) IsUser() bool {
	return r == User
}

func (r Role) IsStaff() bool {
	return r == Staff
}

func (r Role) IsSystem() bool {
	return r == System
}

func (r Role) String() string {
	return string(r)
}
