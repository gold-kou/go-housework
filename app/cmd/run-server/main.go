package main

import (
	"net/http"

	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/server"
	"github.com/gold-kou/go-housework/app/server/middleware"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Printf("Server started")

	// DBと接続する
	db, err := common.OpenGDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	common.SetDB(db)

	router := server.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", middleware.CORS(router)))
}
