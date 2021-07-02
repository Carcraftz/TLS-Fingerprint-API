package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Danny-Dasilva/CycleTLS/cycletls"
)

const addr = "127.0.0.1:8082" // HOST:PORT

var client = cycletls.Init()

func main() {
	fmt.Println("Hosting TLS API on http://" + addr)
	http.HandleFunc("/", handleReq)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalln("Error starting the HTTP server:", err)
	}
}

func handleReq(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Ensure page URL header is provided
	pageURL := r.Header.Get("Poptls-Url")
	if pageURL == "" {
		http.Error(w, "ERROR: No Page URL Provided", http.StatusBadRequest)
		return
	}
	// Remove header to ignore later
	r.Header.Del("Poptls-Url")

	// Ensure user agent header is provided
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		http.Error(w, "ERROR: No User Agent Provided", http.StatusBadRequest)
		return
	}

	// Extract proxy from header
	proxy := r.Header.Get("Poptls-Proxy")
	if proxy != "" {
		r.Header.Del("Poptls-Proxy")
	}

	// Change JA3
	// Use Chrome JA3 by default
	tlsClient := "771,4865-4867-4866-49195-49199-52393-52392-49196-49200-49162-49161-49171-49172-51-57-47-53-10,0-23-65281-10-11-35-16-5-51-43-13-45-28-21,29-23-24-25-256-257,0"
	if strings.Contains(strings.ToLower(userAgent), "firefox") {
		tlsClient = "771,4865-4867-4866-49195-49199-52393-52392-49196-49200-49162-49161-49171-49172-156-157-47-53-10,0-23-65281-10-11-35-16-5-51-43-13-45-28-21,29-23-24-25-256-257,0"
	}

	// Forward body
	body, _ := ioutil.ReadAll(r.Body)

	// Forward headers
	headers := make(map[string]string)
	for k := range r.Header {
		headers[k] = r.Header.Get(k)
	}

	// Forward query params
	var addedQuery string
	for k, v := range r.URL.Query() {
		addedQuery += "&" + k + "=" + v[0]
	}

	endpoint := pageURL + "?" + addedQuery
	if strings.Contains(pageURL, "?") {
		endpoint = pageURL + addedQuery
	} else if addedQuery != "" {
		endpoint = pageURL + "?" + addedQuery[1:]
	}

	// Set options
	opts := cycletls.Options{
		Body:      string(body), // ideally, body should be passed as bytes
		Ja3:       tlsClient,
		UserAgent: userAgent,
		Headers:   headers,
		Method:    r.Method,
		Proxy:     proxy,
	}

	// Perform request
	res, err := client.Do(endpoint, opts, r.Method)
	if err != nil {
		log.Println("Request failed:", err)
	}

	// Forward response headers
	for k, v := range res.Response.Headers {
		// Do not forward the content length and encoding headers, as we will handle the content ourselves
		if k != "Content-Length" && k != "Content-Encoding" {
			continue
		}
		w.Header().Set(k, v)
	}
	// Forward respone status
	w.WriteHeader(res.Response.Status)

	// Decode response content
	encoding := res.Response.Headers["Content-Encoding"]
	decodedBody := decodeFile(encoding, []byte(res.Response.Body))

	// Write forwarded proxy request
	if _, err := fmt.Fprint(w, decodedBody); err != nil {
		log.Println("Error writing body:", err)
	}
}

// TODO: decode to bytes instead of string
// TODO: add other content encoding types besides gzip
func decodeFile(encoding string, body []byte) string {
	switch encoding {
	case "gzip":
		unz, err := gUnzipData(body)
		if err != nil {
			log.Println("INDEX 8")
			log.Println("Error occured:", err)
		}
		return string(unz)
	default:
		return string(body)
	}
}

// Not needed now, flush is also not required (see http.Error() implementation)
// func sendRes(w http.ResponseWriter, s string) {
// 	n, err := fmt.Fprint(w, s)
// 	if err != nil {
// 		log.Println("Error writing body:", err)
// 	}

// 	// Only flush if response writer implements http.Flusher and wrote more than 0 bytes
// 	if flush, ok := w.(http.Flusher); ok && n > 0 {
// 		flush.Flush()
// 	}
// }

func gUnzipData(data []byte) ([]byte, error) {
	// Unzip reader
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return []byte{}, err
	}

	// Copy from unzip reader to dst buffer
	dst := bytes.NewBuffer([]byte{})
	_, err = io.Copy(dst, r)
	return dst.Bytes(), err
}
