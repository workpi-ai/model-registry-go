package registry

const (
	ProviderTypeAPI          ProviderType = "api"
	ProviderTypeSubscription ProviderType = "subscription"
	ProviderTypeAggregator   ProviderType = "aggregator"
)

const (
	AuthTypeAPIKey         AuthType = "api_key"
	AuthTypeOAuth2         AuthType = "oauth2"
	AuthTypeAWSCredentials AuthType = "aws_credentials"
)

const (
	APIFormatOpenAI    APIFormat = "openai"
	APIFormatAnthropic APIFormat = "anthropic"
	APIFormatGemini    APIFormat = "gemini"
	APIFormatCodex     APIFormat = "codex"
	APIFormatBedrock   APIFormat = "bedrock"
)

const (
	ProviderNameOpenAI       = "openai"
	ProviderNameOpenAISub    = "openai-sub"
	ProviderNameAnthropic    = "anthropic"
	ProviderNameAnthropicSub = "anthropic-sub"
	ProviderNameGemini       = "gemini"
	ProviderNameGeminiSub    = "gemini-sub"
	ProviderNameDeepSeek     = "deepseek"
	ProviderNameXai          = "xai"
	ProviderNameZai          = "zai"
	ProviderNameZaiCN        = "zai-cn"
	ProviderNameMoonshot     = "moonshot"
	ProviderNameMoonshotCN   = "moonshot-cn"
	ProviderNameMiniMax      = "minimax"
	ProviderNameMiniMaxCN    = "minimax-cn"
	ProviderNameQwen         = "qwen"
	ProviderNameQwenCN       = "qwen-cn"
	ProviderNameOpenRouter   = "openrouter"
	ProviderNameBedrock      = "bedrock"
)

const (
	repoOwner    = "workpi-ai"
	repoName     = "model-registry"
	providersDir = "providers"
	versionFile  = "metadata.json"
	providerYAML = "provider.yaml"
	yamlExt      = ".yaml"
)
