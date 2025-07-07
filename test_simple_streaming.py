#!/usr/bin/env python3
"""
测试简单的流式请求（不包含工具）
"""

import json
import requests

# 配置
BASE_URL = "http://localhost:3000"
API_KEY = "sk-QgFaYx2G1HSoyrJP3426C5Ff16E64421A4Ca9c38E37e117e"
MODEL = "claude-3-5-sonnet-latest"

def test_simple_streaming():
    """测试简单的流式请求（不包含工具）"""
    print("🌊 测试简单的流式请求...")
    
    messages = [
        {
            "role": "user",
            "content": "请简单介绍一下人工智能"
        }
    ]
    
    payload = {
        "model": MODEL,
        "messages": messages,
        "stream": True,
        "max_tokens": 500
    }
    
    headers = {
        "Authorization": f"Bearer {API_KEY}",
        "Content-Type": "application/json"
    }
    
    try:
        response = requests.post(
            f"{BASE_URL}/v1/chat/completions",
            headers=headers,
            json=payload,
            stream=True,
            timeout=30
        )
        
        if response.status_code == 200:
            print("✅ 流式请求成功")
            
            # 解析流式响应
            content_received = False
            for line in response.iter_lines():
                if line:
                    line = line.decode('utf-8')
                    if line.startswith('data: '):
                        data = line[6:]
                        if data.strip() == '[DONE]':
                            break
                        try:
                            chunk = json.loads(data)
                            if chunk.get('choices') and chunk['choices'][0].get('delta', {}).get('content'):
                                content_received = True
                                content = chunk['choices'][0]['delta']['content']
                                print(f"📝 收到内容: {content}", end='', flush=True)
                        except json.JSONDecodeError:
                            continue
            
            print()  # 换行
            
            if content_received:
                print("✅ 简单流式请求测试成功")
                return True
            else:
                print("❌ 简单流式请求测试失败：未收到内容")
                return False
        else:
            print(f"❌ 流式请求失败: {response.status_code} - {response.text}")
            return False
            
    except Exception as e:
        print(f"❌ 流式请求异常: {e}")
        return False

def test_simple_non_streaming():
    """测试简单的非流式请求（不包含工具）"""
    print("\n📝 测试简单的非流式请求...")
    
    messages = [
        {
            "role": "user",
            "content": "请用一句话介绍Python编程语言"
        }
    ]
    
    payload = {
        "model": MODEL,
        "messages": messages,
        "stream": False,
        "max_tokens": 200
    }
    
    headers = {
        "Authorization": f"Bearer {API_KEY}",
        "Content-Type": "application/json"
    }
    
    try:
        response = requests.post(
            f"{BASE_URL}/v1/chat/completions",
            headers=headers,
            json=payload,
            timeout=30
        )
        
        if response.status_code == 200:
            result = response.json()
            content = result.get('choices', [{}])[0].get('message', {}).get('content', '')
            if content and 'Python' in content:
                print(f"✅ 非流式请求成功: {content}")
                return True
            else:
                print(f"❌ 非流式请求失败：内容不符合预期: {content}")
                return False
        else:
            print(f"❌ 非流式请求失败: {response.status_code} - {response.text}")
            return False
            
    except Exception as e:
        print(f"❌ 非流式请求异常: {e}")
        return False

if __name__ == "__main__":
    print("🚀 开始测试AWS适配器基本功能...")
    
    # 测试简单的非流式请求
    success1 = test_simple_non_streaming()
    
    # 测试简单的流式请求
    success2 = test_simple_streaming()
    
    print(f"\n📊 测试结果:")
    print(f"  - 简单非流式请求: {'✅ 成功' if success1 else '❌ 失败'}")
    print(f"  - 简单流式请求: {'✅ 成功' if success2 else '❌ 失败'}")
    
    if success1 and success2:
        print("\n🎉 所有基本功能测试通过！")
    else:
        print("\n⚠️  部分测试失败，需要进一步检查")
