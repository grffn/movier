package models

type DocModel struct {
	Name     string   `json:"name" binding:"required"`
	Category string   `json:"category" binding:"required"`
	Tags     []string `json:"tags"`
	Authors  []string `json:"authors"`
	URL      string   `json:"url" binding:"required"`
	MimeType string   `json:"mimetype"`
}
