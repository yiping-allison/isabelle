package daisymaebot_test

import (
	"reflect"
	"testing"

	"github.com/yiping-allison/daisymae/daisymaebot"
)

func TestNew(t *testing.T) {
	tests := map[string]struct {
		bot_key string
		want    daisymaebot.Bot
		err     error
	}{
		"need to have bot key": {
			bot_key: "",
			want:    daisymaebot.Bot{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			bot, err := daisymaebot.New(tc.bot_key)
			if err == nil {
				t.Errorf("New() Need to have error stating user needs BotKey")
			}
			if !reflect.DeepEqual(bot, &tc.want) {
				t.Errorf("New() not empty; got = %v; want = %v", bot, tc.want)
			}
		})
	}
}
