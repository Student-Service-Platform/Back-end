> [!WARNING]
> 最终结果以最终具体api为准！此文档仅供参考！

# 学生服务平台 API 文档

> [!IMPORTANT]  
> 1. 未标注返回值的接口即只返回  
>    {"code":"{status}",  data: null,msg: "{err}",}
> 2. POST返回值默认指的是返回结构中的data项

## 1. 认证

### POST /api/auth/login

- **描述**:
  - 用户登录接口。
  - 相关错误代码参考`状态码查询`页

- **参数**:
  - `username`: 用户名
  - `password`: 密码
  - `isadmin`:是否管理员

### POST /api/auth/adminlogin

- **描述**:
  - 管理员登录接口。
  - 相关错误代码参考`状态码查询`页

- **参数**:
  - `username`: 用户名
  - `password`: 密码
  - `isadmin`:是否管理员

## 2. 注册

### POST /api/reg

- **描述**:
  - 注册新用户
  - 保存至普通用户
  - 注意前后端检测学号、手机号、邮箱、密码格式
  - 学号即为用户ID
- **参数**:
  - `student_id`: 学号
  - `username`: 用户名
  - `password`: 密码
  - `mail_auth`: 邮箱验证成功返回值
  - `phone`: 手机号

## 3. 用户

### GET /api/user/profile

- **描述**:
  - 获取当前用户的个人资料。
  - authMiddleware

### PUT /api/user/profile/:name

- **描述**:
  - 更新当前用户的个人资料中的`name`项。
  - authMiddleware

- **参数**:
  - `value`

### GET /api/user/feedbacks

- **描述**:
  - 获取当前用户的反馈列表。
  - authMiddleware
- **参数**:
  - `limit`: 每页显示数量
  - `offset`: 分页偏移量

### POST /api/user/resetpwd

- **描述**:
  - 找回密码
- **参数**:
  - `mail_auth`
  - `password`
  - `user_id`

## 4. 反馈

### POST /api/feedback

- **描述**:
  - 提交新的反馈。
  - authMiddleware

- **参数**:
  - `user_id`: 用户ID
  - `title`: 反馈标题
  - `description`: 反馈内容
  - `category`: 反馈类别
  - `is_urgent`: 是否紧急
  - `is_anonymous`: 是否匿名
  - `images[]`: 图片附件

### GET /api/feedback

- **描述**:
  - 获取反馈列表。
  - authMiddleware

- **参数**:
  - `user_id`: 请求用户id
  - `status`: 状态过滤
  - `category`: 类别过滤
  - `limit`: 每页显示数量
  - `offset`: 分页偏移量

### GET /api/feedback/:id

- **描述**:
  - 获取指定`id`的反馈详情。

- **参数**:
  - `user_id`: 请求用户id

### POST /api/feedback/:id/reply

- **描述**:
  - 允许普通用户或管理员回复反馈。
  - authMiddleware

- **参数**:
  - `user_id`:  回复者ID
  - `reply`: 回复内容

### POST /api/feedback/:id/admin

- **描述**:
  - 接单/取消接单（仅限管理员）。
  - authMiddleware

- **参数**:
  - `action`: 动作类型 (`0`accept 或 `1`cancel)

### POST /api/feedback/:id/mark

- **描述**:
  - 标记为垃圾信息（仅限超级管理员）。

- **参数**:
  - `confirmation`: 确认标记 (`true` 或 `false`)

## 5. 邮件

### POST /api/mail/

- **描述**:
  - 发送邮件通知给用户。
  - authMiddleware

- **参数**:
  - `recipient`: 收件人
  - `subject`: 主题
  - `body`: 正文

### POST /api/mail/reg

- **描述**:
  - 请求邮箱验证码
  - 注意防止同ip滥用，前后端都应做"冷却"
  - 需要另建一个数据库
  - 此数据库中的条例将会在5min后自动删除
  - 结构可为 `ID` `MAIL` `CODE` `AUTH`
  - 其中`ID`为标识符, `AUTH`为uuid字符，用于传递验证成功消息
- **参数**:
  - `mail`: 邮箱地址

### POST /api/mail/auth

- **描述**:
  - 检查验证码，验证邮箱可用
- **参数**:
  - `mail`: 邮箱地址
  - `code`: 验证码
- **返回**:
  - `auth`: 验证成功的uuid

## 6. 预设回复

### GET /api/admin/replies

- **描述**:
  - 获取预设回复列表（仅限管理员）。
  - authMiddleware

### POST /api/admin/replies

- **描述**:
  - 添加新的预设回复（仅限管理员）。
  - authMiddleware

- **参数**:
  - `text`: 回复文本

## 10. 超级管理员特有部分

### POST /api/suadmin/add

- **描述**:
  - 添加管理员（仅限超级管理员）。
  - authMiddleware

- **requires**
  - `pwd` 普通管理员密码
  - `admin_nam`e 普通管理员用户名

# 以下尚未修订完成

## 7. 错误处理

- **描述**:
  - 使用全局错误处理中间件来统一处理错误，并返回相应的错误码和错误消息。

## 8. 页面加速

- **描述**:
  - 使用CDN来托管静态资源如JavaScript, CSS等。

## 9. 数据可视化

### GET /api/admin/dashboard

- **描述**:
  - 获取用于超级管理员界面的数据统计信息。
