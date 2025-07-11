package browse

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mangaweb4/mangaweb4-backend/database"
	ent_meta "github.com/mangaweb4/mangaweb4-backend/ent/meta"
	ent_tag "github.com/mangaweb4/mangaweb4-backend/ent/tag"
	"github.com/mangaweb4/mangaweb4-backend/handler"
	"github.com/mangaweb4/mangaweb4-backend/meta"
	"github.com/mangaweb4/mangaweb4-backend/tag"
	"github.com/mangaweb4/mangaweb4-backend/user"
	"github.com/rs/zerolog/log"
)

const (
	PathBrowse = "/browse"
)

type browseRequest struct {
	User        string         `json:"user"`
	Tag         string         `json:"tag"`
	Filter      meta.Filter    `json:"filter"`
	Page        int            `json:"page"`
	ItemPerPage int            `json:"item_per_page" default:"30"`
	Search      string         `json:"search"`
	Sort        meta.SortField `json:"sort"`
	Order       meta.SortOrder `json:"order"`
}

type browseResponse struct {
	Request     browseRequest `json:"request"`
	TagFavorite bool          `json:"tag_favorite"`
	TotalPage   int           `json:"total_page"`
	Items       []browseItem  `json:"items"`
}

type browseItem struct {
	// ID of the item
	ID int `json:"id,omitempty"`
	// Name of the item
	Name string `json:"name,omitempty"`
	// Favorite this item is a favorite
	Favorite bool `json:"favorite,omitempty"`
	// Read this item has been read before.
	Read bool `json:"read,omitempty"`
	// PageCount the number of pages.
	PageCount int `json:"page_count"`
	// TagFavorite this item contains favorite tags
	TagFavorite bool `json:"tag_favorite,omitempty"`
}

func createDefaultBrowseRequest() browseRequest {
	return browseRequest{
		Tag:         "",
		Filter:      "",
		Page:        0,
		Search:      "",
		Sort:        meta.SortFieldCreateTime,
		Order:       meta.SortOrderDescending,
		ItemPerPage: 30,
	}
}

// @accept json
// @Param request body browse.browseRequest false "request"
// @Success      200  {object}  browse.browseResponse
// @Failure      500  {object}  errors.Error
// @Router /browse [post]
func Handler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	req := createDefaultBrowseRequest()

	if err := handler.ParseInput(r.Body, &req); err != nil {
		handler.WriteResponse(w, err)
		return
	}

	client := database.CreateEntClient()
	defer client.Close()

	u, err := user.GetUser(r.Context(), client, req.User)
	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	allMeta, err := meta.ReadPage(r.Context(),
		client,
		u,
		meta.QueryParams{
			SearchName:  req.Search,
			Filter:      req.Filter,
			SearchTag:   req.Tag,
			SortBy:      req.Sort,
			SortOrder:   req.Order,
			Page:        req.Page,
			ItemPerPage: req.ItemPerPage,
		})
	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	items := make([]browseItem, len(allMeta))
	for i, m := range allMeta {
		items[i] = browseItem{
			ID:        m.ID,
			Name:      m.Name,
			Favorite:  u.QueryFavoriteItems().Where(ent_meta.ID(m.ID)).ExistX(r.Context()),
			Read:      u.QueryHistories().QueryItem().Where(ent_meta.ID(m.ID)).ExistX(r.Context()),
			PageCount: len(m.FileIndices),
		}

		tags, err := m.QueryTags().All(r.Context())
		if err != nil {
			handler.WriteResponse(w, err)
			return
		}

		for _, t := range tags {
			if u.QueryFavoriteTags().Where(ent_tag.ID(t.ID)).ExistX(r.Context()) {
				items[i].TagFavorite = true
				break
			}
		}
	}

	count, err := meta.Count(r.Context(),
		client,
		u,
		meta.QueryParams{
			SearchName:  req.Search,
			Filter:      req.Filter,
			SearchTag:   req.Tag,
			SortBy:      req.Sort,
			SortOrder:   req.Order,
			Page:        0,
			ItemPerPage: 0,
		})
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

	data := browseResponse{
		Request:   req,
		Items:     items,
		TotalPage: pageCount,
	}

	if req.Tag != "" {
		tagObj, err := tag.Read(r.Context(), client, req.Tag)
		if err != nil {
			handler.WriteResponse(w, err)
			return
		}

		data.TagFavorite = u.QueryFavoriteTags().Where(ent_tag.ID(tagObj.ID)).ExistX(r.Context())
	}

	handler.WriteResponse(w, data)
}
