package kandinsky

import (
	"os"
	"testing"
)

// TestToByte test converting Image instance to byte slice
func TestToByte(t *testing.T) {
	image, err := GetImage(key, secret, params)
	if err != nil {
		t.Errorf("create image error > %s", err)
	}

	emptyImage := new(Image)

	testCases := []struct {
		desc string
		i    *Image
		want error
	}{
		{
			desc: "Successful convert Image to byte",
			i:    image,
			want: nil,
		},
		{
			desc: "Empty Image",
			i:    emptyImage,
			want: ErrEmptyImage,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			b, err := tC.i.ToByte()
			if err == nil {
				if len(b) == 0 {
					t.Errorf("%s: len of byte Image is 0 > %s", tC.desc, err)
				}
			}

			if err != nil {
				if err != tC.want {
					t.Errorf("%s: want: %s, got: %s", tC.desc, tC.want, err)
				}
			}
		})
	}
}

// TestToFile test converting Image instance to os.File
func TestToFile(t *testing.T) {
	image, err := GetImage(key, secret, params)
	if err != nil {
		t.Errorf("create image error > %s", err)
	}

	emptyImage := new(Image)

	testCases := []struct {
		desc string
		i    *Image
		want error
	}{
		{
			desc: "Successful convert Image to file",
			i:    image,
			want: nil,
		},
		{
			desc: "Empty Image",
			i:    emptyImage,
			want: ErrEmptyImage,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			f, err := tC.i.ToFile()
			if err == nil {
				stat, err := f.Stat()
				if err != nil {
					t.Errorf("%s: get file stat error 0 > %s", tC.desc, err)
				}
				if stat.Size() == 0 {
					t.Errorf("%s: size of file is 0 > %s", tC.desc, err)
				}
			}

			if err != nil {
				if err != tC.want {
					t.Errorf("%s: want: %s, got: %s", tC.desc, tC.want, err)
				}
			}
		})
	}
}

// TestSavePNGTo test saving Image to path/name.png
func TestSavePNGTo(t *testing.T) {
	image, err := GetImage(key, secret, params)
	if err != nil {
		t.Errorf("create image error > %s", err)
	}

	emptyImage := new(Image)

	testCases := []struct {
		desc string
		name string
		path string
		i    *Image
		want error
	}{
		{
			desc: "Successful convert Image to PNG",
			name: "name",
			path: "path/",
			i:    image,
			want: nil,
		},
		{
			desc: "Empty file name",
			name: "",
			path: "path/",
			i:    image,
			want: ErrEmptyFileName,
		},
		{
			desc: "Empty file path",
			name: "name",
			path: "",
			i:    image,
			want: ErrEmptyFilePath,
		},
		{
			desc: "Empty Image instance",
			name: "name",
			path: "path/",
			i:    emptyImage,
			want: ErrEmptyImage,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := tC.i.SavePNGTo(tC.name, tC.path)
			ext := ".png"
			if err == nil {
				// check file exists
				if _, err := os.Stat(tC.path + tC.name + ext); err != nil {
					t.Errorf("%s: file %s not created > %s", tC.desc, tC.name, err)
				}
			}

			if err != nil {
				if err != tC.want {
					t.Errorf("%s: want: %s, got: %s", tC.desc, tC.want, err)
				}
			}
		})
	}
}

// TestSavePNGTo test saving Image to path/name.jpg
func TestSaveJPGTo(t *testing.T) {
	image, err := GetImage(key, secret, params)
	if err != nil {
		t.Errorf("create image error > %s", err)
	}

	emptyImage := new(Image)

	testCases := []struct {
		desc string
		name string
		path string
		i    *Image
		want error
	}{
		{
			desc: "Successful convert Image to JPG",
			name: "name",
			path: "path/",
			i:    image,
			want: nil,
		},
		{
			desc: "Empty file name",
			name: "",
			path: "path/",
			i:    image,
			want: ErrEmptyFileName,
		},
		{
			desc: "Empty file path",
			name: "name",
			path: "",
			i:    image,
			want: ErrEmptyFilePath,
		},
		{
			desc: "Empty Image instance",
			name: "name",
			path: "path/",
			i:    emptyImage,
			want: ErrEmptyImage,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := tC.i.SaveJPGTo(tC.name, tC.path)
			ext := ".jpg"
			if err == nil {
				if _, err := os.Stat(tC.path + tC.name + ext); err != nil {
					t.Errorf("%s: file %s not created > %s", tC.desc, tC.name, err)
				}
			}

			if err != nil {
				if err != tC.want {
					t.Errorf("%s: want: %s, got: %s", tC.desc, tC.want, err)
				}
			}
		})
	}
}
