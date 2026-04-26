package errs

import "net/http"

var ErrTenantNotFound = newErr("TENANT_NOT_FOUND", http.StatusNotFound)
