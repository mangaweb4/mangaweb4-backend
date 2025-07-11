package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"strconv"

	"entgo.io/ent/dialect"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/mangaweb4/mangaweb4-backend/configuration"
	"github.com/mangaweb4/mangaweb4-backend/database"
	_ "github.com/mangaweb4/mangaweb4-backend/docs"
	"github.com/mangaweb4/mangaweb4-backend/handler/browse"
	maintenanceHandler "github.com/mangaweb4/mangaweb4-backend/handler/maintenance"
	"github.com/mangaweb4/mangaweb4-backend/handler/tag"
	"github.com/mangaweb4/mangaweb4-backend/handler/view"
	"github.com/mangaweb4/mangaweb4-backend/maintenance"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

var versionString string = "development"

//go:generate go run -mod=mod github.com/swaggo/swag/cmd/swag@latest init

// @title           Mangaweb3 API
// @version         3.0
// @description     API Server for Mangaweb

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	ctx := context.Background()

	flag.Usage = func() {
		os.Stderr.WriteString("Usage: mangaweb3-backend [options]\n\n")
		flag.PrintDefaults()
		os.Exit(0)
	}

	envFlag := flag.String("environment", "",
		"Choose the environment the server run as.\n"+
			"The {{environment}}.env will be loaded and override the environment variables set on the system.")

	helpFlag := flag.Bool("help", false, "Show this help message.")
	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	envFile := *envFlag + ".env"

	useEnvFile := false
	if err := godotenv.Overload(envFile); err == nil {
		useEnvFile = true
	}

	debugMode := false
	if value, valid := os.LookupEnv("MANGAWEB_DEBUG"); valid {
		debugMode, _ = strconv.ParseBool(value)
		if debugMode {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).
				Level(zerolog.DebugLevel)
		}
	}

	log.Info().Str("environment", *envFlag).Msg("")

	if !useEnvFile {
		log.Info().Str("file", envFile).Msg("Environment file not found.")
	}

	address := ":8972"
	if value, valid := os.LookupEnv("MANGAWEB_ADDRESS"); valid {
		address = value
	}

	dataPath := "./data"
	if value, valid := os.LookupEnv("MANGAWEB_DATA_PATH"); valid {
		dataPath = value
	}

	cachePath := "./cache"
	if value, valid := os.LookupEnv("MANGAWEB_CACHE_PATH"); valid {
		cachePath = value
	}

	connectionStr := "postgres://postgres:password@localhost:5432/manga"
	if value, valid := os.LookupEnv("MANGAWEB_DB"); valid {
		connectionStr = value
	}

	dbType := dialect.Postgres
	if value, valid := os.LookupEnv("MANGAWEB_DB_TYPE"); valid {
		dbType = value
	}

	log.Info().
		Bool("debugMode", debugMode).
		Str("version", versionString).
		Str("dataPath", dataPath).
		Str("cachePath", cachePath).
		Msg("Server initializes.")

	configuration.Init(configuration.Config{
		DebugMode:     debugMode,
		VersionString: versionString,
		DataPath:      dataPath,
		CachePath:     cachePath,
	})

	log.Info().Str("dbType", dbType).Str("dbConnection", connectionStr).Msg("Database open.")
	if err := database.Open(ctx, dbType, connectionStr); err != nil {
		log.Error().AnErr("error", err).Msg("Connect to Database fails")
		return
	} else {
		defer database.Close()
	}

	if err := database.CreateSchema(ctx); err != nil {
		log.Error().AnErr("error", err).Msg("Database creating schema fails.")
		return
	}

	go maintenance.UpdateLibrary(context.Background())

	router := httprouter.New()
	RegisterHandler(router, debugMode)

	log.Info().Msg("Server starts.")

	handler := cors.AllowAll().Handler(router)
	if err := http.ListenAndServe(address, handler); err != nil {
		log.Error().AnErr("error", err).Msg("Starting server fails")
		return
	}

	log.Info().Msg("Server shutdown.")
}

func RegisterHandler(router *httprouter.Router, debugMode bool) {
	browse.Register(router)
	tag.Register(router)
	view.Register(router)
	maintenanceHandler.Register(router)

	if debugMode {
		router.GET("/doc/:any", swaggerHandler)
	}
}

func swaggerHandler(res http.ResponseWriter, req *http.Request, p httprouter.Params) {
	httpSwagger.WrapHandler(res, req)
}
