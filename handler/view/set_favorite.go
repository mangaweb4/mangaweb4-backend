package view

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mangaweb4/mangaweb4-backend/database"
	"github.com/mangaweb4/mangaweb4-backend/handler"
	"github.com/mangaweb4/mangaweb4-backend/meta"
	"github.com/mangaweb4/mangaweb4-backend/user"
	"github.com/rs/zerolog/log"
)

type setFavoriteRequest struct {
	User     string `json:"user"`
	Name     string `json:"name"`
	Favorite bool   `json:"favorite"`
}

type setFavoriteResponse struct {
	Request  setFavoriteRequest `json:"request"`
	Favorite bool               `json:"favorite"`
}

const (
	PathFavorite = "/view/set_favorite"
)

// @accept json
// @Param request body view.setFavoriteRequest true "request"
// @Success      200  {object}  view.setFavoriteResponse
// @Failure      500  {object}  errors.Error
// @Router /view/set_favorite [post]
func SetFavoriteHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	req := setFavoriteRequest{}
	var reqBody []byte

	if read, err := io.ReadAll(r.Body); err != nil {
		handler.WriteResponse(w, err)
		return
	} else {
		reqBody = read
	}

	if err := json.Unmarshal(reqBody, &req); err != nil {
		handler.WriteResponse(w, err)
		return
	}

	item := req.Name

	log.Info().
		Interface("request", req).
		Msg("Set Favorite Item.")

	client := database.CreateEntClient()
	defer client.Close()

	u, err := user.GetUser(r.Context(), client, req.User)
	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	m, err := meta.Read(r.Context(), client, item)
	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	if req.Favorite {
		_, err = u.Update().AddFavoriteItems(m).Save(r.Context())
	} else {
		_, err = u.Update().RemoveFavoriteItems(m).Save(r.Context())
	}

	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	response := setFavoriteResponse{
		Request:  req,
		Favorite: req.Favorite,
	}

	handler.WriteResponse(w, response)
}
