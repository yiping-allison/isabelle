package cmd

import "testing"

func TestToLowerAndFormat(t *testing.T) {
	tests := map[string]struct {
		input []string
		want  string
	}{
		"simple correct case": {
			input: []string{"Bee"},
			want:  "bee",
		},
		"multiple arguments": {
			input: []string{"Darner", "Dragonfly"},
			want:  "darner_dragonfly",
		},
		"mixed caps case": {
			input: []string{"dIving", "BeEtle"},
			want:  "diving_beetle",
		},
		"apostrophe in name": {
			input: []string{"Rajah", "Brooke's", "Birdwing"},
			want:  "rajah_brooke's_birdwing",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := toLowerAndFormat(tc.input)
			if got != tc.want {
				t.Errorf("toLowerAndFormat() got = %s; want %s", got, tc.want)
			}
		})
	}
}

func TestFormatName(t *testing.T) {
	tests := map[string]struct {
		input []string
		want  string
	}{
		"simple case": {
			input: []string{"tarantula"},
			want:  "Tarantula",
		},
		"multiple args": {
			input: []string{"orchid", "mantis"},
			want:  "Orchid Mantis",
		},
		"mixed case args": {
			input: []string{"gIaNt", "WAter", "bug"},
			want:  "Giant Water Bug",
		},
		"with apostrophes": {
			input: []string{"rajah", "Brooke's", "butterfly"},
			want:  "Rajah Brooke's Butterfly",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := formatName(tc.input)
			if got != tc.want {
				t.Errorf("formatName() got = %s; want %s", got, tc.want)
			}
		})
	}
}
