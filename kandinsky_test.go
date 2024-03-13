package kandinsky

// For testing create .env file included:
// KAND_API_KEY=your_key
// KAND_API_SECRET=your_secret

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

// global variables for tests
var (
	// Kandinsky API key
	key string
	// Kandinsky API secret
	secret string
	// auth URL
	aURL string
	// generate image URL
	gURL string
	// url for check image
	cURL string
)

// TestMain runs before that the below test functions
func TestMain(m *testing.M) {
	// load key secret from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("wrong '.env' file is not exists")
	}

	key = os.Getenv("KAND_API_KEY")
	secret = os.Getenv("KAND_API_SECRET")
	aURL = os.Getenv("KAND_API_AUTH_URL")
	gURL = os.Getenv("KAND_API_GEN_URL")
	cURL = os.Getenv("KAND_API_CHECK_URL")

	if key == "" || secret == "" {
		log.Fatal("empty key or secret")
	}

	// run all tests
	code := m.Run()

	// exit
	os.Exit(code)
}

// TestNew common test
func TestNew(t *testing.T) {
	testCases := []struct {
		desc   string
		key    string
		secret string
		kand   interface{}
		err    error
	}{
		{
			desc:   "Successful create Kandinsky instance",
			key:    key,
			secret: secret,
			kand:   &Kand{},
			err:    nil,
		},
		{
			desc:   "Empty Key",
			key:    "",
			secret: secret,
			kand:   nil,
			err:    ErrEmptyKey,
		},
		{
			desc:   "Empty Secret",
			key:    key,
			secret: "",
			kand:   nil,
			err:    ErrEmptySecret,
		},
		{
			desc:   "Empty Key and Secret",
			key:    "",
			secret: "",
			kand:   nil,
			err:    ErrEmptyKey,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			k, err := New(tc.key, tc.secret)
			if err != tc.err {
				t.Errorf("%s: want error '%v', got '%v'", tc.desc, tc.err, err)
				return
			}

			// Check type
			if tc.kand == nil {
				if k != nil {
					t.Errorf("%s: want nil result, got non-nil result", tc.desc)
				}
			} else {
				expectedType := reflect.TypeOf(tc.kand)
				resultType := reflect.TypeOf(k)
				if resultType != expectedType {
					t.Errorf("%s: want type '%s', got type '%s'", tc.desc, expectedType, resultType)
				}
			}
		})
	}
}

// TestSetModel common test
func TestSetModel(t *testing.T) {
	k, err := New(key, secret)
	if err != nil {
		t.Errorf("create Kandinsky instance error > %s", err)
	}

	id, err := k.SetModel(aURL)
	if err != nil {
		t.Errorf("set model error > %s", err)
	}
	if id < 1 {
		t.Errorf("set model error, wrong model id == %d ", id)
	}
}

// TestSetModelEmptyURL tests the empty url parameter passed to the SetModel() function
func TestSetModelURL(t *testing.T) {
	k, err := New(key, secret)
	if err != nil {
		t.Errorf("create Kandinsky instance error > %s", err)
	}

	testCases := []struct {
		desc string
		url  string
		want string
	}{
		{
			desc: "empty url",
			url:  "",
			want: "kandinsky url is not exists in env or in config file",
		},
		{
			desc: "incorrect url",
			url:  "incorrect",
			want: "Get \"incorrect\": unsupported protocol scheme \"\"",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, err = k.SetModel(tC.url)
			fmt.Println(err.Error())
			if err.Error() != tC.want {
				t.Errorf("want: '%s' got: '%v'", tC.want, err)
			}
		})
	}
}

// TestSetModelEmptyURL tests the empty url parameter passed to the SetModel() function
func TestSetModelStatusCodeFromKandinsky(t *testing.T) {
	k401, err := New("wrongKey", "wrongSecret")
	if err != nil {
		t.Fatalf("error create kandinsky instance > %v", err)
	}

	k404, err := New(key, secret)
	if err != nil {
		t.Fatalf("error create kandinsky instance > %v", err)
	}

	testCases := []struct {
		desc string
		url  string
		k    Kandinsky
		e    error
	}{
		{
			desc: "status unauthorized",
			url:  aURL,
			k:    k401,
			e:    ErrUnauthorized,
		},
		{
			desc: "status unauthorized",
			url:  aURL + "wrongSuffix",
			k:    k404,
			e:    ErrNotFound,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if _, err := tC.k.SetModel(tC.url); err != tC.e {
				t.Errorf("want: '%s' got: '%v'", tC.e, err)
			}
		})
	}
}

// TestSetModel common test
func TestGetImageUUID(t *testing.T) {
	k, err := New(key, secret)
	if err != nil {
		t.Errorf("create Kandinsky instance error > %s", err)
	}

	_, err = k.SetModel(aURL)
	if err != nil {
		t.Errorf("set model error > %s", err)
	}

	testCases := []struct {
		desc string
		p    Params
		want string
	}{
		{
			desc: "Successful GetImageUUID",
			p: Params{
				Width:          1024,
				Height:         1024,
				NumImages:      1,
				Type:           "GENERATE",
				Style:          "KANDINSKY",
				NegativePrompt: "",
				GenerateParams: struct {
					Query string "json:\"query\""
				}{
					Query: "black cat",
				},
			},
			want: "",
		},
		{
			desc: "Width or Height less than 128",
			p: Params{
				Width:          127,
				Height:         127,
				NumImages:      1,
				Type:           "GENERATE",
				Style:          "KANDINSKY",
				NegativePrompt: "",
				GenerateParams: struct {
					Query string "json:\"query\""
				}{
					Query: "black cat",
				},
			},
			want: "status 400 Bad Request",
		},
		{
			desc: "Empty Width, Height, NumImages, Type",
			p: Params{
				Width:          127,
				Height:         127,
				NumImages:      1,
				Type:           "GENERATE",
				Style:          "KANDINSKY",
				NegativePrompt: "",
				GenerateParams: struct {
					Query string "json:\"query\""
				}{
					Query: "black cat",
				},
			},
			want: "status 400 Bad Request",
		},
		{
			desc: "NumImages more than 1",
			p: Params{
				Width:          1024,
				Height:         1024,
				NumImages:      2,
				Type:           "GENERATE",
				Style:          "KANDINSKY",
				NegativePrompt: "",
				GenerateParams: struct {
					Query string "json:\"query\""
				}{
					Query: "black cat",
				},
			},
			want: "status 400 Bad Request",
		},
		{
			desc: "Wrong style",
			p: Params{
				Width:          1024,
				Height:         1024,
				NumImages:      1,
				Type:           "GENERATE",
				Style:          "WRONG",
				NegativePrompt: "",
				GenerateParams: struct {
					Query string "json:\"query\""
				}{
					Query: "black cat",
				},
			},
			want: "",
		},
		{
			desc: "Wrong type",
			p: Params{
				Width:          1024,
				Height:         1024,
				NumImages:      1,
				Type:           "WRONG",
				Style:          "KANDINSKY",
				NegativePrompt: "",
				GenerateParams: struct {
					Query string "json:\"query\""
				}{
					Query: "black cat",
				},
			},
			want: "",
		},
		{
			desc: "Empty query",
			p: Params{
				Width:          1024,
				Height:         1024,
				NumImages:      1,
				Type:           "GENERATE",
				Style:          "KANDINSKY",
				NegativePrompt: "",
				GenerateParams: struct {
					Query string "json:\"query\""
				}{
					Query: "",
				},
			},
			want: "kandinsky prompt is empty",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			u, err := k.GetImageUUID(gURL, tC.p)
			if err == nil {
				if u.ID == "" {
					t.Errorf("empty UUID struct > %s", err)
				}
			} else {
				if !strings.Contains(err.Error(), tC.want) {
					t.Errorf("want: %s, got: %s", tC.want, err)
				}
			}
		})
	}
}

// TestCheckImage common test
func TestCheckImage(t *testing.T) {
	k, err := New(key, secret)
	if err != nil {
		t.Errorf("create Kandinsky instance error > %s", err)
	}

	_, err = k.SetModel(aURL)
	if err != nil {
		t.Errorf("set model error > %s", err)
	}

	p := Params{
		Width:          1024,
		Height:         1024,
		NumImages:      1,
		Type:           "GENERATE",
		Style:          "KANDINSKY",
		NegativePrompt: "",
		GenerateParams: struct {
			Query string "json:\"query\""
		}{
			Query: "black cat",
		},
	}

	u, err := k.GetImageUUID(gURL, p)
	if err != nil {
		t.Errorf("get image UUID model error > %s", err)
	}

	time.Sleep(time.Second * 15)

	testCases := []struct {
		desc string
		url  string
		u    UUID
		want error
	}{
		{
			desc: "Successful CheckImage",
			url: cURL,
			u: u,
			want: nil,
		},
		{
			desc: "Empty ULR",
			url: "",
			u: u,
			want: ErrEmptyURL,
		},
		{
			desc: "Empty UUID",
			url: cURL,
			u: UUID{},
			want: ErrEmptyUUID,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			i, err := k.CheckImage(tC.url, tC.u)
			if err == nil {
				if i.Status != "DONE" {
					t.Errorf("error status image > %s", err)
				}
			} else {
				if err != tC.want {
					t.Errorf("want: %s, got: %s", tC.want, err)
				}
			}
		})
	}
}

func TestGetImage(t *testing.T) {
	testCases := []struct {
		desc	string
		key string
		secret string
		p Params
		want error
	}{
		{
			desc: "Successful create Image",
			key: key,
			secret: secret,
			p: Params{
				Width:          1024,
				Height:         1024,
				NumImages:      1,
				Type:           "GENERATE",
				Style:          "KANDINSKY",
				NegativePrompt: "",
				GenerateParams: struct {
					Query string "json:\"query\""
				}{
					Query: "black cat",
				},
			},
			want: nil,
		},
		{
			desc: "Empty key",
			key: "",
			secret: secret,
			p: Params{
				Width:          1024,
				Height:         1024,
				NumImages:      1,
				Type:           "GENERATE",
				Style:          "KANDINSKY",
				NegativePrompt: "",
				GenerateParams: struct {
					Query string "json:\"query\""
				}{
					Query: "black cat",
				},
			},
			want: ErrEmptyKey,
		},
		{
			desc: "Empty Secret",
			key: key,
			secret: "",
			p: Params{
				Width:          1024,
				Height:         1024,
				NumImages:      1,
				Type:           "GENERATE",
				Style:          "KANDINSKY",
				NegativePrompt: "",
				GenerateParams: struct {
					Query string "json:\"query\""
				}{
					Query: "black cat",
				},
			},
			want: ErrEmptySecret,
		},
		{
			desc: "Empty Prompt",
			key: key,
			secret: secret,
			p: Params{
				GenerateParams: struct {
					Query string "json:\"query\""
				}{
					Query: "",
				},
			},
			want: ErrEmptyPrompt,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			i, err := GetImage(tC.key, tC.secret, tC.p)
			if err == nil {
				if i.UUID == "" {
					t.Errorf("Image instance is empty > %v", i)
				}
			}
			if err != tC.want {
				t.Errorf("want: '%s' got: '%v'", tC.want, err)
			}
		})
	}
}
