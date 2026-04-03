# GMCC API-ONLY 版本

这是 GMCC 的纯 API 版本，移除了所有前端界面，只提供 RESTful API 服务。

## 启动方式

```bash
# 构建
go build -o gmcc.exe ./cmd/gmcc

# 运行
./gmcc.exe
```

## API 接口

### 公开接口（无需认证）

- `GET /api/status` - 获取集群状态
- `GET /api/accounts` - 获取账号列表
- `GET /api/accounts/:id` - 获取单个账号详情
- `POST /api/auth/verify` - 验证操作密码
- `POST /api/auth/microsoft/init` - 初始化 Microsoft 认证
- `POST /api/auth/microsoft/poll` - 轮询 Microsoft 认证状态

### 受保护接口（需要密码验证）

所有受保护接口需要在请求体中包含 `password` 字段：

```json
{
    "password": "your_password",
    "...": "其他参数"
}
```

- `POST /api/instances/:id/start` - 启动实例
- `POST /api/instances/:id/stop` - 停止实例
- `POST /api/instances/:id/restart` - 重启实例
- `DELETE /api/instances/:id` - 删除实例
- `POST /api/accounts` - 创建账号
- `DELETE /api/accounts/:id` - 删除账号
- `POST /api/passwords` - 创建密码
- `DELETE /api/passwords/:id` - 删除密码
- `GET /api/logs/operations` - 获取操作日志

## 配置文件

配置文件 `config.yaml` 保持不变，但 Web 相关配置仍然有效：

```yaml
web:
    bind: "0.0.0.0:8080"
    auth:
        token_expiry: 5m0s
        audit_log_retention_days: 30
        passwords:
            - id: admin
              password: password123
              enabled: true
              note: 管理员密码
    token_vault:
        scrypt_n: 1048576
        scrypt_r: 8
        scrypt_p: 1
```

## 使用示例

### 获取集群状态

```bash
curl http://localhost:8080/api/status
```

### 启动实例（需要密码）

```bash
curl -X POST http://localhost:8080/api/instances/account1/start \
  -H "Content-Type: application/json" \
  -d '{"password": "password123"}'
```

### 获取账号列表

```bash
curl http://localhost:8080/api/accounts
```

## 注意事项

1. 此版本不提供任何 Web 界面，仅提供 API 服务
2. 所有前端相关的路由和静态文件服务已被移除
3. API 功能保持完整，包括认证、审计和集群管理
4. 适合用于集成到其他系统或自定义前端实现
