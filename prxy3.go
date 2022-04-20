package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

type transport struct {
	http.RoundTripper
}

func processRequest(request *http.Request) (*http.Request, error) {
	fmt.Print("\n\nREQUEST\n")
	fmt.Println("url: ", request.URL)
	fmt.Println("method: ", request.Method)

	if request.Body != nil {
		reqB, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		err = request.Body.Close()
		if err != nil {
			return nil, err
		}

		fmt.Println("begin request body")
		fmt.Println(string(reqB))
		fmt.Println("end request body")

		request.Body = ioutil.NopCloser(bytes.NewReader(reqB))
	}

	return request, nil
}

func processResponse(response *http.Response) (*http.Response, error) {
	respB, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = response.Body.Close()
	if err != nil {
		return nil, err
	}

	fmt.Print("\n\nRESPONSE\n")
	fmt.Println("status code: ", response.StatusCode)
	fmt.Println("begin response body")
	fmt.Println(string(respB))
	fmt.Println("end response body")

	body := ioutil.NopCloser(bytes.NewReader(respB))

	response.Body = body
	response.ContentLength = int64(len(respB))
	response.Header.Set("Content-Length", strconv.Itoa(len(respB)))
	return response, nil
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	req, err = processRequest(req)
	if err != nil {
		return nil, err
	}

	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	resp, err = processResponse(resp)
	return resp, err
}

func main() {
	dst := flag.String("d", "localhost:8080", "the address of the target")
	prx := flag.String("p", "localhost:9090", "the address of the proxy")
	flag.Parse()

	reverseUrlStr := "http://" + *dst + "/"
	reverseUrl, _ := url.Parse(reverseUrlStr)

	handler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			// log.Println(r.URL)
			r.Host = reverseUrl.Host
			//w.Header().Set("X-PRXY", "rad")
			p.ServeHTTP(w, r)
		}
	}

	proxy := httputil.NewSingleHostReverseProxy(reverseUrl)
	proxy.Transport = &transport{http.DefaultTransport}

	http.HandleFunc("/", handler(proxy))

	err := http.ListenAndServe(*prx, nil)
	if err != nil {
		panic(err)
	}
}
