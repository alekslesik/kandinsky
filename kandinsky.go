package kandinsky

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	ErrEmptyURL             = errors.New("kandinsky url is not exists in env or in config file")
	ErrEmptyKey             = errors.New("kandinsky auth key is not exists in env or config file")
	ErrEmptySecret          = errors.New("kandinsky auth secret is not exists in env or config file")
	ErrAuth                 = errors.New("kandinsky authentication error, check your key and secret")
	ErrStatusNot200         = errors.New("kandinsky status is not 200")
	ErrTaskNotCompleted     = errors.New("kandinsky the task could not be completed")
	ErrNotFound             = errors.New("kandinsky resource not found")
	ErrUnauthorized         = errors.New("kandinsky authentication error, check your key and secret")
	ErrInternalServerError  = errors.New("kandinsky server error")
	ErrUnsupportedMediaType = errors.New("kandinsky is not support format")
	ErrBadRequest           = errors.New("kandinsky wrong request parameters or prompt too long ")
)

const (
	StatusBadRequest           = 400
	StatusUnauthorized         = 401
	StatusNotFound             = 404
	StatusInternalServerError  = 500
	StatusUnsupportedMediaType = 415
)

const (
	URLMODEL = "https://api-key.fusionbrain.ai/key/api/v1/models"
	URLUUID  = "https://api-key.fusionbrain.ai/key/api/v1/text2image/run"
	URLCHECK = "https://api-key.fusionbrain.ai/key/api/v1/text2image/status/"
)

// Kandinsky struct, all fields are required
// https://fusionbrain.ai/docs/ru/doc/api-dokumentaciya/
type Kandinsky struct {
	key    string
	secret string
	model  Model
}

// Model is the message from kandinsky API after auth
// [
//
//	{
//	    "id": 4,
//	    "name": "Kandinsky",
//	      "version": 3.0,
//	      "type": "TEXT2IMAGE"
//	}
//
// ]
type Model struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Version float32 `json:"version"`
	Type    string  `json:"type"`
}

// Params for generate image
//
//	{
//		"type": "GENERATE",
//		"style": "string",
//		"width": 1024,
//		"height": 1024,
//		"num_images": 1,
//		"negativePromptUnclip": "яркие цвета, кислотность, высокая контрастность",
//		"generateParams": {
//			"query": "Пушистый кот в очках",
//		}
//	}
type Params struct {
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	NumImages      int    `json:"num_images"`
	Type           string `json:"type"`
	Style          string `json:"style"`
	NegativePrompt string `json:"negativePromptUnclip"`
	GenerateParams struct {
		Query string `json:"query"`
	} `json:"generateParams"`
}

// UUID response with UUID from Kandinsky API
//
//	{
//		"uuid": "string",
//		"status": "INITIAL"
//	}
type UUID struct {
	ID     string `json:"uuid"`
	Status string `json:"status"`
}

// ErrResponse from Kandinsky API
//
//	{
//		"timestamp": "2024-03-04T13:46:55.473+00:00",
//		"status": 400,
//		"error": "Bad Request",
//		"message": "Failed to convert value of type 'java.lang.String' to required type 'int'; For input string: \"\"4\"\"",
//		"path": "/key/api/v1/text2image/run"
//	}
type ErrResponse struct {
	Timestamp string `json:"timestamp"`
	Status    int    `json:"status"`
	Error     string `json:"error"`
	Message   string `json:"message"`
	Path      string `json:"path"`
}

// GetImage return Image struct, generated by Kandinsky API
func GetImage(key, secret string, params Params) (Image, error) {
	i := Image{}

	k, err := New(key, secret)
	if err != nil {
		return i, err
	}

	err = k.SetModel(URLMODEL)
	if err != nil {
		return i, err
	}

	u, err := k.GetImageUUID(URLUUID, params)
	if err != nil {
		return i, err
	}

	i, err = k.Check(URLCHECK, u)
	if err != nil {
		return i, err
	}

	return i, nil
}

// New return Kandinsky instance
func New(key, secret string) (*Kandinsky, error) {
	if key == "" {
		return nil, ErrEmptyKey
	}

	if secret == "" {
		return nil, ErrEmptySecret
	}

	k := &Kandinsky{
		key:    key,
		secret: secret,
		model:  Model{},
	}

	return k, nil
}

// SetModel send auth request to url and set image UUID to Kandinsky instance from json response:
// [
//
//	{
//	    "id": 4,
//	    "name": "Kandinsky",
//	      "version": 3.0,
//	      "type": "TEXT2IMAGE"
//	}
//
// ]
func (k *Kandinsky) SetModel(url string) error {
	// create GET request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	// create auth header
	req.Header.Add("X-Key", "Key "+k.key)
	req.Header.Add("X-Secret", "Secret "+k.secret)

	// create client
	client := http.Client{}
	// Do request to Kandinsky API
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// check status code from received from API
	err = checkStatusCode(res.StatusCode)
	if err != nil {
		return err
	}

	// unmarshal response
	m := []Model{}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return err
	}

	if m[0].Id == 0 {
		return ErrAuth
	}

	k.model = m[0]

	return nil
}

// GetImageUUID send POST request with params to url and return response:
//
//	{
//		"uuid": "string",
//		"status": "INITIAL"
//	}
func (k *Kandinsky) GetImageUUID(url string, params Params) (UUID, error) {
	u := UUID{}

	if k.model.Id == 0 {
		k.model.Id = 4
	}

	// marshall params to json bytes
	b, err := json.Marshal(&params)
	if err != nil {
		return u, err
	}

	// generate command string
	curlCommand := fmt.Sprintf(`curl --location --request POST 'https://api-key.fusionbrain.ai/key/api/v1/text2image/run' --header 'X-Key: Key %s' --header 'X-Secret: Secret %s' -F 'params=%s
	};type=application/json' --form 'model_id="%d"'`, k.key, k.secret, string(b), k.model.Id)

	// create command
	cmd := exec.Command("sh", "-c", curlCommand)

	// buffers for standard out and error
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// run command
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return u, err
	}

	// out to string
	s := out.String()
	// if response status not 200
	if strings.Contains(s, "error") {
		e := ErrResponse{}
		err = json.Unmarshal(out.Bytes(), &e)
		if err != nil {
			return u, err
		}

		return u, errors.New("error from Kandinsky API: status " + strconv.Itoa(e.Status) + " " + e.Error + " > " + e.Message)
	}

	// unmarshal out data to UUID struct
	err = json.Unmarshal(out.Bytes(), &u)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (k *Kandinsky) Check(url string, u UUID) (Image, error) {
	image := Image{}

	for {
		// create GET request
		req, err := http.NewRequest(http.MethodGet, url+u.ID, nil)
		if err != nil {
			return image, err
		}

		// create auth header
		req.Header.Add("X-Key", "Key "+k.key)
		req.Header.Add("X-Secret", "Secret "+k.secret)

		// create client
		client := http.Client{}

		// Do request to Kandinsky API
		res, err := client.Do(req)
		if err != nil {
			return image, err
		}

		// check status code from received from API
		err = checkStatusCode(res.StatusCode)
		if err != nil {
			return image, err
		}

		b, err := io.ReadAll(res.Body)
		if err != nil {
			return image, err
		}

		err = json.Unmarshal(b, &image)
		if err != nil {
			return image, err
		}

		if image.Status == "DONE" {
			return image, nil
		} else if image.Status == "FAIL" {
			return image, ErrTaskNotCompleted
		}

		time.Sleep(time.Second * 10)
	}
}

// checkStatusCode check response code from kandinsky
func checkStatusCode(code int) error {
	switch code {
	case StatusBadRequest:
		return ErrBadRequest
	case StatusUnauthorized:
		return ErrUnauthorized
	case StatusNotFound:
		return ErrNotFound
	case StatusInternalServerError:
		return ErrInternalServerError
	case StatusUnsupportedMediaType:
		return ErrUnsupportedMediaType
	default:
		return nil
	}
}
