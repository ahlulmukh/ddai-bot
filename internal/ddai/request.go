package ddai

import (
	"bytes"
	"crypto/tls"
	"ddai-bot/internal/utils"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type HTTPClient struct {
	proxy      string
	currentNum int
	total      int
	client     *http.Client
}

func NewHTTPClient(proxy string, currentNum, total int) *HTTPClient {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err == nil {
			transport := &http.Transport{
				Proxy:           http.ProxyURL(proxyUrl),
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			client.Transport = transport
		} else {
			utils.LogMessage(currentNum, total, fmt.Sprintf("Invalid proxy URL: %v", err), "warning")
		}
	}

	return &HTTPClient{
		proxy:      proxy,
		currentNum: currentNum,
		total:      total,
		client:     client,
	}
}

func (h *HTTPClient) MakeRequest(method, urlPath string) ([]byte, error) {
	req, err := http.NewRequest(method, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	return body, nil
}

func (h *HTTPClient) MakeRequestWithBody(method, urlPath string, body []byte, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(method, urlPath, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return respBody, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	return respBody, nil
}
