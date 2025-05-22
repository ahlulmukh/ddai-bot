package utils

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/fatih/color"
)

func LogMessage(currentNum, total int, message, messageType string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	accountStatus := ""
	if currentNum > 0 && total > 0 {
		accountStatus = fmt.Sprintf("[%d/%d] ", currentNum, total)
	}

	var colorPrinter *color.Color
	var symbol string

	switch messageType {
	case "info":
		colorPrinter = color.New(color.FgCyan)
		symbol = "[i]"
	case "success":
		colorPrinter = color.New(color.FgGreen)
		symbol = "[✓]"
	case "error":
		colorPrinter = color.New(color.FgRed)
		symbol = "[-]"
	case "warning":
		colorPrinter = color.New(color.FgYellow)
		symbol = "[!]"
	case "process":
		colorPrinter = color.New(color.FgHiCyan)
		symbol = "[>]"
	default:
		colorPrinter = color.New(color.Reset)
		symbol = "[*]"
	}

	logText := fmt.Sprintf("%s %s", symbol, message)
	fmt.Printf("[%s] %s", timestamp, accountStatus)
	colorPrinter.Println(logText)
}

func GeneratePassword() string {
	rand.Seed(time.Now().UnixNano())

	firstLetter := string(rune(rand.Intn(26) + 65))

	otherLetters := make([]rune, 4)
	for i := range otherLetters {
		otherLetters[i] = rune(rand.Intn(26) + 97)
	}

	numbers := make([]rune, 3)
	for i := range numbers {
		numbers[i] = rune(rand.Intn(10) + 48)
	}

	return fmt.Sprintf("%s%s@%s!", firstLetter, string(otherLetters), string(numbers))
}

func GenerateUsername() string {
	rand.Seed(time.Now().UnixNano())

	firstnameLen := rand.Intn(4) + 3
	firstnameChars := make([]rune, firstnameLen)
	for i := range firstnameChars {
		firstnameChars[i] = rune(rand.Intn(26) + 97)
	}
	firstname := string(firstnameChars)

	lastnameLen := rand.Intn(4) + 3
	lastnameChars := make([]rune, lastnameLen)
	for i := range lastnameChars {
		lastnameChars[i] = rune(rand.Intn(26) + 97)
	}
	lastname := string(lastnameChars)

	return fmt.Sprintf("%s%s", firstname, lastname)
}

func GenerateEmail() string {
	rand.Seed(time.Now().UnixNano())

	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	length := rand.Intn(5) + 8

	email := make([]byte, length)
	for i := range email {
		email[i] = charset[rand.Intn(len(charset))]
	}

	return string(email) + "@gmail.com"
}

func GenerateEmailTemp(domain string) string {
	firstnameLen := rand.Intn(4) + 3
	firstnameChars := make([]rune, firstnameLen)
	for i := range firstnameChars {
		firstnameChars[i] = rune(rand.Intn(26) + 97)
	}
	firstname := string(firstnameChars)

	lastnameLen := rand.Intn(4) + 3
	lastnameChars := make([]rune, lastnameLen)
	for i := range lastnameChars {
		lastnameChars[i] = rune(rand.Intn(26) + 97)
	}
	lastname := string(lastnameChars)
	randomNums := rand.Intn(900) + 100
	separator := ""
	if rand.Intn(2) > 0 {
		separator = "."
	}
	email := fmt.Sprintf("%s%s%s%d@%s", firstname, separator, lastname, randomNums, domain)

	return email
}

func ClearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

var (
	fileMutex sync.Mutex
)

func SaveAccountToFile(email, password string) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	file, err := os.OpenFile("accounts.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	accountLine := fmt.Sprintf("%s:%s\n", email, password)
	if _, err := file.WriteString(accountLine); err != nil {
		return err
	}

	return nil
}
