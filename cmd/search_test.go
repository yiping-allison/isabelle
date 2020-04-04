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

func TestParseHemi(t *testing.T) {
	tests := map[string]struct {
		northSt   string
		northEnd  string
		southSt   string
		southEnd  string
		wantNorth string
		wantSouth string
	}{
		"simple case": {
			northSt:   "3",
			northEnd:  "6",
			southSt:   "7",
			southEnd:  "2",
			wantNorth: "March to June",
			wantSouth: "July to February",
		},
		"split case": {
			northSt:   "3|9",
			northEnd:  "4|12",
			southSt:   "5|1",
			southEnd:  "8|2",
			wantNorth: "March to April, September to December",
			wantSouth: "May to August, January to February",
		},
		"single month case": {
			northSt:   "3|10",
			northEnd:  "6|10",
			southSt:   "9|4",
			southEnd:  "12|4",
			wantNorth: "March to June, October",
			wantSouth: "September to December, April",
		},
		"only one month": {
			northSt:   "3",
			northEnd:  "3",
			southSt:   "7",
			southEnd:  "7",
			wantNorth: "March",
			wantSouth: "July",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotN, gotS := parseHemi(tc.northSt, tc.northEnd, tc.southSt, tc.southEnd)
			if gotN != tc.wantNorth {
				t.Errorf("parseHemi() north = %v; want %v", gotN, tc.wantNorth)
			}
			if gotS != tc.wantSouth {
				t.Errorf("parseHemi() south = %v; want %v", gotS, tc.wantSouth)
			}
		})
	}
}
