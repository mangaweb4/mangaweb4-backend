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

func (s *MaintenanceServer) PurgeCache(
	ctx context.Context,
	req *grpc.MaintenancePurgeCacheRequest,
) (resp *grpc.MaintenancePurgeCacheResponse, err error) {
	defer func() { log.Err(err).Interface("request", req).Msg("MaintenanceServer.PurgeCache") }()

	go maintenance.PurgeCache()

	resp = &grpc.MaintenancePurgeCacheResponse{
		IsSuccess: true,
	}
	err = nil

	return

}

func (s *MaintenanceServer) UpdateLibrary(
	ctx context.Context,
	req *grpc.MaintenanceUpdateLibraryRequest,
) (resp *grpc.MaintenanceUpdateLibraryResponse, err error) {
	defer func() { log.Err(err).Interface("request", req).Msg("MaintenanceServer.UpdateLibrary") }()

	go maintenance.UpdateLibrary(context.Background())

	resp = &grpc.MaintenanceUpdateLibraryResponse{
		IsSuccess: true,
	}

	err = nil
	return
}

func (s *MaintenanceServer) PopulateTags(
	ctx context.Context,
	req *grpc.MaintenancePopulateTagsRequest,
) (resp *grpc.MaintenancePopulateTagsResponse, err error) {
	defer func() { log.Err(err).Interface("request", req).Msg("MaintenanceServer.UpdateLibrary") }()

	go maintenance.PopulateTags(context.Background())

	resp = &grpc.MaintenancePopulateTagsResponse{
		IsSuccess: true,
	}

	err = nil
	return
}
