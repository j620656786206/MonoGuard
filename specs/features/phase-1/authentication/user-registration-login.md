# Authentication: 用戶註冊與登入功能規格

## 概述

提供安全的用戶註冊和登入功能，支援 Email/密碼認證和 OAuth 第三方登入（GitHub、GitLab）。

## 功能細節

### 1. Email/密碼註冊

#### API 端點

```
POST /api/v1/auth/register
```

#### 請求格式

```typescript
interface RegisterRequest {
  email: string;
  password: string;
  name: string;
  company?: string;
  agreeToTerms: boolean;
}
```

#### 驗證規則

- Email: 有效的 email 格式，唯一性檢查
- 密碼: 最少 8 字元，包含大小寫字母、數字和特殊字元
- 名稱: 2-50 字元
- 必須同意服務條款

#### 響應格式

```typescript
interface RegisterResponse {
  user: {
    id: string;
    email: string;
    name: string;
    createdAt: string;
  };
  tokens: {
    accessToken: string;
    refreshToken: string;
  };
  message: string;
}
```

#### 流程

```
1. 用戶提交註冊表單
2. 後端驗證輸入
3. 密碼 hash (bcrypt, cost=12)
4. 創建用戶記錄
5. 生成 JWT tokens
6. 發送歡迎郵件（異步）
7. 返回用戶資料和 tokens
```

### 2. Email/密碼登入

#### API 端點

```
POST /api/v1/auth/login
```

#### 請求格式

```typescript
interface LoginRequest {
  email: string;
  password: string;
  rememberMe?: boolean;
}
```

#### 響應格式

```typescript
interface LoginResponse {
  user: {
    id: string;
    email: string;
    name: string;
    role: string;
  };
  tokens: {
    accessToken: string;
    refreshToken: string;
  };
}
```

#### 安全措施

- 登入失敗計數 (5 次後鎖定帳號 15 分鐘)
- bcrypt 密碼驗證
- IP 地址記錄
- 異常登入偵測

### 3. Token 管理

#### Access Token

```typescript
interface AccessTokenPayload {
  userId: string;
  email: string;
  role: string;
  iat: number;
  exp: number; // 1 hour
}
```

#### Refresh Token

```typescript
interface RefreshTokenPayload {
  userId: string;
  tokenId: string;
  iat: number;
  exp: number; // 7 days or 30 days (remember me)
}
```

#### Token 刷新

```
POST /api/v1/auth/refresh
Authorization: Bearer <refresh_token>

Response:
{
  accessToken: string;
  refreshToken: string;
}
```

### 4. 登出

#### 單一裝置登出

```
POST /api/v1/auth/logout
Authorization: Bearer <access_token>

# 使當前 refresh token 失效
```

#### 所有裝置登出

```
POST /api/v1/auth/logout-all
Authorization: Bearer <access_token>

# 使所有 refresh tokens 失效
```

### 5. 前端頁面

#### 註冊頁面 (`/register`)

```tsx
interface RegisterFormProps {
  onSubmit: (data: RegisterRequest) => void;
  isLoading: boolean;
  errors?: Record<string, string>;
}

// 表單欄位
- Email (必填，email 驗證)
- 密碼 (必填，強度指示器)
- 確認密碼 (必填，匹配驗證)
- 名稱 (必填)
- 公司名稱 (選填)
- 同意條款 (checkbox，必填)

// 驗證
- 即時 email 格式驗證
- 密碼強度實時顯示
- 密碼匹配檢查
- 後端錯誤顯示

// 功能
- "已有帳號？登入" 連結
- "使用 GitHub 註冊" 按鈕
- "使用 GitLab 註冊" 按鈕
```

#### 登入頁面 (`/login`)

```tsx
interface LoginFormProps {
  onSubmit: (data: LoginRequest) => void;
  isLoading: boolean;
  errors?: Record<string, string>;
}

// 表單欄位
- Email (必填)
- 密碼 (必填)
- 記住我 (checkbox)

// 功能
- "忘記密碼？" 連結
- "建立新帳號" 連結
- "使用 GitHub 登入" 按鈕
- "使用 GitLab 登入" 按鈕
- 顯示登入失敗次數警告
```

## User Stories

### User Story 1: 用戶註冊

**As a** 新用戶
**I want to** 使用 email 和密碼註冊帳號
**So that** 我可以開始使用 MonoGuard

**Acceptance Criteria:**

- [ ] 註冊表單驗證即時反饋
- [ ] 密碼強度指示器顯示
- [ ] 註冊成功後自動登入
- [ ] 發送歡迎郵件
- [ ] 重複 email 顯示友善錯誤訊息

### User Story 2: 用戶登入

**As a** 註冊用戶
**I want to** 使用 email 和密碼登入
**So that** 我可以訪問我的專案

**Acceptance Criteria:**

- [ ] 登入成功後跳轉到儀表板
- [ ] "記住我" 功能延長登入時間
- [ ] 登入失敗顯示清楚的錯誤訊息
- [ ] 多次失敗後顯示鎖定警告
- [ ] 支援從任何頁面登入後返回原頁面

### User Story 3: 安全登出

**As a** 登入用戶
**I want to** 安全登出我的帳號
**So that** 其他人無法訪問我的資料

**Acceptance Criteria:**

- [ ] 登出後清除所有本地 tokens
- [ ] 登出後無法使用舊 token 訪問 API
- [ ] "所有裝置登出" 使所有會話失效
- [ ] 登出後跳轉到登入頁
- [ ] 確認對話框（選用）

### User Story 4: Token 自動刷新

**As a** 登入用戶
**I want to** 長時間工作不需要重新登入
**So that** 我的工作流程不被中斷

**Acceptance Criteria:**

- [ ] Access token 過期前自動刷新
- [ ] 刷新失敗時提示重新登入
- [ ] 不中斷用戶正在進行的操作
- [ ] Refresh token 過期前提前通知
- [ ] 支援多 tab 同步

## 測試項目

### 單元測試

#### 1. 註冊驗證測試

```typescript
describe('User Registration', () => {
  test('should validate email format', () => {
    const result = validateEmail('invalid-email');
    expect(result.isValid).toBe(false);
    expect(result.error).toContain('valid email');
  });

  test('should enforce password strength', () => {
    const weakPassword = validatePassword('weak');
    expect(weakPassword.isValid).toBe(false);

    const strongPassword = validatePassword('Strong@Pass123');
    expect(strongPassword.isValid).toBe(true);
  });

  test('should check email uniqueness', async () => {
    await createUser({ email: 'test@example.com' });

    const result = await register({
      email: 'test@example.com',
      password: 'Test@123',
      name: 'Test User',
      agreeToTerms: true,
    });

    expect(result.success).toBe(false);
    expect(result.error).toContain('already registered');
  });

  test('should hash password securely', async () => {
    const user = await register({
      email: 'new@example.com',
      password: 'SecurePass@123',
      name: 'New User',
      agreeToTerms: true,
    });

    const dbUser = await getUserById(user.id);
    expect(dbUser.password).not.toBe('SecurePass@123');
    expect(dbUser.password).toMatch(/^\$2[ayb]\$.{56}$/); // bcrypt hash
  });
});
```

#### 2. 登入測試

```typescript
describe('User Login', () => {
  test('should login with valid credentials', async () => {
    const user = await createTestUser();

    const result = await login({
      email: user.email,
      password: 'TestPass@123',
    });

    expect(result.success).toBe(true);
    expect(result.tokens.accessToken).toBeDefined();
    expect(result.tokens.refreshToken).toBeDefined();
  });

  test('should reject invalid credentials', async () => {
    const user = await createTestUser();

    const result = await login({
      email: user.email,
      password: 'WrongPassword',
    });

    expect(result.success).toBe(false);
    expect(result.error).toContain('Invalid credentials');
  });

  test('should lock account after 5 failed attempts', async () => {
    const user = await createTestUser();

    for (let i = 0; i < 5; i++) {
      await login({ email: user.email, password: 'wrong' });
    }

    const result = await login({
      email: user.email,
      password: 'TestPass@123', // Correct password
    });

    expect(result.success).toBe(false);
    expect(result.error).toContain('locked');
  });
});
```

#### 3. Token 測試

```typescript
describe('Token Management', () => {
  test('should generate valid JWT tokens', () => {
    const token = generateAccessToken({
      userId: '123',
      email: 'test@example.com',
      role: 'user',
    });

    const decoded = verifyToken(token);
    expect(decoded.userId).toBe('123');
    expect(decoded.email).toBe('test@example.com');
  });

  test('should refresh tokens', async () => {
    const { refreshToken } = await login({
      email: 'test@example.com',
      password: 'pass',
    });

    const result = await refreshTokens(refreshToken);

    expect(result.accessToken).toBeDefined();
    expect(result.refreshToken).toBeDefined();
    expect(result.refreshToken).not.toBe(refreshToken); // New refresh token
  });

  test('should invalidate old refresh token after refresh', async () => {
    const { refreshToken } = await login({
      email: 'test@example.com',
      password: 'pass',
    });

    await refreshTokens(refreshToken);

    // Try to use old refresh token
    const result = await refreshTokens(refreshToken);
    expect(result.success).toBe(false);
  });

  test('should reject expired tokens', async () => {
    const expiredToken = generateAccessToken(
      { userId: '123' },
      { expiresIn: '-1h' } // Expired 1 hour ago
    );

    const result = verifyToken(expiredToken);
    expect(result).toBeNull();
  });
});
```

### 整合測試

#### 1. 完整註冊流程

```typescript
describe('E2E Registration Flow', () => {
  test('should complete registration and login', async () => {
    // 1. Register
    const registerResponse = await apiPost('/auth/register', {
      email: 'newuser@example.com',
      password: 'NewUser@123',
      name: 'New User',
      agreeToTerms: true,
    });

    expect(registerResponse.status).toBe(201);
    expect(registerResponse.data.tokens).toBeDefined();

    // 2. Use access token to access protected route
    const dashboardResponse = await apiGet('/projects', {
      headers: {
        Authorization: `Bearer ${registerResponse.data.tokens.accessToken}`,
      },
    });

    expect(dashboardResponse.status).toBe(200);

    // 3. Logout
    const logoutResponse = await apiPost(
      '/auth/logout',
      {},
      {
        headers: {
          Authorization: `Bearer ${registerResponse.data.tokens.accessToken}`,
        },
      }
    );

    expect(logoutResponse.status).toBe(200);

    // 4. Verify token is invalidated
    const retryResponse = await apiGet('/projects', {
      headers: {
        Authorization: `Bearer ${registerResponse.data.tokens.accessToken}`,
      },
    });

    expect(retryResponse.status).toBe(401);
  });
});
```

#### 2. 前端整合測試

```typescript
describe('Frontend Auth Flow', () => {
  test('should register via UI', async () => {
    await page.goto('/register');

    await page.fill('[name="email"]', 'test@example.com');
    await page.fill('[name="password"]', 'Test@123');
    await page.fill('[name="confirmPassword"]', 'Test@123');
    await page.fill('[name="name"]', 'Test User');
    await page.check('[name="agreeToTerms"]');

    await page.click('button[type="submit"]');

    // Should redirect to dashboard
    await page.waitForURL('/dashboard');
    expect(page.url()).toContain('/dashboard');
  });

  test('should show validation errors', async () => {
    await page.goto('/register');

    await page.fill('[name="email"]', 'invalid-email');
    await page.fill('[name="password"]', 'weak');

    await page.click('button[type="submit"]');

    await expect(page.locator('.error')).toContainText('valid email');
    await expect(page.locator('.error')).toContainText('at least 8 characters');
  });
});
```

### 安全測試

```typescript
describe('Security Tests', () => {
  test('should prevent SQL injection', async () => {
    const result = await login({
      email: "admin'--",
      password: 'anything',
    });

    expect(result.success).toBe(false);
  });

  test('should prevent timing attacks', async () => {
    // Measure time for existing user
    const start1 = Date.now();
    await login({ email: 'existing@example.com', password: 'wrong' });
    const time1 = Date.now() - start1;

    // Measure time for non-existing user
    const start2 = Date.now();
    await login({ email: 'nonexisting@example.com', password: 'wrong' });
    const time2 = Date.now() - start2;

    // Times should be similar (within 100ms)
    expect(Math.abs(time1 - time2)).toBeLessThan(100);
  });

  test('should rate limit registration attempts', async () => {
    const attempts = [];

    for (let i = 0; i < 20; i++) {
      attempts.push(
        register({
          email: `user${i}@example.com`,
          password: 'Test@123',
          name: `User ${i}`,
          agreeToTerms: true,
        })
      );
    }

    const results = await Promise.all(attempts);
    const rateLimited = results.some((r) => r.status === 429);

    expect(rateLimited).toBe(true);
  });
});
```

## 技術實作細節

### 後端實作 (Go)

```go
// models/user.go
type User struct {
    ID           string    `json:"id" gorm:"primaryKey"`
    Email        string    `json:"email" gorm:"unique;not null"`
    Password     string    `json:"-" gorm:"not null"`
    Name         string    `json:"name"`
    Company      string    `json:"company"`
    Role         string    `json:"role" gorm:"default:'user'"`
    IsVerified   bool      `json:"is_verified" gorm:"default:false"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

type RefreshToken struct {
    ID        string    `gorm:"primaryKey"`
    UserID    string    `gorm:"not null"`
    Token     string    `gorm:"unique;not null"`
    ExpiresAt time.Time `gorm:"not null"`
    CreatedAt time.Time
}

// services/auth_service.go
type AuthService struct {
    db     *gorm.DB
    jwt    *JWTService
    mailer *MailService
}

func (s *AuthService) Register(req RegisterRequest) (*RegisterResponse, error) {
    // 1. Validate input
    if err := validateRegisterRequest(req); err != nil {
        return nil, err
    }

    // 2. Check email uniqueness
    var existingUser User
    if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
        return nil, errors.New("email already registered")
    }

    // 3. Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
    if err != nil {
        return nil, err
    }

    // 4. Create user
    user := User{
        ID:       uuid.NewString(),
        Email:    req.Email,
        Password: string(hashedPassword),
        Name:     req.Name,
        Company:  req.Company,
        Role:     "user",
    }

    if err := s.db.Create(&user).Error; err != nil {
        return nil, err
    }

    // 5. Generate tokens
    accessToken, err := s.jwt.GenerateAccessToken(user.ID, user.Email, user.Role)
    if err != nil {
        return nil, err
    }

    refreshToken, err := s.jwt.GenerateRefreshToken(user.ID)
    if err != nil {
        return nil, err
    }

    // 6. Send welcome email (async)
    go s.mailer.SendWelcomeEmail(user.Email, user.Name)

    return &RegisterResponse{
        User: user,
        Tokens: Tokens{
            AccessToken:  accessToken,
            RefreshToken: refreshToken,
        },
    }, nil
}
```

### 前端實作 (React + TypeScript)

```typescript
// hooks/useAuth.ts
export const useAuth = () => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const register = async (data: RegisterRequest) => {
    setIsLoading(true);
    try {
      const response = await api.post('/auth/register', data);
      const { user, tokens } = response.data;

      // Store tokens
      localStorage.setItem('accessToken', tokens.accessToken);
      localStorage.setItem('refreshToken', tokens.refreshToken);

      setUser(user);
      return { success: true };
    } catch (error) {
      return { success: false, error: error.message };
    } finally {
      setIsLoading(false);
    }
  };

  const login = async (data: LoginRequest) => {
    // Similar implementation
  };

  const logout = async () => {
    await api.post('/auth/logout');
    localStorage.removeItem('accessToken');
    localStorage.removeItem('refreshToken');
    setUser(null);
  };

  return { user, register, login, logout, isLoading };
};

// pages/register.tsx
export default function RegisterPage() {
  const { register } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (data: RegisterRequest) => {
    const result = await register(data);
    if (result.success) {
      navigate('/dashboard');
    }
  };

  return (
    <RegisterForm onSubmit={handleSubmit} />
  );
}
```

## 完成標準 (Definition of Done)

- [ ] 所有 API 端點實作完成
- [ ] 密碼安全 hash (bcrypt)
- [ ] JWT token 生成和驗證
- [ ] 註冊/登入頁面完成
- [ ] 表單驗證完整
- [ ] 錯誤處理完善
- [ ] 所有測試通過 (覆蓋率 ≥ 90%)
- [ ] 安全測試通過
- [ ] 登入失敗鎖定機制
- [ ] 歡迎郵件發送
- [ ] API 文件完整
