#!/bin/bash

# 批量更新管理员页面使用标准API调用

echo "=== 更新管理员页面使用标准API调用 ==="

# 管理员页面目录
ADMIN_DIR="/Users/zhaolianpeng/code/Goproject/src/miniprogram/pages/admin"

# 需要修改的文件
FILES=(
    "login.js"
    "home.js" 
    "admins.js"
    "orders.js"
    "users.js"
)

# 为每个文件添加标准API导入
for file in "${FILES[@]}"; do
    file_path="$ADMIN_DIR/$file"
    if [ -f "$file_path" ]; then
        echo "处理文件: $file"
        
        # 检查是否已经导入了标准API
        if ! grep -q "require.*cloud-container-standard" "$file_path"; then
            # 在文件开头添加导入
            sed -i '' '1i\
const { api } = require("../../utils/cloud-container-standard")
' "$file_path"
            echo "  ✅ 已添加标准API导入"
        else
            echo "  ⚠️  已存在标准API导入"
        fi
    else
        echo "  ❌ 文件不存在: $file"
    fi
done

echo ""
echo "=== 手动修改说明 ==="
echo "请手动将以下调用方式："
echo "  app.callContainer('/api/admin/xxx', 'GET', data)"
echo "修改为："
echo "  api.adminXxx(data)"
echo ""
echo "具体映射："
echo "  /api/admin/login -> api.adminLogin"
echo "  /api/admin/check-status -> api.adminCheckStatus"
echo "  /api/admin/users -> api.adminUsers"
echo "  /api/admin/orders -> api.adminOrders"
echo "  /api/admin/stats -> api.adminStats"
echo "  /api/admin/admins -> api.adminAdmins"
echo "  /api/admin/set-admin -> api.adminSetAdmin"
echo "  /api/admin/remove-admin -> api.adminRemoveAdmin"
echo "  /api/admin/order/update-amount -> api.adminUpdateOrderAmount"
echo "  /api/admin/order/refund -> api.adminRefundOrder"
echo "  /api/consultation/stats -> api.consultationStats"

echo ""
echo "=== 更新完成 ==="
