package listener

import (
	"bytes"
	"crypto/tls"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
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

func (l *Listener) handleRequest(req *http.Request, resp http.ResponseWriter) {
	result := newHandleResult(req)
	RequestHandlingChanged(result)

	upstreams := upstream.GetUpstreams()

	for idx, u := range upstreams {
		e := &event.HandleBodyAndHeaders{
			Status: "Pending",
		}
		result.Responses[idx] = e
		RequestHandlingChanged(result)

		statusCode, _, headers, body := l.proxyRequest(req, u)

		if headers == nil {
			headers = make(map[string][]string)
		}

		RequestHandlingChanged(result)
		var bodyBytes []byte

		if idx == 0 {
			bodyBytes, _ = handleResponseCloser(body, resp, statusCode, headers)
		} else {
			bodyBytes, _ = retrieveResponseBytes(body)
		}

		e.Body = bodyBytes
		e.Headers = headers
		e.Status = "finished"
		RequestHandlingChanged(result)
	}
}

func newHandleResult(req *http.Request) event.HandleResult {
	return event.HandleResult{
		ID:     uuid.New().String(),
		Method: req.Method,
		URL:    req.URL.String(),
		Request: event.HandleBodyAndHeaders{
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

func (l *Listener) proxyRequest(req *http.Request, u upstream.Upstream) (int, string, map[string][]string, io.ReadCloser) {
	defer req.Body.Close()

	newURL, _ := url.Parse(u.Url)

	concatPath := req.URL.Path
	if concatPath != "" {
		newURL.Path = newURL.Path + concatPath
	}

	bodyInBytes, _ := ioutil.ReadAll(req.Body)

	proxiedReq, _ := http.NewRequest(req.Method, newURL.String(), bytes.NewBuffer(bodyInBytes))

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

	return proxiedResp.StatusCode, proxiedReq.URL.String(), headers, proxiedResp.Body
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
