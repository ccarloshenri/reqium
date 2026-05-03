package enums

type HTTPMethod string

const (
	MethodGET    HTTPMethod = "GET"
	MethodPOST   HTTPMethod = "POST"
	MethodPUT    HTTPMethod = "PUT"
	MethodPATCH  HTTPMethod = "PATCH"
	MethodDELETE HTTPMethod = "DELETE"
)

func (m HTTPMethod) String() string {
	return string(m)
}

func ValidHTTPMethod(method string) bool {
	switch HTTPMethod(method) {
	case MethodGET, MethodPOST, MethodPUT, MethodPATCH, MethodDELETE:
		return true
	default:
		return false
	}
}

func MethodAllowsBody(method string) bool {
	switch HTTPMethod(method) {
	case MethodPOST, MethodPUT, MethodPATCH:
		return true
	default:
		return false
	}
}
