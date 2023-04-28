// Package goal helps you achieve your goals by using strategy and tactics.
package goal

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	Once Interval = iota
	Daily
	Weekly
	Monthly
)

type Interval int

func (i Interval) String() string {
	return [...]string{"once", "daily", "weekly", "monthly"}[i]
}

// Goal is where you want to get or what you want to achieve.
type Goal struct {
	Path        string // filesystem path; filename is goal name
	Description string
	Strategy    string    // high-level plan to reach your goal
	Tactics     []Tactic  // implementation of the strategy
	Updated     CivilTime // last update
}

// CivilTime represents time in the format "2006-01-02".
type CivilTime time.Time

// UnmarshalYAML implements yaml.Unmarshaler so CivilTime can be unmarshaled
// from a YAML document.
func (c *CivilTime) UnmarshalYAML(n *yaml.Node) error {
	value := strings.Trim(string(n.Value), `"`) // get rid of "
	switch value {
	case "", "null":
		return nil
	case "unknown":
		value = "1970-01-01"
	case "never":
		value = "0001-01-01"
	}

	t, err := time.Parse("2006-01-02", value) //parse time
	if err != nil {
		return err
	}
	*c = CivilTime(t) // set result using the pointer
	return nil
}

func (c *CivilTime) String() string {
	t := time.Time(*c)
	switch {
	case t.IsZero():
		return "never"
	case t.Equal(time.Unix(0, 0)):
		return "unknown"
	default:
		days := time.Since(t).Hours() / 24
		return fmt.Sprintf("%.0fd ago", days)
		// return time.Time(*c).Format("2006-01-02")
	}
}

// UnmarshalYAML implements yaml.Unmarshaler so Interval can be unmarshaled from
// a YAML document.
func (i *Interval) UnmarshalYAML(n *yaml.Node) error {
	switch n.Value {
	case "", "once":
		*i = Once
		return nil
	case "daily":
		*i = Daily
		return nil
	case "weekly":
		*i = Weekly
		return nil
	case "monthly":
		*i = Monthly
		return nil
	default:
		return fmt.Errorf("unknown interval: %s", n.Value)
	}
}

// Tactic defines what to do and whether it's already done.
type Tactic struct {
	Do       string
	Done     CivilTime `yaml:"done,omitempty"`     // defaults to 0001-01-01
	Interval Interval  `yaml:"interval,omitempty"` // defaults to once
}

// Parse recursively parses files in dir into goals.
func Parse(dir string) ([]Goal, error) {
	var goals []Goal

	visit := func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			ext := filepath.Ext(entry.Name())
			if ext != ".yaml" && ext != ".yml" {
				return nil
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			b, err := io.ReadAll(file)
			if err != nil {
				return err
			}
			g, err := parse(b)
			if err != nil {
				return fmt.Errorf("parsing %s: %w", entry.Name(), err)
			}
			g.Path = path
			goals = append(goals, g)

		}
		return nil
	}

	if err := filepath.WalkDir(dir, visit); err != nil {
		return goals, err
	}

	return goals, nil
}

// Example returns sample YAML file content.
func Example() string {
	return `description: Become a black belt martial artist in under five years.
strategy: Get a personal trainer and train consistently.
tactics:
- do: Find a personal trainer.
  interval: once 	# default, can be omitted
  done: never 	 	# default, can be ommitted
- do: Train daily, 2 hours per session.
  interval: daily  	# or weekly, monthly
  done: 2023-04-27 	# will expire in a day because of daily interval
- do: Have a health/diet plan focused on mind, body and spirit.
  done: unknown    	# I've done this but don't know the date`
}

func parse(yamlData []byte) (Goal, error) {
	var goal Goal
	if err := yaml.Unmarshal(yamlData, &goal); err != nil {
		return goal, err
	}

	var updated time.Time
	for _, t := range goal.Tactics {
		if time.Time(t.Done).After(updated) {
			updated = time.Time(t.Done)
		}
	}

	goal.Updated = CivilTime(updated)
	return goal, nil
}

func (t Tactic) isDone() bool {
	switch t.Interval {
	case Once:
		return !time.Time(t.Done).IsZero()
	case Daily:
		return time.Since(time.Time(t.Done)) < time.Hour*24
	case Weekly:
		return time.Since(time.Time(t.Done)) < time.Hour*24*7
	case Monthly:
		return time.Since(time.Time(t.Done)) < time.Hour*24*7*30
	}
	return false // should never get here
}
func printTactic(t Tactic, verbose bool) {
	if !verbose && t.isDone() {
		return
	}
	if t.isDone() {
		fmt.Print("✅ ")
	} else {
		fmt.Print("👉 ")
	}
	fmt.Printf("%s (do: %s, done: %s)", t.Do, t.Interval, &t.Done)
	fmt.Println()
}

func Print(goals []Goal, verbose bool) {
	// const format = "%v\t%v\n"
	// tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	// fmt.Fprintf(tw, format, "Goal", "Status")
	// fmt.Fprintf(tw, format, "----", "------")
	// for name, g := range goals {
	// 	fmt.Fprintf(tw, format, name, g.Status())
	// }
	// tw.Flush()

	circledNumbers := map[int]string{
		0: "⓪", 1: "①", 2: "②", 3: "③", 4: "④",
		5: "⑤", 6: "⑥", 7: "⑦", 8: "⑧", 9: "⑨",
	}

	sortGoals(goals)

	const sep = "--------------------------------------------------------------------------------"

	for i, g := range goals {
		fmt.Printf("%s %s", circledNumbers[i+1], g.Path)
		if verbose {
			fmt.Printf(" (updated: %s)\n", &g.Updated)
			fmt.Printf("🏁 %s\n", g.Description)
			fmt.Printf("🧭 %s", g.Strategy)
		}
		fmt.Println()
		for _, t := range g.Tactics {
			printTactic(t, verbose)
		}
		fmt.Println(sep)
	}
}

type customSort struct {
	goals []Goal
	less  func(x, y Goal) bool
}

func (x customSort) Len() int           { return len(x.goals) }
func (x customSort) Less(i, j int) bool { return x.less(x.goals[i], x.goals[j]) }
func (x customSort) Swap(i, j int)      { x.goals[i], x.goals[j] = x.goals[j], x.goals[i] }

func sortGoals(goals []Goal) {
	sort.Sort(customSort{goals, func(x, y Goal) bool {
		if x.Updated != y.Updated {
			return time.Time(x.Updated).After(time.Time(y.Updated))
		}
		if x.Path != y.Path {
			return x.Path < y.Path
		}
		return false
	}})
}
