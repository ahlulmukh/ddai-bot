package captcha

import (
	"bytes"
	"ddai-bot/internal/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type CaptchaServices struct {
	sitekey           string
	pageUrl           string
	cfSolved          string
	antiCaptchaApiUrl string
}

type Config struct {
	CaptchaServices struct {
		CaptchaUsing      string   `json:"captchaUsing"`
		UrlPrivate        string   `json:"urlPrivate"`
		AntiCaptchaApikey []string `json:"antiCaptchaApikey"`
		Captcha2Apikey    []string `json:"captcha2Apikey"`
	} `json:"captchaServices"`
}

func NewCaptchaServices() *CaptchaServices {
	config := LoadConfig()
	return &CaptchaServices{
		sitekey:           "0x4AAAAAABdw7Ezbqw4v6Kr1",
		pageUrl:           "https://app.ddai.network/",
		cfSolved:          config.CaptchaServices.UrlPrivate + "/turnstiler",
		antiCaptchaApiUrl: "https://api.anti-captcha.com",
	}
}

func LoadConfig() *Config {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}

	return &config
}

func (cs *CaptchaServices) SolveCaptcha(currentNum, total int) (string, error) {
	config := LoadConfig()
	provider := config.CaptchaServices.CaptchaUsing

	switch provider {
	case "2captcha":
		return cs.solveCaptcha2(currentNum, total)
	case "antiCaptcha":
		return cs.antiCaptcha(currentNum, total)
	case "private":
		return cs.solvedPrivate(currentNum, total)
	default:
		utils.LogMessage(currentNum, total, "Invalid captcha provider.", "error")
		return "", fmt.Errorf("invalid captcha provider")
	}
}

func (cs *CaptchaServices) solvedPrivate(currentNum, total int) (string, error) {
	utils.LogMessage(currentNum, total, "Trying to solved captcha cloudflare...", "process")

	data := map[string]interface{}{
		"url":     "https://app.ddai.network/",
		"siteKey": "0x4AAAAAABdw7Ezbqw4v6Kr1",
		"mode":    "turnstile-min",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(cs.cfSolved, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	if token, ok := result["token"].(string); ok {
		utils.LogMessage(currentNum, total, "Captcha solved successfully!", "success")
		return token, nil
	}

	return "", fmt.Errorf("failed to get token from Cloudflare")
}

func (cs *CaptchaServices) antiCaptcha(currentNum, total int) (string, error) {
	utils.LogMessage(currentNum, total, "Trying solving captcha Turnstile...", "process")
	config := LoadConfig()
	apiKey := config.CaptchaServices.AntiCaptchaApikey[0]

	taskData := map[string]interface{}{
		"clientKey": apiKey,
		"task": map[string]interface{}{
			"type":       "TurnstileTaskProxyless",
			"websiteURL": cs.pageUrl,
			"websiteKey": cs.sitekey,
		},
		"softId": 0,
	}

	jsonData, err := json.Marshal(taskData)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(cs.antiCaptchaApiUrl+"/createTask", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var createTaskResp struct {
		ErrorId int    `json:"errorId"`
		TaskId  int    `json:"taskId"`
		Status  string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&createTaskResp); err != nil {
		return "", err
	}

	if createTaskResp.ErrorId != 0 || createTaskResp.TaskId == 0 {
		return "", fmt.Errorf("failed to create task")
	}

	utils.LogMessage(currentNum, total, fmt.Sprintf("Task created with ID: %d", createTaskResp.TaskId), "process")

	getTaskData := map[string]interface{}{
		"clientKey": apiKey,
		"taskId":    createTaskResp.TaskId,
	}

	var result string
	for i := 0; i < 10; i++ {
		time.Sleep(5 * time.Second)

		jsonData, err = json.Marshal(getTaskData)
		if err != nil {
			continue
		}

		resp, err = http.Post(cs.antiCaptchaApiUrl+"/getTaskResult", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			continue
		}

		var taskResult struct {
			ErrorId  int    `json:"errorId"`
			Status   string `json:"status"`
			Solution struct {
				Token string `json:"token"`
			} `json:"solution"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&taskResult); err != nil {
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		if taskResult.Status == "ready" {
			result = taskResult.Solution.Token
			utils.LogMessage(currentNum, total, "Captcha solved successfully!", "success")
			break
		}
	}

	if result == "" {
		return "", fmt.Errorf("failed to get captcha solution")
	}

	return result, nil
}

func (cs *CaptchaServices) solveCaptcha2(currentNum, total int) (string, error) {
	utils.LogMessage(currentNum, total, "Trying solving captcha Turnstile...", "process")
	config := LoadConfig()
	apiKey := config.CaptchaServices.Captcha2Apikey[0]

	client := &http.Client{Timeout: 120 * time.Second}
	reqData := map[string]string{
		"key":     apiKey,
		"method":  "turnstile",
		"sitekey": cs.sitekey,
		"pageurl": cs.pageUrl,
		"json":    "1",
	}

	form := url.Values{}
	for k, v := range reqData {
		form.Add(k, v)
	}

	resp, err := client.PostForm("http://2captcha.com/in.php", form)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var inResponse struct {
		Status  int    `json:"status"`
		Request string `json:"request"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&inResponse); err != nil {
		return "", err
	}

	if inResponse.Status != 1 {
		return "", fmt.Errorf("failed to submit captcha to 2captcha")
	}

	captchaID := inResponse.Request
	utils.LogMessage(currentNum, total, fmt.Sprintf("Captcha submitted, ID: %s", captchaID), "process")

	for i := 0; i < 20; i++ {
		time.Sleep(5 * time.Second)
		resp, err := client.Get(fmt.Sprintf("http://2captcha.com/res.php?key=%s&action=get&id=%s&json=1", apiKey, captchaID))
		if err != nil {
			continue
		}

		var resResponse struct {
			Status  int    `json:"status"`
			Request string `json:"request"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&resResponse); err != nil {
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		if resResponse.Status == 1 {
			utils.LogMessage(currentNum, total, "Captcha solved successfully!", "success")
			return resResponse.Request, nil
		}
	}

	return "", fmt.Errorf("failed to solve captcha with 2captcha")
}
