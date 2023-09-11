package response

import (
	"net/http"

	"github.com/go-chi/render"
)

func OK(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	render.Status(r, status)
	render.JSON(w, r, data)
}

func Error(w http.ResponseWriter, r *http.Request, status int, message string) {
	render.Status(r, status)
	render.JSON(w, r, map[string]string{"error": message})
}
