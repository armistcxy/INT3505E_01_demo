package access

func CanAccess(role string, resource string) bool {
	if role == "admin" {
		return true
	}
	if role == "user" && resource != "admin_panel" {
		return true
	}
	return false
}
