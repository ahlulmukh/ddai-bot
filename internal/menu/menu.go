package menu

import (
	"bufio"
	"ddai-bot/internal/utils"
	"fmt"
	"os"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

type MenuHandler struct {
	reader *bufio.Reader
}

func NewMenuHandler() *MenuHandler {
	return &MenuHandler{
		reader: bufio.NewReader(os.Stdin),
	}
}

func (m *MenuHandler) ShowMainMenu(version string) {
	m.showBanner(version)

	if !m.configExists() {
		m.createConfig()
	}

	for {
		choice := m.showMenuOptions()
		switch choice {
		case "1":
			m.RunReferralProgram()
			m.showBanner(version)
		case "2":
			//m.RunAutoBot()
			m.showBanner(version)
		case "3":
			m.EditConfig()
			m.showBanner(version)
		case "4":
			m.ShowFileInfo()
			m.showBanner(version)
		case "5":
			os.Exit(0)
		default:
			utils.LogMessage(0, 0, "Invalid choice", "error")
		}
	}
}

func (m *MenuHandler) showBanner(version string) {
	//utils.ClearScreen()
	myFigure := figure.NewFigure("Ddai Depin", "", true)
	figureStr := myFigure.String()
	fmt.Println(color.HiYellowString(figureStr))

	fmt.Println(color.HiBlueString("🔹 Made by : ") + color.WhiteString("El Puqus Airdrop") + "  |  " +
		color.HiBlueString("📦 GitHub: ") + color.WhiteString("github.com/ahlulmukh") + "  |  " +
		color.HiBlueString("✅ Version: ") + color.WhiteString(version))
}

func (m *MenuHandler) showMenuOptions() string {
	fmt.Println(color.CyanString("\nMain Menu:"))
	fmt.Println(color.HiGreenString("1. Auto Referral + Auto Claim Task"))
	fmt.Println(color.YellowString("2. Auto Bot (multiple modes)"))
	fmt.Println(color.HiBlueString("3. Edit Config"))
	fmt.Println(color.MagentaString("4. Information"))
	fmt.Println(color.RedString("5. Exit"))
	fmt.Print(color.HiCyanString("Enter your choice (1-5): "))

	choice, _ := m.reader.ReadString('\n')
	return strings.TrimSpace(choice)
}

func (m *MenuHandler) ShowFileInfo() {
	utils.ClearScreen()
	fmt.Println(color.CyanString("\nFile Requirements:"))
	fmt.Println(color.YellowString("1. accounts.txt"))
	fmt.Println("   - This file contains all your accounts referral")
	fmt.Println()

	fmt.Println(color.YellowString("2. runaccounts.txt"))
	fmt.Println("   - This file is used to run the bot")
	fmt.Println("   - Format: email:password")
	fmt.Println("   - Example: user1@gmail.com:password123")
	fmt.Println()

	fmt.Println(color.YellowString("3. proxy.txt (optional)"))
	fmt.Println("   - Format: user:pass@host:port")
	fmt.Println("   - Example: http://puqus:gaming@yesyes.com:8080")
	fmt.Println()

	fmt.Println(color.HiGreenString("Usage Instructions:"))
	fmt.Println("1. For auto referral, accounts saved to accounts.txt")
	fmt.Println("2. For auto bot, add accounts to runaccounts.txt")

	m.waitForEnter()
}

func (m *MenuHandler) showBotModeMenu() string {
	utils.ClearScreen()
	fmt.Println(color.CyanString("\nBot Running Mode:"))
	fmt.Println(color.HiGreenString("1. Concurrent Mode (all accounts run simultaneously)"))
	fmt.Println(color.YellowString("2. Queue Mode (accounts run one by one)"))
	fmt.Println(color.RedString("3. Back to Main Menu"))
	fmt.Print(color.HiCyanString("Choose mode (1-3): "))

	choice, _ := m.reader.ReadString('\n')
	return strings.TrimSpace(choice)
}

func (m *MenuHandler) waitForEnter() {
	fmt.Print("Press Enter to continue...")
	m.reader.ReadBytes('\n')
}
