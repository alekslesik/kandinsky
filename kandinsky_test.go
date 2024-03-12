package kandinsky

// For testing create .env file included:
// KAND_API_KEY=your_key
// KAND_API_SECRET=your_secret

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

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

	if key == "" || secret == "" {
		log.Fatal("empty key or secret")
	}

	// run all tests
	code := m.Run()

	// exit
	os.Exit(code)
}

// TestNew tests the New() function
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

// TestSetModel do common test
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

// Тем не менее, табличные тесты всё же могут быть использованы для организации некоторых частей тестирования,
// особенно когда дело доходит до тестирования различных статус-кодов ответа или других сценариев, которые легко параметризовать.
// Например, вы можете создать табличные тесты для проверки различных ответов сервера
// (успешные, ошибки клиента, ошибки сервера, невалидный JSON и т.д.), используя mock HTTP-сервер или клиент.

// tests := []struct {
// 	name           string
// 	mockStatusCode int
// 	mockBody       string
// 	expectedError  error
// }{
// 	{"Valid Response", http.StatusOK, `[{"id":4,"name":"Kandinsky","version":3.0,"type":"TEXT2IMAGE"}]`, nil},
// 	{"Unauthorized", http.StatusUnauthorized, ``, ErrAuth},
// 	{"Not Found", http.StatusNotFound, ``, checkStatusCode(http.StatusNotFound)},
// 	{"Bad Gateway", http.StatusBadGateway, ``, checkStatusCode(http.StatusBadGateway)},
// 	{"Invalid JSON", http.StatusOK, `Invalid JSON`, json.UnmarshalTypeError{}},
// }

// for _, tt := range tests {
// 	t.Run(tt.name, func(t *testing.T) {
// 			// Настройте ваш mock HTTP-сервер или клиент здесь для возврата `tt.mockStatusCode` и `tt.mockBody`

// 			// Вызовите SetModel и проверьте ожидаемый результат
// 			k := Kandinsky{}
// 			err := k.SetModel("mockURL")

// 			// Проверка возвращаемой ошибки в зависимости от сценария
// 			if !errors.As(err, &tt.expectedError) {
// 					t.Errorf("Expected error %v, got %v", tt.expectedError, err)
// 			}
// 	})
// }

// В этом примере используется параметризация для проверки различных ответов и ожидаемых ошибок.
// Вы должны настроить mock HTTP-сервер или клиент,
// чтобы он возвращал соответствующие статус-коды и тела ответов в соответствии с каждым тестовым случаем.
