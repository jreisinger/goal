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

	d = flag.String("d", filepath.Join(home, "goal"), "directory holding yaml files")
	e = flag.Bool("e", false, "print example yaml file content and exit")
	v = flag.Bool("v", false, "be verbose")
)

func main() {
	flag.Parse()

	if *e {
		fmt.Println(goal.Example())
		os.Exit(0)
	}

	goals, err := goal.Parse(*d)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parsing yaml files in %s: %v\n", *d, err)
		os.Exit(1)
	}
	goal.Print(goals, *v)
}
