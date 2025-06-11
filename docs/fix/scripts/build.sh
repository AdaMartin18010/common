#!/bin/bash

# Golang Common 库构建脚本
# 作者: AI Assistant
# 版本: 1.0.0

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 配置变量
PROJECT_NAME="golang-common"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_DIR="build"
DIST_DIR="dist"
DOCKER_IMAGE="golang-common"
DOCKER_TAG="${VERSION}"

# 显示构建信息
show_build_info() {
    log_info "构建信息:"
    echo "  项目名称: ${PROJECT_NAME}"
    echo "  版本: ${VERSION}"
    echo "  构建时间: ${BUILD_TIME}"
    echo "  Git提交: ${GIT_COMMIT}"
    echo "  构建目录: ${BUILD_DIR}"
    echo "  分发目录: ${DIST_DIR}"
}

# 清理构建目录
clean_build() {
    log_info "清理构建目录..."
    rm -rf ${BUILD_DIR}
    rm -rf ${DIST_DIR}
    mkdir -p ${BUILD_DIR}
    mkdir -p ${DIST_DIR}
    log_success "构建目录清理完成"
}

# 检查依赖
check_dependencies() {
    log_info "检查依赖..."
    
    # 检查Go版本
    if ! command -v go &> /dev/null; then
        log_error "Go未安装，请先安装Go"
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}')
    log_info "Go版本: ${GO_VERSION}"
    
    # 检查Docker
    if command -v docker &> /dev/null; then
        DOCKER_VERSION=$(docker --version)
        log_info "Docker版本: ${DOCKER_VERSION}"
    else
        log_warning "Docker未安装，跳过Docker构建"
    fi
    
    # 检查代码格式化工具
    if ! command -v gofmt &> /dev/null; then
        log_warning "gofmt未安装"
    fi
    
    # 检查代码检查工具
    if ! command -v golint &> /dev/null; then
        log_warning "golint未安装，运行: go install golang.org/x/lint/golint@latest"
    fi
    
    # 检查测试覆盖率工具
    if ! command -v go test &> /dev/null; then
        log_error "go test不可用"
        exit 1
    fi
    
    log_success "依赖检查完成"
}

# 代码格式化
format_code() {
    log_info "格式化代码..."
    
    # 使用gofmt格式化
    if command -v gofmt &> /dev/null; then
        gofmt -s -w .
        log_success "代码格式化完成"
    else
        log_warning "跳过代码格式化，gofmt未安装"
    fi
}

# 代码检查
lint_code() {
    log_info "检查代码质量..."
    
    # 使用golint检查
    if command -v golint &> /dev/null; then
        golint ./... || true
        log_success "代码检查完成"
    else
        log_warning "跳过代码检查，golint未安装"
    fi
    
    # 使用go vet检查
    go vet ./...
    log_success "go vet检查完成"
}

# 运行测试
run_tests() {
    log_info "运行测试..."
    
    # 运行单元测试
    go test -v ./...
    
    # 运行基准测试
    go test -bench=. ./...
    
    # 生成测试覆盖率报告
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    
    log_success "测试完成，覆盖率报告: coverage.html"
}

# 构建二进制文件
build_binary() {
    log_info "构建二进制文件..."
    
    # 设置构建标志
    LDFLAGS="-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}"
    
    # 构建主程序
    go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/${PROJECT_NAME} .
    
    # 构建测试程序
    go test -c -o ${BUILD_DIR}/${PROJECT_NAME}-test .
    
    log_success "二进制文件构建完成"
}

# 构建Docker镜像
build_docker() {
    if ! command -v docker &> /dev/null; then
        log_warning "Docker未安装，跳过Docker构建"
        return
    fi
    
    log_info "构建Docker镜像..."
    
    # 创建Dockerfile
    cat > Dockerfile << EOF
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${PROJECT_NAME} .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/${PROJECT_NAME} .
COPY --from=builder /app/config ./config
EXPOSE 8080
CMD ["./${PROJECT_NAME}"]
EOF
    
    # 构建镜像
    docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} .
    docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:latest
    
    log_success "Docker镜像构建完成: ${DOCKER_IMAGE}:${DOCKER_TAG}"
}

# 创建发布包
create_release() {
    log_info "创建发布包..."
    
    # 创建发布目录
    RELEASE_DIR="${DIST_DIR}/${PROJECT_NAME}-${VERSION}"
    mkdir -p ${RELEASE_DIR}
    
    # 复制文件
    cp ${BUILD_DIR}/${PROJECT_NAME} ${RELEASE_DIR}/
    cp -r config ${RELEASE_DIR}/
    cp -r docs ${RELEASE_DIR}/
    cp README.md ${RELEASE_DIR}/
    cp LICENSE ${RELEASE_DIR}/
    cp go.mod ${RELEASE_DIR}/
    cp go.sum ${RELEASE_DIR}/
    
    # 创建启动脚本
    cat > ${RELEASE_DIR}/start.sh << 'EOF'
#!/bin/bash
./golang-common -config ./config/config.yaml
EOF
    chmod +x ${RELEASE_DIR}/start.sh
    
    # 创建停止脚本
    cat > ${RELEASE_DIR}/stop.sh << 'EOF'
#!/bin/bash
pkill -f golang-common
EOF
    chmod +x ${RELEASE_DIR}/stop.sh
    
    # 创建压缩包
    cd ${DIST_DIR}
    tar -czf ${PROJECT_NAME}-${VERSION}.tar.gz ${PROJECT_NAME}-${VERSION}
    zip -r ${PROJECT_NAME}-${VERSION}.zip ${PROJECT_NAME}-${VERSION}
    
    log_success "发布包创建完成: ${DIST_DIR}/${PROJECT_NAME}-${VERSION}.tar.gz"
}

# 安全扫描
security_scan() {
    log_info "执行安全扫描..."
    
    # 检查Go模块漏洞
    if command -v gosec &> /dev/null; then
        gosec ./...
        log_success "安全扫描完成"
    else
        log_warning "gosec未安装，跳过安全扫描"
        log_info "安装gosec: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"
    fi
}

# 性能分析
performance_analysis() {
    log_info "执行性能分析..."
    
    # 运行性能基准测试
    go test -bench=. -benchmem ./... > ${BUILD_DIR}/benchmark.txt
    
    # 生成CPU分析
    go test -cpuprofile=${BUILD_DIR}/cpu.prof -bench=. ./...
    
    # 生成内存分析
    go test -memprofile=${BUILD_DIR}/mem.prof -bench=. ./...
    
    log_success "性能分析完成"
}

# 显示帮助信息
show_help() {
    echo "Golang Common 库构建脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help          显示帮助信息"
    echo "  -c, --clean         清理构建目录"
    echo "  -f, --format        格式化代码"
    echo "  -l, --lint          检查代码质量"
    echo "  -t, --test          运行测试"
    echo "  -b, --build         构建二进制文件"
    echo "  -d, --docker        构建Docker镜像"
    echo "  -r, --release       创建发布包"
    echo "  -s, --security      执行安全扫描"
    echo "  -p, --performance   执行性能分析"
    echo "  -a, --all           执行所有步骤"
    echo ""
    echo "示例:"
    echo "  $0 --all            # 执行完整构建流程"
    echo "  $0 --test --build   # 运行测试并构建"
    echo "  $0 --docker         # 仅构建Docker镜像"
}

# 主函数
main() {
    # 解析命令行参数
    CLEAN=false
    FORMAT=false
    LINT=false
    TEST=false
    BUILD=false
    DOCKER=false
    RELEASE=false
    SECURITY=false
    PERFORMANCE=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -c|--clean)
                CLEAN=true
                shift
                ;;
            -f|--format)
                FORMAT=true
                shift
                ;;
            -l|--lint)
                LINT=true
                shift
                ;;
            -t|--test)
                TEST=true
                shift
                ;;
            -b|--build)
                BUILD=true
                shift
                ;;
            -d|--docker)
                DOCKER=true
                shift
                ;;
            -r|--release)
                RELEASE=true
                shift
                ;;
            -s|--security)
                SECURITY=true
                shift
                ;;
            -p|--performance)
                PERFORMANCE=true
                shift
                ;;
            -a|--all)
                CLEAN=true
                FORMAT=true
                LINT=true
                TEST=true
                BUILD=true
                DOCKER=true
                RELEASE=true
                SECURITY=true
                PERFORMANCE=true
                shift
                ;;
            *)
                log_error "未知选项: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # 如果没有指定任何选项，显示帮助
    if [[ "$CLEAN" == false && "$FORMAT" == false && "$LINT" == false && "$TEST" == false && "$BUILD" == false && "$DOCKER" == false && "$RELEASE" == false && "$SECURITY" == false && "$PERFORMANCE" == false ]]; then
        show_help
        exit 0
    fi
    
    # 显示构建信息
    show_build_info
    
    # 检查依赖
    check_dependencies
    
    # 执行构建步骤
    if [[ "$CLEAN" == true ]]; then
        clean_build
    fi
    
    if [[ "$FORMAT" == true ]]; then
        format_code
    fi
    
    if [[ "$LINT" == true ]]; then
        lint_code
    fi
    
    if [[ "$TEST" == true ]]; then
        run_tests
    fi
    
    if [[ "$SECURITY" == true ]]; then
        security_scan
    fi
    
    if [[ "$PERFORMANCE" == true ]]; then
        performance_analysis
    fi
    
    if [[ "$BUILD" == true ]]; then
        build_binary
    fi
    
    if [[ "$DOCKER" == true ]]; then
        build_docker
    fi
    
    if [[ "$RELEASE" == true ]]; then
        create_release
    fi
    
    log_success "构建流程完成！"
    
    # 显示结果
    echo ""
    log_info "构建结果:"
    if [[ "$BUILD" == true ]]; then
        echo "  二进制文件: ${BUILD_DIR}/${PROJECT_NAME}"
    fi
    if [[ "$DOCKER" == true ]]; then
        echo "  Docker镜像: ${DOCKER_IMAGE}:${DOCKER_TAG}"
    fi
    if [[ "$RELEASE" == true ]]; then
        echo "  发布包: ${DIST_DIR}/${PROJECT_NAME}-${VERSION}.tar.gz"
    fi
    if [[ "$TEST" == true ]]; then
        echo "  测试覆盖率: coverage.html"
    fi
    if [[ "$PERFORMANCE" == true ]]; then
        echo "  性能分析: ${BUILD_DIR}/benchmark.txt"
    fi
}

# 执行主函数
main "$@" 