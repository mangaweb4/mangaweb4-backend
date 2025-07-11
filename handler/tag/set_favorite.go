package tag

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mangaweb4/mangaweb4-backend/database"
	"github.com/mangaweb4/mangaweb4-backend/handler"
	"github.com/mangaweb4/mangaweb4-backend/tag"
	"github.com/mangaweb4/mangaweb4-backend/user"
	"github.com/rs/zerolog/log"
)

type setFavoriteRequest struct {
	User     string `json:"user"`
	Tag      string `json:"tag"`
	Favorite bool   `json:"favorite"`
}

type setFavoriteResponse struct {
	Request  setFavoriteRequest `json:"request"`
	Favorite bool               `json:"favorite"`
}

const (
	PathSetFavorite = "/tag/set_favorite"
)

// @accept json
// @Param request body tag.setFavoriteRequest true "request"
// @Success      200  {object}  tag.setFavoriteResponse
// @Failure      500  {object}  errors.Error
// @Router /tag/set_favorite [post]
func SetFavoriteHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	req := setFavoriteRequest{}
	if err := handler.ParseInput(r.Body, &req); err != nil {
		handler.WriteResponse(w, err)
	}

	log.Info().Interface("request", req).Msg("Set favorite tag.")

	client := database.CreateEntClient()
	defer client.Close()

	m, err := tag.Read(r.Context(), client, req.Tag)
	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	u, err := user.GetUser(r.Context(), client, req.User)
	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	if req.Favorite {
		_, err = u.Update().AddFavoriteTags(m).Save(r.Context())
	} else {
		_, err = u.Update().RemoveFavoriteTags(m).Save(r.Context())
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
