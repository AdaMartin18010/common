# 01-在线商城平台 (Online Marketplace Platform)

## 1. 概述

### 1.1 定义与目标

在线商城平台是电子商务领域的核心系统，为买家和卖家提供商品交易、支付处理、物流配送等完整的电商服务。

**形式化定义**：
设 $U$ 为用户集合，$P$ 为商品集合，$O$ 为订单集合，$T$ 为交易集合，则商城平台函数 $f$ 定义为：

$$f: U \times P \rightarrow O \times T$$

其中：

- $U = \{u_1, u_2, ..., u_n\}$ 为用户集合（买家+卖家）
- $P = \{p_1, p_2, ..., p_m\}$ 为商品集合
- $O = \{o_1, o_2, ..., o_k\}$ 为订单集合
- $T = \{t_1, t_2, ..., t_r\}$ 为交易集合

### 1.2 业务模型

**商品推荐模型**：
$$\text{RecommendationScore}(u, p) = \alpha \cdot \text{Popularity}(p) + \beta \cdot \text{Relevance}(u, p) + \gamma \cdot \text{Profit}(p)$$

其中 $\alpha + \beta + \gamma = 1$ 为权重系数。

**定价优化模型**：
$$\text{OptimalPrice}(p) = \arg\max_{price} \text{Revenue}(price) \times \text{Demand}(price)$$

## 2. 架构设计

### 2.1 微服务架构

```text
┌─────────────────────────────────────┐
│           用户服务层                  │
├─────────────────────────────────────┤
│           商品服务层                  │
├─────────────────────────────────────┤
│           订单服务层                  │
├─────────────────────────────────────┤
│           支付服务层                  │
├─────────────────────────────────────┤
│           库存服务层                  │
├─────────────────────────────────────┤
│           推荐服务层                  │
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
    address     AddressService
    payment     PaymentMethodService
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
    Addresses   []*Address             `json:"addresses"`
    PaymentMethods []*PaymentMethod    `json:"payment_methods"`
    Status      UserStatus             `json:"status"`
    CreatedAt   time.Time              `json:"created_at"`
    UpdatedAt   time.Time              `json:"updated_at"`
    Metadata    map[string]interface{} `json:"metadata"`
}

// UserRole 用户角色
type UserRole string

const (
    UserRoleBuyer  UserRole = "buyer"
    UserRoleSeller UserRole = "seller"
    UserRoleAdmin  UserRole = "admin"
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
    FirstName   string     `json:"first_name"`
    LastName    string     `json:"last_name"`
    Phone       string     `json:"phone"`
    DateOfBirth *time.Time `json:"date_of_birth"`
    AvatarURL   string     `json:"avatar_url"`
    Language    string     `json:"language"`
    Currency    string     `json:"currency"`
    Preferences *Preferences `json:"preferences"`
}

// Preferences 用户偏好
type Preferences struct {
    Theme       string `json:"theme"`
    Notifications NotificationSettings `json:"notifications"`
    Privacy     PrivacySettings        `json:"privacy"`
}

// Address 地址模型
type Address struct {
    ID          string `json:"id"`
    Type        AddressType `json:"type"`
    FirstName   string `json:"first_name"`
    LastName    string `json:"last_name"`
    Company     string `json:"company"`
    Address1    string `json:"address1"`
    Address2    string `json:"address2"`
    City        string `json:"city"`
    State       string `json:"state"`
    PostalCode  string `json:"postal_code"`
    Country     string `json:"country"`
    Phone       string `json:"phone"`
    IsDefault   bool   `json:"is_default"`
}

// AddressType 地址类型
type AddressType string

const (
    AddressTypeBilling  AddressType = "billing"
    AddressTypeShipping AddressType = "shipping"
)

// PaymentMethod 支付方式
type PaymentMethod struct {
    ID          string           `json:"id"`
    Type        PaymentType      `json:"type"`
    Provider    string           `json:"provider"`
    LastFour    string           `json:"last_four"`
    ExpiryMonth int              `json:"expiry_month"`
    ExpiryYear  int              `json:"expiry_year"`
    IsDefault   bool             `json:"is_default"`
    Metadata    map[string]interface{} `json:"metadata"`
}

// PaymentType 支付类型
type PaymentType string

const (
    PaymentTypeCreditCard PaymentType = "credit_card"
    PaymentTypeDebitCard  PaymentType = "debit_card"
    PaymentTypePayPal     PaymentType = "paypal"
    PaymentTypeBankTransfer PaymentType = "bank_transfer"
)

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
    
    // 设置地址
    if len(req.Addresses) > 0 {
        user.Addresses = req.Addresses
    }
    
    // 设置支付方式
    if len(req.PaymentMethods) > 0 {
        user.PaymentMethods = req.PaymentMethods
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

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
    Email         string           `json:"email"`
    Username      string           `json:"username"`
    Role          UserRole         `json:"role"`
    Profile       *UserProfile     `json:"profile,omitempty"`
    Addresses     []*Address       `json:"addresses,omitempty"`
    PaymentMethods []*PaymentMethod `json:"payment_methods,omitempty"`
}

// sendWelcomeEmail 发送欢迎邮件
func (us *UserService) sendWelcomeEmail(user *User) {
    // 实现邮件发送逻辑
    log.Printf("Sending welcome email to %s", user.Email)
}
```

#### 2.2.2 商品服务

```go
// ProductService 商品服务
type ProductService struct {
    repo        ProductRepository
    category    CategoryService
    inventory   InventoryService
    search      SearchService
    config      *ProductServiceConfig
}

// ProductServiceConfig 商品服务配置
type ProductServiceConfig struct {
    MaxProductsPerPage int     `json:"max_products_per_page"`
    EnableSearch       bool    `json:"enable_search"`
    EnableRecommendation bool  `json:"enable_recommendation"`
    CacheTTL           time.Duration `json:"cache_ttl"`
}

// Product 商品模型
type Product struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    CategoryID  string                 `json:"category_id"`
    Brand       string                 `json:"brand"`
    SKU         string                 `json:"sku"`
    Price       *Money                 `json:"price"`
    SalePrice   *Money                 `json:"sale_price,omitempty"`
    Inventory   *InventoryInfo         `json:"inventory"`
    Images      []*ProductImage        `json:"images"`
    Attributes  map[string]interface{} `json:"attributes"`
    Variants    []*ProductVariant      `json:"variants"`
    Status      ProductStatus          `json:"status"`
    SellerID    string                 `json:"seller_id"`
    Rating      float64                `json:"rating"`
    ReviewCount int                    `json:"review_count"`
    ViewCount   int64                  `json:"view_count"`
    CreatedAt   time.Time              `json:"created_at"`
    UpdatedAt   time.Time              `json:"updated_at"`
    Metadata    map[string]interface{} `json:"metadata"`
}

// Money 金额模型
type Money struct {
    Amount   int64  `json:"amount"`   // 以分为单位
    Currency string `json:"currency"`
}

// ProductStatus 商品状态
type ProductStatus string

const (
    ProductStatusDraft     ProductStatus = "draft"
    ProductStatusActive    ProductStatus = "active"
    ProductStatusInactive  ProductStatus = "inactive"
    ProductStatusOutOfStock ProductStatus = "out_of_stock"
    ProductStatusDiscontinued ProductStatus = "discontinued"
)

// InventoryInfo 库存信息
type InventoryInfo struct {
    Quantity     int    `json:"quantity"`
    Reserved     int    `json:"reserved"`
    Available    int    `json:"available"`
    LowStockThreshold int `json:"low_stock_threshold"`
    TrackQuantity bool  `json:"track_quantity"`
}

// ProductImage 商品图片
type ProductImage struct {
    ID       string `json:"id"`
    URL      string `json:"url"`
    AltText  string `json:"alt_text"`
    Order    int    `json:"order"`
    IsPrimary bool  `json:"is_primary"`
}

// ProductVariant 商品变体
type ProductVariant struct {
    ID          string                 `json:"id"`
    SKU         string                 `json:"sku"`
    Price       *Money                 `json:"price"`
    SalePrice   *Money                 `json:"sale_price,omitempty"`
    Inventory   *InventoryInfo         `json:"inventory"`
    Attributes  map[string]interface{} `json:"attributes"`
    Images      []*ProductImage        `json:"images"`
}

// CreateProduct 创建商品
func (ps *ProductService) CreateProduct(ctx context.Context, req *CreateProductRequest) (*Product, error) {
    // 验证请求
    if err := ps.validateCreateProductRequest(req); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }
    
    // 创建商品
    product := &Product{
        ID:          uuid.New().String(),
        Name:        req.Name,
        Description: req.Description,
        CategoryID:  req.CategoryID,
        Brand:       req.Brand,
        SKU:         req.SKU,
        Price:       req.Price,
        SalePrice:   req.SalePrice,
        Inventory:   req.Inventory,
        Images:      req.Images,
        Attributes:  req.Attributes,
        Variants:    req.Variants,
        Status:      ProductStatusDraft,
        SellerID:    req.SellerID,
        Rating:      0.0,
        ReviewCount: 0,
        ViewCount:   0,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
        Metadata:    make(map[string]interface{}),
    }
    
    // 保存商品
    if err := ps.repo.Create(ctx, product); err != nil {
        return nil, fmt.Errorf("failed to create product: %w", err)
    }
    
    // 更新搜索索引
    if ps.config.EnableSearch {
        go ps.search.IndexProduct(product)
    }
    
    return product, nil
}

// validateCreateProductRequest 验证创建商品请求
func (ps *ProductService) validateCreateProductRequest(req *CreateProductRequest) error {
    if req.Name == "" {
        return errors.New("name is required")
    }
    
    if req.Description == "" {
        return errors.New("description is required")
    }
    
    if req.CategoryID == "" {
        return errors.New("category ID is required")
    }
    
    if req.Brand == "" {
        return errors.New("brand is required")
    }
    
    if req.SKU == "" {
        return errors.New("SKU is required")
    }
    
    if req.Price == nil {
        return errors.New("price is required")
    }
    
    if req.SellerID == "" {
        return errors.New("seller ID is required")
    }
    
    return nil
}

// SearchProducts 搜索商品
func (ps *ProductService) SearchProducts(ctx context.Context, query *SearchQuery) (*SearchResult, error) {
    if !ps.config.EnableSearch {
        return ps.repo.Search(ctx, query)
    }
    
    return ps.search.Search(ctx, query)
}

// SearchQuery 搜索查询
type SearchQuery struct {
    Query       string                 `json:"query"`
    CategoryID  string                 `json:"category_id,omitempty"`
    Brand       string                 `json:"brand,omitempty"`
    MinPrice    *Money                 `json:"min_price,omitempty"`
    MaxPrice    *Money                 `json:"max_price,omitempty"`
    Attributes  map[string]interface{} `json:"attributes,omitempty"`
    SortBy      string                 `json:"sort_by,omitempty"`
    SortOrder   string                 `json:"sort_order,omitempty"`
    Page        int                    `json:"page"`
    PageSize    int                    `json:"page_size"`
}

// SearchResult 搜索结果
type SearchResult struct {
    Products    []*Product `json:"products"`
    Total       int64      `json:"total"`
    Page        int        `json:"page"`
    PageSize    int        `json:"page_size"`
    Facets      *Facets    `json:"facets,omitempty"`
}

// Facets 分面信息
type Facets struct {
    Categories  []*FacetValue `json:"categories"`
    Brands      []*FacetValue `json:"brands"`
    PriceRanges []*FacetValue `json:"price_ranges"`
    Attributes  map[string][]*FacetValue `json:"attributes"`
}

// FacetValue 分面值
type FacetValue struct {
    Value string `json:"value"`
    Count int64  `json:"count"`
}

// CreateProductRequest 创建商品请求
type CreateProductRequest struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    CategoryID  string                 `json:"category_id"`
    Brand       string                 `json:"brand"`
    SKU         string                 `json:"sku"`
    Price       *Money                 `json:"price"`
    SalePrice   *Money                 `json:"sale_price,omitempty"`
    Inventory   *InventoryInfo         `json:"inventory"`
    Images      []*ProductImage        `json:"images,omitempty"`
    Attributes  map[string]interface{} `json:"attributes,omitempty"`
    Variants    []*ProductVariant      `json:"variants,omitempty"`
    SellerID    string                 `json:"seller_id"`
}
```

#### 2.2.3 订单服务

```go
// OrderService 订单服务
type OrderService struct {
    repo        OrderRepository
    product     ProductService
    payment     PaymentService
    inventory   InventoryService
    shipping    ShippingService
    config      *OrderServiceConfig
}

// OrderServiceConfig 订单服务配置
type OrderServiceConfig struct {
    MaxOrdersPerPage int           `json:"max_orders_per_page"`
    OrderTimeout     time.Duration `json:"order_timeout"`
    EnableAutoCancel bool          `json:"enable_auto_cancel"`
    CancelThreshold  time.Duration `json:"cancel_threshold"`
}

// Order 订单模型
type Order struct {
    ID              string                 `json:"id"`
    UserID          string                 `json:"user_id"`
    OrderNumber     string                 `json:"order_number"`
    Status          OrderStatus            `json:"status"`
    Items           []*OrderItem           `json:"items"`
    Subtotal        *Money                 `json:"subtotal"`
    Tax             *Money                 `json:"tax"`
    Shipping        *Money                 `json:"shipping"`
    Discount        *Money                 `json:"discount"`
    Total           *Money                 `json:"total"`
    BillingAddress  *Address               `json:"billing_address"`
    ShippingAddress *Address               `json:"shipping_address"`
    PaymentMethod   *PaymentMethod         `json:"payment_method"`
    ShippingMethod  *ShippingMethod        `json:"shipping_method"`
    Notes           string                 `json:"notes"`
    CreatedAt       time.Time              `json:"created_at"`
    UpdatedAt       time.Time              `json:"updated_at"`
    Metadata        map[string]interface{} `json:"metadata"`
}

// OrderStatus 订单状态
type OrderStatus string

const (
    OrderStatusPending   OrderStatus = "pending"
    OrderStatusConfirmed OrderStatus = "confirmed"
    OrderStatusPaid      OrderStatus = "paid"
    OrderStatusProcessing OrderStatus = "processing"
    OrderStatusShipped   OrderStatus = "shipped"
    OrderStatusDelivered OrderStatus = "delivered"
    OrderStatusCancelled OrderStatus = "cancelled"
    OrderStatusRefunded  OrderStatus = "refunded"
)

// OrderItem 订单项
type OrderItem struct {
    ID          string     `json:"id"`
    ProductID   string     `json:"product_id"`
    VariantID   string     `json:"variant_id,omitempty"`
    Name        string     `json:"name"`
    SKU         string     `json:"sku"`
    Quantity    int        `json:"quantity"`
    UnitPrice   *Money     `json:"unit_price"`
    TotalPrice  *Money     `json:"total_price"`
    Attributes  map[string]interface{} `json:"attributes"`
}

// ShippingMethod 配送方式
type ShippingMethod struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Carrier     string `json:"carrier"`
    Service     string `json:"service"`
    Cost        *Money `json:"cost"`
    EstimatedDays int  `json:"estimated_days"`
}

// CreateOrder 创建订单
func (os *OrderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*Order, error) {
    // 验证请求
    if err := os.validateCreateOrderRequest(req); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }
    
    // 检查库存
    if err := os.checkInventory(ctx, req.Items); err != nil {
        return nil, fmt.Errorf("inventory check failed: %w", err)
    }
    
    // 计算订单金额
    calculation, err := os.calculateOrder(req)
    if err != nil {
        return nil, fmt.Errorf("calculation failed: %w", err)
    }
    
    // 创建订单
    order := &Order{
        ID:              uuid.New().String(),
        UserID:          req.UserID,
        OrderNumber:     os.generateOrderNumber(),
        Status:          OrderStatusPending,
        Items:           req.Items,
        Subtotal:        calculation.Subtotal,
        Tax:             calculation.Tax,
        Shipping:        calculation.Shipping,
        Discount:        calculation.Discount,
        Total:           calculation.Total,
        BillingAddress:  req.BillingAddress,
        ShippingAddress: req.ShippingAddress,
        PaymentMethod:   req.PaymentMethod,
        ShippingMethod:  req.ShippingMethod,
        Notes:           req.Notes,
        CreatedAt:       time.Now(),
        UpdatedAt:       time.Now(),
        Metadata:        make(map[string]interface{}),
    }
    
    // 保存订单
    if err := os.repo.Create(ctx, order); err != nil {
        return nil, fmt.Errorf("failed to create order: %w", err)
    }
    
    // 预留库存
    if err := os.reserveInventory(ctx, order); err != nil {
        return nil, fmt.Errorf("failed to reserve inventory: %w", err)
    }
    
    // 设置订单超时
    if os.config.EnableAutoCancel {
        go os.scheduleOrderCancel(order.ID)
    }
    
    return order, nil
}

// validateCreateOrderRequest 验证创建订单请求
func (os *OrderService) validateCreateOrderRequest(req *CreateOrderRequest) error {
    if req.UserID == "" {
        return errors.New("user ID is required")
    }
    
    if len(req.Items) == 0 {
        return errors.New("order items are required")
    }
    
    if req.BillingAddress == nil {
        return errors.New("billing address is required")
    }
    
    if req.ShippingAddress == nil {
        return errors.New("shipping address is required")
    }
    
    if req.PaymentMethod == nil {
        return errors.New("payment method is required")
    }
    
    if req.ShippingMethod == nil {
        return errors.New("shipping method is required")
    }
    
    return nil
}

// checkInventory 检查库存
func (os *OrderService) checkInventory(ctx context.Context, items []*OrderItem) error {
    for _, item := range items {
        available, err := os.inventory.GetAvailableQuantity(ctx, item.ProductID, item.VariantID)
        if err != nil {
            return fmt.Errorf("failed to check inventory for product %s: %w", item.ProductID, err)
        }
        
        if available < item.Quantity {
            return fmt.Errorf("insufficient inventory for product %s", item.ProductID)
        }
    }
    
    return nil
}

// calculateOrder 计算订单金额
func (os *OrderService) calculateOrder(req *CreateOrderRequest) (*OrderCalculation, error) {
    calculation := &OrderCalculation{}
    
    // 计算小计
    var subtotal int64
    for _, item := range req.Items {
        itemTotal := item.UnitPrice.Amount * int64(item.Quantity)
        item.TotalPrice = &Money{
            Amount:   itemTotal,
            Currency: item.UnitPrice.Currency,
        }
        subtotal += itemTotal
    }
    
    calculation.Subtotal = &Money{
        Amount:   subtotal,
        Currency: req.Items[0].UnitPrice.Currency,
    }
    
    // 计算税费
    taxRate := os.getTaxRate(req.ShippingAddress)
    taxAmount := int64(float64(subtotal) * taxRate)
    calculation.Tax = &Money{
        Amount:   taxAmount,
        Currency: calculation.Subtotal.Currency,
    }
    
    // 计算运费
    calculation.Shipping = req.ShippingMethod.Cost
    
    // 计算折扣
    calculation.Discount = &Money{
        Amount:   0,
        Currency: calculation.Subtotal.Currency,
    }
    
    // 计算总计
    total := subtotal + taxAmount + calculation.Shipping.Amount - calculation.Discount.Amount
    calculation.Total = &Money{
        Amount:   total,
        Currency: calculation.Subtotal.Currency,
    }
    
    return calculation, nil
}

// getTaxRate 获取税率
func (os *OrderService) getTaxRate(address *Address) float64 {
    // 简化实现，实际应该根据地址查询税率
    return 0.08 // 8% 税率
}

// reserveInventory 预留库存
func (os *OrderService) reserveInventory(ctx context.Context, order *Order) error {
    for _, item := range order.Items {
        if err := os.inventory.Reserve(ctx, item.ProductID, item.VariantID, item.Quantity); err != nil {
            return fmt.Errorf("failed to reserve inventory for product %s: %w", item.ProductID, err)
        }
    }
    
    return nil
}

// generateOrderNumber 生成订单号
func (os *OrderService) generateOrderNumber() string {
    timestamp := time.Now().Format("20060102150405")
    random := fmt.Sprintf("%04d", rand.Intn(10000))
    return fmt.Sprintf("ORD%s%s", timestamp, random)
}

// scheduleOrderCancel 调度订单取消
func (os *OrderService) scheduleOrderCancel(orderID string) {
    timer := time.NewTimer(os.config.OrderTimeout)
    <-timer.C
    
    // 检查订单状态
    ctx := context.Background()
    order, err := os.repo.GetByID(ctx, orderID)
    if err != nil {
        log.Printf("Failed to get order %s: %v", orderID, err)
        return
    }
    
    // 如果订单仍为待支付状态，则取消
    if order.Status == OrderStatusPending {
        if err := os.CancelOrder(ctx, orderID, "Order timeout"); err != nil {
            log.Printf("Failed to cancel order %s: %v", orderID, err)
        }
    }
}

// CancelOrder 取消订单
func (os *OrderService) CancelOrder(ctx context.Context, orderID, reason string) error {
    order, err := os.repo.GetByID(ctx, orderID)
    if err != nil {
        return fmt.Errorf("failed to get order: %w", err)
    }
    
    if order.Status != OrderStatusPending {
        return errors.New("order cannot be cancelled")
    }
    
    // 更新订单状态
    order.Status = OrderStatusCancelled
    order.UpdatedAt = time.Now()
    order.Metadata["cancel_reason"] = reason
    
    if err := os.repo.Update(ctx, order); err != nil {
        return fmt.Errorf("failed to update order: %w", err)
    }
    
    // 释放库存
    if err := os.releaseInventory(ctx, order); err != nil {
        log.Printf("Failed to release inventory for order %s: %v", orderID, err)
    }
    
    return nil
}

// releaseInventory 释放库存
func (os *OrderService) releaseInventory(ctx context.Context, order *Order) error {
    for _, item := range order.Items {
        if err := os.inventory.Release(ctx, item.ProductID, item.VariantID, item.Quantity); err != nil {
            return fmt.Errorf("failed to release inventory for product %s: %w", item.ProductID, err)
        }
    }
    
    return nil
}

// OrderCalculation 订单计算
type OrderCalculation struct {
    Subtotal *Money `json:"subtotal"`
    Tax      *Money `json:"tax"`
    Shipping *Money `json:"shipping"`
    Discount *Money `json:"discount"`
    Total    *Money `json:"total"`
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
    UserID          string           `json:"user_id"`
    Items           []*OrderItem     `json:"items"`
    BillingAddress  *Address         `json:"billing_address"`
    ShippingAddress *Address         `json:"shipping_address"`
    PaymentMethod   *PaymentMethod   `json:"payment_method"`
    ShippingMethod  *ShippingMethod  `json:"shipping_method"`
    Notes           string           `json:"notes,omitempty"`
}
```

#### 2.2.4 支付服务

```go
// PaymentService 支付服务
type PaymentService struct {
    providers   map[string]PaymentProvider
    processor   PaymentProcessor
    config      *PaymentServiceConfig
}

// PaymentServiceConfig 支付服务配置
type PaymentServiceConfig struct {
    DefaultProvider string          `json:"default_provider"`
    RetryAttempts   int             `json:"retry_attempts"`
    RetryDelay      time.Duration   `json:"retry_delay"`
    EnableWebhooks  bool            `json:"enable_webhooks"`
}

// PaymentProvider 支付提供商接口
type PaymentProvider interface {
    ID() string
    Name() string
    ProcessPayment(ctx context.Context, req *PaymentRequest) (*PaymentResult, error)
    RefundPayment(ctx context.Context, req *RefundRequest) (*RefundResult, error)
    ValidateWebhook(payload []byte, signature string) error
}

// PaymentRequest 支付请求
type PaymentRequest struct {
    OrderID       string        `json:"order_id"`
    Amount        *Money        `json:"amount"`
    Currency      string        `json:"currency"`
    PaymentMethod *PaymentMethod `json:"payment_method"`
    Description   string        `json:"description"`
    Metadata      map[string]interface{} `json:"metadata"`
}

// PaymentResult 支付结果
type PaymentResult struct {
    ID            string                 `json:"id"`
    Status        PaymentStatus          `json:"status"`
    Amount        *Money                 `json:"amount"`
    Currency      string                 `json:"currency"`
    ProviderID    string                 `json:"provider_id"`
    TransactionID string                 `json:"transaction_id"`
    CreatedAt     time.Time              `json:"created_at"`
    Metadata      map[string]interface{} `json:"metadata"`
}

// PaymentStatus 支付状态
type PaymentStatus string

const (
    PaymentStatusPending   PaymentStatus = "pending"
    PaymentStatusProcessing PaymentStatus = "processing"
    PaymentStatusSucceeded PaymentStatus = "succeeded"
    PaymentStatusFailed    PaymentStatus = "failed"
    PaymentStatusCancelled PaymentStatus = "cancelled"
)

// ProcessPayment 处理支付
func (ps *PaymentService) ProcessPayment(ctx context.Context, req *PaymentRequest) (*PaymentResult, error) {
    // 选择支付提供商
    provider, err := ps.selectProvider(req)
    if err != nil {
        return nil, fmt.Errorf("failed to select provider: %w", err)
    }
    
    // 处理支付
    result, err := ps.processor.ProcessWithRetry(ctx, provider, req, ps.config.RetryAttempts, ps.config.RetryDelay)
    if err != nil {
        return nil, fmt.Errorf("payment processing failed: %w", err)
    }
    
    return result, nil
}

// selectProvider 选择支付提供商
func (ps *PaymentService) selectProvider(req *PaymentRequest) (PaymentProvider, error) {
    // 根据支付方式选择提供商
    switch req.PaymentMethod.Type {
    case PaymentTypeCreditCard, PaymentTypeDebitCard:
        if provider, exists := ps.providers["stripe"]; exists {
            return provider, nil
        }
    case PaymentTypePayPal:
        if provider, exists := ps.providers["paypal"]; exists {
            return provider, nil
        }
    case PaymentTypeBankTransfer:
        if provider, exists := ps.providers["bank"]; exists {
            return provider, nil
        }
    }
    
    // 使用默认提供商
    if provider, exists := ps.providers[ps.config.DefaultProvider]; exists {
        return provider, nil
    }
    
    return nil, errors.New("no suitable payment provider found")
}

// PaymentProcessor 支付处理器
type PaymentProcessor struct{}

// ProcessWithRetry 带重试的支付处理
func (pp *PaymentProcessor) ProcessWithRetry(ctx context.Context, provider PaymentProvider, req *PaymentRequest, attempts int, delay time.Duration) (*PaymentResult, error) {
    var lastErr error
    
    for i := 0; i < attempts; i++ {
        result, err := provider.ProcessPayment(ctx, req)
        if err == nil {
            return result, nil
        }
        
        lastErr = err
        
        // 如果不是最后一次尝试，则等待后重试
        if i < attempts-1 {
            select {
            case <-ctx.Done():
                return nil, ctx.Err()
            case <-time.After(delay):
                // 继续重试
            }
        }
    }
    
    return nil, fmt.Errorf("payment failed after %d attempts: %w", attempts, lastErr)
}
```

## 3. 数学建模

### 3.1 推荐算法模型

**协同过滤评分**：
$$r_{u,i} = \frac{\sum_{v \in N(u)} sim(u,v) \cdot r_{v,i}}{\sum_{v \in N(u)} |sim(u,v)|}$$

其中：

- $r_{u,i}$ 为用户 $u$ 对商品 $i$ 的预测评分
- $N(u)$ 为用户 $u$ 的邻居集合
- $sim(u,v)$ 为用户 $u$ 和 $v$ 的相似度

**商品相似度**：
$$sim(i,j) = \frac{\sum_{u} r_{u,i} \cdot r_{u,j}}{\sqrt{\sum_{u} r_{u,i}^2} \cdot \sqrt{\sum_{u} r_{u,j}^2}}$$

### 3.2 库存优化模型

**经济订货量模型**：
$$EOQ = \sqrt{\frac{2DS}{H}}$$

其中：

- $D$ 为年需求量
- $S$ 为订货成本
- $H$ 为库存持有成本

**安全库存**：
$$SS = z \cdot \sigma_L \cdot \sqrt{L}$$

其中：

- $z$ 为服务水平系数
- $\sigma_L$ 为提前期需求标准差
- $L$ 为提前期

### 3.3 定价优化模型

**价格弹性**：
$$\epsilon = \frac{\Delta Q / Q}{\Delta P / P}$$

**最优定价**：
$$P^* = \frac{MC}{1 - \frac{1}{|\epsilon|}}$$

其中：

- $MC$ 为边际成本
- $\epsilon$ 为价格弹性

## 4. 性能优化

### 4.1 缓存策略

```go
// ProductCache 商品缓存
type ProductCache struct {
    cache *lru.Cache
    ttl   time.Duration
}

// GetProduct 获取商品
func (pc *ProductCache) GetProduct(productID string) (*Product, bool) {
    if entry, found := pc.cache.Get(productID); found {
        cacheEntry := entry.(*CacheEntry)
        if time.Now().Before(cacheEntry.ExpiresAt) {
            return cacheEntry.Data.(*Product), true
        } else {
            pc.cache.Remove(productID)
        }
    }
    
    return nil, false
}

// SetProduct 设置商品
func (pc *ProductCache) SetProduct(productID string, product *Product) {
    entry := &CacheEntry{
        Data:      product,
        ExpiresAt: time.Now().Add(pc.ttl),
    }
    
    pc.cache.Add(productID, entry)
}
```

### 4.2 异步处理

```go
// OrderProcessor 订单处理器
type OrderProcessor struct {
    orderQueue chan *Order
    workers    int
    ctx        context.Context
    cancel     context.CancelFunc
    wg         sync.WaitGroup
}

// NewOrderProcessor 创建订单处理器
func NewOrderProcessor(workers int) *OrderProcessor {
    return &OrderProcessor{
        orderQueue: make(chan *Order, 1000),
        workers:    workers,
    }
}

// Start 启动处理器
func (op *OrderProcessor) Start() {
    op.ctx, op.cancel = context.WithCancel(context.Background())
    
    for i := 0; i < op.workers; i++ {
        op.wg.Add(1)
        go op.worker()
    }
}

// Stop 停止处理器
func (op *OrderProcessor) Stop() {
    op.cancel()
    op.wg.Wait()
}

// worker 工作协程
func (op *OrderProcessor) worker() {
    defer op.wg.Done()
    
    for {
        select {
        case <-op.ctx.Done():
            return
        case order := <-op.orderQueue:
            op.processOrder(order)
        }
    }
}

// processOrder 处理订单
func (op *OrderProcessor) processOrder(order *Order) {
    // 处理订单逻辑
    // 1. 验证订单
    // 2. 处理支付
    // 3. 更新库存
    // 4. 发送通知
}
```

## 5. 监控与分析

### 5.1 业务指标

```go
// ECommerceMetrics 电商指标
type ECommerceMetrics struct {
    totalOrders      int64
    totalRevenue     int64
    avgOrderValue    float64
    conversionRate   float64
    cartAbandonment  float64
    customerLifetimeValue float64
    mu               sync.RWMutex
}

// RecordOrder 记录订单
func (em *ECommerceMetrics) RecordOrder(order *Order) {
    em.mu.Lock()
    defer em.mu.Unlock()
    
    em.totalOrders++
    em.totalRevenue += order.Total.Amount
    
    // 更新平均订单价值
    em.avgOrderValue = float64(em.totalRevenue) / float64(em.totalOrders)
}

// GetMetrics 获取指标
func (em *ECommerceMetrics) GetMetrics() map[string]interface{} {
    em.mu.RLock()
    defer em.mu.RUnlock()
    
    return map[string]interface{}{
        "total_orders":           em.totalOrders,
        "total_revenue":          em.totalRevenue,
        "avg_order_value":        em.avgOrderValue,
        "conversion_rate":        em.conversionRate,
        "cart_abandonment":       em.cartAbandonment,
        "customer_lifetime_value": em.customerLifetimeValue,
    }
}
```

## 6. 总结

在线商城平台是电子商务领域的核心系统，通过用户管理、商品管理、订单处理和支付处理等核心功能，为买家和卖家提供完整的电商服务。本模块提供了：

1. **完整的平台架构**：用户服务、商品服务、订单服务、支付服务
2. **智能推荐系统**：协同过滤、内容推荐、个性化推荐
3. **高性能实现**：缓存策略、异步处理、并发控制
4. **业务分析**：销售指标、用户行为分析、库存优化

通过Go语言的高性能和并发特性，实现了高效、可扩展的在线商城平台，为电子商务发展提供了强有力的技术支撑。

---

**相关链接**：

- [02-支付处理系统](../02-Payment-Processing-System/README.md)
- [03-库存管理系统](../03-Inventory-Management-System/README.md)
- [04-推荐引擎](../04-Recommendation-Engine/README.md)
