package meta

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/mangaweb4/mangaweb4-backend/configuration"
	"github.com/mangaweb4/mangaweb4-backend/container"
	"github.com/mangaweb4/mangaweb4-backend/ent"
	"github.com/mangaweb4/mangaweb4-backend/ent/meta"
	tag_util "github.com/mangaweb4/mangaweb4-backend/tag"
	"github.com/rs/zerolog/log"

	"golang.org/x/exp/slices"
	_ "golang.org/x/image/webp"
)

const (
	META_THUMB_LOCATION        = "meta"
	THUMBNAIL_FILENAME_PATTERN = "%d.jpg"

	THUMBNAIL_HEIGHT = 510
)

func NewItem(ctx context.Context, client *ent.Client, name string, ct meta.ContainerType) (i *ent.Meta, err error) {
	createTime := time.Now()

	c := configuration.Get()
	if stat, e := fs.Stat(os.DirFS(c.DataPath), name); e == nil {
		createTime = stat.ModTime()
	}

	i = &ent.Meta{
		Name:          name,
		CreateTime:    createTime,
		Favorite:      false,
		ContainerType: ct,
	}

	if err = GenerateImageIndices(i); err != nil {
		return
	}

	return client.Meta.Create().
		SetName(i.Name).
		SetCreateTime(i.CreateTime).
		SetFileIndices(i.FileIndices).
		SetContainerType(ct).
		Save(ctx)
}

func Open(m *ent.Meta) (reader io.ReadCloser, err error) {
	mutex := new(sync.Mutex)
	mutex.Lock()
	defer mutex.Unlock()

	c := configuration.Get()

	fullpath := filepath.Join(c.DataPath, m.Name)

	reader, err = os.Open(fullpath)
	return
}

type CropDetails struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

func CreateThumbnail(m *ent.Meta) (thumbnail image.Image, err error) {
	mutex := new(sync.Mutex)
	mutex.Lock()
	defer mutex.Unlock()

	c, err := container.CreateContainer(m)
	if err != nil {
		return
	}

	stream, _, err := c.OpenItem(context.Background(), m.ThumbnailIndex)
	if err != nil {
		return
	}

	defer func() { log.Err(stream.Close()).Msg("close thumbnail stream.") }()

	img, err := imaging.Decode(stream, imaging.AutoOrientation(true))
	if err != nil {
		return
	}

	if m.ThumbnailWidth > 0 && m.ThumbnailHeight > 0 {
		img = imaging.Crop(img, image.Rectangle{
			Min: image.Point{
				X: m.ThumbnailX,
				Y: m.ThumbnailY,
			},
			Max: image.Point{
				X: m.ThumbnailX + m.ThumbnailWidth,
				Y: m.ThumbnailY + m.ThumbnailHeight,
			},
		})
	}

	if img.Bounds().Dy() > THUMBNAIL_HEIGHT {
		resized := imaging.Resize(img, 0, THUMBNAIL_HEIGHT, imaging.MitchellNetravali)
		img = resized
	}

	thumbnail = img
	return
}

func CreateThumbnailPath(id int) string {
	c := configuration.Get()
	return filepath.Join(c.CachePath, META_THUMB_LOCATION, fmt.Sprintf(THUMBNAIL_FILENAME_PATTERN, id))
}

func GetThumbnailBytes(m *ent.Meta) (thumbnail []byte, err error) {
	thumbfile := CreateThumbnailPath(m.ID)
	file, err := os.Open(thumbfile)
	buffer := bytes.Buffer{}

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = os.MkdirAll(filepath.Dir(thumbfile), fs.ModePerm)
			if err != nil {
				return
			}
			img, e := CreateThumbnail(m)
			if e != nil {
				err = e
				return
			}

			err = imaging.Save(img, thumbfile, imaging.JPEGQuality(75))
			if err != nil {
				return
			}

			err = imaging.Encode(&buffer, img, imaging.JPEG, imaging.JPEGQuality(75))
			if err != nil {
				return
			}
		} else {
			return
		}
	}
	_, err = io.Copy(&buffer, file)
	if err != nil {
		return
	}

	thumbnail = bytes.Clone(buffer.Bytes())

	return
}

func DeleteThumbnail(m *ent.Meta) error {
	thumbfile := CreateThumbnailPath(m.ID)
	err := os.Remove(thumbfile)

	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return err
}

func GenerateImageIndices(m *ent.Meta) error {
	mutex := new(sync.Mutex)
	mutex.Lock()
	defer mutex.Unlock()

	c, err := container.CreateContainer(m)
	if err != nil {
		return err
	}

	return c.PopulateImageIndices(context.Background())
}

func PopulateTags(ctx context.Context, client *ent.Client, m *ent.Meta) (out *ent.Meta, tags []*ent.Tag, err error) {
	tagStrs := tag_util.ParseTag(m.Name)
	currentTags, _ := m.QueryTags().All(ctx)

	newTags := make([]*ent.Tag, 0)
	for _, t := range tagStrs {

		if slices.ContainsFunc(currentTags, func(tag *ent.Tag) bool {
			return tag.Name == t
		}) {
			continue
		}

		var tag *ent.Tag
		if temp, err := tag_util.Read(ctx, client, t); err != nil {
			tag = &ent.Tag{
				Name: t,
			}

			tag, _ = client.Tag.Create().
				SetName(tag.Name).
				SetHidden(tag.Hidden).
				Save(ctx)

		} else {
			tag = temp
		}
		newTags = append(newTags, tag)
	}

	m, _ = m.Update().
		AddTags(newTags...).
		Save(ctx)

	out = m
	tags = append(currentTags, newTags...)

	return
}
