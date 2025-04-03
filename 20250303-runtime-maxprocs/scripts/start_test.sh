#!/bin/bash

# 颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT=$(cd "$(dirname "$0")/.." && pwd)

echo -e "${BLUE}开始GOMAXPROCS性能测试流程${NC}"

# 1. 设置脚本为可执行
echo "设置测试脚本为可执行..."
chmod +x ${PROJECT_ROOT}/scripts/run_test.sh

# 2. 构建和启动Docker容器
echo -e "\n${BLUE}构建和启动Docker容器...${NC}"
cd ${PROJECT_ROOT}
docker-compose -f ${PROJECT_ROOT}/deployments/docker-compose.yml down
docker-compose -f ${PROJECT_ROOT}/deployments/docker-compose.yml up --build -d

# 3. 等待容器启动
echo -e "\n${BLUE}等待容器启动...${NC}"
sleep 5

# 4. 运行测试
echo -e "\n${BLUE}开始执行测试...${NC}"
${PROJECT_ROOT}/scripts/run_test.sh

# 5. 完成后保留容器运行以便检查
echo -e "\n${GREEN}测试完成!${NC}"
echo "容器仍在运行。完成后可以使用以下命令停止容器:"
echo "docker-compose -f ${PROJECT_ROOT}/deployments/docker-compose.yml down"