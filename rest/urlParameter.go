package rest

type UrlParameter interface {
	GetUrlParameter() string
}

// Helper function for syntax of url parameters, useful for implementers of getUrlParameter
func BuildUrlParameter(key string, value string) string {
	return key + "=" + value
}
