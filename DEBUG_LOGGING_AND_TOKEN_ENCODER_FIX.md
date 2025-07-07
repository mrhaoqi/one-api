# 调试日志和Token编码器修复

## 🚨 问题分析

### Token编码器错误
```
[ERROR] failed to get token encoder for model claude-3-5-sonnet-latest: no encoding for model claude-3-5-sonnet-latest, using encoder for gpt-3.5-turbo
```

**问题原因**：
- one-api主要为OpenAI模型设计
- Claude模型没有专用的tiktoken编码器
- 系统自动回退到GPT-3.5-turbo编码器

## ✅ 修复内容

### 1. 添加详细调试日志

在关键位置添加了调试日志，帮助追踪工具调用流程：

#### Handler入口日志
```go
fmt.Printf("[DEBUG] Handler called with modelID: %s\n", modelID)
fmt.Printf("[DEBUG] Request has %d messages\n", len(request.Messages))
fmt.Printf("[DEBUG] Request has %d tools\n", len(request.Tools))
if len(request.Tools) > 0 {
    fmt.Printf("[DEBUG] Tools in request:\n")
    for i, tool := range request.Tools {
        fmt.Printf("[DEBUG]   Tool %d: name=%s, type=%s, description=%s\n", 
            i, tool.Function.Name, tool.Type, tool.Function.Description)
    }
}
```

#### 工具检测日志
```go
fmt.Printf("[DEBUG] convertToConverseRequest called with modelID: %s\n", modelID)
fmt.Printf("[DEBUG] Request has %d tools\n", len(request.Tools))
if len(request.Tools) > 0 {
    fmt.Printf("[DEBUG] Tools detected, will switch to InvokeModel API\n")
    for i, tool := range request.Tools {
        fmt.Printf("[DEBUG] Tool %d: name=%s, type=%s\n", i, tool.Function.Name, tool.Type)
    }
}
```

#### API切换日志
```go
fmt.Printf("[DEBUG] Switching to InvokeModel API for tool calls with modelID: %s\n", invokeModelID)
fmt.Printf("[DEBUG] InvokeModel API call succeeded\n")
```

#### InvokeModel处理日志
```go
fmt.Printf("[DEBUG] handleInvokeModelRequest called with modelID: %s\n", modelID)
fmt.Printf("[DEBUG] Request has %d tools for InvokeModel API\n", len(request.Tools))
fmt.Printf("[DEBUG] Successfully converted request to Anthropic format, body length: %d bytes\n", len(requestBody))
```

### 2. 修复Token编码器错误

为Claude模型添加了编码器映射：

```go
for model := range billingratio.ModelRatio {
    if strings.HasPrefix(model, "gpt-3.5") {
        tokenEncoderMap[model] = gpt35TokenEncoder
    } else if strings.HasPrefix(model, "gpt-4o") {
        tokenEncoderMap[model] = gpt4oTokenEncoder
    } else if strings.HasPrefix(model, "gpt-4") {
        tokenEncoderMap[model] = gpt4TokenEncoder
    } else if strings.HasPrefix(model, "claude") {
        // 为Claude模型使用gpt-3.5-turbo编码器（近似）
        tokenEncoderMap[model] = gpt35TokenEncoder
    } else {
        tokenEncoderMap[model] = nil
    }
}
```

## 🔍 调试日志输出示例

### 无工具调用的请求
```
[DEBUG] Handler called with modelID: us.anthropic.claude-3-5-sonnet-20241022-v2:0
[DEBUG] Request has 1 messages
[DEBUG] Request has 0 tools
[DEBUG] convertToConverseRequest called with modelID: us.anthropic.claude-3-5-sonnet-20241022-v2:0
[DEBUG] Request has 0 tools
```

### 有工具调用的请求
```
[DEBUG] Handler called with modelID: us.anthropic.claude-3-5-sonnet-20241022-v2:0
[DEBUG] Request has 1 messages
[DEBUG] Request has 1 tools
[DEBUG] Tools in request:
[DEBUG]   Tool 0: name=sequentialthinking_sequential-thinking, type=function, description=A detailed tool for dynamic and reflective problem-solving
[DEBUG] convertToConverseRequest called with modelID: us.anthropic.claude-3-5-sonnet-20241022-v2:0
[DEBUG] Request has 1 tools
[DEBUG] Tools detected, will switch to InvokeModel API
[DEBUG] Tool 0: name=sequentialthinking_sequential-thinking, type=function
[DEBUG] Tools detected in convertToConverseRequest, switching to InvokeModel API
[DEBUG] Number of tools: 1
[DEBUG] Switching to InvokeModel API for tool calls with modelID: us.anthropic.claude-3-5-sonnet-20241022-v2:0
[DEBUG] handleInvokeModelRequest called with modelID: us.anthropic.claude-3-5-sonnet-20241022-v2:0
[DEBUG] Request has 1 tools for InvokeModel API
[DEBUG] Successfully converted request to Anthropic format, body length: 1234 bytes
[DEBUG] InvokeModel API call succeeded
```

## 🧪 测试和调试

### 重新编译和启动
```bash
go build -o one-api
./one-api --port 3000
```

### 测试工具调用
发送包含工具的请求，观察控制台输出：

```json
{
  "model": "claude-3-5-sonnet-latest",
  "messages": [
    {"role": "user", "content": "深入思考一个问题，有关于设计一个cad图纸生成3d模型的方案"}
  ],
  "tools": [
    {
      "type": "function",
      "function": {
        "name": "sequentialthinking_sequential-thinking",
        "description": "A detailed tool for dynamic and reflective problem-solving through thoughts",
        "parameters": {
          "type": "object",
          "properties": {
            "thought": {"type": "string"},
            "nextThoughtNeeded": {"type": "boolean"},
            "thoughtNumber": {"type": "integer"},
            "totalThoughts": {"type": "integer"}
          },
          "required": ["thought", "nextThoughtNeeded", "thoughtNumber", "totalThoughts"]
        }
      }
    }
  ]
}
```

### 预期的调试输出

1. **请求接收**：显示模型ID、消息数量、工具数量
2. **工具检测**：显示检测到的工具详情
3. **API切换**：显示从Converse切换到InvokeModel
4. **请求转换**：显示Anthropic格式转换成功
5. **API调用**：显示InvokeModel API调用结果

## 🔧 故障排除

### 如果工具没有被调用

检查调试日志中的关键信息：

1. **工具数量**：`Request has X tools`
   - 如果是0，说明请求中没有工具定义
   - 如果>0，说明工具被正确识别

2. **API切换**：`Switching to InvokeModel API`
   - 如果没有这条日志，说明没有检测到工具
   - 如果有，说明API切换正常

3. **请求转换**：`Successfully converted request to Anthropic format`
   - 如果失败，会显示具体错误信息
   - 如果成功，显示请求体大小

4. **API调用结果**：`InvokeModel API call succeeded/failed`
   - 成功：工具调用应该正常工作
   - 失败：会显示具体错误信息

### 常见问题诊断

1. **工具定义缺失**：
   ```
   [DEBUG] Request has 0 tools
   ```
   → 检查请求中是否包含tools字段

2. **工具名称错误**：
   ```
   [DEBUG] Tool 0: name=wrong-name, type=function
   ```
   → 检查工具名称是否与MCP服务器一致

3. **API调用失败**：
   ```
   [DEBUG] InvokeModel API call failed: ValidationException
   ```
   → 检查请求格式和模型ID

## 💡 优化建议

### 生产环境
在生产环境中，可以通过环境变量控制调试日志：

```go
if os.Getenv("DEBUG_AWS_BEDROCK") == "true" {
    fmt.Printf("[DEBUG] ...")
}
```

### 日志级别
可以考虑使用结构化日志：

```go
logger.Debug("AWS Bedrock request", 
    "modelID", modelID,
    "toolCount", len(request.Tools),
    "messageCount", len(request.Messages))
```

现在您可以通过详细的调试日志准确诊断工具调用问题，并且不再看到Token编码器的错误信息！🎉
