package model

const (
	ResizeFileDir = "resize"
	ImagesFileDir = "images"
)

type GetUploadURLRequest struct {
	Bucket      string `json:"bucket,omitempty"`
	File        string `json:"file,omitempty"`
	ContentType string `json:"contentType,omitempty"`
	Width       int64  `json:"width,omitempty"`
	Height      int64  `json:"height,omitempty"`
}

type GetUploadURLResponse struct {
	URL     string            `json:"url,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

type GetDownloadURLRequest struct {
	Scheme string `json:"scheme,omitempty"`
	Domain string `json:"domain,omitempty"`
	File   string `json:"file,omitempty"`
}

type GetDownloadURLResponse struct {
	URL string `json:"url,omitempty"`
}
