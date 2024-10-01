package listener

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"sync"
)

var (
	httpClient *http.Client
)

func init() {
	httpClient = &http.Client{
		// Do not auto-follow redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

type responseResult struct {
	idx      int
	response response
}

type response struct {
	statusCode int
	url        string
	headers    map[string][]string
	body       io.ReadCloser
}

func (l *Listener) handleRequest(req *http.Request, resp http.ResponseWriter) {
	defer req.Body.Close()

	bodyInBytes, _ := io.ReadAll(req.Body)

	ch := make(chan responseResult)
	var wg sync.WaitGroup

	for idx, u := range l.upstreams {
		wg.Add(1)

		go func(idx int, upstream string) {
			defer wg.Done()

			res := l.proxyRequest(req, upstream, bodyInBytes)

			ch <- responseResult{
				idx:      idx,
				response: res,
			}
		}(idx, u)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for resResult := range ch {
		defer resResult.response.body.Close()

		if resResult.idx == 0 {
			handleResponseCloser(resResult.response.body, resp, resResult.response.statusCode, resResult.response.headers)
		}
	}
}

func handleResponseCloser(readCloser io.ReadCloser, resp http.ResponseWriter, statusCode int, headers map[string][]string) {
	defer readCloser.Close()

	for name, values := range headers {
		for _, value := range values {
			resp.Header().Add(name, value)
		}
	}
	resp.WriteHeader(statusCode)

	io.Copy(resp, readCloser)
}

func (l *Listener) proxyRequest(req *http.Request, upstream string, body []byte) response {
	newUrl, _ := url.Parse(upstream)

	concatPath := req.URL.Path
	if concatPath != "" {
		newUrl.Path = newUrl.Path + concatPath
	}

	newQuery := newUrl.Query()

	for name, values := range req.URL.Query() {
		for _, value := range values {
			newQuery.Add(name, value)
		}
	}

	newUrl.RawQuery = newQuery.Encode()

	proxiedReq, _ := http.NewRequest(req.Method, newUrl.String(), bytes.NewBuffer(body))

	for name, values := range req.Header {
		if name == "X-Forwarded-Host" {
			continue
		}

		for _, value := range values {
			proxiedReq.Header.Add(name, value)
		}
	}

	if l.rewriteHost {
		proxiedReq.Host = l.originHostName
	}

	proxiedResp, _ := httpClient.Do(proxiedReq)

	headers := make(map[string][]string)

	// Copy response headers
	for name, values := range proxiedResp.Header {
		for _, value := range values {

			currentValues, exists := headers[name]
			if !exists {
				currentValues = make([]string, 0)
			}
			currentValues = append(currentValues, value)
			headers[name] = currentValues
		}
	}

	if headers == nil {
		headers = make(map[string][]string)
	}

	return response{
		statusCode: proxiedResp.StatusCode,
		url:        proxiedReq.URL.String(),
		headers:    headers,
		body:       proxiedResp.Body,
	}
}
