package view

import (
	"bytes"
	"context"
	"io"

	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/julienschmidt/httprouter"
	"github.com/mangaweb4/mangaweb4-backend/container"
	"github.com/mangaweb4/mangaweb4-backend/database"
	"github.com/mangaweb4/mangaweb4-backend/ent/progress"
	"github.com/mangaweb4/mangaweb4-backend/handler"
	"github.com/mangaweb4/mangaweb4-backend/meta"
	"github.com/mangaweb4/mangaweb4-backend/user"
	"github.com/rs/zerolog/log"

	_ "golang.org/x/image/webp"
)

const (
	PathGetImage = "/view/get_image"
)

// @Param name query string true "name of the item"
// @Param width query int false "width"
// @Param height query int false "height"
// @Param i query int true "index"
// @Param user query string false "user"
// @Success      200  {body}  file
// @Failure      500  {object}  errors.Error
// @Router /view/get_image [get]
// GetImage returns an image with specific width/height while retains aspect ratio.
func GetImage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	query := r.URL.Query()
	item := query.Get("name")

	var width, height int64 = 0, 0
	if w, e := strconv.ParseInt(query.Get("width"), 10, 64); e == nil {
		width = w
		height = width
	}

	if h, e := strconv.ParseInt(query.Get("height"), 10, 64); e == nil {
		height = h
	}

	var index = 0
	if i, err := strconv.Atoi(query.Get("i")); err == nil {
		index = i
	}

	log.Info().
		Interface("request", query).
		Msg("Get image")

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

	steam, f, err := c.OpenItem(context.Background(), index)
	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	data, err := io.ReadAll(steam)
	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	userName := query.Get("user")

	u, err := user.GetUser(r.Context(), client, userName)
	if err == nil {
		progressRec, _ := client.Progress.Query().Where(progress.UserID(u.ID), progress.ItemID(m.ID)).Only(r.Context())

		if progressRec == nil {
			_, err = client.Progress.Create().
				SetPage(index).
				SetItem(m).
				SetUser(u).
				Save(r.Context())
		} else {
			_, err = progressRec.Update().
				SetPage(index).
				SetItem(m).
				SetUser(u).
				Save(r.Context())
		}

		if err != nil {
			handler.WriteResponse(w, err)
			return
		}
	}
	if width == 0 && height == 0 {
		var contentType string
		switch filepath.Ext(strings.ToLower(f)) {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".webp":
			contentType = "image/webp"
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
		w.Header().Set("Content-Type", contentType)

		return
	}

	reader := bytes.NewBuffer(data)

	img, err := imaging.Decode(reader, imaging.AutoOrientation(true))
	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	if img.Bounds().Dx() > int(width) || img.Bounds().Dy() > int(height) {
		resized := imaging.Fit(img, int(width), int(height), imaging.MitchellNetravali)
		img = resized
	}

	w.WriteHeader(http.StatusOK)
	err = imaging.Encode(w, img, imaging.JPEG)

	if err != nil {
		handler.WriteResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
}
