package console

import (
	"bufio"
	"gonet"
	"os"
	"strings"
)

var consoleAgent uint32
var loop bool = true

func SetAgent(agent uint32) {
	consoleAgent = agent
}

func Begin() {
	if !loop {
		loop = true
		go consoleLoop()
	}
}

func Stop() {
	loop = false
}

func consoleLoop() {
	reader := bufio.NewReader(os.Stdin)

	for loop {
		data, _, _ := reader.ReadLine()
		command := string(data)

		args := strings.Split(command, " ")

		if consoleAgent != 0 && loop {
			gonet.Send(consoleAgent, "ConsoleCommand", args)
		}
	}
}

func init() {
	go consoleLoop()
}
