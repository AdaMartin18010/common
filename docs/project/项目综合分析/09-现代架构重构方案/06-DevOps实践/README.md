# 6. DevOps实践

## 6.1 DevOps概述

### 6.1.1 DevOps理念

DevOps是一种将开发(Development)、运维(Operations)和质量保证(QA)整合的文化和实践。对于葛洲坝船闸导航系统，DevOps实践将提供：

- **持续集成/持续部署(CI/CD)**：自动化构建、测试和部署
- **基础设施即代码(IaC)**：版本化的基础设施管理
- **GitOps**：基于Git的声明式部署
- **自动化测试**：全面的测试覆盖
- **监控和反馈**：快速发现和解决问题

### 6.1.2 DevOps价值

```mermaid
graph LR
    A[开发] --> B[构建]
    B --> C[测试]
    C --> D[部署]
    D --> E[监控]
    E --> F[反馈]
    F --> A
```

## 6.2 CI/CD流水线

### 6.2.1 GitLab CI/CD配置

```yaml
# .gitlab-ci.yml
stages:
  - validate
  - test
  - build
  - security
  - deploy

variables:
  DOCKER_REGISTRY: "registry.example.com"
  IMAGE_TAG: $CI_COMMIT_SHA
  GO_VERSION: "1.24"

# 代码验证阶段
code-validation:
  stage: validate
  image: golang:1.24-alpine
  script:
    - go mod download
    - go mod verify
    - go vet ./...
    - golangci-lint run
  only:
    - merge_requests
    - main

# 单元测试阶段
unit-tests:
  stage: test
  image: golang:1.24-alpine
  services:
    - postgres:13
    - redis:6-alpine
  variables:
    POSTGRES_DB: navlock_test
    POSTGRES_USER: navlock
    POSTGRES_PASSWORD: test_password
    REDIS_URL: redis://redis:6379
  script:
    - go mod download
    - go test -v -race -coverprofile=coverage.out ./...
    - go tool cover -html=coverage.out -o coverage.html
  coverage: '/coverage: \d+\.\d+% of statements/'
  artifacts:
    reports:
      coverage_report:
        coverage_format: cobertura
        path: coverage.xml
    paths:
      - coverage.html
    expire_in: 1 week
  only:
    - merge_requests
    - main

# 集成测试阶段
integration-tests:
  stage: test
  image: golang:1.24-alpine
  services:
    - postgres:13
    - redis:6-alpine
    - kafka:2.8
  variables:
    POSTGRES_DB: navlock_integration_test
    POSTGRES_USER: navlock
    POSTGRES_PASSWORD: test_password
    REDIS_URL: redis://redis:6379
    KAFKA_BROKERS: kafka:9092
  script:
    - go mod download
    - go test -v -tags=integration ./tests/integration/...
  only:
    - main

# 构建阶段
build:
  stage: build
  image: docker:20.10.16
  services:
    - docker:20.10.16-dind
  variables:
    DOCKER_TLS_CERTDIR: "/certs"
  script:
    - docker build -t $DOCKER_REGISTRY/navlock-service:$IMAGE_TAG .
    - docker build -t $DOCKER_REGISTRY/navlock-service:latest .
    - docker push $DOCKER_REGISTRY/navlock-service:$IMAGE_TAG
    - docker push $DOCKER_REGISTRY/navlock-service:latest
  only:
    - main

# 安全扫描阶段
security-scan:
  stage: security
  image: aquasec/trivy:latest
  script:
    - trivy image --exit-code 0 --severity HIGH,CRITICAL $DOCKER_REGISTRY/navlock-service:$IMAGE_TAG
    - trivy fs --exit-code 0 --severity HIGH,CRITICAL .
  allow_failure: true
  only:
    - main

# 部署到开发环境
deploy-dev:
  stage: deploy
  image: alpine:latest
  before_script:
    - apk add --no-cache curl
  script:
    - curl -X POST $DEV_WEBHOOK_URL -H "Content-Type: application/json" -d "{\"image\": \"$DOCKER_REGISTRY/navlock-service:$IMAGE_TAG\"}"
  environment:
    name: development
    url: https://navlock-dev.example.com
  only:
    - main

# 部署到测试环境
deploy-test:
  stage: deploy
  image: alpine:latest
  before_script:
    - apk add --no-cache curl
  script:
    - curl -X POST $TEST_WEBHOOK_URL -H "Content-Type: application/json" -d "{\"image\": \"$DOCKER_REGISTRY/navlock-service:$IMAGE_TAG\"}"
  environment:
    name: testing
    url: https://navlock-test.example.com
  when: manual
  only:
    - main

# 部署到生产环境
deploy-prod:
  stage: deploy
  image: alpine:latest
  before_script:
    - apk add --no-cache curl
  script:
    - curl -X POST $PROD_WEBHOOK_URL -H "Content-Type: application/json" -d "{\"image\": \"$DOCKER_REGISTRY/navlock-service:$IMAGE_TAG\"}"
  environment:
    name: production
    url: https://navlock.example.com
  when: manual
  only:
    - main
```

### 6.2.2 GitHub Actions配置

```yaml
# .github/workflows/ci-cd.yml
name: CI/CD Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: [1.23, 1.24]
    
    steps:
    - name: Checkout code
      uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: ${{ matrix.go-version }}
    
    - name: Cache dependencies
      uses: actions/cache@v3
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-
    
    - name: Download dependencies
      run: go mod download
    
    - name: Run linter
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest
    
    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out

  build:
    needs: test
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    
    steps:
    - name: Checkout code
      uses: actions/checkout@v4
    
    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v3
    
    - name: Log in to Container Registry
      uses: docker/login-action@v3
      with:
        registry: ${{ env.REGISTRY }}
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}
    
    - name: Extract metadata
      id: meta
      uses: docker/metadata-action@v5
      with:
        images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
        tags: |
          type=ref,event=branch
          type=ref,event=pr
          type=sha,prefix={{branch}}-
          type=raw,value=latest,enable={{is_default_branch}}
    
    - name: Build and push Docker image
      uses: docker/build-push-action@v5
      with:
        context: .
        push: true
        tags: ${{ steps.meta.outputs.tags }}
        labels: ${{ steps.meta.outputs.labels }}
        cache-from: type=gha
        cache-to: type=gha,mode=max

  security:
    needs: build
    runs-on: ubuntu-latest
    
    steps:
    - name: Run Trivy vulnerability scanner
      uses: aquasecurity/trivy-action@master
      with:
        image-ref: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:latest
        format: 'sarif'
        output: 'trivy-results.sarif'
    
    - name: Upload Trivy scan results
      uses: github/codeql-action/upload-sarif@v2
      if: always()
      with:
        sarif_file: 'trivy-results.sarif'

  deploy-dev:
    needs: [build, security]
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/develop'
    environment: development
    
    steps:
    - name: Deploy to development
      run: |
        echo "Deploying to development environment"
        # 部署脚本

  deploy-prod:
    needs: [build, security]
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    environment: production
    
    steps:
    - name: Deploy to production
      run: |
        echo "Deploying to production environment"
        # 部署脚本
```

## 6.3 GitOps实践

### 6.3.1 ArgoCD配置

```yaml
# ArgoCD应用配置
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: navlock-system
  namespace: argocd
  finalizers:
    - resources-finalizer.argocd.argoproj.io
spec:
  project: default
  source:
    repoURL: https://github.com/example/navlock-k8s-manifests.git
    targetRevision: HEAD
    path: navlock-system
  destination:
    server: https://kubernetes.default.svc
    namespace: navlock-system
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
      - PrunePropagationPolicy=foreground
      - PruneLast=true
    retry:
      limit: 5
      backoff:
        duration: 5s
        factor: 2
        maxDuration: 3m
---
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: navlock-devices
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/example/navlock-k8s-manifests.git
    targetRevision: HEAD
    path: navlock-devices
  destination:
    server: https://kubernetes.default.svc
    namespace: navlock-devices
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
---
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: navlock-monitoring
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/example/navlock-k8s-manifests.git
    targetRevision: HEAD
    path: navlock-monitoring
  destination:
    server: https://kubernetes.default.svc
    namespace: navlock-monitoring
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
```

### 6.3.2 Helm Charts

```yaml
# Chart.yaml
apiVersion: v2
name: navlock-system
description: A Helm chart for Navlock System
type: application
version: 1.0.0
appVersion: "1.0.0"

# values.yaml
replicaCount: 3

image:
  repository: registry.example.com/navlock-service
  tag: "latest"
  pullPolicy: IfNotPresent

imagePullSecrets:
  - name: navlock-registry-secret

serviceAccount:
  create: true
  annotations: {}
  name: ""

podAnnotations: {}

podSecurityContext: {}

securityContext: {}

service:
  type: ClusterIP
  port: 8080

ingress:
  enabled: true
  className: "istio"
  annotations:
    kubernetes.io/ingress.class: istio
  hosts:
    - host: navlock-api.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: navlock-tls-secret
      hosts:
        - navlock-api.example.com

resources:
  limits:
    cpu: 1000m
    memory: 2Gi
  requests:
    cpu: 500m
    memory: 1Gi

autoscaling:
  enabled: true
  minReplicas: 3
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70
  targetMemoryUtilizationPercentage: 80

nodeSelector: {}

tolerations: []

affinity: {}

# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "navlock-system.fullname" . }}
  labels:
    {{- include "navlock-system.labels" . | nindent 4 }}
spec:
  {{- if not .Values.autoscaling.enabled }}
  replicas: {{ .Values.replicaCount }}
  {{- end }}
  selector:
    matchLabels:
      {{- include "navlock-system.selectorLabels" . | nindent 6 }}
  template:
    metadata:
      labels:
        {{- include "navlock-system.selectorLabels" . | nindent 8 }}
      {{- with .Values.podAnnotations }}
      annotations:
        {{- toYaml . | nindent 8 }}
      {{- end }}
    spec:
      {{- with .Values.imagePullSecrets }}
      imagePullSecrets:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      serviceAccountName: {{ include "navlock-system.serviceAccountName" . }}
      securityContext:
        {{- toYaml .Values.podSecurityContext | nindent 8 }}
      containers:
        - name: {{ .Chart.Name }}
          securityContext:
            {{- toYaml .Values.securityContext | nindent 12 }}
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag | default .Chart.AppVersion }}"
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          ports:
            - name: http
              containerPort: 8080
              protocol: TCP
          livenessProbe:
            httpGet:
              path: /health
              port: http
          readinessProbe:
            httpGet:
              path: /ready
              port: http
          resources:
            {{- toYaml .Values.resources | nindent 12 }}
      {{- with .Values.nodeSelector }}
      nodeSelector:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.affinity }}
      affinity:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.tolerations }}
      tolerations:
        {{- toYaml . | nindent 8 }}
      {{- end }}
```

## 6.4 自动化测试

### 6.4.1 单元测试

```go
// 单元测试示例
package service

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "go.uber.org/zap"
)

// Mock依赖
type MockDeviceRepository struct {
    mock.Mock
}

func (m *MockDeviceRepository) GetByID(ctx context.Context, deviceID string) (*Device, error) {
    args := m.Called(ctx, deviceID)
    return args.Get(0).(*Device), args.Error(1)
}

func (m *MockDeviceRepository) Save(ctx context.Context, device *Device) error {
    args := m.Called(ctx, device)
    return args.Error(0)
}

// 测试设备管理服务
func TestDeviceManagementService_RegisterDevice(t *testing.T) {
    // 准备测试数据
    device := &Device{
        ID:       "test-device-001",
        Type:     "radar",
        Status:   DeviceStatusActive,
        Metadata: map[string]string{"location": "lock-1"},
    }

    // 创建Mock
    mockRepo := new(MockDeviceRepository)
    mockRepo.On("Save", mock.Anything, device).Return(nil)

    // 创建服务实例
    service := &deviceManagementService{
        deviceRepo: mockRepo,
        logger:     zap.NewNop(),
    }

    // 执行测试
    err := service.RegisterDevice(context.Background(), device)

    // 验证结果
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}

func TestDeviceManagementService_GetDeviceStatus(t *testing.T) {
    // 准备测试数据
    deviceID := "test-device-001"
    expectedDevice := &Device{
        ID:     deviceID,
        Type:   "radar",
        Status: DeviceStatusActive,
    }

    // 创建Mock
    mockRepo := new(MockDeviceRepository)
    mockRepo.On("GetByID", mock.Anything, deviceID).Return(expectedDevice, nil)

    // 创建服务实例
    service := &deviceManagementService{
        deviceRepo: mockRepo,
        logger:     zap.NewNop(),
    }

    // 执行测试
    device, err := service.GetDeviceStatus(context.Background(), deviceID)

    // 验证结果
    assert.NoError(t, err)
    assert.Equal(t, expectedDevice, device)
    mockRepo.AssertExpectations(t)
}

// 基准测试
func BenchmarkDeviceManagementService_RegisterDevice(b *testing.B) {
    device := &Device{
        ID:       "bench-device",
        Type:     "radar",
        Status:   DeviceStatusActive,
        Metadata: map[string]string{"location": "lock-1"},
    }

    mockRepo := new(MockDeviceRepository)
    mockRepo.On("Save", mock.Anything, device).Return(nil)

    service := &deviceManagementService{
        deviceRepo: mockRepo,
        logger:     zap.NewNop(),
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        device.ID = fmt.Sprintf("bench-device-%d", i)
        service.RegisterDevice(context.Background(), device)
    }
}
```

### 6.4.2 集成测试

```go
// 集成测试示例
package integration

import (
    "context"
    "database/sql"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    _ "github.com/lib/pq"
)

type DeviceManagementIntegrationTestSuite struct {
    suite.Suite
    db     *sql.DB
    ctx    context.Context
    cancel context.CancelFunc
}

func (suite *DeviceManagementIntegrationTestSuite) SetupSuite() {
    // 连接测试数据库
    db, err := sql.Open("postgres", "postgres://navlock:test@localhost:5432/navlock_test?sslmode=disable")
    assert.NoError(suite.T(), err)
    suite.db = db

    // 创建测试表
    _, err = suite.db.Exec(`
        CREATE TABLE IF NOT EXISTS devices (
            id VARCHAR(50) PRIMARY KEY,
            type VARCHAR(50) NOT NULL,
            status VARCHAR(20) NOT NULL,
            metadata JSONB,
            created_at TIMESTAMPTZ DEFAULT NOW(),
            updated_at TIMESTAMPTZ DEFAULT NOW()
        )
    `)
    assert.NoError(suite.T(), err)
}

func (suite *DeviceManagementIntegrationTestSuite) TearDownSuite() {
    if suite.db != nil {
        suite.db.Close()
    }
}

func (suite *DeviceManagementIntegrationTestSuite) SetupTest() {
    suite.ctx, suite.cancel = context.WithTimeout(context.Background(), 30*time.Second)
    
    // 清理测试数据
    _, err := suite.db.Exec("DELETE FROM devices")
    assert.NoError(suite.T(), err)
}

func (suite *DeviceManagementIntegrationTestSuite) TearDownTest() {
    if suite.cancel != nil {
        suite.cancel()
    }
}

func (suite *DeviceManagementIntegrationTestSuite) TestDeviceRegistration() {
    // 准备测试数据
    device := &Device{
        ID:       "integration-test-device",
        Type:     "radar",
        Status:   DeviceStatusActive,
        Metadata: map[string]string{"location": "lock-1"},
    }

    // 创建服务实例
    repo := NewPostgreSQLDeviceRepository(suite.db)
    service := NewDeviceManagementService(repo)

    // 执行测试
    err := service.RegisterDevice(suite.ctx, device)
    assert.NoError(suite.T(), err)

    // 验证数据已保存
    savedDevice, err := service.GetDeviceStatus(suite.ctx, device.ID)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), device.ID, savedDevice.ID)
    assert.Equal(suite.T(), device.Type, savedDevice.Type)
    assert.Equal(suite.T(), device.Status, savedDevice.Status)
}

func (suite *DeviceManagementIntegrationTestSuite) TestDeviceStatusUpdate() {
    // 准备测试数据
    device := &Device{
        ID:       "status-test-device",
        Type:     "radar",
        Status:   DeviceStatusActive,
        Metadata: map[string]string{"location": "lock-1"},
    }

    repo := NewPostgreSQLDeviceRepository(suite.db)
    service := NewDeviceManagementService(repo)

    // 注册设备
    err := service.RegisterDevice(suite.ctx, device)
    assert.NoError(suite.T(), err)

    // 更新设备状态
    device.Status = DeviceStatusInactive
    err = service.UpdateDeviceStatus(suite.ctx, device.ID, device.Status)
    assert.NoError(suite.T(), err)

    // 验证状态已更新
    updatedDevice, err := service.GetDeviceStatus(suite.ctx, device.ID)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), DeviceStatusInactive, updatedDevice.Status)
}

func TestDeviceManagementIntegrationTestSuite(t *testing.T) {
    suite.Run(t, new(DeviceManagementIntegrationTestSuite))
}
```

### 6.4.3 端到端测试

```go
// 端到端测试示例
package e2e

import (
    "context"
    "net/http"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

type NavlockE2ETestSuite struct {
    suite.Suite
    baseURL string
    client  *http.Client
    ctx     context.Context
    cancel  context.CancelFunc
}

func (suite *NavlockE2ETestSuite) SetupSuite() {
    suite.baseURL = "http://localhost:8080"
    suite.client = &http.Client{
        Timeout: 30 * time.Second,
    }
}

func (suite *NavlockE2ETestSuite) SetupTest() {
    suite.ctx, suite.cancel = context.WithTimeout(context.Background(), 60*time.Second)
}

func (suite *NavlockE2ETestSuite) TearDownTest() {
    if suite.cancel != nil {
        suite.cancel()
    }
}

func (suite *NavlockE2ETestSuite) TestDeviceManagementFlow() {
    // 1. 注册设备
    deviceID := "e2e-test-device"
    registerReq := map[string]interface{}{
        "device_id": deviceID,
        "type":      "radar",
        "status":    "active",
        "metadata": map[string]string{
            "location": "lock-1",
        },
    }

    resp, err := suite.makeRequest("POST", "/api/v1/devices", registerReq)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

    // 2. 查询设备状态
    resp, err = suite.makeRequest("GET", "/api/v1/devices/"+deviceID+"/status", nil)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

    // 3. 控制设备
    controlReq := map[string]interface{}{
        "command": "start_scanning",
        "params": map[string]interface{}{
            "frequency": "24GHz",
            "range":     "1000m",
        },
    }

    resp, err = suite.makeRequest("POST", "/api/v1/devices/"+deviceID+"/control", controlReq)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func (suite *NavlockE2ETestSuite) TestShipTrackingFlow() {
    // 1. 开始船舶跟踪
    shipID := "e2e-test-ship"
    trackReq := map[string]interface{}{
        "ship_id": shipID,
        "position": map[string]float64{
            "x": 100.0,
            "y": 200.0,
            "z": 0.0,
        },
        "speed": 10.5,
    }

    resp, err := suite.makeRequest("POST", "/api/v1/ships/"+shipID+"/track", trackReq)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

    // 2. 查询船舶位置
    resp, err = suite.makeRequest("GET", "/api/v1/ships/"+shipID+"/position", nil)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

    // 3. 检查船舶是否在禁停区域
    resp, err = suite.makeRequest("GET", "/api/v1/ships/"+shipID+"/restricted-area-check", nil)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func (suite *NavlockE2ETestSuite) makeRequest(method, path string, body interface{}) (*http.Response, error) {
    var req *http.Request
    var err error

    if body != nil {
        jsonBody, err := json.Marshal(body)
        if err != nil {
            return nil, err
        }
        req, err = http.NewRequestWithContext(suite.ctx, method, suite.baseURL+path, bytes.NewBuffer(jsonBody))
        req.Header.Set("Content-Type", "application/json")
    } else {
        req, err = http.NewRequestWithContext(suite.ctx, method, suite.baseURL+path, nil)
    }

    if err != nil {
        return nil, err
    }

    return suite.client.Do(req)
}

func TestNavlockE2ETestSuite(t *testing.T) {
    suite.Run(t, new(NavlockE2ETestSuite))
}
```

## 6.5 蓝绿部署

### 6.5.1 蓝绿部署配置

```yaml
# 蓝绿部署配置
apiVersion: argoproj.io/v1alpha1
kind: Rollout
metadata:
  name: navlock-service-rollout
  namespace: navlock-system
spec:
  replicas: 3
  selector:
    matchLabels:
      app: navlock-service
  template:
    metadata:
      labels:
        app: navlock-service
    spec:
      containers:
      - name: navlock-service
        image: registry.example.com/navlock-service:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: navlock-secrets
              key: database-url
        - name: REDIS_URL
          valueFrom:
            configMapKeyRef:
              name: navlock-config
              key: redis-url
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
  strategy:
    blueGreen:
      activeService: navlock-service-active
      previewService: navlock-service-preview
      autoPromotionEnabled: false
      scaleDownDelaySeconds: 30
      prePromotionAnalysis:
        templates:
        - templateName: success-rate
        args:
        - name: service-name
          value: navlock-service-preview
      postPromotionAnalysis:
        templates:
        - templateName: success-rate
        args:
        - name: service-name
          value: navlock-service-active
---
apiVersion: v1
kind: Service
metadata:
  name: navlock-service-active
  namespace: navlock-system
spec:
  type: ClusterIP
  ports:
  - port: 8080
    targetPort: 8080
    protocol: TCP
    name: http
  selector:
    app: navlock-service
---
apiVersion: v1
kind: Service
metadata:
  name: navlock-service-preview
  namespace: navlock-system
spec:
  type: ClusterIP
  ports:
  - port: 8080
    targetPort: 8080
    protocol: TCP
    name: http
  selector:
    app: navlock-service
```

### 6.5.2 蓝绿部署脚本

```bash
#!/bin/bash
# blue-green-deploy.sh

set -e

# 配置变量
SERVICE_NAME="navlock-service"
NAMESPACE="navlock-system"
NEW_VERSION="v1.2.0"
REGISTRY="registry.example.com"

echo "开始蓝绿部署..."

# 1. 构建新版本镜像
echo "构建新版本镜像..."
docker build -t $REGISTRY/$SERVICE_NAME:$NEW_VERSION .
docker push $REGISTRY/$SERVICE_NAME:$NEW_VERSION

# 2. 更新预览服务
echo "更新预览服务..."
kubectl set image deployment/$SERVICE_NAME $SERVICE_NAME=$REGISTRY/$SERVICE_NAME:$NEW_VERSION -n $NAMESPACE

# 3. 等待预览服务就绪
echo "等待预览服务就绪..."
kubectl rollout status deployment/$SERVICE_NAME -n $NAMESPACE --timeout=300s

# 4. 运行预部署测试
echo "运行预部署测试..."
./run-pre-deployment-tests.sh

# 5. 切换到新版本
echo "切换到新版本..."
kubectl patch service navlock-service-active -n $NAMESPACE -p '{"spec":{"selector":{"version":"'$NEW_VERSION'"}}}'

# 6. 运行后部署测试
echo "运行后部署测试..."
./run-post-deployment-tests.sh

# 7. 清理旧版本
echo "清理旧版本..."
kubectl delete deployment $SERVICE_NAME-old -n $NAMESPACE --ignore-not-found=true

echo "蓝绿部署完成！"
```

## 6.6 监控和告警

### 6.6.1 部署监控

```yaml
# 部署监控配置
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: navlock-deployment-monitor
  namespace: navlock-monitoring
spec:
  selector:
    matchLabels:
      app: navlock-service
  endpoints:
  - port: http
    path: /metrics
    interval: 30s
    scrapeTimeout: 10s
---
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: navlock-deployment-alerts
  namespace: navlock-monitoring
spec:
  groups:
  - name: navlock.deployment.rules
    rules:
    - alert: NavlockDeploymentFailed
      expr: kube_deployment_status_replicas_unavailable{deployment="navlock-service"} > 0
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "Navlock服务部署失败"
        description: "部署 {{ $labels.deployment }} 有 {{ $value }} 个副本不可用"
    
    - alert: NavlockHighMemoryUsage
      expr: (container_memory_usage_bytes{container="navlock-service"} / container_spec_memory_limit_bytes{container="navlock-service"}) > 0.85
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Navlock服务内存使用率过高"
        description: "内存使用率超过85%，当前值: {{ $value | humanizePercentage }}"
    
    - alert: NavlockHighCPUUsage
      expr: rate(container_cpu_usage_seconds_total{container="navlock-service"}[5m]) > 0.8
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Navlock服务CPU使用率过高"
        description: "CPU使用率超过80%，当前值: {{ $value | humanizePercentage }}"
```

### 6.6.2 部署指标

```go
// 部署指标收集器
type DeploymentMetricsCollector struct {
    deploymentCounter   prometheus.Counter
    deploymentDuration  prometheus.Histogram
    deploymentStatus    prometheus.Gauge
    logger             *zap.Logger
    tracer             trace.Tracer
}

func NewDeploymentMetricsCollector(reg prometheus.Registerer) *DeploymentMetricsCollector {
    deploymentCounter := prometheus.NewCounter(prometheus.CounterOpts{
        Name: "navlock_deployments_total",
        Help: "Total number of deployments",
    })
    
    deploymentDuration := prometheus.NewHistogram(prometheus.HistogramOpts{
        Name:    "navlock_deployment_duration_seconds",
        Help:    "Duration of deployments in seconds",
        Buckets: prometheus.DefBuckets,
    })
    
    deploymentStatus := prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "navlock_deployment_status",
        Help: "Current deployment status (1=success, 0=failure)",
    })
    
    reg.MustRegister(deploymentCounter, deploymentDuration, deploymentStatus)
    
    return &DeploymentMetricsCollector{
        deploymentCounter:  deploymentCounter,
        deploymentDuration: deploymentDuration,
        deploymentStatus:   deploymentStatus,
    }
}

// 记录部署开始
func (dmc *DeploymentMetricsCollector) RecordDeploymentStart(ctx context.Context, version string) {
    ctx, span := dmc.tracer.Start(ctx, "record_deployment_start")
    defer span.End()
    
    span.SetAttributes(attribute.String("version", version))
    
    dmc.deploymentCounter.Inc()
    
    dmc.logger.Info("部署开始",
        zap.String("version", version),
    )
}

// 记录部署完成
func (dmc *DeploymentMetricsCollector) RecordDeploymentComplete(ctx context.Context, version string, duration time.Duration, success bool) {
    ctx, span := dmc.tracer.Start(ctx, "record_deployment_complete")
    defer span.End()
    
    span.SetAttributes(
        attribute.String("version", version),
        attribute.Duration("duration", duration),
        attribute.Bool("success", success),
    )
    
    dmc.deploymentDuration.Observe(duration.Seconds())
    
    if success {
        dmc.deploymentStatus.Set(1)
    } else {
        dmc.deploymentStatus.Set(0)
    }
    
    dmc.logger.Info("部署完成",
        zap.String("version", version),
        zap.Duration("duration", duration),
        zap.Bool("success", success),
    )
}
```

## 6.7 总结

DevOps实践为葛洲坝船闸导航系统提供了现代化的软件交付能力。通过CI/CD流水线、GitOps、自动化测试、蓝绿部署等技术，系统获得了：

**核心优势：**
1. **自动化交付**：快速、可靠的软件发布
2. **质量保证**：全面的测试覆盖和验证
3. **风险控制**：蓝绿部署降低发布风险
4. **可观测性**：部署过程全程监控
5. **团队协作**：开发运维一体化

**实施建议：**
1. 渐进式实施，先基础CI/CD再高级特性
2. 建立完善的测试体系和监控
3. 制定详细的部署策略和回滚方案
4. 培训团队掌握DevOps工具和实践
5. 建立持续改进的文化

这个DevOps实践为系统的快速迭代和稳定运行提供了强有力的支撑，是现代化软件交付的重要组成部分。 