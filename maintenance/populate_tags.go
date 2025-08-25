package maintenance

import (
	"context"

	"github.com/mangaweb4/mangaweb4-backend/database"
	"github.com/mangaweb4/mangaweb4-backend/meta"
	"github.com/rs/zerolog/log"
)

func PopulateTags(ctx context.Context) error {

	client := database.CreateEntClient()
	defer func() { log.Err(client.Close()).Msg("Populate tags close client.") }()

	allMeta, err := meta.ReadAll(context.Background(), client)
	if err != nil {
		return err
	}

	for _, m := range allMeta {
		log.Info().Str("item", m.Name).Msg("Populate tags.")
		_, _, err := meta.PopulateTags(context.Background(), client, m)
		if err != nil {
			log.Err(err).Msg("fails to populate tags.")
		}
	}
	return nil
}
