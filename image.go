package kandinsky

import (
	"encoding/base64"
	"errors"
	"os"
	"strings"
)

// Image represents the image data returned by the Kandinsky API.
type Image struct {
	// The unique identifier of the image generation task.
	UUID string `json:"uuid"`
	// The status of the image generation task.
	Status string `json:"status"`
	// A slice containing the generated images encoded in Base64.
	Images []string `json:"images"`
	// Indicates whether the image has been censored.
	Censored bool `json:"censored"`
}

var (
	ErrEmptyImage      = errors.New("kandinsky image is empty")
	ErrEmptyFileName   = errors.New("kandinsky file name is empty")
	ErrEmptyFilePath   = errors.New("kandinsky file path is empty")
	ErrEmptyBase       = errors.New("kandinsky base is empty")
	ErrNotBase64Format = errors.New("kandinsky string is not base64 format")
)

// AddBase64 add base64 to Image.
func (i *Image) AddBase64(base string) error {
	if base == "" {
		return ErrEmptyBase
	}

	if !isValidBase64(base) {
		return ErrNotBase64Format
	}

	if len(i.Images) == 0 {
		i.Images = append(i.Images, "")
	}

	i.Images[0] = base

	return nil
}

// ToByte Converts the image to a byte slice.
func (i *Image) ToByte() ([]byte, error) {
	if len(i.Images) == 0 {
		return nil, ErrEmptyImage
	}

	l := len(i.Images[0])

	var b = make([]byte, l)

	_, err := base64.StdEncoding.Decode(b, []byte(i.Images[0]))
	if err != nil {
		return nil, err
	}

	return b, nil
}

// ToFile Converts the image to an os.File.
func (i *Image) ToFile() (*os.File, error) {
	if len(i.Images) == 0 {
		return nil, ErrEmptyImage
	}

	f, err := os.OpenFile(".temp.png", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return f, err
	}

	data, err := base64.StdEncoding.DecodeString(i.Images[0])
	if err != nil {
		return f, err
	}

	_, err = f.Write(data)
	if err != nil {
		return f, err
	}

	return f, nil
}

// SavePNGTo saves the image as a PNG file to the specified path.
func (i *Image) SavePNGTo(name, path string) error {
	if len(i.Images) == 0 {
		return ErrEmptyImage
	}

	if name == "" {
		return ErrEmptyFileName
	}

	if path == "" {
		return ErrEmptyFilePath
	}

	ext := ".png"

	trimName := strings.Trim(path+name+ext, "\"")

	f, err := os.OpenFile(trimName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := base64.StdEncoding.DecodeString(i.Images[0])
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// SaveJPGTo saves image as a JPG file to the specified path.
func (i *Image) SaveJPGTo(name, path string) error {
	if len(i.Images) == 0 {
		return ErrEmptyImage
	}

	if name == "" {
		return ErrEmptyFileName
	}

	if path == "" {
		return ErrEmptyFilePath
	}

	ext := ".jpg"

	f, err := os.OpenFile(path+name+ext, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := base64.StdEncoding.DecodeString(i.Images[0])
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// isValidBase64 check that s is base64
func isValidBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

// trimName trims quotes, spaces and dashes
func trimName(name string) string {
	trimFunc := make(map[string]func(name string) string)

	trimFunc["quotes"] = func(name string) string {
		return strings.Trim(name, "\"")
	}

	trimFunc["spaces"] = func(name string) string {
		return strings.Trim(name, " ")
	}

	trimFunc["dash"] = func(name string) string {
		return strings.Trim(name, "-")
	}

	for _, v := range trimFunc {
		name = v(name)
	}

	return name
}
