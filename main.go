package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"mime"
)

func main() {
	// Config
	host := os.Getenv("HOST")
	if host == "" {
		fmt.Printf("missing host")
	}
	timeout := os.Getenv("TIMEOUT")
	scheme := os.Getenv("SCHEME")
	if scheme == "" {
		scheme = "http"
	}
	proxyTimeout, err := strconv.Atoi(timeout)
	if err != nil {
		fmt.Printf("invalid timeout: %v\n", err)
		return
	}

	// Server
	e := echo.New()
	e.Any("/*", timeoutFixerFor(host, scheme, proxyTimeout))
	e.Start(":8080")
}

func timeoutFixerFor(host, scheme string, timeout int) echo.HandlerFunc {
	return func(c echo.Context) error {
		director := func(req *http.Request) {
			r := c.Request()
			req = r
			req.URL.Scheme = scheme
			req.URL.Host = r.URL.Host
		}
		req := c.Request()
		req.URL.Host = host
		req.Host = host
		req.URL.Scheme = scheme

		proxy := &httputil.ReverseProxy{Director: director}
		recorder := httptest.NewRecorder()
		done := make(chan bool)
		// Replace accept-encoding header
		req.Header.Set("Accept-Encoding", "identity")

		// This might take up to an hour, so spawn a go routine
		go func() {
			proxy.ServeHTTP(recorder, req)
			done <- true
		}()

		// Set correct mime-type
		ext := filepath.Ext(req.URL.Path)
		mimeType := mime.TypeByExtension(ext)
		if mimeType == "" {
			mimeType = "text/html" // Default
		}
		c.Response().Header().Set("Content-Type", mimeType)
		c.Response().Header().Set("X-Proxy-Pass", "timeout-tricker")

		// Write headers
		// We may have to use some heuristics based on the request
		// to send the correct headers
		headersSent := false

		for {
			select {
			case <-done:
				// Upstream request is done!
				// Write out the original body
				fmt.Printf("remote request finished\n")
				// and header if we did not sent anything yet!
				if !headersSent {
					writer := c.Response().Writer
					for k, v := range recorder.Header() {
						writer.Header()[k] = v
					}
					c.Response().WriteHeader(recorder.Result().StatusCode)
				}
				_, err := c.Response().Writer.Write(recorder.Body.Bytes())
				return err
			case <-time.After(time.Duration(timeout) * time.Second):
				writer := c.Response().Writer
				_, _ = writer.Write([]byte(" "))
				if f, ok := writer.(http.Flusher); ok {
					fmt.Printf("flushing\n")
					f.Flush()
				}
				headersSent = true
			}
		}
	}
}