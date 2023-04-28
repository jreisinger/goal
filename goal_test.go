package goal

import (
	"testing"
	"time"
)

func TestCivilTimeString(t *testing.T) {
	tests := []struct {
		c    CivilTime
		want string
	}{
		{CivilTime{}, "never"},
		{CivilTime(time.Now()), "0d ago"},
		{CivilTime(time.Now().Add(-time.Hour * 23)), "0d ago"},
		{CivilTime(time.Now().Add(-time.Hour * 24)), "1d ago"},
		{CivilTime(time.Now().Add(-time.Hour * 47)), "1d ago"},
		{CivilTime(time.Now().Add(-time.Hour * 48)), "2d ago"},
		{CivilTime(time.Now().Add(time.Hour)), "time travel, huh?"},
	}
	for i, test := range tests {
		got := test.c.String()
		if got != test.want {
			t.Errorf("test %d: got %q, want %q", i, got, test.want)
		}
	}
}
