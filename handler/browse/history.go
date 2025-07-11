package browse

import (
	"net/http"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/julienschmidt/httprouter"
	"github.com/mangaweb4/mangaweb4-backend/database"
	"github.com/mangaweb4/mangaweb4-backend/ent/history"
	ent_meta "github.com/mangaweb4/mangaweb4-backend/ent/meta"
	"github.com/mangaweb4/mangaweb4-backend/handler"
	"github.com/mangaweb4/mangaweb4-backend/user"
	"github.com/rs/zerolog/log"
)

const (
	PathHistory = "/history"
)

type historyRequest struct {
	User        string `json:"user"`
	Page        int    `json:"page"`
	ItemPerPage int    `json:"item_per_page" default:"30"`
}

type historyResponse struct {
	Request   historyRequest `json:"request"`
	TotalPage int            `json:"total_page"`
	Items     []historyItem  `json:"items"`
}

type historyItem struct {
	// ID of the item
	ID int `json:"id,omitempty"`
	// Name of the item
	Name string `json:"name,omitempty"`
	// Favorite this item is a favorite
	Favorite bool `json:"favorite,omitempty"`
	// Read this item has been read before.
	Read bool `json:"read,omitempty"`
	// PageCount the number of pages.
	PageCount int `json:"page_count,omitempty"`
	// AccessTime the time the item is accessed.
	AccessTime time.Time `json:"access_time,omitempty"`
	// TagFavorite this item contains favorite tags
	TagFavorite bool `json:"tag_favorite,omitempty"`
}

func createDefaultHistoryRequest() historyRequest {
	return historyRequest{
		Page:        0,
		ItemPerPage: 30,
	}
}

// @accept json
// @Param request body browse.historyRequest false "request"
// @Success      200  {object}  browse.historyResponse
// @Failure      500  {object}  errors.Error
// @Router /history [post]
func historyHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	req := createDefaultHistoryRequest()

	if err := handler.ParseInput(r.Body, &req); err != nil {
		handler.WriteResponse(w, err)
	}

	client := database.CreateEntClient()
	defer client.Close()

	u, err := user.GetUser(r.Context(), client, req.User)
	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	histories, err := client.User.QueryHistories(u).
		Order(history.ByCreateTime(sql.OrderDesc())).
		Limit(req.ItemPerPage).
		Offset(req.ItemPerPage * req.Page).All(r.Context())

	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	items := make([]historyItem, len(histories))

	for i, h := range histories {
		m, err := h.QueryItem().Only(r.Context())

		if err != nil {
			handler.WriteResponse(w, err)
			return
		}

		items[i] = historyItem{
			ID:         m.ID,
			Name:       m.Name,
			Favorite:   u.QueryFavoriteItems().Where(ent_meta.ID(m.ID)).ExistX(r.Context()),
			Read:       true,
			PageCount:  len(m.FileIndices),
			AccessTime: h.CreateTime,
		}

		tags, err := m.QueryTags().All(r.Context())
		if err != nil {
			handler.WriteResponse(w, err)
			return
		}

		for _, t := range tags {
			if t.Favorite {
				items[i].TagFavorite = true
				break
			}
		}
	}

	count, err := client.History.Query().Count(r.Context())

	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	pageCount := int(count) / req.ItemPerPage
	if int(count)%req.ItemPerPage > 0 {
		pageCount++
	}

	if req.Page > pageCount || req.Page < 0 {
		req.Page = 0
	}

	log.Info().
		Interface("request", req).
		Msg("Browse")

	data := historyResponse{
		Request:   req,
		Items:     items,
		TotalPage: pageCount,
	}

	handler.WriteResponse(w, data)
}
