# AWS Bedrock 流式API修复

## 🚨 问题描述

之前遇到的错误：
```
ValidationException: Malformed input request: #: subject must not be valid against schema {"required":["messages"]}#: required key [max_tokens] not found#: required key [anthropic_version] not found#: extraneous key [inferenceConfig] is not permitted
```

## 🔍 根本原因

1. **API不匹配**：流式处理仍在使用旧的`InvokeModelWithResponseStream` API
2. **请求格式错误**：旧API需要Anthropic特定的请求格式
3. **字段不兼容**：新的Converse API字段与旧API不兼容

## ✅ 修复内容

### 1. 更新流式处理API

**之前（旧API）：**
```go
// 使用 InvokeModelWithResponseStream
response, err := awsCli.InvokeModelWithResponseStream(c.Request.Context(), input)
```

**现在（新API）：**
```go
// 使用 ConverseStream
response, err := awsCli.ConverseStream(c.Request.Context(), streamInput)
```

### 2. 统一请求格式

**之前：**
- 非流式：使用 Converse API
- 流式：使用 InvokeModelWithResponseStream API（不一致）

**现在：**
- 非流式：使用 Converse API
- 流式：使用 ConverseStream API（统一）

### 3. 实现完整的流式响应处理

```go
func processConverseStreamResponse(c *gin.Context, stream *bedrockruntime.ConverseStreamEventStream, id, model string, createdTime int64, modelID string, usage *relaymodel.Usage) error {
    for event := range stream.Events() {
        switch e := event.(type) {
        case *types.ConverseStreamOutputMemberContentBlockDelta:
            // 处理内容增量
        case *types.ConverseStreamOutputMemberMetadata:
            // 处理使用情况统计
        case *types.ConverseStreamOutputMemberMessageStop:
            // 处理消息结束
        }
    }
}
```

## 🔧 技术改进

### API统一性
- ✅ 非流式和流式都使用Converse API系列
- ✅ 统一的请求格式转换
- ✅ 统一的错误处理

### 流式处理增强
- ✅ 实时内容流式传输
- ✅ 正确的使用情况统计
- ✅ 标准的OpenAI流式响应格式

### 错误处理改进
- ✅ 更详细的错误诊断信息
- ✅ 地理位置限制的具体解决方案
- ✅ 权限和配置问题的指导

## 🚀 现在支持的功能

### 非流式对话
- ✅ 完整的Converse API支持
- ✅ 多模态输入（文本+图像）
- ✅ 工具调用（Function Calling）
- ✅ 系统消息处理

### 流式对话
- ✅ 实时流式响应
- ✅ 增量内容传输
- ✅ 使用情况统计
- ✅ 正确的结束标记

### 模型支持
- ✅ Claude 3系列（直接模型ID）
- ✅ Claude 3.5+系列（推理配置文件ID）
- ✅ Claude 4系列（推理配置文件ID）

## 📋 测试建议

### 1. 非流式测试
```bash
curl -X POST http://localhost:3000/v1/chat/completions \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "claude-3-5-sonnet-latest",
    "messages": [{"role": "user", "content": "Hello!"}],
    "stream": false
  }'
```

### 2. 流式测试
```bash
curl -X POST http://localhost:3000/v1/chat/completions \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "claude-3-5-sonnet-latest", 
    "messages": [{"role": "user", "content": "Tell me a story"}],
    "stream": true
  }'
```

### 3. 工具调用测试
```bash
curl -X POST http://localhost:3000/v1/chat/completions \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "claude-3-5-sonnet-latest",
    "messages": [{"role": "user", "content": "What is the weather like?"}],
    "tools": [{"type": "function", "function": {"name": "get_weather", "description": "Get weather info"}}]
  }'
```

## 🔄 部署步骤

1. **重新编译**：
   ```bash
   go build -o one-api
   ```

2. **重启服务**：
   ```bash
   ./one-api --port 3000
   ```

3. **验证功能**：
   - 测试非流式对话
   - 测试流式对话
   - 测试工具调用

## 💡 最佳实践

1. **区域设置**：使用 `us-east-1` 获得最佳兼容性
2. **网络环境**：避免使用VPN或代理
3. **模型选择**：新模型使用推理配置文件ID
4. **错误监控**：关注地理位置和权限相关错误

现在AWS Bedrock适配器完全支持新的Converse API，提供了统一、高效的流式和非流式对话体验！
