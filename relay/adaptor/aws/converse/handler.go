package converse

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/common"
	"github.com/songquanpeng/one-api/common/helper"
	"github.com/songquanpeng/one-api/common/random"
	relaymodel "github.com/songquanpeng/one-api/relay/model"
)

// StreamHandler handles streaming requests using AWS Bedrock Converse API
func StreamHandler(c *gin.Context, awsCli *bedrockruntime.Client) (usage *relaymodel.Usage, err *relaymodel.ErrorWithStatusCode) {
	fmt.Printf("[DEBUG] StreamHandler called\n")

	// Parse request
	request, parseErr := getRequest(c)
	if parseErr != nil {
		return nil, &relaymodel.ErrorWithStatusCode{
			Error: relaymodel.Error{
				Message: parseErr.Error(),
				Type:    "invalid_request_error",
			},
			StatusCode: http.StatusBadRequest,
		}
	}

	// Get AWS model ID from context
	awsModelID, exists := c.Get("aws_model_id")
	if !exists {
		return nil, &relaymodel.ErrorWithStatusCode{
			Error: relaymodel.Error{
				Message: "AWS model ID not found in context",
				Type:    "invalid_request_error",
			},
			StatusCode: http.StatusBadRequest,
		}
	}

	modelID := awsModelID.(string)
	fmt.Printf("[DEBUG] Request has %d messages\n", len(request.Messages))
	if request.Tools != nil {
		fmt.Printf("[DEBUG] Request has %d tools\n", len(request.Tools))
	}

	// For now, always use InvokeModel API for all streaming requests
	// This ensures consistent behavior and avoids the unimplemented Converse streaming API
	if request.Tools != nil && len(request.Tools) > 0 {
		fmt.Printf("[DEBUG] Tools detected, using InvokeModel API for streaming\n")
	} else {
		fmt.Printf("[DEBUG] No tools detected, using InvokeModel API for streaming (Converse streaming not implemented)\n")
	}
	return handleInvokeModelStreamRequest(c, awsCli, request, modelID)
}

func handleConverseStreamRequest(c *gin.Context, awsCli *bedrockruntime.Client, request *relaymodel.GeneralOpenAIRequest, modelID string) (usage *relaymodel.Usage, err *relaymodel.ErrorWithStatusCode) {
	// Implementation for Converse API streaming
	return nil, &relaymodel.ErrorWithStatusCode{
		Error: relaymodel.Error{
			Message: "Converse streaming not implemented",
			Type:    "not_implemented",
		},
		StatusCode: http.StatusNotImplemented,
	}
}

func handleInvokeModelStreamRequest(c *gin.Context, awsCli *bedrockruntime.Client, request *relaymodel.GeneralOpenAIRequest, modelID string) (usage *relaymodel.Usage, err *relaymodel.ErrorWithStatusCode) {
	fmt.Printf("[DEBUG] handleInvokeModelStreamRequest called with modelID: %s\n", modelID)

	if request.Tools != nil {
		fmt.Printf("[DEBUG] Request has %d tools for InvokeModelWithResponseStream API\n", len(request.Tools))
	}

	// Convert to Anthropic native format for InvokeModel API
	anthropicRequest, convertErr := convertToAnthropicNativeRequest(request, modelID)
	if convertErr != nil {
		return nil, &relaymodel.ErrorWithStatusCode{
			Error: relaymodel.Error{
				Message: convertErr.Error(),
				Type:    "invalid_request_error",
			},
			StatusCode: http.StatusBadRequest,
		}
	}

	requestBody, marshalErr := json.Marshal(anthropicRequest)
	if marshalErr != nil {
		return nil, &relaymodel.ErrorWithStatusCode{
			Error: relaymodel.Error{
				Message: marshalErr.Error(),
				Type:    "invalid_request_error",
			},
			StatusCode: http.StatusBadRequest,
		}
	}

	fmt.Printf("[DEBUG] Successfully converted request to Anthropic format for streaming, body length: %d bytes\n", len(requestBody))

	// Call InvokeModelWithResponseStream
	input := &bedrockruntime.InvokeModelWithResponseStreamInput{
		ModelId:     aws.String(modelID),
		Body:        requestBody,
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
	}

	stream, streamErr := awsCli.InvokeModelWithResponseStream(context.Background(), input)
	if streamErr != nil {
		return nil, &relaymodel.ErrorWithStatusCode{
			Error: relaymodel.Error{
				Message: fmt.Sprintf("InvokeModelWithResponseStream API调用失败 (模型: %s): %v", modelID, streamErr),
				Type:    "api_error",
			},
			StatusCode: http.StatusInternalServerError,
		}
	}

	// Process streaming response
	usage = &relaymodel.Usage{}
	id := fmt.Sprintf("chatcmpl-%s", random.GetUUID())
	model := modelID // Use modelID directly for now
	createdTime := helper.GetTimestamp()

	processErr := processInvokeModelStreamResponse(c, stream.GetStream(), id, model, int64(createdTime), modelID, usage)
	if processErr != nil {
		return usage, &relaymodel.ErrorWithStatusCode{
			Error: relaymodel.Error{
				Message: processErr.Error(),
				Type:    "api_error",
			},
			StatusCode: http.StatusInternalServerError,
		}
	}

	return usage, nil
}

func handleConverseRequest(c *gin.Context, awsCli *bedrockruntime.Client, request *relaymodel.GeneralOpenAIRequest, modelID string) (usage *relaymodel.Usage, err *relaymodel.ErrorWithStatusCode) {
	fmt.Printf("[DEBUG] handleConverseRequest called with modelID: %s\n", modelID)

	// For now, fall back to InvokeModel API until we fix the Converse API types
	fmt.Printf("[DEBUG] Falling back to InvokeModel API for non-streaming request\n")
	return handleInvokeModelRequest(c, awsCli, request, modelID)
}

func handleInvokeModelRequest(c *gin.Context, awsCli *bedrockruntime.Client, request *relaymodel.GeneralOpenAIRequest, modelID string) (usage *relaymodel.Usage, err *relaymodel.ErrorWithStatusCode) {
	fmt.Printf("[DEBUG] handleInvokeModelRequest called with modelID: %s\n", modelID)

	// Convert to Anthropic native format for InvokeModel API
	anthropicRequest, convertErr := convertToAnthropicNativeRequest(request, modelID)
	if convertErr != nil {
		return nil, &relaymodel.ErrorWithStatusCode{
			Error: relaymodel.Error{
				Message: convertErr.Error(),
				Type:    "invalid_request_error",
			},
			StatusCode: http.StatusBadRequest,
		}
	}

	// Remove stream parameter for non-streaming
	delete(anthropicRequest, "stream")

	requestBody, marshalErr := json.Marshal(anthropicRequest)
	if marshalErr != nil {
		return nil, &relaymodel.ErrorWithStatusCode{
			Error: relaymodel.Error{
				Message: marshalErr.Error(),
				Type:    "invalid_request_error",
			},
			StatusCode: http.StatusBadRequest,
		}
	}

	fmt.Printf("[DEBUG] Successfully converted request to Anthropic format for non-streaming, body length: %d bytes\n", len(requestBody))

	// Call InvokeModel
	input := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(modelID),
		Body:        requestBody,
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
	}

	response, invokeErr := awsCli.InvokeModel(context.Background(), input)
	if invokeErr != nil {
		return nil, &relaymodel.ErrorWithStatusCode{
			Error: relaymodel.Error{
				Message: fmt.Sprintf("InvokeModel API调用失败 (模型: %s): %v", modelID, invokeErr),
				Type:    "api_error",
			},
			StatusCode: http.StatusInternalServerError,
		}
	}

	// Parse response
	var anthropicResponse map[string]interface{}
	if err := json.Unmarshal(response.Body, &anthropicResponse); err != nil {
		return nil, &relaymodel.ErrorWithStatusCode{
			Error: relaymodel.Error{
				Message: fmt.Sprintf("Failed to parse response: %v", err),
				Type:    "api_error",
			},
			StatusCode: http.StatusInternalServerError,
		}
	}

	// Convert to OpenAI format and send response
	openaiResponse := convertAnthropicResponseToOpenAI(anthropicResponse, modelID)

	// Extract usage information
	usage = &relaymodel.Usage{}
	if usageInfo, ok := anthropicResponse["usage"].(map[string]interface{}); ok {
		if inputTokens, ok := usageInfo["input_tokens"].(float64); ok {
			usage.PromptTokens = int(inputTokens)
		}
		if outputTokens, ok := usageInfo["output_tokens"].(float64); ok {
			usage.CompletionTokens = int(outputTokens)
		}
		usage.TotalTokens = usage.PromptTokens + usage.CompletionTokens
	}

	c.JSON(http.StatusOK, openaiResponse)
	return usage, nil
}

// convertAnthropicResponseToOpenAI converts Anthropic response to OpenAI format
func convertAnthropicResponseToOpenAI(anthropicResponse map[string]interface{}, modelID string) map[string]interface{} {
	openaiResponse := map[string]interface{}{
		"id":      fmt.Sprintf("chatcmpl-%s", random.GetUUID()),
		"object":  "chat.completion",
		"created": helper.GetTimestamp(),
		"model":   modelID,
	}

	var choices []map[string]interface{}
	var message map[string]interface{}

	// Extract content
	if content, ok := anthropicResponse["content"].([]interface{}); ok {
		var textContent string
		var toolCalls []map[string]interface{}

		for _, item := range content {
			if contentItem, ok := item.(map[string]interface{}); ok {
				contentType, _ := contentItem["type"].(string)

				switch contentType {
				case "text":
					if text, ok := contentItem["text"].(string); ok {
						textContent = text
					}
				case "tool_use":
					if toolId, ok := contentItem["id"].(string); ok {
						if toolName, ok := contentItem["name"].(string); ok {
							var arguments string
							if input, ok := contentItem["input"]; ok {
								if inputBytes, err := json.Marshal(input); err == nil {
									arguments = string(inputBytes)
								}
							}

							toolCall := map[string]interface{}{
								"id":   toolId,
								"type": "function",
								"function": map[string]interface{}{
									"name":      toolName,
									"arguments": arguments,
								},
							}
							toolCalls = append(toolCalls, toolCall)
						}
					}
				}
			}
		}

		message = map[string]interface{}{
			"role": "assistant",
		}

		if len(toolCalls) > 0 {
			message["tool_calls"] = toolCalls
			if textContent != "" {
				message["content"] = textContent
			}
		} else {
			message["content"] = textContent
		}
	}

	choice := map[string]interface{}{
		"index":   0,
		"message": message,
	}

	// Map stop reason
	if stopReason, ok := anthropicResponse["stop_reason"].(string); ok {
		choice["finish_reason"] = mapStopReason(stopReason)
	} else {
		choice["finish_reason"] = "stop"
	}

	choices = append(choices, choice)
	openaiResponse["choices"] = choices

	// Add usage information
	if usage, ok := anthropicResponse["usage"].(map[string]interface{}); ok {
		openaiResponse["usage"] = map[string]interface{}{
			"prompt_tokens":     usage["input_tokens"],
			"completion_tokens": usage["output_tokens"],
			"total_tokens":      int(usage["input_tokens"].(float64)) + int(usage["output_tokens"].(float64)),
		}
	}

	return openaiResponse
}

// convertToAnthropicNativeRequest converts OpenAI format to Anthropic native format
func convertToAnthropicNativeRequest(request *relaymodel.GeneralOpenAIRequest, modelID string) (map[string]interface{}, error) {
	anthropicRequest := map[string]interface{}{
		"anthropic_version": "bedrock-2023-05-31",
		"max_tokens":        4096,
	}

	// Add tools if present
	if request.Tools != nil && len(request.Tools) > 0 {
		var anthropicTools []map[string]interface{}
		for _, tool := range request.Tools {
			if tool.Type == "function" {
				anthropicTool := map[string]interface{}{
					"name":        tool.Function.Name,
					"description": tool.Function.Description,
				}
				if tool.Function.Parameters != nil {
					anthropicTool["input_schema"] = tool.Function.Parameters
				}
				anthropicTools = append(anthropicTools, anthropicTool)
			}
		}
		anthropicRequest["tools"] = anthropicTools
	}

	// Convert messages with proper tool message handling
	var anthropicMessages []map[string]interface{}
	var systemPrompt string

	for _, msg := range request.Messages {
		if msg.Role == "system" {
			systemPrompt = msg.StringContent()
			continue
		}

		// Handle tool messages properly
		if msg.Role == "tool" {
			// Convert tool message to Anthropic format
			anthropicMsg := map[string]interface{}{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type":        "tool_result",
						"tool_use_id": msg.ToolCallId,
						"content":     msg.StringContent(),
					},
				},
			}
			anthropicMessages = append(anthropicMessages, anthropicMsg)
			continue
		}

		// Handle assistant messages with tool calls
		if msg.Role == "assistant" && msg.ToolCalls != nil && len(msg.ToolCalls) > 0 {
			var contentBlocks []map[string]interface{}

			// Add text content if present
			if msg.StringContent() != "" {
				contentBlocks = append(contentBlocks, map[string]interface{}{
					"type": "text",
					"text": msg.StringContent(),
				})
			}

			// Add tool use blocks
			for _, toolCall := range msg.ToolCalls {
				if toolCall.Type == "function" {
					var input interface{}
					if argumentsStr, ok := toolCall.Function.Arguments.(string); ok && argumentsStr != "" {
						if err := json.Unmarshal([]byte(argumentsStr), &input); err != nil {
							// If JSON parsing fails, use the raw string
							input = argumentsStr
						}
					} else {
						// If Arguments is not a string, use it directly
						input = toolCall.Function.Arguments
					}

					contentBlocks = append(contentBlocks, map[string]interface{}{
						"type":  "tool_use",
						"id":    toolCall.Id,
						"name":  toolCall.Function.Name,
						"input": input,
					})
				}
			}

			anthropicMsg := map[string]interface{}{
				"role":    "assistant",
				"content": contentBlocks,
			}
			anthropicMessages = append(anthropicMessages, anthropicMsg)
			continue
		}

		// Handle regular messages
		anthropicMsg := map[string]interface{}{
			"role":    msg.Role,
			"content": msg.StringContent(),
		}
		anthropicMessages = append(anthropicMessages, anthropicMsg)
	}

	anthropicRequest["messages"] = anthropicMessages

	if systemPrompt != "" {
		anthropicRequest["system"] = systemPrompt
	}

	return anthropicRequest, nil
}

// TODO: Implement convertToConverseInput when AWS SDK types are properly resolved

func getRequest(c *gin.Context) (*relaymodel.GeneralOpenAIRequest, error) {
	var request relaymodel.GeneralOpenAIRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		return nil, err
	}
	return &request, nil
}

// TODO: Implement Converse API conversion functions when AWS SDK types are properly resolved

// Handler handles non-streaming requests
func Handler(c *gin.Context, awsCli *bedrockruntime.Client) (usage *relaymodel.Usage, err *relaymodel.ErrorWithStatusCode) {
	fmt.Printf("[DEBUG] Handler called (non-streaming)\n")

	// Parse request
	request, parseErr := getRequest(c)
	if parseErr != nil {
		return nil, &relaymodel.ErrorWithStatusCode{
			Error: relaymodel.Error{
				Message: parseErr.Error(),
				Type:    "invalid_request_error",
			},
			StatusCode: http.StatusBadRequest,
		}
	}

	// Get AWS model ID from context
	awsModelID, exists := c.Get("aws_model_id")
	if !exists {
		return nil, &relaymodel.ErrorWithStatusCode{
			Error: relaymodel.Error{
				Message: "AWS model ID not found in context",
				Type:    "invalid_request_error",
			},
			StatusCode: http.StatusBadRequest,
		}
	}

	modelID := awsModelID.(string)
	fmt.Printf("[DEBUG] Non-streaming request for model: %s\n", modelID)

	// For now, always use InvokeModel API for all requests
	// This ensures consistent behavior and proper tool message handling
	if request.Tools != nil && len(request.Tools) > 0 {
		fmt.Printf("[DEBUG] Tools detected, using InvokeModel API for non-streaming\n")
	} else {
		fmt.Printf("[DEBUG] No tools detected, using InvokeModel API for non-streaming (consistent with streaming)\n")
	}
	return handleInvokeModelRequest(c, awsCli, request, modelID)
}

// processInvokeModelStreamResponse processes streaming response from InvokeModelWithResponseStream
func processInvokeModelStreamResponse(c *gin.Context, stream *bedrockruntime.InvokeModelWithResponseStreamEventStream, id, model string, createdTime int64, modelID string, usage *relaymodel.Usage) error {
	// Set streaming headers
	common.SetEventStreamHeaders(c)

	for event := range stream.Events() {
		switch e := event.(type) {
		case *types.ResponseStreamMemberChunk:
			if e.Value.Bytes != nil {
				// Parse Anthropic streaming format
				var anthropicEvent map[string]interface{}
				if err := json.Unmarshal(e.Value.Bytes, &anthropicEvent); err != nil {
					continue // Skip malformed events
				}

				// Convert to OpenAI format
				openaiChunk := convertAnthropicStreamToOpenAI(anthropicEvent, id, model, createdTime)
				if openaiChunk != nil {
					// Send OpenAI format streaming response
					responseBytes, _ := json.Marshal(openaiChunk)
					c.Writer.Write([]byte("data: "))
					c.Writer.Write(responseBytes)
					c.Writer.Write([]byte("\n\n"))
					c.Writer.Flush()
				}

				// Extract usage information from message_stop event
				if eventType, ok := anthropicEvent["type"].(string); ok && eventType == "message_stop" {
					if metrics, ok := anthropicEvent["amazon-bedrock-invocationMetrics"].(map[string]interface{}); ok {
						if inputTokens, ok := metrics["inputTokenCount"].(float64); ok {
							usage.PromptTokens = int(inputTokens)
						}
						if outputTokens, ok := metrics["outputTokenCount"].(float64); ok {
							usage.CompletionTokens = int(outputTokens)
						}
						usage.TotalTokens = usage.PromptTokens + usage.CompletionTokens
					}
				}
			}
		default:
			// Handle other event types or errors
			continue
		}
	}

	// Send final [DONE] message
	c.Writer.Write([]byte("data: [DONE]\n\n"))
	c.Writer.Flush()

	return nil
}

// convertAnthropicStreamToOpenAI converts Anthropic streaming events to OpenAI format
func convertAnthropicStreamToOpenAI(anthropicEvent map[string]interface{}, id, model string, createdTime int64) map[string]interface{} {
	eventType, ok := anthropicEvent["type"].(string)
	if !ok {
		return nil
	}

	switch eventType {
	case "message_start":
		// Message start - send role
		return map[string]interface{}{
			"id":      id,
			"object":  "chat.completion.chunk",
			"created": createdTime,
			"model":   model,
			"choices": []map[string]interface{}{
				{
					"index": 0,
					"delta": map[string]interface{}{
						"role": "assistant",
					},
					"finish_reason": nil,
				},
			},
		}

	case "content_block_start":
		// Start of content block - check if it's a tool use
		if contentBlock, ok := anthropicEvent["content_block"].(map[string]interface{}); ok {
			if blockType, ok := contentBlock["type"].(string); ok && blockType == "tool_use" {
				// Tool call start
				if toolName, ok := contentBlock["name"].(string); ok {
					if toolId, ok := contentBlock["id"].(string); ok {
						return map[string]interface{}{
							"id":      id,
							"object":  "chat.completion.chunk",
							"created": createdTime,
							"model":   model,
							"choices": []map[string]interface{}{
								{
									"index": 0,
									"delta": map[string]interface{}{
										"tool_calls": []map[string]interface{}{
											{
												"index": 0,
												"id":    toolId,
												"type":  "function",
												"function": map[string]interface{}{
													"name":      toolName,
													"arguments": "",
												},
											},
										},
									},
									"finish_reason": nil,
								},
							},
						}
					}
				}
			}
		}
		return nil

	case "content_block_delta":
		// Content delta
		if delta, ok := anthropicEvent["delta"].(map[string]interface{}); ok {
			deltaType, ok := delta["type"].(string)
			if !ok {
				return nil
			}

			switch deltaType {
			case "text_delta":
				// Text content delta
				if text, ok := delta["text"].(string); ok {
					return map[string]interface{}{
						"id":      id,
						"object":  "chat.completion.chunk",
						"created": createdTime,
						"model":   model,
						"choices": []map[string]interface{}{
							{
								"index": 0,
								"delta": map[string]interface{}{
									"content": text,
								},
								"finish_reason": nil,
							},
						},
					}
				}

			case "input_json_delta":
				// Tool call input delta
				if partialJson, ok := delta["partial_json"].(string); ok {
					return map[string]interface{}{
						"id":      id,
						"object":  "chat.completion.chunk",
						"created": createdTime,
						"model":   model,
						"choices": []map[string]interface{}{
							{
								"index": 0,
								"delta": map[string]interface{}{
									"tool_calls": []map[string]interface{}{
										{
											"index": 0,
											"function": map[string]interface{}{
												"arguments": partialJson,
											},
										},
									},
								},
								"finish_reason": nil,
							},
						},
					}
				}
			}
		}
		return nil

	case "content_block_stop":
		// Content block stop - no output needed
		return nil

	case "message_delta":
		// Message delta - check for stop reason
		if delta, ok := anthropicEvent["delta"].(map[string]interface{}); ok {
			if stopReason, ok := delta["stop_reason"].(string); ok {
				finishReason := mapStopReason(stopReason)
				return map[string]interface{}{
					"id":      id,
					"object":  "chat.completion.chunk",
					"created": createdTime,
					"model":   model,
					"choices": []map[string]interface{}{
						{
							"index":         0,
							"delta":         map[string]interface{}{},
							"finish_reason": finishReason,
						},
					},
				}
			}
		}
		return nil

	case "message_stop":
		// Message stop - send final finish reason
		return map[string]interface{}{
			"id":      id,
			"object":  "chat.completion.chunk",
			"created": createdTime,
			"model":   model,
			"choices": []map[string]interface{}{
				{
					"index":         0,
					"delta":         map[string]interface{}{},
					"finish_reason": "stop",
				},
			},
		}

	default:
		// Unknown event type
		return nil
	}
}

// mapStopReason maps Anthropic stop reasons to OpenAI format
func mapStopReason(anthropicReason string) string {
	switch anthropicReason {
	case "end_turn":
		return "stop"
	case "max_tokens":
		return "length"
	case "stop_sequence":
		return "stop"
	case "tool_use":
		return "tool_calls"
	default:
		return "stop"
	}
}