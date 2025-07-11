package maintenance

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mangaweb4/mangaweb4-backend/handler"
	"github.com/mangaweb4/mangaweb4-backend/maintenance"
	"github.com/rs/zerolog/log"
)

const (
	PathPurgeCache = "/maintenance/purge_cache"
)

type PurgeCacheResponse struct {
	Result bool `json:"result"`
}

// @Success      200  {object}  maintenance.PurgeCacheResponse
// @Failure      500  {object}  errors.Error
// @Router /maintenance/purge_cache [get]
func PurgeCacheHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Info().Msg("Purge Cache")

	go maintenance.PurgeCache()

	response := PurgeCacheResponse{
		Result: true,
	}

	handler.WriteResponse(w, response)
}
