package main

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/process"
)

var playerID = "355473"
var port = "52949"

func findPacketloggerProcesses() string {
	processes, err := process.Processes()
	if err != nil {
		fmt.Printf("Error getting processes: %v\n", err)
		return ""
	}
	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			continue // Skip this process if we can't get its name
		}
		//fmt.Println(name)

		if strings.Contains(name, "NostaleClientX.exe") {
			fmt.Printf("Process found: %s\n", name)

			connections, err := p.Connections()
			if err != nil {
				fmt.Printf("Error getting connections for process %s: %v\n", name, err)
				continue
			}

			for _, conn := range connections {
				if conn.Status == "LISTEN" {
					fmt.Printf("Port: %d\n", conn.Laddr.Port)
					if conn.Laddr.Port > 0 {
						return fmt.Sprintf("%d", conn.Laddr.Port)
					}
				}
			}
			fmt.Println("---")
		}
	}
	return ""
}

func main() {
	port = findPacketloggerProcesses()
	fmt.Println("starting bot")
	fBot, err := initBot(port, playerID)
	if err != nil {
		panic(fmt.Sprintf("err: %s", err))
	}
	fBot.run()
}
