# 字符映射功能 - 快速使用指南

## 第一步：查看服务器发送的特殊字符

1. 启动客户端，`config.yaml`中设置`debug: true`
2. 连接到服务器，观察聊天消息中的Unicode信息：

```
[INFO] 收到聊天消息: type=system sender= plain="消息内容"
[INFO] 特殊字符: ''(U+E000,pos=0), ''(U+E001,pos=1)
```

记录这些Unicode编码（如`U+E000`、`U+E001`）。

## 第二步：创建字符映射

### 方法一：使用命令行工具

```bash
# 交互模式（推荐）
./charmap.exe interactive

# 在交互模式中：
> add \uE000 "左上角" "┌"
已添加映射: \uE000 () -> "┌" [左上角]
> add \uE001 "横线" "─"
已添加映射: \uE001 () -> "─" [横线]
> enable
已启用字符替换功能
> save
> exit
```

### 方法二：直接编辑配置文件

编辑 `.charmap/charmap.json`：

```json
{
  "enable_replace": true,
  "show_unicode_info": true,
  "mappings": {
    "\\uE000": {
      "description": "GUI左上角",
      "replace_with": "┌"
    },
    "\\uE001": {
      "description": "GUI横线",
      "replace_with": "─"
    }
  }
}
```

## 第三步：验证映射效果

```bash
# 测试文本替换
./charmap.exe analyze "测试字符E000和E001"

# 预期输出：
原文: "测试字符E000和E001"
替换: "测试字符┌和─"
特殊字符: ...
```

## 第四步：运行客户端

```bash
./gmcc.exe
```

客户端会自动加载字符映射配置，聊天消息中的特殊字符将被替换。

## 常用Unicode字符参考

### GUI边框字符
```
╔ ═ ╗ ║ ╚ ╝  - 双线边框
┌ ─ ┐ │ └ ┘  - 单线边框
┏ ━ ┓ ┃ ┗ ┛  - 粗线边框
```

### 图标字符
```
● ○ ◐ ◑  - 圆形
★ ☆     - 星形
♥ ♡     - 心形
→ ← ↑ ↓  - 箭头
```

### 货币符号
```
$ € £ ¥ ¢  - 货币
● ◆ ■ ▲   - 资源图标
```

## 故障排除

### 问题：看不到Unicode信息

**解决方法：**
```bash
# 确保启用显示
./charmap.exe show-unicode on

# 或在配置文件中设置
{
  "show_unicode_info": true
}
```

### 问题：字符替换不生效

**解决方法：**
```bash
# 确保启用替换
./charmap.exe enable

# 检查映射是否正确
./charmap.exe list

# 确保Unicode编码正确格式为 \uXXXX
```

### 问题：保存映射失败

**解决方法：**
```bash
# 手动创建目录
mkdir .charmap

# 检查权限
# Windows: 右键 -> 属性 -> 安全
# Linux/Mac: chmod 755 .charmap
```

## 实际案例

### 案例：服务器经济系统

服务器发送消息：`价格: 100金币`

观察到：`特殊字符: ''(U+E010,pos=5)`

添加映射：
```bash
./charmap.exe add \uE010 "金币图标" "●"
./charmap.exe enable
```

替换结果：`价格: 100●金币`

### 案例：服务器GUI菜单

服务器发送菜单：
```
╔═══════════╗
║  商店菜单  ║
╚═══════════╝
```

观察到多个Unicode编码，批量添加：

```bash
./charmap.exe interactive
> add \uE000 "左上角" "╔"
> add \uE001 "横线" "═"
> add \uE002 "右上角" "╗"
> add \uE003 "竖线" "║"
> add \uE004 "左下角" "╚"
> add \uE005 "右下角" "╝"
> enable
> exit
```

## 高级技巧

### 批量导入映射

创建JSON文件并导入：

```bash
# 创建 custom_mappings.json
{
  "enable_replace": true,
  "show_unicode_info": true,
  "mappings": {
    "\\uE000": {"description": "图标1", "replace_with": "●"},
    "\\uE001": {"description": "图标2", "replace_with": "○"}
  }
}

# 复制到配置目录
cp custom_mappings.json .charmap/charmap.json
```

### 临时禁用映射

```bash
# 禁用替换但保留配置
./charmap.exe disable

# 需要时重新启用
./charmap.exe enable
```

### 调试模式

查看详细的字符信息：

```bash
# 分析包含特殊字符的文本
./charmap.exe analyze "Test字符U+E000"

# 输出更详细的Unicode信息：
原文: "Test字符U+E000"
替换: "Test字符"
特殊字符: '字'(U+5B57,pos=4), '符'(U+7B26,pos=5), ''(U+E000,pos=7)
```

## 下一步

- 阅读 [完整文档](charmap.md)
- 查看 [配置示例](../charmap_template.json)
- 参考 [命令列表](#命令行工具)

## 获取帮助

```bash
# 查看所有命令
./charmap.exe help

# 交互模式中查看帮助
./charmap.exe interactive
> help
```