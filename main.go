package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func main() {

	// Get cli args
	args := os.Args[1:]
	port := "8080"

	if len(args) > 0 {
		if port != "" {
			port = args[0]
		}
	}

	listenAddr := fmt.Sprintf(":%s", port)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "application/json"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Public Routes
	r.Group(func(r chi.Router) {
		r.Get("/", index)
		r.Get("/healthz", healthz)
		r.Get("/ws", ws)
	})

	log.Fatal(http.ListenAndServe(listenAddr, r))
}

func index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	io.WriteString(w, `{"msg": go-hello}`)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	io.WriteString(w, `{"alive": true}`)
}

func ws(w http.ResponseWriter, r *http.Request) {
	// Create ws connexion
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		log.Printf("%s", err)
		return
	}

	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*30)
	defer cancel()

	var v interface{}
	err = wsjson.Read(ctx, c, &v)
	if err != nil {
		log.Printf("%s", err)
		return
	}

	log.Printf("received: %v", v)
	err = wsjson.Write(ctx, c, v)
	if err != nil {
		log.Printf("%s", err)
		return
	}

	c.Close(websocket.StatusNormalClosure, "")
}
