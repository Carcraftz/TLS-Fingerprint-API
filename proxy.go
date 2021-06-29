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

const (
	Host = "127.0.0.1"
	Port = "8082"
)

var client = cycletls.Init()

func main() {

	fmt.Println("Hosting TLS API on http://" + Host + ":" + Port)
	http.HandleFunc("/", handleReq)
	err := http.ListenAndServe(Host+":"+Port, nil)
	if err != nil {
		log.Fatal("Error Starting the HTTP Server : ", err)
		return
	}

}
func handleReq(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	fetchmode := "GET"
	switch r.Method {
		case http.MethodPost:
			fetchmode = "POST"
		case http.MethodPut:
			fetchmode = "PUT"
		case http.MethodDelete:
			fetchmode = "DELETE"
		case http.MethodOptions:
			fetchmode = "OPTIONS"
		default:
			fetchmode = "GET"
	}
	fmt.Println(fetchmode)

	useproxy := len(r.Header["Poptls-Proxy"]) > 0
	//make sure required info is provided
	if (len(r.Header["Poptls-Url"]) == 0) {
		fmt.Println("NO PAGE URL")
		sendRes(w, "ERROR: No Page URL Provided")

	}
	if (len(r.Header["User-Agent"]) == 0) {
		sendRes(w, "ERROR: No User Agent Provided")
	}
	cond1 := (len(r.Header["Poptls-Url"]) > 0)
	cond2 := (len(r.Header["User-Agent"]) > 0)
	proceed := cond1 && cond2
	if proceed {
		pageurl := r.Header["Poptls-Url"][0]
		uagent := r.Header["User-Agent"][0]
		var tlsclient string

		//change ja3
		if strings.Contains(strings.ToLower(uagent), "chrome") {
			tlsclient = "771,4865-4867-4866-49195-49199-52393-52392-49196-49200-49162-49161-49171-49172-51-57-47-53-10,0-23-65281-10-11-35-16-5-51-43-13-45-28-21,29-23-24-25-256-257,0"
		} else if strings.Contains(strings.ToLower(uagent), "firefox") {
			tlsclient = "771,4865-4867-4866-49195-49199-52393-52392-49196-49200-49162-49161-49171-49172-156-157-47-53-10,0-23-65281-10-11-35-16-5-51-43-13-45-28-21,29-23-24-25-256-257,0"
		} else {
			tlsclient = "771,4865-4867-4866-49195-49199-52393-52392-49196-49200-49162-49161-49171-49172-51-57-47-53-10,0-23-65281-10-11-35-16-5-51-43-13-45-28-21,29-23-24-25-256-257,0"
		}

		//forward body
		body, err := ioutil.ReadAll(r.Body)
		bodyString := string(body)

		//forward headers
		headerMap := make(map[string]string)
		for name, values := range r.Header {
			// Loop over all values for the name.
			for _, value := range values {
				if name != "Poptls-Url" && name != "Poptls-Proxy" {
					headerMap[name] = value
				}

			}
		}
		//forward query params
		addedquery := ""
		for k, v := range r.URL.Query() {
			addedquery = addedquery + "&" + k + "=" + v[0]
		}
		endpoint := ""
		if strings.Contains(pageurl, "?") {
			endpoint = pageurl + addedquery
		} else if addedquery != "" {
			endpoint = pageurl + "?" + addedquery[1:]
		} else {
			endpoint = pageurl + "?" + addedquery
		}

		opts := cycletls.Options{
			Body:      bodyString,
			Ja3:       tlsclient,
			UserAgent: uagent,
			Headers:   headerMap,
			Method:    fetchmode,
		}
		//use proxy
		if useproxy {
			opts.Proxy = r.Header["Poptls-Proxy"][0]
		}
		response, err := client.Do(endpoint, opts, fetchmode)
		if err != nil {
			log.Print("Request Failed: " + err.Error())
		}
		for k, v := range response.Response.Headers {
			//we want go to handle decoding the files, so don't forward the content encoding and length headers
			if k != "Content-Length" && k != "Content-Encoding" {
				w.Header().Set(k, v)
			}
		}
		encoding := response.Response.Headers["Content-Encoding"]
		bodyStr := decodeFile(encoding, []byte(response.Response.Body))
		w.WriteHeader(response.Response.Status)
		sendRes(w, bodyStr)
		//fmt.Println(bodyStr)

	}
}

//TODO: add other content encoding types besides gzip
func decodeFile(encoding string, body []byte) string {
	if encoding == "gzip" {
		unz, err := gUnzipData(body)
		if err != nil {
			log.Println("INDEX 8")
			defer func() {
				if err := recover(); err != nil {
					log.Println("panic occurred:", err)
				}
			}()
			panic(err)
		}
		return string(unz)
	} else if encoding == "add more types here" {
		return string(body)
	} else {
		return string(body)
	}
	return string(body)
}
func sendRes(w http.ResponseWriter, s string) {
	n, err := fmt.Fprintf(w, s)
	if err != nil {
		fmt.Println(err)
	}
	n = n + 1
	w.(http.Flusher).Flush()

}
func gUnzipData(data []byte) (resData []byte, err error) {
	b := bytes.NewBuffer(data)

	var r io.Reader
	r, err = gzip.NewReader(b)
	if err != nil {
		return
	}

	var resB bytes.Buffer
	_, err = resB.ReadFrom(r)
	if err != nil {
		return
	}

	resData = resB.Bytes()

	return
}
