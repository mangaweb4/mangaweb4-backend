package maintenance

import (
	"context"

	"github.com/mangaweb4/mangaweb4-backend/database"
	"github.com/rs/zerolog/log"
)

func UpdateLibrary(ctx context.Context) {
	client := database.CreateEntClient()
	defer client.Close()

	log.Info().Msg("Update metadata set.")
	ScanLibrary(ctx, client)
}
