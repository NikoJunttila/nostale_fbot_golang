package main

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

type FishBot struct {
	socket      net.Conn
	expBuff     bool
	lineBuff    bool
	baitSkill   bool
	castLine    bool
	outOfBaits  bool
	proCastLine bool
	playerID    string
	covert      bool
	stop        bool
}

func initBot(port, pID string) (*FishBot, error) {
	ip4 := fmt.Sprintf("127.0.0.1:%s", port)
	conn, err := net.Dial("tcp", ip4)
	if err != nil {
		fmt.Println("Failed to connect:", err)
		return nil, err
	}

	return &FishBot{
		socket:      conn,
		expBuff:     true,
		lineBuff:    true,
		baitSkill:   true,
		castLine:    true,
		proCastLine: true,
		covert:      true,
		playerID:    pID,
	}, nil
}

var fishCount int

func (bot *FishBot) handleIN(words []string) {
	if len(words) < 10 {
		return
	}
	if words[9] != "2" {
		return
	}
	fmt.Println("\n\nADMIN ID DETECTED ALERT!!!!")
	bot.stop = true
	err := playSound("allu.mp3")
	if err != nil {
		fmt.Println("Error playing sound:", err)
	}
	bot.rs(4000, 1000)
}

func (bot *FishBot) handleCMap() {
	fmt.Println("Map change!!!")
	fmt.Println("ADMIN ALERT!!!!")
	bot.stop = true
	err := playSound("allu.mp3")
	if err != nil {
		fmt.Println("Error playing sound:", err)
	}
	bot.rs(4000, 1000)
}

func (bot *FishBot) handleGURI(line []string) {
	switch line[5] {
	case "30":
		fishCount++
		fmt.Println("fish: ", fishCount)
		bot.rs(400, 300)
		bot.castSkill("2")
	case "31":
		fishCount++
		fmt.Println("fish: ", fishCount)
		fmt.Println("legendary fish!!")
		bot.rs(400, 300)
		bot.castSkill("2")
	case "0":
		go bot.checkBuffs()
	}
}

func (bot *FishBot) rs(min, random int) {
	sleepDuration := time.Duration(rand.Intn(random)+min) * time.Millisecond
	time.Sleep(sleepDuration)
}

func (bot *FishBot) handleSayi(line []string) {
	if len(line) < 5 {
		return
	}
	if line[2] != "1" {
		return
	}
	if line[3] != bot.playerID {
		return
	}
	if line[5] == "2497" {
		bot.outOfBaits = true
		fmt.Println("out of baits")
	}
}

func (bot *FishBot) handleSR(skillID string) {
	switch skillID {
	case "1":
		bot.castLine = true
	case "3":
		bot.baitSkill = true
		if bot.outOfBaits {
			go bot.checkBuffs()
		}
	case "8":
		bot.expBuff = true
	case "5":
		bot.covert = true
	case "9":
		bot.lineBuff = true
	case "10":
		bot.proCastLine = true
	}
}

func (bot *FishBot) checkBuffs() {
	bot.rs(1000, 300)

	if bot.lineBuff {
		bot.castSkill("9")
		bot.lineBuff = false
		bot.rs(2000, 1000)
	}
	if bot.baitSkill {
		bot.castSkill("3")
		bot.baitSkill = false
		bot.rs(1500, 1000)
		if bot.outOfBaits {
			if bot.proCastLine {
				bot.castSkill("10")
			} else {
				bot.castSkill("1")
			}
			return
		}
	}
	if bot.outOfBaits {
		return
	}
	if bot.proCastLine {
		bot.castSkill("10")
		bot.proCastLine = false
		return
	}
	bot.castSkill("1")
}

func (bot *FishBot) castSkill(skillID string) {
	skill := fmt.Sprintf("1 u_s %s 1 %s", skillID, bot.playerID)
	_, err := bot.socket.Write([]byte(skill))
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}
}

func (bot *FishBot) run() {
	defer bot.socket.Close()

	buf := make([]byte, 1024*2)
	for {
		bytesRead, err := bot.socket.Read(buf)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}
		receivedData := string(buf[:bytesRead])
		splitted := strings.Split(receivedData, "\n")
		for _, line := range splitted {
			if line != "" {
				if bot.stop {
					fmt.Println("stopping bot")
					return
				}
				words := strings.Fields(line)
				if len(words) < 2 {
					continue
				}
				switch words[1] {
				case "guri":
					if len(words) < 6 {
						continue
					}
					if words[4] != bot.playerID {
						continue
					}
					bot.handleGURI(words)
				case "sr":
					if len(words) == 3 {
						bot.handleSR(words[2])
					}
				case "sayi":
					bot.handleSayi(words)
				case "c_map":
					bot.handleCMap()
				case "in":
					bot.handleIN(words)
				}
			}
		}
	}
}

/* func main() {
	fmt.Println("starting bot")
	fBot, err := initBot(port, playerID)
	if err != nil {
		panic(fmt.Sprintf("err: %s", err))
	}
	fBot.run()
}
*/
