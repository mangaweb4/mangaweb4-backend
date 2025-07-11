package view

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/mangaweb4/mangaweb4-backend/container"
	"github.com/mangaweb4/mangaweb4-backend/database"
	"github.com/mangaweb4/mangaweb4-backend/handler"
	"github.com/mangaweb4/mangaweb4-backend/meta"
	"github.com/rs/zerolog/log"
)

const (
	PathDownload = "/view/download"
)

// @Param name query string true "name of the file"
// @Success      200  {body}  file
// @Failure      500  {object}  errors.Error
// @Router /view/download [get]
func Download(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	item := r.URL.Query().Get("name")

	log.Info().Str("name", item).Msg("Download")

	client := database.CreateEntClient()
	defer client.Close()

	m, err := meta.Read(r.Context(), client, item)
	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	c, err := container.CreateContainer(m)
	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	reader, filename, err := c.Download(r.Context())
	if err != nil {
		handler.WriteResponse(w, err)
		return
	}
	defer reader.Close()

	bytes, err := io.ReadAll(reader)
	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Length", strconv.Itoa(len(bytes)))
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.WriteHeader(http.StatusOK)

	w.Write(bytes)
}
