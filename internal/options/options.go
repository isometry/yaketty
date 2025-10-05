package options

type ModelOptions struct {
	NumCtx        int      `mapstructure:"num_ctx" default:"8192"`
	RepeatLastN   int      `mapstructure:"repeat_last_n" default:"-1"`
	RepeatPenalty float32  `mapstructure:"repeat_penalty" default:"1.1"`
	Temperature   float32  `mapstructure:"temperature" default:"0.8"`
	Stop          []string `mapstructure:"stop"`
	TopK          int      `mapstructure:"top_k" default:"40"`
	TopP          float32  `mapstructure:"top_p" default:"0.9"`
}

func (o ModelOptions) AsMap() map[string]any {
	return map[string]any{
		"num_ctx":        o.NumCtx,
		"repeat_last_n":  o.RepeatLastN,
		"repeat_penalty": o.RepeatPenalty,
		"temperature":    o.Temperature,
		"stop":           o.Stop,
		"top_k":          o.TopK,
		"top_p":          o.TopP,
	}
}
