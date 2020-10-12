package utility

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// Testing ...
var (
	MockDo             func(req *http.Request) (*http.Response, error)
	IsMock             bool
	HTTPRequestTimeout = time.Duration(5) * time.Second
)

func do(req *http.Request) (*http.Response, error) {
	return MockDo(req)
}

//SendRequest sent the all http api request
func SendRequest(method string, url string, headers map[string]string, body interface{}) ([]byte, http.Header, int, error) {

	return send(method, url, headers, body, false, nil, "", "")
}

//SendRequestWithSSL ...
func SendRequestWithSSL(method string, url string, headers map[string]string, body interface{}) ([]byte, http.Header, int, error) {
	return send(method, url, headers, body, true, nil, "", "")
}

//SendRequestWithInsecure ...
func SendRequestWithInsecure(method string, url string, headers map[string]string, body interface{}, certFile string) ([]byte, http.Header, int, error) {
	return send(method, url, headers, body, true, &certFile, "", "")
}

// SendRequestWithBasicAuth ...
func SendRequestWithBasicAuth(method string, url string, headers map[string]string, body interface{}, username string, password string) ([]byte, http.Header, int, error) {
	return send(method, url, headers, body, false, nil, username, password)
}

func send(method string, url string, headers map[string]string, body interface{}, withSSL bool, certFile *string, username string, password string) ([]byte, http.Header, int, error) {
	request, err := bundleRequest(method, url, headers, body)
	if err != nil {
		return nil, nil, -1, err
	}

	if username != "" && password != "" {
		request.SetBasicAuth(username, password)
	}

	resp, err := bundleClient(request, false, nil)

	if err != nil {
		return nil, nil, -1, err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, resp.StatusCode, err
	}

	return b, resp.Header, resp.StatusCode, err
}

func bundleClient(req *http.Request, withSSL bool, certFile *string) (resp *http.Response, err error) {
	var client *http.Client
	if withSSL {
		var trans *http.Transport
		if certFile != nil {
			trans = withX509File(certFile)
		} else {
			trans = withInsecureVierify()
		}
		client = &http.Client{
			Transport: trans,
			Timeout:   HTTPRequestTimeout,
		}
	} else {
		client = &http.Client{
			Timeout: HTTPRequestTimeout,
		}
	}

	if IsMock {
		resp, err = do(req)
	} else {
		resp, err = client.Do(req)
	}
	return
}

func withX509File(certFile *string) *http.Transport {
	f, err := ioutil.ReadFile(*certFile)
	if err != nil {
		log.Fatal(err)
	}
	cert := x509.NewCertPool()
	if ok := cert.AppendCertsFromPEM(f); ok {
		return &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: cert},
		}
	}
	return nil
}

func withInsecureVierify() *http.Transport {
	return &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

func bundleRequest(method string, url string, headers map[string]string, body interface{}) (request *http.Request, err error) {
	var b io.Reader
	if body != nil {
		switch body.(type) {
		case *bytes.Buffer:
			b = bytes.NewBuffer([]byte(fmt.Sprintf("%v", body)))
		default:
			if b, err = convertBodyType(&body); err != nil {
				return nil, err
			}
		}
	}
	if request, err = http.NewRequest(method, url, b); err != nil {
		return nil, err
	}
	if headers != nil {
		for key, val := range headers {
			request.Header.Add(key, val)
		}
	}
	return request, nil
}

func convertBodyType(body *interface{}) (*strings.Reader, error) {
	if body == nil {
		return nil, nil
	}
	newBodyType, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}
	return strings.NewReader(string(newBodyType)), nil
}
