// Package goal helps you achieve your goals by using strategy and tactics.
package goal

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"text/tabwriter"

	"gopkg.in/yaml.v3"
)

// Goal is where you want to get or what you want to achieve.
type Goal struct {
	Description string
	Strategy    string   // high-level plan to reach your goal
	Tactics     []Tactic // implementation of the strategy
}

// Tactic defines what to do and whether it's already done.
type Tactic struct {
	Do   string
	Done bool `yaml:"done,omitempty"`
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
				return err
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
  - do: Find a personal trainer.
    done: true
  - do: Set annual, monthly and weekly goals.
    done: false # can be omitted
  - do: Have a health/diet plan focused on mind, body and spirit.
  - do: Develop a series of minor milestones (to stay motivated).
  - do: Research martial arts instructors in this area.
  - do: Find a ‘training buddy’.
  - do: Find an online community to share ideas and get tips.
  - do: Train on Monday, Tuesday, Thursday and Friday (2 hours per session).
  - do: Write a diet plan.
  - do: Buy training equipment for home use.
  - do: Meditate daily (10 – 30 minutes).
  - do: Develop a ‘rewards’ scheme for minor milestones achieved.`
}

func parse(yamlData []byte) (Goal, error) {
	var goal Goal
	if err := yaml.Unmarshal(yamlData, &goal); err != nil {
		return goal, err
	}
	return goal, nil
}

// Done returns percentage and number of the steps done out of all the steps.
func (g Goal) Done() string {
	var total, done int
	for _, step := range g.Tactics {
		total++
		if step.Done {
			done++
		}
	}
	return fmt.Sprintf("%02.0f%% (%d/%d)", float64(done)/float64(total)*100, done, total)
}

func Print(w io.Writer, goals map[string]Goal) {
	const format = "%v\t%v\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Goal", "Done")
	fmt.Fprintf(tw, format, "----", "----")
	for name, g := range goals {
		fmt.Fprintf(tw, format, name, g.Done())
	}
	tw.Flush()
}
