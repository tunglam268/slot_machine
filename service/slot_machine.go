package service

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var chips uint32 = 10

type SymbolType struct {
	Id     byte
	icon   string
	Reward uint32
}

type SlotMachineType struct {
	Reels   [][]SymbolType
	Credits uint32
}

var reward = 0

var Symbols = []SymbolType{
	{
		Id:     0,
		icon:   "🍎",
		Reward: 5,
	},
	{
		Id:     1,
		icon:   "🥭",
		Reward: 10,
	},
	{
		Id:     2,
		icon:   "🍓",
		Reward: 15,
	},
	{
		Id:     3,
		icon:   "🍉",
		Reward: 20,
	},
	{
		Id:     4,
		icon:   "🍇",
		Reward: 25,
	},
	{
		Id:     5,
		icon:   "🍒",
		Reward: 30,
	},
	{
		Id:     6,
		icon:   "💎",
		Reward: 50,
	},
	{
		Id:     7,
		icon:   "🍏",
		Reward: 100,
	},
}

func NewSlotMachine() *SlotMachineType {
	return &SlotMachineType{
		Reels:   [][]SymbolType{},
		Credits: 1000,
	}
}

func (s *SlotMachineType) Spin(chips uint32) {
	if s.Credits < 10 {
		fmt.Println("You dont have enough credits to continue")
		return
	}

	if s.Credits > 1000 {
		chips = 100
	}

	s.Credits -= chips
	s.Reels = make([][]SymbolType, 3)
	for i := range s.Reels {
		s.Reels[i] = make([]SymbolType, 3)
		for j := range s.Reels[i] {
			s.Reels[i][j] = generateRandomSymbol(0, 7)
		}
	}
}

func (s *SlotMachineType) Display(status string) {
	displayCred := 10
	if s.Credits > 1000 {
		displayCred = 100
	}

	CallClear()

	fmt.Println("-------------------------------------------------")
	fmt.Println("|🍎 = 10|🥭 = 10|🍉 = 10|🍓 = 15|🍇 = 20|🍒 = 20|💎 = (credits x 3)|🍏 = 100|")
	fmt.Printf("|spin = %d|\n", displayCred)

	for i := range s.Reels {
		fmt.Print("\n                 ")
		fmt.Printf("%s | %s | %s ", s.Reels[i][0].icon, s.Reels[i][1].icon, s.Reels[i][2].icon)
		fmt.Print("\n                 ")
		// for j := range s.Reels[i] {
		// 	fmt.Print(strings.TrimSpace(s.Reels[i][j].icon), " | ")
		// }
	}
	fmt.Print("\n")

	if s.Credits > 1000 {
		fmt.Println("| Higher stakes 🧈🧈🧈 |")
	}
	fmt.Printf("\nCredits: %d \n\n", s.Credits)
	fmt.Printf("Reward: %d \n\n", reward)
	fmt.Println(status)
	fmt.Println("Press (Enter) key to play more ")
	fmt.Println("Press e to exit")
	fmt.Println("-------------------------------------------------")
}

func (s *SlotMachineType) CheckWin() (bool, SymbolType) {
	for i := 0; i < 3; i++ {
		if checkLine(s.Reels[i][0], s.Reels[i][1], s.Reels[i][2]) {
			s.Credits += s.Reels[i][0].Reward
			reward = int(s.Reels[i][0].Reward)
			return true, s.Reels[i][0]
		}
		if checkLine(s.Reels[0][i], s.Reels[1][i], s.Reels[2][i]) {
			s.Credits += s.Reels[0][i].Reward
			reward = int(s.Reels[0][i].Reward)
			return true, s.Reels[0][i]
		}
	}

	if checkLine(s.Reels[0][0], s.Reels[1][1], s.Reels[2][2]) {
		s.Credits += s.Reels[1][1].Reward
		reward = int(s.Reels[1][1].Reward)
		return true, s.Reels[1][1]
	}
	if checkLine(s.Reels[0][2], s.Reels[1][1], s.Reels[2][0]) {
		s.Credits += s.Reels[1][1].Reward
		reward = int(s.Reels[1][1].Reward)
		return true, s.Reels[1][1]
	}

	reward = 0
	return false, SymbolType{}
}

func checkLine(a, b, c SymbolType) bool {
	return a == b && b == c
}

func generateRandomSymbol(min, max int) SymbolType {
	//? 	THIS IS BASED ON TIME - GO TOO FAST IT GIVES SAME NUMBER ID
	// source := rand.NewSource(time.Now().UnixNano())
	// rng := rand.New(source)

	return Symbols[rand.Intn(max-min)+1]
}

var clear map[string]func() //create a map for storing clear funcs

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func SlotMachine() {
	slot := NewSlotMachine()
	fmt.Println("Press enter to start")

	for {
		reader := bufio.NewReader(os.Stdin)
		playerInput, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(playerInput)) == "e" {
			break
		}
		if strings.ToLower(strings.TrimSpace(playerInput)) == "r" {
			slot.Credits += 100
		}

		slot.Spin(chips)
		status, multipier := slot.CheckWin()
		if status {
			if multipier.Id == 6 {
				multipliedReward := slot.Credits * 3
				if multipliedReward > 1000 {
					slot.Credits += 1000
					slot.Display("DIAMOND: +1000 REACHED MAX WIN!")
					slot.Display("You Won")
				} else {
					slot.Display("DIAMOND: CREDITS MULTIPLIED!")
					slot.Credits *= 3
				}
			} else {
				slot.Display("You Won")
			}
		} else {
			slot.Display("Unlucky try again?")
		}
		time.Sleep(3 * time.Second)
	}
}
