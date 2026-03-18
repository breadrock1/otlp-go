package provider

var (
	ExcludedPaths = []string{
		"/health",
		"/metrics",
		"/favicon.ico",
		"/static/",
		"/api/swagger",
	}
)
