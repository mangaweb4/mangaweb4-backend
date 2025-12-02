package meta

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ThumbnailTestSuite struct {
	suite.Suite
}

func TestThumbnailTestSuite(t *testing.T) {
	suite.Run(t, new(ThumbnailTestSuite))
}

func (t *ThumbnailTestSuite) TestDefaultCropHorizontalNarrow() {
	width := 200
	height := 150

	crop := DefaultThumbnailCrop(width, height)
	t.Assertions.Equal(106, crop.Dx())
	t.Assertions.Equal(150, crop.Dy())
	t.Assertions.Equal(0, crop.Min.X)
	t.Assertions.Equal(0, crop.Min.Y)
}

func (t *ThumbnailTestSuite) TestDefaultCropHorizontalWide() {
	width := 500
	height := 150

	crop := DefaultThumbnailCrop(width, height)
	t.Assertions.Equal(106, crop.Dx())
	t.Assertions.Equal(150, crop.Dy())
	t.Assertions.Equal(72, crop.Min.X)
	t.Assertions.Equal(0, crop.Min.Y)
}

func (t *ThumbnailTestSuite) TestDefaultCropVertical() {
	width := 200
	height := 500

	crop := DefaultThumbnailCrop(width, height)
	t.Assertions.Equal(200, crop.Dx())
	t.Assertions.Equal(283, crop.Dy())
	t.Assertions.Equal(0, crop.Min.X)
	t.Assertions.Equal(109, crop.Min.Y)
}

func (t *ThumbnailTestSuite) TestDefaultCropASeries() {
	width := 200
	height := 283

	crop := DefaultThumbnailCrop(width, height)
	t.Assertions.Equal(200, crop.Dx())
	t.Assertions.Equal(283, crop.Dy())
	t.Assertions.Equal(0, crop.Min.X)
	t.Assertions.Equal(0, crop.Min.Y)
}
