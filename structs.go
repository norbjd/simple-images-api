package main

type ImageMetadata struct {
	Name        string
	Description string
}

type ImageContent struct {
	Content []byte
}

type Image struct {
	Content  ImageContent
	Metadata ImageMetadata
}

type ImageIDWithMetadata struct {
	ID       string
	Metadata ImageMetadata
}
