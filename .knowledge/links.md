[协议简介](https://minecraft.wiki/w/Java_Edition_protocol)
[协议简介中文版](https://zh.minecraft.wiki/w/Java%E7%89%88%E7%BD%91%E7%BB%9C%E5%8D%8F%E8%AE%AE?variant=zh-cn)
[NBT格式简介](https://zh.minecraft.wiki/w/NBT%E6%A0%BC%E5%BC%8F)
[SNBT格式简介](https://zh.minecraft.wiki/w/SNBT%E6%A0%BC%E5%BC%8F)
[NBT路径简介](https://zh.minecraft.wiki/w/NBT%E8%B7%AF%E5%BE%84)
[文本组件](https://zh.minecraft.wiki/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6)
[java 1.21.11 介绍](https://zh.minecraft.wiki/w/Java%E7%89%881.21.11)
[Solt Data介绍](https://minecraft.wiki/w/Java_Edition_protocol/Slot_data)
[数据/物品组件](https://zh.minecraft.wiki/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?variant=zh)
[物品堆叠组件](https://zh.minecraft.wiki/w/Tutorial:%E7%89%A9%E5%93%81%E5%A0%86%E5%8F%A0%E7%BB%84%E4%BB%B6)
[SNBT 教程](https://zh.minecraft.wiki/w/Tutorial:SNBT)

## 协议查看器（可视化数据包结构）

[协议查看器](https://tools.minecraft.wiki/static/tools/protocol/) - 可交互查看各版本协议数据包结构

### 使用方法
1. 选择版本：点击左上角下拉菜单选择 Minecraft 版本（如 1.21.11 协议版本 774）
2. 选择数据包：点击右侧下拉菜单选择要查看的数据包

### 数据包分类
- 握手阶段 (Handshake): intention
- 状态阶段 (Status): status_request, status_response, ping_request, pong_response
- 登录阶段 (Login): hello, key, login_disconnect, login_finished, login_compression
- 配置阶段 (Configuration): registry_data, custom_payload, finish_configuration
- 游戏阶段 (Play): 大量游戏相关数据包

### 自动化交互方法（agent-browser）
该网站使用 Codex Vue 组件库，下拉菜单需要特殊处理：
```bash
# 1. 点击打开版本下拉菜单
agent-browser eval --stdin <<'EOF'
(function() {
  var selectVue = document.querySelector('.cdx-select-vue');
  if (!selectVue) return 'Select not found';
  selectVue.click();
  return 'Clicked to open';
})();
EOF

# 2. 等待菜单展开
agent-browser wait 500

# 3. 查找并点击目标选项（如 1.21.11）
agent-browser eval --stdin <<'EOF'
(function() {
  var menuItems = document.querySelectorAll('.cdx-menu__item, [role="option"]');
  for (var i = 0; i < menuItems.length; i++) {
    var text = menuItems[i].innerText || '';
    if (text.indexOf('1.21.11') !== -1 && text.indexOf('774') !== -1) {
      menuItems[i].click();
      return 'Clicked: ' + text;
    }
  }
  return 'Not found';
})();
EOF
```

### 数据源
协议数据来自 GitHub 仓库：https://github.com/Nickid2018/MC_Protocol_Data

---

## 本地文档

项目 `docs/` 目录下已有中文文档，可离线查阅：

| 文档 | 说明 |
|------|------|
| [NBT格式](docs/nbt_format.md) | NBT 二进制格式规范 |
| [SNBT格式](docs/snbt_format.md) | SNBT 文本格式规范 |
| [文本组件](docs/text_component.md) | 文本组件 JSON 格式 |
| [数据组件 1.21.11](docs/data_components_1.21.11.md) | 物品数据组件规范 |

> 注：上述链接为项目内相对路径，需要在项目根目录 `C:\Users\Yhrza\Desktop\python\gmcc\` 下使用。
