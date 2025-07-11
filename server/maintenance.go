package server

import (
	"context"

	"github.com/mangaweb4/mangaweb4-backend/grpc"
	"github.com/mangaweb4/mangaweb4-backend/maintenance"
	"github.com/rs/zerolog/log"
)

type MaintenanceServer struct {
	grpc.UnimplementedMaintenanceServer
}

func (s *MaintenanceServer) PurgeCache(ctx context.Context, req *grpc.MaintenancePurgeCacheRequest) (resp *grpc.MaintenancePurgeCacheResponse, err error) {
	log.Info().Msg("Purge Cache")

	go maintenance.PurgeCache()

	resp = &grpc.MaintenancePurgeCacheResponse{
		IsSuccess: true,
	}
	err = nil

	return

}

func (s *MaintenanceServer) UpdateLibrary(ctx context.Context, req *grpc.MaintenanceUpdateLibraryRequest) (resp *grpc.MaintenanceUpdateLibraryResponse, err error) {
	log.Info().Msg("Update library")

	go maintenance.UpdateLibrary(context.Background())

	resp = &grpc.MaintenanceUpdateLibraryResponse{
		IsSuccess: true,
	}

	err = nil
	return
}
