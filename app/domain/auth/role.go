package auth

type Role string

const (
	User       Role = "user"
	Backoffice Role = "backoffice"
	System     Role = "system"
)

func (r Role) IsUser() bool {
	return r == User
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
