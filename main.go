package main

import (
	"log"
	"net/http"

	"githhub.com/alekslesik/learn/kandinsky"
)

// https://api-key.fusionbrain.ai/key/api/v1/text2image/run
// FFDB0757E5E5A2FF2D7A297CE95BDA2D
// F9C4980BC3166E19C9A42607358B8DA6

func main() {
	key := "FFDB0757E5E5A2FF2D7A297CE95BDA2D"
	secret := "F9C4980BC3166E19C9A42607358B8DA6"

	p := kandinsky.Params{
		Type:           "GENERATE",
		Style:          "UHD",
		Width:          1024,
		Height:         1024,
		NumImages:      1,
		NegativePrompt: "",
		GenerateParams: struct {
			Query string "json:\"query\""
		}{"Обнаружены доказательства гипотезы РНК-мира"},
	}

	image := make(chan kandinsky.Image)

	go func() {
		for {
			i, err := kandinsky.GetImage(key, secret, p)
			if err != nil {
				log.Fatal(err)
			}

			image <- i
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	errCh := make(chan error)

	go func() {
		err := http.ListenAndServe("localhost:6666", mux)
		if err != nil {
			errCh <- err
		}
	}()

	log.Println("start on http://localhost:6666")

	for {
		select {
		case err := <-errCh:
			if err == kandinsky.ErrTaskNotCompleted {
				log.Println(err)
				continue
			}
			log.Fatal(err)
		case  <-image:
			i := <-image
			err := i.SaveJPGTo("name", "./")
			if err != nil {
				log.Fatal(err)
			}

			log.Println("Image created")
			break
		}
	}

}
