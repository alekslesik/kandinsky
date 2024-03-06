package kandinsky

import (
	"encoding/base64"
	"os"
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

// ToByte Converts the image to a byte slice.
func (i *Image) ToByte() ([]byte, error) {
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
	ext := ".png"

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

// SaveJPGTo saves image as a JPG file to the specified path.
func (i *Image) SaveJPGTo(name, path string) error {
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
