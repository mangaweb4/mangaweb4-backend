package maintenance

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mangaweb4/mangaweb4-backend/handler"
	"github.com/mangaweb4/mangaweb4-backend/maintenance"
	"github.com/rs/zerolog/log"
)

type UpdateLibraryResponse struct {
	Result bool `json:"result"`
}

const (
	PathUpdateLibrary = "/maintenance/update_library"
)

// @Success      200  {object}  maintenance.UpdateLibraryResponse
// @Failure      500  {object}  errors.Error
// @Router /maintenance/update_library [get]
func UpdateLibraryHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Info().Msg("Update library")

	go maintenance.UpdateLibrary(context.Background())

	response := UpdateLibraryResponse{
		Result: true,
	}

	handler.WriteResponse(w, response)
}
