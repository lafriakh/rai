package internal

type Config struct {
	Gemini    ProviderConfig `mapstructure:"gemini"`
	Anthropic ProviderConfig `mapstructure:"anthropic"`
}

type ProviderConfig struct {
	ModelID          string `mapstructure:"model_id"`
	APIKey           string `mapstructure:"api_key"`
	SystemPromptPath string `mapstructure:"system_prompt_path"`
}
