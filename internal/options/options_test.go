package options

import "testing"

func TestAsMapRepeatLastN(t *testing.T) {
	tests := []struct {
		name        string
		repeatLastN int
		numCtx      int
		want        int
	}{
		{"negative sentinel resolves to num_ctx", -1, 8192, 8192},
		{"zero is passed through (disabled)", 0, 8192, 0},
		{"positive value is passed through", 64, 8192, 64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := ModelOptions{RepeatLastN: tt.repeatLastN, NumCtx: tt.numCtx}.AsMap()
			got, ok := m["repeat_last_n"].(int)
			if !ok {
				t.Fatalf("repeat_last_n has type %T, want int", m["repeat_last_n"])
			}
			if got != tt.want {
				t.Errorf("repeat_last_n = %d, want %d", got, tt.want)
			}
		})
	}
}
