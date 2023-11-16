package constants

import api "github.com/save-abandoned-projects/ignite/pkg/apis/ignite"

const (
	// Path to directory containing a subdirectory for each image
	IMAGE_DIR = DATA_DIR + api.KindImage

	// Filename for the image file containing the image filesystem
	IMAGE_FS = "image.ext4"
)
