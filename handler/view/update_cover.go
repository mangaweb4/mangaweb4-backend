package view

import (
	_ "image/png"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mangaweb4/mangaweb4-backend/database"
	"github.com/mangaweb4/mangaweb4-backend/handler"
	"github.com/mangaweb4/mangaweb4-backend/meta"
	"github.com/rs/zerolog/log"
)

type updateCoverRequest struct {
	Name    string           `json:"name"`
	Index   int              `json:"index"`
	Details meta.CropDetails `json:"crop_details"`
}

type updateCoverResponse struct {
	Request updateCoverRequest `json:"request"`
	Success bool               `json:"success"`
}

const (
	PathUpdateCover = "/view/update_cover"
)

// @accept json
// @Param request body view.updateCoverRequest true "request"
// @Success      200  {object}  view.updateCoverResponse
// @Failure      500  {object}  errors.Error
// @Router /view/update_cover [post]
// UpdateCover a handler to update the cover to specific image
func UpdateCover(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	req := updateCoverRequest{}
	if err := handler.ParseInput(r.Body, &req); err != nil {
		handler.WriteResponse(w, err)
		return
	}

	log.Info().
		Interface("request", req).
		Msg("Update cover.")

	client := database.CreateEntClient()
	defer client.Close()

	m, err := meta.Read(r.Context(), client, req.Name)
	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	m.ThumbnailIndex = req.Index
	m.ThumbnailHeight = req.Details.Height
	m.ThumbnailWidth = req.Details.Width
	m.ThumbnailX = req.Details.X
	m.ThumbnailY = req.Details.Y

	err = meta.DeleteThumbnail(m)
	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	err = meta.Write(r.Context(), client, m)
	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	handler.WriteResponse(w, updateCoverResponse{
		Request: req,
		Success: true,
	})
}
