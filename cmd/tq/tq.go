package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/4ever9/tq"
	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

var rootCMD = &cobra.Command{
	Use:   "tq [flags] <filter> [file]",
	Short: "tq - toml tool like jq",
	Args:  cobra.RangeArgs(0, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := ""
		selector := ""
		pipeMode := false
		if isPipeMode() {
			s := bufio.NewScanner(os.Stdin)
			s.Split(bufio.ScanLines)
			for s.Scan() {
				input = fmt.Sprintf("%s\n", s.Text())
			}
			pipeMode = true
		} else {
			if len(args) == 0 {
				return cmd.Help()
			}

			data, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}

			input = string(data)
		}

		selector, value := verifyArgs(args, pipeMode)

		// decode toml file
		m := make(map[string]interface{})
		md, err := toml.Decode(input, &m)
		if err != nil {
			return fmt.Errorf("toml decode: %s", err)
		}

		ret, err := tq.Handle(m, md, selector, value)
		if err != nil {
			return err
		}

		fmt.Println(ret)

		return nil
	},
}

func isPipeMode() bool {
	f, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if (f.Mode() & os.ModeNamedPipe) != os.ModeNamedPipe {
		return false
	}

	return true
}

func verifyArgs(args []string, pipeMode bool) (string, string) {
	selector := ""
	value := ""
	if pipeMode {
		if len(args) > 0 {
			selector = args[0]
		}

		if len(args) > 1 {
			value = args[1]
		}

		return selector, value
	}

	if len(args) > 1 {
		selector = args[1]
	}

	if len(args) > 2 {
		value = args[2]
	}

	return selector, value
}

func main() {
	rootCMD.AddCommand(versionCMD)

	if err := rootCMD.Execute(); err != nil {
		os.Exit(1)
	}
}
