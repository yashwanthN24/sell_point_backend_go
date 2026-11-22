package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/YashwanthN24/sell_point/internal/constants"
	"github.com/YashwanthN24/sell_point/internal/dtos"
	"github.com/YashwanthN24/sell_point/internal/httpx"
	"github.com/YashwanthN24/sell_point/internal/middlewares"
	"github.com/YashwanthN24/sell_point/internal/storage"
)

type UploadSHander struct {
	logger  *slog.Logger
	storage *storage.Client
}

func NewUploadsHandler(logger *slog.Logger, store *storage.Client) *UploadSHander {
	return &UploadSHander{
		logger:  logger,
		storage: store,
	}
}

// This pattern is called closure factory as inner func forms a closure with outer func
func (uh UploadSHander) Presign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middlewares.RequestIDFromContext(ctx)
	log := uh.logger.With("request_id", requestId)

	userId, ok := middlewares.UserIDFromContext(ctx)

	if !ok {
		log.Error("no userid found in context")
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	// decode the ayload from the post request
	var req dtos.PresignRequest

	// handles case when body is not passed deocde error
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode", "err", err)
		httpx.Error(w, http.StatusBadRequest, "invalid_body", httpx.CodeMalformedJSON)
		return
	}

	// {"files":[]} case no files sent should handle
	if len(req.Files) == 0 {
		log.Error("no files sent")
		httpx.Error(w, http.StatusBadRequest, "files must not be empty", httpx.CodeValidationFailed)
		return

	}

	//? we must even limit the number of max images that can be uploaded
	if len(req.Files) > constants.MAX_IMAGES_PER_LISTING {
		log.Error("Maximum images limit exceeded for a listing",
			"max_images", constants.MAX_IMAGES_PER_LISTING)
		httpx.Error(w, http.StatusBadRequest, fmt.Sprintf("Max %d images can be uploaded per listing", constants.MAX_IMAGES_PER_LISTING), httpx.CodeValidationFailed)
		return
	}

	// validate the content_type and size_byt5es of the8images

	for _, f := range req.Files {
		// check if content type of images is allowed list
		ext, ok := constants.Allowed_Content_Types[f.ContentType]

		// if content-type is not present in allowed content-type list NOT ALLOWED FURTHER
		if !ok {
			log.Error("Only image/jpeg , image/png , image/webp allowed", "allowed-list", constants.GetAllowedTypes())
			httpx.Error(w, http.StatusBadRequest, fmt.Sprintf("Max %d images can be uploaded per listing", constants.MAX_IMAGES_PER_LISTING), httpx.CodeValidationFailed)
			return
		}
		// validate maxallowed size of image based on bytes of the image

		if f.SizeBytes <= 0 || f.SizeBytes > constants.MAX_IMAGE_BYTES {
			log.Error("Max Image Size Reached", "allowed-size-in-bytes", constants.MAX_IMAGE_BYTES)
			httpx.Error(w, http.StatusBadRequest, fmt.Sprintf("size bytes must be between 1 and %d", constants.MAX_IMAGE_BYTES), httpx.CodeValidationFailed)
			return
		}

		key := constants.UploadKey(userId, ext)

		// actually signing presigned url where user can directly uplad images directly to object storage

		fmt.Println(ext, key)
	}

	fmt.Println(req)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("ok"))
}
