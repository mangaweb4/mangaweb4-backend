package maintenance

import (
	"context"

	"github.com/mangaweb4/mangaweb4-backend/ent"
	"github.com/mangaweb4/mangaweb4-backend/meta"
	"github.com/rs/zerolog/log"
)

func RebuildThumbnail(client *ent.Client) error {
	allMeta, err := meta.ReadAll(context.Background(), client)
	if err != nil {
		return err
	}

	for _, m := range allMeta {
		err := meta.DeleteThumbnail(m)
		if err != nil {
			log.Warn().Err(err).Str("meta", m.Name).Msg("unable to delete thumbnail file")
		}
	}

	return nil
}
