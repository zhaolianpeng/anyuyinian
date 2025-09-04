#!/bin/bash

# 云函数部署脚本
# 适用于腾讯云云函数

echo "☁️ 构建云函数部署包..."

# 1. 创建部署目录
mkdir -p cloud_function_package
cd cloud_function_package

# 2. 复制必要文件
echo "📦 复制文件..."
cp ../main.go .
cp ../go.mod .
cp ../go.sum .
cp -r ../config .
cp -r ../db .
cp -r ../service .
cp -r ../handler .
cp -r ../middleware .
cp -r ../utils .

# 3. 创建云函数入口文件
echo "🔧 创建云函数入口..."
cat > scf.go << 'EOF'
package main

import (
    "context"
    "github.com/tencentyun/scf-go-lib/cloudfunction"
    "net/http"
)

// 云函数入口
func main() {
    // 初始化应用
    initApp()
    
    // 启动HTTP服务器
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // 处理所有请求
        handleRequest(w, r)
    })
    
    // 启动云函数
    cloudfunction.Start(handleRequest)
}

// 处理请求
func handleRequest(w http.ResponseWriter, r *http.Request) {
    // 设置CORS头
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    
    // 处理OPTIONS请求
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }
    
    // 路由处理
    switch r.URL.Path {
    case "/api/health":
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("healthy"))
    case "/api/payment/notify":
        // 处理支付回调
        handlePaymentNotify(w, r)
    default:
        // 其他API处理
        handleAPI(w, r)
    }
}

// 处理支付回调
func handlePaymentNotify(w http.ResponseWriter, r *http.Request) {
    // 实现支付回调逻辑
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("success"))
}

// 处理其他API
func handleAPI(w http.ResponseWriter, r *http.Request) {
    // 实现其他API逻辑
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("API response"))
}
EOF

# 4. 创建部署配置文件
echo "⚙️ 创建部署配置..."
cat > template.yaml << 'EOF'
Resources:
  default:
    Type: TencentCloud::Scf::Function
    Properties:
      CodeUri: .
      Description: 安语一年小程序后端API
      Environment:
        Variables:
          WECHAT_PAY_APP_ID: wx101090677bd5219e
          WECHAT_PAY_MCH_ID: 1726638701
          WECHAT_PAY_MCH_KEY: JQzOCB8doIdgaUjAobELsk9nTyxdKhat
          WECHAT_PAY_NOTIFY_URL: https://your-function-url.com/api/payment/notify
          WECHAT_PAY_ENVIRONMENT: production
      FunctionName: anyuyinian-api
      Handler: scf.main
      MemorySize: 512
      Runtime: Go1
      Timeout: 30
      Triggers:
        - Type: APIGW
          Properties:
            Name: anyuyinian-api-gateway
            Description: 安语一年API网关
            Protocol: https
            ServiceName: anyuyinian-service
EOF

# 5. 创建部署脚本
echo "🚀 创建部署脚本..."
cat > deploy.sh << 'EOF'
#!/bin/bash

# 安装Serverless Framework
npm install -g serverless

# 配置腾讯云凭证
serverless config credentials --provider tencent --key YOUR_SECRET_ID --secret YOUR_SECRET_KEY

# 部署云函数
serverless deploy

echo "✅ 云函数部署完成！"
echo "🌐 API网关地址: https://your-function-url.com"
echo "📱 支付回调地址: https://your-function-url.com/api/payment/notify"
EOF

chmod +x deploy.sh

echo "✅ 云函数部署包构建完成！"
echo "📁 部署包位置: cloud_function_package/"
echo "📝 下一步："
echo "1. 修改 template.yaml 中的配置"
echo "2. 运行 ./deploy.sh 部署云函数"
echo "3. 获取云函数URL并配置到微信商户平台"
