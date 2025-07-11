package view

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mangaweb4/mangaweb4-backend/configuration"
	"github.com/mangaweb4/mangaweb4-backend/database"
	"github.com/mangaweb4/mangaweb4-backend/ent"
	ent_meta "github.com/mangaweb4/mangaweb4-backend/ent/meta"
	"github.com/mangaweb4/mangaweb4-backend/ent/progress"
	"github.com/mangaweb4/mangaweb4-backend/handler"
	"github.com/mangaweb4/mangaweb4-backend/meta"
	"github.com/mangaweb4/mangaweb4-backend/user"
	"github.com/rs/zerolog/log"
)

type viewRequest struct {
	User string `json:"user"`
	Name string `json:"name"`
}

type viewResponse struct {
	Request     viewRequest `json:"request"`
	Name        string      `json:"name"`
	Version     string      `json:"version"`
	Favorite    bool        `json:"favorite"`
	PageCount   int         `json:"page_count"`
	CurrentPage int         `json:"current_page"`
	Tags        []*ent.Tag  `json:"tags"`
}

const (
	PathView = "/view"
)

// @accept json
// @Param request body view.viewRequest true "request"
// @Success      200  {object}  view.viewResponse
// @Failure      500  {object}  errors.Error
// @Router /view [post]
func Handler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	req := viewRequest{}

	if err := handler.ParseInput(r.Body, &req); err != nil {
		handler.WriteResponse(w, err)
		return
	}

	item := req.Name

	client := database.CreateEntClient()
	defer client.Close()

	m, err := meta.Read(r.Context(), client, item)
	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	if !m.Read {
		m.Read = true
		meta.Write(r.Context(), client, m)
	}

	log.Info().
		Interface("request", req).
		Msg("View Item")

	tags, err := m.QueryTags().All(r.Context())
	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	u, err := user.GetUser(r.Context(), client, req.User)
	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	progress, _ := client.Progress.Query().Where(progress.UserID(u.ID), progress.ItemID(m.ID)).Only(r.Context())

	currentPage := 0
	if progress != nil {
		currentPage = progress.Page
	}

	c := configuration.Get()
	data := viewResponse{
		Request:     req,
		Name:        item,
		Version:     c.VersionString,
		Favorite:    u.QueryFavoriteItems().Where(ent_meta.ID(m.ID)).ExistX(r.Context()),
		Tags:        tags,
		PageCount:   len(m.FileIndices),
		CurrentPage: currentPage,
	}

	client.History.Create().
		SetUser(u).
		SetItem(m).
		Save(r.Context())

	handler.WriteResponse(w, data)
}
