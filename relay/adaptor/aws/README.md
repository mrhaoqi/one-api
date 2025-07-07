# AWS Bedrock Adapter

这是one-api的AWS Bedrock适配器，使用最新的Converse API统一接口。

## 特性

- ✅ **统一API**: 使用AWS Bedrock的新一代Converse API
- ✅ **广泛支持**: 支持70+个模型，包括Claude、Llama、Nova等
- ✅ **多模态**: 支持文本+图像输入
- ✅ **工具调用**: 完整的function calling支持
- ✅ **流式响应**: 支持实时流式输出
- ✅ **模型映射**: 友好的模型名称映射

## 支持的模型

### Anthropic Claude系列
- Claude 3 (Haiku, Sonnet, Opus)
- Claude 3.5 (Haiku, Sonnet)
- Claude 3.7 Sonnet
- Claude 4 (Sonnet, Opus)

### Meta Llama系列
- Llama 3.0 (8B, 70B)
- Llama 3.1 (8B, 70B, 405B)
- Llama 3.2 (1B, 3B, 11B, 90B)
- Llama 3.3 (70B)
- Llama 4 (Maverick, Scout)

### Amazon Nova系列
- Nova Micro
- Nova Lite
- Nova Pro
- Nova Premier

### 其他模型
- Cohere Command系列
- Mistral系列
- AI21 Jamba系列
- Amazon Titan系列
- DeepSeek R1

## 配置示例

### 渠道配置
```json
{
  "region": "us-east-1",
  "ak": "your-access-key",
  "sk": "your-secret-key"
}
```

### 模型重定向示例
```json
{
  "gpt-4": "claude-sonnet-4-20250514",
  "gpt-3.5-turbo": "claude-3-5-haiku-20241022",
  "llama-70b": "llama3-1-70b-instruct",
  "nova": "nova-pro"
}
```

## API升级

### 从旧API到新API
- **旧**: InvokeModel / InvokeModelWithResponseStream
- **新**: Converse / ConverseStream

### 优势
1. **统一接口**: 所有模型使用相同的请求格式
2. **更简洁**: 不需要模型特定的payload转换
3. **功能丰富**: 支持更多高级功能
4. **官方推荐**: AWS官方推荐的新API

## 使用方式

### 基本聊天
```bash
curl -X POST http://localhost:3000/v1/chat/completions \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "claude-3-5-sonnet-20241022",
    "messages": [
      {"role": "user", "content": "Hello!"}
    ]
  }'
```

### 多模态输入
```bash
curl -X POST http://localhost:3000/v1/chat/completions \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "claude-3-5-sonnet-20241022",
    "messages": [
      {
        "role": "user",
        "content": [
          {"type": "text", "text": "What is in this image?"},
          {"type": "image_url", "image_url": {"url": "data:image/jpeg;base64,..."}}
        ]
      }
    ]
  }'
```

### 流式响应
```bash
curl -X POST http://localhost:3000/v1/chat/completions \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "claude-3-5-sonnet-20241022",
    "messages": [{"role": "user", "content": "Tell me a story"}],
    "stream": true
  }'
```

## 迁移指南

### 从旧适配器迁移
1. 旧的claude和llama3适配器已被移除
2. 所有模型现在使用统一的Converse适配器
3. 模型映射保持不变，无需修改配置
4. 新适配器自动支持更多模型

### 兼容性
- ✅ 完全向后兼容
- ✅ 保持现有的模型映射
- ✅ 支持现有的渠道配置
- ✅ 无需修改客户端代码

## 故障排除

### 常见问题
1. **模型不支持**: 检查模型名称是否在支持列表中
2. **认证失败**: 检查AWS凭证配置
3. **区域错误**: 确保模型在指定区域可用
4. **配额限制**: 检查AWS账户配额

### 调试
启用调试日志查看详细信息：
```bash
export LOG_LEVEL=debug
```

## 开发

### 添加新模型
1. 在`model.go`中的`ModelMapping`添加映射
2. 如果是多模态模型，在`IsMultimodalModel`中添加
3. 运行测试确保正常工作

### 测试
```bash
go test ./relay/adaptor/aws/converse/...
```
