package server

import (
	"context"

	"github.com/mangaweb4/mangaweb4-backend/database"
	"github.com/mangaweb4/mangaweb4-backend/ent/meta"
	"github.com/mangaweb4/mangaweb4-backend/ent/tag"
	"github.com/mangaweb4/mangaweb4-backend/grpc"
	"github.com/mangaweb4/mangaweb4-backend/system"

	"github.com/rs/zerolog/log"
)

type SystemServer struct {
	grpc.UnimplementedSystemServer
}

func (s *SystemServer) Info(
	ctx context.Context,
	req *grpc.SystemInfoRequest,
) (resp *grpc.SystemInfoResponse, err error) {
	defer func() { log.Err(err).Msg("SystemServer.Info") }()
	client := database.CreateEntClient()

	defer func() { log.Err(client.Close()).Msg("database client close on MangaServer.List") }()

	countManga, err := client.Meta.Query().Where(meta.Active(true), meta.Hidden(false)).Count(ctx)
	if err != nil {
		return
	}

	countTag, err := client.Tag.Query().Where(tag.Hidden(false)).Count(ctx)
	if err != nil {
		return
	}

	resp = &grpc.SystemInfoResponse{
		Version:   system.VersionString,
		ItemCount: int32(countManga),
		TagCount:  int32(countTag),
	}

	return
}
