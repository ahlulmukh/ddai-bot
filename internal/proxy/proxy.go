package proxy

import (
	"bufio"
	"crypto/tls"
	"ddai-bot/internal/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	proxyList    []string
	lastUsedIP   string
	proxyEnabled bool
)

type IPResponse struct {
	IP string `json:"ip"`
}

func LoadProxies() bool {
	file, err := os.Open("proxy.txt")
	if err != nil {
		utils.LogMessage(0, 0, "Running without proxy: "+err.Error(), "warning")
		proxyEnabled = false
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxy := strings.TrimSpace(scanner.Text())
		if proxy != "" {
			if !strings.Contains(proxy, "://") {
				proxy = "http://" + proxy
			}
			proxyList = append(proxyList, proxy)
		}
	}

	if len(proxyList) == 0 {
		utils.LogMessage(0, 0, "No valid proxies found, running without proxy", "warning")
		proxyEnabled = false
		return false
	}

	proxyEnabled = true
	utils.LogMessage(0, 0, fmt.Sprintf("Loaded %d proxies", len(proxyList)), "success")
	return true
}

func getProxyTransport(proxyUrl string) (*http.Transport, error) {
	if proxyUrl == "" {
		return &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}, nil
	}

	proxy, err := url.Parse(proxyUrl)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(strings.ToLower(proxy.Scheme), "socks") {
		dialer := &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}
		return &http.Transport{
			Dial:                dialer.Dial,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			Proxy:               http.ProxyURL(proxy),
			DisableKeepAlives:   true,
			MaxIdleConnsPerHost: -1,
		}, nil
	}

	return &http.Transport{
		Proxy:               http.ProxyURL(proxy),
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives:   true,
		MaxIdleConnsPerHost: -1,
	}, nil
}

func CheckIP(currentNum, total int) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	if proxyEnabled {
		proxy := proxyList[currentNum%len(proxyList)]
		transport, err := getProxyTransport(proxy)
		if err != nil {
			return "", fmt.Errorf("proxy error: %v", err)
		}
		client.Transport = transport
	}

	resp, err := client.Get("https://api.ipify.org?format=json")
	if err != nil {
		return "", fmt.Errorf("IP check failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var ipResp IPResponse
	if err := json.Unmarshal(body, &ipResp); err != nil {
		return "", err
	}

	if ipResp.IP != lastUsedIP {
		utils.LogMessage(currentNum, total,
			fmt.Sprintf("New IP: %s (Proxy: %v)", ipResp.IP, proxyEnabled), "info")
		lastUsedIP = ipResp.IP
	}

	return ipResp.IP, nil
}

func GetRandomProxy(currentNum, total int) (string, error) {
	if !proxyEnabled || len(proxyList) == 0 {
		_, err := CheckIP(currentNum, total)
		if err != nil {
			return "", err
		}
		return "", nil
	}

	proxy := proxyList[currentNum%len(proxyList)]
	_, err := CheckIP(currentNum, total)
	if err != nil {
		return "", fmt.Errorf("proxy %s failed: %v", proxy, err)
	}

	return proxy, nil
}
