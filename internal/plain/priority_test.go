package plain

import "testing"

func TestParsePriority(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{name: "low", input: "low", want: 0},
		{name: "normal", input: "normal", want: 1},
		{name: "high", input: "high", want: 2},
		{name: "urgent", input: "urgent", want: 3},
		{name: "trimmed numeric", input: " 2 ", want: 2},
		{name: "mixed case", input: "UrGent", want: 3},
		{name: "invalid text", input: "critical", wantErr: true},
		{name: "out of range", input: "9", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePriority(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %d, got %d", tt.want, got)
			}
		})
	}
}
