package tag

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mangaweb4/mangaweb4-backend/database"
	ent_tag "github.com/mangaweb4/mangaweb4-backend/ent/tag"
	"github.com/mangaweb4/mangaweb4-backend/handler"
	"github.com/mangaweb4/mangaweb4-backend/tag"
	"github.com/mangaweb4/mangaweb4-backend/user"
	"github.com/rs/zerolog/log"
)

type listRequest struct {
	User         string        `json:"user"`
	FavoriteOnly bool          `json:"favorite_only"`
	Search       string        `json:"search"`
	Page         int           `json:"page"`
	ItemPerPage  int           `json:"item_per_page"  default:"30"`
	Sort         tag.SortField `json:"sort"`
	Order        tag.SortOrder `json:"order"`
}

type Tag struct {
	Name      string `json:"name,omitempty"`
	Favorite  bool   `json:"favorite,omitempty"`
	ItemCount int    `json:"item_count,omitempty"`
}

type listResponse struct {
	Request   listRequest `json:"request"`
	Tags      []Tag       `json:"tags"`
	TotalPage int         `json:"total_page"`
}

const (
	PathList = "/tag/list"
)

const (
	DefaultItemPerPage = 30
)

// @accept json
// @Param request body tag.listRequest true "request"
// @Success      200  {object}  tag.listResponse
// @Failure      500  {object}  errors.Error
// @Router /tag/list [post]
func ListHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	req := listRequest{
		FavoriteOnly: false,
		Search:       "",
		Page:         0,
		ItemPerPage:  DefaultItemPerPage,
		Sort:         tag.SortFieldName,
		Order:        tag.SortOrderAscending,
	}

	if err := handler.ParseInput(r.Body, &req); err != nil {
		handler.WriteResponse(w, err)
		return
	}

	log.Info().Interface("request", req).Msg("Tag list")

	client := database.CreateEntClient()
	defer client.Close()

	u, err := user.GetUser(r.Context(), client, req.User)
	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	allTags, err := tag.ReadPage(r.Context(), client, u,
		tag.QueryParams{
			FavoriteOnly: req.FavoriteOnly,
			Search:       req.Search,
			Page:         req.Page,
			ItemPerPage:  req.ItemPerPage,
			Sort:         req.Sort,
			Order:        req.Order,
		})

	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	total, err := tag.Count(r.Context(), client, u,
		tag.QueryParams{
			FavoriteOnly: req.FavoriteOnly,
			Search:       req.Search,
			Page:         0,
			ItemPerPage:  0,
			Sort:         req.Sort,
			Order:        req.Order,
		})

	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	data := listResponse{
		Request:   req,
		TotalPage: (total / req.ItemPerPage) + 1,
	}

	data.Tags = make([]Tag, len(allTags))
	for i, t := range allTags {
		items, err := t.QueryMeta().All(r.Context())
		if err != nil {
			handler.WriteResponse(w, err)
			return
		}
		data.Tags[i] = Tag{
			Name:      t.Name,
			Favorite:  u.QueryFavoriteTags().Where(ent_tag.ID(t.ID)).ExistX(r.Context()),
			ItemCount: len(items),
		}
	}

	handler.WriteResponse(w, data)
}
