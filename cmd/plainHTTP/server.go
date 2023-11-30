package plainHTTP

import (
	"go-users-service/cmd/plainHTTP/handlers"
	"go-users-service/internal/core/user"
	"log"
	"net/http"
)

type HTTPServer struct {
}

func NewServer(cockroachRepository user.Repository) *HTTPServer {
	userHandlers := handlers.NewUserHandlers(cockroachRepository)

	http.HandleFunc("/user", userHandlers.Handle)

	return &HTTPServer{}
}

func (srv HTTPServer) ListenAndServe(port string) {
	log.Println("Running HTTP in :", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Println(err)
	}
}
