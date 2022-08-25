package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
)

type transport struct {
	http.RoundTripper
}

func processRequest(request *http.Request) (*http.Request, error) {
	fmt.Print("\nREQUEST\n")
	fmt.Println(request.Method, request.URL)
	request.Header.Write(os.Stdout)

	if request.Body != nil {
		reqB, err := io.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		err = request.Body.Close()
		if err != nil {
			return nil, err
		}

		fmt.Print("\nREQUEST BODY\n")
		fmt.Println(string(reqB))
		fmt.Print("\n")

		request.Body = io.NopCloser(bytes.NewReader(reqB))
	}

	return request, nil
}

func processResponse(response *http.Response) (*http.Response, error) {
	respB, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = response.Body.Close()
	if err != nil {
		return nil, err
	}

	fmt.Print("\nRESPONSE\n")
	fmt.Println(response.Proto, response.Status)
	response.Header.Write(os.Stdout)
	fmt.Print("\nRESPONSE BODY\n")
	fmt.Println(string(respB))
	fmt.Print("\n")

	body := io.NopCloser(bytes.NewReader(respB))

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
