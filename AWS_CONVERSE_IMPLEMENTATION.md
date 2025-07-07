# AWS Bedrock Converse API 实现

## 概述

本实现为 One-API 项目提供了基于 AWS Bedrock Converse API 的统一适配器，专门支持 Anthropic Claude 模型。

## 架构设计

### 目录结构
```
relay/adaptor/aws/
├── converse/          # 新的统一适配器
│   ├── adapter.go     # 适配器接口实现
│   ├── handler.go     # API处理器
│   └── model.go       # 模型映射
├── utils/             # 工具函数
├── adaptor.go         # 主适配器
└── registry.go        # 注册器
```

### 支持的模型

当前实现专门支持 Anthropic Claude 模型：
- Claude 3 系列
- Claude 3.5 系列  
- Claude 3.7 系列
- Claude 4 系列

## 功能特性

### ✅ 已实现
1. **非流式对话**: 完整的 Converse API 支持
2. **消息转换**: OpenAI 格式到 Bedrock Converse 格式的转换
3. **工具调用**: 支持 Function Calling
4. **图像输入**: 支持 base64 图像处理
5. **系统消息**: 正确处理系统提示
6. **参数映射**: temperature, top_p, max_tokens, stop_sequences

### 🚧 部分实现
1. **流式响应**: 框架已就绪，但需要更新的 AWS SDK 版本

### ❌ 未实现
1. **其他模型**: 目前只支持 Anthropic Claude

## 技术实现

### 请求转换
```go
// OpenAI 格式 -> Bedrock Converse 格式
func convertToClaudeConverseRequest(request *relaymodel.GeneralOpenAIRequest) ([]byte, error)
```

### 响应转换
```go
// Bedrock 响应 -> OpenAI 格式
func convertFromBedrockResponse(responseBody []byte, modelName, modelID string) (*openai.TextResponse, *relaymodel.Usage, error)
```

### 工具调用支持
```go
// OpenAI 工具 -> Converse 工具格式
func convertToolsToConverseFormat(tools []relaymodel.Tool, toolChoice interface{}) (map[string]interface{}, error)
```

## 使用方法

### 1. 配置 AWS 凭证
确保设置了正确的 AWS 凭证和区域。

### 2. 创建渠道
在 One-API 管理界面中创建 AWS 渠道，选择 Anthropic Claude 模型。

### 3. 模型映射
系统会自动将友好名称映射到正确的 AWS 模型 ID：

**Claude 3 系列**（支持直接模型ID）：
- `claude-3-haiku-20240307` -> `anthropic.claude-3-haiku-20240307-v1:0`
- `claude-3-sonnet-20240229` -> `anthropic.claude-3-sonnet-20240229-v1:0`
- `claude-3-opus-20240229` -> `anthropic.claude-3-opus-20240229-v1:0`

**Claude 4 系列**（强制要求推理配置文件ID）：
- `claude-sonnet-4-20250514` -> `us.anthropic.claude-sonnet-4-20250514-v1:0`
- `claude-opus-4-20250514` -> `us.anthropic.claude-opus-4-20250514-v1:0`

**重要**: Claude 4 系列模型采用新的部署架构，**必须**使用推理配置文件 ID，不支持直接模型 ID。

## 错误处理

### 推理配置文件错误（已修复）
```
Invocation of model ID anthropic.claude-sonnet-4-20250514-v1:0 with on-demand throughput isn't supported.
Retry your request with the ID or ARN of an inference profile that contains this model.
```
**解决方案**: 使用推理配置文件 ID `us.anthropic.claude-sonnet-4-20250514-v1:0`

### 地理位置限制错误
```
Access to Anthropic models is not allowed from unsupported countries, regions, or territories.
Please refer to https://www.anthropic.com/supported-countries for more information.
```
**可能原因和解决方案**:
1. **AWS区域问题**: 确保使用支持Claude的AWS区域
   - 推荐: `us-east-1` (N. Virginia)
   - 其他: `us-west-2`, `eu-central-1`, `ap-northeast-1`, `ap-southeast-1`
2. **网络问题**: 检查是否使用VPN或代理，可能影响地理位置检测
3. **账户区域**: 确认AWS账户注册的区域在支持列表中
4. **模型访问权限**: 在AWS控制台申请Claude模型的访问权限

### 不支持的模型
```
unsupported model: meta.llama3-8b-instruct-v1:0. Only Anthropic Claude models are supported with Converse API
```

### 流式响应
```
streaming not yet implemented for this AWS SDK version
```

## 推理配置文件说明

AWS Bedrock 的新版本 Claude 模型需要使用推理配置文件（Inference Profile）而不是直接的模型 ID。

### 推理配置文件的优势
1. **跨区域推理**: 自动路由到可用的区域
2. **负载均衡**: 在多个区域间分配请求
3. **高可用性**: 提供更好的服务可用性

### 支持的推理配置文件
| 模型名称 | 推理配置文件 ID |
|---------|----------------|
| Claude 3 Haiku | `us.anthropic.claude-3-haiku-20240307-v1:0` |
| Claude 3 Sonnet | `us.anthropic.claude-3-sonnet-20240229-v1:0` |
| Claude 3 Opus | `us.anthropic.claude-3-opus-20240229-v1:0` |
| Claude 3.5 Sonnet v1 | `us.anthropic.claude-3-5-sonnet-20240620-v1:0` |
| Claude 3.5 Sonnet v2 | `us.anthropic.claude-3-5-sonnet-20241022-v2:0` |
| Claude 3.5 Haiku | `us.anthropic.claude-3-5-haiku-20241022-v1:0` |
| Claude 3.7 Sonnet | `us.anthropic.claude-3-7-sonnet-20250219-v1:0` |
| Claude Opus 4 | `us.anthropic.claude-opus-4-20250514-v1:0` |
| Claude Sonnet 4 | `us.anthropic.claude-sonnet-4-20250514-v1:0` |

## 升级路径

### 短期 (当前版本)
- ✅ 基本的 Claude 对话功能
- ✅ 工具调用支持
- ✅ 图像输入支持

### 中期 (下个版本)
- 🔄 完整的流式响应支持
- 🔄 更新 AWS SDK 到支持完整 Converse API 的版本

### 长期 (未来版本)
- 📋 支持更多模型 (Meta Llama, Amazon Nova 等)
- 📋 高级功能 (文档聊天, 引用等)

## 注意事项

1. **推理配置文件**: 新版本的 Claude 模型（Claude 3+）需要使用推理配置文件 ID 而不是直接的模型 ID
2. **AWS SDK 版本**: 当前使用的 AWS SDK 版本可能不支持最新的 Converse API 特性
3. **模型限制**: 只支持 Anthropic Claude 模型
4. **流式响应**: 需要更新的 SDK 版本才能完全支持
5. **错误处理**: 提供了清晰的错误信息指导用户

## 贡献指南

如果您想贡献代码：

1. **完善流式响应**: 更新 AWS SDK 并实现完整的流式处理
2. **添加模型支持**: 扩展支持其他 Bedrock 模型
3. **优化性能**: 改进请求/响应转换的性能
4. **增强测试**: 添加更多的单元测试和集成测试

## 相关文档

- [AWS Bedrock Converse API 文档](https://docs.aws.amazon.com/bedrock/latest/userguide/conversation-inference.html)
- [Anthropic Claude 模型文档](https://docs.aws.amazon.com/bedrock/latest/userguide/model-parameters-claude.html)
- [支持的模型和功能](https://docs.aws.amazon.com/bedrock/latest/userguide/conversation-inference-supported-models-features.html)
