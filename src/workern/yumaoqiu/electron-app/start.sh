#!/bin/bash

echo "========================================"
echo "羽毛球场地预订桌面应用启动脚本"
echo "========================================"
echo

# 检查Node.js是否安装
if ! command -v node &> /dev/null; then
    echo "[错误] 未检测到Node.js，请先安装Node.js (版本16或更高)"
    echo "下载地址: https://nodejs.org/"
    exit 1
fi

echo "[信息] Node.js版本:"
node --version
echo

# 检查是否已安装依赖
if [ ! -d "node_modules" ]; then
    echo "[信息] 首次运行，正在安装依赖包..."
    echo "这可能需要几分钟时间，请耐心等待..."
    echo
    npm install
    if [ $? -ne 0 ]; then
        echo "[错误] 依赖安装失败，请检查网络连接"
        exit 1
    fi
    echo
    echo "[成功] 依赖安装完成！"
    echo
fi

# 检查Go程序是否存在
GO_EXE_PATH="../src/workern/yumaoqiu/fetch_and_order"
GO_SRC_PATH="../src/workern/yumaoqiu/fetch_and_order.go"

if [ -f "$GO_EXE_PATH" ]; then
    echo "[信息] 找到Go可执行文件: $GO_EXE_PATH"
elif [ -f "$GO_SRC_PATH" ]; then
    echo "[信息] 找到Go源码文件: $GO_SRC_PATH"
    echo "[提示] 正在编译Go程序以获得更好的性能..."
    
    # 检查Go是否安装
    if command -v go &> /dev/null; then
        cd "../src/workern/yumaoqiu"
        go build -o fetch_and_order fetch_and_order.go
        if [ $? -eq 0 ]; then
            echo "[成功] Go程序编译完成"
        else
            echo "[警告] Go程序编译失败，将使用源码运行"
        fi
        cd - > /dev/null
    else
        echo "[警告] 未安装Go环境，将使用源码运行（需要Go环境）"
    fi
else
    echo "[警告] 未找到Go程序文件，请确保以下文件存在："
    echo "  - $GO_EXE_PATH"
    echo "  - $GO_SRC_PATH"
fi
echo

# 启动应用
echo "[信息] 正在启动应用..."
echo "应用启动后，此终端将保持打开状态"
echo "按 Ctrl+C 可以关闭应用"
echo
echo "========================================"
echo

npm start

echo
echo "[信息] 应用已关闭"