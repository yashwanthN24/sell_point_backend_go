package dtos

// {
// 	"files":[
// 		{
// 			"content_type": "image/jpeg",
// 			"size_bytes": 12000
// 		}
// 	]
// }

// type PresignRequest struct {
// 	Files []struct {
// 		ContentType string `json:"content_type"`
// 		SizeBytes   int64  `json:"size_bytes"`
// 	} `json:"files"`
// }

// or 


type PresignFileRequest struct {
    ContentType string `json:"content_type"` // MIME Type of the meida file
    SizeBytes   int64  `json:"size_bytes"` // size in bytes 
}

type PresignRequest struct {
    Files []PresignFileRequest `json:"files"`
}

// https://chatgpt.com/c/6afd68d8-4494-83ee-92db-1585e2244905