package kandinsky

// For testing create .env file included:
// KAND_API_KEY=your_key
// KAND_API_SECRET=your_secret

import (
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
	// correct instance of Params
	params Params
)

// TestMain runs before that the below test functions
func TestMain(m *testing.M) {
	// load key secret from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("wrong '.env' file is not exists")
	}

	key = os.Getenv("KAND_API_KEY")
	secret = os.Getenv("KAND_API_SECRET")

	if key == "" || secret == "" {
		log.Fatal("empty key or secret")
	}

	params = Params{
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

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			k, err := New(tC.key, tC.secret)
			if err != tC.err {
				t.Errorf("%s: want error '%v', got '%v'", tC.desc, tC.err, err)
				return
			}

			// Check type
			if tC.kand == nil {
				if k != nil {
					t.Errorf("%s: want nil result, got non-nil result", tC.desc)
				}
			} else {
				expectedType := reflect.TypeOf(tC.kand)
				resultType := reflect.TypeOf(k)
				if resultType != expectedType {
					t.Errorf("%s: want type '%s', got type '%s'", tC.desc, expectedType, resultType)
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

	id, err := k.SetModel()
	if err != nil {
		t.Errorf("set model error > %s", err)
	}
	if id < 1 {
		t.Errorf("set model error, wrong model id == %d ", id)
	}
}

// TestSetModel common test
func TestGetImageUUID(t *testing.T) {
	k, err := New(key, secret)
	if err != nil {
		t.Errorf("create Kandinsky instance error > %s", err)
	}

	_, err = k.SetModel()
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
			u, err := k.GetImageUUID(tC.p)
			if err == nil {
				if u.ID == "" {
					t.Errorf("%s: empty UUID struct > %s", tC.desc, err)
				}
			}

			if err != nil {
				if !strings.Contains(err.Error(), tC.want) {
					t.Errorf("\n%s:\n\twant:\n\t\t%s, \n\tgot:\n\t\t%s\n", tC.desc, tC.want, err)
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

	_, err = k.SetModel()
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

	u, err := k.GetImageUUID(p)
	if err != nil {
		t.Errorf("get image UUID model error > %s", err)
	}

	pCens := Params{
		Width:          1024,
		Height:         1024,
		NumImages:      1,
		Type:           "GENERATE",
		Style:          "KANDINSKY",
		NegativePrompt: "",
		GenerateParams: struct {
			Query string "json:\"query\""
		}{
			Query: "murder",
		},
	}

	uCens, err := k.GetImageUUID(pCens)
	if err != nil {
		t.Errorf("get image UUID model error > %s", err)
	}

	time.Sleep(time.Second * 15)

	testCases := []struct {
		desc string
		url  string
		u    *UUID
		want error
	}{
		{
			desc: "Successful CheckImage",
			u:    u,
			want: nil,
		},
		{
			desc: "Empty UUID",
			u:    &UUID{},
			want: ErrEmptyUUID,
		},
		{
			desc: "Censored UUID",
			u:    uCens,
			want: ErrCensored,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			i, err := k.CheckImage(tC.u)
			if err == nil {
				if i.Status != "DONE" {
					t.Errorf("%s: error status image > %s", tC.desc, err)
				}
			}

			if err != nil {
				if err != tC.want {
					t.Errorf("\n%s:\n\twant:\n\t\t%s, \n\tgot:\n\t\t%s\n", tC.desc, tC.want, err)
				}
			}
		})
	}
}

func TestGetImage(t *testing.T) {
	testCases := []struct {
		desc   string
		key    string
		secret string
		p      Params
		want   error
	}{
		{
			desc:   "Successful create Image",
			key:    key,
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
			desc:   "Empty key",
			key:    "",
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
			desc:   "Empty Secret",
			key:    key,
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
			desc:   "Empty Prompt",
			key:    key,
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
					t.Errorf("%s: Image instance is empty > %v", tC.desc, i)
				}
			}

			if err != nil {
				if err != tC.want {
					t.Errorf("\n%s:\n\twant:\n\t\t%s, \n\tgot:\n\t\t%v\n", tC.desc, tC.want, err)
				}
			}
		})
	}
}
