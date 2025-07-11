package browse

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mangaweb4/mangaweb4-backend/database"
	"github.com/mangaweb4/mangaweb4-backend/handler"
	"github.com/mangaweb4/mangaweb4-backend/maintenance"
	"github.com/rs/zerolog/log"
)

type recreateThumbnailsResponse struct {
	Result bool `json:"result"`
}

const (
	PathRecreateThumbnails = "/browse/recreate_thumbnails"
)

// @Success      200  {object}  browse.recreateThumbnailsResponse
// @Failure      500  {object}  errors.Error
// @Router /browse/recreate_thumbnails [get]
func RecreateThumbnailHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Info().Msg("Rescan library")

	client := database.CreateEntClient()
	defer client.Close()

	go maintenance.RebuildThumbnail(client)

	response := recreateThumbnailsResponse{
		Result: true,
	}

	handler.WriteResponse(w, response)
}
