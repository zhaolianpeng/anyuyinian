#!/bin/bash

# 测试时间槽接口
echo "=== 测试时间槽接口 ==="

# 设置基础URL（请根据实际情况修改）
BASE_URL="http://localhost:80"

# 测试1: 获取明天的可用时间槽
echo "测试1: 获取明天的可用时间槽"
TOMORROW=$(date -v+1d +%Y-%m-%d 2>/dev/null || date -d "tomorrow" +%Y-%m-%d)
echo "日期: $TOMORROW"

curl -X POST "$BASE_URL/api/order/time_slots" \
  -H "Content-Type: application/json" \
  -d "{\"date\": \"$TOMORROW\"}" \
  | jq '.'

echo -e "\n"

# 测试2: 获取3天后的可用时间槽
echo "测试2: 获取3天后的可用时间槽"
DAY_AFTER_3=$(date -v+3d +%Y-%m-%d 2>/dev/null || date -d "3 days" +%Y-%m-%d)
echo "日期: $DAY_AFTER_3"

curl -X POST "$BASE_URL/api/order/time_slots" \
  -H "Content-Type: application/json" \
  -d "{\"date\": \"$DAY_AFTER_3\"}" \
  | jq '.'

echo -e "\n"

# 测试3: 测试无效日期（今天）
echo "测试3: 测试无效日期（今天）"
TODAY=$(date +%Y-%m-%d)
echo "日期: $TODAY"

curl -X POST "$BASE_URL/api/order/time_slots" \
  -H "Content-Type: application/json" \
  -d "{\"date\": \"$TODAY\"}" \
  | jq '.'

echo -e "\n"

# 测试4: 测试无效日期（8天后）
echo "测试4: 测试无效日期（8天后）"
DAY_AFTER_8=$(date -v+8d +%Y-%m-%d 2>/dev/null || date -d "8 days" +%Y-%m-%d)
echo "日期: $DAY_AFTER_8"

curl -X POST "$BASE_URL/api/order/time_slots" \
  -H "Content-Type: application/json" \
  -d "{\"date\": \"$DAY_AFTER_8\"}" \
  | jq '.'

echo -e "\n"

# 测试5: 测试日期格式错误
echo "测试5: 测试日期格式错误"
curl -X POST "$BASE_URL/api/order/time_slots" \
  -H "Content-Type: application/json" \
  -d "{\"date\": \"2024-13-45\"}" \
  | jq '.'

echo -e "\n"

# 测试6: 测试空日期
echo "测试6: 测试空日期"
curl -X POST "$BASE_URL/api/order/time_slots" \
  -H "Content-Type: application/json" \
  -d "{\"date\": \"\"}" \
  | jq '.'

echo -e "\n"

echo "=== 测试完成 ===" 