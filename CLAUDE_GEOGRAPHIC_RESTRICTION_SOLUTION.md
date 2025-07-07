# Claude 地理位置限制问题解决方案

## 🚨 问题现状

**错误信息**:
```
Access to Anthropic models is not allowed from unsupported countries, regions, or territories.
```

**当前配置**:
- IP位置: 美国加利福尼亚州 ✅ (支持的国家)
- 模型ID: `anthropic.claude-sonnet-4-20250514-v1:0` ✅ (正确的直接模型ID)
- AWS区域: us-east-2 ✅ (官方文档确认支持)

## 🔍 根本原因分析

既然地理位置、模型ID和AWS区域都正确，问题很可能是：

### 1. 模型访问权限未申请 ⭐ 最可能
AWS Bedrock需要为每个模型单独申请访问权限

### 2. AWS凭证配置问题
- 凭证可能指向错误的区域
- 凭证权限不足

### 3. 网络路由问题
- VPN/代理影响地理位置检测
- DNS解析问题

## 🛠️ 解决步骤

### 步骤1: 申请Claude模型访问权限 (必须)

#### 1.1 登录AWS控制台
```bash
# 确保使用正确的区域
https://us-east-2.console.aws.amazon.com/bedrock/
```

#### 1.2 申请模型访问
1. 在左侧菜单点击 **"Model access"**
2. 点击 **"Request model access"** 或 **"Manage model access"**
3. 找到 **Anthropic** 部分
4. 勾选以下模型：
   - ✅ Claude 3 Haiku
   - ✅ Claude 3 Sonnet
   - ✅ Claude 3 Opus
   - ✅ Claude 3.5 Sonnet
   - ✅ Claude 3.5 Haiku
   - ✅ Claude 3.7 Sonnet
   - ✅ **Claude Sonnet 4** (重点)
   - ✅ **Claude Opus 4**
5. 填写使用案例信息
6. 提交申请

#### 1.3 等待审批
- 通常几分钟到几小时
- 检查邮件通知
- 在控制台查看状态

### 步骤2: 验证AWS配置

#### 2.1 检查当前配置
```bash
aws configure list
aws configure get region
aws sts get-caller-identity
```

#### 2.2 确保区域正确
```bash
# 设置为us-east-2
aws configure set region us-east-2

# 或使用环境变量
export AWS_DEFAULT_REGION=us-east-2
export AWS_REGION=us-east-2
```

#### 2.3 测试AWS连接
```bash
# 测试基本连接
aws bedrock list-foundation-models --region us-east-2

# 测试Claude Sonnet 4模型
aws bedrock get-foundation-model \
  --model-identifier anthropic.claude-sonnet-4-20250514-v1:0 \
  --region us-east-2
```

### 步骤3: 测试模型调用

#### 3.1 使用AWS CLI测试
```bash
# 创建测试请求
cat > test-request.json << 'EOF'
{
  "anthropic_version": "bedrock-2023-05-31",
  "max_tokens": 100,
  "messages": [
    {
      "role": "user",
      "content": "Hello, what is your model name?"
    }
  ]
}
EOF

# 调用模型
aws bedrock-runtime invoke-model \
  --model-id anthropic.claude-sonnet-4-20250514-v1:0 \
  --body file://test-request.json \
  --cli-binary-format raw-in-base64-out \
  --region us-east-2 \
  test-response.json

# 查看响应
cat test-response.json
```

#### 3.2 检查响应
如果成功，应该看到类似：
```json
{
  "content": [
    {
      "text": "I am Claude Sonnet 4...",
      "type": "text"
    }
  ],
  "id": "msg_...",
  "model": "claude-sonnet-4-20250514",
  "role": "assistant",
  "stop_reason": "end_turn",
  "stop_sequence": null,
  "type": "message",
  "usage": {
    "input_tokens": 12,
    "output_tokens": 25
  }
}
```

### 步骤4: 更新One-API配置

#### 4.1 确认渠道配置
在One-API管理界面：
1. 编辑AWS渠道
2. 确保区域设置为 `us-east-2`
3. 确认AWS凭证正确
4. 保存配置

#### 4.2 重新测试
重新测试Claude Sonnet 4模型

## 🔧 高级故障排除

### 如果仍然失败

#### 1. 检查IAM权限
确保AWS用户/角色有以下权限：
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "bedrock:InvokeModel",
        "bedrock:InvokeModelWithResponseStream",
        "bedrock:GetFoundationModel",
        "bedrock:ListFoundationModels"
      ],
      "Resource": [
        "arn:aws:bedrock:us-east-2::foundation-model/anthropic.claude-sonnet-4-20250514-v1:0",
        "arn:aws:bedrock:us-east-2::foundation-model/anthropic.*"
      ]
    }
  ]
}
```

#### 2. 尝试其他区域
如果us-east-2仍有问题，尝试：
- `us-east-1` (N. Virginia)
- `us-west-2` (Oregon)

#### 3. 使用推理配置文件
如果直接模型ID不工作，尝试推理配置文件：
```
us.anthropic.claude-sonnet-4-20250514-v1:0
```

#### 4. 联系AWS支持
如果所有步骤都完成但仍有问题：
1. 创建AWS技术支持案例
2. 提供详细的错误日志
3. 说明已完成的故障排除步骤

## 📋 检查清单

- [ ] 已申请Claude Sonnet 4模型访问权限
- [ ] 模型访问请求已获批准
- [ ] AWS区域设置为us-east-2
- [ ] AWS凭证配置正确
- [ ] IAM权限包含Bedrock访问
- [ ] 网络连接正常
- [ ] One-API渠道配置正确
- [ ] 使用正确的模型ID

## 💡 预防措施

1. **提前申请**: 在使用前申请所有需要的模型访问权限
2. **监控配额**: 定期检查模型使用配额和限制
3. **备用区域**: 配置多个区域的渠道作为备份
4. **权限管理**: 确保IAM权限足够但不过度

## 🎯 预期结果

完成以上步骤后，您应该能够：
1. 成功调用Claude Sonnet 4模型
2. 在One-API中正常使用AWS渠道
3. 不再看到地理位置限制错误
