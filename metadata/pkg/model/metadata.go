package model

import "movieexample.com/gen"

type Metadata struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Director    string `json:"director"`
}

// MetadataToProto Convert a generated metadata struct to
// a generated metadata proto contrapart.
func MetadataToProto(m *Metadata) *gen.Metadata {

	return &gen.Metadata{
		Id:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Director:    m.Director,
	}
}

// MetadataFromProto converts a proto generated metadata to
// a metadata struct
func MetadataFromProto(m *gen.Metadata) *Metadata {
	return &Metadata{
		ID:          m.Id,
		Title:       m.Title,
		Description: m.Description,
		Director:    m.Director,
	}
}
