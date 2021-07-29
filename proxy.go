package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"

	http "github.com/useflyent/fhttp"

	"github.com/Carcraftz/CycleTLS/cycletls"
)

var client = cycletls.Init()

func main() {
	port := flag.String("port", "8082", "A port number (default 8082)")
	flag.Parse()
	addr := "127.0.0.1:" + string(*port)
	fmt.Println("Hosting a TLS API on http://" + addr)
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

	redirectVal := r.Header.Get("Poptls-Allowredirect")
	allowRedirect := true
	if redirectVal != "" {
		if redirectVal == "false" {
			fmt.Println(("redirects disabled"))
			allowRedirect = false
		}
	}
	if redirectVal != "" {
		r.Header.Del("Poptls-Allowredirect")
	}
	// Change JA3
	// Use Chrome JA3 by default
	//TODO: ADD MORE JA3 HASHES
	tlsClient := "771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-21,29-23-24,0"
	if strings.Contains(strings.ToLower(userAgent), "firefox") {
		tlsClient = "771,4865-4867-4866-49195-49199-52393-52392-49196-49200-49162-49161-49171-49172-156-157-47-53-10,0-23-65281-10-11-35-16-5-51-43-13-45-28-21,29-23-24-25-256-257,0"
	}

	// Forward body
	body, _ := ioutil.ReadAll(r.Body)

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
	u, err := url.Parse(endpoint)
	if err != nil {
		panic(err)
	}
	// Forward headers
	headers := make(map[string]string)
	for k := range r.Header {
		if k != "Content-Length" {
			headers[k] = r.Header.Get(k)
		}
	}
	headers["Host"] = u.Host
	// Set options
	opts := cycletls.Options{
		Body:          string(body), // ideally, body should be passed as bytes
		Ja3:           tlsClient,
		UserAgent:     userAgent,
		Headers:       headers,
		Method:        r.Method,
		Proxy:         proxy,
		AllowRedirect: allowRedirect,
	}
	// Perform request
	res, err := client.Do(endpoint, opts, r.Method)
	if err != nil {
		log.Println("Request failed:", err)
	}
	fmt.Println(res.Response.Status)
	// Forward response headers
	for k, v := range res.Response.Headers {
		// Do not forward the content length and encoding headers, as we will handle the decoding ourselves
		if k != "Content-Length" && k != "Content-Encoding" {
			w.Header().Set(k, v)
		}

	}
	// Forward respone status
	w.WriteHeader(res.Response.Status)

	// Write forwarded proxy request
	if _, err := fmt.Fprint(w, res.Response.Body); err != nil {
		log.Println("Error writing body:", err)
	}

}
