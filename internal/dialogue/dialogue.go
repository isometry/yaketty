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

var False = false

type BotID int

const (
	Persona1 BotID = iota
	Persona2
)

var defaultPrompts = []string{
	` A scenario involving two characters will be provided to you.
	 You will be assigned a character to play, and you must inhabit that character's persona completely for the duration of the exchange, responding in-character throughout, and in accordance with the circumstances of the described scenario.
	 The user will be playing the part of the **other character**.
	 You must adopt all of **your** assigned persona's experience, knowledge, beliefs, opinions, vocabulary, mannerisms and communication style.
	 You may infer additional characteristics of your persona, but they should be consistent with known aspects of **your assigned persona**.
	 Unless specifically instructed otherwise, you should **never** introduce yourself.
	 Do not copy the other character's style, stay true to your own.
	 You must **never** break character.`,
}

func (b BotID) Opponent() BotID {
	return 1 - b
}

type Dialogue struct {
	scenario.Scenario
	ExtraPrompts []string
	Personas     [2]*persona.Persona
	Messages     []*Message
	Output       output.OutputStyle
	ctx          context.Context
	client       *api.Client
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

func NewDialogue(ctx context.Context, cfg *config.Config) *Dialogue {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		panic(err)
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
	}
}

func (c *Dialogue) AddMessage(botID BotID, content string) {
	c.Output.Render(c.Personas[botID].Name, content)
	c.Messages = append(c.Messages, &Message{botID, content})
}

func (c *Dialogue) FromPerspective(botID BotID) api.ChatRequest {
	prompts := make([]string, 0, 4+len(defaultPrompts)+len(c.ExtraPrompts))
	prompts = append(prompts, defaultPrompts...)
	prompts = append(prompts,
		c.Scenario.Scenario,
		c.Personas[botID].Persona,
		// fmt.Sprintf("The **other** persona/character in this scenario is %s.", c.Personas[botID.Opponent()].Name),
		c.Scenario.Roles[botID],
	)

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
		Stream:   &False,
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
	c.AddMessage(Persona1, strings.TrimSpace(c.Scenario.Opening))
	return c.SendRequest(Persona2)
}
