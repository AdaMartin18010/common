# 临床系统

## 概述

临床系统是医疗信息系统的核心组成部分，包括电子病历系统(EMR)、临床决策支持系统(CDSS)、医疗设备集成等。本章将详细介绍临床系统的设计、实现和集成方案。

## 1. 电子病历系统 (EMR)

### 1.1 病历数据结构

```go
// 电子病历核心数据结构
type ElectronicMedicalRecord struct {
    ID          string    `json:"id"`
    PatientID   string    `json:"patientId"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
    Status      string    `json:"status"` // active, archived, deleted
    
    // 基本信息
    Demographics *Demographics `json:"demographics"`
    
    // 临床信息
    Allergies    []Allergy     `json:"allergies"`
    Medications  []Medication  `json:"medications"`
    Diagnoses    []Diagnosis   `json:"diagnoses"`
    Procedures   []Procedure   `json:"procedures"`
    LabResults   []LabResult   `json:"labResults"`
    Imaging      []Imaging     `json:"imaging"`
    VitalSigns   []VitalSigns  `json:"vitalSigns"`
    Notes        []ClinicalNote `json:"notes"`
    
    // 安全信息
    AccessLog    []AccessLog   `json:"accessLog"`
    AuditTrail   []AuditEvent  `json:"auditTrail"`
}

type Demographics struct {
    Name         string    `json:"name"`
    DateOfBirth  time.Time `json:"dateOfBirth"`
    Gender       string    `json:"gender"`
    Race         string    `json:"race"`
    Ethnicity    string    `json:"ethnicity"`
    Address      Address   `json:"address"`
    Phone        string    `json:"phone"`
    Email        string    `json:"email"`
    EmergencyContact EmergencyContact `json:"emergencyContact"`
}

type Allergy struct {
    ID          string    `json:"id"`
    Allergen    string    `json:"allergen"`
    Reaction    string    `json:"reaction"`
    Severity    string    `json:"severity"` // mild, moderate, severe
    OnsetDate   time.Time `json:"onsetDate"`
    Status      string    `json:"status"` // active, resolved
}

type ClinicalNote struct {
    ID          string    `json:"id"`
    Type        string    `json:"type"` // progress, assessment, plan
    Author      string    `json:"author"`
    Content     string    `json:"content"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
    Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
    ID       string `json:"id"`
    Type     string `json:"type"` // image, document, audio
    Filename string `json:"filename"`
    URL      string `json:"url"`
    Size     int64  `json:"size"`
}
```

### 1.2 EMR服务实现

```go
// EMR服务
type EMRService struct {
    db          *sql.DB
    cache       *redis.Client
    fileStorage FileStorage
    encryption  EncryptionService
}

type FileStorage interface {
    Upload(data []byte, filename string) (string, error)
    Download(url string) ([]byte, error)
    Delete(url string) error
}

type EncryptionService interface {
    Encrypt(data []byte) ([]byte, error)
    Decrypt(data []byte) ([]byte, error)
}

func NewEMRService(db *sql.DB, cache *redis.Client, storage FileStorage, encryption EncryptionService) *EMRService {
    return &EMRService{
        db:         db,
        cache:      cache,
        fileStorage: storage,
        encryption: encryption,
    }
}

// 创建病历
func (s *EMRService) CreateEMR(patientID string, demographics *Demographics) (*ElectronicMedicalRecord, error) {
    emr := &ElectronicMedicalRecord{
        ID:        generateUUID(),
        PatientID: patientID,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Status:    "active",
        Demographics: demographics,
    }
    
    // 加密敏感数据
    encryptedData, err := s.encryption.Encrypt([]byte(emr.Demographics.Name))
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt patient name: %w", err)
    }
    
    // 存储到数据库
    query := `
        INSERT INTO emr (id, patient_id, demographics, created_at, updated_at, status)
        VALUES (```latex
1,
```2, ```latex
3,
```4, ```latex
5,
```6)
    `
    
    demographicsJSON, _ := json.Marshal(demographics)
    _, err = s.db.Exec(query, emr.ID, emr.PatientID, demographicsJSON, emr.CreatedAt, emr.UpdatedAt, emr.Status)
    if err != nil {
        return nil, fmt.Errorf("failed to create EMR: %w", err)
    }
    
    // 缓存病历
    s.cache.Set(context.Background(), fmt.Sprintf("emr:%s", emr.ID), emr, 24*time.Hour)
    
    return emr, nil
}

// 获取病历
func (s *EMRService) GetEMR(emrID string) (*ElectronicMedicalRecord, error) {
    // 先从缓存获取
    cached, err := s.cache.Get(context.Background(), fmt.Sprintf("emr:%s", emrID)).Result()
    if err == nil {
        var emr ElectronicMedicalRecord
        if err := json.Unmarshal([]byte(cached), &emr); err == nil {
            return &emr, nil
        }
    }
    
    // 从数据库获取
    query := `
        SELECT id, patient_id, demographics, created_at, updated_at, status
        FROM emr WHERE id = $1
    `
    
    var emr ElectronicMedicalRecord
    var demographicsJSON []byte
    
    err = s.db.QueryRow(query, emrID).Scan(
        &emr.ID, &emr.PatientID, &demographicsJSON, &emr.CreatedAt, &emr.UpdatedAt, &emr.Status,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to get EMR: %w", err)
    }
    
    // 解析人口统计学信息
    if err := json.Unmarshal(demographicsJSON, &emr.Demographics); err != nil {
        return nil, fmt.Errorf("failed to unmarshal demographics: %w", err)
    }
    
    // 获取相关数据
    if err := s.loadEMRRelatedData(&emr); err != nil {
        return nil, fmt.Errorf("failed to load related data: %w", err)
    }
    
    // 更新缓存
    emrJSON, _ := json.Marshal(emr)
    s.cache.Set(context.Background(), fmt.Sprintf("emr:%s", emrID), emrJSON, 24*time.Hour)
    
    return &emr, nil
}

func (s *EMRService) loadEMRRelatedData(emr *ElectronicMedicalRecord) error {
    // 加载过敏信息
    allergies, err := s.getAllergies(emr.ID)
    if err != nil {
        return err
    }
    emr.Allergies = allergies
    
    // 加载用药信息
    medications, err := s.getMedications(emr.ID)
    if err != nil {
        return err
    }
    emr.Medications = medications
    
    // 加载诊断信息
    diagnoses, err := s.getDiagnoses(emr.ID)
    if err != nil {
        return err
    }
    emr.Diagnoses = diagnoses
    
    return nil
}

// 添加临床记录
func (s *EMRService) AddClinicalNote(emrID string, note *ClinicalNote) error {
    note.ID = generateUUID()
    note.CreatedAt = time.Now()
    note.UpdatedAt = time.Now()
    
    // 处理附件
    for i := range note.Attachments {
        if note.Attachments[i].ID == "" {
            note.Attachments[i].ID = generateUUID()
        }
    }
    
    query := `
        INSERT INTO clinical_notes (id, emr_id, type, author, content, created_at, updated_at)
        VALUES (```latex
1,
```2, ```latex
3,
```4, ```latex
5,
```6, $7)
    `
    
    _, err := s.db.Exec(query, note.ID, emrID, note.Type, note.Author, note.Content, note.CreatedAt, note.UpdatedAt)
    if err != nil {
        return fmt.Errorf("failed to add clinical note: %w", err)
    }
    
    // 清除缓存
    s.cache.Del(context.Background(), fmt.Sprintf("emr:%s", emrID))
    
    return nil
}
```

## 2. 临床决策支持系统 (CDSS)

### 2.1 规则引擎

```go
// 临床决策支持系统
type ClinicalDecisionSupportSystem struct {
    rulesEngine *RulesEngine
    knowledgeBase *KnowledgeBase
    mlModels     map[string]*MLModel
}

type RulesEngine struct {
    rules []ClinicalRule
}

type ClinicalRule struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Category    string                 `json:"category"`
    Priority    int                    `json:"priority"`
    Conditions  []Condition            `json:"conditions"`
    Actions     []Action               `json:"actions"`
    Evidence    []Evidence             `json:"evidence"`
    Metadata    map[string]interface{} `json:"metadata"`
}

type Condition struct {
    Field    string      `json:"field"`
    Operator string      `json:"operator"` // eq, ne, gt, lt, gte, lte, in, contains
    Value    interface{} `json:"value"`
    Logic    string      `json:"logic"` // AND, OR
}

type Action struct {
    Type        string                 `json:"type"` // alert, recommendation, order
    Message     string                 `json:"message"`
    Severity    string                 `json:"severity"` // low, medium, high, critical
    Data        map[string]interface{} `json:"data"`
    AutoExecute bool                   `json:"autoExecute"`
}

type Evidence struct {
    Source      string  `json:"source"`
    Level       string  `json:"level"` // A, B, C, D
    Description string  `json:"description"`
    URL         string  `json:"url"`
}

func (cdss *ClinicalDecisionSupportSystem) EvaluatePatient(patient *Patient, context *ClinicalContext) ([]Decision, error) {
    var decisions []Decision
    
    // 基于规则的评估
    ruleDecisions, err := cdss.rulesEngine.Evaluate(patient, context)
    if err != nil {
        return nil, fmt.Errorf("failed to evaluate rules: %w", err)
    }
    decisions = append(decisions, ruleDecisions...)
    
    // 基于知识库的评估
    kbDecisions, err := cdss.knowledgeBase.Query(patient, context)
    if err != nil {
        return nil, fmt.Errorf("failed to query knowledge base: %w", err)
    }
    decisions = append(decisions, kbDecisions...)
    
    // 基于机器学习的评估
    for modelName, model := range cdss.mlModels {
        mlDecisions, err := model.Predict(patient, context)
        if err != nil {
            log.Printf("Failed to get predictions from model %s: %v", modelName, err)
            continue
        }
        decisions = append(decisions, mlDecisions...)
    }
    
    // 去重和排序
    decisions = cdss.deduplicateAndSort(decisions)
    
    return decisions, nil
}

type Decision struct {
    ID          string                 `json:"id"`
    Type        string                 `json:"type"`
    Category    string                 `json:"category"`
    Priority    int                    `json:"priority"`
    Message     string                 `json:"message"`
    Severity    string                 `json:"severity"`
    Confidence  float64                `json:"confidence"`
    Source      string                 `json:"source"`
    Evidence    []Evidence             `json:"evidence"`
    Actions     []Action               `json:"actions"`
    Metadata    map[string]interface{} `json:"metadata"`
}

type ClinicalContext struct {
    CurrentVisit    *Visit     `json:"currentVisit"`
    RecentVisits    []Visit    `json:"recentVisits"`
    CurrentMedications []Medication `json:"currentMedications"`
    LabResults      []LabResult `json:"labResults"`
    VitalSigns      *VitalSigns `json:"vitalSigns"`
    Allergies       []Allergy   `json:"allergies"`
}

func (re *RulesEngine) Evaluate(patient *Patient, context *ClinicalContext) ([]Decision, error) {
    var decisions []Decision
    
    for _, rule := range re.rules {
        if re.evaluateConditions(rule.Conditions, patient, context) {
            decision := Decision{
                ID:         generateUUID(),
                Type:       "rule_based",
                Category:   rule.Category,
                Priority:   rule.Priority,
                Message:    rule.Description,
                Severity:   "medium",
                Confidence: 1.0,
                Source:     "rules_engine",
                Evidence:   rule.Evidence,
                Actions:    rule.Actions,
                Metadata:   rule.Metadata,
            }
            decisions = append(decisions, decision)
        }
    }
    
    return decisions, nil
}

func (re *RulesEngine) evaluateConditions(conditions []Condition, patient *Patient, context *ClinicalContext) bool {
    if len(conditions) == 0 {
        return true
    }
    
    result := re.evaluateCondition(conditions[0], patient, context)
    
    for i := 1; i < len(conditions); i++ {
        condition := conditions[i]
        conditionResult := re.evaluateCondition(condition, patient, context)
        
        if condition.Logic == "AND" {
            result = result && conditionResult
        } else if condition.Logic == "OR" {
            result = result || conditionResult
        }
    }
    
    return result
}

func (re *RulesEngine) evaluateCondition(condition Condition, patient *Patient, context *ClinicalContext) bool {
    // 获取字段值
    fieldValue := re.getFieldValue(condition.Field, patient, context)
    
    // 根据操作符进行比较
    switch condition.Operator {
    case "eq":
        return reflect.DeepEqual(fieldValue, condition.Value)
    case "ne":
        return !reflect.DeepEqual(fieldValue, condition.Value)
    case "gt":
        return re.compareValues(fieldValue, condition.Value, "gt")
    case "lt":
        return re.compareValues(fieldValue, condition.Value, "lt")
    case "gte":
        return re.compareValues(fieldValue, condition.Value, "gte")
    case "lte":
        return re.compareValues(fieldValue, condition.Value, "lte")
    case "in":
        return re.isIn(fieldValue, condition.Value)
    case "contains":
        return re.contains(fieldValue, condition.Value)
    default:
        return false
    }
}

func (re *RulesEngine) getFieldValue(field string, patient *Patient, context *ClinicalContext) interface{} {
    // 根据字段路径获取值
    switch field {
    case "patient.age":
        if patient.BirthDate != "" {
            birthDate, _ := time.Parse("2006-01-02", patient.BirthDate)
            return int(time.Since(birthDate).Hours() / 24 / 365.25)
        }
        return nil
    case "patient.gender":
        return patient.Gender
    case "vitals.heartRate":
        if context.VitalSigns != nil {
            return context.VitalSigns.HeartRate
        }
        return nil
    case "vitals.bloodPressureSystolic":
        if context.VitalSigns != nil {
            return context.VitalSigns.BloodPressureSystolic
        }
        return nil
    case "medications.count":
        return len(context.CurrentMedications)
    default:
        return nil
    }
}
```

### 2.2 药物相互作用检查

```go
// 药物相互作用检查器
type DrugInteractionChecker struct {
    interactions map[string][]DrugInteraction
    db           *sql.DB
}

type DrugInteraction struct {
    ID              string   `json:"id"`
    Drug1           string   `json:"drug1"`
    Drug2           string   `json:"drug2"`
    Severity        string   `json:"severity"` // major, moderate, minor
    Description     string   `json:"description"`
    Mechanism       string   `json:"mechanism"`
    Management      string   `json:"management"`
    Evidence        string   `json:"evidence"`
    References      []string `json:"references"`
}

func (dic *DrugInteractionChecker) CheckInteractions(medications []Medication) ([]DrugInteraction, error) {
    var interactions []DrugInteraction
    
    // 检查所有药物对之间的相互作用
    for i := 0; i < len(medications); i++ {
        for j := i + 1; j < len(medications); j++ {
            drug1 := medications[i].Name
            drug2 := medications[j].Name
            
            // 查询相互作用
            interaction, err := dic.getInteraction(drug1, drug2)
            if err != nil {
                log.Printf("Failed to get interaction between %s and %s: %v", drug1, drug2, err)
                continue
            }
            
            if interaction != nil {
                interactions = append(interactions, *interaction)
            }
        }
    }
    
    // 按严重程度排序
    sort.Slice(interactions, func(i, j int) bool {
        severityOrder := map[string]int{"major": 3, "moderate": 2, "minor": 1}
        return severityOrder[interactions[i].Severity] > severityOrder[interactions[j].Severity]
    })
    
    return interactions, nil
}

func (dic *DrugInteractionChecker) getInteraction(drug1, drug2 string) (*DrugInteraction, error) {
    query := `
        SELECT id, drug1, drug2, severity, description, mechanism, management, evidence, references
        FROM drug_interactions 
        WHERE (drug1 = ```latex
1 AND drug2 =
```2) OR (drug1 = ```latex
2 AND drug2 =
```1)
    `
    
    var interaction DrugInteraction
    var referencesJSON []byte
    
    err := dic.db.QueryRow(query, drug1, drug2).Scan(
        &interaction.ID, &interaction.Drug1, &interaction.Drug2, &interaction.Severity,
        &interaction.Description, &interaction.Mechanism, &interaction.Management,
        &interaction.Evidence, &referencesJSON,
    )
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, fmt.Errorf("failed to query drug interaction: %w", err)
    }
    
    // 解析参考文献
    if err := json.Unmarshal(referencesJSON, &interaction.References); err != nil {
        return nil, fmt.Errorf("failed to unmarshal references: %w", err)
    }
    
    return &interaction, nil
}

// 过敏检查
func (dic *DrugInteractionChecker) CheckAllergies(medications []Medication, allergies []Allergy) ([]AllergyAlert, error) {
    var alerts []AllergyAlert
    
    for _, medication := range medications {
        for _, allergy := range allergies {
            if allergy.Status == "active" && dic.isAllergic(medication.Name, allergy.Allergen) {
                alert := AllergyAlert{
                    Medication: medication.Name,
                    Allergen:   allergy.Allergen,
                    Reaction:   allergy.Reaction,
                    Severity:   allergy.Severity,
                    Message:    fmt.Sprintf("患者对 %s 过敏，反应: %s", allergy.Allergen, allergy.Reaction),
                }
                alerts = append(alerts, alert)
            }
        }
    }
    
    return alerts, nil
}

type AllergyAlert struct {
    Medication string `json:"medication"`
    Allergen   string `json:"allergen"`
    Reaction   string `json:"reaction"`
    Severity   string `json:"severity"`
    Message    string `json:"message"`
}

func (dic *DrugInteractionChecker) isAllergic(medication, allergen string) bool {
    // 简化的过敏检查逻辑
    // 实际应用中需要更复杂的药物分类和过敏原匹配
    return strings.Contains(strings.ToLower(medication), strings.ToLower(allergen))
}
```

## 3. 医疗设备集成

### 3.1 设备通信协议

```go
// 医疗设备集成
type MedicalDeviceIntegration struct {
    devices map[string]MedicalDevice
    protocols map[string]DeviceProtocol
}

type MedicalDevice struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Type        string            `json:"type"`
    Manufacturer string           `json:"manufacturer"`
    Model       string            `json:"model"`
    Protocol    string            `json:"protocol"`
    Connection  DeviceConnection  `json:"connection"`
    Status      string            `json:"status"`
    Capabilities []string         `json:"capabilities"`
    Metadata    map[string]interface{} `json:"metadata"`
}

type DeviceConnection struct {
    Type     string `json:"type"` // tcp, serial, usb, bluetooth
    Address  string `json:"address"`
    Port     int    `json:"port"`
    BaudRate int    `json:"baudRate,omitempty"`
    Timeout  int    `json:"timeout"`
}

type DeviceProtocol interface {
    Connect(connection DeviceConnection) error
    Disconnect() error
    SendCommand(command []byte) ([]byte, error)
    ReadData() ([]byte, error)
    IsConnected() bool
}

// HL7协议实现
type HL7Protocol struct {
    conn net.Conn
    connected bool
}

func (p *HL7Protocol) Connect(connection DeviceConnection) error {
    var err error
    p.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", connection.Address, connection.Port))
    if err != nil {
        return fmt.Errorf("failed to connect to device: %w", err)
    }
    
    p.connected = true
    return nil
}

func (p *HL7Protocol) Disconnect() error {
    if p.conn != nil {
        p.connected = false
        return p.conn.Close()
    }
    return nil
}

func (p *HL7Protocol) SendCommand(command []byte) ([]byte, error) {
    if !p.connected {
        return nil, fmt.Errorf("device not connected")
    }
    
    // 发送HL7消息
    _, err := p.conn.Write(command)
    if err != nil {
        return nil, fmt.Errorf("failed to send command: %w", err)
    }
    
    // 读取响应
    response := make([]byte, 1024)
    n, err := p.conn.Read(response)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }
    
    return response[:n], nil
}

func (p *HL7Protocol) ReadData() ([]byte, error) {
    if !p.connected {
        return nil, fmt.Errorf("device not connected")
    }
    
    data := make([]byte, 1024)
    n, err := p.conn.Read(data)
    if err != nil {
        return nil, fmt.Errorf("failed to read data: %w", err)
    }
    
    return data[:n], nil
}

func (p *HL7Protocol) IsConnected() bool {
    return p.connected
}

// DICOM协议实现
type DICOMProtocol struct {
    conn net.Conn
    connected bool
}

func (p *DICOMProtocol) Connect(connection DeviceConnection) error {
    var err error
    p.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", connection.Address, connection.Port))
    if err != nil {
        return fmt.Errorf("failed to connect to DICOM device: %w", err)
    }
    
    p.connected = true
    return nil
}

func (p *DICOMProtocol) Disconnect() error {
    if p.conn != nil {
        p.connected = false
        return p.conn.Close()
    }
    return nil
}

func (p *DICOMProtocol) SendCommand(command []byte) ([]byte, error) {
    if !p.connected {
        return nil, fmt.Errorf("DICOM device not connected")
    }
    
    // 发送DICOM命令
    _, err := p.conn.Write(command)
    if err != nil {
        return nil, fmt.Errorf("failed to send DICOM command: %w", err)
    }
    
    // 读取响应
    response := make([]byte, 4096)
    n, err := p.conn.Read(response)
    if err != nil {
        return nil, fmt.Errorf("failed to read DICOM response: %w", err)
    }
    
    return response[:n], nil
}

func (p *DICOMProtocol) ReadData() ([]byte, error) {
    if !p.connected {
        return nil, fmt.Errorf("DICOM device not connected")
    }
    
    data := make([]byte, 4096)
    n, err := p.conn.Read(data)
    if err != nil {
        return nil, fmt.Errorf("failed to read DICOM data: %w", err)
    }
    
    return data[:n], nil
}

func (p *DICOMProtocol) IsConnected() bool {
    return p.connected
}
```

### 3.2 设备数据采集

```go
// 设备数据采集器
type DeviceDataCollector struct {
    devices map[string]MedicalDevice
    protocols map[string]DeviceProtocol
    dataProcessor *DeviceDataProcessor
    storage       DeviceDataStorage
}

type DeviceDataProcessor struct {
    parsers map[string]DataParser
}

type DataParser interface {
    Parse(data []byte) (interface{}, error)
    Validate(data []byte) error
}

type DeviceDataStorage interface {
    Store(deviceID string, data interface{}) error
    Query(deviceID string, startTime, endTime time.Time) ([]interface{}, error)
}

func (ddc *DeviceDataCollector) CollectData(deviceID string) error {
    device, exists := ddc.devices[deviceID]
    if !exists {
        return fmt.Errorf("device %s not found", deviceID)
    }
    
    protocol, exists := ddc.protocols[deviceID]
    if !exists {
        return fmt.Errorf("protocol for device %s not found", deviceID)
    }
    
    // 读取设备数据
    data, err := protocol.ReadData()
    if err != nil {
        return fmt.Errorf("failed to read data from device: %w", err)
    }
    
    // 解析数据
    parsedData, err := ddc.dataProcessor.Parse(device.Type, data)
    if err != nil {
        return fmt.Errorf("failed to parse device data: %w", err)
    }
    
    // 存储数据
    if err := ddc.storage.Store(deviceID, parsedData); err != nil {
        return fmt.Errorf("failed to store device data: %w", err)
    }
    
    return nil
}

func (ddp *DeviceDataProcessor) Parse(deviceType string, data []byte) (interface{}, error) {
    parser, exists := ddp.parsers[deviceType]
    if !exists {
        return nil, fmt.Errorf("parser for device type %s not found", deviceType)
    }
    
    return parser.Parse(data)
}

// 生命体征数据解析器
type VitalSignsParser struct{}

func (p *VitalSignsParser) Parse(data []byte) (interface{}, error) {
    var vitals VitalSigns
    
    // 解析设备数据格式
    // 这里假设数据是JSON格式
    if err := json.Unmarshal(data, &vitals); err != nil {
        return nil, fmt.Errorf("failed to unmarshal vital signs: %w", err)
    }
    
    // 验证数据
    if err := p.Validate(data); err != nil {
        return nil, fmt.Errorf("vital signs validation failed: %w", err)
    }
    
    return &vitals, nil
}

func (p *VitalSignsParser) Validate(data []byte) error {
    var vitals VitalSigns
    if err := json.Unmarshal(data, &vitals); err != nil {
        return err
    }
    
    // 验证心率范围
    if vitals.HeartRate < 30 || vitals.HeartRate > 200 {
        return fmt.Errorf("heart rate out of range: %d", vitals.HeartRate)
    }
    
    // 验证血压范围
    if vitals.BloodPressureSystolic < 70 || vitals.BloodPressureSystolic > 200 {
        return fmt.Errorf("systolic blood pressure out of range: %d", vitals.BloodPressureSystolic)
    }
    
    // 验证体温范围
    if vitals.Temperature < 30 || vitals.Temperature > 45 {
        return fmt.Errorf("temperature out of range: %.1f", vitals.Temperature)
    }
    
    return nil
}

// 实验室数据解析器
type LabDataParser struct{}

func (p *LabDataParser) Parse(data []byte) (interface{}, error) {
    var labResult LabResult
    
    if err := json.Unmarshal(data, &labResult); err != nil {
        return nil, fmt.Errorf("failed to unmarshal lab result: %w", err)
    }
    
    if err := p.Validate(data); err != nil {
        return nil, fmt.Errorf("lab result validation failed: %w", err)
    }
    
    return &labResult, nil
}

func (p *LabDataParser) Validate(data []byte) error {
    var labResult LabResult
    if err := json.Unmarshal(data, &labResult); err != nil {
        return err
    }
    
    // 验证必要字段
    if labResult.TestName == "" {
        return fmt.Errorf("test name is required")
    }
    
    if labResult.Value == 0 {
        return fmt.Errorf("test value cannot be zero")
    }
    
    return nil
}
```

## 4. 临床工作流

### 4.1 工作流引擎

```go
// 临床工作流引擎
type ClinicalWorkflowEngine struct {
    workflows map[string]*ClinicalWorkflow
    executor  *WorkflowExecutor
    monitor   *WorkflowMonitor
}

type ClinicalWorkflow struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Version     string                 `json:"version"`
    Steps       []WorkflowStep         `json:"steps"`
    Triggers    []WorkflowTrigger      `json:"triggers"`
    Variables   map[string]interface{} `json:"variables"`
    Status      string                 `json:"status"`
}

type WorkflowStep struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Type        string                 `json:"type"` // task, decision, parallel, subprocess
    Action      string                 `json:"action"`
    Parameters  map[string]interface{} `json:"parameters"`
    Conditions  []Condition            `json:"conditions"`
    NextSteps   []string               `json:"nextSteps"`
    Timeout     int                    `json:"timeout"`
    RetryPolicy *RetryPolicy           `json:"retryPolicy"`
}

type WorkflowTrigger struct {
    Type        string                 `json:"type"` // event, schedule, manual
    Event       string                 `json:"event"`
    Schedule    string                 `json:"schedule"`
    Conditions  []Condition            `json:"conditions"`
    Parameters  map[string]interface{} `json:"parameters"`
}

type RetryPolicy struct {
    MaxAttempts int           `json:"maxAttempts"`
    Delay       time.Duration `json:"delay"`
    Backoff     string        `json:"backoff"` // linear, exponential
}

func (cwe *ClinicalWorkflowEngine) ExecuteWorkflow(workflowID string, context *WorkflowContext) (*WorkflowExecution, error) {
    workflow, exists := cwe.workflows[workflowID]
    if !exists {
        return nil, fmt.Errorf("workflow %s not found", workflowID)
    }
    
    execution := &WorkflowExecution{
        ID:         generateUUID(),
        WorkflowID: workflowID,
        Status:     "running",
        StartTime:  time.Now(),
        Context:    context,
        Steps:      make(map[string]*StepExecution),
    }
    
    // 执行工作流
    go cwe.executor.Execute(workflow, execution)
    
    return execution, nil
}

type WorkflowContext struct {
    PatientID   string                 `json:"patientId"`
    UserID      string                 `json:"userId"`
    Variables   map[string]interface{} `json:"variables"`
    Metadata    map[string]interface{} `json:"metadata"`
}

type WorkflowExecution struct {
    ID         string                    `json:"id"`
    WorkflowID string                    `json:"workflowId"`
    Status     string                    `json:"status"`
    StartTime  time.Time                 `json:"startTime"`
    EndTime    *time.Time                `json:"endTime,omitempty"`
    Context    *WorkflowContext          `json:"context"`
    Steps      map[string]*StepExecution `json:"steps"`
    Error      string                    `json:"error,omitempty"`
}

type StepExecution struct {
    ID         string                 `json:"id"`
    Status     string                 `json:"status"`
    StartTime  time.Time              `json:"startTime"`
    EndTime    *time.Time             `json:"endTime,omitempty"`
    Result     interface{}            `json:"result,omitempty"`
    Error      string                 `json:"error,omitempty"`
    Attempts   int                    `json:"attempts"`
    Metadata   map[string]interface{} `json:"metadata"`
}

// 工作流执行器
type WorkflowExecutor struct {
    taskHandlers map[string]TaskHandler
}

type TaskHandler func(context *WorkflowContext, parameters map[string]interface{}) (interface{}, error)

func (we *WorkflowExecutor) Execute(workflow *ClinicalWorkflow, execution *WorkflowExecution) {
    defer func() {
        execution.EndTime = &[]time.Time{time.Now()}[0]
        if execution.Status == "running" {
            execution.Status = "completed"
        }
    }()
    
    // 执行工作流步骤
    for _, step := range workflow.Steps {
        stepExecution := &StepExecution{
            ID:        step.ID,
            Status:    "running",
            StartTime: time.Now(),
        }
        execution.Steps[step.ID] = stepExecution
        
        // 执行步骤
        result, err := we.executeStep(step, execution.Context)
        if err != nil {
            stepExecution.Status = "failed"
            stepExecution.Error = err.Error()
            stepExecution.EndTime = &[]time.Time{time.Now()}[0]
            execution.Status = "failed"
            execution.Error = err.Error()
            return
        }
        
        stepExecution.Status = "completed"
        stepExecution.Result = result
        stepExecution.EndTime = &[]time.Time{time.Now()}[0]
    }
}

func (we *WorkflowExecutor) executeStep(step WorkflowStep, context *WorkflowContext) (interface{}, error) {
    handler, exists := we.taskHandlers[step.Action]
    if !exists {
        return nil, fmt.Errorf("task handler for action %s not found", step.Action)
    }
    
    return handler(context, step.Parameters)
}

// 注册任务处理器
func (we *WorkflowExecutor) RegisterHandler(action string, handler TaskHandler) {
    we.taskHandlers[action] = handler
}

// 示例任务处理器
func (we *WorkflowExecutor) registerDefaultHandlers() {
    // 创建实验室订单
    we.RegisterHandler("create_lab_order", func(context *WorkflowContext, parameters map[string]interface{}) (interface{}, error) {
        patientID := context.PatientID
        testName := parameters["testName"].(string)
        
        // 创建实验室订单的逻辑
        order := &LabOrder{
            ID:        generateUUID(),
            PatientID: patientID,
            TestName:  testName,
            Status:    "pending",
            CreatedAt: time.Now(),
        }
        
        return order, nil
    })
    
    // 发送通知
    we.RegisterHandler("send_notification", func(context *WorkflowContext, parameters map[string]interface{}) (interface{}, error) {
        recipient := parameters["recipient"].(string)
        message := parameters["message"].(string)
        
        // 发送通知的逻辑
        notification := &Notification{
            ID:        generateUUID(),
            Recipient: recipient,
            Message:   message,
            Status:    "sent",
            SentAt:    time.Now(),
        }
        
        return notification, nil
    })
}

type LabOrder struct {
    ID        string    `json:"id"`
    PatientID string    `json:"patientId"`
    TestName  string    `json:"testName"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"createdAt"`
}

type Notification struct {
    ID        string    `json:"id"`
    Recipient string    `json:"recipient"`
    Message   string    `json:"message"`
    Status    string    `json:"status"`
    SentAt    time.Time `json:"sentAt"`
}
```

## 总结

临床系统是医疗信息化的核心，通过电子病历、临床决策支持、设备集成和工作流管理，为医疗工作者提供全面的技术支持。

### 关键特性

1. **标准化**: 支持HL7、DICOM等医疗标准
2. **智能化**: 集成临床决策支持和药物相互作用检查
3. **集成性**: 支持多种医疗设备的数据采集
4. **工作流**: 提供灵活的临床工作流管理
5. **安全性**: 实现完整的访问控制和审计功能

### 技术栈

- **数据库**: PostgreSQL, Redis
- **消息队列**: Apache Kafka
- **协议**: HL7, DICOM, TCP/IP
- **加密**: AES-256-GCM
- **工作流**: 自定义工作流引擎

---

**相关链接**:

- [医疗数据管理](./01-Medical-Data-Management.md)
- [软件架构层](../02-Software-Architecture/README.md)
- [工作流系统](../10-Workflow-Systems/README.md)
- [数据安全](../08-Cybersecurity/README.md)
