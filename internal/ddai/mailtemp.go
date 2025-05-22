package ddai

import (
	"ddai-bot/internal/utils"
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

type MailTemp struct {
	proxy      string
	currentNum int
	total      int
	httpClient *HTTPClient
}

func NewMailTemp(proxy string, currentNum, total int) *MailTemp {
	return &MailTemp{
		proxy:      proxy,
		currentNum: currentNum,
		total:      total,
		httpClient: NewHTTPClient(proxy, currentNum, total),
	}
}

func (m *MailTemp) MakeRequest(method, urlPath string) ([]byte, error) {
	return m.httpClient.MakeRequest(method, urlPath)
}

func (m *MailTemp) GetRandomDomain() (string, error) {
	utils.LogMessage(m.currentNum, m.total, "Trying to get a random domain...", "process")
	vowels := "aeiou"
	consonants := "bcdfghjklmnpqrstvwxyz"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomConsonant := string(consonants[r.Intn(len(consonants))])
	randomVowel := string(vowels[r.Intn(len(vowels))])
	keyword := randomConsonant + randomVowel

	responseBody, err := m.MakeRequest("GET", fmt.Sprintf("https://generator.email/search.php?key=%s", keyword))
	if err != nil {
		utils.LogMessage(m.currentNum, m.total, fmt.Sprintf("Error getting random domain: %v", err), "error")
		return "", err
	}

	if len(responseBody) == 0 {
		utils.LogMessage(m.currentNum, m.total, "No response from API", "error")
		return "", fmt.Errorf("empty response from API")
	}

	var domains []string
	if err := json.Unmarshal(responseBody, &domains); err != nil {
		utils.LogMessage(m.currentNum, m.total, fmt.Sprintf("Failed to parse domains: %v", err), "error")
		return "", err
	}

	var filteredDomains []string
	asciiPattern := regexp.MustCompile(`^[\x00-\x7F]*$`)
	for _, domain := range domains {
		if asciiPattern.MatchString(domain) {
			filteredDomains = append(filteredDomains, domain)
		}
	}

	if len(filteredDomains) == 0 {
		utils.LogMessage(m.currentNum, m.total, "Could not find valid domain", "error")
		return "", fmt.Errorf("no valid domains found")
	}

	selectedDomain := filteredDomains[r.Intn(len(filteredDomains))]
	utils.LogMessage(m.currentNum, m.total, fmt.Sprintf("Selected domain: %s", selectedDomain), "success")

	return selectedDomain, nil
}
