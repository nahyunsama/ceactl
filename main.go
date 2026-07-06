package main

import (
	"os"

	"github.com/nahyunsama/ceactl/cmd"
)

func main() {
	os.Setenv("GODEBUG", "tlssha1=1")
	cmd.Execute()
}
