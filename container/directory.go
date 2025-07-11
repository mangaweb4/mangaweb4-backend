package container

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/facette/natsort"
	"github.com/mangaweb4/mangaweb4-backend/configuration"
	"github.com/mangaweb4/mangaweb4-backend/ent"
)

func isValidDirectory(name string) bool {
	fullpath := filepath.Join(configuration.Get().DataPath, name)

	dir, err := os.Open(fullpath)

	if err != nil {
		return false
	}

	children, err := dir.Readdir(0)
	if err != nil {
		return false
	}

	for _, f := range children {
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}

		if f.IsDir() {
			continue
		}

		if isValidImageFile(f.Name()) {
			return true
		}
	}
	return false
}

type DirectoryContainer struct {
	Meta *ent.Meta
}

func (c *DirectoryContainer) getChildren(ctx context.Context) (children []fs.FileInfo, err error) {
	m := c.Meta
	fullpath := filepath.Join(configuration.Get().DataPath, m.Name)
	dir, err := os.Open(fullpath)
	if err != nil {
		return
	}

	children, err = dir.Readdir(0)

	return
}

func (c *DirectoryContainer) ListItems(ctx context.Context) (names []string, err error) {
	m := c.Meta

	children, err := c.getChildren(ctx)
	if err != nil {
		return
	}

	names = make([]string, len(m.FileIndices))
	for i, f := range m.FileIndices {
		names[i] = children[f].Name()
	}

	return
}

func (c *DirectoryContainer) OpenItem(ctx context.Context, index int) (reader io.ReadCloser, name string, err error) {
	if index >= len(c.Meta.FileIndices) {
		err = fmt.Errorf("invalid item")
		return
	}

	children, err := c.getChildren(ctx)
	if err != nil {
		return
	}

	zf := children[c.Meta.FileIndices[index]]

	if zf == nil {
		err = fmt.Errorf("file not found : %v", index)
		return
	}

	fullpath := filepath.Join(configuration.Get().DataPath, c.Meta.Name, zf.Name())

	reader, err = os.Open(fullpath)
	name = zf.Name()

	return
}

func (c *DirectoryContainer) PopulateImageIndices(ctx context.Context) error {
	m := c.Meta

	children, err := c.getChildren(ctx)
	if err != nil {
		return err
	}

	type fileIndexPair struct {
		Index    int
		FileName string
	}

	var fileNames []fileIndexPair
	for i, f := range children {
		if isValidImageFile(f.Name()) {
			fileNames = append(fileNames,
				fileIndexPair{
					i, f.Name(),
				})
		}
	}

	sort.Slice(fileNames, func(i, j int) bool {
		return natsort.Compare(fileNames[i].FileName, fileNames[j].FileName)
	})

	m.FileIndices = make([]int, len(fileNames))
	for i, p := range fileNames {
		m.FileIndices[i] = p.Index
	}

	return nil
}

func (c *DirectoryContainer) Download(ctx context.Context) (reader io.ReadCloser, filename string, err error) {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	length := len(c.Meta.FileIndices)

	for i := 0; i < length; i++ {
		img, name, e := c.OpenItem(ctx, i)
		if e != nil {
			err = e
			return
		}
		f, e := w.Create(name)

		if e != nil {
			err = e
			return
		}

		_, e = io.Copy(f, img)
		if e != nil {
			err = e
			return
		}
	}

	w.Close()

	reader = io.NopCloser(buf)

	filename = fmt.Sprintf("%s.zip", filepath.Base(c.Meta.Name))

	return
}
