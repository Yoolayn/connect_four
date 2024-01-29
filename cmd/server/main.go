package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const mongoUri = "mongodb://localhost:27017"

var (
	games = make(map[uuid.UUID]Game, 1)
	creds = options.Credential{
		Username: "root",
		Password: "example",
	}
	client      *mongo.Client
	db          *mongo.Database
	collections = make(map[string]collection)
	logger *log.Logger
)

func corsConfig() cors.Config {
	config := cors.DefaultConfig()
	config.AllowMethods = []string{"PUT", "POST", "GET", "DELETE", "OPTIONS"}
	config.AllowOrigins = []string{"*"}

	return config
}

func addAdmin() {
	hash, err := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	_, err = collections["users"].c.InsertOne(ctx, User{Login: "login", Password: string(hash), IsAdmin: true}.ToDB())
	if err != nil {
		panic(err)
	}
}

func dbConnection() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoUri).SetAuth(creds))
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Fatal("dbSetup", "failed to connect", err)
	}

	db = client.Database("connect_four")
	collections["users"] = collection{c: db.Collection("users"), name: "users"}
}

func newStyle() (style *log.Styles) {
	style = log.DefaultStyles()
	pinkText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffc0cb"))

	grayText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#808080"))

	style.Key = pinkText
	style.Value = grayText
	return
}

const (
	LevelsDebug   = "debug"
	LevelsInfo    = "info"
	LevelsWarning = "warn"
	LevelsError   = "error"
	LevelsFatal   = "fatal"
)

func setLevel() {
	switch level := os.Getenv("LOG"); level {
	case LevelsDebug:
		logger.SetLevel(log.DebugLevel)
	case LevelsInfo:
		logger.SetLevel(log.InfoLevel)
	case LevelsWarning:
		logger.SetLevel(log.WarnLevel)
	case LevelsError:
		logger.SetLevel(log.ErrorLevel)
	case LevelsFatal:
		logger.SetLevel(log.FatalLevel)
	default:
		logger.SetLevel(log.InfoLevel)
	}
}

func newLogger() *os.File {
	val := os.Getenv("LOGLOC")
	clean := filepath.Clean(val)

	file, err := os.Create(clean)
	if err != nil {
		logger = log.New(os.Stdout)
		log.Error(err)
		return nil
	}
	logger = log.New(file)
	return file
}

func main() {
	file := newLogger()
	if file != nil {
		gin.DisableConsoleColor()
		gin.DefaultWriter = io.MultiWriter(file)
		defer file.Close()
	}
	logger.SetStyles(newStyle())
	setLevel()
	logger.Info("starting")

	r := gin.Default()
	r.Use(cors.New(corsConfig()))
	addHandlers(r)

	dbConnection()
	addAdmin()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
		defer cancel()

		for _, v := range collections {
			err := v.c.Drop(ctx)
			if err != nil {
				panic(err)
			}
		}
	}()

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

	err := srv.Shutdown(ctx)
	if err != nil {
		fmt.Println("error, Shutdown:", err)
	}
}
