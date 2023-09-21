package storage

import "github.com/save-abandoned-projects/libgitops/pkg/serializer"

// ContentTypes describes the connection between
// file extensions and a content types.
var ContentTypes = map[string]serializer.ContentType{
	".json": serializer.ContentTypeJSON,
	".yaml": serializer.ContentTypeYAML,
	".yml":  serializer.ContentTypeYAML,
}

func extForContentType(wanted serializer.ContentType) string {
	for ext, ct := range ContentTypes {
		if ct == wanted {
			return ext
		}
	}
	return ""
}
