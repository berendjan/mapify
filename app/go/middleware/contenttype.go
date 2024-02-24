package middleware

import (
	"net/http"
	"strings"
)

var contentTypeMap = map[string]string{
	"css": "text/css",
	"js":  "text/javascript",
}

func ContentTypeHeader(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path
	_, fileType, found := strings.Cut(filePath, ".")
	if found {
		if contentType, ok := contentTypeMap[fileType]; ok {
			w.Header().Add("Content-Type", contentType)
		}
	}
}
