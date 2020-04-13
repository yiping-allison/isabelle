package cmd

import (
	"reflect"
	"testing"
)

func TestParseEvent(t *testing.T) {
	tests := map[string]struct {
		cmd   string
		img   string
		name  string
		event *newEvent
	}{
		"not enough args": {
			cmd:   "msg=\"testing\"",
			img:   errThumbURL,
			name:  "Error",
			event: nil,
		},
		"wrong args": {
			cmd:   "msg=\"testing\" lol=\"hi\"",
			img:   errThumbURL,
			name:  "Error",
			event: nil,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := parseCmd(tc.cmd, tc.name, tc.img)
			if !reflect.DeepEqual(tc.event, got) {
				t.Errorf("parseEvent() got = %v; want %v", got, tc.event)
			}
		})
	}
}
