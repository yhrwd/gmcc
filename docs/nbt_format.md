 NBT格式 - 中文 Minecraft Wiki      

                             

# NBT格式

来自Minecraft Wiki

[跳转到导航](#mw-head) [跳转到搜索](#searchInput)

![](/images/Disambig_gray.svg?1bb41)本条目介绍的是NBT网络传输和文件格式。关于用文本表示的NBT结构，请见“**[SNBT格式](/w/SNBT%E6%A0%BC%E5%BC%8F "SNBT格式")**”；关于在命令中检索特定NBT标签的方法，请见“**[NBT路径](/w/NBT%E8%B7%AF%E5%BE%84 "NBT路径")**”。

**NBT（Named Binary Tag，又称“二进制命名标签”）**是一种用带名称的二进制标签表示的树状数据结构。

## 目录

-   [1 概述](#概述)
    -   [1.1 常见游戏对象NBT](#常见游戏对象NBT)
-   [2 结构](#结构)
    -   [2.1 标签类型](#标签类型)
    -   [2.2 存储格式](#存储格式)
    -   [2.3 传输格式](#传输格式)
-   [3 转换](#转换)
    -   [3.1 程序对象](#程序对象)
    -   [3.2 SNBT](#SNBT)
    -   [3.3 JSON](#JSON)
-   [4 NBT程序对象](#NBT程序对象)
    -   [4.1 基于NBT修改对象](#基于NBT修改对象)
    -   [4.2 测试NBT标签](#测试NBT标签)
-   [5 历史](#历史)
-   [6 参考](#参考)
-   [7 导航](#导航)

## 概述

\[[编辑](/w/NBT%E6%A0%BC%E5%BC%8F?section=1&veaction=edit "编辑章节：概述") | [编辑源代码](/w/NBT%E6%A0%BC%E5%BC%8F?action=edit&section=1 "编辑章节的源代码： 概述")\]

NBT将游戏中的数据结构转换成通用的二进制数据包，以便存储或传输。游戏的大量数据文件都使用NBT书写；游戏在网络通讯时也常会传输NBT标签而非直接传输程序对象。

玩家在游戏中不会直接查看或修改NBT数据，而是在[命令](/w/%E5%91%BD%E4%BB%A4 "命令")中通过[SNBT](/w/SNBT "SNBT")这一文本形式的中介间接修改相应NBT数据。在游戏外，也可以用特制的程序直接查看或编辑NBT文件。

### 常见游戏对象NBT

\[[编辑](/w/NBT%E6%A0%BC%E5%BC%8F?section=2&veaction=edit "编辑章节：常见游戏对象NBT") | [编辑源代码](/w/NBT%E6%A0%BC%E5%BC%8F?action=edit&section=2 "编辑章节的源代码： 常见游戏对象NBT")\]

这里列出了游戏中可通过SNBT描述或允许被存储为NBT的常见游戏对象页面：

-   [方块实体数据格式](/w/%E6%96%B9%E5%9D%97%E5%AE%9E%E4%BD%93%E6%95%B0%E6%8D%AE%E6%A0%BC%E5%BC%8F "方块实体数据格式")
-   [实体数据格式](/w/%E5%AE%9E%E4%BD%93%E6%95%B0%E6%8D%AE%E6%A0%BC%E5%BC%8F "实体数据格式")
-   [物品格式](/w/%E7%89%A9%E5%93%81%E6%A0%BC%E5%BC%8F "物品格式")

## 结构

\[[编辑](/w/NBT%E6%A0%BC%E5%BC%8F?section=3&veaction=edit "编辑章节：结构") | [编辑源代码](/w/NBT%E6%A0%BC%E5%BC%8F?action=edit&section=3 "编辑章节的源代码： 结构")\]

NBT的基本结构是带名称的**标签**。标签可以大致分为数值标签和结构标签；前者表示单个值，后者则容纳其他标签，通过嵌套构建树状数据结构。

除了结束标签以外，每个标签都由**标签ID**、**标签名称**和**[负载](https://zh.wikipedia.org/wiki/%E8%B4%9F%E8%BD%BD_\(%E8%AE%A1%E7%AE%97%E6%9C%BA\) "wzh:负载 (计算机)")**组成。ID是表示该标签类型的字节；名称是一个带长字符串，包含一个按大端序\[仅[JE](/w/Java%E7%89%88 "Java版")\]/小端序\[仅[BE](/w/%E5%9F%BA%E5%B2%A9%E7%89%88 "基岩版")\]存储的无符号短整型\[[需要更多信息](/w/Special:TalkPage/NBT%E6%A0%BC%E5%BC%8F "Special:TalkPage/NBT格式")\]/有符号短整型\[仅[BE](/w/%E5%9F%BA%E5%B2%A9%E7%89%88 "基岩版")\]长度，以及按Java使用的[变种UTF-8（Modified UTF-8）编码](https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/io/DataInput.html#modified-utf-8)书写的名称；负载表示该标签所承载的数据，具体布局因标签而异。

带名称标签的布局如下：

组分

标签ID

标签名称

负载

字节偏移

0

1~2

3~2+L

3+L~?

数据

标签ID

名称长度

名称（变种UTF-8）

*见后*

结束标签不含名称与负载，仅由单字节ID构成。

### 标签类型

\[[编辑](/w/NBT%E6%A0%BC%E5%BC%8F?section=4&veaction=edit "编辑章节：标签类型") | [编辑源代码](/w/NBT%E6%A0%BC%E5%BC%8F?action=edit&section=4 "编辑章节的源代码： 标签类型")\]

NBT共有13\[仅[JE](/w/Java%E7%89%88 "Java版")\] 或 12\[仅[BE](/w/%E5%9F%BA%E5%B2%A9%E7%89%88 "基岩版")\]种标签类型（计入结束标签）。详见下表。

ID

ID (HEX)

标签类型

负载

负载长度（字节）

0

00

结束标签

–

–

1

01

![字节型](/images/Data_node_byte.svg?eb0e0)字节型

单个有符号字节

1

2

02

![短整型](/images/Data_node_short.svg?c1f72)短整型

单个有符号短整型，大端序\[仅[JE](/w/Java%E7%89%88 "Java版")\]/小端序\[仅[BE](/w/%E5%9F%BA%E5%B2%A9%E7%89%88 "基岩版")\]

2

3

03

![整型](/images/Data_node_int.svg?8d24f)整型

单个有符号整型，大端序\[仅[JE](/w/Java%E7%89%88 "Java版")\]/小端序\[仅[BE](/w/%E5%9F%BA%E5%B2%A9%E7%89%88 "基岩版")\]

4

4

04

![长整型](/images/Data_node_long.svg?dde3f)长整型

单个有符号长整型，大端序\[仅[JE](/w/Java%E7%89%88 "Java版")\]/小端序\[仅[BE](/w/%E5%9F%BA%E5%B2%A9%E7%89%88 "基岩版")\]

8

5

05

![单精度浮点数](/images/Data_node_float.svg?ae55e)单精度浮点型

单个单精度浮点数（IEEE754），大端序\[仅[JE](/w/Java%E7%89%88 "Java版")\]/小端序\[仅[BE](/w/%E5%9F%BA%E5%B2%A9%E7%89%88 "基岩版")\]

4

6

06

![双精度浮点数](/images/Data_node_double.svg?14320)双精度浮点型

单个双精度浮点数（IEEE754），大端序\[仅[JE](/w/Java%E7%89%88 "Java版")\]/小端序\[仅[BE](/w/%E5%9F%BA%E5%B2%A9%E7%89%88 "基岩版")\]

8

7

07

![字节型数组](/images/Data_node_byte-array.svg?2b418)字节数组

按大端序\[仅[JE](/w/Java%E7%89%88 "Java版")\]/小端序\[仅[BE](/w/%E5%9F%BA%E5%B2%A9%E7%89%88 "基岩版")\]存储的有符号整型长度，后跟相应个数的有符号字节

4+<长度>

8

08

![字符串](/images/Data_node_string.svg?42545)字符串标签

一个带长字符串，包含按大端序\[仅[JE](/w/Java%E7%89%88 "Java版")\]/小端序\[仅[BE](/w/%E5%9F%BA%E5%B2%A9%E7%89%88 "基岩版")\]存储的无符号短整型\[[需要更多信息](/w/Special:TalkPage/NBT%E6%A0%BC%E5%BC%8F "Special:TalkPage/NBT格式")\]/有符号短整型\[仅[BE](/w/%E5%9F%BA%E5%B2%A9%E7%89%88 "基岩版")\]长度。字符串是按[变种UTF-8编码](https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/io/DataInput.html#modified-utf-8)\[仅[JE](/w/Java%E7%89%88 "Java版")\]或标准UTF-8编码\[仅[BE](/w/%E5%9F%BA%E5%B2%A9%E7%89%88 "基岩版")\]书写的

2+<长度>

9

09

![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)列表

列表中标签的标签ID，后跟按大端序\[仅[JE](/w/Java%E7%89%88 "Java版")\]/小端序\[仅[BE](/w/%E5%9F%BA%E5%B2%A9%E7%89%88 "基岩版")\]存储的有符号整型长度，后跟相应个数的子标签**负载**

5+<长度1>+<长度2>+...+<长度n>（n为包含的子标签数）

10

0A

![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)复合标签

若干个子标签**整体**（ID、名称、负载），后跟空字节（可视为结束标签）

?+1

11

0B

![整型数组](/images/Data_node_int-array.svg?546e8)整型数组

按大端序\[仅[JE](/w/Java%E7%89%88 "Java版")\]/小端序\[仅[BE](/w/%E5%9F%BA%E5%B2%A9%E7%89%88 "基岩版")\]存储的有符号整型长度，后跟相应个数的有符号整型

4+<长度>×4

12

0C

![长整型数组](/images/Data_node_long-array.svg?92504)长整型数组\[仅[JE](/w/Java%E7%89%88 "Java版")\]

按大端序存储的有符号整型长度，后跟相应个数的有符号长整型

4+<长度>×8

### 存储格式

\[[编辑](/w/NBT%E6%A0%BC%E5%BC%8F?section=5&veaction=edit "编辑章节：存储格式") | [编辑源代码](/w/NBT%E6%A0%BC%E5%BC%8F?action=edit&section=5 "编辑章节的源代码： 存储格式")\]

NBT文件包含单个![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)复合标签或![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)列表\[仅[基岩版](/w/%E5%9F%BA%E5%B2%A9%E7%89%88 "基岩版")\]（“根标签”），其包含单个子标签。该文件所存储的数据结构即在该子标签中。在NBT文件中，根标签可能未压缩、用GZip压缩或用Zlib压缩\[仅[基岩版](/w/%E5%9F%BA%E5%B2%A9%E7%89%88 "基岩版")\]。具体采用何种压缩方式与数据类型有关，于相应数据类型的页面中注明。

### 传输格式

\[[编辑](/w/NBT%E6%A0%BC%E5%BC%8F?section=6&veaction=edit "编辑章节：传输格式") | [编辑源代码](/w/NBT%E6%A0%BC%E5%BC%8F?action=edit&section=6 "编辑章节的源代码： 传输格式")\]

在网络上传输时，NBT以流式传输而不被压缩。传输的NBT数据形式上是单个NBT标签（“根标签”），为一个仅含单个子标签的![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)复合标签。

在[Java版](/w/Java%E7%89%88 "Java版")中，传输格式的NBT中，根标签的名称（包含长度和字符串）将被省略，也就是根标签的标签ID后直接跟负载。

在[基岩版](/w/%E5%9F%BA%E5%B2%A9%E7%89%88 "基岩版")中，传输格式的NBT会使用变长数字格式（Varint）及其变种编码标签中的所有数字，包括整型与浮点数标签，以及列表、数组、复合标签与字符串的长度字段。\[[需要更多信息](/w/Special:TalkPage/NBT%E6%A0%BC%E5%BC%8F "Special:TalkPage/NBT格式")\]

## 转换

\[[编辑](/w/NBT%E6%A0%BC%E5%BC%8F?section=7&veaction=edit "编辑章节：转换") | [编辑源代码](/w/NBT%E6%A0%BC%E5%BC%8F?action=edit&section=7 "编辑章节的源代码： 转换")\]

### 程序对象

\[[编辑](/w/NBT%E6%A0%BC%E5%BC%8F?section=8&veaction=edit "编辑章节：程序对象") | [编辑源代码](/w/NBT%E6%A0%BC%E5%BC%8F?action=edit&section=8 "编辑章节的源代码： 程序对象")\]

游戏在运行时，有关数据直接存储在程序对象中，而非NBT结构中。在需要传输或存储时，以及玩家用[命令](/w/%E5%91%BD%E4%BB%A4 "命令")修改时，游戏会根据相应数据的**编码格式**选择性地从程序对象生成相应的NBT结构，或从NBT重新构建相应程序对象。这种转换有复杂的规则，例如某些数据有意不写入NBT，或者在NBT中采用其他表示方法；这些规则因数据类型而异，在相应数据格式的页面下有所介绍。

### SNBT

\[[编辑](/w/NBT%E6%A0%BC%E5%BC%8F?section=9&veaction=edit "编辑章节：SNBT") | [编辑源代码](/w/NBT%E6%A0%BC%E5%BC%8F?action=edit&section=9 "编辑章节的源代码： SNBT")\]

[SNBT](/w/SNBT "SNBT")作为一种基于文本的数据结构，是NBT与玩家的中介。在使用部分[命令](/w/%E5%91%BD%E4%BB%A4 "命令")修改、查看数据，或使用[客户端核心文件](/w/%E5%AE%A2%E6%88%B7%E7%AB%AF%E6%A0%B8%E5%BF%83%E6%96%87%E4%BB%B6 "客户端核心文件")中的数据生成器时，都涉及NBT与SNBT的转换。

结构模板有特殊转换规则，见[结构存储格式](/w/%E7%BB%93%E6%9E%84%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "结构存储格式")。

**SNBT转NBT**

SNBT包含了NBT的各种标签类型，但还包含了一些不直接属于NBT的标签类型。下面给出各种标签的转换规则。

转换规则表

SNBT标签

NBT表示

![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)复合标签  
![字节型数组](/images/Data_node_byte-array.svg?2b418)字节数组  
![整型数组](/images/Data_node_int-array.svg?546e8)整型数组  
![长整型数组](/images/Data_node_long-array.svg?92504)长整型数组  
有符号![字节型](/images/Data_node_byte.svg?eb0e0)字节型  
有符号![短整型](/images/Data_node_short.svg?c1f72)短整型  
有符号![整型](/images/Data_node_int.svg?8d24f)整型  
有符号![长整型](/images/Data_node_long.svg?dde3f)长整型  
![单精度浮点数](/images/Data_node_float.svg?ae55e)单精度浮点数  
![双精度浮点数](/images/Data_node_double.svg?14320)双精度浮点数

直接转换至相应NBT标签

无符号![字节型](/images/Data_node_byte.svg?eb0e0)字节型  
无符号![短整型](/images/Data_node_short.svg?c1f72)短整型  
无符号![整型](/images/Data_node_int.svg?8d24f)整型  
无符号![长整型](/images/Data_node_long.svg?dde3f)长整型

按补码解释，将无符号数转换成同种类型的有符号数标签

![布尔型](/images/Data_node_bool.svg?77754)布尔型

true转换为字节型1b，false转换为字节型0b

![字符串](/images/Data_node_string.svg?42545)字符串

执行所有转义后，按Java使用的[变种UTF-8编码](https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/io/DataInput.html#modified-utf-8)存储

![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)列表

如列表元素类型相同，转换为元素类型为子标签唯一类型的列表。否则，转换为元素类型为复合标签的列表；列表元素中复合标签元素原样存储，非复合标签元素分别包装在`{"": <值>}`标签中存储

**NBT转SNBT**

NBT标签在SNBT中几乎都有对应物。在同一标签的多种表示形式中，转换得到的SNBT通常会选择固定的一种。

转换规则表

NBT标签

SNBT表示

备注

![字节型](/images/Data_node_byte.svg?eb0e0)字节型

`<值>b`

`<值>`以十进制表示，不含下划线。

![短整型](/images/Data_node_short.svg?c1f72)短整型

`<值>s`

![整型](/images/Data_node_int.svg?8d24f)整型

`<值>`

![长整型](/images/Data_node_long.svg?dde3f)长整型

`<值>l`

![单精度浮点数](/images/Data_node_float.svg?ae55e)单精度浮点型

`<值>f`

`<值>`通过[Float.toString](https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/lang/Float.html#toString\(float\))和[Double.toString](https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/lang/Double.html#toString\(double\))方法转换为文本，可能为整数、小数或科学计数法。无穷数和NaN也会转换为相应的文本表示（`Infinity`和`NaN`），但不再能被SNBT推断为浮点型。[\[1\]](#cite_note-1)

![双精度浮点数](/images/Data_node_double.svg?14320)双精度浮点型

`<值>d`

![字节型数组](/images/Data_node_byte-array.svg?2b418)字节数组

`[B;<值>B,...]`

`<值>`的转换与整型的规则一致。元素列表末尾不附加逗号。

![整型数组](/images/Data_node_int-array.svg?546e8)整型数组

`[I;<值>,...]`

![长整型数组](/images/Data_node_long-array.svg?92504)长整型数组

`[L;<值>L,...]`

![字符串](/images/Data_node_string.svg?42545)字符串

`"<值>"`或`'<值>'`

如文本中首个引号为双引号则使用单引号字符串，否则均使用双引号字符串。`<值>`中的反斜杠、同种引号、有[固定转义序列](/w/SNBT%E6%A0%BC%E5%BC%8F#固定转义 "SNBT格式")的控制符以及码位在U+0000到U+001F的字符（使用\\x转义）均会被转义。

![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)列表

`[<值>,...]`

`<值>`按相应标签的规则转换。元素列表末尾不附加逗号。

![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)复合标签

`{<键>:<值>,...}`

元素依键按字典序升序排列。`<键>`如能用无引号字符串表示（含忽略大小写的true或false），则会转换为无引号形式，否则按字符串标签的规则转换。`<值>`按相应标签的规则转换。元素列表末尾不附加逗号。

**结束标签**

`END`

仅在不与![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)复合标签配对时出现。在聊天栏显示中不会输出。

在聊天栏输出中，上述SNBT会进一步转换为带语法着色的[文本组件](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6 "文本组件")，且会用`<...>`截断过长的列表、数组和复合标签。

### JSON

\[[编辑](/w/NBT%E6%A0%BC%E5%BC%8F?section=10&veaction=edit "编辑章节：JSON") | [编辑源代码](/w/NBT%E6%A0%BC%E5%BC%8F?action=edit&section=10 "编辑章节的源代码： JSON")\]

[JSON](/w/JSON "JSON")格式与NBT和SNBT完全不同；由于语法和基本类型不同，JSON与NBT并不兼容。要精确地在NBT中嵌入JSON对象，通常需要将JSON文本放在字符串里传入。

游戏的部分数据，例如[生物群系数据格式](/w/%E7%94%9F%E7%89%A9%E7%BE%A4%E7%B3%BB%E6%95%B0%E6%8D%AE%E6%A0%BC%E5%BC%8F "生物群系数据格式")，以JSON格式存储，却需要转换为NBT格式。此时，游戏会试图在NBT与JSON之间转换，但这一过程会损失部分信息。

**JSON转NBT**

由于存在`null`值以及元素类型不同的列表，JSON转NBT不一定能够成功。

转换规则表

JSON类型

NBT标签

备注

JsonString

![字符串](/images/Data_node_string.svg?42545)字符串

JsonBoolean

![字节型](/images/Data_node_byte.svg?eb0e0)字节型

true转换为1b，false转换为0b

JsonNumber

![字节型](/images/Data_node_byte.svg?eb0e0)字节型  
![短整型](/images/Data_node_short.svg?c1f72)短整型  
![整型](/images/Data_node_int.svg?8d24f)整型  
![长整型](/images/Data_node_long.svg?dde3f)长整型  
![单精度浮点数](/images/Data_node_float.svg?ae55e)单精度浮点型  
![双精度浮点数](/images/Data_node_double.svg?14320)双精度浮点型

按照数值范围从小到大（即所列顺序）依次尝试，如在该类型可表达的数值范围内即转换为该类型

JsonNull

无法转换

JsonArray

![字节型数组](/images/Data_node_byte-array.svg?2b418)字节型数组  
![整型数组](/images/Data_node_int-array.svg?546e8)整型数组  
![长整型数组](/images/Data_node_long-array.svg?92504)长整型数组  
![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)列表

先将列表中的元素按各自规则转换。随后，如元素转换得到的类型不同（即便相互兼容），则无法转换；否则，优先尝试转换成相应类型的数组，如不存在类型匹配的数组则转换成列表

JsonObject

![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)复合标签

**NBT转JSON**

NBT所支持的类型基本都可转换到JSON。但由于JSON的数字类型单一，转换过程会丢失NBT的数字类型信息。

转换规则表

NBT标签

JSON类型

![字节型](/images/Data_node_byte.svg?eb0e0)字节型  
![短整型](/images/Data_node_short.svg?c1f72)短整型  
![整型](/images/Data_node_int.svg?8d24f)整型  
![长整型](/images/Data_node_long.svg?dde3f)长整型  
![单精度浮点数](/images/Data_node_float.svg?ae55e)单精度浮点型  
![双精度浮点数](/images/Data_node_double.svg?14320)双精度浮点型

JsonNumber

![字节型数组](/images/Data_node_byte-array.svg?2b418)字节型数组  
![整型数组](/images/Data_node_int-array.svg?546e8)整型数组  
![长整型数组](/images/Data_node_long-array.svg?92504)长整型数组  
![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)列表

JsonArray

![字符串](/images/Data_node_string.svg?42545)字符串

JsonString

![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)复合标签

JsonObject

## NBT程序对象

\[[编辑](/w/NBT%E6%A0%BC%E5%BC%8F?section=11&veaction=edit "编辑章节：NBT程序对象") | [编辑源代码](/w/NBT%E6%A0%BC%E5%BC%8F?action=edit&section=11 "编辑章节的源代码： NBT程序对象")\]

游戏在运行时其数据以程序对象而非NBT格式存在，通常游戏只会在构建程序对象或保存程序对象时进行转换，例如存档的读写，刷怪笼生成实体等。

在[Java版](/w/Java%E7%89%88 "Java版")中，可以通过命令等方式在游戏内主动进行程序对象与NBT格式的转换，例如命令`/[data](/w/%E5%91%BD%E4%BB%A4/data "命令/data")`和[目标选择器](/w/%E7%9B%AE%E6%A0%87%E9%80%89%E6%8B%A9%E5%99%A8 "目标选择器")等。

### 基于NBT修改对象

\[[编辑](/w/NBT%E6%A0%BC%E5%BC%8F?section=12&veaction=edit "编辑章节：基于NBT修改对象") | [编辑源代码](/w/NBT%E6%A0%BC%E5%BC%8F?action=edit&section=12 "编辑章节的源代码： 基于NBT修改对象")\]

[![](/images/Information_icon.svg?eefcf)](/w/File:Information_icon.svg)

**本段落所述内容仅适用于[Java版](/w/Java%E7%89%88 "Java版")。**

在修改方块实体或实体等对象之前，游戏将传入的SNBT或JSON等转换为NBT实例后再进行修改。只有某些特定属性可以被传入的NBT修改，例如方块实体的坐标不可修改。传入不被程序对象使用的NBT属性会被丢弃，例如实体不会保存`nonExist`标签。向具有特定编码格式要求的NBT属性传入错误的数据会报错，例如实体自定义名称需要[文本组件](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6 "文本组件")，尽管它以复合标签存储传输，但传入空标签会报错。

若某属性需要一个[命名空间ID](/w/%E5%91%BD%E5%90%8D%E7%A9%BA%E9%97%B4ID "命名空间ID")，但传入了一个不带命名空间前缀的字符串，则该值将按照[字符串到命名空间ID的转换](/w/%E5%91%BD%E5%90%8D%E7%A9%BA%E9%97%B4ID#字符串转换为命名空间ID "命名空间ID")规则进行转换。

若某属性需要一个布尔值，但传入了数值类型的值，那么该值将则向下取整转换为字节型数字，若不是`0b`则为`1b`。

若某属性需要一个布尔值，但传入了非数值也非布尔类型的值，那么此属性值为`0b`。

若某属性需要某数值类型的数字作为其值，但传入了一个与之所需类型不符的数值类型标签，那么需要转换为所需要的类型。如果属性需要整型，则先向下取整后再转换。

若某属性需要某数值类型的值，但传入了一个非数值类型的值，那么该属性将被赋值为0（具体数据类型与属性而异）。

若某属性需要一个字符串，但传入了一个非字符串，那么该属性值将会是一个空字符串。

若某属性需要一个列表或是某类型的数组，但传入的标签类型与需要的类型不符，则该属性值将为一个空列表或空数组。

若某属性需要一个复合标签，但传入了一个非复合标签，则该属性值最终为一个空复合标签。

### 测试NBT标签

\[[编辑](/w/NBT%E6%A0%BC%E5%BC%8F?section=13&veaction=edit "编辑章节：测试NBT标签") | [编辑源代码](/w/NBT%E6%A0%BC%E5%BC%8F?action=edit&section=13 "编辑章节的源代码： 测试NBT标签")\]

[![](/images/Information_icon.svg?eefcf)](/w/File:Information_icon.svg)

**本段落所述内容仅适用于[Java版](/w/Java%E7%89%88 "Java版")。**

当游戏需要测试NBT标签，例如在[目标选择器中使用NBT标签](/w/%E7%9B%AE%E6%A0%87%E9%80%89%E6%8B%A9%E5%99%A8#实体的NBT标签 "目标选择器")筛选实体时，游戏会将传入的SNBT等转换为NBT对象，再从待测试的目标对象还原出另一个NBT对象，然后测试目标对象的NBT标签中是否具有提供的NBT标签。

只要提供的标签确实存在于目标对象中，测试就会成功，也即无需真正完整提供目标对象所有的NBT属性。对于列表的测试也如此，只要提供的列表元素全部存在于目标对象的列表中，即使提供的列表中的元素顺序和元素数量不符合，游戏也会测试成功。

例如，对于拥有NBT标签`Pos: [1d, 2d, 3d], data: {tag1: {name: test}}`的实体而言：

-   `@e[nbt={data: {}}]`可以选中，因为目标实体存在复合标签`data.tag1`。
-   `@e[nbt={Pos: [2d, 3d, 1d]}]`可以选中，因为列表匹配不考虑元素顺序。
-   `@e[nbt={Pos: [1d]}]`可以选中，因为目标实体的此列表中存在元素`1d`。
-   `@e[nbt={Pos: []}]`无法选中，因为**空列表只能匹配空列表**。

在测试NBT标签的具体值时，提供的标签名称及数据类型必须与目标对象的相应标签及数据类型**完全匹配**，否则此测试无效。数组也在此范畴内，因此不能像列表一样忽略内部元素的个数和顺序。例如`1d`不会被`1`匹配，因为双精度浮点数无法与整数匹配；`[L; 1L, 3L]`不会被`[L; 3L]`或`[L; 3L, 1L]`匹配，因为数组要求完全匹配。

SNBT向NBT的[§ 转换](#转换)依然会进行。这意味着若测试值是否是`1b`，在命令中提供`true`或`1ub`都会测试成功。但由于游戏会各自整理提供的NBT对象和目标对象的NBT对象，因此并不会执行上述[§ 基于NBT修改对象](#基于NBT修改对象)的转换，因为这些转换需要目标对象的编码格式。一个典型的例子是测试命名空间ID，即使目标对象的命名空间ID使用`minecraft`命名空间，也必须在命令中提供此命名空间，因为游戏不会自动将字符串转换为符合命名空间ID格式的字符串，例如石头的物品实体不会被`@e[nbt={Item: {id: stone}}]`选中而会被`@e[nbt={Item: {id: "minecraft:stone"}}]`选中。

## 历史

\[[编辑](/w/NBT%E6%A0%BC%E5%BC%8F?section=14&veaction=edit "编辑章节：历史") | [编辑源代码](/w/NBT%E6%A0%BC%E5%BC%8F?action=edit&section=14 "编辑章节的源代码： 历史")\]

[2011年9月28日](https://twitter.com/notch/status/119296531592515584)

Notch致力于“用物品实例来存储任意信息”。

[Java版](/w/Java%E7%89%88%E7%89%88%E6%9C%AC%E8%AE%B0%E5%BD%95 "Java版版本记录")

?

加入了NBT格式。

[1.12](/w/Java%E7%89%881.12 "Java版1.12")

?

加入了![长整型数组](/images/Data_node_long-array.svg?92504)长整型数组标签。

[1.16](/w/Java%E7%89%881.16 "Java版1.16")

[20w21a](/w/Java%E7%89%8820w21a "Java版20w21a")

加入了NBT格式和JSON格式文件的转换功能。

[携带版](/w/%E6%90%BA%E5%B8%A6%E7%89%88 "携带版")

[0.2.0](/w/%E6%90%BA%E5%B8%A6%E7%89%880.2.0 "携带版0.2.0")

加入了NBT格式。

## 参考

\[[编辑](/w/NBT%E6%A0%BC%E5%BC%8F?section=15&veaction=edit "编辑章节：参考") | [编辑源代码](/w/NBT%E6%A0%BC%E5%BC%8F?action=edit&section=15 "编辑章节的源代码： 参考")\]

1.  [↑](#cite_ref-1) [MC-200070](https://bugs.mojang.com/browse/MC-200070 "mojira:MC-200070")

## 导航

\[[编辑](/w/NBT%E6%A0%BC%E5%BC%8F?section=16&veaction=edit "编辑章节：导航") | [编辑源代码](/w/NBT%E6%A0%BC%E5%BC%8F?action=edit&section=16 "编辑章节的源代码： 导航")\]

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

-   [文本组件](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6 "文本组件")
-   NBT格式
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
-   NBT格式
-   [动画](/w/%E5%8A%A8%E7%94%BB "动画")
-   [附加包](/w/%E9%99%84%E5%8A%A0%E5%8C%85 "附加包")
    -   [Molang](/w/Molang "Molang")
-   [游戏测试](/w/%E6%B8%B8%E6%88%8F%E6%B5%8B%E8%AF%95 "游戏测试")
-   [配方](/w/%E9%85%8D%E6%96%B9 "配方")
-   [方块实体](/w/%E6%96%B9%E5%9D%97%E5%AE%9E%E4%BD%93 "方块实体")
-   [命名空间ID](/w/%E5%91%BD%E5%90%8D%E7%A9%BA%E9%97%B4ID "命名空间ID")
-   [文本组件](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6 "文本组件")
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
-   NBT格式
    -   [NBT路径](/w/NBT%E8%B7%AF%E5%BE%84 "NBT路径")
-   [SNBT格式](/w/SNBT%E6%A0%BC%E5%BC%8F "SNBT格式")
-   [JSON](/w/JSON "JSON")
-   [文本组件](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6 "文本组件")
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

检索自“[https://zh.minecraft.wiki/w/NBT格式?oldid=1309757](https://zh.minecraft.wiki/w/NBT格式?oldid=1309757)”

[分类](/w/Special:Categories "Special:Categories")：​

-   [需要信息](/w/Category:%E9%9C%80%E8%A6%81%E4%BF%A1%E6%81%AF "Category:需要信息")
-   [Java版](/w/Category:Java%E7%89%88 "Category:Java版")
-   [基岩版](/w/Category:%E5%9F%BA%E5%B2%A9%E7%89%88 "Category:基岩版")

隐藏分类：​

-   [Java版独有信息](/w/Category:Java%E7%89%88%E7%8B%AC%E6%9C%89%E4%BF%A1%E6%81%AF "Category:Java版独有信息")
-   [基岩版独有信息](/w/Category:%E5%9F%BA%E5%B2%A9%E7%89%88%E7%8B%AC%E6%9C%89%E4%BF%A1%E6%81%AF "Category:基岩版独有信息")
-   [有未知版本的History模板的页面/Java版](/w/Category:%E6%9C%89%E6%9C%AA%E7%9F%A5%E7%89%88%E6%9C%AC%E7%9A%84History%E6%A8%A1%E6%9D%BF%E7%9A%84%E9%A1%B5%E9%9D%A2/Java%E7%89%88 "Category:有未知版本的History模板的页面/Java版")

## 导航菜单

### 个人工具

-   [创建账号](/w/Special:CreateAccount?returnto=NBT%E6%A0%BC%E5%BC%8F&returntoquery=variant%3Dzh-cn "我们推荐您创建账号并登录，但这不是强制性的")
-   [登录](/w/Special:UserLogin?returnto=NBT%E6%A0%BC%E5%BC%8F&returntoquery=variant%3Dzh-cn "我们推荐您登录，但这不是强制性的​[o]")

### 命名空间

-   [页面](/w/NBT%E6%A0%BC%E5%BC%8F "查看内容页面​[c]")
-   [讨论](/w/Talk:NBT%E6%A0%BC%E5%BC%8F "有关内容页面的讨论​[t]")

 大陆简体

-   [不转换](/w/NBT%E6%A0%BC%E5%BC%8F?variant=zh)
-   [简体](/w/NBT%E6%A0%BC%E5%BC%8F?variant=zh-hans)
-   [繁體](/w/NBT%E6%A0%BC%E5%BC%8F?variant=zh-hant)
-   [大陆简体](/w/NBT%E6%A0%BC%E5%BC%8F?variant=zh-cn)
-   [香港繁體](/w/NBT%E6%A0%BC%E5%BC%8F?variant=zh-hk)
-   [臺灣正體](/w/NBT%E6%A0%BC%E5%BC%8F?variant=zh-tw)

### 查看

-   [阅读](/w/NBT%E6%A0%BC%E5%BC%8F)
-   [编辑](/w/NBT%E6%A0%BC%E5%BC%8F?veaction=edit "编辑该页面​[v]")
-   [编辑源代码](/w/NBT%E6%A0%BC%E5%BC%8F?action=edit "编辑该页面的源代码​[e]")
-   [查看历史](/w/NBT%E6%A0%BC%E5%BC%8F?action=history "此页面过去的修订​[h]")

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

-   [链入页面](/w/Special:WhatLinksHere/NBT%E6%A0%BC%E5%BC%8F "所有链接至本页面的wiki页面列表​[j]")
-   [相关更改](/w/Special:RecentChangesLinked/NBT%E6%A0%BC%E5%BC%8F "链自本页的页面的最近更改​[k]")
-   [可打印版](javascript:print\(\); "本页面的可打印版本​[p]")
-   [固定链接](/w/NBT%E6%A0%BC%E5%BC%8F?oldid=1309757 "此页面该修订版本的固定链接")
-   [页面信息](/w/NBT%E6%A0%BC%E5%BC%8F?action=info "关于此页面的更多信息")
-   [特殊页面](/w/Special:SpecialPages)
-   [查看存储桶](/w/NBT%E6%A0%BC%E5%BC%8F?action=bucket "Bucket")

[](/hp/1773763518)

### 其他语言

-   [Deutsch](https://de.minecraft.wiki/w/NBT "NBT – Deutsch")
-   [English](https://minecraft.wiki/w/NBT_format "NBT format – English")
-   [Español](https://es.minecraft.wiki/w/Formato_NBT "Formato NBT – español")
-   [Français](https://fr.minecraft.wiki/w/Format_NBT "Format NBT – français")
-   [日本語](https://ja.minecraft.wiki/w/NBT%E3%83%95%E3%82%A9%E3%83%BC%E3%83%9E%E3%83%83%E3%83%88 "NBTフォーマット – 日本語")
-   [Nederlands](https://nl.minecraft.wiki/w/NBT_formaat "NBT formaat – Nederlands")
-   [Português](https://pt.minecraft.wiki/w/Formato_NBT "Formato NBT – português")
-   [Русский](https://ru.minecraft.wiki/w/%D0%A4%D0%BE%D1%80%D0%BC%D0%B0%D1%82_NBT "Формат NBT – русский")

-   此页面最后编辑于2026年3月8日 (星期日) 08:04。
-   本网站内容采用[CC BY-NC-SA 3.0](https://creativecommons.org/licenses/by-nc-sa/3.0/)授权，[附加条款亦可能应用](https://meta.weirdgloop.org/w/Licensing "wgmeta:Licensing")。  
    本站并非Minecraft官方网站，与Mojang和微软亦无从属关系。

-   [隐私政策](https://weirdgloop.org/privacy)
-   [关于Minecraft Wiki](/w/Minecraft_Wiki:%E5%85%B3%E4%BA%8E)
-   [免责声明](https://meta.minecraft.wiki/w/General_disclaimer/zh)
-   [使用条款](https://weirdgloop.org/terms)
-   [联系Weird Gloop](/w/Special:Contact)
-   [移动版视图](https://zh.minecraft.wiki/w/NBT%E6%A0%BC%E5%BC%8F?mobileaction=toggle_view_mobile&variant=zh-cn)

-   [![CC BY-NC-SA 3.0](https://meta.weirdgloop.org/images/Creative_Commons_footer.png)](https://creativecommons.org/licenses/by-nc-sa/3.0/)
-   [![Hosted by Weird Gloop](https://meta.weirdgloop.org/images/Weird_Gloop_footer_hosted.png)](https://weirdgloop.org)