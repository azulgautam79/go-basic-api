package scalar

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	html := `
<!DOCTYPE html>
<html>
<head>
	<title>Go Task API</title>
	<meta charset="utf-8" />

	<link
		rel="icon"
		type="image/svg+xml"
		href="/static/lemon_purple.svg"
	/>
</head>
<body>
	<script
		id="api-reference"
		data-url="/openapi.json"
		data-configuration='{ "theme": "kepler", "layout": "modern", "showSidebar": true, "hideModels": false, "defaultOpenFirstTag": true }'>
	</script>

	<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>
`

	_, _ = w.Write([]byte(html))
}
