# AWS Bedrock 推理配置文件修复

## 问题描述

在测试AWS Bedrock Claude模型时遇到以下错误：

```
ValidationException: Invocation of model ID anthropic.claude-3-5-sonnet-20241022-v2:0 with on-demand throughput isn't supported. Retry your request with the ID or ARN of an inference profile that contains this model.
```

## 根本原因

AWS Bedrock对新版本的Claude模型（特别是Claude 3.5+和Claude 4系列）有新的要求：

1. **旧模型**：可以直接使用模型ID（如 `anthropic.claude-3-haiku-20240307-v1:0`）
2. **新模型**：必须使用推理配置文件ID（如 `us.anthropic.claude-3-5-sonnet-20241022-v2:0`）

## 解决方案

### 修复内容

更新了 `relay/adaptor/aws/converse/model.go` 中的模型映射：

```go
// 旧模型（使用直接模型ID）
"claude-3-haiku-20240307":    "anthropic.claude-3-haiku-20240307-v1:0",
"claude-3-sonnet-20240229":   "anthropic.claude-3-sonnet-20240229-v1:0",
"claude-3-opus-20240229":     "anthropic.claude-3-opus-20240229-v1:0",
"claude-3-5-sonnet-20240620": "anthropic.claude-3-5-sonnet-20240620-v1:0",

// 新模型（使用推理配置文件ID - AWS要求）
"claude-3-5-sonnet-20241022": "us.anthropic.claude-3-5-sonnet-20241022-v2:0",
"claude-3-5-sonnet-latest":   "us.anthropic.claude-3-5-sonnet-20241022-v2:0",
"claude-3-5-haiku-20241022":  "us.anthropic.claude-3-5-haiku-20241022-v1:0",
"claude-3-7-sonnet-20250219": "us.anthropic.claude-3-7-sonnet-20250219-v1:0",
"claude-opus-4-20250514":     "us.anthropic.claude-opus-4-20250514-v1:0",
"claude-sonnet-4-20250514":   "us.anthropic.claude-sonnet-4-20250514-v1:0",
```

### 关键变化

1. **前缀变化**：从 `anthropic.` 改为 `us.anthropic.`
2. **推理配置文件**：新模型使用推理配置文件ID而不是直接模型ID
3. **向后兼容**：旧模型仍然使用原来的直接模型ID

## 推理配置文件 vs 直接模型ID

### 推理配置文件的优势

1. **跨区域支持**：推理配置文件可以在多个AWS区域使用
2. **自动路由**：AWS自动将请求路由到最佳可用区域
3. **更好的可用性**：提供更高的服务可用性
4. **成本优化**：AWS可以更好地优化资源使用

### 使用规则

| 模型系列 | 使用方式 | 示例 |
|---------|---------|------|
| Claude 1-2 | 直接模型ID | `anthropic.claude-v2` |
| Claude 3 (早期) | 直接模型ID | `anthropic.claude-3-haiku-20240307-v1:0` |
| Claude 3.5+ | 推理配置文件ID | `us.anthropic.claude-3-5-sonnet-20241022-v2:0` |
| Claude 4 | 推理配置文件ID | `us.anthropic.claude-sonnet-4-20250514-v1:0` |

## 测试验证

修复后，请重新测试您的AWS渠道：

1. **重新编译**：
   ```bash
   go build -o one-api
   ```

2. **重启服务**：
   ```bash
   ./one-api --port 3000
   ```

3. **测试渠道**：在管理界面中重新测试AWS渠道

## 支持的区域

推理配置文件支持以下AWS区域：
- `us-east-1` (北弗吉尼亚)
- `us-west-2` (俄勒冈)
- `eu-central-1` (法兰克福)
- `ap-southeast-2` (悉尼)
- `eu-west-3` (巴黎)

## 注意事项

1. **区域限制**：确保您的AWS区域支持所选的Claude模型
2. **权限要求**：确保您的AWS凭证有访问Bedrock推理配置文件的权限
3. **网络连接**：如果使用VPN，可能会影响地理位置检测

## 错误排查

如果仍然遇到问题，请检查：

1. **AWS区域**：确认使用支持的区域
2. **模型权限**：在AWS控制台中确认已启用相应模型
3. **网络环境**：检查是否有VPN或代理影响
4. **凭证配置**：确认AWS凭证配置正确

这个修复确保了与AWS Bedrock最新要求的兼容性，提供了更好的服务可用性和性能。
