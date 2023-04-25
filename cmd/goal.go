// Goal helps you achieve your goals by using strategy and tactics.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jreisinger/goal"
)

var (
	home, _ = os.UserHomeDir()
	dir     = flag.String("dir", filepath.Join(home, "goal"), "directory holding yaml files")
	example = flag.Bool("example", false, "print example yaml file content and exit")
)

func main() {
	flag.Parse()

	if *example {
		fmt.Println(goal.Example())
		os.Exit(0)
	}

	goals, err := goal.Parse(*dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parsing yaml files in %s: %v\n", *dir, err)
		os.Exit(1)
	}
	goal.Print(os.Stdout, goals)
}
