# 01-在线学习平台 (Online Learning Platform)

## 1. 概述

### 1.1 定义与目标

在线学习平台是教育科技领域的核心系统，为学习者提供个性化、交互式的数字化学习体验。

**形式化定义**：
设 $U$ 为用户集合，$C$ 为课程集合，$L$ 为学习活动集合，$P$ 为进度集合，则学习平台函数 $f$ 定义为：

$$f: U \times C \rightarrow L \times P$$

其中：
- $U = \{u_1, u_2, ..., u_n\}$ 为学习者集合
- $C = \{c_1, c_2, ..., c_m\}$ 为课程集合
- $L = \{l_1, l_2, ..., l_k\}$ 为学习活动集合
- $P = \{p_1, p_2, ..., p_r\}$ 为学习进度集合

### 1.2 学习理论模型

**个性化学习模型**：
$$\text{PersonalizedScore}(u, c) = \alpha \cdot \text{Interest}(u, c) + \beta \cdot \text{Ability}(u, c) + \gamma \cdot \text{Progress}(u, c)$$

其中 $\alpha + \beta + \gamma = 1$ 为权重系数。

**学习路径优化**：
$$\text{OptimalPath}(u) = \arg\min_{p \in P} \sum_{i=1}^{n} \text{Difficulty}(p_i) - \text{Ability}(u, p_i)$$

## 2. 架构设计

### 2.1 微服务架构

```
┌─────────────────────────────────────┐
│           用户服务层                  │
├─────────────────────────────────────┤
│           课程服务层                  │
├─────────────────────────────────────┤
│           学习引擎层                  │
├─────────────────────────────────────┤
│           内容管理层                  │
├─────────────────────────────────────┤
│           分析推荐层                  │
└─────────────────────────────────────┘
```

### 2.2 核心组件

#### 2.2.1 用户服务

```go
// UserService 用户服务
type UserService struct {
    repo        UserRepository
    auth        AuthService
    profile     ProfileService
    preferences PreferenceService
    config      *UserServiceConfig
}

// UserServiceConfig 用户服务配置
type UserServiceConfig struct {
    MaxUsersPerPage int           `json:"max_users_per_page"`
    SessionTimeout  time.Duration `json:"session_timeout"`
    EnableMFA       bool          `json:"enable_mfa"`
    PasswordPolicy  PasswordPolicy `json:"password_policy"`
}

// User 用户模型
type User struct {
    ID          string                 `json:"id"`
    Email       string                 `json:"email"`
    Username    string                 `json:"username"`
    Role        UserRole               `json:"role"`
    Profile     *UserProfile           `json:"profile"`
    Preferences *UserPreferences       `json:"preferences"`
    Status      UserStatus             `json:"status"`
    CreatedAt   time.Time              `json:"created_at"`
    UpdatedAt   time.Time              `json:"updated_at"`
    Metadata    map[string]interface{} `json:"metadata"`
}

// UserRole 用户角色
type UserRole string

const (
    UserRoleStudent      UserRole = "student"
    UserRoleTeacher      UserRole = "teacher"
    UserRoleAdministrator UserRole = "administrator"
    UserRoleParent       UserRole = "parent"
)

// UserStatus 用户状态
type UserStatus string

const (
    UserStatusActive   UserStatus = "active"
    UserStatusInactive UserStatus = "inactive"
    UserStatusSuspended UserStatus = "suspended"
    UserStatusPending  UserStatus = "pending"
)

// UserProfile 用户档案
type UserProfile struct {
    FirstName    string    `json:"first_name"`
    LastName     string    `json:"last_name"`
    AvatarURL    string    `json:"avatar_url"`
    Bio          string    `json:"bio"`
    GradeLevel   *GradeLevel `json:"grade_level"`
    Subjects     []string  `json:"subjects"`
    LearningGoals []string `json:"learning_goals"`
    DateOfBirth  *time.Time `json:"date_of_birth"`
    Location     string    `json:"location"`
    Timezone     string    `json:"timezone"`
}

// GradeLevel 年级级别
type GradeLevel struct {
    Level     int    `json:"level"`
    Name      string `json:"name"`
    Category  string `json:"category"` // elementary, middle, high, college
}

// UserPreferences 用户偏好
type UserPreferences struct {
    Language      string                 `json:"language"`
    Theme         string                 `json:"theme"`
    Notifications NotificationSettings   `json:"notifications"`
    Accessibility AccessibilitySettings  `json:"accessibility"`
    Privacy       PrivacySettings        `json:"privacy"`
    Custom        map[string]interface{} `json:"custom"`
}

// NotificationSettings 通知设置
type NotificationSettings struct {
    EmailEnabled     bool `json:"email_enabled"`
    PushEnabled      bool `json:"push_enabled"`
    SMSEnabled       bool `json:"sms_enabled"`
    CourseUpdates    bool `json:"course_updates"`
    AssignmentDue    bool `json:"assignment_due"`
    GradeUpdates     bool `json:"grade_updates"`
    DiscussionReplies bool `json:"discussion_replies"`
}

// AccessibilitySettings 无障碍设置
type AccessibilitySettings struct {
    HighContrast    bool `json:"high_contrast"`
    LargeText       bool `json:"large_text"`
    ScreenReader    bool `json:"screen_reader"`
    KeyboardOnly    bool `json:"keyboard_only"`
    ReducedMotion   bool `json:"reduced_motion"`
}

// PrivacySettings 隐私设置
type PrivacySettings struct {
    ProfileVisible  bool `json:"profile_visible"`
    ProgressVisible bool `json:"progress_visible"`
    ActivityVisible bool `json:"activity_visible"`
    DataSharing     bool `json:"data_sharing"`
}

// CreateUser 创建用户
func (us *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // 验证请求
    if err := us.validateCreateRequest(req); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }
    
    // 检查用户是否已存在
    if exists, _ := us.repo.ExistsByEmail(ctx, req.Email); exists {
        return nil, errors.New("user already exists")
    }
    
    // 创建用户
    user := &User{
        ID:        uuid.New().String(),
        Email:     req.Email,
        Username:  req.Username,
        Role:      req.Role,
        Status:    UserStatusPending,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Metadata:  make(map[string]interface{}),
    }
    
    // 设置档案
    if req.Profile != nil {
        user.Profile = req.Profile
    }
    
    // 设置偏好
    if req.Preferences != nil {
        user.Preferences = req.Preferences
    } else {
        user.Preferences = us.getDefaultPreferences()
    }
    
    // 保存用户
    if err := us.repo.Create(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    // 发送欢迎邮件
    go us.sendWelcomeEmail(user)
    
    return user, nil
}

// validateCreateRequest 验证创建请求
func (us *UserService) validateCreateRequest(req *CreateUserRequest) error {
    if req.Email == "" {
        return errors.New("email is required")
    }
    
    if !us.isValidEmail(req.Email) {
        return errors.New("invalid email format")
    }
    
    if req.Username == "" {
        return errors.New("username is required")
    }
    
    if len(req.Username) < 3 {
        return errors.New("username must be at least 3 characters")
    }
    
    if req.Role == "" {
        return errors.New("role is required")
    }
    
    return nil
}

// isValidEmail 验证邮箱格式
func (us *UserService) isValidEmail(email string) bool {
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return emailRegex.MatchString(email)
}

// getDefaultPreferences 获取默认偏好
func (us *UserService) getDefaultPreferences() *UserPreferences {
    return &UserPreferences{
        Language: "en",
        Theme:    "light",
        Notifications: NotificationSettings{
            EmailEnabled:     true,
            PushEnabled:      true,
            SMSEnabled:       false,
            CourseUpdates:    true,
            AssignmentDue:    true,
            GradeUpdates:     true,
            DiscussionReplies: true,
        },
        Accessibility: AccessibilitySettings{
            HighContrast:  false,
            LargeText:     false,
            ScreenReader:  false,
            KeyboardOnly:  false,
            ReducedMotion: false,
        },
        Privacy: PrivacySettings{
            ProfileVisible:  true,
            ProgressVisible: true,
            ActivityVisible: false,
            DataSharing:     false,
        },
        Custom: make(map[string]interface{}),
    }
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
    Email       string           `json:"email"`
    Username    string           `json:"username"`
    Role        UserRole         `json:"role"`
    Profile     *UserProfile     `json:"profile,omitempty"`
    Preferences *UserPreferences `json:"preferences,omitempty"`
}

// sendWelcomeEmail 发送欢迎邮件
func (us *UserService) sendWelcomeEmail(user *User) {
    // 实现邮件发送逻辑
    log.Printf("Sending welcome email to %s", user.Email)
}
```

#### 2.2.2 课程服务

```go
// CourseService 课程服务
type CourseService struct {
    repo        CourseRepository
    content     ContentService
    enrollment  EnrollmentService
    progress    ProgressService
    config      *CourseServiceConfig
}

// CourseServiceConfig 课程服务配置
type CourseServiceConfig struct {
    MaxCoursesPerPage int           `json:"max_courses_per_page"`
    EnablePrerequisites bool        `json:"enable_prerequisites"`
    EnableCertificates  bool        `json:"enable_certificates"`
    AutoEnrollment      bool        `json:"auto_enrollment"`
}

// Course 课程模型
type Course struct {
    ID              string                 `json:"id"`
    Title           string                 `json:"title"`
    Description     string                 `json:"description"`
    InstructorID    string                 `json:"instructor_id"`
    Category        string                 `json:"category"`
    Level           CourseLevel            `json:"level"`
    Duration        time.Duration          `json:"duration"`
    Price           *Price                 `json:"price"`
    Prerequisites   []string               `json:"prerequisites"`
    LearningOutcomes []string              `json:"learning_outcomes"`
    Modules         []*Module              `json:"modules"`
    Status          CourseStatus           `json:"status"`
    EnrollmentCount int                    `json:"enrollment_count"`
    Rating          float64                `json:"rating"`
    ReviewCount     int                    `json:"review_count"`
    CreatedAt       time.Time              `json:"created_at"`
    UpdatedAt       time.Time              `json:"updated_at"`
    Metadata        map[string]interface{} `json:"metadata"`
}

// CourseLevel 课程级别
type CourseLevel string

const (
    CourseLevelBeginner CourseLevel = "beginner"
    CourseLevelIntermediate CourseLevel = "intermediate"
    CourseLevelAdvanced CourseLevel = "advanced"
)

// CourseStatus 课程状态
type CourseStatus string

const (
    CourseStatusDraft     CourseStatus = "draft"
    CourseStatusPublished CourseStatus = "published"
    CourseStatusArchived  CourseStatus = "archived"
    CourseStatusSuspended CourseStatus = "suspended"
)

// Price 价格模型
type Price struct {
    Amount   float64 `json:"amount"`
    Currency string  `json:"currency"`
    Type     PriceType `json:"type"`
}

// PriceType 价格类型
type PriceType string

const (
    PriceTypeFree  PriceType = "free"
    PriceTypePaid  PriceType = "paid"
    PriceTypeSubscription PriceType = "subscription"
)

// Module 课程模块
type Module struct {
    ID          string        `json:"id"`
    Title       string        `json:"title"`
    Description string        `json:"description"`
    Order       int           `json:"order"`
    Duration    time.Duration `json:"duration"`
    Lessons     []*Lesson     `json:"lessons"`
    Quiz        *Quiz         `json:"quiz,omitempty"`
    Assignment  *Assignment   `json:"assignment,omitempty"`
}

// Lesson 课程章节
type Lesson struct {
    ID          string        `json:"id"`
    Title       string        `json:"title"`
    Description string        `json:"description"`
    Order       int           `json:"order"`
    Duration    time.Duration `json:"duration"`
    Type        LessonType    `json:"type"`
    Content     *Content      `json:"content"`
    Resources   []*Resource   `json:"resources"`
}

// LessonType 章节类型
type LessonType string

const (
    LessonTypeVideo     LessonType = "video"
    LessonTypeText      LessonType = "text"
    LessonTypeAudio     LessonType = "audio"
    LessonTypeInteractive LessonType = "interactive"
    LessonTypeQuiz      LessonType = "quiz"
)

// Content 内容模型
type Content struct {
    ID       string                 `json:"id"`
    Type     ContentType            `json:"type"`
    URL      string                 `json:"url"`
    Data     []byte                 `json:"data"`
    Metadata map[string]interface{} `json:"metadata"`
}

// ContentType 内容类型
type ContentType string

const (
    ContentTypeVideo     ContentType = "video"
    ContentTypeDocument  ContentType = "document"
    ContentTypeImage     ContentType = "image"
    ContentTypeAudio     ContentType = "audio"
    ContentTypeHTML      ContentType = "html"
    ContentTypeMarkdown  ContentType = "markdown"
)

// Resource 资源模型
type Resource struct {
    ID          string      `json:"id"`
    Title       string      `json:"title"`
    Description string      `json:"description"`
    Type        ResourceType `json:"type"`
    URL         string      `json:"url"`
    Size        int64       `json:"size"`
    Format      string      `json:"format"`
}

// ResourceType 资源类型
type ResourceType string

const (
    ResourceTypeDocument ResourceType = "document"
    ResourceTypeVideo    ResourceType = "video"
    ResourceTypeAudio    ResourceType = "audio"
    ResourceTypeImage    ResourceType = "image"
    ResourceTypeLink     ResourceType = "link"
)

// CreateCourse 创建课程
func (cs *CourseService) CreateCourse(ctx context.Context, req *CreateCourseRequest) (*Course, error) {
    // 验证请求
    if err := cs.validateCreateCourseRequest(req); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }
    
    // 创建课程
    course := &Course{
        ID:              uuid.New().String(),
        Title:           req.Title,
        Description:     req.Description,
        InstructorID:    req.InstructorID,
        Category:        req.Category,
        Level:           req.Level,
        Duration:        req.Duration,
        Price:           req.Price,
        Prerequisites:   req.Prerequisites,
        LearningOutcomes: req.LearningOutcomes,
        Modules:         req.Modules,
        Status:          CourseStatusDraft,
        EnrollmentCount: 0,
        Rating:          0.0,
        ReviewCount:     0,
        CreatedAt:       time.Now(),
        UpdatedAt:       time.Now(),
        Metadata:        make(map[string]interface{}),
    }
    
    // 保存课程
    if err := cs.repo.Create(ctx, course); err != nil {
        return nil, fmt.Errorf("failed to create course: %w", err)
    }
    
    return course, nil
}

// validateCreateCourseRequest 验证创建课程请求
func (cs *CourseService) validateCreateCourseRequest(req *CreateCourseRequest) error {
    if req.Title == "" {
        return errors.New("title is required")
    }
    
    if req.Description == "" {
        return errors.New("description is required")
    }
    
    if req.InstructorID == "" {
        return errors.New("instructor ID is required")
    }
    
    if req.Category == "" {
        return errors.New("category is required")
    }
    
    if req.Level == "" {
        return errors.New("level is required")
    }
    
    return nil
}

// CreateCourseRequest 创建课程请求
type CreateCourseRequest struct {
    Title            string      `json:"title"`
    Description      string      `json:"description"`
    InstructorID     string      `json:"instructor_id"`
    Category         string      `json:"category"`
    Level            CourseLevel `json:"level"`
    Duration         time.Duration `json:"duration"`
    Price            *Price      `json:"price,omitempty"`
    Prerequisites    []string    `json:"prerequisites,omitempty"`
    LearningOutcomes []string    `json:"learning_outcomes,omitempty"`
    Modules          []*Module   `json:"modules,omitempty"`
}
```

#### 2.2.3 学习引擎

```go
// LearningEngine 学习引擎
type LearningEngine struct {
    progress    ProgressService
    analytics   AnalyticsService
    recommendation RecommendationService
    adaptive    AdaptiveLearningService
    config      *LearningEngineConfig
}

// LearningEngineConfig 学习引擎配置
type LearningEngineConfig struct {
    EnableAdaptiveLearning bool    `json:"enable_adaptive_learning"`
    ProgressThreshold      float64 `json:"progress_threshold"`
    RecommendationCount    int     `json:"recommendation_count"`
    AnalyticsEnabled       bool    `json:"analytics_enabled"`
}

// LearningSession 学习会话
type LearningSession struct {
    ID          string                 `json:"id"`
    UserID      string                 `json:"user_id"`
    CourseID    string                 `json:"course_id"`
    ModuleID    string                 `json:"module_id"`
    LessonID    string                 `json:"lesson_id"`
    StartTime   time.Time              `json:"start_time"`
    EndTime     *time.Time             `json:"end_time,omitempty"`
    Duration    time.Duration          `json:"duration"`
    Progress    float64                `json:"progress"`
    Score       *float64               `json:"score,omitempty"`
    Events      []*LearningEvent       `json:"events"`
    Metadata    map[string]interface{} `json:"metadata"`
}

// LearningEvent 学习事件
type LearningEvent struct {
    ID        string           `json:"id"`
    Type      LearningEventType `json:"type"`
    Timestamp time.Time        `json:"timestamp"`
    Data      interface{}      `json:"data"`
}

// LearningEventType 学习事件类型
type LearningEventType string

const (
    LearningEventTypeStart     LearningEventType = "start"
    LearningEventTypePause     LearningEventType = "pause"
    LearningEventTypeResume    LearningEventType = "resume"
    LearningEventTypeComplete  LearningEventType = "complete"
    LearningEventTypeQuiz      LearningEventType = "quiz"
    LearningEventTypeAssignment LearningEventType = "assignment"
    LearningEventTypeDiscussion LearningEventType = "discussion"
)

// StartSession 开始学习会话
func (le *LearningEngine) StartSession(ctx context.Context, req *StartSessionRequest) (*LearningSession, error) {
    // 创建会话
    session := &LearningSession{
        ID:        uuid.New().String(),
        UserID:    req.UserID,
        CourseID:  req.CourseID,
        ModuleID:  req.ModuleID,
        LessonID:  req.LessonID,
        StartTime: time.Now(),
        Progress:  0.0,
        Events:    make([]*LearningEvent, 0),
        Metadata:  make(map[string]interface{}),
    }
    
    // 记录开始事件
    session.Events = append(session.Events, &LearningEvent{
        ID:        uuid.New().String(),
        Type:      LearningEventTypeStart,
        Timestamp: session.StartTime,
        Data:      req,
    })
    
    // 保存会话
    if err := le.progress.SaveSession(ctx, session); err != nil {
        return nil, fmt.Errorf("failed to save session: %w", err)
    }
    
    // 更新学习分析
    if le.config.AnalyticsEnabled {
        go le.analytics.RecordSessionStart(session)
    }
    
    return session, nil
}

// UpdateProgress 更新学习进度
func (le *LearningEngine) UpdateProgress(ctx context.Context, sessionID string, progress float64) error {
    // 获取会话
    session, err := le.progress.GetSession(ctx, sessionID)
    if err != nil {
        return fmt.Errorf("failed to get session: %w", err)
    }
    
    // 更新进度
    session.Progress = progress
    session.Duration = time.Since(session.StartTime)
    
    // 记录进度事件
    session.Events = append(session.Events, &LearningEvent{
        ID:        uuid.New().String(),
        Type:      LearningEventTypeComplete,
        Timestamp: time.Now(),
        Data:      map[string]interface{}{"progress": progress},
    })
    
    // 保存更新
    if err := le.progress.UpdateSession(ctx, session); err != nil {
        return fmt.Errorf("failed to update session: %w", err)
    }
    
    // 检查是否完成
    if progress >= 100.0 {
        return le.completeSession(ctx, session)
    }
    
    return nil
}

// completeSession 完成会话
func (le *LearningEngine) completeSession(ctx context.Context, session *LearningSession) error {
    endTime := time.Now()
    session.EndTime = &endTime
    session.Duration = endTime.Sub(session.StartTime)
    session.Progress = 100.0
    
    // 记录完成事件
    session.Events = append(session.Events, &LearningEvent{
        ID:        uuid.New().String(),
        Type:      LearningEventTypeComplete,
        Timestamp: endTime,
        Data:      map[string]interface{}{"completed": true},
    })
    
    // 保存完成状态
    if err := le.progress.UpdateSession(ctx, session); err != nil {
        return fmt.Errorf("failed to update session: %w", err)
    }
    
    // 更新学习分析
    if le.config.AnalyticsEnabled {
        go le.analytics.RecordSessionComplete(session)
    }
    
    // 生成推荐
    if le.config.EnableAdaptiveLearning {
        go le.generateRecommendations(session.UserID, session.CourseID)
    }
    
    return nil
}

// generateRecommendations 生成推荐
func (le *LearningEngine) generateRecommendations(userID, courseID string) {
    // 获取用户学习历史
    history, err := le.progress.GetUserHistory(context.Background(), userID)
    if err != nil {
        log.Printf("Failed to get user history: %v", err)
        return
    }
    
    // 分析学习模式
    patterns := le.analytics.AnalyzeLearningPatterns(history)
    
    // 生成推荐
    recommendations := le.recommendation.GenerateRecommendations(userID, courseID, patterns)
    
    // 保存推荐
    if err := le.recommendation.SaveRecommendations(userID, recommendations); err != nil {
        log.Printf("Failed to save recommendations: %v", err)
    }
}

// StartSessionRequest 开始会话请求
type StartSessionRequest struct {
    UserID   string `json:"user_id"`
    CourseID string `json:"course_id"`
    ModuleID string `json:"module_id"`
    LessonID string `json:"lesson_id"`
}
```

#### 2.2.4 推荐系统

```go
// RecommendationService 推荐服务
type RecommendationService struct {
    collaborative CollaborativeFilter
    content      ContentBasedFilter
    hybrid       HybridFilter
    config       *RecommendationConfig
}

// RecommendationConfig 推荐配置
type RecommendationConfig struct {
    Algorithm        string  `json:"algorithm"`
    MinScore         float64 `json:"min_score"`
    MaxRecommendations int   `json:"max_recommendations"`
    EnablePersonalization bool `json:"enable_personalization"`
}

// Recommendation 推荐结果
type Recommendation struct {
    ID          string                 `json:"id"`
    UserID      string                 `json:"user_id"`
    CourseID    string                 `json:"course_id"`
    Score       float64                `json:"score"`
    Type        RecommendationType     `json:"type"`
    Reason      string                 `json:"reason"`
    CreatedAt   time.Time              `json:"created_at"`
    Metadata    map[string]interface{} `json:"metadata"`
}

// RecommendationType 推荐类型
type RecommendationType string

const (
    RecommendationTypeCollaborative RecommendationType = "collaborative"
    RecommendationTypeContent       RecommendationType = "content"
    RecommendationTypeHybrid        RecommendationType = "hybrid"
    RecommendationTypePopular       RecommendationType = "popular"
    RecommendationTypeTrending      RecommendationType = "trending"
)

// GenerateRecommendations 生成推荐
func (rs *RecommendationService) GenerateRecommendations(userID, courseID string, patterns *LearningPatterns) ([]*Recommendation, error) {
    var recommendations []*Recommendation
    
    switch rs.config.Algorithm {
    case "collaborative":
        recs, err := rs.collaborative.Recommend(userID, rs.config.MaxRecommendations)
        if err != nil {
            return nil, err
        }
        recommendations = append(recommendations, recs...)
        
    case "content":
        recs, err := rs.content.Recommend(userID, patterns, rs.config.MaxRecommendations)
        if err != nil {
            return nil, err
        }
        recommendations = append(recommendations, recs...)
        
    case "hybrid":
        recs, err := rs.hybrid.Recommend(userID, patterns, rs.config.MaxRecommendations)
        if err != nil {
            return nil, err
        }
        recommendations = append(recommendations, recs...)
        
    default:
        return nil, fmt.Errorf("unknown algorithm: %s", rs.config.Algorithm)
    }
    
    // 过滤低分推荐
    filtered := rs.filterRecommendations(recommendations)
    
    // 个性化调整
    if rs.config.EnablePersonalization {
        filtered = rs.personalizeRecommendations(filtered, patterns)
    }
    
    return filtered, nil
}

// filterRecommendations 过滤推荐
func (rs *RecommendationService) filterRecommendations(recommendations []*Recommendation) []*Recommendation {
    var filtered []*Recommendation
    
    for _, rec := range recommendations {
        if rec.Score >= rs.config.MinScore {
            filtered = append(filtered, rec)
        }
    }
    
    return filtered
}

// personalizeRecommendations 个性化推荐
func (rs *RecommendationService) personalizeRecommendations(recommendations []*Recommendation, patterns *LearningPatterns) []*Recommendation {
    for _, rec := range recommendations {
        // 基于学习模式调整分数
        if patterns.PreferredTime != "" {
            rec.Score *= rs.getTimePreferenceBonus(patterns.PreferredTime)
        }
        
        if patterns.PreferredDuration > 0 {
            rec.Score *= rs.getDurationPreferenceBonus(patterns.PreferredDuration)
        }
        
        if len(patterns.PreferredSubjects) > 0 {
            rec.Score *= rs.getSubjectPreferenceBonus(rec.CourseID, patterns.PreferredSubjects)
        }
    }
    
    // 重新排序
    sort.Slice(recommendations, func(i, j int) bool {
        return recommendations[i].Score > recommendations[j].Score
    })
    
    return recommendations
}

// getTimePreferenceBonus 获取时间偏好奖励
func (rs *RecommendationService) getTimePreferenceBonus(preferredTime string) float64 {
    // 简化实现
    return 1.1
}

// getDurationPreferenceBonus 获取时长偏好奖励
func (rs *RecommendationService) getDurationPreferenceBonus(preferredDuration time.Duration) float64 {
    // 简化实现
    return 1.05
}

// getSubjectPreferenceBonus 获取科目偏好奖励
func (rs *RecommendationService) getSubjectPreferenceBonus(courseID string, preferredSubjects []string) float64 {
    // 简化实现
    return 1.2
}

// CollaborativeFilter 协同过滤接口
type CollaborativeFilter interface {
    Recommend(userID string, count int) ([]*Recommendation, error)
}

// ContentBasedFilter 基于内容的过滤接口
type ContentBasedFilter interface {
    Recommend(userID string, patterns *LearningPatterns, count int) ([]*Recommendation, error)
}

// HybridFilter 混合过滤接口
type HybridFilter interface {
    Recommend(userID string, patterns *LearningPatterns, count int) ([]*Recommendation, error)
}

// LearningPatterns 学习模式
type LearningPatterns struct {
    PreferredTime     string        `json:"preferred_time"`
    PreferredDuration time.Duration `json:"preferred_duration"`
    PreferredSubjects []string      `json:"preferred_subjects"`
    LearningStyle     string        `json:"learning_style"`
    Difficulty        string        `json:"difficulty"`
}
```

## 3. 数学建模

### 3.1 推荐算法模型

**协同过滤评分**：
$$r_{u,i} = \frac{\sum_{v \in N(u)} sim(u,v) \cdot r_{v,i}}{\sum_{v \in N(u)} |sim(u,v)|}$$

其中：
- $r_{u,i}$ 为用户 $u$ 对项目 $i$ 的预测评分
- $N(u)$ 为用户 $u$ 的邻居集合
- $sim(u,v)$ 为用户 $u$ 和 $v$ 的相似度

**余弦相似度**：
$$sim(u,v) = \frac{\sum_{i} r_{u,i} \cdot r_{v,i}}{\sqrt{\sum_{i} r_{u,i}^2} \cdot \sqrt{\sum_{i} r_{v,i}^2}}$$

### 3.2 学习进度模型

**进度计算**：
$$P(u,c) = \frac{\sum_{i=1}^{n} w_i \cdot p_i}{\sum_{i=1}^{n} w_i}$$

其中：
- $P(u,c)$ 为用户 $u$ 在课程 $c$ 中的总进度
- $w_i$ 为模块 $i$ 的权重
- $p_i$ 为模块 $i$ 的完成进度

### 3.3 学习效果评估

**学习效果评分**：
$$E(u,c) = \alpha \cdot P(u,c) + \beta \cdot Q(u,c) + \gamma \cdot T(u,c)$$

其中：
- $P(u,c)$ 为进度分数
- $Q(u,c)$ 为质量分数
- $T(u,c)$ 为时间效率分数
- $\alpha + \beta + \gamma = 1$

## 4. 性能优化

### 4.1 缓存策略

```go
// LearningCache 学习缓存
type LearningCache struct {
    userCache     *lru.Cache
    courseCache   *lru.Cache
    progressCache *lru.Cache
    ttl           time.Duration
}

// CacheEntry 缓存条目
type CacheEntry struct {
    Data      interface{}
    ExpiresAt time.Time
}

// GetUserProgress 获取用户进度
func (lc *LearningCache) GetUserProgress(userID, courseID string) (*UserProgress, bool) {
    key := fmt.Sprintf("progress:%s:%s", userID, courseID)
    
    if entry, found := lc.progressCache.Get(key); found {
        cacheEntry := entry.(*CacheEntry)
        if time.Now().Before(cacheEntry.ExpiresAt) {
            return cacheEntry.Data.(*UserProgress), true
        } else {
            lc.progressCache.Remove(key)
        }
    }
    
    return nil, false
}

// SetUserProgress 设置用户进度
func (lc *LearningCache) SetUserProgress(userID, courseID string, progress *UserProgress) {
    key := fmt.Sprintf("progress:%s:%s", userID, courseID)
    
    entry := &CacheEntry{
        Data:      progress,
        ExpiresAt: time.Now().Add(lc.ttl),
    }
    
    lc.progressCache.Add(key, entry)
}
```

### 4.2 异步处理

```go
// AsyncProcessor 异步处理器
type AsyncProcessor struct {
    eventQueue   chan *LearningEvent
    workers      int
    ctx          context.Context
    cancel       context.CancelFunc
    wg           sync.WaitGroup
}

// NewAsyncProcessor 创建异步处理器
func NewAsyncProcessor(workers int) *AsyncProcessor {
    return &AsyncProcessor{
        eventQueue: make(chan *LearningEvent, 1000),
        workers:    workers,
    }
}

// Start 启动处理器
func (ap *AsyncProcessor) Start() {
    ap.ctx, ap.cancel = context.WithCancel(context.Background())
    
    for i := 0; i < ap.workers; i++ {
        ap.wg.Add(1)
        go ap.worker()
    }
}

// Stop 停止处理器
func (ap *AsyncProcessor) Stop() {
    ap.cancel()
    ap.wg.Wait()
}

// worker 工作协程
func (ap *AsyncProcessor) worker() {
    defer ap.wg.Done()
    
    for {
        select {
        case <-ap.ctx.Done():
            return
        case event := <-ap.eventQueue:
            ap.processEvent(event)
        }
    }
}

// processEvent 处理事件
func (ap *AsyncProcessor) processEvent(event *LearningEvent) {
    // 处理学习事件
    switch event.Type {
    case LearningEventTypeComplete:
        ap.handleCompletion(event)
    case LearningEventTypeQuiz:
        ap.handleQuiz(event)
    case LearningEventTypeAssignment:
        ap.handleAssignment(event)
    }
}

// handleCompletion 处理完成事件
func (ap *AsyncProcessor) handleCompletion(event *LearningEvent) {
    // 更新学习进度
    // 生成推荐
    // 发送通知
}

// handleQuiz 处理测验事件
func (ap *AsyncProcessor) handleQuiz(event *LearningEvent) {
    // 评分
    // 更新进度
    // 生成反馈
}

// handleAssignment 处理作业事件
func (ap *AsyncProcessor) handleAssignment(event *LearningEvent) {
    // 提交处理
    // 评分
    // 反馈生成
}
```

## 5. 监控与分析

### 5.1 学习分析

```go
// LearningAnalytics 学习分析
type LearningAnalytics struct {
    metrics    *LearningMetrics
    insights   *LearningInsights
    reports    *LearningReports
    config     *AnalyticsConfig
}

// LearningMetrics 学习指标
type LearningMetrics struct {
    totalUsers       int64
    activeUsers      int64
    courseCompletions int64
    avgCompletionTime time.Duration
    avgScore         float64
    retentionRate    float64
    mu               sync.RWMutex
}

// RecordUserActivity 记录用户活动
func (la *LearningAnalytics) RecordUserActivity(userID string, activity *UserActivity) {
    la.metrics.mu.Lock()
    defer la.metrics.mu.Unlock()
    
    // 更新指标
    la.metrics.totalUsers++
    
    if activity.IsActive {
        la.metrics.activeUsers++
    }
    
    if activity.CourseCompleted {
        la.metrics.courseCompletions++
    }
    
    // 更新平均完成时间
    if activity.CompletionTime > 0 {
        total := la.metrics.avgCompletionTime * time.Duration(la.metrics.courseCompletions-1)
        la.metrics.avgCompletionTime = (total + activity.CompletionTime) / time.Duration(la.metrics.courseCompletions)
    }
    
    // 更新平均分数
    if activity.Score > 0 {
        total := la.metrics.avgScore * float64(la.metrics.courseCompletions-1)
        la.metrics.avgScore = (total + activity.Score) / float64(la.metrics.courseCompletions)
    }
}

// GetMetrics 获取指标
func (la *LearningAnalytics) GetMetrics() map[string]interface{} {
    la.metrics.mu.RLock()
    defer la.metrics.mu.RUnlock()
    
    return map[string]interface{}{
        "total_users":         la.metrics.totalUsers,
        "active_users":        la.metrics.activeUsers,
        "course_completions":  la.metrics.courseCompletions,
        "avg_completion_time": la.metrics.avgCompletionTime,
        "avg_score":           la.metrics.avgScore,
        "retention_rate":      la.metrics.retentionRate,
    }
}

// UserActivity 用户活动
type UserActivity struct {
    UserID            string        `json:"user_id"`
    IsActive          bool          `json:"is_active"`
    CourseCompleted   bool          `json:"course_completed"`
    CompletionTime    time.Duration `json:"completion_time"`
    Score             float64       `json:"score"`
    SessionDuration   time.Duration `json:"session_duration"`
    LoginCount        int           `json:"login_count"`
}
```

## 6. 总结

在线学习平台是教育科技领域的核心系统，通过个性化学习、智能推荐和实时分析，为学习者提供优质的数字化学习体验。本模块提供了：

1. **完整的平台架构**：用户服务、课程服务、学习引擎、推荐系统
2. **个性化学习**：自适应学习路径、智能推荐算法、学习分析
3. **高性能实现**：缓存策略、异步处理、并发控制
4. **数据分析**：学习指标、用户行为分析、效果评估

通过Go语言的高性能和并发特性，实现了高效、可扩展的在线学习平台，为教育数字化转型提供了强有力的技术支撑。

---

**相关链接**：
- [02-教育管理系统](../02-Education-Management-System/README.md)
- [03-智能评估系统](../03-Intelligent-Assessment-System/README.md)
- [04-内容管理系统](../04-Content-Management-System/README.md)
