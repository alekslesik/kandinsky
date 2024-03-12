# Kandinsky Go Client

This package provides a Go client for the Kandinsky API, enabling the generation of images from text prompts.

## Key Features

- Authenticate with the Kandinsky API using a key and secret.
- Retrieve a list of available models.
- Generate images with specified parameters.
- Check the status of image generation.

## Installation

To use the `kandinsky` package, add it to your Go project:

``` go get github.com/yourusername/kandinsky ```

## Usage Example

To use the Kandinsky client, you need to create an instance of `Kandinsky` with your key and secret:

```go
import "github.com/alekslesik/kandinsky"

func main() {
    key := "your_key_here"
    secret := "your_secret_here"

    k, err := kandinsky.New(key, secret)
    if err != nil {
        log.Fatalf("error creating Kandinsky client: %v", err)
    }

    // Set the model
    err = k.SetModel(kandinsky.URLMODEL)
    if err != nil {
        log.Fatalf("error setting model: %v", err)
    }

    // Parameters for image generation
    params := kandinsky.Params{
        Width:          1024,
        Height:         1024,
        NumImages:      1,
        Type:           "GENERATE",
        Style:          "UHD",
        NegativePrompt: "",
        GenerateParams: struct{
            Query string `json:"query"`
        }{
            Query: "Fluffy cat wearing glasses",
        },
    }

    // Generate the image
    image, err := k.GetImage(key, secret, params)
    if err != nil {
        log.Fatalf("error getting image: %v", err)
    }

    fmt.Println("Generated Image:", image)
}
```
## Info Guide

This guide explains how to ask the Kandinsky API to make an image from words using a POST request. It's just a guide, not a strict rule.

### What You Need in Your Request

For the API to recognize who's asking, your request needs a couple of special notes called headers:

- `X-Key: Key YOUR_KEY` and `X-Secret: Secret YOUR_SECRET`: These are like secret handshakes that prove who you are. You'll swap `YOUR_KEY` and `YOUR_SECRET` with your actual keys.

### Filling Out Your Request

Your request will be a mix of simple info and a file in a format called `multipart/form-data`. This format is perfect for when you need to send a detailed instruction book (like a file) along with your request. Here's what you put in it:

- `model_id="4"`: This tells the API you're choosing model number 1 to create your image.
- `-F 'params=...;type=application/json'`: This part is important because it's like attaching a file to an email. This "file" is written in JSON and tells the API exactly what you want your image to be about. In this case:
  - The `params` "file" tells the API you want to "GENERATE" an image.
  - Inside `params`, `generateParams` describes the image details, like `"query": "kitty"` for an image of a kitty.

### An Easy cURL Command Example

```bash
# This is an example to show you how it's done.
curl --location --request POST 'https://api-key.fusionbrain.ai/key/api/v1/text2image/run' \
--header 'X-Key: Key YOUR_KEY' \
--header 'X-Secret: Secret YOUR_SECRET' \
-F 'params={
  "type": "GENERATE",
  "generateParams": {
    "query": "kitty"
  }
};type=application/json' \
--form 'model_id="4"'
```


## Errors

The `kandinsky` package defines several errors to handle various scenarios of interaction with the API:

- `ErrEmptyURL`: The API URL is not provided.
- `ErrEmptyPrompt`: The API prompt provided.
- `ErrEmptyUUID`: The API UUID provided.
- `ErrEmptyKey`: The authentication key is not provided.
- `ErrEmptySecret`: The authentication secret is not provided.
- `ErrAuth`: Authentication error, indicating issues with the provided key and secret.
- `ErrStatusNot200`: The API response status is not 200, indicating that the request was not successful.
- `ErrTaskNotCompleted`: The task could not be completed, possibly due to an error in processing by the API.
- `ErrNotFound`: The requested resource was not found on the API server.
- `ErrUnauthorized`: Authentication error, similar to `ErrAuth`, indicating issues with the provided credentials.
- `ErrInternalServerError`: The API server encountered an internal error, suggesting a problem on the server-side.
- `ErrUnsupportedMediaType`: The media type provided is not supported by the API, indicating an issue with the format of the request.
- `ErrBadRequest`: The request parameters are incorrect or the prompt is too long, indicating that the client has constructed a bad request.

These errors provide a way to handle specific issues encountered when interacting with the Kandinsky API, allowing for more granular error handling and troubleshooting in client applications.


## Constants

The package includes constants for handling specific HTTP response statuses from the Kandinsky API, providing clear indicators for various outcomes of API requests. These constants are as follows:

- `StatusBadRequest`: Corresponds to HTTP status code 400, indicating that the server could not understand the request due to invalid syntax.
- `StatusUnauthorized`: Corresponds to HTTP status code 401, indicating that the request has not been applied because it lacks valid authentication credentials for the target resource.
- `StatusNotFound`: Corresponds to HTTP status code 404, indicating that the server cannot find the requested resource.
- `StatusInternalServerError`: Corresponds to HTTP status code 500, indicating that the server encountered an unexpected condition that prevented it from fulfilling the request.
- `StatusUnsupportedMediaType`: Corresponds to HTTP status code 415, indicating that the media type of the requested data is not supported by the server, so the server is refusing the request.

Additionally, the package defines URLs used for various API requests:

- `URLMODEL`: The URL endpoint to retrieve the list of available models from the Kandinsky API.
- `URLUUID`: The URL endpoint to initiate a new image generation task and receive a UUID for the task.
- `URLCHECK`: The URL endpoint used to check the status of an image generation task using the task's UUID.

Styles for generate images:

-	`KANDINSKY`: Kandinsky style.
-	`UHD`: Detailed photo.
-	`ANIME`: Anime.
-	`DEFAULT`: No style.

These constants and URL endpoints are integral to the operation of the Kandinsky Go client, streamlining the process of making requests to the Kandinsky API and handling responses.


## API Documentation

For more detailed information about the Kandinsky API and parameters for image generation, [refer to the official Kandinsky API documentation.](https://fusionbrain.ai/docs/ru/doc/api-dokumentaciya/)

This `README.md` file offers a basic guide and usage example of your Go package `kandinsky`. You can expand it with additional setup instructions, requirements, and advanced usage examples according to the capabilities of your API client.

## Structures

### `Kandinsky`

Represents the main client for interacting with the Kandinsky API.

```go
type Kandinsky struct {
	// The API key for authenticating requests to the Kandinsky API.
	key    string
	// The API secret for authenticating requests to the Kandinsky API.
	secret string
	// The current model selected for generating images, represented by the Model structure.
	model  Model
}
```


### `Model`

Represents a model provided by the Kandinsky API. Return model ID.

```go
type Model struct {
	// Unique identifier of the model
	Id      int     `json:"id"`
	// Name of the model
	Name    string  `json:"name"`
	// Version of the model.
	Version float32 `json:"version"`
	// Type of tasks the model is designed for, e.g., "TEXT2IMAGE".
	Type    string  `json:"type"`
}
```

### `Params`
Defines parameters for generating an image.


```go
type Params struct {
	// Desired width of the generated image.
	// Optional. Default 1024. Perfect if multiple of 8. Must be more than 128.
	Width          int    `json:"width"`
	// Desired height of the generated image.
	// Optional. Default 1024. Perfect if multiple of 8. Must be more than 128.
	Height         int    `json:"height"`
	// Number of images to generate.
	// Optional. Must be equal 1.
	NumImages      int    `json:"num_images"`
	// Type of generation.
	// Optional. Default "GENERATE" Must be equal "GENERATE".
	Type           string `json:"type"`
	// Style of the generated image.
	// KANDINSKY - kandinsky style
	// UHD - detailed photo
	// ANIME - anime
	// DEFAULT - No style
	Style          string `json:"style"`
	// Negative prompts to avoid in the image generation.
	NegativePrompt string `json:"negativePromptUnclip"`
	// Parameters for the generation, including the Query for the image content.
	GenerateParams struct {
		// Requirement prompt to generate image
		Query string `json:"query"`
	} `json:"generateParams"`
}
```
### `UUID`

Represents a response containing a UUID from the Kandinsky API.

```go
type Model struct {
	// Unique identifier of the model
	Id      int     `json:"id"`
	// Name of the model
	Name    string  `json:"name"`
	// Version of the model.
	Version float32 `json:"version"`
	// Type of tasks the model is designed for, e.g., "TEXT2IMAGE".
	Type    string  `json:"type"`
}
```
### `ErrResponse`
Defines the structure of an error response from the Kandinsky API.

```go
type ErrResponse struct {
	// The time at which the error occurred.
	Timestamp string `json:"timestamp"`
	// HTTP status code of the error.
	Status    int    `json:"status"`
	// Short description of the error.
	Error     string `json:"error"`
	// Detailed message about the error.
	Message   string `json:"message"`
	// API endpoint at which the error occurred.
	Path      string `json:"path"`
}

```

### `Image`

Represents the image data returned by the Kandinsky API.

```go
type Image struct {
    UUid     string   `json:"uuid"`
    Status   string   `json:"status"`
    Images   []string `json:"images"`
    Censored bool     `json:"censored"`
}
```

## Methods

### `New`

```go
func New(key, secret string) (*Kandinsky, error)
```

Creates a new instance of the Kandinsky client.

- `key`: The API key for authentication.
- `secret`: The API secret for authentication.
- Returns a new Kandinsky instance or an error.

### `SetModel`

```go
func (k *Kandinsky) SetModel(url string) error
```

Sets the model to be used by the Kandinsky client.

- `url`:  The URL to send the generation request to.
- `params`: The parameters for image generation.
- Returns a UUID struct with the task details or an error.

### `GetImageUUID`

Sends a POST request with parameters to generate an image and returns the UUID.

```go
func (k *Kandinsky) GetImageUUID(url string, p Params) (UUID, error)
```
- `url`: The URL to send the generation request to.
- `params`: The parameters for image generation.
- `Returns` a UUID struct with the task details or an error.


### `CheckImage`
Checks the status of an image generation task.
```go
func (k *Kandinsky) CheckImage(url string, u UUID) (Image, error)
```
- `url`: The URL to check the task status.
- `u`: The UUID of the task to check.
- Returns an Image struct with the generated image details or an error.

### `ToByte`
Converts the image to a byte slice.

```go
func (i *Image) ToByte() ([]byte, error)
```
Returns:
- A byte slice representation of the image.
- An error if the Base64 decoding fails.

### `ToFile`
Converts image to an os.File.

```go
func (i *Image) ToFile() (*os.File, error)
```
Returns:
- A file pointer to the newly created file containing the image.
- An error if file creation or Base64 decoding fails.


### `SavePNGTo`
Saves the image as a PNG file to the specified path.
```go
func (i *Image) SavePNGTo(name, path string) error
```
Parameters:
- `name`: The name for the saved file (without extension).
- `path`: The directory path where the file should be saved.

Returns:
- An error if file creation, Base64 decoding, or file writing fails.
SaveJPGTo
- Saves the first image in the Images slice as a JPG file to the specified path.

### `SaveJPGTo`
Saves the image as a JPG file to the specified path.
```go
func (i *Image) SaveJPGTo(name, path string) error
```
Parameters:
- name: The name for the saved file (without extension).
- path: The directory path where the file should be saved.

Returns:
An error if file creation, Base64 decoding, or file writing fails.


This documentation provides users with comprehensive information on how to handle the image data returned from the Kandinsky API, offering flexibility in how they can use or store the generated images.











