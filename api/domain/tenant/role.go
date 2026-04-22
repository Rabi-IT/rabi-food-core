package tenant

type Role string

const (
	Owner Role = "tenant_owner"
	Staff Role = "tenant_staff"
)

func (r Role) String() string {
	return string(r)
}
