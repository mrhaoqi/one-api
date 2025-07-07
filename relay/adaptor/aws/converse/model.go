package converse

// ConverseRequest represents the unified request structure for Converse API
// Note: This is just a placeholder struct. The actual AWS SDK types are used in the adapter.
type ConverseRequest struct {
	ModelId         string                 `json:"modelId"`
	Messages        interface{}            `json:"messages"`
	System          interface{}            `json:"system,omitempty"`
	InferenceConfig interface{}            `json:"inferenceConfig,omitempty"`
	ToolConfig      interface{}            `json:"toolConfig,omitempty"`
	GuardrailConfig interface{}            `json:"guardrailConfig,omitempty"`
	AdditionalModelRequestFields map[string]interface{} `json:"additionalModelRequestFields,omitempty"`
}

// ModelMapping contains all supported AWS Bedrock models
// Claude 3 series: Uses direct model IDs for regional support
// Claude 4 series: Uses inference profile IDs (required by AWS)
var ModelMapping = map[string]string{
	// Anthropic Claude models
	// Legacy models (use direct model IDs)
	"claude-instant-1.2":         "anthropic.claude-instant-v1",
	"claude-2.0":                 "anthropic.claude-v2",
	"claude-2.1":                 "anthropic.claude-v2:1",
	"claude-3-haiku-20240307":    "anthropic.claude-3-haiku-20240307-v1:0",
	"claude-3-sonnet-20240229":   "anthropic.claude-3-sonnet-20240229-v1:0",
	"claude-3-opus-20240229":     "anthropic.claude-3-opus-20240229-v1:0",
	"claude-3-5-sonnet-20240620": "anthropic.claude-3-5-sonnet-20240620-v1:0",

	// New models (use inference profile IDs - required by AWS)
	"claude-3-5-sonnet-20241022": "us.anthropic.claude-3-5-sonnet-20241022-v2:0",
	"claude-3-5-sonnet-latest":   "us.anthropic.claude-3-5-sonnet-20241022-v2:0",
	"claude-3-5-haiku-20241022":  "us.anthropic.claude-3-5-haiku-20241022-v1:0",
	"claude-3-7-sonnet-20250219": "us.anthropic.claude-3-7-sonnet-20250219-v1:0",
	"claude-opus-4-20250514":     "us.anthropic.claude-opus-4-20250514-v1:0",
	"claude-sonnet-4-20250514":   "us.anthropic.claude-sonnet-4-20250514-v1:0",

	// Meta Llama models
	"llama3-8b-8192":             "meta.llama3-8b-instruct-v1:0",
	"llama3-70b-8192":            "meta.llama3-70b-instruct-v1:0",
	"llama3-1-8b-instruct":       "meta.llama3-1-8b-instruct-v1:0",
	"llama3-1-70b-instruct":      "meta.llama3-1-70b-instruct-v1:0",
	"llama3-1-405b-instruct":     "meta.llama3-1-405b-instruct-v1:0",
	"llama3-2-1b-instruct":       "meta.llama3-2-1b-instruct-v1:0",
	"llama3-2-3b-instruct":       "meta.llama3-2-3b-instruct-v1:0",
	"llama3-2-11b-instruct":      "meta.llama3-2-11b-instruct-v1:0",
	"llama3-2-90b-instruct":      "meta.llama3-2-90b-instruct-v1:0",
	"llama3-3-70b-instruct":      "meta.llama3-3-70b-instruct-v1:0",
	"llama4-maverick-17b-instruct": "meta.llama4-maverick-17b-instruct-v1:0",
	"llama4-scout-17b-instruct":  "meta.llama4-scout-17b-instruct-v1:0",

	// Amazon Nova models
	"nova-micro":                 "amazon.nova-micro-v1:0",
	"nova-lite":                  "amazon.nova-lite-v1:0",
	"nova-pro":                   "amazon.nova-pro-v1:0",
	"nova-premier":               "amazon.nova-premier-v1:0",

	// Amazon Titan models
	"titan-text-express":         "amazon.titan-text-express-v1",
	"titan-text-lite":            "amazon.titan-text-lite-v1",
	"titan-text-premier":         "amazon.titan-text-premier-v1:0",

	// Cohere models
	"command-r":                  "cohere.command-r-v1:0",
	"command-r-plus":             "cohere.command-r-plus-v1:0",
	"command-light":              "cohere.command-light-text-v14",
	"command":                    "cohere.command-text-v14",

	// Mistral models
	"mistral-7b-instruct":        "mistral.mistral-7b-instruct-v0:2",
	"mistral-large-2402":         "mistral.mistral-large-2402-v1:0",
	"mistral-large-2407":         "mistral.mistral-large-2407-v1:0",
	"mistral-small-2402":         "mistral.mistral-small-2402-v1:0",
	"mixtral-8x7b-instruct":      "mistral.mixtral-8x7b-instruct-v0:1",
	"pixtral-large-2502":         "mistral.pixtral-large-2502-v1:0",

	// AI21 Jamba models
	"jamba-1-5-large":            "ai21.jamba-1-5-large-v1:0",
	"jamba-1-5-mini":             "ai21.jamba-1-5-mini-v1:0",

	// DeepSeek models
	"deepseek-r1":                "deepseek.r1-v1:0",
}

// GetAwsModelID returns the AWS Bedrock model ID for a given model name
func GetAwsModelID(modelName string) (string, bool) {
	awsModelID, exists := ModelMapping[modelName]
	return awsModelID, exists
}

// GetSupportedModels returns all supported model names
func GetSupportedModels() []string {
	models := make([]string, 0, len(ModelMapping))
	for model := range ModelMapping {
		models = append(models, model)
	}
	return models
}

// IsMultimodalModel checks if a model supports multimodal input
func IsMultimodalModel(modelName string) bool {
	multimodalModels := map[string]bool{
		"claude-3-haiku-20240307":    true,
		"claude-3-sonnet-20240229":   true,
		"claude-3-opus-20240229":     true,
		"claude-3-5-sonnet-20240620": true,
		"claude-3-5-sonnet-20241022": true,
		"claude-3-7-sonnet-20250219": true,
		"claude-opus-4-20250514":     true,
		"claude-sonnet-4-20250514":   true,
		"llama3-2-11b-instruct":      true,
		"llama3-2-90b-instruct":      true,
		"llama4-maverick-17b-instruct": true,
		"llama4-scout-17b-instruct":  true,
		"nova-lite":                  true,
		"nova-pro":                   true,
		"nova-premier":               true,
		"pixtral-large-2502":         true,
	}
	return multimodalModels[modelName]
}
