package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

func cached(duration string, handler func(w http.ResponseWriter, r *http.Request) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		content := storage.Get(r.RequestURI)
		if content != nil {
			fmt.Print("Cache Hit!\n")
			enableCors(&w)
			w.Write(content)
		} else {
			c := httptest.NewRecorder()
			err := handler(c, r)

			for k, v := range c.HeaderMap {
				w.Header()[k] = v
			}

			w.WriteHeader(c.Code)
			content := c.Body.Bytes()

			if err == nil {
				if d, err := time.ParseDuration(duration); err == nil {
					fmt.Printf("New page cached: %s for %s\n", r.RequestURI, duration)
					storage.Set(r.RequestURI, content, d)
				} else {
					fmt.Printf("Page not cached. err: %s\n", err)
				}
			}

			enableCors(&w)
			w.Write(content)
		}

	}
}
