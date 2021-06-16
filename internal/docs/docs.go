package docs

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

type DocsServer struct {
	http.Server
}

func Init(addr string) {

	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	redocHandler := middleware.Redoc(opts, nil)

	http.Handle(addr, redocHandler)
	http.HandleFunc("/swagger.yaml", swaggerHandler)
}

// Serve swagger documentation.
func swaggerHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "swagger.yaml")
}
