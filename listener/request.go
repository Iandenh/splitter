package listener

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/url"
	"splitter/event"
	"splitter/upstream"
	"sync"
)

var (
	httpClient      *http.Client
	httpClientMutex sync.Mutex
)

const (
	bufferSize        = 4 * 1024 * 1024        // 4 KBs
	maxStoredBodySize = 5 * 1024 * 1024 * 1024 // 5 MBs
)

type responseResult struct {
	idx      int
	e        *event.HandleBodyAndHeaders
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

	result := newHandleResult(req, bodyInBytes)
	RequestHandlingChanged(result)

	upstreams := upstream.GetUpstreams()

	ch := make(chan responseResult)
	var wg sync.WaitGroup

	for idx, u := range upstreams {
		wg.Add(1)

		go func(idx int, u upstream.Upstream) {
			defer wg.Done()
			e := &event.HandleBodyAndHeaders{
				Status: "Pending",
			}
			result.Responses[idx] = e

			fmt.Println("Starting %s", u.Url)

			res := l.proxyRequest(req, u, bodyInBytes)

			fmt.Println("Done %s", u.Url)
			ch <- responseResult{
				idx:      idx,
				e:        e,
				response: res,
			}

		}(idx, u)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for resResult := range ch {
		var bodyBytes []byte

		if resResult.idx == 0 {
			bodyBytes, _ = handleResponseCloser(resResult.response.body, resp, resResult.response.statusCode, resResult.response.headers)
			result.Response = resResult.e
		} else {
			bodyBytes, _ = retrieveResponseBytes(resResult.response.body)
		}

		resResult.e.Body = bodyBytes
		resResult.e.Headers = resResult.response.headers
		resResult.e.Status = "finished"
		RequestHandlingChanged(result)
	}
}

func newHandleResult(req *http.Request, body []byte) event.HandleResult {
	return event.HandleResult{
		ID:     uuid.New().String(),
		Method: req.Method,
		URL:    req.URL.String(),
		Request: &event.HandleBodyAndHeaders{
			Body:    body,
			Headers: req.Header,
		},
		Responses: make(map[int]*event.HandleBodyAndHeaders),
	}
}

func handleResponseCloser(readCloser io.ReadCloser, resp http.ResponseWriter, statusCode int, headers map[string][]string) ([]byte, error) {
	defer readCloser.Close()

	for name, values := range headers {
		for _, value := range values {
			resp.Header().Add(name, value)
		}
	}
	resp.WriteHeader(statusCode)

	responseBytes := make([]byte, 0)
	buffer := make([]byte, bufferSize)
	for {
		bytesRead, readError := readCloser.Read(buffer)

		if readError != nil && readError != io.EOF {
			return nil, readError
		}

		if bytesRead == 0 {
			break
		}

		if len(responseBytes) < maxStoredBodySize {
			responseBytes = append(responseBytes, buffer[:bytesRead]...)
		}
		resp.Write(buffer[:bytesRead])
	}

	return responseBytes, nil
}

func retrieveResponseBytes(readCloser io.ReadCloser) ([]byte, error) {
	defer readCloser.Close()

	responseBytes := make([]byte, 0)
	buffer := make([]byte, bufferSize)
	for {
		bytesRead, readError := readCloser.Read(buffer)

		if readError != nil && readError != io.EOF {
			return nil, readError
		}

		if bytesRead == 0 {
			break
		}

		if len(responseBytes) < maxStoredBodySize {
			responseBytes = append(responseBytes, buffer[:bytesRead]...)
		}
	}

	return responseBytes, nil
}

func (l *Listener) proxyRequest(req *http.Request, u upstream.Upstream, body []byte) response {
	newUrl, _ := url.Parse(u.Url)

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
		for _, value := range values {
			proxiedReq.Header.Add(name, value)
		}
	}

	if l.RewriteHost {
		proxiedReq.Host = l.OriginHostName
	}

	proxiedResp, _ := createHTTPClient().Do(proxiedReq)

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

func createHTTPClient() *http.Client {
	httpClientMutex.Lock()
	defer httpClientMutex.Unlock()

	if httpClient == nil {
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
	return httpClient
}
