package maintenance

import (
	"os"

	"github.com/mangaweb4/mangaweb4-backend/configuration"
	"github.com/rs/zerolog/log"
)

func PurgeCache() {
	c := configuration.Get()

	err := os.RemoveAll(c.CachePath)
	if err != nil {
		log.Error().AnErr("error", err).Msg("Purge cache error.")
	}
}
