package maintenance

import (
	"context"

	"github.com/mangaweb4/mangaweb4-backend/ent"
	"github.com/mangaweb4/mangaweb4-backend/meta"
)

func RebuildThumbnail(client *ent.Client) error {
	allMeta, err := meta.ReadAll(context.Background(), client)
	if err != nil {
		return err
	}

	for _, m := range allMeta {
		meta.DeleteThumbnail(m)
	}

	return nil
}
