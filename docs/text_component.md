 文本组件 - 中文 Minecraft Wiki      

                             

# 文本组件

来自Minecraft Wiki

[跳转到导航](#mw-head) [跳转到搜索](#searchInput)

 **![](/images/Comment_information.svg?eab4a) Wiki上有与该主题相关的教程！**

见[教程:文本组件](/w/Tutorial:%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6 "Tutorial:文本组件")。

 **![](/images/Comment_information.svg?eab4a) Wiki上有与该主题相关的教程！**

见[教程:文本组件](/w/Tutorial:%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6 "Tutorial:文本组件")。

## 目录

-   [1 上下文与解析](#上下文与解析)
-   [2 Java版](#Java版)
    -   [2.1 基础结构](#基础结构)
    -   [2.2 组件继承](#组件继承)
    -   [2.3 预解析触发](#预解析触发)
    -   [2.4 组件类型](#组件类型)
        -   [2.4.1 纯文本组件](#纯文本组件)
        -   [2.4.2 本地化文本组件](#本地化文本组件)
        -   [2.4.3 键位绑定组件](#键位绑定组件)
        -   [2.4.4 记分板数据组件](#记分板数据组件)
        -   [2.4.5 实体名称组件](#实体名称组件)
        -   [2.4.6 NBT组件](#NBT组件)
        -   [2.4.7 精灵图组件](#精灵图组件)
            -   [2.4.7.1 纹理图集精灵图](#纹理图集精灵图)
            -   [2.4.7.2 玩家皮肤精灵图](#玩家皮肤精灵图)
    -   [2.5 组件样式](#组件样式)
        -   [2.5.1 文字颜色](#文字颜色)
        -   [2.5.2 文字字体](#文字字体)
        -   [2.5.3 文字样式](#文字样式)
        -   [2.5.4 文本插入](#文本插入)
        -   [2.5.5 点击事件](#点击事件)
        -   [2.5.6 悬停事件](#悬停事件)
-   [3 基岩版](#基岩版)
    -   [3.1 数据格式](#数据格式)
    -   [3.2 内容组件](#内容组件)
        -   [3.2.1 纯文本（Text）](#纯文本（Text）)
        -   [3.2.2 翻译文本（Translate）](#翻译文本（Translate）)
        -   [3.2.3 记分板分数（Score）](#记分板分数（Score）)
        -   [3.2.4 实体名称（Selector）](#实体名称（Selector）)
    -   [3.3 组件解析](#组件解析)
    -   [3.4 内容的定义](#内容的定义)
    -   [3.5 字符串与格式化字符串](#字符串与格式化字符串)
        -   [3.5.1 格式化字符串](#格式化字符串)
    -   [3.6 编写规范](#编写规范)
-   [4 历史](#历史)
-   [5 参见](#参见)
-   [6 注释](#注释)
-   [7 参考](#参考)
-   [8 导航](#导航)

**文本组件（Text Component）**，过去亦作**原始JSON文本（Raw JSON Text）**，Minecraft通过它向玩家发送和显示富文本。

## 上下文与解析

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=1&veaction=edit "编辑章节：上下文与解析") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=1 "编辑章节的源代码： 上下文与解析")\]

在游戏使用文本组件进行渲染时，会根据当前环境的**上下文（Context）**进行**解析（Parsing）**，成为**格式文本（Formatted Text）**，并最终使用格式文本渲染。

文本组件解析上下文包括下列数据：

-   静态上下文：当前的语言、键位等。
-   动态上下文：世界[记分板](/w/%E8%AE%B0%E5%88%86%E6%9D%BF "记分板")、实体、方块实体、[命令存储](/w/%E5%91%BD%E4%BB%A4%E5%AD%98%E5%82%A8 "命令存储")及触发预解析行为的实体、位置与朝向。

游戏解析文本组件分为两步：

1.  预解析动态组件：读取当前世界信息的组件都是动态组件，这些组件在发送到客户端前就需要在服务端使用动态上下文进行预解析成为预解析文本组件，以固定内部的动态组件内容。
2.  解析静态组件：发送到客户端的文本组件仅含静态组件，游戏客户端根据静态上下文进行最终解析，获得的格式文本仅包含文本及其样式信息，以让客户端渲染。

当游戏上下文变动时，动态组件不会跟随变化，因为这些组件已被预解析，即动态组件体现为动态上下文的快照（Snapshot）；而静态组件可以跟随静态上下文的变动而重新解析。

## Java版

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=2&veaction=edit "编辑章节：Java版") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=2 "编辑章节的源代码： Java版")\]

### 基础结构

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=3&veaction=edit "编辑章节：基础结构") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=3 "编辑章节的源代码： 基础结构")\]

文本组件有非常复杂的格式以支持各种高级用法，并在[NBT](/w/NBT "NBT")、[SNBT](/w/SNBT "SNBT")和[JSON](/w/JSON "JSON")中都有对应的格式。

-   NBT格式用于序列化与持久化，存档数据中保存与网络传输中使用的文本组件就使用了NBT格式。
-   SNBT格式用于方便输入，主要用于各种需要输入文本组件的命令、函数等。
-   JSON格式用于各种注册定义格式，例如各种数据包注册项。

下文中将主要以SNBT格式介绍；同时，JSON格式与SNBT格式相同，仅在语法上有差异，因此两者基本上可以互相转换。

文本组件主要分为三种基础结构：![字符串](/images/Data_node_string.svg?42545)字符串形式、![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)列表形式和![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)复合标签形式。游戏不会使用列表形式进行序列化，仅使用字符串形式和复合标签形式存储与传输。

复合标签格式作为文本组件的最基础格式，具有下列结构：

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597) 复合标签格式根标签
    -   ![字符串](/images/Data_node_string.svg?42545)type：文本组件类型。多数情况下此标签可以被省略，见下文[§ 组件类型](#组件类型)。
    -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)extra：此文本组件的子组件。见下文[§ 组件继承](#组件继承)。
        -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)：子文本组件。
    -   其余可选项，见[§ 组件样式](#组件样式)。
    -   与组件类型相关的项，见[§ 组件类型](#组件类型)中对应的组件格式。

字符串格式是纯文本组件的简写形式。在实际使用中，字符串格式等价于复合标签格式的`{text: <*字符串*>}`。如果一个文本组件所有样式均为默认值、不包含子组件且为纯文本组件，游戏在序列化或持久化此组件时会直接使用字符串格式而非复合标签格式。

列表格式是多个组件拼接的简写形式，列表中每个元素都是一个有效的文本组件，格式可以为复合标签、列表或字符串。由于SNBT与JSON都支持异构列表或数组，所以列表格式中的元素允许混搭，即文本组件的三种格式可以同时存在于一个列表中。游戏在读取列表格式时会自动将其转换为复合标签格式，转换遵照下列规则：

-   游戏将以第一个文本组件作为根组件，并将剩余所有组件作为子组件写入第一个文本组件的![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)extra中。
-   如果第一个文本组件的![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)extra不为空，游戏自动将剩余组件拼接在第一个文本组件![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)extra列表的后方。
-   在转换过程中，所有子组件将以原样直接写入子组件列表中，不会修改子组件的任何信息。

### 组件继承

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=4&veaction=edit "编辑章节：组件继承") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=4 "编辑章节的源代码： 组件继承")\]

文本组件以树状形式保存，以复用样式并增强结构性。

具有非空![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)extra的文本组件即为文本组件树节点。为简化说明，这里定义![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)extra内的文本组件即为此组件的**子组件**，而此组件为![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)extra内的所有文本组件的**父组件**，游戏在渲染时使用的最外层父组件称为**根组件**。

子组件会自动继承父组件内的样式信息，但不会继承组件类型及其他组件数据；如果子组件内定义了与父组件内相同的标签，则子组件定义的信息将替代继承自父组件的信息。如果继承链上没有任何组件定义了某项标签，游戏会使用当前渲染环境的样式默认值。例如，文本组件{text: 'A', color: 'red', extra: \['B', {text: 'C', color: 'yellow'}\]}中，`AB`将以红色渲染，`C`将以黄色渲染。

当文本组件解析成为格式文本时，父组件解析结果会在整个结果的最前方，子组件按照![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)extra内的顺序一个一个拼接到后方，即渲染时文本组件树将以[深度优先搜索](https://zh.wikipedia.org/wiki/%E6%B7%B1%E5%BA%A6%E4%BC%98%E5%85%88%E6%90%9C%E7%B4%A2 "wzh:深度优先搜索")的方式进行渲染。例如，文本组件{text: 'A', extra: \['B', {text: 'C', extra:\['E', 'F'\]}, {text: 'D', extra: \['G'\]}\]}将渲染为`ABCEFDG`。

根据文本组件列表格式的解析机制，第一个列表元素将成为根组件，同时也代表了第一个列表元素的样式会继承到所有后续组件中，即成为了整个文本组件的全局样式。如果使用列表格式时不想让第一个列表元素的样式控制整个文本组件，可以在第一个元素前插入空纯文本组件""以防止样式污染。

### 预解析触发

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=5&veaction=edit "编辑章节：预解析触发") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=5 "编辑章节的源代码： 预解析触发")\]

对于动态组件，游戏在发送到客户端前需要进行预解析。下列行为可以触发预解析：

-   传入文本组件的命令会进行预解析，但行为并不一致：
    -   使用`/[tellraw](/w/%E5%91%BD%E4%BB%A4/tellraw "命令/tellraw")`、`/[title](/w/%E5%91%BD%E4%BB%A4/title "命令/title")`显示文本组件时，文本组件将在发送到各个客户端前进行预解析，触发预解析的实体是将要发送到的客户端的玩家。
    -   其他所有命令，在接收文本组件参数时，都会立刻进行预解析，触发预解析的实体是执行此命令的玩家。
-   打开未预解析的[成书](/w/%E6%88%90%E4%B9%A6 "成书")（![布尔型](/images/Data_node_bool.svg?77754)resolved不存在或为`false`）。
    -   如果在[讲台](/w/%E8%AE%B2%E5%8F%B0 "讲台")上打开未解析的成书，则缺失触发预解析行为的实体，造成部分动态组件行为异常。[\[1\]](#cite_note-1)
    -   如果由玩家打开未解析的成书，则触发实体为打开此书的玩家，动态组件可正常预解析。
-   加载或设置告示牌的文本时会立刻触发预解析，但这种触发不存在触发预解析的实体。[\[2\]](#cite_note-2)
-   加载或设置[文本展示实体](/w/%E6%96%87%E6%9C%AC%E5%B1%95%E7%A4%BA%E5%AE%9E%E4%BD%93 "文本展示实体")的![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)text时将预解析内部的文本组件，触发实体为文本展示实体自身。
    -   实体**不会**在加载或设置![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)CustomName时进行预解析，此类预解析由`selector`类型的文本组件创建此文本组件副本进行预解析后使用。
-   使用[物品修饰器](/w/%E7%89%A9%E5%93%81%E4%BF%AE%E9%A5%B0%E5%99%A8 "物品修饰器")`set_lore`和`set_name`时，如果指定了![字符串](/images/Data_node_string.svg?42545)entity，且![字符串](/images/Data_node_string.svg?42545)entity所指代的实体存在，游戏将以此实体作为触发实体进行预解析。

更详细的动态上下文信息，见[命令上下文 § 文本组件解析](/w/%E5%91%BD%E4%BB%A4%E4%B8%8A%E4%B8%8B%E6%96%87#文本组件解析 "命令上下文")。

当文本组件进行预解析时，游戏将从根组件开始，遍历树中所有子组件以及样式、类型中定义的其他文本组件进行预解析。如果树的深度超过100，则深度超过100的部分将无法预解析。在下文中，如果没有明确提及，则默认所有文本组件都能被有效预解析。

[Java版26.1](/w/Java%E7%89%8826.1 "Java版26.1")起，如果在MOTD中显示文本组件[\[注 1\]](#cite_note-motd-3)，则遍历树的深度超过16的部分将替换为省略号`...`，若解析失败则强制显示为空字符串。

### 组件类型

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=6&veaction=edit "编辑章节：组件类型") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=6 "编辑章节的源代码： 组件类型")\]

文本组件具有多种类型，游戏内共有下列7种组件类型，这些类型也是![字符串](/images/Data_node_string.svg?42545)type的可选值：

-   `text`：静态组件，纯文本组件类型。
-   `translatable`：静态组件，本地化文本组件类型。做标签时名称为`translate`。
-   `keybind`：静态组件，键位绑定组件类型。
-   `score`：动态组件，记分板数据组件类型。
-   `selector`：动态组件，实体名称组件类型。
-   `nbt`：动态组件，NBT组件类型。
-   `object`：静态组件，精灵图组件类型。做标签时的名称见[§ 精灵图组件](#精灵图组件)。

当游戏解析组件时，首先会读取![字符串](/images/Data_node_string.svg?42545)type确定组件类型，并根据组件类型定义的必须项和可选项读取组件内容并进行相应解析，如果必选项不存在或组件的附加条件不满足，解析时将报错；如果![字符串](/images/Data_node_string.svg?42545)type不存在，游戏会按照上文顺序，找到第一个具有对应名称且类型对应的标签，作为此组件的组件类型；如果![字符串](/images/Data_node_string.svg?42545)type不存在或无效、且上文定义的所有标签都不存在或类型不对应，文本组件解析失败。

例如，{text: {}, keybind: 'key.inventory', translate: 'addServer.add'}在进行解析时会判定为本地化文本组件，而非纯文本组件或键位绑定组件：`text`组件要求使用字符串标签类型，而此处为复合标签，因而条件不满足而跳过解析；`translate`（即`translatable`类型）的顺序高于`keybind`，且标签类型正确对应为字符串，所以满足条件而成为此组件的类型。

换言之，![字符串](/images/Data_node_string.svg?42545)type是用于严格校验组件正确性的标签。因为游戏保持了文本组件的旧版本兼容性，此项为非必须项，且在序列化与持久化时![字符串](/images/Data_node_string.svg?42545)type不会保存。

#### 纯文本组件

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=7&veaction=edit "编辑章节：纯文本组件") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=7 "编辑章节的源代码： 纯文本组件")\]

纯文本组件直接定义了组件解析后的文本。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597) 复合标签格式根标签
    -   ![字符串](/images/Data_node_string.svg?42545)type：`text`。
    -   ![字符串](/images/Data_node_string.svg?42545)\*  
        \*text：要渲染的文字。

#### 本地化文本组件

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=8&veaction=edit "编辑章节：本地化文本组件") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=8 "编辑章节的源代码： 本地化文本组件")\]

本地化文本组件定义了本地化键名，并提供了替换参数。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597) 复合标签格式根标签
    -   ![字符串](/images/Data_node_string.svg?42545)type：`translatable`。
    -   ![字符串](/images/Data_node_string.svg?42545)\*  
        \*translate：本地化键名，游戏将使用此键名查找对应本地化文本。查询算法见下文。
    -   ![字符串](/images/Data_node_string.svg?42545)fallback：回落文本。当本地化文本无法查询得到时使用此文本作为本地化文本使用。具体算法见下文。
    -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)with：替换本地化文本内的参数。参数替换算法见下文。
        -   ![任意类型](/images/Data_node_any.svg?d406c)：一个文本组件、字符串、数字或布尔值。作为JSON格式时不可以为`null`。

作为静态组件，此组件将以原样（不包括![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)with及其递归子组件内的动态组件）发送到客户端，并由客户端在渲染前解析为格式文本渲染，因此此组件在不同的客户端上渲染可以不一致。

客户端在解析此组件时会先读取本地化键名![字符串](/images/Data_node_string.svg?42545)\*  
\*translate，从当前语言中查找此键名对应的本地化文本。如果没有找到，则寻找默认语言（英语（美国），即`en_us`）中是否有此键名的对应本地化文本。如果仍然没有找到，则尝试使用回落文本![字符串](/images/Data_node_string.svg?42545)fallback。如果![字符串](/images/Data_node_string.svg?42545)fallback没有定义，则使用键名直接作为本地化文本使用。

与UI中使用的本地化文本类似，文本组件内的本地化文本也可以使用`%%`、`%s`及其带序号的形式。但是，如果本地化文本不在语言文件中而回落使用![字符串](/images/Data_node_string.svg?42545)\*  
\*translate或者![字符串](/images/Data_node_string.svg?42545)fallback时，其内部包含的`%d`、`%f`等不支持的格式字符不会自动转换为`%s`，使得本地化文本解析失败。对于成功解析的本地化文本，游戏会计算有多少个参数需要填入，如果![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)with中定义的参数数量小于需要的参数数量，解析也会失败。无论解析因为何种原因失败，游戏都会直接使用解析过程中查找到的本地化文本作为解析结果。

对于填入的参数，如果参数不是文本组件，那么游戏事先将这些内容转换为无样式的纯文本组件。填入的参数类似文本组件子组件行为，即继承父组件（本地化文本组件）的样式，并允许自行覆盖样式。例如，本地化文本组件{translate: '%s%s', with: \[{text: 'A', color: 'red'}, 'B'\], color: 'yellow'}中`A`将渲染为红色而`B`将渲染为黄色。

#### 键位绑定组件

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=9&veaction=edit "编辑章节：键位绑定组件") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=9 "编辑章节的源代码： 键位绑定组件")\]

渲染当前玩家客户端所绑定的[键位](/w/%E9%94%AE%E4%BD%8D "键位")名称。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597) 复合标签格式根标签
    -   ![字符串](/images/Data_node_string.svg?42545)type：`keybind`。
    -   ![字符串](/images/Data_node_string.svg?42545)\*  
        \*keybind：一个[绑定键位标识符](/w/%E6%8E%A7%E5%88%B6#可设置的键位 "控制")。将显示为当前所绑定的按键名称。比如，`{keybind: "key.inventory"}`将显示"E"（默认物品栏打开按键）。若找不到相应按键标识符将尝试显示为对应的翻译名称。

与本地化文本组件类似，键位绑定组件也是静态组件，将以原样发送到客户端，并由客户端在渲染前解析为格式文本渲染，因此此组件在不同的客户端上渲染可以不一致。

#### 记分板数据组件

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=10&veaction=edit "编辑章节：记分板数据组件") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=10 "编辑章节的源代码： 记分板数据组件")\]

记分板数据组件可以读取服务端内的记分板数据，并在服务端预解析成为静态组件后发送到客户端。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597) 复合标签格式根标签
    -   ![字符串](/images/Data_node_string.svg?42545)type：`score`。
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)\*  
        \*score：要获取的分数信息。
        -   ![字符串](/images/Data_node_string.svg?42545)\*  
            \*name：分数持有者信息。此项的详细格式见下文。
        -   ![字符串](/images/Data_node_string.svg?42545)\*  
            \*objective：要获取的记分项信息。

游戏在预解析此组件时，会先计算分数持有者信息。![字符串](/images/Data_node_string.svg?42545)\*  
\*name可以为下列3种格式，分别对应了不同的计算方式：

-   有效的[目标选择器](/w/%E7%9B%AE%E6%A0%87%E9%80%89%E6%8B%A9%E5%99%A8 "目标选择器")：游戏会根据目标选择器在世界中查找符合条件的实体，将查找到的唯一实体作为分数持有者。如果查找到的实体不止一个，预解析将直接报错（只允许一个实体，但提供的选择器允许多个实体）；如果没有查找到实体，则认为此目标选择器的文本本身为玩家名称或UUID，进行第二次预解析。
-   玩家名称或UUID：游戏直接将对应实体作为分数持有者。
-   分数持有者通配符`*`：游戏将使用当前动态上下文内触发预解析行为的实体作为分数持有者。

如果分数持有者为空，或分数持有者的对应记分项信息不存在，则游戏预解析后此项直接转变为空纯文本组件；如果不为空，游戏先根据记分项返回的分数使用记分项对应的数字格式进行格式化。

如果此动态组件未被成功预解析，在客户端渲染此组件时将直接渲染为空文本。

#### 实体名称组件

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=11&veaction=edit "编辑章节：实体名称组件") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=11 "编辑章节的源代码： 实体名称组件")\]

实体名称组件可以读取世界内的实体数据，并在预解析时处理为对应的静态组件发送到客户端。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597) 复合标签格式根标签
    -   ![字符串](/images/Data_node_string.svg?42545)type：`selector`。
    -   ![字符串](/images/Data_node_string.svg?42545)\*  
        \*selector：有效的目标选择器、玩家名称或UUID。
    -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)separator：多个实体名称间的分隔符。

此组件进行预解析时，游戏会先预解析![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)separator内定义的名称分隔符；如果![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)separator不存在，则游戏默认使用灰色的逗号（{text: ', ', color: 'gray'}）作为名称分隔符。

在预解析名称分隔符后，游戏根据目标选择器、玩家名称或UUID在世界中查找匹配的实体，并获取所有匹配实体的显示名称。其中，显示名称文本组件的获取算法如下：

1.  获取实体未格式化的名称：
    -   如果实体是玩家，获取玩家名称。玩家名称为纯文本组件，内部为玩家档案名称。
    -   如果实体是带有自定义名称的实体，获取![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)CustomName。
        -   如果![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)CustomName内定义了点击事件样式，则抹除此点击事件样式，其他样式将保留。[\[3\]](#cite_note-4)
    -   如果实体不满足上述两个条件，则使用实体类型名称作为实体名称，即`{translate: 'entity.<*实体类型命名空间ID*>'}`。
2.  添加[队伍](/w/%E9%98%9F%E4%BC%8D "队伍")样式，如果实体不在任何队伍中则不进行此步骤。
    1.  为名称添加前缀（![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)MemberNamePrefix）和尾缀（![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)MemberNameSuffix）。此步骤中，游戏会自动创建一个新的空纯文本组件，并将前缀、名称、尾缀作为此组件的子组件放入。
    2.  为上一步修饰好的名称添加样式，由队伍中的![字符串](/images/Data_node_string.svg?42545)TeamColor决定。如果此值为`reset`（重置样式），则不进行此步骤。
3.  添加悬浮事件样式，此样式为`show_entity`类型，并添加此实体的类型、UUID和未格式化的名称数据。
4.  添加插入事件样式，此样式信息为实体的带分隔符形式的UUID。
5.  如果此实体为玩家，继续添加点击事件样式，此样式为`suggest_command`类型，补全信息为`/tell <*玩家档案名称*>` 。

获取显示名称文本组件列表后，游戏将进行预解析的最后步骤，拼接出静态组件：

-   如果列表为空，即未找到匹配实体，则返回空纯文本组件。
-   如果列表内只有一个元素，则直接返回此元素。
-   如果列表内有多个元素，游戏会先创建一个空纯文本组件，并按照列表顺序取出元素作为子组件加入到纯文本组件中，同时两个元素中间会插入名称分隔符文本组件进行分隔。

如果此动态组件未被成功预解析，在客户端渲染此组件时将直接渲染为![字符串](/images/Data_node_string.svg?42545)\*  
\*selector。

#### NBT组件

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=12&veaction=edit "编辑章节：NBT组件") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=12 "编辑章节的源代码： NBT组件")\]

NBT组件能获取游戏内方块实体、实体及命令存储中的信息，并预解析为对应的静态组件发送到客户端。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597) 复合标签格式根标签
    -   ![字符串](/images/Data_node_string.svg?42545)type：`nbt`。
    -   ![字符串](/images/Data_node_string.svg?42545)\*  
        \*nbt：[NBT路径](/w/NBT%E8%B7%AF%E5%BE%84 "NBT路径")。详细使用见下文。
    -   ![布尔型](/images/Data_node_bool.svg?77754)interpret：（默认为`false`，[Java版26.1](/w/Java%E7%89%8826.1 "Java版26.1")起，![布尔型](/images/Data_node_bool.svg?77754)plain为`true`时必须为`false`）是否将获取到的NBT数据解析为文本组件。[Java版26.1](/w/Java%E7%89%8826.1 "Java版26.1")起，此值为`false`时，游戏会将获取到的NBT数据整理为与`/[data](/w/%E5%91%BD%E4%BB%A4/data "命令/data") get`一致的经过语法高亮的[SNBT](/w/SNBT "SNBT")。
    -   ![布尔型](/images/Data_node_bool.svg?77754)plain\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]：（默认为`false`，![布尔型](/images/Data_node_bool.svg?77754)interpret为`true`时必须为`false`）![布尔型](/images/Data_node_bool.svg?77754)interpret为`false`时，游戏是否将获取到的数据输出为简单的单一文本，而非经过语法高亮处理的富文本。
    -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)separator：多个NBT数据展示时使用的分隔符。
    -   ![字符串](/images/Data_node_string.svg?42545)source：NBT数据源类型。详情见下文。
    
    -   下列标签必须至少包含一个，作为NBT数据源使用：
    -   ![字符串](/images/Data_node_string.svg?42545)entity：将实体数据作为NBT数据源使用。此标签需要为有效的目标选择器、玩家名称或UUID。
    -   ![字符串](/images/Data_node_string.svg?42545)block：将指定位置的方块实体数据作为NBT数据源使用。此标签需要为有效的[方块位置参数](/w/%E5%9D%90%E6%A0%87#命令 "坐标")，相对坐标与局部坐标解析时使用触发预解析行为时的坐标和朝向。
    -   ![字符串](/images/Data_node_string.svg?42545)storage：将指定命名空间ID的命令存储作为NBT数据源使用。

![字符串](/images/Data_node_string.svg?42545)source与文本组件中的基础标签![字符串](/images/Data_node_string.svg?42545)type类似，作用是用于检验数据结构的正确性，可选值为`entity`、​`block`和​`storage`，分别对应字符串类型标签。当![字符串](/images/Data_node_string.svg?42545)source不存在时，游戏按照`entity`、`block`、`storage`的顺序尝试读取。如果这三项都不存在，此文本组件解析无效。

这三种数据源能获取的NBT数据如下：

-   ![字符串](/images/Data_node_string.svg?42545)block获取方块实体的数据。数据包括方块实体的位置、ID、数据组件信息，详情见[方块实体数据格式](/w/%E6%96%B9%E5%9D%97%E5%AE%9E%E4%BD%93%E6%95%B0%E6%8D%AE%E6%A0%BC%E5%BC%8F "方块实体数据格式")。此数据与`/[data](/w/%E5%91%BD%E4%BB%A4/data "命令/data") block`获取的数据一致。
-   ![字符串](/images/Data_node_string.svg?42545)entity获取一系列实体的数据。数据包括实体的所有数据，但不包括实体ID。如果实体为玩家且选中物品非空，则还额外包括![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)SelectedItem数据。此数据与`/[data](/w/%E5%91%BD%E4%BB%A4/data "命令/data") entity`获取的数据一致。
-   ![字符串](/images/Data_node_string.svg?42545)storage获取命令存储中的数据。此数据与`/[data](/w/%E5%91%BD%E4%BB%A4/data "命令/data") storage`获取的数据一致。

如果数据源找不到相应数据（坐标位置上不存在方块实体、目标选择器未选中到实体，指定命令存储不存在），则直接预解析为空纯文本组件；如果具有数据，则将所有数据放入数组中（实体可能具有多个数据，方块和命令存储仅可能有一个数据），并根据指定的NBT路径选中所有元素进行映射后进行扁平化（Flatten Mapping）。

例如，如果一个实体数据源`@e`获取到了两个实体，并且位置数据![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)Motion分别为\[1d, 0d, \-1d\]和\[\-2d, 0d, 2d\]：

-   如果NBT路径为`Motion`，则获取到的数据为\[1d, 0d, \-1d\]和\[\-2d, 0d, 2d\]共2项数据。
-   如果NBT路径为`Motion[]`，则获取到的数据为1d、0d、\-1d、\-2d、0d、2d共6项数据。
-   如果NBT路径为`Motion[0]`，则获取到的数据为1d和\-2d共2项数据。

获取到数据数组后，根据![布尔型](/images/Data_node_bool.svg?77754)interpret的取值，游戏会进行不同的预解析步骤：

-   如果![布尔型](/images/Data_node_bool.svg?77754)interpret为`false`，游戏先将所有数据转换为SNBT字符串。[Java版26.1](/w/Java%E7%89%8826.1 "Java版26.1")前，如果![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)separator不存在，游戏将直接使用字符串`,` 拼接字符串，并合并为一个纯文本组件；如果存在，则将第一个数据SNBT字符串作为根组件，![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)separator和剩余组件依次放入根组件作为子组件拼接为最终的文本组件。[Java版26.1](/w/Java%E7%89%8826.1 "Java版26.1")起，如果日志级别为调试级别，则游戏会对标签按标签名进行排序整理。若![布尔型](/images/Data_node_bool.svg?77754)plain为`true`，则直接拼接为简单的单一字符串进行输出，否则游戏会将各字符串重新整理为更复杂的文本组件，使其按照SNBT的语法结构渲染。各SNBT字符串之间的拼接字符由![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)separator决定，不存在时游戏将使用无样式的逗号`,` 拼接字符串。
-   如果![布尔型](/images/Data_node_bool.svg?77754)interpret为`true`，游戏先将所有数据转换为文本组件，如果不能转换，则直接删除对应数据。与为`false`时不同，如果![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)separator不存在游戏将以无样式的逗号（', '）作为默认分隔符。数据的第一个文本组件将作为根组件，![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)separator和剩余组件依次放入根组件作为子组件拼接为最终的文本组件。这种拼接方式代表了数据的第一个文本组件将控制影响所有后续文本组件的样式数据。

例如，如果一个实体数据源`@e`获取到了三个实体，并且![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)data分别为{color: 'red', text: 'A'}、{test: 'B'}和{text: 'C'}，且NBT路径为`data`：

-   如果![布尔型](/images/Data_node_bool.svg?77754)interpret为`false`，则预解析后的文本组件为'{color:"red",text:"A"}, {test:"B"}, {text:"C"}'。
-   如果![布尔型](/images/Data_node_bool.svg?77754)interpret为`true`，则预解析后的文本组件为{color: 'red', text: 'A', extra: \[', ', 'C'\]}，显示为红色的`A, C`。

如果此动态组件未被成功预解析，在客户端渲染此组件时将直接渲染为空文本。

#### 精灵图组件

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=13&veaction=edit "编辑章节：精灵图组件") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=13 "编辑章节的源代码： 精灵图组件")\]

精灵图组件，亦称**对象组件**，可以指定游戏客户端渲染精灵图。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597) 复合标签格式根标签
    -   ![字符串](/images/Data_node_string.svg?42545)type：`object`。
    -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)fallback\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]：当精灵图组件渲染失败时，作为回落的文本组件。
    -   ![字符串](/images/Data_node_string.svg?42545)object：精灵图类型，取值见下文。
    -   剩余标签与精灵图类型相关。

精灵图组件自身具有多种类型，目前游戏内的精灵图类型如下：

-   `atlas`：渲染[纹理图集](/w/%E7%BA%B9%E7%90%86%E5%9B%BE%E9%9B%86 "纹理图集")中的精灵图。做标签时名称为`sprite`。
-   `player`：渲染玩家的头的正面图像。

与文本组件解析组件类型类似，当游戏需要解析精灵图组件时，会先读取![字符串](/images/Data_node_string.svg?42545)object字段确定精灵图类型，使用此精灵图类型进行解析；如果未指定![字符串](/images/Data_node_string.svg?42545)object字段，则自动按照上述顺序进行解析；如果![字符串](/images/Data_node_string.svg?42545)object未指定或且上述标签均不存在，则此精灵图组件无效。

游戏在渲染精灵图组件时，将自动替换精灵图组件所在位置的字符为`U+FFFC`（Object Replacement Character），并将此文字的渲染替换为对应要渲染的精灵图。

精灵图渲染时，游戏会自动将原始的精灵图转换为字体度量中的8×8像素大小，而非将整个精灵图缩放到8×8像素渲染。

组件样式中强行指定精灵图组件的字体无效，粗体、斜体和随机字符样式会被忽略。

##### 纹理图集精灵图

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=14&veaction=edit "编辑章节：纹理图集精灵图") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=14 "编辑章节的源代码： 纹理图集精灵图")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597) 复合标签格式根标签
    -   ![字符串](/images/Data_node_string.svg?42545)type：`object`。
    -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)fallback\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]：（默认为`[<*精灵图命名空间ID*>]`（默认值）或`[<*精灵图命名空间ID*>@<*纹理图集命名空间ID*>]`）当纹理图集精灵图组件渲染失败时，作为回落的文本组件。
    -   ![字符串](/images/Data_node_string.svg?42545)object：`atlas`。
    -   ![字符串](/images/Data_node_string.svg?42545)atlas：（默认为`blocks`）使用的纹理图集。
    -   ![字符串](/images/Data_node_string.svg?42545)\*  
        \*sprite：要渲染的精灵图在指定纹理图集中的命名空间ID。

##### 玩家皮肤精灵图

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=15&veaction=edit "编辑章节：玩家皮肤精灵图") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=15 "编辑章节的源代码： 玩家皮肤精灵图")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597) 复合标签格式根标签
    -   ![字符串](/images/Data_node_string.svg?42545)type：`object`。
    -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)fallback\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]：（默认为`[<*玩家名称*> head]`（存在玩家名称时）或`[unknown player head]`）当玩家皮肤精灵图组件渲染失败时，作为回落的文本组件。
    -   ![字符串](/images/Data_node_string.svg?42545)object：`player`。
    -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)\*  
        \*player：要使用的玩家皮肤。设置为字符串时相当于直接设置![字符串](/images/Data_node_string.svg?42545)name。
        
        -   游戏档案，见[Template:Nbt inherit/resolvable profile/source](/w/Template:Nbt_inherit/resolvable_profile/source "Template:Nbt inherit/resolvable profile/source")
    -   ![布尔型](/images/Data_node_bool.svg?77754)hat：（默认为`true`）是否渲染皮肤的“帽子”部分。

在MOTD中使用玩家皮肤精灵图会被强制渲染为回落文本。[\[注 1\]](#cite_note-motd-3)

### 组件样式

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=16&veaction=edit "编辑章节：组件样式") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=16 "编辑章节的源代码： 组件样式")\]

文本组件可以附加样式，以修改文本组件渲染时的效果，或给文本增加相应事件。

#### 文字颜色

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=17&veaction=edit "编辑章节：文字颜色") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=17 "编辑章节的源代码： 文字颜色")\]

文本组件可以指定文字渲染的颜色和阴影颜色，标签如下：

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597) 复合标签格式根标签
    -   ![字符串](/images/Data_node_string.svg?42545)color：文本渲染颜色。
    -   ![整型](/images/Data_node_int.svg?8d24f)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)shadow\_color：文本的阴影颜色。游戏只以整数形式保存。
        
        -   ARGB颜色，见[Template:Nbt inherit/argb color/source](/w/Template:Nbt_inherit/argb_color/source "Template:Nbt inherit/argb color/source")

文本渲染颜色![字符串](/images/Data_node_string.svg?42545)color可以使用两种格式：

-   直接使用16进制颜色，此时字符串需要以`#`开头，后方为16进制数字。其中16进制数字不能大于0xFFFFFF（16777215），也不能小于0，不包含A（透明度）通道。例如`#FF0000`。
-   使用颜色代码，有效的颜色代码及相应的颜色见[格式化代码 § 颜色代码](/w/%E6%A0%BC%E5%BC%8F%E5%8C%96%E4%BB%A3%E7%A0%81#颜色代码 "格式化代码")，使用其名称作为有效值。例如`yellow`。

文本阴影颜色![整型](/images/Data_node_int.svg?8d24f)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)shadow\_color与渲染颜色不同，可以指定A通道。由于着色器限制，当A通道小于0.1（整数形式为0x1A）时，文本阴影无法渲染。

#### 文字字体

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=18&veaction=edit "编辑章节：文字字体") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=18 "编辑章节的源代码： 文字字体")\]

[![](/images/thumb/Different_Fonts_in_Minecraft.png/300px-Different_Fonts_in_Minecraft.png?27f63)](/w/File:Different_Fonts_in_Minecraft.png)

使用游戏中存在的4种字体显示文本“Minecraft Wiki”，从上到下依次是预设字体、Unicode字体、标准银河字母和“illageralt”字体

文字字体也可以使用文本组件样式设置：

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597) 复合标签格式根标签
    -   ![字符串](/images/Data_node_string.svg?42545)font：命名空间ID，指定渲染使用的字体。默认使用`minecraft:default`。

文字可以使用的[字体](/w/%E8%87%AA%E5%AE%9A%E4%B9%89%E5%AD%97%E4%BD%93 "自定义字体")由[资源包](/w/%E8%B5%84%E6%BA%90%E5%8C%85 "资源包")定义在`assets/<*命名空间*>/font/<*路径*>.json`中。如果字体不存在，则游戏直接将所有字符替换为缺失的字形渲染。

#### 文字样式

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=19&veaction=edit "编辑章节：文字样式") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=19 "编辑章节的源代码： 文字样式")\]

文本组件也可以设置文字的显示样式和装饰：

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597) 复合标签格式根标签
    -   ![布尔型](/images/Data_node_bool.svg?77754)bold：将文字渲染为粗体。
    -   ![布尔型](/images/Data_node_bool.svg?77754)italic：将文字渲染为斜体。
    -   ![布尔型](/images/Data_node_bool.svg?77754)underlined：为文字渲染添加下划线。
    -   ![布尔型](/images/Data_node_bool.svg?77754)strikethrough：为文字渲染添加删除线。
    -   ![布尔型](/images/Data_node_bool.svg?77754)obfuscated：渲染为随机字符。

各种文字样式的渲染效果如下：

标签

渲染效果

渲染原理（像素以未缩放计）

![布尔型](/images/Data_node_bool.svg?77754)bold

hello

文本重复渲染两次，两次渲染中存在渲染偏移

![布尔型](/images/Data_node_bool.svg?77754)italic

hello

修改文字顶部顶点偏移，倾斜角度为14°（arctan(0.25)）

![布尔型](/images/Data_node_bool.svg?77754)underlined

hello

在字身框顶下方8像素处渲染1像素的线

![布尔型](/images/Data_node_bool.svg?77754)strikethrough

hello

在字身框顶下方3.5像素处渲染1像素的线

![布尔型](/images/Data_node_bool.svg?77754)obfuscated

[![](/images/thumb/TextObfuscated.gif/200px-TextObfuscated.gif?59f3d)](/w/File:TextObfuscated.gif)

从字体中挑选**相同宽度**的字符替代渲染原字符

#### 文本插入

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=20&veaction=edit "编辑章节：文本插入") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=20 "编辑章节的源代码： 文本插入")\]

文本组件可以指定被玩家按住⇧ Shift并进行有效点击时插入的文本。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597) 复合标签格式根标签
    -   ![字符串](/images/Data_node_string.svg?42545)insertion：当按住⇧ Shift并进行有效点击时插入文本。

文本插入只在聊天屏幕中生效，点击时将替换聊天栏内选中文本，如果未选中文本则在光标处插入指定文本。

#### 点击事件

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=21&veaction=edit "编辑章节：点击事件") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=21 "编辑章节的源代码： 点击事件")\]

在不同位置渲染的文本组件可以执行点击事件，当鼠标指针点击文本后即可执行对应行为。虽然文本插入也需要玩家点击但游戏内部不认为其属于点击事件。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597) 复合标签格式根标签
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)click\_event：定义点击事件，有效具体行为及条件见下文。
        -   ![字符串](/images/Data_node_string.svg?42545)\*  
            \*action：点击后的执行行为。
        
        -   如果![字符串](/images/Data_node_string.svg?42545)\*  
            \*action为`change_page`：
        -   ![整型](/images/Data_node_int.svg?8d24f)\*  
            \*page：（值>0）要跳转到的书页。
        
        -   如果![字符串](/images/Data_node_string.svg?42545)\*  
            \*action为`copy_to_clipboard`：
        -   ![字符串](/images/Data_node_string.svg?42545)\*  
            \*value：要复制到剪贴板的字符串。
        
        -   如果![字符串](/images/Data_node_string.svg?42545)\*  
            \*action为`custom`：
        -   ![字符串](/images/Data_node_string.svg?42545)\*  
            \*id：要发送的自定义网络负载命名空间ID。
        -   ![任意类型](/images/Data_node_any.svg?d406c)payload：要发送的自定义网络负载。嵌套不允许超过16层，且序列化后长度不能超过32768字节。
        
        -   如果![字符串](/images/Data_node_string.svg?42545)\*  
            \*action为`open_file`：
        -   ![字符串](/images/Data_node_string.svg?42545)\*path：要打开的文件。
        
        -   如果![字符串](/images/Data_node_string.svg?42545)\*  
            \*action为`open_url`：
        -   ![字符串](/images/Data_node_string.svg?42545)\*  
            \*url：要打开的URL。URL的协议（Scheme）必须为`http`或`https`，其他协议无法解析。
        
        -   如果![字符串](/images/Data_node_string.svg?42545)\*  
            \*action为`run_command`：
        -   ![字符串](/images/Data_node_string.svg?42545)\*  
            \*command：要发送到服务端执行的命令，可以不带前导正斜杠（`/`）。此字符串中不允许出现`\u00a7`（分节符）、`\u007f`（DEL）和小于`\u0020`的字符。
        
        -   如果![字符串](/images/Data_node_string.svg?42545)\*  
            \*action为`show_dialog`：
        -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)\*  
            \*dialog：要打开的[对话框](/w/%E5%AF%B9%E8%AF%9D%E6%A1%86 "对话框")。可以指定对话框的命名空间ID，也可以内联定义。
        
        -   如果![字符串](/images/Data_node_string.svg?42545)\*  
            \*action为`suggest_command`：
        -   ![字符串](/images/Data_node_string.svg?42545)\*  
            \*command：要替换聊天栏中的命令。此字符串中不允许出现`\u00a7`（分节符）、`\u007f`（DEL）和小于`\u0020`的字符。

文本组件的有效点击事件只有在下列情况下才可以触发：

-   点击死亡屏幕中的死亡消息文本组件，且点击事件类型必须为`open_url`。
-   点击告示牌，要求点击事件必须在根组件上，且类型必须为`run_command`、`custom`或`show_dialog`。
-   点击聊天屏幕中的任何文本组件（不包括悬停文本框）。
-   成书书预览屏幕中的书内部的文本组件（不包括悬停文本框）。

下列是各种点击事件的详细行为：

-   **`change_page`**必须在成书的预览屏幕中才能生效，点击后游戏会将成书翻开到指定页数。如果页数超出成书页码范围，则跳转到第一页或最后一页。
-   **`copy_to_clipboard`**在聊天屏幕和成书预览屏幕中生效，将指定文本写入系统剪贴板中。
-   **`custom`**在聊天屏幕、成书预览屏幕和告示牌中生效，点击后客户端将向服务端发送指定命名空间ID和负载的`custom_click_action`网络数据包。
    -   预留给自定义服务端，原版服务端只会以调试等级日志输出一条信息`Received custom click action <*网络负载ID*> with payload <*网络负载*>`。
-   **`open_file`**在聊天屏幕和成书预览屏幕中生效，且此点击事件类型无法被序列化与反序列化，仅用于客户端内部使用。此事件用于打开指定文件，在不同系统下，此点击事件行为不同：
    -   Windows系统下，游戏实际上调用`rundll32 url.dll,FileProtocolHandler file:<*文件路径*>`。
    -   macOS下，游戏实际上调用`open file:<*文件路径*>`。
    -   其他操作系统下（例如Linux），游戏实际上调用`xdg-open file://<*文件路径*>`。
-   **`open_url`**在死亡屏幕、聊天屏幕和成书预览屏幕中生效，用于打开指定的URL。
    -   如果[客户端选项文件](/w/%E5%AE%A2%E6%88%B7%E7%AB%AF%E9%80%89%E9%A1%B9%E6%96%87%E4%BB%B6%E6%A0%BC%E5%BC%8F "客户端选项文件格式")中的`chatLinks`为false，则此点击事件不生效。
    -   如果客户端选项文件中的`chatLinksPrompt`为true，则在点击后弹出询问是否打开此URL的屏幕，确认后进行跳转；如果为false，则直接跳转。
    -   不同系统下此点击事件行为不一致，与打开文件类似使用相同命令（`rundll32/open/xdg-open`），区别仅在最后参数的URL。
-   **`run_command`**在告示牌、聊天屏幕和成书预览屏幕中生效，用于执行指定命令。
    -   在告示牌内点击生效时，由服务端处理行为，命令上下文权限等级与触发者无关。详细上下文见[命令上下文 § 告示牌](/w/%E5%91%BD%E4%BB%A4%E4%B8%8A%E4%B8%8B%E6%96%87#告示牌 "命令上下文")。
    -   在聊天屏幕和成书预览屏幕生效时，由客户端处理行为，命令上下文与客户端在聊天栏内执行命令一致。
    -   在网络阶段为配置（Configuration）阶段中的对话框内此事件不生效。
-   **`show_dialog`**在聊天屏幕、成书预览屏幕和告示牌中生效，点击后将打开指定的对话框。
-   **`suggest_command`**在聊天屏幕中生效，点击时将直接替代当前聊天栏内的信息。

#### 悬停事件

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=22&veaction=edit "编辑章节：悬停事件") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=22 "编辑章节的源代码： 悬停事件")\]

与点击事件类似，当鼠标指针悬停在渲染的文本组件区域内时，会触发悬停事件并使游戏展示指定的数据。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597) 复合标签格式根标签
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)hover\_event：定义悬停事件。
        -   ![字符串](/images/Data_node_string.svg?42545)\*  
            \*action：悬停事件的类型，控制悬停渲染的行为。
        
        -   如果![字符串](/images/Data_node_string.svg?42545)\*  
            \*action为`show_entity`：
        -   ![字符串](/images/Data_node_string.svg?42545)\*  
            \*id：（命名空间ID）实体类型。
        -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)name：（文本组件）实体的显示名称。**此文本组件无法被预解析**。
        -   ![字符串](/images/Data_node_string.svg?42545)![整型数组](/images/Data_node_int-array.svg?546e8)\*  
            \*uuid：实体的UUID。可以为带有连字符的UUID字符串，也可以为4个整数组成UUID。
        
        -   如果![字符串](/images/Data_node_string.svg?42545)\*  
            \*action为`show_item`：
        -   ![字符串](/images/Data_node_string.svg?42545)\*  
            \*id：（[命名空间ID](/w/%E5%91%BD%E5%90%8D%E7%A9%BA%E9%97%B4ID "命名空间ID")）表示某种类的物品堆叠。若未指定，游戏会在加载区块或者生成物品时将其变更为空气。
        -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components：当前物品的组件修订，将修改物品的[数据组件](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6 "数据组件")信息。
            -   ![任意类型](/images/Data_node_any.svg?d406c)<*数据组件ID*\>：一项组件和其对应的数据，代表物品拥有此组件。设置组件数据时可以不写命名空间，但游戏在导出时会自行加上`minecraft:`前缀。
            -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)!<*数据组件ID*\>：存在时，使一个数据组件失效。此复合标签的内容不影响行为。设置组件数据时可以不写命名空间，但游戏在导出时会自行加上`minecraft:`前缀。
        -   ![整型](/images/Data_node_int.svg?8d24f)count：（0<值≤物品最大堆叠数量）[物品](/w/%E7%89%A9%E5%93%81 "物品")的堆叠数。不存在或无效时则默认为1。
        
        -   如果![字符串](/images/Data_node_string.svg?42545)\*  
            \*action为`show_text`：
        -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)\*  
            \*value：（文本组件）要显示的文本组件。[\[4\]](#cite_note-5)

文本组件的悬停事件在下列情况下才可触发：

-   死亡屏幕中的死亡消息文本组件。
-   聊天屏幕中的所有文本组件（不包括悬停文本框）。
-   成书书预览屏幕中的书内部的文本组件（不包括悬停文本框）。

下列是各种悬停事件的详细行为：

-   **`show_entity`**将展示实体的信息，包括实体类型、名称和UUID。此悬停事件必须在客户端选项文件中`advancedItemTooltips`（高级提示框，可使用F3 + H切换）为true时才能生效。
    -   游戏将按照三行显示数据，分别为实体的名称、类型和UUID。如果悬停事件中未定义实体名称，则名称行不存在，只渲染类型和UUID两行。
    -   `show_entity`中![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)name无法被预解析。例如{text: 'test', hover\_event: {action: 'show\_entity', uuid: \[I; 0, 0, 0, 0\], id: 'player', name: {selector: '@e'}}}在悬停时实体名称直接显示为`@e`而不是预解析后的具体实体名称列表。
-   **`show_item`**将展示物品的信息，渲染结果与物品栏内物品悬浮提示框一致。此渲染也受到客户端选项文件中`advancedItemTooltips`值的影响。
-   **`show_text`**将直接展示指定的文本组件。

## 基岩版

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=23&veaction=edit "编辑章节：基岩版") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=23 "编辑章节的源代码： 基岩版")\]

[![](/images/thumb/Knowledge_Book_JE2.png/16px-Knowledge_Book_JE2.png?632b9)](/w/File:Knowledge_Book_JE2.png)

**此章节缺失以下信息：需要更多基岩版文本组件信息**

请协助补充相关内容的描述，[讨论页](/w/Talk:%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6 "Talk:文本组件")可能有更多细节。

[基岩版](/w/%E5%9F%BA%E5%B2%A9%E7%89%88 "基岩版")有着较为精简的文本组件格式，它以文本显示为主，目前为止不具备任何交互功能。文本的字体样式可通过[格式化代码](/w/%E6%A0%BC%E5%BC%8F%E5%8C%96%E4%BB%A3%E7%A0%81 "格式化代码")进行修饰。

在[基岩版](/w/%E5%9F%BA%E5%B2%A9%E7%89%88 "基岩版")中，文本组件被应用在`/[tellraw](/w/%E5%91%BD%E4%BB%A4/tellraw "命令/tellraw")`的内容、`/[titleraw](/w/%E5%91%BD%E4%BB%A4/titleraw "命令/titleraw")`的标题、[NPC](/w/NPC "NPC")的名字、[书与笔](/w/%E4%B9%A6%E4%B8%8E%E7%AC%94 "书与笔")的文本（题目和作者除外）、[告示牌](/w/%E5%91%8A%E7%A4%BA%E7%89%8C "告示牌")的文本以及大部分富文本信息中。但记分板分数组件（Score）和实体名称组件（Selector）只能在流动的文本信息中生效，即聊天室讯息和屏幕标题，否则将不显示内容。

### 数据格式

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=24&veaction=edit "编辑章节：数据格式") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=24 "编辑章节的源代码： 数据格式")\]

基岩版的文本组件的根节点只能是![字符串](/images/Data_node_string.svg?42545)纯文本或具有![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)rawtext属性的![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)对象。玩家在游戏中一般只能使用到![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)文本组件。

-   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597) 文本组件根节点。
    
    -   如果根节点为![字符串](/images/Data_node_string.svg?42545)：纯文本形式，一般只能在文件代码中碰见。
    
    -   如果根节点为![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)：对象形式的文本组件。
    -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)rawtext：必须包含一个或多个内容组件。若为空列表将出错。
        -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)：一个内容组件。“没有内容”的![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)内容组件将被忽略。

### 内容组件

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=25&veaction=edit "编辑章节：内容组件") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=25 "编辑章节的源代码： 内容组件")\]

内容组件是![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)文本组件的一部分，无法独立存在，用于定义一个![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)文本组件所显示的内容。一个内容组件只能显示一种内容类型和一个扩展的内容组件列表![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)rawtext，内容组件内容的类型由对象内出现的必要字段决定。各类内容组件内容对应的必要字段如下：

-   纯文本组件内容：![字符串](/images/Data_node_string.svg?42545)text
-   翻译文本组件内容：![字符串](/images/Data_node_string.svg?42545)translate
-   记分板分数组件内容：![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)score
-   实体名称组件内容：![字符串](/images/Data_node_string.svg?42545)selector

若内容组件中没有任何必要字段和扩展的内容组件列表，则[不显示任何内容](#内容的定义)。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)：内容组件根节点。
    -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)rawtext：（可选）扩展的内容组件列表，将始终显示在内容组件的基础内容之后。可以包含零个或多个内容组件，若为空则不显示内容。
        -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)：一个内容组件。
    -   根据内容类型而指定的额外字段。详见下方后续。

#### 纯文本（Text）

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=26&veaction=edit "编辑章节：纯文本（Text）") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=26 "编辑章节的源代码： 纯文本（Text）")\]

一段纯文本内容。

纯文本内容类型的内容组件有以下字段：

-   -   ![字符串](/images/Data_node_string.svg?42545)text：值是允许转义字符的字符串。

示例：

这是最简单的纯文本组件。会显示hello。

{"rawtext": \[{"text": "hello"}\]}

hello

文本组件中的换行必须通过`\n`完成。

{"rawtext": \[{"text": "第一行\\n第二行\\n\\n第四行"}\]}

第一行  
第二行  
  
第四行

多个内容组件将依序显示。

{"rawtext": \[{"text": "hello "}, {"text": "world"}\]}

hello world

#### 翻译文本（Translate）

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=27&veaction=edit "编辑章节：翻译文本（Translate）") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=27 "编辑章节的源代码： 翻译文本（Translate）")\]

参见：[Tutorial:自定义附加包语言文件](/w/Tutorial:%E8%87%AA%E5%AE%9A%E4%B9%89%E9%99%84%E5%8A%A0%E5%8C%85%E8%AF%AD%E8%A8%80%E6%96%87%E4%BB%B6 "Tutorial:自定义附加包语言文件")

一段可根据所选[语言](/w/%E8%AF%AD%E8%A8%80 "语言")翻译成译文的本地化键名或一段[格式化字符串](#格式化字符串)（Format String）。翻译文本会被转译成格式化字符串，可以包含纯文本的格式说明符（Format Specifier），使文本内容可变。

翻译文本内容类型的内容组件有以下字段：

-   -   ![字符串](/images/Data_node_string.svg?42545)translate：一个本地化键名，它将会以玩家所选语言显示对应的[语言文件](/w/Tutorial:%E8%87%AA%E5%AE%9A%E4%B9%89%E9%99%84%E5%8A%A0%E5%8C%85%E8%AF%AD%E8%A8%80%E6%96%87%E4%BB%B6 "Tutorial:自定义附加包语言文件")下相应的译文。若未找到当前语言下有该本地化键名，则默认检查`en_us.lang`中是否有对应的译文。\[[需要在基岩版上验证](/w/Special:TalkPage/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6 "Special:TalkPage/文本组件")\]若未找到，则该值将被视为允许格式说明符的[格式化字符串](#格式化字符串)。
    
    -   **with字段**（可选）可以是![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)列表或![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)对象。它的值将用作![字符串](/images/Data_node_string.svg?42545)translate的匹配项列表，而不是直接显示。若未提供with字段，则![字符串](/images/Data_node_string.svg?42545)translate内的格式说明符会显示自身。
    -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)with：（可选）with字段的常见类型，只能使用纯文本。
        -   ![字符串](/images/Data_node_string.svg?42545)：一个文本匹配项，值是允许转义字符的字符串，可以包含纯文本的格式说明符（`%x %`）。载入![字符串](/images/Data_node_string.svg?42545)translate时会被转译成格式化字符串和真正的格式说明符`%%x %%`。
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)with：（可选）比![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)with支持更多的文本组件类型。值必然是对象![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)形式的文本组件，组件中![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)rawtext的值将用作匹配项列表。匹配项可以包含纯文本的格式说明符（`%x %`），载入![字符串](/images/Data_node_string.svg?42545)translate时会被转译成格式化字符串和真正的格式说明符`%%x %%`。

示例：

翻译文本组件中的本地化键名会被翻译。

{"rawtext": \[{"translate": "commands.op.success"}\]}

已将 %s 设为管理员

如果提供了with字段内容，会进一步对格式说明符进行匹配。

{"rawtext": \[{"translate": "commands.op.success", "with": \["朋友", "外人"\]}\]}

已将 朋友 设为管理员

翻译文本组件的`%%`最终会输出为纯文本`%`，而纯文本`%`载入格式化字符串时会转译成`%%`。

{"rawtext": \[{"translate": "原文 %% 纯文本 %%1 翻译文本 %%s", "with": {"rawtext": \[{"translate": "%%"}, {"text": "%"}\]}}\]}

原文 % 纯文本 % 翻译文本 %

#### 记分板分数（Score）

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=28&veaction=edit "编辑章节：记分板分数（Score）") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=28 "编辑章节的源代码： 记分板分数（Score）")\]

参见：[记分板](/w/%E8%AE%B0%E5%88%86%E6%9D%BF "记分板")和[命令/scoreboard](/w/%E5%91%BD%E4%BB%A4/scoreboard "命令/scoreboard")

显示[记分板](/w/%E8%AE%B0%E5%88%86%E6%9D%BF "记分板")中，目标分数持有者（Score Holder）在指定记分项（Objective）中的分数。

记分板分数内容类型的内容组件有以下字段：

-   -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)score：要显示的目标分数来源。分数不存在时不显示内容，多个分数以逗号加一个空格`,` 分隔。
        -   ![字符串](/images/Data_node_string.svg?42545)name：目标分数持有者。可以使用[目标选择器](/w/%E7%9B%AE%E6%A0%87%E9%80%89%E6%8B%A9%E5%99%A8 "目标选择器")、玩家名称、虚拟玩家名称（通过一个自定义键名直接在记分项中创建的分数持有者，`#`开头的虚拟玩家名称会在记分板中隐藏显示）、实体的记分板ID或分数持有者通配符`*`。若为目标选择器，允许匹配多个目标。若为`*`则表示为信息读者本身。
        -   ![字符串](/images/Data_node_string.svg?42545)objective：指定记分项的名称。必须与[创建记分项](/w/%E5%91%BD%E4%BB%A4/scoreboard#基岩版 "命令/scoreboard")时使用的名称一致。

示例：

当记分板分数组件同时显示多个分数时，情况如下。

{"rawtext": \[{"score": {"name": "@e", "objective": "obj"}}\]}

2, -3, 15

记分板分数组件会输出纯文本，并与`§`组合成格式化代码。

{"rawtext": \[{"text": "颜色：§"}, {"score": {"name": "\*", "objective": "obj"}}\]}

假如读者在obj中的分数是3456.

颜色：456

无论分数或记分项是否存在，当![字符串](/images/Data_node_string.svg?42545)name的值是`*`时，将必定显示一个内容。

{"rawtext": \[{"score": {"name": "\*", "objective": ""}}\]}

下方显示的是一个空内容。

#### 实体名称（Selector）

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=29&veaction=edit "编辑章节：实体名称（Selector）") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=29 "编辑章节的源代码： 实体名称（Selector）")\]

显示被选中的实体的名称，若该实体没有[自定义名称](/w/%E5%91%BD%E5%90%8D "命名")（Custom Name）则显示其[实体类型](/w/%E5%AE%9E%E4%BD%93%E7%B1%BB%E5%9E%8B "实体类型")的本地化名称（例如：狼）。

实体名称内容类型的内容组件有以下字段：

-   -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)selector：目标实体。可以使用[目标选择器](/w/%E7%9B%AE%E6%A0%87%E9%80%89%E6%8B%A9%E5%99%A8 "目标选择器")、玩家名称或实体通配符`*`，若为`*`则表示为信息读者本身。没有匹配的目标时不显示内容，多个匹配目标的名称以逗号加一个空格`,` 分隔。

示例：

实体名称组件会显示实体的名称或类型。

{"rawtext": \[{"selector": "@e\[type=wolf\]"}\]}

假如有一只被命名为小白的狼和一只没有被命名的狼。

小白, 狼

文本中的目标选择器必须经过转义。如特殊字符`"`必须写成`\"`。

{"rawtext": \[{"selector": "@a\[name=\\"Alex\\"\]"}\]}

假如玩家Alex在线。

Alex

当实体名称组件没有匹配到目标时，不会显示内容。

{"rawtext": \[{"selector": "@e\[c=0\]"}\]}

以上是一个无意义的文本组件。

### 组件解析

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=30&veaction=edit "编辑章节：组件解析") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=30 "编辑章节的源代码： 组件解析")\]

参见：[命令上下文 § 文本组件解析](/w/%E5%91%BD%E4%BB%A4%E4%B8%8A%E4%B8%8B%E6%96%87#文本组件解析 "命令上下文")

组件的解析遵循固定的规律和优先级。 一个文本组件中包含的多个内容组件，将依照排列顺序解析成文本后组合起来。当遇到嵌套的内容组件或文本组件时，子级组件会先被解析成文本后再回传给父级组件。

一个完整的解析过程：

步骤1：解析第一个内容组件的rawtext的内容。

{"rawtext": \[
  // rawtext的内容是子级组件
  {"text": "A", "rawtext": \[{"text": "B"}, {"text": "C"}\]},
  {"translate":"%%s", "with": {"rawtext": \[
    {"text":"D"},
    {          },
    {"text":"E"}
  \]}}
\]}

步骤2：合并第一个内容组件的rawtext。

{"rawtext": \[
  {"text": "A", "rawtext": \["B", "C"\]},  // 列表内容会合并
  {"translate":"%%s", "with": {"rawtext": \[
    {"text":"D"},
    {          },
    {"text":"E"}
  \]}}
\]}

步骤3：合并第一个内容组件。

{"rawtext": \[
  {"text": "A", "rawtext": "BC"},  //主内容和扩展内容会合并
  {"translate":"%%s", "with": {"rawtext": \[
    {"text":"D"},
    {          },
    {"text":"E"}
  \]}}
\]}

步骤4：解析第二个内容组件的with字段的rawtext。

{"rawtext": \[
  "ABC",
  {"translate":"%%s", "with": {"rawtext": \[  // rawtext是子级组件
    {"text":"D"},
    {          },  // 无内容组件会被忽略
    {"text":"E"}
  \]}}
\]}

步骤5：生成匹配项列表。

{"rawtext": \[
  "ABC",
  // rawtext会解析成匹配项列表
  {"translate":"%%s", "with": {"rawtext": \["D", "E"\]}}
\]}

步骤6：解析第二个内容组件。

{"rawtext": \[
  "ABC",
  // translate的格式化字符串开始解析
  {"translate":"%%s", "with": \["D", "E"\]}
\]}

步骤7：合并所有内容。

{"rawtext": \[  // 所有内容合并
  "ABC",
  "D"
\]}

最终的显示结果如下。

ABCD

最终一个文本组件会生成一个纯文本，只有当真正显示出来时[格式化代码](/w/%E6%A0%BC%E5%BC%8F%E5%8C%96%E4%BB%A3%E7%A0%81 "格式化代码")才会被解析。这意味着记分板分数组件的内容可以与`§`组合成格式化代码。

不停变换色彩的标题：

创建一个记分项`/[scoreboard](/w/%E5%91%BD%E4%BB%A4/scoreboard "命令/scoreboard") objectives add var dummy`，然后不断地依次执行下方命令：

scoreboard players add color var 1
execute if score color var matches 10.. run scoreboard players set color var 0
titleraw @a title {"rawtext":\[{"text":"§"},{"score":{"name":"color","objective":"var"}},{"text":"Rainbow"}\]}

当一个内容组件内包含多个必要字段时，内容组件的内容类型将由优先级最大的字段决定而无关顺序。必要字段优先级如下：![字符串](/images/Data_node_string.svg?42545)translate\>![字符串](/images/Data_node_string.svg?42545)text\>![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)score\>![字符串](/images/Data_node_string.svg?42545)selector。这意味着translate将覆盖text的内容。

示例：

{"rawtext": \[{"translate":"内容A", "translate": "内容B", "text": "内容C"}\]}

内容组件中优先级最大的字段是![字符串](/images/Data_node_string.svg?42545)translate，并且第二个translate字段将覆盖先前的translate字段。最终将显示内容B。

内容B

### 内容的定义

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=31&veaction=edit "编辑章节：内容的定义") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=31 "编辑章节的源代码： 内容的定义")\]

**内容**在文本组件上下文中具有特定的含义，“内容”只能由四个基本字段定义（![字符串](/images/Data_node_string.svg?42545)translate、![字符串](/images/Data_node_string.svg?42545)text、![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)score和![字符串](/images/Data_node_string.svg?42545)selector）。而在特定条件下![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)score和![字符串](/images/Data_node_string.svg?42545)selector可能无法显示内容。

值得注意的是，“不显示内容”与“显示空内容”是不同的概念，前者将导致难以预料或不符合预期的结果。

-   不显示内容（无内容）：可以被视为是一次失败的结果。当一个内容组件不显示内容时，它将直接被忽略；当文本组件不显示内容时，将“无事发生”。
-   显示空内容（空内容）：空内容可以通过`{"text": ""}`轻易获取。显示空内容的内容组件不会被忽略，因此可以通过组合可能不显示内容的内容组件和空内容以避免不可预期的结果。

### 字符串与格式化字符串

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=32&veaction=edit "编辑章节：字符串与格式化字符串") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=32 "编辑章节的源代码： 字符串与格式化字符串")\]

文本组件中的所有字符串都是纯文本![字符串](/images/Data_node_string.svg?42545)形式的文本组件，因此可以接受部分转义字符和[Unicode编码](https://zh.wikipedia.org/wiki/Unicode "wzh:Unicode")。所有的转义字符都会在进行复杂的组件解析前完成转义。

**特殊字符**：尽管`\0`并不能被成功转义，但与之对应的`\u0000`却能正常地将当前字符串强行截断。

允许的转义字符及对应结果

转义字符

结果

对应编码

`\n`

换行

`\u000A`

`\r`

`õ`

`\u00F5`

`\t`

`õ`

`\u00F5`

`\b`

`ô`

`\u00F4`

`\f`

`ã`

`\u00E3`

`\\`

`\`

`\u005C`

`\"`

`"`

`\u0022`

#### 格式化字符串

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=33&veaction=edit "编辑章节：格式化字符串") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=33 "编辑章节的源代码： 格式化字符串")\]

**格式化字符串（Format String）**被应用在了翻译文本组件中，是一种存在变量的文本，具有普通字符串的所有特性。实际上它是源于[Java](https://zh.wikipedia.org/wiki/Java "wzh:Java")编程语言的一种字符串规范[\[5\]](#cite_note-6)，而在基岩版中它是通过[C++](https://zh.wikipedia.org/wiki/C%2B%2B "wzh:C++")编程语言模仿实现的。\[[需要验证](/w/Special:TalkPage/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6 "Special:TalkPage/文本组件")\]

格式化字符串中的变量被称为**格式说明符(Format Specifier)**，又分为按序匹配和按位匹配。

-   按序匹配：`%%<数据类型>`会匹配到匹配项列表中的“下一项”内容，并将基准点后移一位。优先级高于按位匹配，且按序匹配间会从左到右逐一进行匹配。
-   按位匹配：`%%<索引位置>`会匹配到从基准点算起的第n项内容，`0`即为基准点的位置。所有按位匹配是同时完成的。
-   带类型的按位匹配：`%%<索引位置>$<数据类型>`一般会在语言文件中看到，但在翻译文本组件的**with**字段内不生效。

参数

-   `数据类型`：只能是`s`表示文本格式（String）或`d`表示数字格式（Digit）。在文本组件中两者没有区别。
-   `索引位置`：可以是`0`到`9`之间任意整数。

任何纯文本被载入格式化字符串时会被转译为格式化字符串，其中所有`%`将转译为`%%`。同理当格式化字符串输出内容成纯文本时，`%%`会转译为`%`。这也是区分纯文本和格式化字符串的主要方法。

一个完整的解析过程：

{"rawtext": \[{"translate":"translation\\u002Etest\\u002Ecomplex", "with": \["$s%d", "你", "我", "他", "你们", "好%0"\]}\]}

字符串解析阶段

`"translation**\u002E**test**\u002E**complex"`

翻译阶段

`"translation.test.complex"`

格式说明符整理阶段

`"前缀，%%s%%2**$s**，然后是 %%s 和 %%1**$s**，最后是 %%s，还有 %%1**$s**！"`

按序匹配阶段

`"前缀，**%%s**%%2，然后是 %%s 和 %%1，最后是 %%s，还有 %%1！"`

`"前缀，$s**%%d**%%2，然后是 %%s 和 %%1，最后是 %%s，还有 %%1！"`

`"前缀，$s你%%2，然后是 **%%s** 和 %%1，最后是 %%s，还有 %%1！"`

`"前缀，$s你%%2，然后是 我 和 %%1，最后是 **%%s**，还有 %%1！"`

按位匹配阶段

`"前缀，$s你**%%2**，然后是 我 和 **%%1**，最后是 他，还有 **%%1**！"`

输出阶段

`"前缀，$s你好**%%**0，然后是 我 和 你们，最后是 他，还有 你们！"`

`"前缀，$s你好%0，然后是 我 和 你们，最后是 他，还有 你们！"`

结果

前缀，$s你好%0，然后是 我 和 你们，最后是 他，还有 你们！

### 编写规范

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=34&veaction=edit "编辑章节：编写规范") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=34 "编辑章节的源代码： 编写规范")\]

由于文本组件属于较底层的元素，实际上许多特性并未经过官方确认，解析器随时可能被优化。为确保文本组件的长期有效性，你应当如此：

-   避免在文本组件中使用扩展内容组件列表![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)rawtext。你可以通过翻译文本组件替代它，或直接拆分成单一内容组件。
-   避免在一个内容组件中使用多个必要字段。
-   在[基岩版](/w/%E5%9F%BA%E5%B2%A9%E7%89%88 "基岩版")中，你应当避免在任何可能参与到格式化字符串解析的内容中使用`$s` `$d`。
-   避免使用可能具备特殊含义的转义字符（例如：`\t` `\r` `\u0000`）。
-   避免显示无内容的文本组件。

## 历史

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=35&veaction=edit "编辑章节：历史") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=35 "编辑章节的源代码： 历史")\]

[Java版](/w/Java%E7%89%88%E7%89%88%E6%9C%AC%E8%AE%B0%E5%BD%95 "Java版版本记录")

[1.7.2](/w/Java%E7%89%881.7.2 "Java版1.7.2")

[13w37a](/w/Java%E7%89%8813w37a "Java版13w37a")

加入了文本组件和`/[tellraw](/w/%E5%91%BD%E4%BB%A4/tellraw "命令/tellraw")`。

[1.8](/w/Java%E7%89%881.8 "Java版1.8")

[14w02a](/w/Java%E7%89%8814w02a "Java版14w02a")

加入了`insertion`。

[14w07a](/w/Java%E7%89%8814w07a "Java版14w07a")

加入了`score`。

[14w20a](/w/Java%E7%89%8814w20a "Java版14w20a")

加入了`/[title](/w/%E5%91%BD%E4%BB%A4/title "命令/title")`，其使用文本组件。

加入了`selector`。

[14w25a](/w/Java%E7%89%8814w25a "Java版14w25a")

现在支持在[告示牌](/w/%E5%91%8A%E7%A4%BA%E7%89%8C "告示牌")和[成书](/w/%E6%88%90%E4%B9%A6 "成书")内使用。

[1.12](/w/Java%E7%89%881.12 "Java版1.12")

[17w16a](/w/Java%E7%89%8817w16a "Java版17w16a")

加入了`keybind`。

[1.13](/w/Java%E7%89%881.13 "Java版1.13")

[18w01a](/w/Java%E7%89%8818w01a "Java版18w01a")

现在支持在自定义名称内使用。

[18w05a](/w/Java%E7%89%8818w05a "Java版18w05a")

加入了`/[bossbar](/w/%E5%91%BD%E4%BB%A4/bossbar "命令/bossbar")`，参数`<name>`使用文本组件。

[1.14](/w/Java%E7%89%881.14 "Java版1.14")

[18w43a](/w/Java%E7%89%8818w43a "Java版18w43a")

加入了`nbt`、`block`和`entity`。

现在支持在物品描述标签内使用。

[18w44a](/w/Java%E7%89%8818w44a "Java版18w44a")

加入了`interpret`。

[1.15](/w/Java%E7%89%881.15 "Java版1.15")

[19w39a](/w/Java%E7%89%8819w39a "Java版19w39a")

加入了`storage`。

[19w41a](/w/Java%E7%89%8819w41a "Java版19w41a")

为`clickEvent`加入了`copy_to_clipboard`。

[1.16](/w/Java%E7%89%881.16 "Java版1.16")

[20w17a](/w/Java%E7%89%8820w17a "Java版20w17a")

加入了`font`。

为`hoverEvent`加入了`contents`。`value`不再使用，但仍受支持。

`color`现在可以使用十六进制颜色码来自定义颜色。

`score`中的`value`被移除。

[1.17](/w/Java%E7%89%881.17 "Java版1.17")

[21w15a](/w/Java%E7%89%8821w15a "Java版21w15a")

加入了`separator`。

[1.19.1](/w/Java%E7%89%881.19.1 "Java版1.19.1")

[rc1](/w/Java%E7%89%881.19.1-rc1 "Java版1.19.1-rc1")

`clickEvent`的`run_command`事件现在不再支持直接发送聊天信息。这意味着现在所有的值都需要以`/`为前缀。

[pre6](/w/Java%E7%89%881.19.1-pre6 "Java版1.19.1-pre6")

`clickEvent`的`run_command`事件现在不再支持任何可发送聊天消息的命令。

[1.19.4](/w/Java%E7%89%881.19.4 "Java版1.19.4")

[23w03a](/w/Java%E7%89%8823w03a "Java版23w03a")

加入了`fallback`。

`translate`格式中的越界参数不再被静默忽略。

[1.20.3](/w/Java%E7%89%881.20.3 "Java版1.20.3")

[23w40a](/w/Java%E7%89%8823w40a "Java版23w40a")

加入了`type`，用于提升解析与错误检查的速度。

纯文本聊天组件（只有文本内容，无并列的组件，无格式）现在总会被序列化成字符串，而非`{"text": "*字符串*"}`。

聊天组件现在会在通过网络发送时序列化。

`show_entity`的`id`字段现在可接受4个整型值所构成的数组形式的UUID。

`translate`组件内的数值与布尔型参数不再被转换成字符串。

不再支持`null`和`[]`JSON文本表达式。

若`color`、`clickEvent`和`hoverEvent`类型字段中出现错误，现在将不再被静默忽略。

[1.21.4](/w/Java%E7%89%881.21.4 "Java版1.21.4")

[24w44a](/w/Java%E7%89%8824w44a "Java版24w44a")

加入了`shadow_color`用于指定文本阴影颜色。

[1.21.5](/w/Java%E7%89%881.21.5 "Java版1.21.5")

[25w02a](/w/Java%E7%89%8825w02a "Java版25w02a")

文本组件现在以NBT形式存储，不再存储为JSON字符串。

将点击事件的标签名由`clickEvent`重命名为`click_event`，并修改了除`copy_to_clipboard`外所有点击事件的格式：

-   将`change_page`的`value`重命名为`page`，现在需要数字而不是数字的字符串表示页码。
-   将`open_url`的`value`重命名为`url`，且不再静默忽略非`http`或`https`协议的URI。
-   将`open_file`的`value`重命名为`path`。
-   将`run_command`的`value`重命名为`command`，且不再必须带`/`前缀、包含非法字符时不再静默忽略。
-   将`suggest_command`的`value`重命名为`command`，且包含非法字符时不再静默忽略。

将悬停事件的标签名由`hoverEvent`重命名为`hover_event`，并修改了所有悬停事件的格式：

-   将`show_item`的`contents`移到根标签，并移除了早已弃用的`value`。
-   将`show_entity`的`contents`移到根标签，且重命名`id`为`uuid`、重命名`type`为`id`，并移除了早已弃用的`value`。
-   将`show_text`的`contents`重命名为`text`，并移除了早已弃用的`value`。

[25w03a](/w/Java%E7%89%8825w03a "Java版25w03a")

将悬停事件`show_text`的`text`字段重命名为`value`。

[25w05a](/w/Java%E7%89%8825w05a "Java版25w05a")

现在`/[bossbar](/w/%E5%91%BD%E4%BB%A4/bossbar "命令/bossbar")`、`/[scoreboard](/w/%E5%91%BD%E4%BB%A4/scoreboard "命令/scoreboard")`和`/[team](/w/%E5%91%BD%E4%BB%A4/team "命令/team")`中的文本组件以命令执行者`@s`解析。

[1.21.6](/w/Java%E7%89%881.21.6 "Java版1.21.6")

[25w20a](/w/Java%E7%89%8825w20a "Java版25w20a")

加入了`custom`和`show_dialog`点击事件。

现在由客户端处理的`run_command`点击事件指定的命令无法解析或所需的权限等级大于0时，玩家点击文本后会出现一个确认屏幕。

[pre1](/w/Java%E7%89%881.21.6-pre1 "Java版1.21.6-pre1")

现在`custom`点击事件发送的负载可以是任意NBT标签。

[1.21.9](/w/Java%E7%89%881.21.9 "Java版1.21.9")

[25w32a](/w/Java%E7%89%8825w32a "Java版25w32a")

加入了`object`文本组件类型。

[25w33a](/w/Java%E7%89%8825w33a "Java版25w33a")

现在`run_command`点击事件中的命令若需要聊天参数，则会进入一个确认屏幕而非静默忽略。

[25w35a](/w/Java%E7%89%8825w35a "Java版25w35a")

扩充了`object`文本组件类型的行为，加入了`player`类型。

[Java版（即将到来）](/w/%E8%AE%A1%E5%88%92%E7%89%88%E6%9C%AC#Java版 "计划版本")

[26.1](/w/Java%E7%89%8826.1 "Java版26.1")

[snapshot-5](/w/Java%E7%89%8826.1-snapshot-5 "Java版26.1-snapshot-5")

现在NBT组件若`interpret`设置为`false`，则各SNBT文本会经过语法高亮渲染，而非单一扁平的简单文本。

[snapshot-8](/w/Java%E7%89%8826.1-snapshot-8 "Java版26.1-snapshot-8")

在NBT组件的语法结构中加入了`plain`字段。

[pre-1](/w/Java%E7%89%8826.1-pre-1 "Java版26.1-pre-1")

向精灵图组件加入了`fallback`字段，以在无法渲染精灵图的场合显示出回落文本

现在若在MOTD中指定玩家精灵图，则游戏会强制渲染回落文本。

[pre-2](/w/Java%E7%89%8826.1-pre-2 "Java版26.1-pre-2")

MOTD中嵌套超过16层的文本组件现在将被丢弃，并用省略号替代。

现在MOTD的会将非玩家精灵图强制渲染为回落文本，而非将玩家精灵图渲染为回落文本。

[pre-3](/w/Java%E7%89%8826.1-pre-3 "Java版26.1-pre-3")

现在MOTD的会将玩家精灵图强制渲染为回落文本，而非将非玩家精灵图渲染为回落文本。

[基岩版](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E7%89%88%E6%9C%AC%E8%AE%B0%E5%BD%95 "基岩版版本记录")

[1.9.0](/w/%E5%9F%BA%E5%B2%A9%E7%89%881.9.0 "基岩版1.9.0")

[1.9.0.0](/w/%E5%9F%BA%E5%B2%A9%E7%89%881.9.0.0 "基岩版1.9.0.0")

加入了`/[tellraw](/w/%E5%91%BD%E4%BB%A4/tellraw "命令/tellraw")`，文本组件用于支持该命令。

[1.16.100](/w/%E5%9F%BA%E5%B2%A9%E7%89%881.16.100 "基岩版1.16.100")

[1.16.100.55](/w/%E5%9F%BA%E5%B2%A9%E7%89%881.16.100.55 "基岩版1.16.100.55")

加入了`score`和`selector`。

## 参见

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=36&veaction=edit "编辑章节：参见") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=36 "编辑章节的源代码： 参见")\]

-   [JSON](/w/JSON "JSON")
-   [命令](/w/%E5%91%BD%E4%BB%A4 "命令")
-   [格式化代码](/w/%E6%A0%BC%E5%BC%8F%E5%8C%96%E4%BB%A3%E7%A0%81 "格式化代码")
-   [Tutorial:文本组件](/w/Tutorial:%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6 "Tutorial:文本组件")

## 注释

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=37&veaction=edit "编辑章节：注释") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=37 "编辑章节的源代码： 注释")\]

1.  ↑ [1.0](#cite_ref-motd_3-0) [1.1](#cite_ref-motd_3-1) 客户端具有将MOTD按文本组件解析的能力，尽管原版服务端只会发送纯字符串。

## 参考

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=38&veaction=edit "编辑章节：参考") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=38 "编辑章节的源代码： 参考")\]

1.  [↑](#cite_ref-1) [MC-103171](https://bugs.mojang.com/browse/MC-103171 "mojira:MC-103171")
2.  [↑](#cite_ref-2) [MC-177273](https://bugs.mojang.com/browse/MC-177273 "mojira:MC-177273")
3.  [↑](#cite_ref-4) [MC-124024](https://bugs.mojang.com/browse/MC-124024 "mojira:MC-124024") — 漏洞状态为“已修复”。
4.  [↑](#cite_ref-5) [MC-56373](https://bugs.mojang.com/browse/MC-56373 "mojira:MC-56373") — 漏洞状态为“已修复”。
5.  [↑](#cite_ref-6) [https://www.geeksforgeeks.org/format-specifiers-in-java/](https://www.geeksforgeeks.org/format-specifiers-in-java/)

## 导航

\[[编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?section=39&veaction=edit "编辑章节：导航") | [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit&section=39 "编辑章节的源代码： 导航")\]

-   [查](/w/Template:Navbox_Java_Edition "Template:Navbox Java Edition")
-   [论](/w/Special:TalkPage/Template:Navbox_Java_Edition "Special:TalkPage/Template:Navbox Java Edition")
-   [编](/w/Special:EditPage/Template:Navbox_Java_Edition "Special:EditPage/Template:Navbox Java Edition")

 [![](/images/thumb/Java_Edition_icon_2.png/18px-Java_Edition_icon_2.png?84f96)](/w/Java%E7%89%88 "Java版") [Java版](/w/Java%E7%89%88 "Java版")

版本

-   [演示版](/w/%E6%BC%94%E7%A4%BA%E6%A8%A1%E5%BC%8F "演示模式")
    -   [地点](/w/%E6%BC%94%E7%A4%BA%E6%A8%A1%E5%BC%8F/%E5%9C%B0%E7%82%B9 "演示模式/地点")
-   [PC Gamer演示版](/w/PC_Gamer%E6%BC%94%E7%A4%BA%E7%89%88 "PC Gamer演示版")（[Beta 1.3](/w/Java%E7%89%88Beta_1.3 "Java版Beta 1.3")）

开发周期

[版本记录](/w/Java%E7%89%88%E7%89%88%E6%9C%AC%E8%AE%B0%E5%BD%95 "Java版版本记录")

-   [![](/images/BlockSprite_grass-block-top-revision-1.png?c2740)](/w/Java%E7%89%88pre-Classic "Java版pre-Classic")[pre-Classic](/w/Java%E7%89%88pre-Classic "Java版pre-Classic")
-   [![](/images/BlockSprite_bricks-revision-1.png?5c857)](/w/Java%E7%89%88Classic "Java版Classic")[Classic](/w/Java%E7%89%88Classic "Java版Classic")
    -   [![](/images/BlockSprite_grass-block-revision-1.png?dfa18)](/w/%E6%97%A9%E6%9C%9F%E5%88%9B%E9%80%A0 "早期创造")[早期创造](/w/%E6%97%A9%E6%9C%9F%E5%88%9B%E9%80%A0 "早期创造")
    -   [![](/images/EnvSprite_player.png?b2666)](/w/%E5%A4%9A%E4%BA%BA%E6%B5%8B%E8%AF%95 "多人测试")[多人测试](/w/%E5%A4%9A%E4%BA%BA%E6%B5%8B%E8%AF%95 "多人测试")
    -   [![](/images/Heart.svg?e7b69)](/w/%E7%94%9F%E5%AD%98%E6%B5%8B%E8%AF%95 "生存测试") [生存测试](/w/%E7%94%9F%E5%AD%98%E6%B5%8B%E8%AF%95 "生存测试")
    -   [![](/images/BlockSprite_gold-block-side.png?e4651)](/w/%E5%90%8E%E6%9C%9F%E5%88%9B%E9%80%A0 "后期创造")[后期创造](/w/%E5%90%8E%E6%9C%9F%E5%88%9B%E9%80%A0 "后期创造")
-   [![](/images/BlockSprite_cobblestone-revision-2.png?6c723)](/w/Java%E7%89%88Indev "Java版Indev")[Indev](/w/Java%E7%89%88Indev "Java版Indev")
-   [![](/images/thumb/Oak_Door_%28item%29_JE1.png/16px-Oak_Door_%28item%29_JE1.png?fe0b2)](/w/Java%E7%89%88Infdev "Java版Infdev") [Infdev](/w/Java%E7%89%88Infdev "Java版Infdev")
-   [![](/images/EnvSprite_nether-portal.png?77ce5)](/w/Java%E7%89%88Alpha "Java版Alpha")[Alpha](/w/Java%E7%89%88Alpha "Java版Alpha")
-   [![](/images/ItemSprite_cookie.png?75eb2)](/w/Java%E7%89%88Beta "Java版Beta")[Beta](/w/Java%E7%89%88Beta "Java版Beta")
-   [![](/images/BlockSprite_grass-block.png?9bd1e)](/w/%E6%AD%A3%E5%BC%8F%E7%89%88 "正式版")[正式版](/w/%E6%AD%A3%E5%BC%8F%E7%89%88 "正式版")
-   [![](/images/BlockSprite_crafting-table.png?7c62b)](/w/Java%E7%89%88%E7%89%88%E6%9C%AC%E8%AE%B0%E5%BD%95/%E5%BC%80%E5%8F%91%E7%89%88%E6%9C%AC "Java版版本记录/开发版本")[开发版本](/w/Java%E7%89%88%E7%89%88%E6%9C%AC%E8%AE%B0%E5%BD%95/%E5%BC%80%E5%8F%91%E7%89%88%E6%9C%AC "Java版版本记录/开发版本")

-   [实验性内容](/w/%E5%AE%9E%E9%AA%8C%E6%80%A7%E5%86%85%E5%AE%B9 "实验性内容")
-   [![](/images/BlockSprite_gear.png?2cf84)](/w/Java%E7%89%88%E5%B7%B2%E7%A7%BB%E9%99%A4%E7%89%B9%E6%80%A7 "Java版已移除特性")[已移除特性](/w/Java%E7%89%88%E5%B7%B2%E7%A7%BB%E9%99%A4%E7%89%B9%E6%80%A7 "Java版已移除特性")
    -   [方块](/w/Java%E7%89%88%E5%B7%B2%E7%A7%BB%E9%99%A4%E6%96%B9%E5%9D%97 "Java版已移除方块")
    -   [物品](/w/Java%E7%89%88%E5%B7%B2%E7%A7%BB%E9%99%A4%E7%89%A9%E5%93%81 "Java版已移除物品")
    -   [配方](/w/Java%E7%89%88%E5%B7%B2%E7%A7%BB%E9%99%A4%E9%85%8D%E6%96%B9 "Java版已移除配方")
-   [![](/images/ItemSprite_awkward-potion.png?398bd)](/w/Java%E7%89%88%E6%9C%AA%E4%BD%BF%E7%94%A8%E7%89%B9%E6%80%A7 "Java版未使用特性")[未使用特性](/w/Java%E7%89%88%E6%9C%AA%E4%BD%BF%E7%94%A8%E7%89%B9%E6%80%A7 "Java版未使用特性")
-   [![](/images/ItemSprite_knowledge-book.png?793c1)](/w/Java%E7%89%88%E7%8B%AC%E6%9C%89%E7%89%B9%E6%80%A7 "Java版独有特性")[独有特性](/w/Java%E7%89%88%E7%8B%AC%E6%9C%89%E7%89%B9%E6%80%A7 "Java版独有特性")
-   [![](/images/EnvSprite_sky-dimension.png?f3ab9)](/w/Java%E7%89%88%E6%8F%90%E5%8F%8A%E7%89%B9%E6%80%A7 "Java版提及特性")[提及特性](/w/Java%E7%89%88%E6%8F%90%E5%8F%8A%E7%89%B9%E6%80%A7 "Java版提及特性")
    -   [插件API](/w/Java%E7%89%88%E6%8F%90%E5%8F%8A%E7%89%B9%E6%80%A7/%E6%8F%92%E4%BB%B6API "Java版提及特性/插件API")
-   [计划版本](/w/%E8%AE%A1%E5%88%92%E7%89%88%E6%9C%AC#Java版 "计划版本")

技术

-   [已知漏洞](https://bugs.mojang.com/browse/MC "mojira:MC")
    -   [*启动器*](https://bugs.mojang.com/browse/MCL "mojira:MCL")
-   [硬件需求](/w/Java%E7%89%88%E7%A1%AC%E4%BB%B6%E9%9C%80%E6%B1%82 "Java版硬件需求")
-   [方块实体](/w/%E6%96%B9%E5%9D%97%E5%AE%9E%E4%BD%93 "方块实体")
-   [命令](/w/%E5%91%BD%E4%BB%A4 "命令")
    -   [函数](/w/Java%E7%89%88%E5%87%BD%E6%95%B0 "Java版函数")
-   [崩溃](/w/%E5%B4%A9%E6%BA%83 "崩溃")
-   [数据值](/w/Java%E7%89%88%E6%95%B0%E6%8D%AE%E5%80%BC "Java版数据值")
    -   [Classic](/w/Java%E7%89%88%E6%95%B0%E6%8D%AE%E5%80%BC/Classic "Java版数据值/Classic")
    -   [Indev](/w/Java%E7%89%88%E6%95%B0%E6%8D%AE%E5%80%BC/Indev "Java版数据值/Indev")
    -   [扁平化前](/w/Java%E7%89%88%E6%95%B0%E6%8D%AE%E5%80%BC/%E6%89%81%E5%B9%B3%E5%8C%96%E5%89%8D "Java版数据值/扁平化前")
-   [数据版本](/w/%E6%95%B0%E6%8D%AE%E7%89%88%E6%9C%AC "数据版本")
-   [调试屏幕](/w/%E8%B0%83%E8%AF%95%E5%B1%8F%E5%B9%95 "调试屏幕")
-   [格式化代码](/w/%E6%A0%BC%E5%BC%8F%E5%8C%96%E4%BB%A3%E7%A0%81 "格式化代码")
-   [高度图](/w/%E9%AB%98%E5%BA%A6%E5%9B%BE "高度图")
-   [键控代码](/w/%E9%94%AE%E6%8E%A7%E4%BB%A3%E7%A0%81 "键控代码")
-   [启动器](/w/Minecraft%E5%90%AF%E5%8A%A8%E5%99%A8 "Minecraft启动器")
    -   [快速进入游戏](/w/%E5%BF%AB%E9%80%9F%E8%BF%9B%E5%85%A5%E6%B8%B8%E6%88%8F "快速进入游戏")
-   [注册表](/w/%E6%B3%A8%E5%86%8C%E8%A1%A8 "注册表")
-   [命名空间ID](/w/%E5%91%BD%E5%90%8D%E7%A9%BA%E9%97%B4ID "命名空间ID")
-   [标签](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE "Java版标签")
-   [兴趣点](/w/%E5%85%B4%E8%B6%A3%E7%82%B9 "兴趣点")
-   [协议版本](/w/%E5%8D%8F%E8%AE%AE%E7%89%88%E6%9C%AC "协议版本")
-   [种子](/w/%E7%A7%8D%E5%AD%90%EF%BC%88%E4%B8%96%E7%95%8C%E7%94%9F%E6%88%90%EF%BC%89 "种子（世界生成）")
-   [粒子](/w/Java%E7%89%88%E7%B2%92%E5%AD%90 "Java版粒子")
-   [统计信息](/w/%E7%BB%9F%E8%AE%A1%E4%BF%A1%E6%81%AF "统计信息")
-   [遥测](/w/%E9%81%A5%E6%B5%8B "遥测")
-   [刻](/w/%E5%88%BB "刻")
-   [UUID](/w/%E9%80%9A%E7%94%A8%E5%94%AF%E4%B8%80%E8%AF%86%E5%88%AB%E7%A0%81 "通用唯一识别码")
-   [出生点保护](/w/%E5%87%BA%E7%94%9F%E7%82%B9%E4%BF%9D%E6%8A%A4 "出生点保护")
-   [坐标](/w/%E5%9D%90%E6%A0%87 "坐标")
-   [世界加载屏幕](/w/%E4%B8%96%E7%95%8C%E5%8A%A0%E8%BD%BD%E5%B1%8F%E5%B9%95 "世界加载屏幕")
-   [社交屏幕](/w/%E7%A4%BE%E4%BA%A4%E5%B1%8F%E5%B9%95 "社交屏幕")

[开发资源](/w/%E5%BC%80%E5%8F%91%E8%B5%84%E6%BA%90 "开发资源")

-   文本组件
-   [NBT格式](/w/NBT%E6%A0%BC%E5%BC%8F "NBT格式")
-   [战利品表](/w/%E6%88%98%E5%88%A9%E5%93%81%E8%A1%A8 "战利品表")
-   [Mojang API](/w/Mojang_API "Mojang API")
-   [网络协议](/w/Java%E7%89%88%E7%BD%91%E7%BB%9C%E5%8D%8F%E8%AE%AE "Java版网络协议")
-   [远程控制台协议](/w/%E8%BF%9C%E7%A8%8B%E6%8E%A7%E5%88%B6%E5%8F%B0%E5%8D%8F%E8%AE%AE "远程控制台协议")
-   [服务端管理协议](/w/%E6%9C%8D%E5%8A%A1%E7%AB%AF%E7%AE%A1%E7%90%86%E5%8D%8F%E8%AE%AE "服务端管理协议")
-   [游戏测试](/w/%E6%B8%B8%E6%88%8F%E6%B5%8B%E8%AF%95 "游戏测试")
-   [混淆映射表](/w/%E6%B7%B7%E6%B7%86%E6%98%A0%E5%B0%84%E8%A1%A8 "混淆映射表")
-   [调试工具](/w/%E8%B0%83%E8%AF%95%E5%B7%A5%E5%85%B7 "调试工具")
-   [Brigadier](/w/Brigadier "Brigadier")
-   `[.minecraft](/w/.minecraft ".minecraft")`
-   [存档格式](/w/Java%E7%89%88%E5%AD%98%E6%A1%A3%E6%A0%BC%E5%BC%8F "Java版存档格式")
-   [结构存储格式](/w/%E7%BB%93%E6%9E%84%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "结构存储格式")（[Schematic文件格式](/w/Schematic%E6%96%87%E4%BB%B6%E6%A0%BC%E5%BC%8F "Schematic文件格式")）
-   `[version_manifest.json](/w/Version_manifest.json "Version manifest.json")`

过时开发资源

-   [Classic服务器协议](/w/Classic%E6%9C%8D%E5%8A%A1%E5%99%A8%E5%8D%8F%E8%AE%AE "Classic服务器协议")
-   [al\_version](/w/Al_version "Al version")
-   [无限世界预览器](/w/%E6%97%A0%E9%99%90%E4%B8%96%E7%95%8C%E9%A2%84%E8%A7%88%E5%99%A8 "无限世界预览器")
-   [旧版Minecraft验证](/w/%E6%97%A7%E7%89%88Minecraft%E9%AA%8C%E8%AF%81 "旧版Minecraft验证")
-   [材料](/w/Java%E7%89%88%E6%9D%90%E6%96%99 "Java版材料")
-   [出生点区块](/w/%E5%87%BA%E7%94%9F%E7%82%B9%E5%8C%BA%E5%9D%97 "出生点区块")
-   [已配置的地表生成器](/w/%E5%B7%B2%E9%85%8D%E7%BD%AE%E7%9A%84%E5%9C%B0%E8%A1%A8%E7%94%9F%E6%88%90%E5%99%A8 "已配置的地表生成器")

多人游戏

-   [服务器](/w/%E6%9C%8D%E5%8A%A1%E5%99%A8 "服务器")
-   [Minecraft Realms](/w/Realms "Realms")
    -   [内容创作者计划](/w/Java%E7%89%88Realms%E5%86%85%E5%AE%B9%E5%88%9B%E4%BD%9C%E8%80%85%E8%AE%A1%E5%88%92 "Java版Realms内容创作者计划")
-   [服务器列表](/w/%E6%9C%8D%E5%8A%A1%E5%99%A8%E5%88%97%E8%A1%A8 "服务器列表")
-   [服务端配置文件格式](/w/%E6%9C%8D%E5%8A%A1%E7%AB%AF%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6%E6%A0%BC%E5%BC%8F "服务端配置文件格式")
-   [服务器需求](/w/%E6%9C%8D%E5%8A%A1%E5%99%A8/%E9%9C%80%E6%B1%82 "服务器/需求")
-   [定制服务器](/w/%E5%AE%9A%E5%88%B6%E6%9C%8D%E5%8A%A1%E5%99%A8 "定制服务器")
-   [在线验证](/w/%E5%9C%A8%E7%BA%BF%E9%AA%8C%E8%AF%81 "在线验证")

游戏订制

-   [皮肤](/w/%E7%9A%AE%E8%82%A4 "皮肤")
-   [披风](/w/%E6%8A%AB%E9%A3%8E "披风")
-   [资源包](/w/%E8%B5%84%E6%BA%90%E5%8C%85 "资源包")
-   [数据包](/w/%E6%95%B0%E6%8D%AE%E5%8C%85 "数据包")
    -   [洞穴与山崖预览数据包](/w/%E6%B4%9E%E7%A9%B4%E4%B8%8E%E5%B1%B1%E5%B4%96%E9%A2%84%E8%A7%88%E6%95%B0%E6%8D%AE%E5%8C%85 "洞穴与山崖预览数据包")
    -   [实验性内容](/w/%E5%AE%9E%E9%AA%8C%E6%80%A7%E5%86%85%E5%AE%B9 "实验性内容")

-   [查](/w/Template:Navbox_Bedrock_Edition "Template:Navbox Bedrock Edition")
-   [论](/w/Special:TalkPage/Template:Navbox_Bedrock_Edition "Special:TalkPage/Template:Navbox Bedrock Edition")
-   [编](/w/Special:EditPage/Template:Navbox_Bedrock_Edition "Special:EditPage/Template:Navbox Bedrock Edition")

 [![](/images/thumb/Bedrock_Edition_icon_2.png/16px-Bedrock_Edition_icon_2.png?80b87)](/w/%E5%9F%BA%E5%B2%A9%E7%89%88 "基岩版") [基岩版](/w/%E5%9F%BA%E5%B2%A9%E7%89%88 "基岩版")

版本

已合并

-   [![](/images/Smartphone.svg?51df4)](/w/%E6%90%BA%E5%B8%A6%E7%89%88 "携带版") [携带版](/w/%E6%90%BA%E5%B8%A6%E7%89%88 "携带版")
-   [![Windows 10](/images/Windows.svg?7d510)](https://zh.wikipedia.org/wiki/Windows_10 "Windows 10") [Windows 10版](/w/Windows_10%E7%89%88 "Windows 10版")

移植到主机

-   [![Xbox One](/images/Xbox_One.svg?4165c)](https://zh.wikipedia.org/wiki/Xbox_One "Xbox One") [Xbox One版](/w/Xbox_One%E7%89%88 "Xbox One版")
-   [![Nintendo Switch](/images/Nintendo_Switch.svg?904e3)](https://zh.wikipedia.org/wiki/%E4%BB%BB%E5%A4%A9%E5%A0%82Switch "Nintendo Switch") [Nintendo Switch版](/w/Nintendo_Switch%E7%89%88 "Nintendo Switch版")
-   [![PlayStation 4](/images/PS4.svg?5e20e)](https://zh.wikipedia.org/wiki/PlayStation_4 "PlayStation 4") [PlayStation 4版](/w/PlayStation_4%E7%89%88 "PlayStation 4版")

![BlockSprite barrier.png：Minecraft中barrier的精灵图](/images/BlockSprite_barrier.png?7d049) 已终止

-   [![Apple TV](/images/AppleTVLogo.svg?ff27e)](https://zh.wikipedia.org/wiki/Apple_TV "Apple TV") [Apple TV版](/w/Apple_TV%E7%89%88 "Apple TV版")
-   [![Gear VR](/images/GearVR.svg?6ffef)](https://zh.wikipedia.org/wiki/%E4%B8%89%E6%98%9FGear_VR "Gear VR") [Gear VR版](/w/Gear_VR%E7%89%88 "Gear VR版")
-   [![Fire TV](/images/FireTV.svg?a900e)](https://zh.wikipedia.org/wiki/Amazon_Fire_TV "Fire TV") [Fire TV版](/w/Fire_TV%E7%89%88 "Fire TV版")

开发

[版本记录](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E7%89%88%E6%9C%AC%E8%AE%B0%E5%BD%95 "基岩版版本记录")

-   [![](/images/BlockSprite_bricks-revision-1.png?5c857)](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E7%89%88%E6%9C%AC%E8%AE%B0%E5%BD%95#Alpha "基岩版版本记录")[Alpha](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E7%89%88%E6%9C%AC%E8%AE%B0%E5%BD%95#Alpha "基岩版版本记录")
-   [![](/images/EntitySprite_ender-dragon.png?89e49)](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E7%89%88%E6%9C%AC%E8%AE%B0%E5%BD%95#1.0 "基岩版版本记录")[正式版](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E7%89%88%E6%9C%AC%E8%AE%B0%E5%BD%95#1.0 "基岩版版本记录")
-   [![](/images/BlockSprite_crafting-table.png?7c62b)](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E7%89%88%E6%9C%AC%E8%AE%B0%E5%BD%95/%E5%BC%80%E5%8F%91%E7%89%88%E6%9C%AC "基岩版版本记录/开发版本")[开发版本](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E7%89%88%E6%9C%AC%E8%AE%B0%E5%BD%95/%E5%BC%80%E5%8F%91%E7%89%88%E6%9C%AC "基岩版版本记录/开发版本")

-   [![](/images/EnvSprite_nether-reactor.png?03c80)](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E5%B7%B2%E7%A7%BB%E9%99%A4%E7%89%B9%E6%80%A7 "基岩版已移除特性")[已移除特性](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E5%B7%B2%E7%A7%BB%E9%99%A4%E7%89%B9%E6%80%A7 "基岩版已移除特性")
    -   [方块](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E5%B7%B2%E7%A7%BB%E9%99%A4%E6%96%B9%E5%9D%97 "基岩版已移除方块")
    -   [配方](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E5%B7%B2%E7%A7%BB%E9%99%A4%E9%85%8D%E6%96%B9 "基岩版已移除配方")
    -   [VR](/w/%E8%99%9A%E6%8B%9F%E7%8E%B0%E5%AE%9E "虚拟现实")
-   [![](/images/EntitySprite_camera.png?d4e4b)](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E6%9C%AA%E4%BD%BF%E7%94%A8%E7%89%B9%E6%80%A7 "基岩版未使用特性")[未使用特性](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E6%9C%AA%E4%BD%BF%E7%94%A8%E7%89%B9%E6%80%A7 "基岩版未使用特性")
-   [![](/images/ItemSprite_potion-of-decay.png?79d0e)](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E7%8B%AC%E6%9C%89%E7%89%B9%E6%80%A7 "基岩版独有特性")[独有特性](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E7%8B%AC%E6%9C%89%E7%89%B9%E6%80%A7 "基岩版独有特性")
-   [提及特性](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E6%8F%90%E5%8F%8A%E7%89%B9%E6%80%A7 "基岩版提及特性")
    -   [Super Duper图形包](/w/Super_Duper%E5%9B%BE%E5%BD%A2%E5%8C%85 "Super Duper图形包")
-   [计划版本](/w/%E8%AE%A1%E5%88%92%E7%89%88%E6%9C%AC#基岩版 "计划版本")
-   [![](/images/thumb/Minecraft_Preview_icon_2.png/16px-Minecraft_Preview_icon_2.png?622c0)](/w/Minecraft_Preview "Minecraft Preview") [Minecraft Preview](/w/Minecraft_Preview "Minecraft Preview")

技术性

-   [已知漏洞](https://bugs.mojang.com/browse/MCPE "mojira:MCPE")
    -   [*启动器*](https://bugs.mojang.com/browse/MCL "mojira:MCL")
-   [快速进入游戏](/w/%E5%BF%AB%E9%80%9F%E8%BF%9B%E5%85%A5%E6%B8%B8%E6%88%8F "快速进入游戏")
-   [RenderDragon](/w/RenderDragon "RenderDragon")
-   [数据值](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E6%95%B0%E6%8D%AE%E5%80%BC "基岩版数据值")
    -   [Alpha 0.2.0前](/w/%E6%90%BA%E5%B8%A6%E7%89%88%E6%95%B0%E6%8D%AE%E5%80%BC/Alpha_0.2.0%E5%89%8D "携带版数据值/Alpha 0.2.0前")
-   [实体组件](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E5%AE%9E%E4%BD%93%E7%BB%84%E4%BB%B6 "基岩版实体组件")
-   [配置要求](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E7%A1%AC%E4%BB%B6%E6%80%A7%E8%83%BD "基岩版硬件性能")
    -   [灵动视效](/w/Vibrant_Visuals#设备支持 "Vibrant Visuals")
    -   [光线追踪](/w/RenderDragon#光线追踪 "RenderDragon")
-   [构建信息](/w/%E6%9E%84%E5%BB%BA%E4%BF%A1%E6%81%AF "构建信息")
-   [存档格式](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E5%AD%98%E6%A1%A3%E6%A0%BC%E5%BC%8F "基岩版存档格式")
-   [NBT格式](/w/NBT%E6%A0%BC%E5%BC%8F "NBT格式")
-   [动画](/w/%E5%8A%A8%E7%94%BB "动画")
-   [附加包](/w/%E9%99%84%E5%8A%A0%E5%8C%85 "附加包")
    -   [Molang](/w/Molang "Molang")
-   [游戏测试](/w/%E6%B8%B8%E6%88%8F%E6%B5%8B%E8%AF%95 "游戏测试")
-   [配方](/w/%E9%85%8D%E6%96%B9 "配方")
-   [方块实体](/w/%E6%96%B9%E5%9D%97%E5%AE%9E%E4%BD%93 "方块实体")
-   [命名空间ID](/w/%E5%91%BD%E5%90%8D%E7%A9%BA%E9%97%B4ID "命名空间ID")
-   文本组件
-   `[com.mojang](/w/Com.mojang "Com.mojang")`
-   [命令](/w/%E5%91%BD%E4%BB%A4 "命令")
    -   [函数](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E5%87%BD%E6%95%B0 "基岩版函数")
    -   [开发者命令](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E5%BC%80%E5%8F%91%E8%80%85%E5%91%BD%E4%BB%A4 "基岩版开发者命令")
-   [生成事件](/w/%E7%94%9F%E6%88%90%E4%BA%8B%E4%BB%B6 "生成事件")
-   [坐标](/w/%E5%9D%90%E6%A0%87 "坐标")
-   [材料](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E6%9D%90%E6%96%99 "基岩版材料")
-   [种子](/w/%E7%A7%8D%E5%AD%90%EF%BC%88%E4%B8%96%E7%95%8C%E7%94%9F%E6%88%90%EF%BC%89 "种子（世界生成）")
-   [粒子](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E7%B2%92%E5%AD%90 "基岩版粒子")
-   [专用服务器](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E4%B8%93%E7%94%A8%E6%9C%8D%E5%8A%A1%E5%99%A8 "基岩版专用服务器")
-   `[manifest.json](/w/Manifest.json "Manifest.json")`
-   `[sound_definitions.json](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E5%A3%B0%E9%9F%B3%E4%BA%8B%E4%BB%B6 "基岩版声音事件")`
-   `[options.txt](/w/%E5%AE%A2%E6%88%B7%E7%AB%AF%E9%80%89%E9%A1%B9%E6%96%87%E4%BB%B6%E6%A0%BC%E5%BC%8F "客户端选项文件格式")`
-   [标签](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E6%A0%87%E7%AD%BE "基岩版标签")
-   [刻](/w/%E5%88%BB "刻")
-   [常加载区域](/w/%E5%B8%B8%E5%8A%A0%E8%BD%BD%E5%8C%BA%E5%9F%9F "常加载区域")
-   [世界加载屏幕](/w/%E4%B8%96%E7%95%8C%E5%8A%A0%E8%BD%BD%E5%B1%8F%E5%B9%95 "世界加载屏幕")
-   [协议版本](/w/%E5%8D%8F%E8%AE%AE%E7%89%88%E6%9C%AC "协议版本")
-   [族](/w/%E6%97%8F "族")
-   [定义](/w/%E5%AE%9A%E4%B9%89 "定义")
-   [基岩版编辑器](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E7%BC%96%E8%BE%91%E5%99%A8 "基岩版编辑器")

多人游戏

-   [服务器](/w/%E6%9C%8D%E5%8A%A1%E5%99%A8 "服务器")
-   [Realms Plus](/w/Realms_Plus "Realms Plus")
-   [服务器列表](/w/%E6%9C%8D%E5%8A%A1%E5%99%A8%E5%88%97%E8%A1%A8 "服务器列表")
-   `[server.properties](/w/%E6%9C%8D%E5%8A%A1%E7%AB%AF%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6%E6%A0%BC%E5%BC%8F "服务端配置文件格式")`
-   [服务器软件](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E6%9C%8D%E5%8A%A1%E5%99%A8%E8%BD%AF%E4%BB%B6 "基岩版服务器软件")
-   [在线验证](/w/%E5%9C%A8%E7%BA%BF%E9%AA%8C%E8%AF%81 "在线验证")

特色功能

-   [实验性玩法](/w/%E5%AE%9E%E9%AA%8C%E6%80%A7%E7%8E%A9%E6%B3%95 "实验性玩法")
-   [加载提示](/w/%E5%8A%A0%E8%BD%BD%E6%8F%90%E7%A4%BA "加载提示")
-   [种子模板](/w/%E7%A7%8D%E5%AD%90%E6%A8%A1%E6%9D%BF "种子模板")
-   [角色创建器](/w/%E8%A7%92%E8%89%B2%E5%88%9B%E5%BB%BA%E5%99%A8 "角色创建器")
    -   [表情](/w/%E8%A1%A8%E6%83%85 "表情")
-   [市场](/w/%E5%B8%82%E5%9C%BA "市场")
-   [精选服务器](/w/%E7%B2%BE%E9%80%89%E6%9C%8D%E5%8A%A1%E5%99%A8 "精选服务器")
-   [活动服务器](/w/%E6%B4%BB%E5%8A%A8%E6%9C%8D%E5%8A%A1%E5%99%A8 "活动服务器")
-   [分屏](/w/%E5%88%86%E5%B1%8F "分屏")
-   [Ore UI](/w/Ore_UI "Ore UI")

-   [查](/w/Template:Navbox_Java_customizable "Template:Navbox Java customizable")
-   [论](/w/Special:TalkPage/Template:Navbox_Java_customizable "Special:TalkPage/Template:Navbox Java customizable")
-   [编](/w/Special:EditPage/Template:Navbox_Java_customizable "Special:EditPage/Template:Navbox Java customizable")

[Java版](/w/Java%E7%89%88 "Java版")可自定义内容

基本概念

-   [注册表](/w/%E6%B3%A8%E5%86%8C%E8%A1%A8 "注册表")
    -   [命名空间ID](/w/%E5%91%BD%E5%90%8D%E7%A9%BA%E9%97%B4ID "命名空间ID")
    -   [标签](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE "Java版标签")
-   [命令](/w/%E5%91%BD%E4%BB%A4 "命令")
    -   [命令存储](/w/%E5%91%BD%E4%BB%A4%E5%AD%98%E5%82%A8%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "命令存储存储格式")
    -   [命令上下文](/w/%E5%91%BD%E4%BB%A4%E4%B8%8A%E4%B8%8B%E6%96%87 "命令上下文")
-   [NBT格式](/w/NBT%E6%A0%BC%E5%BC%8F "NBT格式")
    -   [NBT路径](/w/NBT%E8%B7%AF%E5%BE%84 "NBT路径")
-   [SNBT格式](/w/SNBT%E6%A0%BC%E5%BC%8F "SNBT格式")
-   [JSON](/w/JSON "JSON")
-   文本组件
-   [格式化代码](/w/%E6%A0%BC%E5%BC%8F%E5%8C%96%E4%BB%A3%E7%A0%81 "格式化代码")
-   [UUID](/w/%E9%80%9A%E7%94%A8%E5%94%AF%E4%B8%80%E8%AF%86%E5%88%AB%E7%A0%81 "通用唯一识别码")

[数据包](/w/%E6%95%B0%E6%8D%AE%E5%8C%85 "数据包")

-   `[pack.mcmeta](/w/Pack.mcmeta "Pack.mcmeta")`
-   [函数](/w/Java%E7%89%88%E5%87%BD%E6%95%B0 "Java版函数")
-   [结构模板](/w/%E7%BB%93%E6%9E%84%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "结构存储格式")
-   [声音事件](/w/Java%E7%89%88%E5%A3%B0%E9%9F%B3%E4%BA%8B%E4%BB%B6 "Java版声音事件")

注册

游戏行为

-   [战利品表](/w/%E6%88%98%E5%88%A9%E5%93%81%E8%A1%A8 "战利品表")
    -   [战利品上下文](/w/%E6%88%98%E5%88%A9%E5%93%81%E4%B8%8A%E4%B8%8B%E6%96%87 "战利品上下文")
    -   [随机序列](/w/%E9%9A%8F%E6%9C%BA%E5%BA%8F%E5%88%97%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "随机序列存储格式")
    -   [物品修饰器](/w/%E7%89%A9%E5%93%81%E4%BF%AE%E9%A5%B0%E5%99%A8 "物品修饰器")
    -   [谓词](/w/%E8%B0%93%E8%AF%8D "谓词")
    -   [槽位源](/w/%E6%A7%BD%E4%BD%8D%E6%BA%90 "槽位源")
-   [配方](/w/%E9%85%8D%E6%96%B9 "配方")
-   [进度](/w/%E8%BF%9B%E5%BA%A6%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "进度定义格式")
    -   [实体谓词](/w/%E5%AE%9E%E4%BD%93%E8%B0%93%E8%AF%8D "实体谓词")

定义格式

-   [旗帜图案](/w/%E6%97%97%E5%B8%9C%E5%9B%BE%E6%A1%88%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "旗帜图案定义格式")
-   [聊天类型](/w/%E8%81%8A%E5%A4%A9%E7%B1%BB%E5%9E%8B%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "聊天类型定义格式")
-   [伤害类型](/w/%E4%BC%A4%E5%AE%B3%E7%B1%BB%E5%9E%8B%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "伤害类型定义格式")
-   [对话框](/w/%E5%AF%B9%E8%AF%9D%E6%A1%86%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "对话框定义格式")
-   [魔咒](/w/%E9%AD%94%E5%92%92%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "魔咒定义格式")
-   [魔咒提供器](/w/%E9%AD%94%E5%92%92%E6%8F%90%E4%BE%9B%E5%99%A8%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "魔咒提供器定义格式")
-   [山羊角乐器](/w/%E5%B1%B1%E7%BE%8A%E8%A7%92%E4%B9%90%E5%99%A8%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "山羊角乐器定义格式")
-   [唱片机曲目](/w/%E5%94%B1%E7%89%87%E6%9C%BA%E6%9B%B2%E7%9B%AE%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "唱片机曲目定义格式")
-   [画变种](/w/%E7%94%BB%E5%8F%98%E7%A7%8D%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "画变种定义格式")
-   [测试环境](/w/%E6%B5%8B%E8%AF%95%E7%8E%AF%E5%A2%83%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "测试环境定义格式")
-   [测试实例](/w/%E6%B5%8B%E8%AF%95%E5%AE%9E%E4%BE%8B%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "测试实例定义格式")
-   [时间线](/w/%E6%97%B6%E9%97%B4%E7%BA%BF%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "时间线定义格式")
-   [交易集](/w/%E4%BA%A4%E6%98%93%E9%9B%86%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "交易集定义格式")\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]
-   [试炼刷怪笼配置](/w/%E8%AF%95%E7%82%BC%E5%88%B7%E6%80%AA%E7%AC%BC%E9%85%8D%E7%BD%AE%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "试炼刷怪笼配置定义格式")
-   [盔甲纹饰](/w/%E7%9B%94%E7%94%B2%E7%BA%B9%E9%A5%B0%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "盔甲纹饰定义格式")
-   [村民交易](/w/%E6%9D%91%E6%B0%91%E4%BA%A4%E6%98%93%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "村民交易定义格式")\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]
-   [世界时钟](/w/%E4%B8%96%E7%95%8C%E6%97%B6%E9%92%9F%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "世界时钟定义格式")\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]

生物变种

-   [猫变种](/w/%E7%8C%AB%E5%8F%98%E7%A7%8D%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "猫变种定义格式")
    -   [音效](/w/%E7%8C%AB%E9%9F%B3%E6%95%88%E5%8F%98%E7%A7%8D%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "猫音效变种定义格式")\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]
-   [牛变种](/w/%E7%89%9B%E5%8F%98%E7%A7%8D%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "牛变种定义格式")
    -   [音效](/w/%E7%89%9B%E9%9F%B3%E6%95%88%E5%8F%98%E7%A7%8D%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "牛音效变种定义格式")\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]
-   [鸡变种](/w/%E9%B8%A1%E5%8F%98%E7%A7%8D%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "鸡变种定义格式")
    -   [音效](/w/%E9%B8%A1%E9%9F%B3%E6%95%88%E5%8F%98%E7%A7%8D%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "鸡音效变种定义格式")\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]
-   [青蛙变种](/w/%E9%9D%92%E8%9B%99%E5%8F%98%E7%A7%8D%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "青蛙变种定义格式")
-   [猪变种](/w/%E7%8C%AA%E5%8F%98%E7%A7%8D%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "猪变种定义格式")
    -   [音效](/w/%E7%8C%AA%E9%9F%B3%E6%95%88%E5%8F%98%E7%A7%8D%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "猪音效变种定义格式")\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]
-   [狼变种](/w/%E7%8B%BC%E5%8F%98%E7%A7%8D%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "狼变种定义格式")
    -   [音效](/w/%E7%8B%BC%E9%9F%B3%E6%95%88%E5%8F%98%E7%A7%8D%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "狼音效变种定义格式")
-   [僵尸鹦鹉螺变种](/w/%E5%83%B5%E5%B0%B8%E9%B9%A6%E9%B9%89%E8%9E%BA%E5%8F%98%E7%A7%8D%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "僵尸鹦鹉螺变种定义格式")

[世界生成](/w/%E8%87%AA%E5%AE%9A%E4%B9%89%E4%B8%96%E7%95%8C%E7%94%9F%E6%88%90 "自定义世界生成")

-   [维度](/w/%E7%BB%B4%E5%BA%A6%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "维度定义格式")
-   [维度类型](/w/%E7%BB%B4%E5%BA%A6%E7%B1%BB%E5%9E%8B "维度类型")
-   [世界预设](/w/%E4%B8%96%E7%95%8C%E9%A2%84%E8%AE%BE%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "世界预设定义格式")
-   [超平坦预设](/w/%E8%B6%85%E5%B9%B3%E5%9D%A6%E4%B8%96%E7%95%8C%E7%94%9F%E6%88%90%E9%A2%84%E8%AE%BE "超平坦世界生成预设")
-   [噪声](/w/%E5%99%AA%E5%A3%B0 "噪声")
-   [噪声设置](/w/%E5%99%AA%E5%A3%B0%E8%AE%BE%E7%BD%AE "噪声设置")
-   [密度函数](/w/%E5%AF%86%E5%BA%A6%E5%87%BD%E6%95%B0 "密度函数")
-   [生物群系](/w/%E7%94%9F%E7%89%A9%E7%BE%A4%E7%B3%BB%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "生物群系定义格式")
-   [雕刻器](/w/%E9%9B%95%E5%88%BB%E5%99%A8%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "雕刻器定义格式")
-   [已配置的地物](/w/%E5%B7%B2%E9%85%8D%E7%BD%AE%E7%9A%84%E5%9C%B0%E7%89%A9 "已配置的地物")
-   [已放置的地物](/w/%E5%B7%B2%E6%94%BE%E7%BD%AE%E7%9A%84%E5%9C%B0%E7%89%A9 "已放置的地物")
-   [结构](/w/%E7%BB%93%E6%9E%84%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "结构定义格式")
-   [结构集](/w/%E7%BB%93%E6%9E%84%E9%9B%86 "结构集")
-   [模板池](/w/%E6%A8%A1%E6%9D%BF%E6%B1%A0 "模板池")
-   [处理器列表](/w/%E5%A4%84%E7%90%86%E5%99%A8%E5%88%97%E8%A1%A8 "处理器列表")
-   [环境属性](/w/%E7%8E%AF%E5%A2%83%E5%B1%9E%E6%80%A7 "环境属性")

[标签](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE "Java版标签")

-   [旗帜图案](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E6%97%97%E5%B8%9C%E5%9B%BE%E6%A1%88 "Java版标签/旗帜图案")
-   [方块](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E6%96%B9%E5%9D%97 "Java版标签/方块")
-   [伤害类型](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E4%BC%A4%E5%AE%B3%E7%B1%BB%E5%9E%8B "Java版标签/伤害类型")
-   [对话框](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E5%AF%B9%E8%AF%9D%E6%A1%86 "Java版标签/对话框")
-   [魔咒](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E9%AD%94%E5%92%92 "Java版标签/魔咒")
-   [实体类型](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E5%AE%9E%E4%BD%93%E7%B1%BB%E5%9E%8B "Java版标签/实体类型")
-   [流体](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E6%B5%81%E4%BD%93 "Java版标签/流体")
-   [函数](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E5%87%BD%E6%95%B0 "Java版标签/函数")
-   [游戏事件](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E6%B8%B8%E6%88%8F%E4%BA%8B%E4%BB%B6 "Java版标签/游戏事件")
-   [山羊角乐器](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E5%B1%B1%E7%BE%8A%E8%A7%92%E4%B9%90%E5%99%A8 "Java版标签/山羊角乐器")
-   [物品](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E7%89%A9%E5%93%81 "Java版标签/物品")
-   [画变种](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E7%94%BB%E5%8F%98%E7%A7%8D "Java版标签/画变种")
-   [兴趣点类型](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E5%85%B4%E8%B6%A3%E7%82%B9%E7%B1%BB%E5%9E%8B "Java版标签/兴趣点类型")
-   [药水效果](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E8%8D%AF%E6%B0%B4%E6%95%88%E6%9E%9C "Java版标签/药水效果")\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]
-   [时间线](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E6%97%B6%E9%97%B4%E7%BA%BF "Java版标签/时间线")
-   [村民交易](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E6%9D%91%E6%B0%91%E4%BA%A4%E6%98%93 "Java版标签/村民交易")\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]
-   世界生成
    -   [生物群系](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E7%94%9F%E7%89%A9%E7%BE%A4%E7%B3%BB "Java版标签/生物群系")
    -   [已配置的地物](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E5%B7%B2%E9%85%8D%E7%BD%AE%E7%9A%84%E5%9C%B0%E7%89%A9 "Java版标签/已配置的地物")\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]
    -   [超平坦预设](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E8%B6%85%E5%B9%B3%E5%9D%A6%E4%B8%96%E7%95%8C%E7%94%9F%E6%88%90%E9%A2%84%E8%AE%BE "Java版标签/超平坦世界生成预设")
    -   [结构](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E7%BB%93%E6%9E%84 "Java版标签/结构")
    -   [世界预设](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E4%B8%96%E7%95%8C%E9%A2%84%E8%AE%BE "Java版标签/世界预设")

[资源包](/w/%E8%B5%84%E6%BA%90%E5%8C%85 "资源包")

-   `[pack.mcmeta](/w/Pack.mcmeta "Pack.mcmeta")`
-   [纹理图集](/w/%E7%BA%B9%E7%90%86%E5%9B%BE%E9%9B%86 "纹理图集")
-   [纹理](/w/%E7%BA%B9%E7%90%86 "纹理")
-   [模型](/w/%E6%A8%A1%E5%9E%8B "模型")
-   [物品模型映射](/w/%E7%89%A9%E5%93%81%E6%A8%A1%E5%9E%8B%E6%98%A0%E5%B0%84 "物品模型映射")
-   [字体](/w/%E8%87%AA%E5%AE%9A%E4%B9%89%E5%AD%97%E4%BD%93 "自定义字体")
-   [着色器](/w/%E7%9D%80%E8%89%B2%E5%99%A8 "着色器")
-   [声音事件](/w/Java%E7%89%88%E5%A3%B0%E9%9F%B3%E4%BA%8B%E4%BB%B6 "Java版声音事件")
-   [装备模型](/w/%E8%A3%85%E5%A4%87%E6%A8%A1%E5%9E%8B "装备模型")
-   [路径点样式](/w/%E8%B7%AF%E5%BE%84%E7%82%B9%E6%A0%B7%E5%BC%8F "路径点样式")

相关条目

-   [属性](/w/%E5%B1%9E%E6%80%A7 "属性")
-   [数据组件](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6 "数据组件")
-   [数据组件谓词](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6%E8%B0%93%E8%AF%8D "数据组件谓词")
-   [粒子数据格式](/w/%E7%B2%92%E5%AD%90%E6%95%B0%E6%8D%AE%E6%A0%BC%E5%BC%8F "粒子数据格式")
-   [实体数据格式](/w/%E5%AE%9E%E4%BD%93%E6%95%B0%E6%8D%AE%E6%A0%BC%E5%BC%8F "实体数据格式")
-   [方块实体数据格式](/w/%E6%96%B9%E5%9D%97%E5%AE%9E%E4%BD%93%E6%95%B0%E6%8D%AE%E6%A0%BC%E5%BC%8F "方块实体数据格式")
-   [物品格式](/w/%E7%89%A9%E5%93%81%E6%A0%BC%E5%BC%8F "物品格式")
-   [存档格式](/w/Java%E7%89%88%E5%AD%98%E6%A1%A3%E6%A0%BC%E5%BC%8F "Java版存档格式")
-   [世界生成](/w/%E4%B8%96%E7%95%8C%E7%94%9F%E6%88%90 "世界生成")
-   [数据生成器](/w/%E6%95%B0%E6%8D%AE%E7%94%9F%E6%88%90%E5%99%A8 "数据生成器")

相关教程

-   [安装数据包](/w/Tutorial:%E5%AE%89%E8%A3%85%E6%95%B0%E6%8D%AE%E5%8C%85 "Tutorial:安装数据包")
-   [制作数据包](/w/Tutorial:%E5%88%B6%E4%BD%9C%E6%95%B0%E6%8D%AE%E5%8C%85 "Tutorial:制作数据包")
-   [优化数据包](/w/Tutorial:%E4%BC%98%E5%8C%96%E6%95%B0%E6%8D%AE%E5%8C%85 "Tutorial:优化数据包")
-   [自定义盔甲纹饰](/w/Tutorial:%E8%87%AA%E5%AE%9A%E4%B9%89%E7%9B%94%E7%94%B2%E7%BA%B9%E9%A5%B0 "Tutorial:自定义盔甲纹饰")

参考实例

官方实例

-   [洞穴与山崖预览数据包](/w/%E6%B4%9E%E7%A9%B4%E4%B8%8E%E5%B1%B1%E5%B4%96%E9%A2%84%E8%A7%88%E6%95%B0%E6%8D%AE%E5%8C%85 "洞穴与山崖预览数据包")
-   [实验性内置数据包](/w/%E5%AE%9E%E9%AA%8C%E6%80%A7%E5%86%85%E5%AE%B9 "实验性内容")
-   [示例数据包](/w/%E7%A4%BA%E4%BE%8B%E6%95%B0%E6%8D%AE%E5%8C%85 "示例数据包")

教程实例

-   [实例：射线投射](/w/Tutorial:%E5%88%B6%E4%BD%9C%E6%95%B0%E6%8D%AE%E5%8C%85/%E5%AE%9E%E4%BE%8B%EF%BC%9A%E5%B0%84%E7%BA%BF%E6%8A%95%E5%B0%84 "Tutorial:制作数据包/实例：射线投射")
-   [实例：视线魔法](/w/Tutorial:%E5%88%B6%E4%BD%9C%E6%95%B0%E6%8D%AE%E5%8C%85/%E5%AE%9E%E4%BE%8B%EF%BC%9A%E8%A7%86%E7%BA%BF%E9%AD%94%E6%B3%95 "Tutorial:制作数据包/实例：视线魔法")

检索自“[https://zh.minecraft.wiki/w/文本组件?oldid=1315109](https://zh.minecraft.wiki/w/文本组件?oldid=1315109)”

[分类](/w/Special:Categories "Special:Categories")：​

-   [需要验证/基岩版](/w/Category:%E9%9C%80%E8%A6%81%E9%AA%8C%E8%AF%81/%E5%9F%BA%E5%B2%A9%E7%89%88 "Category:需要验证/基岩版")
-   [需要验证](/w/Category:%E9%9C%80%E8%A6%81%E9%AA%8C%E8%AF%81 "Category:需要验证")
-   [Java版](/w/Category:Java%E7%89%88 "Category:Java版")
-   [基岩版](/w/Category:%E5%9F%BA%E5%B2%A9%E7%89%88 "Category:基岩版")
-   [开发](/w/Category:%E5%BC%80%E5%8F%91 "Category:开发")

隐藏分类：​

-   [Java版即将到来/26.1](/w/Category:Java%E7%89%88%E5%8D%B3%E5%B0%86%E5%88%B0%E6%9D%A5/26.1 "Category:Java版即将到来/26.1")
-   [Java版即将移除/26.1](/w/Category:Java%E7%89%88%E5%8D%B3%E5%B0%86%E7%A7%BB%E9%99%A4/26.1 "Category:Java版即将移除/26.1")
-   [需要补充的条目](/w/Category:%E9%9C%80%E8%A6%81%E8%A1%A5%E5%85%85%E7%9A%84%E6%9D%A1%E7%9B%AE "Category:需要补充的条目")
-   [基岩版独有信息](/w/Category:%E5%9F%BA%E5%B2%A9%E7%89%88%E7%8B%AC%E6%9C%89%E4%BF%A1%E6%81%AF "Category:基岩版独有信息")

## 导航菜单

### 个人工具

-   [创建账号](/w/Special:CreateAccount?returnto=%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6&returntoquery=variant%3Dzh-cn "我们推荐您创建账号并登录，但这不是强制性的")
-   [登录](/w/Special:UserLogin?returnto=%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6&returntoquery=variant%3Dzh-cn "我们推荐您登录，但这不是强制性的​[o]")

### 命名空间

-   [页面](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6 "查看内容页面​[c]")
-   [讨论](/w/Talk:%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6 "有关内容页面的讨论​[t]")

 大陆简体

-   [不转换](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?variant=zh)
-   [简体](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?variant=zh-hans)
-   [繁體](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?variant=zh-hant)
-   [大陆简体](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?variant=zh-cn)
-   [香港繁體](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?variant=zh-hk)
-   [臺灣正體](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?variant=zh-tw)

### 查看

-   [阅读](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6)
-   [编辑](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?veaction=edit "编辑该页面​[v]")
-   [编辑源代码](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=edit "编辑该页面的源代码​[e]")
-   [查看历史](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=history "此页面过去的修订​[h]")

 更多

搜索

[](/ "访问首页")

### 导航

-   [参与贡献](/w/Help:%E5%8F%82%E4%B8%8E%E8%B4%A1%E7%8C%AE)
-   [最近更改](/w/Special:RecentChanges "本wiki的最近更改列表​[r]")
-   [随机页面](/w/Special:RandomRootpage "随机加载页面​[x]")
-   [在Minecraft中](/w/Special:RandomRootpage/Main)
-   [在Dungeons中](/w/Special:RandomRootpage/MCD)
-   [在Legends中](/w/Special:RandomRootpage/MCL)
-   [在Earth中](/w/Special:RandomRootpage/MCE)
-   [在Story Mode中](/w/Special:RandomRootpage/MCSM)
-   [在中国版中](/w/Special:RandomRootpage/MCCN)
-   [交流群](/w/Minecraft_Wiki:%E4%BA%A4%E6%B5%81%E7%BE%A4)

### 社区

-   [Wiki手册](/w/Help:%E7%BC%96%E8%BE%91%E6%89%8B%E5%86%8C)
-   [社区专页](/w/Minecraft_Wiki:%E7%A4%BE%E5%8C%BA%E4%B8%93%E9%A1%B5 "关于本项目，您可以做什么，以及何处能找到相关资料")
-   [管理员告示板](/w/Minecraft_Wiki:%E7%AE%A1%E7%90%86%E5%91%98%E5%91%8A%E7%A4%BA%E6%9D%BF)
-   [Wiki论坛](/w/Minecraft_Wiki:%E8%AE%BA%E5%9D%9B)
-   [Wiki条例](/w/Minecraft_Wiki:Wiki%E6%9D%A1%E4%BE%8B)
-   [标准译名列表](/w/Minecraft_Wiki:%E8%AF%91%E5%90%8D%E6%A0%87%E5%87%86%E5%8C%96)
-   [格式指导](/w/Minecraft_Wiki:%E6%A0%BC%E5%BC%8F%E6%8C%87%E5%AF%BC)
-   [计划](/w/Minecraft_Wiki:%E8%AE%A1%E5%88%92)
-   [沙盒](/w/Minecraft_Wiki:%E6%B2%99%E7%9B%92)
-   [草稿](/w/Help:%E8%8D%89%E7%A8%BF)

### 游戏及衍生作品

-   [Minecraft](/)
-   [Dungeons](/w/Dungeons:Wiki)
-   [Legends](/w/Legends:Wiki)
-   [Earth](/w/Earth:Wiki)
-   [Story Mode](/w/Story_Mode:Wiki)
-   [中国版](/w/China_Edition:Wiki)

### 版本

-   [Java版](/w/Java%E7%89%88)
-   [正式版：1.21.11](/w/Java%E7%89%881.21.11)
-   [开发版：26.1-pre-3](/w/Java%E7%89%8826.1-pre-3)
-   [基岩版](/w/%E5%9F%BA%E5%B2%A9%E7%89%88)
-   [正式版：26.3](/w/%E5%9F%BA%E5%B2%A9%E7%89%8826.3)
-   [测试版：26.20.20](/w/%E5%9F%BA%E5%B2%A9%E7%89%8826.20.20)

### 常用页面

-   [交易](/w/%E4%BA%A4%E6%98%93)
-   [药水酿造](/w/%E8%8D%AF%E6%B0%B4%E9%85%BF%E9%80%A0)
-   [附魔](/w/%E9%99%84%E9%AD%94)
-   [方块](/w/%E6%96%B9%E5%9D%97)
-   [物品](/w/%E7%89%A9%E5%93%81)
-   [生物](/w/%E7%94%9F%E7%89%A9)
-   [合成](/w/%E5%90%88%E6%88%90)
-   [烧炼](/w/%E7%83%A7%E7%82%BC)
-   [红石电路](/w/%E7%BA%A2%E7%9F%B3%E7%94%B5%E8%B7%AF)
-   [教程](/w/%E6%95%99%E7%A8%8B)
-   [资源包](/w/%E8%B5%84%E6%BA%90%E5%8C%85)

### 常用页面

-   [武器](/w/Dungeons:%E6%AD%A6%E5%99%A8)
-   [附魔](/w/Dungeons:%E9%99%84%E9%AD%94)
-   [盔甲](/w/Dungeons:%E7%9B%94%E7%94%B2)
-   [法器](/w/Dungeons:%E6%B3%95%E5%99%A8)
-   [地点](/w/Dungeons:%E5%9C%B0%E7%82%B9)

### 常用页面

-   [生物](/w/Legends:%E7%94%9F%E7%89%A9)
-   [资源](/w/Legends:%E8%B5%84%E6%BA%90)
-   [生物群系](/w/Legends:%E7%94%9F%E7%89%A9%E7%BE%A4%E7%B3%BB)
-   [建筑结构](/w/Legends:%E5%BB%BA%E7%AD%91%E7%BB%93%E6%9E%84)
-   [失落的传奇](/w/Legends:%E5%A4%B1%E8%90%BD%E7%9A%84%E4%BC%A0%E5%A5%87)

### 工具

-   [链入页面](/w/Special:WhatLinksHere/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6 "所有链接至本页面的wiki页面列表​[j]")
-   [相关更改](/w/Special:RecentChangesLinked/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6 "链自本页的页面的最近更改​[k]")
-   [可打印版](javascript:print\(\); "本页面的可打印版本​[p]")
-   [固定链接](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?oldid=1315109 "此页面该修订版本的固定链接")
-   [页面信息](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=info "关于此页面的更多信息")
-   [特殊页面](/w/Special:SpecialPages)
-   [查看存储桶](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?action=bucket "Bucket")

[](/hp/1773762603)

### 其他语言

-   [Deutsch](https://de.minecraft.wiki/w/JSON-Text "JSON-Text – Deutsch")
-   [English](https://minecraft.wiki/w/Text_component_format "Text component format – English")
-   [日本語](https://ja.minecraft.wiki/w/Raw_JSON%E3%83%86%E3%82%AD%E3%82%B9%E3%83%88%E3%83%95%E3%82%A9%E3%83%BC%E3%83%9E%E3%83%83%E3%83%88 "Raw JSONテキストフォーマット – 日本語")
-   [Português](https://pt.minecraft.wiki/w/Formato_de_texto_JSON_bruto "Formato de texto JSON bruto – português")

-   此页面最后编辑于2026年3月17日 (星期二) 15:32。
-   本网站内容采用[CC BY-NC-SA 3.0](https://creativecommons.org/licenses/by-nc-sa/3.0/)授权，[附加条款亦可能应用](https://meta.weirdgloop.org/w/Licensing "wgmeta:Licensing")。  
    本站并非Minecraft官方网站，与Mojang和微软亦无从属关系。

-   [隐私政策](https://weirdgloop.org/privacy)
-   [关于Minecraft Wiki](/w/Minecraft_Wiki:%E5%85%B3%E4%BA%8E)
-   [免责声明](https://meta.minecraft.wiki/w/General_disclaimer/zh)
-   [使用条款](https://weirdgloop.org/terms)
-   [联系Weird Gloop](/w/Special:Contact)
-   [移动版视图](https://zh.minecraft.wiki/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6?mobileaction=toggle_view_mobile&variant=zh-cn)

-   [![CC BY-NC-SA 3.0](https://meta.weirdgloop.org/images/Creative_Commons_footer.png)](https://creativecommons.org/licenses/by-nc-sa/3.0/)
-   [![Hosted by Weird Gloop](https://meta.weirdgloop.org/images/Weird_Gloop_footer_hosted.png)](https://weirdgloop.org)