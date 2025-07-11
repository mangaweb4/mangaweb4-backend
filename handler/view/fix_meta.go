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

type fixMetaRequest struct {
	Name string `json:"name"`
}

type fixMetaResponse struct {
	Request fixMetaRequest `json:"request"`
	Success bool           `json:"success"`
}

const (
	PathFixMeta = "/view/fix_meta"
)

// @accept json
// @Param request body view.fixMetaRequest true "request"
// @Success      200  {object}  view.fixMetaResponse
// @Failure      500  {object}  errors.Error
// @Router /view/fix_meta [post]
// Fix the input item metadata.
func FixMeta(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	req := fixMetaRequest{}
	if err := handler.ParseInput(r.Body, &req); err != nil {
		handler.WriteResponse(w, err)
		return
	}

	log.Info().
		Interface("request", req).
		Msg("Fix metadata")

	client := database.CreateEntClient()
	defer client.Close()

	m, err := meta.Read(r.Context(), client, req.Name)
	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	if m, _, err = meta.PopulateTags(r.Context(), client, m); err != nil {
		handler.WriteResponse(w, err)
		return
	}

	if err = meta.GenerateImageIndices(m); err != nil {
		handler.WriteResponse(w, err)
		return
	}

	if err = meta.DeleteThumbnail(m); err != nil {
		handler.WriteResponse(w, err)
		return
	}

	if err = meta.Write(r.Context(), client, m); err != nil {
		handler.WriteResponse(w, err)
		return
	}

	handler.WriteResponse(w, fixMetaResponse{
		Request: req,
		Success: true,
	})
}
