# AWS Bedrock 工具调用支持状态

## 🚨 当前状态

**工具调用功能暂时禁用**

由于AWS SDK v2的document interface要求，工具调用功能暂时不可用。

## 🔍 技术原因

### AWS SDK Document Interface 要求

AWS Bedrock Converse API的工具配置需要实现特定的`document.Interface`：

```go
type ToolInputSchemaMemberJson struct {
    Value document.Interface  // 需要实现 MarshalSmithyDocument 方法
}
```

### 当前限制

1. **接口要求**：`document.Interface`需要实现`MarshalSmithyDocument`方法
2. **类型转换**：OpenAI的`interface{}`类型无法直接转换为AWS的document interface
3. **复杂性**：正确实现document interface需要深入了解AWS SDK内部机制

## ✅ 当前支持的功能

### 完全支持
- ✅ **基础对话**：文本对话完全正常
- ✅ **流式对话**：实时流式响应
- ✅ **多模态输入**：文本+图像支持
- ✅ **系统消息**：系统提示正常工作
- ✅ **参数控制**：temperature, top_p, max_tokens, stop_sequences

### 暂时不支持
- 🚧 **工具调用**：Function Calling暂时禁用
- 🚧 **函数调用**：需要等待document interface实现

## 🛠️ 解决方案

### 临时解决方案

**如果您需要工具调用功能，建议：**

1. **使用其他模型**：
   - OpenAI GPT-4系列（完整工具支持）
   - Anthropic Claude（通过官方API）
   - 其他支持工具调用的模型

2. **分步处理**：
   - 先用AWS Claude进行对话
   - 需要工具时切换到其他模型

### 长期解决方案

**我们正在开发：**

1. **Document Interface实现**：
   ```go
   type ParameterDocument struct {
       data map[string]interface{}
   }
   
   func (p *ParameterDocument) MarshalSmithyDocument() ([]byte, error) {
       // 实现正确的序列化
   }
   ```

2. **工具转换器**：
   - OpenAI工具格式 → AWS document interface
   - 参数验证和类型转换
   - 错误处理和回退机制

## 📋 错误处理

### 当前行为

如果请求包含工具：
```json
{
  "model": "claude-3-5-sonnet-latest",
  "messages": [...],
  "tools": [...]  // 这会触发错误
}
```

**返回错误**：
```
工具调用功能暂时不支持，正在开发中。请使用不带工具的请求
```

### 建议的请求格式

**正确的请求**（不包含工具）：
```json
{
  "model": "claude-3-5-sonnet-latest",
  "messages": [
    {"role": "user", "content": "Hello!"}
  ],
  "temperature": 0.7,
  "max_tokens": 1000
}
```

## 🧪 测试建议

### 基础功能测试

```bash
# 测试基础对话
curl -X POST http://localhost:3000/v1/chat/completions \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "claude-3-5-sonnet-latest",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

### 流式对话测试

```bash
# 测试流式响应
curl -X POST http://localhost:3000/v1/chat/completions \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "claude-3-5-sonnet-latest",
    "messages": [{"role": "user", "content": "Tell me a story"}],
    "stream": true
  }'
```

### 多模态测试

```bash
# 测试图像输入
curl -X POST http://localhost:3000/v1/chat/completions \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "claude-3-5-sonnet-latest",
    "messages": [{
      "role": "user",
      "content": [
        {"type": "text", "text": "What is in this image?"},
        {"type": "image_url", "image_url": {"url": "data:image/jpeg;base64,..."}}
      ]
    }]
  }'
```

## 📅 开发计划

### 短期目标（1-2周）
- 🔧 研究AWS SDK document interface最佳实践
- 🔧 实现基础的document interface wrapper
- 🔧 添加简单工具调用支持

### 中期目标（1个月）
- 🚀 完整的工具调用功能
- 🚀 复杂参数schema支持
- 🚀 工具选择策略实现

### 长期目标（持续）
- ⭐ 性能优化
- ⭐ 错误处理增强
- ⭐ 更多AWS Bedrock模型支持

## 💡 当前最佳实践

1. **使用AWS Claude进行**：
   - 基础文本对话
   - 流式对话
   - 多模态输入
   - 内容生成

2. **使用其他模型进行**：
   - 需要工具调用的任务
   - 复杂的函数调用
   - API集成任务

3. **混合使用策略**：
   - 主要对话用AWS Claude（成本优势）
   - 工具调用用OpenAI GPT-4（功能完整）

## 🆘 获取帮助

如果您有工具调用的紧急需求：

1. **联系开发团队**：报告具体的使用场景
2. **提供反馈**：帮助我们优先开发最需要的功能
3. **参与测试**：当工具支持就绪时，帮助测试验证

现在AWS Bedrock适配器在基础对话功能上已经完全稳定，工具调用功能正在积极开发中！
