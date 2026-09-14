package scalar

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	html := `
<!DOCTYPE html>
<html>
<head>
	<title>Go SQLite Task API</title>
	<meta charset="utf-8" />
</head>
<body>
	<script
		id="api-reference"
		data-url="/openapi.json">
	</script>

	<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>
`

	_, _ = w.Write([]byte(html))
}
