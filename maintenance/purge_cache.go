package maintenance

import (
	"os"

	"github.com/mangaweb4/mangaweb4-backend/configuration"
	"github.com/rs/zerolog/log"
)

func PurgeCache() {
	c := configuration.Get()

	err := os.RemoveAll(c.CachePath)
	log.Err(err).Msg("Purge cache")
}
