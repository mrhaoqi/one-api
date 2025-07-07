# AWS Bedrock Claude 地理位置限制故障排除指南

## 🚨 错误信息
```
Access to Anthropic models is not allowed from unsupported countries, regions, or territories. 
Please refer to https://www.anthropic.com/supported-countries for more information.
```

## 🔍 问题分析

### 当前检测到的信息
- **IP地址**: 23.172.200.71
- **显示位置**: 美国加利福尼亚州圣何塞
- **国家**: 美国 (支持列表中)

### 可能的根本原因

1. **AWS区域配置问题** ⭐ 最可能
2. **模型访问权限问题**
3. **网络路由问题**
4. **AWS账户区域限制**

## 🛠️ 解决方案

### 方案1: 检查并更改AWS区域 (推荐)

#### 1.1 检查当前AWS配置
```bash
aws configure list
aws configure get region
```

#### 1.2 支持Claude的AWS区域
| 区域代码 | 区域名称 | Claude支持 | 推荐度 |
|---------|---------|-----------|--------|
| `us-east-1` | N. Virginia | ✅ 完整支持 | ⭐⭐⭐ |
| `us-west-2` | Oregon | ✅ 完整支持 | ⭐⭐⭐ |
| `eu-central-1` | Frankfurt | ✅ 部分支持 | ⭐⭐ |
| `ap-northeast-1` | Tokyo | ✅ 部分支持 | ⭐⭐ |
| `ap-southeast-1` | Singapore | ✅ 部分支持 | ⭐⭐ |

#### 1.3 更改AWS区域
```bash
# 方法1: 使用AWS CLI
aws configure set region us-east-1

# 方法2: 设置环境变量
export AWS_DEFAULT_REGION=us-east-1
export AWS_REGION=us-east-1

# 方法3: 在One-API中配置
# 在渠道配置中设置区域参数
```

### 方案2: 申请模型访问权限

#### 2.1 登录AWS控制台
1. 访问 [AWS Bedrock控制台](https://console.aws.amazon.com/bedrock/)
2. 选择正确的区域 (us-east-1推荐)

#### 2.2 申请模型访问
1. 点击左侧菜单 "Model access"
2. 找到 "Anthropic" 部分
3. 勾选需要的Claude模型：
   - Claude 3 Haiku
   - Claude 3 Sonnet  
   - Claude 3 Opus
   - Claude 3.5 Sonnet
   - Claude 3.5 Haiku
   - Claude 4 Sonnet
   - Claude 4 Opus
4. 点击 "Request model access"
5. 填写使用案例信息
6. 等待审批 (通常几分钟到几小时)

### 方案3: 检查网络配置

#### 3.1 确认真实IP位置
```bash
# 检查多个IP查询服务
curl ipinfo.io
curl ifconfig.me/all.json
curl ipapi.co/json
```

#### 3.2 如果使用VPN/代理
- 尝试断开VPN连接
- 使用美国或欧盟的VPN节点
- 确保VPN提供商的IP在支持列表中

### 方案4: AWS账户检查

#### 4.1 确认账户区域
1. 登录AWS控制台
2. 检查右上角的区域选择器
3. 确认账户是在支持的区域注册的

#### 4.2 检查账户状态
- 确认账户已完成验证
- 检查是否有任何限制或暂停

## 🧪 测试步骤

### 1. 基础连接测试
```bash
# 测试AWS连接
aws bedrock list-foundation-models --region us-east-1

# 测试特定模型
aws bedrock get-foundation-model \
  --model-identifier us.anthropic.claude-3-haiku-20240307-v1:0 \
  --region us-east-1
```

### 2. One-API配置测试
1. 在One-API管理界面中
2. 编辑AWS渠道配置
3. 确保区域设置为 `us-east-1`
4. 保存并重新测试

### 3. 模型访问测试
```bash
# 使用AWS CLI测试模型调用
aws bedrock-runtime invoke-model \
  --model-id us.anthropic.claude-3-haiku-20240307-v1:0 \
  --body '{"anthropic_version":"bedrock-2023-05-31","max_tokens":100,"messages":[{"role":"user","content":"Hello"}]}' \
  --cli-binary-format raw-in-base64-out \
  --region us-east-1 \
  output.json
```

## 📋 检查清单

- [ ] AWS区域设置为支持的区域 (推荐us-east-1)
- [ ] 已申请Claude模型访问权限
- [ ] 模型访问请求已获批准
- [ ] 网络连接正常，无VPN干扰
- [ ] AWS凭证配置正确
- [ ] One-API渠道配置正确
- [ ] 使用正确的推理配置文件ID

## 🆘 如果仍然无法解决

### 联系支持
1. **AWS支持**: 创建技术支持案例
2. **Anthropic支持**: 发送邮件到支持团队
3. **One-API社区**: 在GitHub提交Issue

### 提供的信息
- 错误的完整日志
- AWS区域配置
- 模型访问权限状态
- 网络环境信息
- One-API版本信息

## 💡 预防措施

1. **使用推荐区域**: 优先选择us-east-1
2. **提前申请权限**: 在使用前申请所有需要的模型
3. **监控配额**: 定期检查模型使用配额
4. **备用区域**: 配置多个区域的渠道作为备份
