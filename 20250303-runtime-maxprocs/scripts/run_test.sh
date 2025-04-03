#!/bin/bash

# 颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
TEST_DURATION=60 # 测试时长(秒)
CONNECTIONS=100   # 并发连接数
THREADS=4        # 线程数
ENDPOINT="/fibonacci?n=35" # 测试端点，可选 /fibonacci 或 /matrix
DEFAULT_PORT=8080
AUTO_PORT=8081

echo -e "${BLUE}开始GOMAXPROCS性能对比测试${NC}"
echo "=========================================="
echo "测试配置:"
echo "- 持续时间: ${TEST_DURATION}秒"
echo "- 并发连接: ${CONNECTIONS}"
echo "- 线程数: ${THREADS}"
echo "- 测试端点: ${ENDPOINT}"
echo "=========================================="

# 等待服务就绪
wait_for_service() {
  local port=$1
  local attempts=0
  local max_attempts=30
  
  echo -n "等待服务在端口 ${port} 就绪..."
  
  while ! curl -s "http://localhost:${port}/health" > /dev/null; do
    attempts=$((attempts + 1))
    if [ $attempts -ge $max_attempts ]; then
      echo -e "${YELLOW}超时!${NC}"
      return 1
    fi
    echo -n "."
    sleep 1
  done
  
  echo -e "${GREEN}服务就绪!${NC}"
  # 重置统计数据
  curl -s "http://localhost:${port}/reset" > /dev/null
  return 0
}

# 显示服务统计信息
show_stats() {
  local port=$1
  local description=$2
  
  echo -e "\n${BLUE}获取${description}统计信息:${NC}"
  curl -s "http://localhost:${port}/stats" | jq '.'
  
  # 从Docker获取容器CPU throttling信息
  if [ "$port" == "$DEFAULT_PORT" ]; then
    container="default-maxprocs"
  else
    container="auto-maxprocs"
  fi
  
  echo -e "\n${BLUE}Docker容器 CPU throttling 统计 (${container}):${NC}"
  docker stats --no-stream $container
  docker exec $container cat /sys/fs/cgroup/cpu/cpu.stat 2>/dev/null || docker exec $container cat /sys/fs/cgroup/cpu.stat
}

# 运行测试
run_test() {
  local port=$1
  local description=$2
  
  echo -e "\n${BLUE}开始测试 ${description} (端口 ${port})...${NC}"
  
  # 使用wrk进行负载测试
  wrk -t${THREADS} -c${CONNECTIONS} -d${TEST_DURATION}s "http://localhost:${port}${ENDPOINT}"
  
  # 显示服务统计信息
  show_stats $port "$description"
}

# 主测试逻辑
echo -e "\n${BLUE}第1步: 测试默认GOMAXPROCS配置 (宿主机CPU核数)${NC}"
if wait_for_service $DEFAULT_PORT; then
  run_test $DEFAULT_PORT "默认GOMAXPROCS"
else
  echo "跳过默认GOMAXPROCS测试，服务未就绪"
fi

echo -e "\n${BLUE}第2步: 测试自动调整GOMAXPROCS配置 (匹配容器CPU限制)${NC}"
if wait_for_service $AUTO_PORT; then
  run_test $AUTO_PORT "自动调整GOMAXPROCS"
else
  echo "跳过自动调整GOMAXPROCS测试，服务未就绪"
fi

echo -e "\n${GREEN}测试完成!${NC}"
echo "请对比上述测试结果分析性能差异"