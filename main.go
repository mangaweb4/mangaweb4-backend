package main

import (
	"context"
	"flag"
	"net"
	"os"
	"strconv"

	"entgo.io/ent/dialect"
	"github.com/joho/godotenv"
	"github.com/mangaweb4/mangaweb4-backend/configuration"
	"github.com/mangaweb4/mangaweb4-backend/database"
	m4_grpc "github.com/mangaweb4/mangaweb4-backend/grpc"
	"github.com/mangaweb4/mangaweb4-backend/maintenance"
	"github.com/mangaweb4/mangaweb4-backend/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

var versionString string = "development"

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

	log.Info().Msg("Server starts.")

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Error().AnErr("error", err).Str("address", address).Msg("failed to listen")
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	m4_grpc.RegisterHistoryServer(grpcServer, &server.HistoryServer{})
	m4_grpc.RegisterMaintenanceServer(grpcServer, &server.MaintenanceServer{})
	m4_grpc.RegisterMangaServer(grpcServer, &server.MangaServer{})
	m4_grpc.RegisterTagServer(grpcServer, &server.TagServer{})

	if err := grpcServer.Serve(listener); err != nil {
		log.Error().AnErr("error", err).Msg("Starting server fails")
		return
	}

	log.Info().Msg("Server shutdown.")
}
