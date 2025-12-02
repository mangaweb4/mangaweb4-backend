package meta

import (
	"image"
	"math"
)

func DefaultThumbnailCrop(width, height int) image.Rectangle {
	// Check for invalid dimensions
	if width <= 0 || height <= 0 {
		return image.Rect(0, 0, 0, 0)
	}

	const targetAspectRatio = 1.0 / 1.41421356237 // ISO A-series aspect ratio
	aspectRatio := float64(width) / float64(height)

	if aspectRatio > targetAspectRatio {
		// the image will be cropped on the x-axis
		// The center of new image will be at 25% of the original image.
		newWidth := math.Round(float64(height) * targetAspectRatio)
		newX := math.Max(math.Round(float64(width)*0.25)-newWidth/2, 0)

		return image.Rect(int(newX), 0, int(newX+newWidth), height)

	} else if aspectRatio < targetAspectRatio {
		// the image will be cropped on the y-axis
		// Calculate the new height based on the target aspect ratio
		newHeight := int(math.Round(float64(width) / targetAspectRatio))
		newY := int(math.Round(float64(height-newHeight) / 2))

		return image.Rect(0, newY, width, newY+newHeight)
	} else {
		// this is unlikely to happen, but just in case
		// return the original size
		return image.Rect(0, 0, width, height)
	}
}
