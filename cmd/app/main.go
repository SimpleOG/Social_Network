package main

import (
	"context"
	"github.com/SimpleOG/Social_Network/internal/api/controllers/server"
	db "github.com/SimpleOG/Social_Network/internal/repositories/database/postgresql/sqlc"
	"github.com/SimpleOG/Social_Network/internal/repositories/redis"
	"github.com/SimpleOG/Social_Network/internal/service"
	"github.com/SimpleOG/Social_Network/pkg/jwt"
	"github.com/SimpleOG/Social_Network/pkg/util/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log"
)

var (
	ReadBufferSize  = 1024
	WriteBufferSize = 1024
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	//gin.SetMode(gin.ReleaseMode)
	config, err := config.InitConfig("config/env")
	if err != nil {
		log.Fatalln(err)
	}
	connPool, err := pgxpool.New(context.Background(), config.DbSource)
	if err != nil {
		log.Fatalf("cannot connect to db %s", err)
	}
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  ReadBufferSize,
		WriteBufferSize: WriteBufferSize,
	}
	AuthInterface := jwt.NewJwtAuth(config.SecretKey)
	client, err := redis.NewRedisClient()
	if err != nil {
		log.Fatalln(err)
	}
	Service := service.NewService(db.New(connPool), AuthInterface, client)
	newServer, err := server.NewServer(router, Service, upgrader)

	if err != nil {
		log.Fatalln(err)
	}
	mirationPath := "file://internal/repositories/database/postgresql/migrations/"
	runDBMigration(mirationPath, config.DbSource)

	err = newServer.Start(":8080")
	if err != nil {
		log.Fatalln(err)
	}
	StopDBMigration(mirationPath, config.DbSource)

}

func runDBMigration(migrationURL, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatalln("cannot find migration to up", err)
	}
	if err = migration.Up(); err != nil {
		if err.Error() == "no change" {
			log.Println("уже заполнено")
			return
		}
		log.Fatalln("cannot start migration", err)
	}

}
func StopDBMigration(migrationURL, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatalln("cannot find migration to down", err)
	}
	if err = migration.Down(); err != nil {
		log.Fatalln("cannot stop migration", err)
	}
	log.Println("stopped")
}
