# .knowledge 目录索引

本目录包含 Minecraft 协议和数据处理相关的知识库，供 AI 开发参考。

## 目录结构

### 1. MC_Protocol_Data/
Minecraft Java Edition 协议数据定义

**来源**: 基于 [PrismarineJS/minecraft-data](https://github.com/PrismarineJS/minecraft-data) 项目

**用途**:
- 用于 Wireshark 的 Minecraft 协议解析器 ([MC_Dissector](https://github.com/Nickid2018/MC_Dissector))
- Minecraft Wiki 协议查看器 ([mcw-calc](https://github.com/mc-wiki/mcw-calc))

**目录内容**:
- `java_edition/` - Java 版协议相关数据 (数据包定义、状态等)
- `.schema/` - 协议数据的 JSON Schema
- `.documentation/` - 协议文档

---

### 2. minecraft-data/
Minecraft 游戏数据规范 (语言无关的数据模块)

**支持版本**: 
- Java Edition: 0.30c ~ 1.21.10+
- Bedrock Edition: 0.14 ~ 1.26.0+

**提供的数据类型**:

| 数据类型 | 说明 |
|---------|------|
| Blocks | 方块数据 |
| Items | 物品数据 |
| Entities | 实体数据 |
| Biomes | 生物群系 |
| Recipes | 合成配方 |
| Effects | 状态/药水效果 |
| Enchantments | 附魔 |
| Particles | 粒子效果 |
| Sounds | 声音 |
| Commands | 命令树结构 |
| Protocol | 完整协议描述 (可自动生成协议实现) |
| Protocol Versions | 协议版本列表 |
| Windows | 窗口类型 |
| Materials | 工具速度属性 |
| Foods | 食物数据 |
| Map icons | 地图图标 |
| Instruments | 方块发出的声音 |

**使用此数据的项目**:
- mineflayer (Node.js bot库)
- node-minecraft-protocol (Node.js 协议库)
- flying-squid (Node.js 服务器库)
- SpockBot (Python bot库)
- ProtocolGen (Java)

**目录内容**:
- `data/` - 各版本游戏数据 JSON 文件
- `doc/` - 数据文档
- `schemas/` - JSON Schema
- `tools/` - 审计和测试工具

---

### 3. prismarine-chat/
Minecraft 聊天组件解析库 (JavaScript/Node.js)

**NPM**: [prismarine-chat](https://www.npmjs.com/package/prismarine-chat)

**功能**: 解析 Minecraft 聊天 JSON 消息格式

**主要 API**:

```js
const registry = require('prismarine-registry')('1.16')
const ChatMessage = require('prismarine-chat')(registry)

const msg = new ChatMessage({"text":"Hello"})
msg.toString()           // 转为纯文本
msg.toAnsi()             // 转为 ANSI 格式 (终端彩色)
msg.toHTML()             // 转为 HTML
msg.toMotd()             // 转为 MOTD 格式
```

**MessageBuilder** (构建聊天消息):
```js
const { MessageBuilder } = require('prismarine-chat')(registry)
const msg = new MessageBuilder()
  .setText('Hello')
  .setColor('red')
  .setBold(true)
  .toJSON()
```

**支持的方法**:
- `setText`, `setColor`, `setBold`, `setItalic`, `setUnderlined`, `setStrikethrough`, `setObfuscated`
- `setTranslate` (翻译键), `setScore`, `setClickEvent`, `setHoverEvent`
- `addExtra`, `addWith` (追加组件)
- `fromNotch()` - 转换旧版聊天格式
- `fromNetwork()` - 转换 1.19+ 网络消息

**目录内容**:
- `index.js` - 主库实现
- `index.d.ts` - TypeScript 类型定义
- `MessageBuilder.js` - 消息构建器
- `examples/` - 使用示例

---

### 4. prismarine-nbt/
NBT 格式解析/序列化库 (JavaScript/Node.js)

**NPM**: [prismarine-nbt](https://www.npmjs.com/package/prismarine-nbt)

**功能**: 解析和序列化 NBT (Named Binary Tag) 格式

**支持格式**:
- Big Endian (Java Edition 使用)
- Little Endian (Bedrock Edition 使用)
- Little Varint

**主要 API**:

```js
const nbt = require('prismarine-nbt')

// 解析 NBT
const buffer = fs.readFileSync('file.nbt')
const { parsed, type } = await nbt.parse(buffer)
console.log(nbt.simplify(parsed))  // 简化输出

// 写入 NBT
const data = nbt.writeUncompressed(value, 'big')
```

**Builder 用法** (构建复杂结构):
```js
const tag = nbt.comp({
  Air: nbt.short(300),
  Armor: nbt.list(nbt.comp([
    { Count: nbt.byte(0), Damage: nbt.short(0), Name: nbt.string('a') }
  ]))
})
```

**主要方法**:
- `parse()` / `parseUncompressed()` - 解析 NBT
- `writeUncompressed()` - 序列化 NBT
- `simplify()` - 简化输出 (只保留值)
- `equal()` - 比较两个 NBT 对象
- `comp()` / `list()` / `string()` / `byte()` / `short()` / `int()` / `long()` / `float()` / `double()` - 构建标签

**目录内容**:
- `nbt.js` - 核心实现
- `nbt.json` / `nbt-varint.json` - NBT 类型定义
- `typings/` - TypeScript 类型
- `sample/` - 示例文件
- `bench/` - 性能测试

---

### 5. links.md
相关文档链接集合

- [协议简介](https://zh.minecraft.wiki/w/Java%E7%89%88%E7%BD%91%E7%BB%9C%E5%8D%8F%E8%AE%AE) - 中文
- [NBT格式简介](https://zh.minecraft.wiki/w/NBT%E6%A0%BC%E5%BC%8F)
- [SNBT格式简介](https://zh.minecraft.wiki/w/SNBT%E6%A0%BC%E5%BC%8F)
- [NBT路径简介](https://zh.minecraft.wiki/w/NBT%E8%B7%AF%E5%BE%84)
- [文本组件](https://zh.minecraft.wiki/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6)

---

## 用途说明

 这些库和数据用于构建 Minecraft 客户端/服务器通信相关的功能：

 | 用途 | 使用的数据 |
 |------|------------|
 | 解析/构造数据包 | MC_Protocol_Data, minecraft-data/protocol |
 | 处理聊天消息 JSON | prismarine-chat |
 | 解析 NBT 数据 (物品、存档等) | prismarine-nbt |
 | 获取游戏数据 (方块、物品ID等) | minecraft-data |

---

## 项目文档 (docs/)

 项目根目录下 `docs/` 目录包含 Minecraft 格式规范的中文文档：

 | 文档 | 说明 |
 |------|------|
 | `docs/nbt_format.md` | NBT 二进制格式规范 |
 | `docs/snbt_format.md` | SNBT 文本格式规范 |
 | `docs/text_component.md` | 文本组件 (JSON 格式) |
 | `docs/data_components_1.21.11.md` | 1.21.11 数据组件规范 |
