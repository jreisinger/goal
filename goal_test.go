package goal_test

import (
	"testing"

	"github.com/jreisinger/goal"
)

func TestDone(t *testing.T) {
	tests := []struct {
		g    goal.Goal
		want string
	}{
		{
			g:    goal.Goal{},
			want: "NaN% (0/0)",
		},
		{
			g: goal.Goal{
				Tactics: []goal.Tactic{
					{Do: "something", Done: false},
				},
			},
			want: "00% (0/1)",
		},
		{
			g: goal.Goal{
				Tactics: []goal.Tactic{
					{Do: "something", Done: false},
					{Do: "something else", Done: true},
				},
			},
			want: "50% (1/2)",
		},
	}

	for _, test := range tests {
		got := test.g.Done()
		if got != test.want {
			t.Errorf("goal.Done() == %q, want %q", got, test.want)
		}
	}
}
