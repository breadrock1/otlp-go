package otlp_go

var (
	ExcludedPaths = []string{
		"/health",
		"/metrics",
		"/favicon.ico",
		"/static/",
		"/api/metrics",
		"/api/swagger",
	}
)
