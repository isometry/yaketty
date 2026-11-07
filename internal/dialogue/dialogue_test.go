package dialogue

import (
	"testing"

	"github.com/isometry/yaketty/internal/persona"
)

func newTestDialogue(think1, think2 bool) *Dialogue {
	return &Dialogue{
		Personas: [2]*persona.Persona{
			{Name: "Jane", Model: "m", Think: think1},
			{Name: "John", Model: "m", Think: think2},
		},
	}
}

func TestFromPerspectiveThink(t *testing.T) {
	tests := []struct {
		name   string
		think1 bool
		think2 bool
		bot    BotID
		want   bool
	}{
		{"default is false", false, false, Persona1, false},
		{"persona1 enabled", true, false, Persona1, true},
		{"persona2 enabled", false, true, Persona2, true},
		{"persona2 unaffected by persona1", true, false, Persona2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := newTestDialogue(tt.think1, tt.think2).FromPerspective(tt.bot)
			if req.Think == nil {
				t.Fatal("Think is nil, want an explicit value")
			}
			if got := req.Think.Bool(); got != tt.want {
				t.Errorf("Think = %v, want %v", got, tt.want)
			}
		})
	}
}
