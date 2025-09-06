@echo off
chcp 65001 >nul
echo ========================================
echo 羽毛球场地预订桌面应用启动脚本
echo ========================================
echo.

:: 检查Node.js是否安装
node --version >nul 2>&1
if %errorlevel% neq 0 (
    echo [错误] 未检测到Node.js，请先安装Node.js (版本16或更高)
    echo 下载地址: https://nodejs.org/
    pause
    exit /b 1
)

echo [信息] Node.js版本:
node --version
echo.

:: 检查是否已安装依赖
if not exist "node_modules" (
    echo [信息] 首次运行，正在安装依赖包...
    echo 这可能需要几分钟时间，请耐心等待...
    echo.
    npm install
    if %errorlevel% neq 0 (
        echo [错误] 依赖安装失败，请检查网络连接
        pause
        exit /b 1
    )
    echo.
    echo [成功] 依赖安装完成！
    echo.
)

:: 检查Go程序是否存在
set GO_EXE_PATH="..\src\workern\yumaoqiu\fetch_and_order.exe"
set GO_SRC_PATH="..\src\workern\yumaoqiu\fetch_and_order.go"

if exist %GO_EXE_PATH% (
    echo [信息] 找到Go可执行文件: %GO_EXE_PATH%
) else if exist %GO_SRC_PATH% (
    echo [信息] 找到Go源码文件: %GO_SRC_PATH%
    echo [提示] 正在编译Go程序以获得更好的性能...
    
    :: 检查Go是否安装
    go version >nul 2>&1
    if %errorlevel% equ 0 (
        cd "..\src\workern\yumaoqiu"
        go build -o fetch_and_order.exe fetch_and_order.go
        if %errorlevel% equ 0 (
            echo [成功] Go程序编译完成
        ) else (
            echo [警告] Go程序编译失败，将使用源码运行
        )
        cd "%~dp0"
    ) else (
        echo [警告] 未安装Go环境，将使用源码运行（需要Go环境）
    )
) else (
    echo [警告] 未找到Go程序文件，请确保以下文件存在：
    echo   - %GO_EXE_PATH%
    echo   - %GO_SRC_PATH%
)
echo.

:: 启动应用
echo [信息] 正在启动应用...
echo 应用启动后，此窗口将保持打开状态
echo 关闭此窗口将同时关闭应用
echo.
echo ========================================
echo.

npm start

echo.
echo [信息] 应用已关闭
pause