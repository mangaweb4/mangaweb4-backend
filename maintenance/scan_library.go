package maintenance

import (
	"context"

	"github.com/mangaweb4/mangaweb4-backend/ent"
	"github.com/mangaweb4/mangaweb4-backend/meta"
	"github.com/rs/zerolog/log"
)

func ScanLibrary(ctx context.Context, client *ent.Client) error {
	allMeta, err := meta.ReadAll(ctx, client)
	if err != nil {
		return err
	}

	files, err := meta.ListDir("")
	if err != nil {
		return err
	}

	for _, file := range files {
		found := false
		for _, m := range allMeta {
			if m.Name == file.Name {
				found = true

				if m.Active {
					break
				}
			}
		}
		if found {
			continue
		}

		log.Info().
			Str("file", file.Name).
			Str("type", file.Type.String()).
			Msg("Creating metadata.")

		if item, err := meta.Read(ctx, client, file.Name); err == nil {
			item.Active = true
			if err := meta.Write(ctx, client, item); err != nil {
				log.Error().Str("name", item.Name).Err(err).Msg("Failed to re-activate meta")
			}
		} else {
			item, err := meta.NewItem(ctx, client, file.Name, file.Type)
			if err != nil {
				log.Error().Err(err).Msg("Failed to create meta data.")

				continue
			}

			item, tags, err := meta.PopulateTags(ctx, client, item)
			if err != nil {
				log.Error().Err(err).Msg("Failed to populate tags.")
				continue
			}

			tagsArray := make([]string, len(tags))

			for i, tag := range tags {
				tagsArray[i] = tag.Name
			}

			log.Info().Str("name", item.Name).Strs("tags", tagsArray).Msg("Created metadata.")
		}
	}

	for _, m := range allMeta {
		found := false
		for _, file := range files {
			if m.Name == file.Name {
				found = true
				break
			}
		}
		if found {
			continue
		}

		log.Info().Str("file", m.Name).Msg("Inactivate metadata.")
		m.Active = false

		if err := meta.Write(ctx, client, m); err != nil {
			log.Error().Str("name", m.Name).Err(err).Msg("Failed to inactivate meta")
		}
	}

	return nil
}
