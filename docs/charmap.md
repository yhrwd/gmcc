# 字符映射功能

## 功能概述

GMCC 提供了强大的字符映射功能，允许您：

1. **查看聊天消息中的Unicode字符编码** - 识别服务器发送的自定义字符
2. **替换特殊字符** - 将服务器材质包的自定义字符替换为终端可显示的字符
3. **管理映射关系** - 灵活配置字符映射规则

这对于需要对接服务器材质包（使用Unicode私用区字符）的场景非常有用。

## 快速开始

### 1. 生成映射模板

```bash
go run cmd/charmap/main.go template
```

这会生成一个 `charmap_template.json` 文件，包含常用字符映射示例：

```json
{
  "enable_replace": true,
  "show_unicode_info": true,
  "mappings": {
    "\\uE000": {
      "description": "方块图标 - 左上角",
      "replace_with": "┌"
    },
    "\\uE001": {
      "description": "方块图标 - 横线",
      "replace_with": "─"
    }
  }
}
```

### 2. 查看当前映射

```bash
go run cmd/charmap/main.go list
```

输出示例：
```
字符替换状态: true
显示Unicode信息: true

当前映射:
  \uE000 () -> "┌" [方块图标 - 左上角]
  \uE001 () -> "─" [方块图标 - 横线]
```

### 3. 分析文本中的特殊字符

```bash
go run cmd/charmap/main.go analyze "测试文本内容"
```

输出示例：
```
原文: "测试文本内容"
替换: "测试文本内容"
特殊字符: '测'(U+6D4B,pos=0), '试'(U+8BD5,pos=1), '文'(U+6587,pos=2), '本'(U+672C,pos=3), '内'(U+5185,pos=4), '容'(U+5BB9,pos=5)
```

## 配置说明

### 配置文件位置

字符映射配置文件默认保存在 `.charmap/charmap.json`。

### 配置项说明

| 配置项 | 类型 | 说明 |
|--------|------|------|
| `enable_replace` | bool | 是否启用字符替换功能 |
| `show_unicode_info` | bool | 是否在聊天消息中显示Unicode字符信息 |
| `mappings` | map | 字符映射规则 |

### 映射规则格式

```json
{
  "\\uE000": {
    "description": "字符描述",
    "replace_with": "替换后的字符"
  }
}
```

## 命令行工具

### 查看所有命令

```bash
go run cmd/charmap/main.go help
```

### 添加映射

```bash
go run cmd/charmap/main.go add \uE000 "左上角" "┌"
```

### 删除映射

```bash
go run cmd/charmap/main.go remove \uE000
```

### 启用/禁用替换

```bash
# 启用
go run cmd/charmap/main.go enable

# 禁用
go run cmd/charmap/main.go disable
```

### 设置是否显示Unicode信息

```bash
# 显示
go run cmd/charmap/main.go show-unicode on

# 隐藏
go run cmd/charmap/main.go show-unicode off
```

### 交互模式

```bash
go run cmd/charmap/main.go interactive
```

在交互模式下可以实时管理映射：

```
> list                    # 列出映射
> add \uE000 描述 替换     # 添加映射
> remove \uE000           # 删除映射
> analyze 测试文本         # 分析文本
> enable                  # 启用替换
> exit                    # 退出
```

## 运行时使用

### 在客户端中启用

GMCC 客户端启动时会自动加载 `.charmap/charmap.json` 配置文件。

### 查看Unicode信息

当启用 `show_unicode_info` 后，客户端接收到包含特殊字符的聊天消息时，会输出类似信息：

```
[INFO] 收到聊天消息: type=system sender= plain="示例文本"
[INFO] 特殊字符: ''(U+E000,pos=0), ''(U+E001,pos=1)
```

这让你能够识别服务器发送的自定义字符编码。

### 字符替换效果

启用 `enable_replace` 后，聊天消息中的特殊字符会被自动替换：

```
原始消息: "玩家名"
替换后: "┌─┐玩家名"
```

## Unicode 私用区

Minecraft 服务器插件通常使用 Unicode 私用区来存储自定义图标：

- **基础多文种平面私用区**: U+E000 ~ U+F8FF
- **补充私用区A**: U+F0000 ~ U+FFFFD
- **补充私用区B**: U+100000 ~ U+10FFFD

常见用途：
- GUI 边框和装饰
- 货币和资源图标
- 状态指示符号
- 自定义表情符号

## 实际应用场景

### 场景 1: 服务器材质包图标

服务器使用材质包自定义GUI，发送包含 U+E000-U+E0FF 范围字符的消息：

1. 观察客户端输出的Unicode信息
2. 为每个字符添加映射
3. 客户端将自动替换为终端可显示字符

### 场景 2: 经济系统显示

服务器发送金币图标（如 U+E010）：

```bash
charmap add \uE010 "金币" "●"
```

经济消息会从 `价格: 100金币` 变为 `价格: 100●金币`

### 场景 3: 状态显示

服务器使用图标显示玩家状态：

```bash
charmap add \uE020 "在线" "●"
charmap add \uE021 "离线" "○"
charmap add \uE022 "忙碌" "◐"
```

## 技术细节

### 字符识别

字符分析器会识别以下类型的特殊字符：

1. **非ASCII字符** (Unicode > 127)
2. **控制字符** (Unicode < 32)
3. **替换字符** (U+FFFD, 表示解码失败)

### 性能优化

- 使用读写锁保护配置访问
- 配置文件修改后立即生效
- 字符分析采用惰性策略，只在需要时执行

### 线程安全

所有公共方法都是线程安全的，可以在并发环境中使用：
- `AnalyzeText()` - 只读操作
- `ReplaceText()` - 只读操作
- `AddMapping()` - 写操作，自动加锁
- `RemoveMapping()` - 写操作，自动加锁

## 故障排除

### 问题：字符替换不生效

检查以下项：
1. `enable_replace` 是否为 `true`
2. 映射的 Unicode 编码是否正确
3. 配置文件是否正确加载

### 问题：Unicode信息不显示

检查以下项：
1. `show_unicode_info` 是否为 `true`
2. 日志级别是否为 Debug 或 Info
3. 消息是否包含特殊字符

### 问题：配置文件无法保存

检查以下项：
1. `.charmap` 目录是否存在且可写
2. 是否有足够的磁盘空间
3. 文件权限是否正确

## 示例配置

以下是一个完整的配置示例，适用于常见的服务器材质包：

```json
{
  "enable_replace": true,
  "show_unicode_info": true,
  "mappings": {
    "\\uE000": {"description": "边框-左上", "replace_with": "╔"},
    "\\uE001": {"description": "边框-横线", "replace_with": "═"},
    "\\uE002": {"description": "边框-右上", "replace_with": "╗"},
    "\\uE003": {"description": "边框-竖线", "replace_with": "║"},
    "\\uE004": {"description": "边框-左下", "replace_with": "╚"},
    "\\uE005": {"description": "边框-右下", "replace_with": "╝"},
    "\\uE010": {"description": "金币", "replace_with": "●"},
    "\\uE011": {"description": "银币", "replace_with": "○"},
    "\\uE020": {"description": "心形", "replace_with": "♥"},
    "\\uE021": {"description": "星形", "replace_with": "★"},
    "\\uE030": {"description": "箭头-右", "replace_with": "→"},
    "\\uE031": {"description": "箭头-左", "replace_with": "←"}
  }
}
```

## 扩展开发

### 编程接口

```go
// 获取字符分析器
analyzer := client.GetCharacterAnalyzer()

// 添加映射
client.AddCharacterMapping("\\uE000", "描述", "替换字符")

// 移除映射
client.RemoveCharacterMapping("\\uE000")

// 启用/禁用替换
client.SetCharacterReplaceEnabled(true)

// 设置显示Unicode信息
client.SetShowUnicodeInfo(true)

// 生成模板
client.GenerateCharacterMappingTemplate("template.json")
```

### 自定义字符分析

```go
// 分析文本中的特殊字符
unicodeInfo := analyzer.AnalyzeText(text)
if unicodeInfo != "" {
    fmt.Println(unicodeInfo)
}

// 替换文本中的字符
replaced := analyzer.ReplaceText(text)
```

## 相关资源

- [Unicode 私用区](https://en.wikipedia.org/wiki/Private_Use_Areas)
- [Minecraft 材质包](https://minecraft.fandom.com/wiki/Resource_Pack)
- [终端字符集](https://en.wikipedia.org/wiki/Code_page_437)

## 更新日志

- **v1.0.0** - 初始版本，支持基本的字符映射和替换功能
- **v1.1.0** - 添加交互模式和命令行工具
- **v1.2.0** - 优化性能，增加线程安全性