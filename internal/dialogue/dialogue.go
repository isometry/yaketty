package dialogue

import (
	"context"
	"log/slog"
	"strings"

	"github.com/ollama/ollama/api"

	"github.com/isometry/yaketty/internal/config"
	"github.com/isometry/yaketty/internal/output"
	"github.com/isometry/yaketty/internal/persona"
	"github.com/isometry/yaketty/internal/scenario"
)

type BotID int

const (
	Persona1 BotID = iota
	Persona2
)

// Configuration constants
const (
	reminderInterval    = 12
	defaultPersona1Name = "Jane"
	defaultPersona2Name = "John"
)

var (
	defaultPrompts = []string{
		`You are playing a character in a dialogue scenario. The user represents the other character.

		 Embody your assigned persona completely - adopt their knowledge, beliefs, vocabulary, mannerisms, and communication style.
		 Build meaningfully on previous exchanges and provide responses that advance the dialogue.
		 Stay authentic to your character's worldview and never break character or make meta-commentary about being AI.
		 Keep your statements and responses brief and relevant to the scenario; avoid monologues.
		 Expect the other character to respond appropriately to the scenario, and remember that you're conversing with them.
		 Never repeat yourself unless explicitly prompted.
		 *SPEAK* as your character.

		 Your character details and scenario context follow.`,
	}

	// Brief periodic reminders (injected every ~10-15 exchanges)
	periodicReminder = `Remember: stay true to your character and the scenario context.`
)

func (b BotID) Opponent() BotID {
	return 1 - b
}

type Dialogue struct {
	// Embedded scenario configuration
	scenario.Scenario

	// Configuration
	ExtraPrompts []string
	Personas     [2]*persona.Persona
	Output       output.OutputStyle

	// Runtime state
	Messages []*Message

	// Internal dependencies
	ctx    context.Context
	client *api.Client
}

const (
	systemRole    = "system"
	userRole      = "user"
	assistantRole = "assistant"
)

type Message struct {
	persona BotID
	content string
}

func newMessage(role, content string) api.Message {
	message := api.Message{Role: role, Content: content}
	return message
}

func systemMessages(contents ...string) []api.Message {
	messages := make([]api.Message, 0, len(contents))
	for _, content := range contents {
		messages = append(messages, newMessage(systemRole, content))
	}
	return messages
}

func NewDialogue(ctx context.Context, cfg *config.Config) (*Dialogue, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, err
	}

	return &Dialogue{
		ctx:          ctx,
		client:       client,
		Scenario:     cfg.Scenario,
		ExtraPrompts: cfg.ExtraPrompts,
		Personas: [2]*persona.Persona{
			&cfg.Persona1,
			&cfg.Persona2,
		},
		Output: output.Text{},
	}, nil
}

func (c *Dialogue) AddMessage(botID BotID, content string) {
	c.Output.Render(c.Personas[botID].Name, content)
	c.Messages = append(c.Messages, &Message{botID, content})
}

func (c *Dialogue) FromPerspective(botID BotID) api.ChatRequest {
	prompts := make([]string, 0, 6+len(defaultPrompts)+len(c.ExtraPrompts))
	prompts = append(prompts, defaultPrompts...)

	prompts = append(prompts,
		c.Scenario.Scenario,
		c.Personas[botID].Persona,
		c.Roles[botID],
	)

	// Add periodic reminder every reminderInterval messages
	if len(c.Messages) > 0 && len(c.Messages)%reminderInterval == 0 {
		prompts = append(prompts, periodicReminder)
	}

	prompts = append(prompts, c.ExtraPrompts...)

	messages := make([]api.Message, 0, len(c.Messages)+len(prompts)+1)
	messages = append(messages, systemMessages(prompts...)...)

	for _, m := range c.Messages {
		if m.persona == botID {
			messages = append(messages, newMessage(assistantRole, m.content))
		} else {
			messages = append(messages, newMessage(userRole, m.content))
		}
	}

	cr := api.ChatRequest{
		Model:    c.Personas[botID].Model,
		Messages: messages,
		Options:  c.Personas[botID].Options.AsMap(),
		Stream:   func() *bool { b := false; return &b }(),
	}

	return cr
}

func (c *Dialogue) SendRequest(botID BotID) error {
	chatRequest := c.FromPerspective(botID)
	slog.Debug("sending chat request", slog.String("perspective", c.Personas[botID].Name), slog.Any("chatRequest", chatRequest))
	return c.client.Chat(c.ctx, &chatRequest, c.HandleResponse(botID))
}

func (c *Dialogue) HandleResponse(botID BotID) func(api.ChatResponse) error {
	return func(cr api.ChatResponse) error {
		message := strings.TrimSpace(cr.Message.Content)
		if len(message) == 0 && len(c.Messages) > 0 {
			if len(c.Messages[len(c.Messages)-1].content) == 0 {
				return nil
			} else {
				c.AddMessage(botID, "...")
			}
		} else {
			c.AddMessage(botID, message)
		}
		return c.SendRequest(botID.Opponent())
	}
}

func (c *Dialogue) Start() error {
	// Send opening prompt to Persona1 as a system instruction
	prompts := make([]string, 0, 6+len(defaultPrompts)+len(c.ExtraPrompts)+1)
	prompts = append(prompts, defaultPrompts...)

	prompts = append(prompts,
		c.Scenario.Scenario,
		c.Personas[Persona1].Persona,
		c.Roles[Persona1],
	)

	prompts = append(prompts, c.ExtraPrompts...)
	prompts = append(prompts, c.OpeningPrompt)

	messages := systemMessages(prompts...)

	chatRequest := api.ChatRequest{
		Model:    c.Personas[Persona1].Model,
		Messages: messages,
		Options:  c.Personas[Persona1].Options.AsMap(),
		Stream:   func() *bool { b := false; return &b }(),
	}

	slog.Debug("sending opening request", slog.String("perspective", c.Personas[Persona1].Name), slog.Any("chatRequest", chatRequest))
	return c.client.Chat(c.ctx, &chatRequest, c.HandleResponse(Persona1))
}
