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
	ErrEmptyPrompt          = errors.New("kandinsky prompt is empty ")
	ErrEmptyUUID            = errors.New("kandinsky UUID struct is empty ")
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

// Generate image styles
const (
	// Kandinsky style
	KANDINSKY = "KANDINSKY"
	// Detailed photo
	UHD = "UHD"
	// Anime
	ANIME = "ANIME"
	// No style
	DEFAULT = "DEFAULT"
)

type Kandinsky interface {
	SetModel() (int, error)
	GetImageUUID(p Params) (UUID, error)
	CheckImage(u UUID) (Image, error)
}

// Kand struct, all fields are required
// https://fusionbrain.ai/docs/ru/doc/api-dokumentaciya/
type Kand struct {
	// The API key for authenticating requests to the Kandinsky API.
	key string
	// The API secret for authenticating requests to the Kandinsky API.
	secret string
	// Authenticate URL for getting Kandinsky API model.
	authURL string
	// Generate URL for getting image UUID.
	genURL string
	// Check URL for getting Image instance
	checkURL string

	// The current Model selected for generating images, represented by the Model structure.
	Model Model
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
	// Unique identifier of the model
	ID int `json:"id"`
	// Name of the model
	Name string `json:"name"`
	// Version of the model.
	Version float32 `json:"version"`
	// Type of tasks the model is designed for, e.g., "TEXT2IMAGE".
	Type string `json:"type"`
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
	// Desired width of the generated image, perfect if multiple of 8.
	Width int `json:"width"`
	// Desired height of the generated image, perfect if multiple of 8.
	Height int `json:"height"`
	// Number of images to generate, always = 1.
	NumImages int `json:"num_images"`
	//Type of generation, always "GENERATE".
	Type string `json:"type"`
	// Style of the generated image
	// KANDINSKY - kandinsky style
	// UHD - detailed photo
	// ANIME - anime
	// DEFAULT - No style
	Style string `json:"style"`
	// Negative prompts to avoid in the image generation.
	NegativePrompt string `json:"negativePromptUnclip"`
	// Parameters for the generation, including the Query for the image content.
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
	// The UUID of the generated image task.
	ID string `json:"uuid"`
	// The status of the task, e.g., "INITIAL".
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
	// The time at which the error occurred.
	Timestamp string `json:"timestamp"`
	// HTTP status code of the error.
	Status int `json:"status"`
	// Short description of the error.
	Error string `json:"error"`
	// Detailed message about the error.
	Message string `json:"message"`
	// API endpoint at which the error occurred.
	Path string `json:"path"`
}

// New creates a new instance of the Kandinsky client.
func New(key, secret string) (Kandinsky, error) {
	if key == "" {
		return nil, ErrEmptyKey
	}

	if secret == "" {
		return nil, ErrEmptySecret
	}

	k := &Kand{
		key:      key,
		secret:   secret,
		authURL:  "https://api-key.fusionbrain.ai/key/api/v1/models",
		genURL:   "https://api-key.fusionbrain.ai/key/api/v1/text2image/run",
		checkURL: "https://api-key.fusionbrain.ai/key/api/v1/text2image/status/",
		Model:    Model{},
	}

	return k, nil
}

// GetImage return Image struct, generated by Kandinsky API
func GetImage(key, secret string, params Params) (Image, error) {
	i := Image{}
	if key == "" {
		return i, ErrEmptyKey
	}

	if secret == "" {
		return i, ErrEmptySecret
	}

	if params.GenerateParams.Query == "" {
		return i, ErrEmptyPrompt
	}

	k, err := New(key, secret)
	if err != nil {
		return i, err
	}

	_, err = k.SetModel()
	if err != nil {
		return i, err
	}

	u, err := k.GetImageUUID(params)
	if err != nil {
		return i, err
	}

	i, err = k.CheckImage(u)
	if err != nil {
		return i, err
	}

	return i, nil
}

// SetModel sets the model to be used by the Kandinsky client. Return model ID.
// Send auth request to url and set image UUID to Kandinsky instance from json response:
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
func (k *Kand) SetModel() (int, error) {
	// create GET request, set auth headers
	req, err := http.NewRequest(http.MethodGet, k.authURL, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("X-Key", "Key "+k.key)
	req.Header.Add("X-Secret", "Secret "+k.secret)

	// create client and do request to Kandinsky API
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	// check status code from received from API
	err = checkStatusCode(res.StatusCode)
	if err != nil {
		return 0, err
	}

	// unmarshal response
	m := []Model{}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return 0, err
	}

	k.Model = m[0]

	return k.Model.ID, nil
}

// GetImageUUID sends a POST request with parameters to generate an image and returns the UUID.
//
//	{
//		"uuid": "string",
//		"status": "INITIAL"
//	}
func (k *Kand) GetImageUUID(p Params) (UUID, error) {
	u := UUID{}

	// set default
	if k.Model.ID == 0 {
		k.Model.ID = 4
	}

	setDefaultParams(&p)

	// prompt must be not empty
	if p.GenerateParams.Query == "" {
		return u, ErrEmptyPrompt
	}

	// marshall params to json bytes
	b, err := json.Marshal(&p)
	if err != nil {
		return u, err
	}

	// generate command string
	curlCommand := fmt.Sprintf(`curl --location --request POST 'https://api-key.fusionbrain.ai/key/api/v1/text2image/run' --header 'X-Key: Key %s' --header 'X-Secret: Secret %s' -F 'params=%s
	};type=application/json' --form 'model_id="%d"'`, k.key, k.secret, string(b), k.Model.ID)

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

// CheckImage image status using image UUID
func (k *Kand) CheckImage(u UUID) (Image, error) {
	image := Image{}

	if u.ID == "" {
		return image, ErrEmptyUUID
	}

	for {
		// create GET request
		req, err := http.NewRequest(http.MethodGet, k.checkURL+u.ID, nil)
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

// setDefaultParams set empty params
func setDefaultParams(p *Params) {
	if p.Width == 0 {
		p.Width = 128
	}

	if p.Height == 0 {
		p.Height = 128
	}

	if p.NumImages == 0 {
		p.NumImages = 1
	}

	p.Type = "GENERATE"
}
