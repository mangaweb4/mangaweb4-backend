package maintenance

import (
	"context"

	"github.com/mangaweb4/mangaweb4-backend/database"
	"github.com/rs/zerolog/log"
)

func UpdateLibrary(ctx context.Context) {
	client := database.CreateEntClient()
	defer func() { log.Err(client.Close()).Msg("Update metadata close client.") }()

	log.Err(ScanLibrary(ctx, client)).Msg("Update metadata set.")
}
