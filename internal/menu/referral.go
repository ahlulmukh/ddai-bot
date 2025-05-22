package menu

import (
	"ddai-bot/internal/ddai"
	"ddai-bot/internal/proxy"
	"ddai-bot/internal/utils"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

func (m *MenuHandler) RunReferralProgram() {
	fmt.Print("Enter referral code: ")
	refCode, _ := m.reader.ReadString('\n')
	refCode = strings.TrimSpace(refCode)

	count := m.getInput("How many accounts? ", 1)
	threads := m.getInput("Threads count? ", 1)

	m.startReferralProcess(refCode, count, threads)
}

func (m *MenuHandler) getInput(prompt string, min int) int {
	fmt.Print(prompt)
	input, _ := m.reader.ReadString('\n')
	val, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || val < min {
		utils.LogMessage(0, 0, "Invalid input", "error")
		return m.getInput(prompt, min)
	}
	return val
}

func (m *MenuHandler) startReferralProcess(refCode string, count, threads int) {
	proxy.LoadProxies()

	var wg sync.WaitGroup
	jobs := make(chan int, count)
	successCh := make(chan int, count)

	for w := 0; w < threads; w++ {
		wg.Add(1)
		go m.referralWorker(w, &wg, jobs, successCh, refCode, count)
	}

	for i := 0; i < count; i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	m.showResults(successCh, count)
}

func (m *MenuHandler) referralWorker(_ int, wg *sync.WaitGroup, jobs <-chan int,
	successCh chan<- int, refCode string, total int) {
	defer wg.Done()

	for idx := range jobs {
		proxy, err := proxy.GetRandomProxy(idx+1, total)
		if err != nil {
			utils.LogMessage(0, 0, fmt.Sprintf("Failed to get proxy for job %d: %v", idx+1, err), "error")
			continue
		}
		dd := ddai.NewDdaiReferral(refCode, proxy, idx+1, total)

		if err := dd.SingleProses(); err == nil {
			successCh <- 1
		}
	}
}

func (m *MenuHandler) showResults(successCh chan int, total int) {
	close(successCh)
	success := len(successCh)

	utils.LogMessage(0, 0, "Process completed!", "success")
	utils.LogMessage(0, 0,
		fmt.Sprintf("Success: %d/%d", success, total), "info")
	m.waitForEnter()
}
