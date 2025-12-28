# PPLX2API

将 Perplexity 会话能力封装为 OpenAI 兼容接口的本地网关，支持流式输出、图片理解、搜索模式与多 Cookie 轮询，并提供可视化管理面板。

## 一、使用前准备

### 1) 获取 Perplexity Cookie
在浏览器登录 https://www.perplexity.ai/

打开开发者工具（F12）-> Application/Storage -> Cookies -> 找到  
`__Secure-next-auth.session-token` 的 value，复制出来。

可以准备多个 Cookie，用于轮询/重试。

### 2) 确认端口
默认端口为 `8482`。如果你想修改端口，后面在 `config.json` 里改 `address`。

---

## 二、本地使用（推荐先用这个）

### 方式 A：直接运行源码
```bash
go run .
```
首次启动会在当前目录生成 `config.json`。

然后访问：
```
http://localhost:8482/
```

### 方式 B：编译后运行
```bash
go build -o pplx2api .
./pplx2api
```

---

## 三、服务器/部署方式

### 方式 A：Docker（推荐）
```bash
docker build -t pplx2api .
docker run -d \
  -p 8482:8482 \
  -v $(pwd)/config.json:/app/config.json \
  --name pplx2api \
  pplx2api
```

说明：
- `-v $(pwd)/config.json:/app/config.json` 用于持久化配置  
- 修改端口时请同时修改 `config.json` 的 `address`

### 方式 B：Docker Compose
```bash
docker compose up -d
```

### 方式 C：服务器裸机运行（不使用 Docker）
```bash
git clone <你的仓库地址>
cd pplx2api
go build -o pplx2api .
./pplx2api
```

确保服务器防火墙/安全组放行 `8482` 端口。

---

## 四、管理面板使用方法
访问 `http://localhost:8482/` 进入控制台。

默认管理密钥：`123456`（建议首次使用后立即修改）。

操作步骤：
1. 填写管理密钥（`config.json` 里的 `apikey`）  
2. 点击“加载”  
3. 在 Cookie 列表中新增/粘贴 Cookie  
4. 点击“刷新列表”拉取模型  
5. 选择默认模型  
6. 保存配置

---

## 五、配置说明（config.json）
示例：
```json
{
  "sessions": ["YOUR_PPLX_SESSION_TOKEN"],
  "address": "0.0.0.0:8482",
  "apikey": "123456",
  "proxy": "",
  "is_incognito": true,
  "max_chat_history_length": 10000,
  "no_role_prefix": false,
  "search_result_compatible": false,
  "prompt_for_file": "You must immerse yourself...",
  "ignore_search_result": false,
  "ignore_model_monitoring": false,
  "is_max_subscribe": false,
  "default_model": "claude-3.7-sonnet",
  "force_model": ""
}
```

字段说明：
| 字段 | 说明 | 默认值 |
| --- | --- | --- |
| `sessions` | Perplexity Cookie（`__Secure-next-auth.session-token`）列表 | `[]` |
| `address` | 服务监听地址 | `0.0.0.0:8482` |
| `apikey` | API/管理密钥 | `123456` |
| `proxy` | 代理地址 | `""` |
| `is_incognito` | 隐身模式 | `true` |
| `max_chat_history_length` | 超过该长度将上传为文件 | `10000` |
| `no_role_prefix` | 禁用角色前缀 | `false` |
| `search_result_compatible` | 搜索结果兼容模式 | `false` |
| `prompt_for_file` | 文件上传时提示词 | 默认内置 |
| `ignore_search_result` | 隐藏搜索结果 | `false` |
| `ignore_model_monitoring` | 忽略模型监控 | `false` |
| `is_max_subscribe` | Max 订阅开关 | `false` |
| `default_model` | 默认模型 | `claude-3.7-sonnet` |
| `force_model` | 强制模型（覆盖请求） | `""` |

---

## 六、API 使用

### 认证
请求头加入：
```
Authorization: Bearer YOUR_API_KEY
```

### 获取模型列表
```bash
curl -H "Authorization: Bearer YOUR_API_KEY" \
  http://localhost:8482/v1/models
```

### Chat Completions
```bash
curl -X POST http://localhost:8482/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -d '{
    "model": "claude-3.7-sonnet",
    "messages": [
      {"role": "user", "content": "你好"}
    ],
    "stream": false
  }'
```

### 流式输出
```bash
curl -N -X POST http://localhost:8482/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -d '{
    "model": "claude-3.7-sonnet",
    "messages": [{"role": "user", "content": "讲个笑话"}],
    "stream": true
  }'
```

### 图片理解
```bash
curl -X POST http://localhost:8482/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -d '{
    "model": "claude-3.7-sonnet",
    "messages": [
      {
        "role": "user",
        "content": [
          {"type": "text", "text": "这张图里有什么？"},
          {"type": "image_url", "image_url": {"url": "data:image/jpeg;base64,..." }}
        ]
      }
    ]
  }'
```

### 搜索模式
模型名后加 `-search`：
```
"model": "grok-4.1-think-search"
```

---

## 七、常见问题
- **页面打不开**：检查 `config.json` 的 `address` 是否与访问端口一致，默认 `8482`。
- **Invalid API key**：使用 `config.json` 的 `apikey`，默认是 `123456`。
- **模型列表为空**：先填写管理密钥，再点击“刷新列表”。
- **config.json 找不到**：确保启动目录与 `config.json` 在同一目录（Docker 需要挂载）。

---
基于https://github.com/yushangxiao/pplx2api.git项目的修改

## License
MIT，详见 `LICENSE`。
