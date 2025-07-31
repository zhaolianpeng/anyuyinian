#!/bin/bash

# 用户创建功能测试脚本

echo "=== 用户创建功能测试 ==="

# 设置测试参数
BASE_URL="http://localhost:80"
LOGIN_URL="$BASE_URL/api/wx/login"
USER_INFO_URL="$BASE_URL/api/user/info"

echo "测试URL: $LOGIN_URL"
echo ""

# 测试登录（应该创建新用户）
echo "测试登录（创建新用户）..."
LOGIN_RESPONSE=$(curl -s -X POST "$LOGIN_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "test_code_123",
    "userInfo": {
      "nickName": "测试用户",
      "avatarUrl": "https://example.com/avatar.jpg",
      "gender": 1,
      "country": "China",
      "province": "Guangdong",
      "city": "Shenzhen",
      "language": "zh_CN"
    }
  }')

echo "登录响应状态码: $?"
echo "登录响应内容: $LOGIN_RESPONSE"
echo ""

# 解析登录响应
if [ $? -eq 0 ]; then
  echo "✅ 登录请求发送成功"
  
  # 提取userId和token
  USER_ID=$(echo "$LOGIN_RESPONSE" | grep -o '"userId":[0-9]*' | cut -d':' -f2)
  TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
  
  if [ -n "$USER_ID" ]; then
    echo "✅ 获取到用户ID: $USER_ID"
    echo "✅ 获取到Token: $TOKEN"
    
    # 测试获取用户信息
    echo ""
    echo "测试获取用户信息 (userId=$USER_ID)..."
    USER_INFO_RESPONSE=$(curl -s -X GET "$USER_INFO_URL?userId=$USER_ID")
    
    echo "用户信息响应状态码: $?"
    echo "用户信息响应内容: $USER_INFO_RESPONSE"
    echo ""
    
    if [ $? -eq 0 ]; then
      echo "✅ 用户信息请求发送成功"
      
      # 检查用户信息是否包含正确的数据
      if echo "$USER_INFO_RESPONSE" | grep -q '"code":0'; then
        echo "✅ 获取用户信息成功"
        
        # 检查是否包含用户信息字段
        if echo "$USER_INFO_RESPONSE" | grep -q '"nickName":"测试用户"'; then
          echo "✅ 用户昵称正确"
        else
          echo "❌ 用户昵称不正确"
        fi
        
        if echo "$USER_INFO_RESPONSE" | grep -q '"avatarUrl":"https://example.com/avatar.jpg"'; then
          echo "✅ 用户头像正确"
        else
          echo "❌ 用户头像不正确"
        fi
        
        if echo "$USER_INFO_RESPONSE" | grep -q '"gender":1'; then
          echo "✅ 用户性别正确"
        else
          echo "❌ 用户性别不正确"
        fi
        
        if echo "$USER_INFO_RESPONSE" | grep -q '"country":"China"'; then
          echo "✅ 用户国家正确"
        else
          echo "❌ 用户国家不正确"
        fi
        
        if echo "$USER_INFO_RESPONSE" | grep -q '"province":"Guangdong"'; then
          echo "✅ 用户省份正确"
        else
          echo "❌ 用户省份不正确"
        fi
        
        if echo "$USER_INFO_RESPONSE" | grep -q '"city":"Shenzhen"'; then
          echo "✅ 用户城市正确"
        else
          echo "❌ 用户城市不正确"
        fi
        
        if echo "$USER_INFO_RESPONSE" | grep -q '"language":"zh_CN"'; then
          echo "✅ 用户语言正确"
        else
          echo "❌ 用户语言不正确"
        fi
        
      else
        echo "❌ 获取用户信息失败"
        echo "$USER_INFO_RESPONSE" | jq '.' 2>/dev/null || echo "$USER_INFO_RESPONSE"
      fi
    else
      echo "❌ 用户信息请求发送失败"
    fi
    
  else
    echo "❌ 未能获取到用户ID"
    echo "$LOGIN_RESPONSE" | jq '.' 2>/dev/null || echo "$LOGIN_RESPONSE"
  fi
else
  echo "❌ 登录请求发送失败"
fi

echo ""
echo "=== 测试完成 ==="
echo ""
echo "如果测试成功，说明："
echo "1. ✅ 登录接口正常工作"
echo "2. ✅ 用户信息被正确保存到数据库"
echo "3. ✅ 用户信息接口能正确查询数据库"
echo "4. ✅ 首次登录用户创建功能正常"
echo ""
echo "如果仍有问题，请检查："
echo "1. 后端服务是否正常运行"
echo "2. 数据库连接是否正常"
echo "3. 数据库表结构是否正确" 