package server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync"

	"github.com/disintegration/imaging"
	"github.com/mangaweb4/mangaweb4-backend/container"
	"github.com/mangaweb4/mangaweb4-backend/database"
	ent_meta "github.com/mangaweb4/mangaweb4-backend/ent/meta"
	"github.com/mangaweb4/mangaweb4-backend/ent/progress"
	ent_tag "github.com/mangaweb4/mangaweb4-backend/ent/tag"
	"github.com/mangaweb4/mangaweb4-backend/grpc"
	"github.com/mangaweb4/mangaweb4-backend/meta"
	"github.com/mangaweb4/mangaweb4-backend/tag"
	"github.com/mangaweb4/mangaweb4-backend/user"
	"github.com/rs/zerolog/log"
	grpclib "google.golang.org/grpc"
)

const MESSAGE_SIZE = 1024 * 1024

type MangaServer struct {
	progressMutex sync.Mutex
	grpc.UnimplementedMangaServer
}

func (s *MangaServer) List(ctx context.Context, req *grpc.MangaListRequest) (
	resp *grpc.MangaListResponse, err error,
) {
	defer func() { log.Err(err).Interface("request", req).Msg("MangaServer.List") }()
	client := database.CreateEntClient()
	defer func() { log.Err(client.Close()).Msg("database client close on MangaServer.List") }()

	u, err := user.GetUser(ctx, client, req.User)
	if err != nil {
		return
	}

	allMeta, err := meta.ReadPage(ctx,
		client,
		u,
		meta.QueryParams{
			SearchName:  req.Search,
			Filter:      req.Filter,
			SearchTag:   req.Tag,
			SortBy:      req.Sort,
			SortOrder:   req.Order,
			Page:        int(req.Page),
			ItemPerPage: int(req.ItemPerPage),
		},
	)

	if err != nil {
		return
	}

	items := make([]*grpc.MangaListResponseItem, len(allMeta))
	for i, m := range allMeta {
		progress, _ := client.Progress.Query().
			Where(progress.UserID(u.ID), progress.ItemID(m.ID)).
			Only(ctx)

		currentPage := 0
		maxProgress := 0
		if progress != nil {
			currentPage = progress.Page
			maxProgress = progress.Max
		}

		items[i] = &grpc.MangaListResponseItem{
			ID:          int32(m.ID),
			Name:        m.Name,
			IsFavorite:  u.QueryFavoriteItems().Where(ent_meta.ID(m.ID)).ExistX(ctx),
			IsRead:      progress != nil,
			PageCount:   int32(len(m.FileIndices)),
			CurrentPage: int32(currentPage),
			MaxProgress: int32(maxProgress),
		}

		tags, e := m.QueryTags().All(ctx)
		if e != nil {
			err = e
			return
		}

		for _, t := range tags {
			if u.QueryFavoriteTags().Where(ent_tag.ID(t.ID)).ExistX(ctx) {
				items[i].HasFavoriteTag = true
				break
			}
		}
	}

	count, err := meta.Count(ctx,
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
		return
	}

	pageCount := int32(count) / req.ItemPerPage
	if int32(count)%req.ItemPerPage > 0 {
		pageCount++
	}

	if req.Page > pageCount || req.Page < 0 {
		req.Page = 0
	}

	log.Info().
		Interface("request", req).
		Msg("Browse")

	resp = &grpc.MangaListResponse{
		Items:     items,
		TotalPage: pageCount,
	}

	if req.Tag != "" {
		tagObj, e := tag.Read(ctx, client, req.Tag)
		if e != nil {
			err = e
			return
		}

		resp.TagFavorite = u.QueryFavoriteTags().Where(ent_tag.ID(tagObj.ID)).ExistX(ctx)
	}

	return
}

func (s *MangaServer) Detail(
	ctx context.Context,
	req *grpc.MangaDetailRequest,
) (resp *grpc.MangaDetailResponse, err error) {
	defer func() { log.Err(err).Interface("request", req).Msg("MangaServer.Detail") }()

	client := database.CreateEntClient()
	defer func() { log.Err(client.Close()).Msg("database client close on MangaServer.Detail") }()

	m, err := client.Meta.Get(ctx, int(req.ID))
	if err != nil {
		return
	}

	log.Info().
		Interface("request", req).
		Msg("View Item")

	tags, err := m.QueryTags().All(ctx)
	if err != nil {
		return
	}

	u, err := user.GetUser(ctx, client, req.User)
	if err != nil {
		return
	}

	progress, _ := client.Progress.Query().
		Where(progress.UserID(u.ID), progress.ItemID(m.ID)).
		Only(ctx)

	currentPage := 0
	if progress != nil {
		currentPage = progress.Page
	}

	grpcTags := make([]*grpc.MangaDetailResponseTagItem, len(tags))
	for i := range tags {
		grpcTags[i] = &grpc.MangaDetailResponseTagItem{
			ID:         int32(tags[i].ID),
			Name:       tags[i].Name,
			IsFavorite: u.QueryFavoriteTags().Where(ent_tag.ID(tags[i].ID)).ExistX(ctx),
			IsHidden:   tags[i].Hidden,
		}
	}

	resp = &grpc.MangaDetailResponse{
		Name:        m.Name,
		Favorite:    u.QueryFavoriteItems().Where(ent_meta.ID(m.ID)).ExistX(ctx),
		Tags:        grpcTags,
		PageCount:   int32(len(m.FileIndices)),
		CurrentPage: int32(currentPage),
	}

	_, err = client.History.Create().
		SetUser(u).
		SetItem(m).
		Save(ctx)

	return
}

func (s *MangaServer) Thumbnail(
	ctx context.Context,
	req *grpc.MangaThumbnailRequest,
) (resp *grpc.MangaThumbnailResponse, err error) {
	defer func() { log.Err(err).Interface("request", req).Msg("MangaServer.Thumbnail") }()

	client := database.CreateEntClient()
	defer func() { log.Err(client.Close()).Msg("database client close on MangaServer.Thumbnail") }()

	m, err := client.Meta.Get(ctx, int(req.ID))
	if err != nil {
		return
	}

	thumbnail, err := meta.GetThumbnailBytes(m)
	if err != nil {
		return
	}

	resp = &grpc.MangaThumbnailResponse{
		ContentType: "image/jpeg",
		Data:        thumbnail,
	}

	return
}

func (s *MangaServer) SetFavorite(
	ctx context.Context,
	req *grpc.MangaSetFavoriteRequest,
) (resp *grpc.MangaSetFavoriteResponse, err error) {
	defer func() { log.Err(err).Interface("request", req).Msg("MangaServer.SetFavorite") }()

	client := database.CreateEntClient()
	defer func() { log.Err(client.Close()).Msg("database client close on MangaServer.SetFavorite") }()

	u, err := user.GetUser(ctx, client, req.User)
	if err != nil {
		return
	}

	m, err := client.Meta.Get(ctx, int(req.ID))
	if err != nil {
		return
	}

	if req.Favorite {
		_, err = u.Update().AddFavoriteItems(m).Save(ctx)
	} else {
		_, err = u.Update().RemoveFavoriteItems(m).Save(ctx)
	}

	if err != nil {
		return
	}

	resp = &grpc.MangaSetFavoriteResponse{
		Favorite: req.Favorite,
		Name:     m.Name,
	}

	return
}

func (s *MangaServer) SetProgress(
	ctx context.Context,
	req *grpc.MangaSetProgressRequest,
) (resp *grpc.MangaSetProgressResponse, err error) {
	defer func() { log.Err(err).Interface("request", req).Msg("MangaServer.SetProgress") }()

	client := database.CreateEntClient()
	defer func() { log.Err(client.Close()).Msg("database client close on MangaServer.SetProgress") }()
	m, err := client.Meta.Get(ctx, int(req.ID))
	if err != nil {
		return
	}

	u, err := user.GetUser(ctx, client, req.User)
	if err == nil {
		s.progressMutex.Lock()
		defer s.progressMutex.Unlock()

		progressRec, _ := client.Progress.Query().
			Where(progress.UserID(u.ID), progress.ItemID(m.ID)).
			Only(ctx)

		if progressRec == nil {
			_, err = client.Progress.Create().
				SetPage(int(req.Page)).
				SetMax(int(0)).
				SetItem(m).
				SetUser(u).
				Save(ctx)
		} else {
			max := max(progressRec.Max, int(req.Page))
			_, err = progressRec.Update().
				SetPage(int(req.Page)).
				SetMax(max).
				SetItem(m).
				SetUser(u).
				Save(ctx)
		}

		if err != nil {
			return
		}
	}

	resp = &grpc.MangaSetProgressResponse{
		Name:    m.Name,
		User:    req.User,
		Page:    req.Page,
		Succeed: true,
	}

	return
}

func (s *MangaServer) UpdateCover(
	ctx context.Context,
	req *grpc.MangaUpdateCoverRequest,
) (resp *grpc.MangaUpdateCoverResponse, err error) {
	defer func() { log.Err(err).Interface("request", req).Msg("MangaServer.UpdateCover") }()

	client := database.CreateEntClient()
	defer func() { log.Err(client.Close()).Msg("database client close on MangaServer.UpdateCover") }()

	m, err := client.Meta.Get(ctx, int(req.ID))
	if err != nil {
		return
	}

	m.ThumbnailIndex = int(req.Index)
	m.ThumbnailHeight = int(req.Height)
	m.ThumbnailWidth = int(req.Width)
	m.ThumbnailX = int(req.X)
	m.ThumbnailY = int(req.Y)

	err = meta.DeleteThumbnail(m)
	if err != nil {
		return
	}

	err = meta.Write(ctx, client, m)
	if err != nil {
		return
	}

	resp = &grpc.MangaUpdateCoverResponse{
		Success: true,
	}

	return
}

func (s *MangaServer) PageImage(
	ctx context.Context,
	req *grpc.MangaPageImageRequest,
) (resp *grpc.MangaPageImageResponse, err error) {
	defer func() { log.Err(err).Interface("request", req).Msg("MangaServer.PageImage") }()

	client := database.CreateEntClient()
	defer func() { log.Err(client.Close()).Msg("database client close on MangaServer.PageImage") }()

	m, err := client.Meta.Get(ctx, int(req.ID))
	if err != nil {
		return
	}

	c, err := container.CreateContainer(m)
	if err != nil {
		return
	}

	steam, f, err := c.OpenItem(context.Background(), int(req.Index))
	if err != nil {
		return
	}

	data, err := io.ReadAll(steam)
	if err != nil {
		return
	}

	if req.Width == 0 && req.Height == 0 {
		var contentType string
		switch filepath.Ext(strings.ToLower(f)) {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".webp":
			contentType = "image/webp"
		}

		resp = &grpc.MangaPageImageResponse{
			ContentType: contentType,
			Data:        data,
		}

		return
	}

	reader := bytes.NewBuffer(data)

	img, err := imaging.Decode(reader, imaging.AutoOrientation(true))
	if err != nil {
		return
	}

	if img.Bounds().Dx() > int(req.Width) || img.Bounds().Dy() > int(req.Height) {
		resized := imaging.Fit(img, int(req.Width), int(req.Height), imaging.MitchellNetravali)
		img = resized
	}

	var buf bytes.Buffer

	err = imaging.Encode(&buf, img, imaging.JPEG)

	if err != nil {
		return
	}

	resp = &grpc.MangaPageImageResponse{
		ContentType: "image/jpeg",
		Data:        buf.Bytes(),
	}

	return
}

func (s *MangaServer) PageImageStream(req *grpc.MangaPageImageRequest,
	stream grpclib.ServerStreamingServer[grpc.MangaPageImageStreamResponse]) error {

	var err error
	var ctx = context.Background()
	var contentType string
	var filename string

	defer func() { log.Err(err).Interface("request", req).Msg("MangaServer.PageImageStream") }()

	client := database.CreateEntClient()
	defer func() { log.Err(client.Close()).Msg("database client close on MangaServer.PageImageStream") }()

	m, err := client.Meta.Get(ctx, int(req.ID))
	if err != nil {
		return err
	}

	c, err := container.CreateContainer(m)
	if err != nil {
		return err
	}

	fstream, filename, err := c.OpenItem(context.Background(), int(req.Index))
	if err != nil {
		return err
	}

	data, err := io.ReadAll(fstream)
	if err != nil {
		return err
	}

	u, err := user.GetUser(ctx, client, req.User)
	if err == nil {
		s.progressMutex.Lock()
		defer s.progressMutex.Unlock()

		progressRec, _ := client.Progress.Query().
			Where(progress.UserID(u.ID), progress.ItemID(m.ID)).
			Only(ctx)

		if progressRec == nil {
			_, err = client.Progress.Create().
				SetPage(int(req.Index)).
				SetMax(int(0)).
				SetItem(m).
				SetUser(u).
				Save(ctx)
		} else {
			max := max(progressRec.Max, int(req.Index))
			_, err = progressRec.Update().
				SetPage(int(req.Index)).
				SetMax(max).
				SetItem(m).
				SetUser(u).
				Save(ctx)
		}

		if err != nil {
			return err
		}
	}
	if req.Width == 0 && req.Height == 0 {
		switch filepath.Ext(strings.ToLower(filename)) {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".webp":
			contentType = "image/webp"
		default:
			contentType = ""
		}

	} else {
		reader := bytes.NewBuffer(data)

		img, err := imaging.Decode(reader, imaging.AutoOrientation(true))
		if err != nil {
			return err
		}

		if img.Bounds().Dx() > int(req.Width) || img.Bounds().Dy() > int(req.Height) {
			resized := imaging.Fit(img, int(req.Width), int(req.Height), imaging.MitchellNetravali)
			img = resized
		}

		var buf bytes.Buffer

		err = imaging.Encode(&buf, img, imaging.JPEG)

		if err != nil {
			return err
		}

		filename = fmt.Sprintf("%s.jpeg", filepath.Base(filename))
	}

	length := len(data)

	for i := 0; i < length; i += MESSAGE_SIZE {
		end := min(i+MESSAGE_SIZE, length)
		err = stream.Send(&grpc.MangaPageImageStreamResponse{
			Filename:    filename,
			ContentType: contentType,
			Data:        data[i:end],
			Size:        int32(end - i),
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *MangaServer) Repair(
	ctx context.Context,
	req *grpc.MangaRepairRequest,
) (resp *grpc.MangaRepairResponse, err error) {
	defer func() { log.Err(err).Interface("request", req).Msg("MangaServer.Repair") }()

	client := database.CreateEntClient()
	defer func() { log.Err(client.Close()).Msg("database client close on MangaServer.Repair") }()

	m, err := client.Meta.Get(ctx, int(req.ID))
	if err != nil {
		return
	}

	if m, _, err = meta.PopulateTags(ctx, client, m); err != nil {
		return
	}

	if err = meta.GenerateImageIndices(m); err != nil {
		return
	}

	if err = meta.DeleteThumbnail(m); err != nil {
		return
	}

	if err = meta.Write(ctx, client, m); err != nil {
		return
	}

	resp = &grpc.MangaRepairResponse{
		Name:      m.Name,
		IsSuccess: true,
	}

	return
}

func (s *MangaServer) Download(
	req *grpc.MangaDownloadRequest,
	stream grpclib.ServerStreamingServer[grpc.MangaDownloadResponse],
) error {
	var err error

	defer func() { log.Err(err).Interface("request", req).Msg("MangaServer.Download") }()
	ctx := context.Background()

	client := database.CreateEntClient()
	defer func() { log.Err(client.Close()).Msg("database client close on MangaServer.Download") }()

	m, err := client.Meta.Get(ctx, int(req.ID))
	if err != nil {

		return err
	}

	c, err := container.CreateContainer(m)
	if err != nil {
		return err
	}

	reader, filename, err := c.Download(ctx)
	if err != nil {
		return err
	}
	defer func() { log.Err(reader.Close()).Msg("container download.") }()

	bytes, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	length := len(bytes)

	for i := 0; i < length; i += MESSAGE_SIZE {
		end := min(i+MESSAGE_SIZE, length)
		err = stream.Send(&grpc.MangaDownloadResponse{
			Filename:    filename,
			ContentType: "application/zip",
			Data:        bytes[i:end],
			Size:        int32(end - i),
		})

		if err != nil {
			return err
		}
	}

	return nil
}
