package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
)

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Run auto in interactive CLI mode",
	Long:  `Run auto in interactive CLI mode`,
	Run:   cliRun,
}

var (
	rl *readline.Instance
)

func init() {
	rootCmd.AddCommand(cliCmd)

	var completer = readline.NewPrefixCompleter()
	for _, child := range rootCmd.Commands() {
		pcFromCommands(completer, child)
	}

	var err error
	rl, err = readline.NewEx(&readline.Config{
		Prompt:            "otto\033[31m»\033[0m ",
		HistoryFile:       "/tmp/readline.tmp",
		AutoComplete:      completer,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
		// FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}
	rl.CaptureExitSignal()
	log.SetOutput(rl.Stderr())

}

func cliRun(cmd *cobra.Command, args []string) {
	defer rl.Close()
	done := false
	for !done {
		done = runLine()
	}
}

func pcFromCommands(parent readline.PrefixCompleterInterface, c *cobra.Command) {
	pc := readline.PcItem(c.Use)
	parent.SetChildren(append(parent.GetChildren(), pc))
	for _, child := range c.Commands() {
		pcFromCommands(pc, child)
	}
}

func runLine() bool {
	line, err := rl.Readline()
	if err == readline.ErrInterrupt {
		if len(line) == 0 {
			return true
		} else {
			return false
		}
	} else if err == io.EOF {
		return true
	}

	line = strings.TrimSpace(line)
	if line == "exit" || line == "quit" {
		return true
	}

	args := strings.Split(line, " ")
	cmd, args, err := rootCmd.Find(args)
	if err != nil {
		fmt.Printf("Error running cmd %q: %s\n", line, err)
	}

	cmd.ParseFlags(args)
	cmd.Run(cmd, args)
	return false
}