package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	game "gitlab.com/Yoolayn/connect_four/internal/logic"
	"golang.org/x/crypto/bcrypt"
)

var (
	games = make(map[uuid.UUID]game.Board, 1)
	users Users
)

func corsConfig() cors.Config {
	config := cors.DefaultConfig()
	config.AllowMethods = []string{"PUT", "POST", "GET", "DELETE", "OPTIONS"}
	config.AllowOrigins = []string{"*"}

	return config
}

func main() {
	r := gin.Default()
	r.Use(cors.New(corsConfig()))
	addHandlers(r)

	hash, err := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	users = append(users, User{Login: "login", Password: string(hash)})

	srv := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Println("error, ListenAndServe:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		fmt.Println("error, Shutdown:", err)
	}
}
