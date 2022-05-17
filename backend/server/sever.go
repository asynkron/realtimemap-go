package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/asynkron/protoactor-go/cluster"
	echo "github.com/labstack/echo/v4"
	"github.com/rs/cors"
)

type Server struct {
	Hub     *AppHub
	handler http.Handler
	ctx     context.Context
}

func NewHttpServer(cluster *cluster.Cluster, ctx context.Context) *Server {
	router := http.NewServeMux()
	hub := serveHub(router, cluster.ActorSystem, ctx)

	echo := echo.New()
	router.Handle("/", echo)

	serveApi(echo, cluster)

	handler := cors.New(cors.Options{

		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(router)

	return &Server{
		Hub:     hub,
		handler: handler,
		ctx:     ctx,
	}
}

func (s *Server) ListenAndServe() <-chan bool {
	done := make(chan bool)

	go listenAndServe(s.handler, done, s.ctx)

	return done
}

func listenAndServe(handler http.Handler, done chan<- bool, ctx context.Context) {
	address := "localhost:5000"
	server := &http.Server{Addr: address, Handler: handler}

	go func() {
		fmt.Printf("Http server starting to listen at http://%s\n", address)
		if err := server.ListenAndServe(); err != nil {
			fmt.Println(err)
		}

		done <- true
	}()

	<-ctx.Done()
	fmt.Println("Shutting down http server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Shutdown(shutdownCtx)
}
