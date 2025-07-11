package maintenance

import (
	"os"

	"github.com/mangaweb4/mangaweb4-backend/configuration"
)

func PurgeCache() error {
	c := configuration.Get()

	return os.RemoveAll(c.CachePath)
}
