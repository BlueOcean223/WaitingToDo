# API 指南

本文档详细描述了 WaitingToDo 后端项目的 API 端点，包括认证、用户、任务、好友和消息等模块。
提示：除 /auth 下的接口外，其余接口均在受保护分组下，调用时需要在请求头中携带 Authorization: Bearer <token>。

## 1. 认证 Auth

- POST /auth/login    - 登录
- POST /auth/register - 注册
- POST /auth/forget   - 忘记密码
- GET  /auth/captcha  - 获取验证码


## 2. 用户 User（需 JWT）

- POST /user/checkCaptcha - 发送验证码
- POST /user/reset        - 重置密码
- POST /user/update       - 修改个人信息
- GET  /user/info         - 获取用户信息


## 3. 文件上传 Upload（需 JWT）

- POST   /upload/img         - 上传图片
- POST   /upload/:id/file    - 上传文件（路径参数 id）
- DELETE /upload/deletefile  - 删除文件


## 4. 任务 Task（需 JWT）

- GET    /task/list             - 查询任务列表
- POST   /task/add              - 新增任务
- DELETE /task/delete           - 删除任务
- PUT    /task/update           - 修改任务
- GET    /task/urgent           - 获取紧急任务列表
- GET    /task/teamList         - 查询小组任务列表
- DELETE /task/team/delete      - 删除小组任务
- POST   /task/team/add         - 添加小组任务
- PUT    /task/team/complete    - 小组成员完成任务
- POST   /task/team/invite      - 邀请成员
- GET    /task/team/inviteCode  - 获取小组任务邀请码
- POST   /task/team/codejoin    - 通过邀请码加入小组任务


## 5. 好友 Friend（需 JWT）

- GET    /friend/info    - 根据 id 查询好友详情
- GET    /friend/list    - 获取好友列表
- GET    /friend/search  - 根据邮箱查询用户信息（用于添加好友）
- POST   /friend/add     - 添加好友请求
- DELETE /friend/delete  - 删除好友


## 6. 消息 Message（需 JWT）

- GET    /message/unreadCount  - 获取用户未读消息数量
- GET    /message/list         - 分页查询消息
- PUT    /message/update       - 更新消息
- DELETE /message/delete       - 删除消息
- PUT    /message/readAll      - 一键已读
- POST   /message/handle       - 处理请求
- POST   /message/add          - 添加消息


## 7. 认证与分组说明

- 放行路由：仅 /auth/*，由SetAuthRoutes 注册。
- 受保护路由：/user、/upload、/task、/friend、/message 等在受保护分组下注册，应用 JWT 中间件（需携带 Authorization: Bearer <token>）。
