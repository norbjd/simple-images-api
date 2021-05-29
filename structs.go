package main

type ImageMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ImageContent struct {
	Content []byte
}

type Image struct {
	Content  ImageContent
	Metadata ImageMetadata
}

func (i Image) getBinaryContent() []byte {
	return i.Content.Content
}

type ImageIDWithMetadata struct {
	ID       string
	Metadata ImageMetadata
}
