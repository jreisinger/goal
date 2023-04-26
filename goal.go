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
	Description string
	Strategy    string   // high-level plan to reach your goal
	Tactics     []Tactic // implementation of the strategy
}

// CivilTime represents time in the format "2006-01-02".
type CivilTime time.Time

// UnmarshalYAML implements yaml.Unmarshaler so CivilTime can be unmarshaled
// from a YAML document.
func (c *CivilTime) UnmarshalYAML(n *yaml.Node) error {
	value := strings.Trim(string(n.Value), `"`) //get rid of "
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse("2006-01-02", value) //parse time
	if err != nil {
		return err
	}
	*c = CivilTime(t) //set result using the pointer
	return nil
}

func (c *CivilTime) String() string {
	return time.Time(*c).Format("2006-01-02")
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

func (t Tactic) String() string {
	return fmt.Sprintf("%s (done: %s, interval: %s)", t.Do, &t.Done, t.Interval)
}

// Parse recursively parses files in dir into name and goal map. Name is the
// path of YAML file holding a goal.
func Parse(dir string) (map[string]Goal, error) {
	goals := make(map[string]Goal)

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
			goals[path] = g

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
strategy: Get a personal trainer and train consistently over the next five years.
tactics:
- do: Find an online community to share ideas and get tips.
  done: 0001-01-01 # can be ommitted
  interval: once # can be omitted 
- do: Find a personal trainer.
- do: Have a health/diet plan focused on mind, body and spirit.
  done: 1970-01-01
- do: Meditate daily 10 – 30 minutes.
  done: 2023-04-25 # will expire in a day because of daily interval
  interval: daily
- do: Train on Monday, Tuesday, Thursday and Friday (2 hours per session).
  interval: weekly`
}

func parse(yamlData []byte) (Goal, error) {
	var goal Goal
	if err := yaml.Unmarshal(yamlData, &goal); err != nil {
		return goal, err
	}
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

func Print(goals map[string]Goal, all bool) {
	// const format = "%v\t%v\n"
	// tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	// fmt.Fprintf(tw, format, "Goal", "Status")
	// fmt.Fprintf(tw, format, "----", "------")
	// for name, g := range goals {
	// 	fmt.Fprintf(tw, format, name, g.Status())
	// }
	// tw.Flush()

	for _, k := range sortKeys(goals) {
		fmt.Println()
		fmt.Println(k)
		g := goals[k]
		for _, t := range g.Tactics {
			if !all && t.isDone() {
				continue
			}
			if t.isDone() {
				fmt.Print("✅ ")
				fmt.Println(t)
			} else {
				fmt.Print("-  ")
				fmt.Println(t)
			}
		}
	}
}

func sortKeys(m map[string]Goal) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
