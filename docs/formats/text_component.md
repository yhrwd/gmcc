# 文本组件规范

## 概述

**文本组件（Text Component）**，也称为原始 JSON 文本，用于 Minecraft 向客户端发送和显示富文本内容。

- 官方文档：[Minecraft Wiki - 文本组件](https://zh.minecraft.wiki/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6)
- 协议版本：774 (1.21.11)

## 基础结构

文本组件支持三种格式：

### 1. 字符串格式

纯文本的简写形式：

```json
"Hello World"
```

### 2. 复合标签格式

完整的组件结构：

```json
{
    "text": "Hello",
    "color": "red",
    "bold": true,
    "extra": [
        {"text": " World", "bold": false}
    ]
}
```

### 3. 列表格式

多个组件拼接：

```json
[
    {"text": "Hello", "color": "red"},
    {"text": " World", "color": "blue"}
]
```

## 组件类型

| 类型 | type 值 | 说明 |
|------|---------|------|
| 纯文本 | `text` | 静态文本内容 |
| 本地化文本 | `translatable` | 使用语言文件中的翻译键 |
| 键位绑定 | `keybind` | 显示当前绑定的按键 |
| 记分板 | `score` | 显示记分板分数 |
| 实体选择器 | `selector` | 显示实体名称 |
| NBT | `nbt` | 显示 NBT 数据 |
| 精灵图 | `object` / `sprite` | 显示纹理图集中的精灵图 |

## 组件样式

### 颜色

```json
{"text": "Red text", "color": "red"}
{"text": "Custom color", "color": "#FF5733"}
```

支持的颜色名称：
- `black`, `dark_blue`, `dark_green`, `dark_aqua`, `dark_red`, `dark_purple`, `gold`
- `gray`, `dark_gray`, `blue`, `green`, `aqua`, `red`, `light_purple`, `yellow`, `white`

### 文字样式

```json
{
    "text": "Styled",
    "bold": true,
    "italic": true,
    "underlined": true,
    "strikethrough": true,
    "obfuscated": false
}
```

### 字体

```json
{"text": "Custom font", "font": "minecraft:default"}
```

## 组件继承

子组件通过 `extra` 继承父组件样式：

```json
{
    "text": "A",
    "color": "red",
    "extra": [
        "B",
        {"text": "C", "color": "yellow"}
    ]
}
```
渲染结果：`A` 为红色，`B` 继承红色，`C` 为黄色。

## 交互事件

### 点击事件

```json
{
    "text": "Click me",
    "click_event": {
        "action": "open_url",
        "url": "https://example.com"
    }
}
```

动作类型：
- `open_url` - 打开 URL
- `run_command` - 执行命令
- `suggest_command` - 建议命令到聊天栏
- `copy_to_clipboard` - 复制到剪贴板
- `change_page` - 翻页（仅限成书）

### 悬停事件

```json
{
    "text": "Hover me",
    "hover_event": {
        "action": "show_text",
        "value": {"text": "Tooltip content"}
    }
}
```

动作类型：
- `show_text` - 显示文本
- `show_item` - 显示物品
- `show_entity` - 显示实体信息

### 插入文本

```json
{
    "text": "Shift+Click me",
    "insertion": "/say inserted text"
}
```

## 动态组件

### 本地化文本

```json
{
    "translate": "chat.type.text",
    "with": [{"selector": "@s"}]
}
```

### 记分板

```json
{
    "score": {
        "name": "@p",
        "objective": "kills"
    }
}
```

### 实体选择器

```json
{
    "selector": "@p",
    "separator": ", "
}
```

## 颜色代码（Legacy）

Minecraft 也支持 § 颜色代码：

| 代码 | 颜色 | 代码 | 颜色 |
|------|------|------|------|
| §0 | 黑色 | §8 | 深灰色 |
| §1 | 深蓝 | §9 | 蓝色 |
| §2 | 深绿 | §a | 绿色 |
| §3 | 深青 | §b | 青色 |
| §4 | 深红 | §c | 红色 |
| §5 | 深紫 | §d | 浅紫 |
| §6 | 金色 | §e | 黄色 |
| §7 | 灰色 | §f | 白色 |

格式代码：
- §k - 随机字符
- §l - 粗体
- §m - 删除线
- §n - 下划线
- §o - 斜体
- §r - 重置

## ANSI 转换

在终端显示时，文本组件可转换为 ANSI 转义序列：

```
§0 → \033[30m (黑色)
§4 → \033[31m (红色)
§a → \033[32m (绿色)
§e → \033[33m (黄色)
§1 → \033[34m (蓝色)
§d → \033[35m (紫色)
§3 → \033[36m (青色)
§f → \033[37m (白色)
§r → \033[0m (重置)
```

## 参考

- [文本组件 - Minecraft Wiki](https://zh.minecraft.wiki/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6)
- [格式化代码 - Minecraft Wiki](https://zh.minecraft.wiki/w/%E6%A0%BC%E5%BC%8F%E5%8C%96%E4%BB%A3%E7%A0%81)
- [Raw JSON Text Format](https://minecraft.wiki/w/Raw_JSON_text_format)
