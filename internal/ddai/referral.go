package ddai

import (
	"ddai-bot/internal/captcha"
	"ddai-bot/internal/utils"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

const (
	retryCount        = 3
	retryDelay        = 3 * time.Second
	maxWaitTime       = 60 * time.Second
	balanceCheckDelay = 5 * time.Second
)

type RegisterResponse struct {
	Status string `json:"status"`
	Data   struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		User         struct {
			ID          int    `json:"_id"`
			Email       string `json:"email"`
			Username    string `json:"username"`
			JoinDate    string `json:"joinDate"`
			Rank        string `json:"rank"`
			RefCode     string `json:"refCode"`
			RefBy       string `json:"refBy"`
			RefCount    int    `json:"refCount"`
			Requests    int    `json:"requests"`
			RequestRate int    `json:"requestRate"`
			Points      int    `json:"points"`
		} `json:"user"`
	} `json:"data"`
	Error map[string]interface{} `json:"error"`
}

type LoginResponse struct {
	Status string `json:"status"`
	Data   struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		User         struct {
			ID          int    `json:"_id"`
			Email       string `json:"email"`
			Username    string `json:"username"`
			JoinDate    string `json:"joinDate"`
			Rank        string `json:"rank"`
			RefCode     string `json:"refCode"`
			RefBy       string `json:"refBy"`
			RefCount    int    `json:"refCount"`
			Requests    int    `json:"requests"`
			RequestRate int    `json:"requestRate"`
			Points      int    `json:"points"`
		} `json:"user"`
	} `json:"data"`
	Error map[string]interface{} `json:"error"`
}

type MissionsResponse struct {
	Status string `json:"status"`
	Data   struct {
		Missions []struct {
			ID          string `json:"_id"`
			Title       string `json:"title"`
			Type        int    `json:"type"`
			Description string `json:"description"`
			Link        string `json:"link"`
			Order       int    `json:"order"`
			Rewards     struct {
				Requests int `json:"requests"`
			} `json:"rewards"`
			Status string `json:"status"`
		} `json:"missions"`
	} `json:"data"`
	Error map[string]interface{} `json:"error"`
}

type ClaimResponse struct {
	Status string `json:"status"`
	Data   struct {
		Claimed   bool   `json:"claimed"`
		MissionID string `json:"missionId"`
		Rewards   struct {
			Sol      int `json:"sol"`
			Requests int `json:"requests"`
		} `json:"rewards"`
	} `json:"data"`
	Error map[string]interface{} `json:"error"`
}

type Mission struct {
	ID    string
	Title string
}

type ddaiReferral struct {
	proxy      string
	mainWallet string
	currentNum int
	total      int
	captcha    *captcha.CaptchaServices
	httpClient *HTTPClient
	mailTemp   *MailTemp
}

func NewDdaiReferral(mainWallet, proxy string, currentNum, total int) *ddaiReferral {
	return &ddaiReferral{
		proxy:      proxy,
		mainWallet: mainWallet,
		currentNum: currentNum,
		total:      total,
		captcha:    captcha.NewCaptchaServices(),
		httpClient: NewHTTPClient(proxy, currentNum, total),
		mailTemp:   NewMailTemp(proxy, currentNum, total),
	}
}

func (m *ddaiReferral) SingleProses() error {
	for attempt := 1; attempt <= retryCount; attempt++ {
		utils.LogMessage(m.currentNum, m.total, fmt.Sprintf("Attempt %d/%d", attempt, retryCount), "process")

		token, err := m.captcha.SolveCaptcha(m.currentNum, m.total)
		if err != nil {
			utils.LogMessage(m.currentNum, m.total, fmt.Sprintf("Failed to solve captcha: %v", err), "error")
			continue
		}

		// domain, err := m.mailTemp.GetRandomDomain()
		// if err != nil {
		// 	utils.LogMessage(m.currentNum, m.total, fmt.Sprintf("Failed to get random domain, using default: %v", err), "warning")
		// 	domain = "gmail.com"
		// }

		email := utils.GenerateEmail()
		password := utils.GeneratePassword()

		username := utils.GenerateUsername()
		err = m.registerAccount(email, username, password, token, m.mainWallet)
		if err != nil {
			utils.LogMessage(m.currentNum, m.total, fmt.Sprintf("%v", err), "warning")
			time.Sleep(retryDelay)
			continue
		}

		token, err = m.captcha.SolveCaptcha(m.currentNum, m.total)
		if err != nil {
			utils.LogMessage(m.currentNum, m.total, fmt.Sprintf("Failed to solve captcha: %v", err), "error")
			continue
		}

		accessToken, err := m.loginAccount(username, password, token)
		if err != nil {
			utils.LogMessage(m.currentNum, m.total, fmt.Sprintf("%v", err), "warning")
			return err
		}

		taskList, err := m.getUserTask(accessToken)
		if err != nil {
			utils.LogMessage(m.currentNum, m.total, fmt.Sprintf("Failed to get tasks: %v", err), "warning")
			continue
		}

		allTasksSuccess := true
		for _, task := range taskList {
			err := m.claimTask(accessToken, task)
			if err != nil {
				utils.LogMessage(m.currentNum, m.total, fmt.Sprintf("Failed to claim task %s: %v", task["name"], err), "warning")
				allTasksSuccess = false
			}
			time.Sleep(1 * time.Second)
		}

		anyTaskSuccess := len(taskList) > 0 && (allTasksSuccess || !allTasksSuccess)

		utils.LogMessage(m.currentNum, m.total, fmt.Sprintf("Tasks found: %d, Tasks claimed: %v", len(taskList), anyTaskSuccess), "info")

		if err := utils.SaveAccountToFile(email, password); err != nil {
			utils.LogMessage(m.currentNum, m.total, fmt.Sprintf("Failed to save account: %v", err), "warning")
		} else {
			utils.LogMessage(m.currentNum, m.total, "Account saved to accounts.txt after successful registration", "success")
		}

		return nil
	}

	return fmt.Errorf("failed after %d attempts", retryCount)
}

func (m *ddaiReferral) registerAccount(email string, username string, password string, token string, referral string) error {
	utils.LogMessage(m.currentNum, m.total, fmt.Sprintf("Registering account with email: %s", email), "process")
	payload := map[string]string{
		"email":        email,
		"username":     username,
		"password":     password,
		"refCode":      referral,
		"captchaToken": token,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	body, err := m.httpClient.MakeRequestWithBody("POST", "https://auth.ddai.space/register", jsonData, headers)
	if err != nil {
		return fmt.Errorf("registration request failed: %v", err)
	}

	var response RegisterResponse

	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to decode response: %v (body: %s)", err, string(body))
	}

	if response.Status == "success" {
		utils.LogMessage(m.currentNum, m.total, "Successfully registered account", "success")
		return nil
	}

	errorMsg := response.Error["message"]
	if errorMsg == nil {
		errorMsg = response.Error
	}
	return fmt.Errorf("registration failed: %v", errorMsg)
}

func (m *ddaiReferral) loginAccount(username string, password string, captcha string) (string, error) {
	utils.LogMessage(m.currentNum, m.total, fmt.Sprintf("Logging in with username: %s", username), "process")
	payload := map[string]string{
		"username":     username,
		"password":     password,
		"captchaToken": captcha,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %v", err)
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	body, err := m.httpClient.MakeRequestWithBody("POST", "https://auth.ddai.space/login", jsonData, headers)
	if err != nil {
		return "", fmt.Errorf("login request failed: %v", err)
	}

	var response LoginResponse

	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to decode response: %v (body: %s)", err, string(body))
	}

	if response.Status == "success" {
		utils.LogMessage(m.currentNum, m.total, "Successfully logged in", "success")
		return response.Data.AccessToken, nil
	}

	errorMsg := response.Error["message"]
	if errorMsg == nil {
		errorMsg = response.Error
	}
	return "", fmt.Errorf("login failed: %v", errorMsg)
}

func (m *ddaiReferral) getUserTask(accessToken string) ([]map[string]string, error) {
	utils.LogMessage(m.currentNum, m.total, "Fetching user tasks...", "process")

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + accessToken,
	}

	body, err := m.httpClient.MakeRequestWithBody("GET", "https://auth.ddai.space/missions", nil, headers)
	if err != nil {
		return nil, fmt.Errorf("get tasks request failed: %v", err)
	}

	var response MissionsResponse

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v (body: %s)", err, string(body))
	}

	utils.LogMessage(m.currentNum, m.total, fmt.Sprintf("Found %d missions", len(response.Data.Missions)), "info")

	var tasks []map[string]string
	for _, task := range response.Data.Missions {
		if task.Status == "PENDING" || task.Status == "idle" || task.Status == "pending" {
			if containsIgnoreCase(task.Title, "invite") {
				continue
			}

			taskInfo := map[string]string{
				"id":   task.ID,
				"name": task.Title,
			}
			tasks = append(tasks, taskInfo)
		}
	}

	return tasks, nil
}

func (m *ddaiReferral) claimTask(accessToken string, task map[string]string) error {
	utils.LogMessage(m.currentNum, m.total, fmt.Sprintf("Claiming task: %s (ID: %s)", task["name"], task["id"]), "process")

	url := fmt.Sprintf("https://auth.ddai.space/missions/claim/%s", task["id"])

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + accessToken,
	}

	body, err := m.httpClient.MakeRequestWithBody("POST", url, nil, headers)
	if err != nil {
		return fmt.Errorf("claim task request failed: %v", err)
	}

	var result ClaimResponse

	if err := json.Unmarshal(body, &result); err != nil {
		utils.LogMessage(m.currentNum, m.total, fmt.Sprintf("Raw response: %s", string(body)), "warning")
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if result.Status != "success" {
		errorMsg := "unknown error"
		if result.Error != nil {
			if msg, ok := result.Error["message"]; ok {
				errorMsg = fmt.Sprintf("%v", msg)
			}
		}
		utils.LogMessage(m.currentNum, m.total, fmt.Sprintf("Claim failed: %s", errorMsg), "warning")
		return fmt.Errorf("claim task failed: %s", errorMsg)
	}

	utils.LogMessage(m.currentNum, m.total, fmt.Sprintf("Successfully claimed task: %s with rewards: %d requests",
		task["name"], result.Data.Rewards.Requests), "success")
	return nil
}
