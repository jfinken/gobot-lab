package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/chzyer/readline"
	"github.com/jfinken/gobot-lab/adabot"
)

func getHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
func main() {

	pi := "\xCE\xA0"
	fmt.Printf("Come to the dork side we have %s\n", pi)
	rl, err := readline.NewEx(&readline.Config{
		Prompt:      "gobot> ",
		HistoryFile: getHomeDir() + "/.gobot_history",
	})

	if err != nil {
		panic(err)
	}
	defer rl.Close()
	// New robot evaluator and start cli loop...
	evaluator := adabot.NewEval()
	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		// detect EOF before parsing
		if line == "" {
			continue
		}
		// evaluate and control
		evaluator.Run(line)
	}
}
