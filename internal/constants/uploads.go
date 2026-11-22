package constants

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	MAX_IMAGES_PER_LISTING = 10
	MAX_IMAGE_BYTES        = 5 * 1024 * 1024 // 5MB because 1 mb = 1024 kb and 1kb = 1024 bytes so 5 * 1024 * 1024 (bytes resemble 5mb)
	uploadPrefix           = "uploads"
)

// key content type value is extension this we casn validate both cntent type and extension
var Allowed_Content_Types = map[string]string{
	"image/jpeg": "jpg",
	"image/png":  "png",
	"image/webp": "webp",
}

var allowedTypes []string

func init() {
	allowedTypes = make([]string, 0, len(Allowed_Content_Types))
	for k := range Allowed_Content_Types {
		allowedTypes = append(allowedTypes, k)
	}
}

func GetAllowedTypes() []string {
	// Return a copy so callers cannot corrupt the global slice
	return append([]string(nil), allowedTypes...)
}

// uploads/{user_id}/{uuid}.jog
func UploadKey(userId uuid.UUID, ext string) string {
	return fmt.Sprintf("%s/%s/%s.%s", uploadPrefix, userId, uuid.New(), ext)
}
