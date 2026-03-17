 数据组件 - 中文 Minecraft Wiki      

                             

# 数据组件

来自Minecraft Wiki

[跳转到导航](#mw-head) [跳转到搜索](#searchInput)

![](/images/Disambig_gray.svg?1bb41)关于基岩版中的物品堆叠组件，请见“**[基岩版物品堆叠组件](/w/%E5%9F%BA%E5%B2%A9%E7%89%88%E7%89%A9%E5%93%81%E5%A0%86%E5%8F%A0%E7%BB%84%E4%BB%B6 "基岩版物品堆叠组件")**”。

![](/images/Disambig_gray.svg?1bb41)“**组件**”重定向至此。关于其他用法，请见“**[组件（消歧义）](/w/%E7%BB%84%E4%BB%B6%EF%BC%88%E6%B6%88%E6%AD%A7%E4%B9%89%EF%BC%89 "组件（消歧义）")**”。

![](/images/Disambig_gray.svg?1bb41)“**元件**”重定向至此。关于建造红石电路的元件，请见“**[红石元件](/w/%E7%BA%A2%E7%9F%B3%E5%85%83%E4%BB%B6 "红石元件")**”。

[![](/images/thumb/Iron_Shovel_JE2_BE2.png/32px-Iron_Shovel_JE2_BE2.png?fadf7)](/w/File:Iron_Shovel_JE2_BE2.png)

**该页面正在草稿中编辑。**

由于页面过旧或未完全翻译等原因，此页面的内容目前位于[草稿](https://zh.minecraft.wiki/w/Draft:%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6)中。  
请移步至草稿查看或编辑最近的版本。  

[![](/images/Information_icon.svg?eefcf)](/w/File:Information_icon.svg)

**本条目所述内容仅适用于[Java版](/w/Java%E7%89%88 "Java版")。**

 **![](/images/Comment_information.svg?eab4a) Wiki上有与该主题相关的教程！**

见[教程:物品堆叠组件](/w/Tutorial:%E7%89%A9%E5%93%81%E5%A0%86%E5%8F%A0%E7%BB%84%E4%BB%B6 "Tutorial:物品堆叠组件")。

 **![](/images/Comment_information.svg?eab4a) Wiki上有与该主题相关的教程！**

见[教程:物品堆叠组件](/w/Tutorial:%E7%89%A9%E5%93%81%E5%A0%86%E5%8F%A0%E7%BB%84%E4%BB%B6 "Tutorial:物品堆叠组件")。

**数据组件（Data Component）**[\[1\]](#cite_note-1)，或简称为**组件（Component）**，是用于定义和存储各种数据属性的结构化数据。

由于[物品堆叠](/w/%E7%89%A9%E5%93%81%E5%A0%86%E5%8F%A0 "物品堆叠")全面使用数据组件格式，故其也被称为**物品堆叠组件（Item Stack Component）**或**物品组件（Item Component）**。[\[2\]](#cite_note-24w09a-2)

## 目录

-   [1 行为](#行为)
    -   [1.1 加载行为](#加载行为)
    -   [1.2 物品堆叠](#物品堆叠)
    -   [1.3 方块实体](#方块实体)
    -   [1.4 实体](#实体)
-   [2 数据组件类型](#数据组件类型)
-   [3 命令格式](#命令格式)
-   [4 数据格式](#数据格式)
    -   [4.1 attack\_range](#attack_range)
    -   [4.2 attribute\_modifiers](#attribute_modifiers)
    -   [4.3 banner\_patterns](#banner_patterns)
    -   [4.4 base\_color](#base_color)
    -   [4.5 bees](#bees)
    -   [4.6 block\_entity\_data](#block_entity_data)
    -   [4.7 block\_state](#block_state)
    -   [4.8 blocks\_attacks](#blocks_attacks)
    -   [4.9 break\_sound](#break_sound)
    -   [4.10 bucket\_entity\_data](#bucket_entity_data)
    -   [4.11 bundle\_contents](#bundle_contents)
    -   [4.12 can\_break和can\_place\_on](#can_break和can_place_on)
    -   [4.13 charged\_projectiles](#charged_projectiles)
    -   [4.14 consumable](#consumable)
    -   [4.15 container](#container)
    -   [4.16 container\_loot](#container_loot)
    -   [4.17 custom\_data](#custom_data)
    -   [4.18 custom\_model\_data](#custom_model_data)
    -   [4.19 custom\_name](#custom_name)
    -   [4.20 damage](#damage)
    -   [4.21 damage\_resistant](#damage_resistant)
    -   [4.22 damage\_type](#damage_type)
    -   [4.23 death\_protection](#death_protection)
    -   [4.24 debug\_stick\_state](#debug_stick_state)
    -   [4.25 dye](#dye)
    -   [4.26 dyed\_color](#dyed_color)
    -   [4.27 enchantable](#enchantable)
    -   [4.28 enchantment\_glint\_override](#enchantment_glint_override)
    -   [4.29 enchantments和stored\_enchantments](#enchantments和stored_enchantments)
    -   [4.30 entity\_data](#entity_data)
    -   [4.31 equippable](#equippable)
    -   [4.32 firework\_explosion](#firework_explosion)
    -   [4.33 fireworks](#fireworks)
    -   [4.34 food](#food)
    -   [4.35 glider](#glider)
    -   [4.36 instrument](#instrument)
    -   [4.37 intangible\_projectile](#intangible_projectile)
    -   [4.38 item\_model](#item_model)
    -   [4.39 item\_name](#item_name)
    -   [4.40 jukebox\_playable](#jukebox_playable)
    -   [4.41 kinetic\_weapon](#kinetic_weapon)
    -   [4.42 lock](#lock)
    -   [4.43 lodestone\_tracker](#lodestone_tracker)
    -   [4.44 lore](#lore)
    -   [4.45 map\_color](#map_color)
    -   [4.46 map\_decorations](#map_decorations)
    -   [4.47 map\_id](#map_id)
    -   [4.48 max\_damage](#max_damage)
    -   [4.49 max\_stack\_size](#max_stack_size)
    -   [4.50 minimum\_attack\_charge](#minimum_attack_charge)
    -   [4.51 note\_block\_sound](#note_block_sound)
    -   [4.52 ominous\_bottle\_amplifier](#ominous_bottle_amplifier)
    -   [4.53 piercing\_weapon](#piercing_weapon)
    -   [4.54 pot\_decorations](#pot_decorations)
    -   [4.55 potion\_contents](#potion_contents)
    -   [4.56 potion\_duration\_scale](#potion_duration_scale)
    -   [4.57 profile](#profile)
    -   [4.58 provides\_banner\_patterns](#provides_banner_patterns)
    -   [4.59 provides\_trim\_material](#provides_trim_material)
    -   [4.60 rarity](#rarity)
    -   [4.61 recipes](#recipes)
    -   [4.62 repair\_cost](#repair_cost)
    -   [4.63 repairable](#repairable)
    -   [4.64 suspicious\_stew\_effects](#suspicious_stew_effects)
    -   [4.65 swing\_animation](#swing_animation)
    -   [4.66 tool](#tool)
    -   [4.67 tooltip\_display](#tooltip_display)
    -   [4.68 tooltip\_style](#tooltip_style)
    -   [4.69 trim](#trim)
    -   [4.70 unbreakable](#unbreakable)
    -   [4.71 use\_cooldown](#use_cooldown)
    -   [4.72 use\_effects](#use_effects)
    -   [4.73 use\_remainder](#use_remainder)
    -   [4.74 weapon](#weapon)
    -   [4.75 writable\_book\_content](#writable_book_content)
    -   [4.76 written\_book\_content](#written_book_content)
    -   [4.77 实体变种组件](#实体变种组件)
        -   [4.77.1 axolotl/variant](#axolotl/variant)
        -   [4.77.2 cat/collar](#cat/collar)
        -   [4.77.3 cat/sound\_variant](#cat/sound_variant)
        -   [4.77.4 cat/variant](#cat/variant)
        -   [4.77.5 chicken/sound\_variant](#chicken/sound_variant)
        -   [4.77.6 chicken/variant](#chicken/variant)
        -   [4.77.7 cow/sound\_variant](#cow/sound_variant)
        -   [4.77.8 cow/variant](#cow/variant)
        -   [4.77.9 fox/variant](#fox/variant)
        -   [4.77.10 frog/variant](#frog/variant)
        -   [4.77.11 horse/variant](#horse/variant)
        -   [4.77.12 llama/variant](#llama/variant)
        -   [4.77.13 mooshroom/variant](#mooshroom/variant)
        -   [4.77.14 painting/variant](#painting/variant)
        -   [4.77.15 parrot/variant](#parrot/variant)
        -   [4.77.16 pig/sound\_variant](#pig/sound_variant)
        -   [4.77.17 pig/variant](#pig/variant)
        -   [4.77.18 rabbit/variant](#rabbit/variant)
        -   [4.77.19 salmon/size](#salmon/size)
        -   [4.77.20 sheep/color](#sheep/color)
        -   [4.77.21 shulker/color](#shulker/color)
        -   [4.77.22 tropical\_fish/base\_color](#tropical_fish/base_color)
        -   [4.77.23 tropical\_fish/pattern](#tropical_fish/pattern)
        -   [4.77.24 tropical\_fish/pattern\_color](#tropical_fish/pattern_color)
        -   [4.77.25 villager/variant](#villager/variant)
        -   [4.77.26 wolf/collar](#wolf/collar)
        -   [4.77.27 wolf/sound\_variant](#wolf/sound_variant)
        -   [4.77.28 wolf/variant](#wolf/variant)
        -   [4.77.29 zombie\_nautilus/variant](#zombie_nautilus/variant)
    -   [4.78 临时组件](#临时组件)
        -   [4.78.1 additional\_trade\_cost](#additional_trade_cost)
        -   [4.78.2 creative\_slot\_lock](#creative_slot_lock)
        -   [4.78.3 map\_post\_processing](#map_post_processing)
-   [5 历史](#历史)
    -   [5.1 已移除的组件](#已移除的组件)
        -   [5.1.1 fire\_resistant](#fire_resistant)
        -   [5.1.2 hide\_additional\_tooltip](#hide_additional_tooltip)
        -   [5.1.3 hide\_tooltip](#hide_tooltip)
-   [6 参考](#参考)
-   [7 导航](#导航)

## 行为

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=1&veaction=edit "编辑章节：行为") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=1 "编辑章节的源代码： 行为")\]

数据组件是结构化的数据，也即每一个组件都有自己独特的编码方式。如果组件的数据格式不正确，则游戏会立即解析失败，对应的命令和文件等全部无效。

由于数据组件的编码解码行为，使得其与通常的[NBT标签](/w/NBT%E6%A0%87%E7%AD%BE "NBT标签")数据不同。通常的NBT标签仅会在游戏尝试序列化为程序对象时才判断其是否符合编码格式，而组件自游戏加载之初就进行了判断。这使得数据组件格式拥有更快的加载性能，可以更早地发现命令和文件中的潜在错误。[\[2\]](#cite_note-24w09a-2)

除了编码方式外，每一个组件都有是否持久化和是否同步两个基本性质。不可持久化的组件通常仅用于网络传输，随游戏计算完毕或内存卸载而移除，不会保存到[存档](/w/%E5%AD%98%E6%A1%A3 "存档")里，强行加载和保存也会导致游戏解析失败。若无特殊说明，下文的组件均指持久化组件。

### 加载行为

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=2&veaction=edit "编辑章节：加载行为") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=2 "编辑章节的源代码： 加载行为")\]

物品堆叠、方块实体和实体可以拥有数据组件。

目前方块实体和实体依然使用了非结构化的NBT标签存储。为了对方块实体或实体套用或获取数据组件，游戏会将对应的组件和对应的NBT标签绑定。这步绑定操作所用的组件在游戏内部被称为隐式组件（Implicit Component）。游戏通常在使用物品放置方块或实体时套用组件，而使用谓词检测或破坏方块时获取组件。

例如：当使用命名过的箱子放置箱子时，箱子物品的`custom_name`组件会套用到箱子的方块实体上。而方块实体会使用方块实体数据![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)CustomName来保存它。而当从箱子获取`custom_name`组件时，游戏会将方块实体数据![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)CustomName视为`custom_name`组件。

### 物品堆叠

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=3&veaction=edit "编辑章节：物品堆叠") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=3 "编辑章节的源代码： 物品堆叠")\]

[![](/images/thumb/The_Component_Count_In_An_Item_Stack_Simplified.png/300px-The_Component_Count_In_An_Item_Stack_Simplified.png?d112b)](/w/File:The_Component_Count_In_An_Item_Stack_Simplified.png)

在物品提示框中会显示当前物品堆叠的组件数

参见：[物品格式](/w/%E7%89%A9%E5%93%81%E6%A0%BC%E5%BC%8F "物品格式")

物品堆叠全面使用数据组件格式。虽然游戏为每个物品定义了依据物品类型的默认组件，但默认组件只在内存中计算，不会保存到存档中。而存档会存储物品的组件修订（Data Component Patch）数据，组件修订中指定的组件会覆盖默认组件的值，且带`!`前缀的组件会移除该物品的默认组件。

绝大多数数据组件都对物品自身有实际作用，决定了物品的诸多性质，例如是否可堆叠、可损坏等，影响了大量的游戏行为。当以组件为单位修改物品时，游戏不允许物品同时具有`damage`组件和值大于1的`max_stack_size`的组件补丁，即物品不可以既可堆叠又可损坏。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597) 物品堆叠数据
    -   ![字符串](/images/Data_node_string.svg?42545)\*  
        \*id：（[命名空间ID](/w/%E5%91%BD%E5%90%8D%E7%A9%BA%E9%97%B4ID "命名空间ID")）表示某种类的物品堆叠。若未指定，游戏会在加载区块或者生成物品时将其变更为空气。
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components：当前物品的组件修订，将修改物品的数据组件信息。
        -   ![任意类型](/images/Data_node_any.svg?d406c)<*数据组件ID*\>：一项组件和其对应的数据，代表物品拥有此组件。设置组件数据时可以不写命名空间，但游戏在导出时会自行加上`minecraft:`前缀。
        -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)!<*数据组件ID*\>：存在时，使一个数据组件失效。此复合标签的内容不影响行为。设置组件数据时可以不写命名空间，但游戏在导出时会自行加上`minecraft:`前缀。
    -   ![整型](/images/Data_node_int.svg?8d24f)count：（0<值≤物品最大堆叠数量）[物品](/w/%E7%89%A9%E5%93%81 "物品")的堆叠数。不存在或无效时则默认为1。

### 方块实体

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=4&veaction=edit "编辑章节：方块实体") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=4 "编辑章节的源代码： 方块实体")\]

参见：[方块实体数据格式](/w/%E6%96%B9%E5%9D%97%E5%AE%9E%E4%BD%93%E6%95%B0%E6%8D%AE%E6%A0%BC%E5%BC%8F "方块实体数据格式")

方块实体部分采用数据组件格式。一部分组件以组件格式原样存储，另一部分则以隐式组件存储于NBT标签中。目前方块实体不支持删除组件。

当使用方块物品放置方块时，`block_state`和`block_entity_data`组件永远不会保存到方块实体中。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597) 方块实体数据
    -   ![整型](/images/Data_node_int.svg?8d24f)\*  
        \*x：当前方块实体的X坐标。
    -   ![整型](/images/Data_node_int.svg?8d24f)\*  
        \*y：当前方块实体的Y坐标。
    -   ![整型](/images/Data_node_int.svg?8d24f)\*  
        \*z：当前方块实体的Z坐标。
    -   ![字符串](/images/Data_node_string.svg?42545)\*  
        \*id：（命名空间ID）方块实体的类型。
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components：方块实体的数据组件信息。当使用此方块实体对应的物品放置此方块实体时，物品额外持有的且不会被继承序列化处理的数据组件会被复制存储入此标签内。
        -   ![任意类型](/images/Data_node_any.svg?d406c)<*数据组件ID*\>：一项数据组件和其对应的数据。

游戏内使用的以隐式组件存储的方块实体组件如下：

命名空间ID

作用方块

方块实体数据

与`block_entity_data`组件间的行为

[banner\_patterns](#banner_patterns)

[![](/images/BlockSprite_all-banners.png?b9a66)](/w/%E6%97%97%E5%B8%9C "旗帜")[旗帜](/w/%E6%97%97%E5%B8%9C "旗帜")

![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)patterns

高优先级，作用的方块实体数据永远无法被`block_entity_data`设置

[bees](#bees)

[![](/images/BlockSprite_bee-nest.png?70047)](/w/%E8%9C%82%E5%B7%A2%EF%BC%88%E6%96%B9%E5%9D%97%EF%BC%89 "蜂巢（方块）")[蜂巢](/w/%E8%9C%82%E5%B7%A2%EF%BC%88%E6%96%B9%E5%9D%97%EF%BC%89 "蜂巢（方块）")  
[![](/images/BlockSprite_beehive.png?0a98d)](/w/%E8%9C%82%E7%AE%B1 "蜂箱")[蜂箱](/w/%E8%9C%82%E7%AE%B1 "蜂箱")

![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)bees

[container](#container)

**容器方块**

[![](/images/BlockSprite_chest.png?05052)](/w/%E7%AE%B1%E5%AD%90 "箱子")[箱子](/w/%E7%AE%B1%E5%AD%90 "箱子")  
[![](/images/BlockSprite_trapped-chest.png?45f4c)](/w/%E9%99%B7%E9%98%B1%E7%AE%B1 "陷阱箱")[陷阱箱](/w/%E9%99%B7%E9%98%B1%E7%AE%B1 "陷阱箱")  
[![](/images/BlockSprite_copper-chest-front.png?2d44c)](/w/%E9%93%9C%E7%AE%B1%E5%AD%90 "铜箱子")[铜箱子](/w/%E9%93%9C%E7%AE%B1%E5%AD%90 "铜箱子")  
[![](/images/BlockSprite_barrel.png?4e239)](/w/%E6%9C%A8%E6%A1%B6 "木桶")[木桶](/w/%E6%9C%A8%E6%A1%B6 "木桶")  
[![](/images/BlockSprite_all-shulker-boxes.png?a78b3)](/w/%E6%BD%9C%E5%BD%B1%E7%9B%92 "潜影盒")[潜影盒](/w/%E6%BD%9C%E5%BD%B1%E7%9B%92 "潜影盒")  
[![](/images/BlockSprite_furnace.png?c6241)](/w/%E7%86%94%E7%82%89 "熔炉")[熔炉](/w/%E7%86%94%E7%82%89 "熔炉")  
[![](/images/BlockSprite_smoker.png?f1a73)](/w/%E7%83%9F%E7%86%8F%E7%82%89 "烟熏炉")[烟熏炉](/w/%E7%83%9F%E7%86%8F%E7%82%89 "烟熏炉")  
[![](/images/BlockSprite_blast-furnace.png?7117e)](/w/%E9%AB%98%E7%82%89 "高炉")[高炉](/w/%E9%AB%98%E7%82%89 "高炉")  
[![](/images/BlockSprite_dispenser.png?555fa)](/w/%E5%8F%91%E5%B0%84%E5%99%A8 "发射器")[发射器](/w/%E5%8F%91%E5%B0%84%E5%99%A8 "发射器")  
[![](/images/BlockSprite_dropper.png?e13bc)](/w/%E6%8A%95%E6%8E%B7%E5%99%A8 "投掷器")[投掷器](/w/%E6%8A%95%E6%8E%B7%E5%99%A8 "投掷器")  
[![](/images/BlockSprite_crafter.png?a29bd)](/w/%E5%90%88%E6%88%90%E5%99%A8 "合成器")[合成器](/w/%E5%90%88%E6%88%90%E5%99%A8 "合成器")  
[![](/images/BlockSprite_brewing-stand.png?6a0a0)](/w/%E9%85%BF%E9%80%A0%E5%8F%B0 "酿造台")[酿造台](/w/%E9%85%BF%E9%80%A0%E5%8F%B0 "酿造台")  
[![](/images/BlockSprite_hopper.png?0a53f)](/w/%E6%BC%8F%E6%96%97 "漏斗")[漏斗](/w/%E6%BC%8F%E6%96%97 "漏斗")

[![](/images/BlockSprite_campfire.png?1c233)](/w/%E8%90%A5%E7%81%AB "营火")[营火](/w/%E8%90%A5%E7%81%AB "营火")  
[![](/images/BlockSprite_soul-campfire.png?8ba70)](/w/%E7%81%B5%E9%AD%82%E8%90%A5%E7%81%AB "灵魂营火")[灵魂营火](/w/%E7%81%B5%E9%AD%82%E8%90%A5%E7%81%AB "灵魂营火")  
[![](/images/BlockSprite_chiseled-bookshelf.png?83d51)](/w/%E9%9B%95%E7%BA%B9%E4%B9%A6%E6%9E%B6 "雕纹书架")[雕纹书架](/w/%E9%9B%95%E7%BA%B9%E4%B9%A6%E6%9E%B6 "雕纹书架")  
[![](/images/BlockSprite_oak-shelf-front.png?2e1fa)](/w/%E5%B1%95%E7%A4%BA%E6%9E%B6 "展示架")[展示架](/w/%E5%B1%95%E7%A4%BA%E6%9E%B6 "展示架")

![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)Items

[![](/images/BlockSprite_decorated-pot.png?73142)](/w/%E9%A5%B0%E7%BA%B9%E9%99%B6%E7%BD%90 "饰纹陶罐")[饰纹陶罐](/w/%E9%A5%B0%E7%BA%B9%E9%99%B6%E7%BD%90 "饰纹陶罐")

![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)Item  
（注：对应组件列表的第一个物品）

[container\_loot](#container_loot)

**战利品容器方块**

[![](/images/BlockSprite_chest.png?05052)](/w/%E7%AE%B1%E5%AD%90 "箱子")[箱子](/w/%E7%AE%B1%E5%AD%90 "箱子")  
[![](/images/BlockSprite_trapped-chest.png?45f4c)](/w/%E9%99%B7%E9%98%B1%E7%AE%B1 "陷阱箱")[陷阱箱](/w/%E9%99%B7%E9%98%B1%E7%AE%B1 "陷阱箱")  
[![](/images/BlockSprite_copper-chest-front.png?2d44c)](/w/%E9%93%9C%E7%AE%B1%E5%AD%90 "铜箱子")[铜箱子](/w/%E9%93%9C%E7%AE%B1%E5%AD%90 "铜箱子")  
[![](/images/BlockSprite_barrel.png?4e239)](/w/%E6%9C%A8%E6%A1%B6 "木桶")[木桶](/w/%E6%9C%A8%E6%A1%B6 "木桶")  
[![](/images/BlockSprite_all-shulker-boxes.png?a78b3)](/w/%E6%BD%9C%E5%BD%B1%E7%9B%92 "潜影盒")[潜影盒](/w/%E6%BD%9C%E5%BD%B1%E7%9B%92 "潜影盒")  
[![](/images/BlockSprite_dispenser.png?555fa)](/w/%E5%8F%91%E5%B0%84%E5%99%A8 "发射器")[发射器](/w/%E5%8F%91%E5%B0%84%E5%99%A8 "发射器")  
[![](/images/BlockSprite_dropper.png?e13bc)](/w/%E6%8A%95%E6%8E%B7%E5%99%A8 "投掷器")[投掷器](/w/%E6%8A%95%E6%8E%B7%E5%99%A8 "投掷器")  
[![](/images/BlockSprite_crafter.png?a29bd)](/w/%E5%90%88%E6%88%90%E5%99%A8 "合成器")[合成器](/w/%E5%90%88%E6%88%90%E5%99%A8 "合成器")

![长整型](/images/Data_node_long.svg?dde3f)LootTableSeed  
![字符串](/images/Data_node_string.svg?42545)LootTable

高优先级，作用的方块实体数据在此组件不存在时可以被`block_entity_data`设置

[custom\_name](#custom_name)

**容器方块**

[![](/images/BlockSprite_chest.png?05052)](/w/%E7%AE%B1%E5%AD%90 "箱子")[箱子](/w/%E7%AE%B1%E5%AD%90 "箱子")  
[![](/images/BlockSprite_trapped-chest.png?45f4c)](/w/%E9%99%B7%E9%98%B1%E7%AE%B1 "陷阱箱")[陷阱箱](/w/%E9%99%B7%E9%98%B1%E7%AE%B1 "陷阱箱")  
[![](/images/BlockSprite_copper-chest-front.png?2d44c)](/w/%E9%93%9C%E7%AE%B1%E5%AD%90 "铜箱子")[铜箱子](/w/%E9%93%9C%E7%AE%B1%E5%AD%90 "铜箱子")  
[![](/images/BlockSprite_barrel.png?4e239)](/w/%E6%9C%A8%E6%A1%B6 "木桶")[木桶](/w/%E6%9C%A8%E6%A1%B6 "木桶")  
[![](/images/BlockSprite_all-shulker-boxes.png?a78b3)](/w/%E6%BD%9C%E5%BD%B1%E7%9B%92 "潜影盒")[潜影盒](/w/%E6%BD%9C%E5%BD%B1%E7%9B%92 "潜影盒")  
[![](/images/BlockSprite_furnace.png?c6241)](/w/%E7%86%94%E7%82%89 "熔炉")[熔炉](/w/%E7%86%94%E7%82%89 "熔炉")  
[![](/images/BlockSprite_smoker.png?f1a73)](/w/%E7%83%9F%E7%86%8F%E7%82%89 "烟熏炉")[烟熏炉](/w/%E7%83%9F%E7%86%8F%E7%82%89 "烟熏炉")  
[![](/images/BlockSprite_blast-furnace.png?7117e)](/w/%E9%AB%98%E7%82%89 "高炉")[高炉](/w/%E9%AB%98%E7%82%89 "高炉")  
[![](/images/BlockSprite_dispenser.png?555fa)](/w/%E5%8F%91%E5%B0%84%E5%99%A8 "发射器")[发射器](/w/%E5%8F%91%E5%B0%84%E5%99%A8 "发射器")  
[![](/images/BlockSprite_dropper.png?e13bc)](/w/%E6%8A%95%E6%8E%B7%E5%99%A8 "投掷器")[投掷器](/w/%E6%8A%95%E6%8E%B7%E5%99%A8 "投掷器")  
[![](/images/BlockSprite_crafter.png?a29bd)](/w/%E5%90%88%E6%88%90%E5%99%A8 "合成器")[合成器](/w/%E5%90%88%E6%88%90%E5%99%A8 "合成器")  
[![](/images/BlockSprite_brewing-stand.png?6a0a0)](/w/%E9%85%BF%E9%80%A0%E5%8F%B0 "酿造台")[酿造台](/w/%E9%85%BF%E9%80%A0%E5%8F%B0 "酿造台")  
[![](/images/BlockSprite_hopper.png?0a53f)](/w/%E6%BC%8F%E6%96%97 "漏斗")[漏斗](/w/%E6%BC%8F%E6%96%97 "漏斗")

[![](/images/BlockSprite_all-banners.png?b9a66)](/w/%E6%97%97%E5%B8%9C "旗帜")[旗帜](/w/%E6%97%97%E5%B8%9C "旗帜")  
[![](/images/BlockSprite_enchanting-table.png?df988)](/w/%E9%99%84%E9%AD%94%E5%8F%B0 "附魔台")[附魔台](/w/%E9%99%84%E9%AD%94%E5%8F%B0 "附魔台")  
[![](/images/BlockSprite_beacon.png?869cc)](/w/%E4%BF%A1%E6%A0%87 "信标")[信标](/w/%E4%BF%A1%E6%A0%87 "信标")  
[![](/images/BlockSprite_command-block.png?df114)](/w/%E5%91%BD%E4%BB%A4%E6%96%B9%E5%9D%97 "命令方块")[命令方块](/w/%E5%91%BD%E4%BB%A4%E6%96%B9%E5%9D%97 "命令方块")  
[![](/images/BlockSprite_chain-command-block.png?82b71)](/w/%E5%91%BD%E4%BB%A4%E6%96%B9%E5%9D%97 "命令方块")[连锁型命令方块](/w/%E5%91%BD%E4%BB%A4%E6%96%B9%E5%9D%97 "命令方块")  
[![](/images/BlockSprite_repeating-command-block.png?44fd4)](/w/%E5%91%BD%E4%BB%A4%E6%96%B9%E5%9D%97 "命令方块")[循环型命令方块](/w/%E5%91%BD%E4%BB%A4%E6%96%B9%E5%9D%97 "命令方块")

![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)CustomName

高优先级，作用的方块实体数据永远无法被`block_entity_data`设置

[![](/images/BlockSprite_all-heads.png?dc1e4)](/w/%E7%94%9F%E7%89%A9%E5%A4%B4%E9%A2%85 "生物头颅")[生物头颅](/w/%E7%94%9F%E7%89%A9%E5%A4%B4%E9%A2%85 "生物头颅")

![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)custom\_name

[lock](#lock)

**容器方块**

[![](/images/BlockSprite_chest.png?05052)](/w/%E7%AE%B1%E5%AD%90 "箱子")[箱子](/w/%E7%AE%B1%E5%AD%90 "箱子")  
[![](/images/BlockSprite_trapped-chest.png?45f4c)](/w/%E9%99%B7%E9%98%B1%E7%AE%B1 "陷阱箱")[陷阱箱](/w/%E9%99%B7%E9%98%B1%E7%AE%B1 "陷阱箱")  
[![](/images/BlockSprite_copper-chest-front.png?2d44c)](/w/%E9%93%9C%E7%AE%B1%E5%AD%90 "铜箱子")[铜箱子](/w/%E9%93%9C%E7%AE%B1%E5%AD%90 "铜箱子")  
[![](/images/BlockSprite_barrel.png?4e239)](/w/%E6%9C%A8%E6%A1%B6 "木桶")[木桶](/w/%E6%9C%A8%E6%A1%B6 "木桶")  
[![](/images/BlockSprite_all-shulker-boxes.png?a78b3)](/w/%E6%BD%9C%E5%BD%B1%E7%9B%92 "潜影盒")[潜影盒](/w/%E6%BD%9C%E5%BD%B1%E7%9B%92 "潜影盒")  
[![](/images/BlockSprite_furnace.png?c6241)](/w/%E7%86%94%E7%82%89 "熔炉")[熔炉](/w/%E7%86%94%E7%82%89 "熔炉")  
[![](/images/BlockSprite_smoker.png?f1a73)](/w/%E7%83%9F%E7%86%8F%E7%82%89 "烟熏炉")[烟熏炉](/w/%E7%83%9F%E7%86%8F%E7%82%89 "烟熏炉")  
[![](/images/BlockSprite_blast-furnace.png?7117e)](/w/%E9%AB%98%E7%82%89 "高炉")[高炉](/w/%E9%AB%98%E7%82%89 "高炉")  
[![](/images/BlockSprite_dispenser.png?555fa)](/w/%E5%8F%91%E5%B0%84%E5%99%A8 "发射器")[发射器](/w/%E5%8F%91%E5%B0%84%E5%99%A8 "发射器")  
[![](/images/BlockSprite_dropper.png?e13bc)](/w/%E6%8A%95%E6%8E%B7%E5%99%A8 "投掷器")[投掷器](/w/%E6%8A%95%E6%8E%B7%E5%99%A8 "投掷器")  
[![](/images/BlockSprite_crafter.png?a29bd)](/w/%E5%90%88%E6%88%90%E5%99%A8 "合成器")[合成器](/w/%E5%90%88%E6%88%90%E5%99%A8 "合成器")  
[![](/images/BlockSprite_brewing-stand.png?6a0a0)](/w/%E9%85%BF%E9%80%A0%E5%8F%B0 "酿造台")[酿造台](/w/%E9%85%BF%E9%80%A0%E5%8F%B0 "酿造台")  
[![](/images/BlockSprite_hopper.png?0a53f)](/w/%E6%BC%8F%E6%96%97 "漏斗")[漏斗](/w/%E6%BC%8F%E6%96%97 "漏斗")

[![](/images/BlockSprite_beacon.png?869cc)](/w/%E4%BF%A1%E6%A0%87 "信标")[信标](/w/%E4%BF%A1%E6%A0%87 "信标")

![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)lock

[note\_block\_sound](#note_block_sound)

[![](/images/BlockSprite_all-heads.png?dc1e4)](/w/%E7%94%9F%E7%89%A9%E5%A4%B4%E9%A2%85 "生物头颅")[生物头颅](/w/%E7%94%9F%E7%89%A9%E5%A4%B4%E9%A2%85 "生物头颅")

![字符串](/images/Data_node_string.svg?42545)note\_block\_sound

[pot\_decorations](#pot_decorations)

[![](/images/BlockSprite_decorated-pot.png?73142)](/w/%E9%A5%B0%E7%BA%B9%E9%99%B6%E7%BD%90 "饰纹陶罐")[饰纹陶罐](/w/%E9%A5%B0%E7%BA%B9%E9%99%B6%E7%BD%90 "饰纹陶罐")

![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)sherds

[profile](#profile)

[![](/images/BlockSprite_all-heads.png?dc1e4)](/w/%E7%94%9F%E7%89%A9%E5%A4%B4%E9%A2%85 "生物头颅")[生物头颅](/w/%E7%94%9F%E7%89%A9%E5%A4%B4%E9%A2%85 "生物头颅")

![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)profile

### 实体

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=5&veaction=edit "编辑章节：实体") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=5 "编辑章节的源代码： 实体")\]

实体的数据组件全部以隐式组件的形式存储于非组件结构的NBT数据中。

若物品同时具有`bucket_entity_data`、`entity_data`组件和其他实体组件，则套用优先级依次为`bucket_entity_data`、`entity_data`、其他组件。

游戏内使用的以隐式组件存储的实体组件如下：

命名空间ID

作用实体

实体数据

[axolotl/variant](#axolotl/variant)

[![](/images/EntitySprite_axolotl.png?2597f)](/w/%E7%BE%8E%E8%A5%BF%E8%9E%88 "美西螈")[美西螈](/w/%E7%BE%8E%E8%A5%BF%E8%9E%88 "美西螈")

![整型](/images/Data_node_int.svg?8d24f)Variant

[cat/collar](#cat/collar)

[![](/images/EntitySprite_cat.png?4d91c)](/w/%E7%8C%AB "猫")[猫](/w/%E7%8C%AB "猫")

![字节型](/images/Data_node_byte.svg?eb0e0)CollarColor

[cat/variant](#cat/variant)

![字符串](/images/Data_node_string.svg?42545)variant

[cat/sound\_variant](#cat/sound_variant)\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]

![字符串](/images/Data_node_string.svg?42545)sound\_variant

[chicken/variant](#chicken/variant)

[![](/images/EntitySprite_chicken.png?ed28c)](/w/%E9%B8%A1 "鸡")[鸡](/w/%E9%B8%A1 "鸡")

![字符串](/images/Data_node_string.svg?42545)variant

[chicken/sound\_variant](#chicken/sound_variant)\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]

![字符串](/images/Data_node_string.svg?42545)sound\_variant

[cow/variant](#cow/variant)

[![](/images/EntitySprite_cow.png?54276)](/w/%E7%89%9B "牛")[牛](/w/%E7%89%9B "牛")

![字符串](/images/Data_node_string.svg?42545)variant

[cow/sound\_variant](#cow/sound_variant)\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]

![字符串](/images/Data_node_string.svg?42545)sound\_variant

[custom\_data](#custom_data)

**所有实体**

![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)data

[custom\_name](#custom_name)

![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)CustomName

[fox/variant](#fox/variant)

[![](/images/EntitySprite_fox.png?c5f58)](/w/%E7%8B%90%E7%8B%B8 "狐狸")[狐狸](/w/%E7%8B%90%E7%8B%B8 "狐狸")

![字符串](/images/Data_node_string.svg?42545)Type

[frog/variant](#frog/variant)

[![](/images/EntitySprite_frog.png?39d41)](/w/%E9%9D%92%E8%9B%99 "青蛙")[青蛙](/w/%E9%9D%92%E8%9B%99 "青蛙")

![字符串](/images/Data_node_string.svg?42545)variant

[horse/variant](#horse/variant)

[![](/images/EntitySprite_creamy-horse.png?80b78)](/w/%E9%A9%AC "马")[马](/w/%E9%A9%AC "马")

![整型](/images/Data_node_int.svg?8d24f)Variant的后8位

[llama/variant](#llama/variant)

[![](/images/EntitySprite_creamy-llama.png?b264d)](/w/%E7%BE%8A%E9%A9%BC "羊驼")[羊驼](/w/%E7%BE%8A%E9%A9%BC "羊驼")  
[![](/images/EntitySprite_creamy-trader-llama.png?80ba1)](/w/%E8%A1%8C%E5%95%86%E7%BE%8A%E9%A9%BC "行商羊驼")[行商羊驼](/w/%E8%A1%8C%E5%95%86%E7%BE%8A%E9%A9%BC "行商羊驼")

![整型](/images/Data_node_int.svg?8d24f)Variant

[mooshroom/variant](#mooshroom/variant)

[![](/images/EntitySprite_mooshroom.png?837a6)](/w/%E5%93%9E%E8%8F%87 "哞菇")[哞菇](/w/%E5%93%9E%E8%8F%87 "哞菇")

![字符串](/images/Data_node_string.svg?42545)Type

[painting/variant](#painting/variant)

[![](/images/EntitySprite_kebab.png?15a6b)](/w/%E7%94%BB "画")[画](/w/%E7%94%BB "画")

![字符串](/images/Data_node_string.svg?42545)variant

[parrot/variant](#parrot/variant)

[![](/images/EntitySprite_parrot.png?ed0cc)](/w/%E9%B9%A6%E9%B9%89 "鹦鹉")[鹦鹉](/w/%E9%B9%A6%E9%B9%89 "鹦鹉")

![整型](/images/Data_node_int.svg?8d24f)Variant

[profile](#profile)

[![](/images/EntitySprite_alex.png?dc657)](/w/%E7%8E%A9%E5%AE%B6%E6%A8%A1%E5%9E%8B "玩家模型")[玩家模型](/w/%E7%8E%A9%E5%AE%B6%E6%A8%A1%E5%9E%8B "玩家模型")

![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)profile

[pig/variant](#pig/variant)

[![](/images/EntitySprite_pig.png?c2459)](/w/%E7%8C%AA "猪")[猪](/w/%E7%8C%AA "猪")

![字符串](/images/Data_node_string.svg?42545)variant

[pig/sound\_variant](#pig/sound_variant)\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]

![字符串](/images/Data_node_string.svg?42545)sound\_variant

[potion\_contents](#potion_contents)

[![](/images/EntitySprite_area-effect-cloud.png?9566c)](/w/%E5%8C%BA%E5%9F%9F%E6%95%88%E6%9E%9C%E4%BA%91 "区域效果云")[区域效果云](/w/%E5%8C%BA%E5%9F%9F%E6%95%88%E6%9E%9C%E4%BA%91 "区域效果云")

![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)potion\_contents

[potion\_duration\_scale](#potion_duration_scale)

![单精度浮点数](/images/Data_node_float.svg?ae55e)potion\_duration\_scale

[rabbit/variant](#rabbit/variant)

[![](/images/EntitySprite_brown-rabbit.png?af02d)](/w/%E5%85%94%E5%AD%90 "兔子")[兔子](/w/%E5%85%94%E5%AD%90 "兔子")

![整型](/images/Data_node_int.svg?8d24f)RabbitType

[salmon/size](#salmon/size)

[![](/images/EntitySprite_salmon.png?913ea)](/w/%E9%B2%91%E9%B1%BC "鲑鱼")[鲑鱼](/w/%E9%B2%91%E9%B1%BC "鲑鱼")

![字符串](/images/Data_node_string.svg?42545)type

[sheep/color](#sheep/color)

[![](/images/EntitySprite_white-sheep.png?07a4a)](/w/%E7%BB%B5%E7%BE%8A "绵羊")[绵羊](/w/%E7%BB%B5%E7%BE%8A "绵羊")

![字节型](/images/Data_node_byte.svg?eb0e0)Color

[shulker/color](#shulker/color)

[![](/images/EntitySprite_shulker.png?a8d4f)](/w/%E6%BD%9C%E5%BD%B1%E8%B4%9D "潜影贝")[潜影贝](/w/%E6%BD%9C%E5%BD%B1%E8%B4%9D "潜影贝")

![字节型](/images/Data_node_byte.svg?eb0e0)Color

[tropical\_fish/base\_color](#tropical_fish/base_color)

[![](/images/EntitySprite_tropical-fish.png?20777)](/w/%E7%83%AD%E5%B8%A6%E9%B1%BC "热带鱼")[热带鱼](/w/%E7%83%AD%E5%B8%A6%E9%B1%BC "热带鱼")

![整型](/images/Data_node_int.svg?8d24f)Variant的低2位

[tropical\_fish/pattern](#tropical_fish/pattern)

![整型](/images/Data_node_int.svg?8d24f)Variant的中4位

[tropical\_fish/pattern\_color](#tropical_fish/pattern_color)

![整型](/images/Data_node_int.svg?8d24f)Variant的高4位

[villager/variant](#villager/variant)

[![](/images/EntitySprite_leatherworker.png?43c3d)](/w/%E6%9D%91%E6%B0%91 "村民")[村民](/w/%E6%9D%91%E6%B0%91 "村民")

![字符串](/images/Data_node_string.svg?42545)VillagerData.type

[![](/images/EntitySprite_zombie-villager.png?a14eb)](/w/%E5%83%B5%E5%B0%B8%E6%9D%91%E6%B0%91 "僵尸村民")[僵尸村民](/w/%E5%83%B5%E5%B0%B8%E6%9D%91%E6%B0%91 "僵尸村民")

[wolf/collar](#wolf/collar)

[![](/images/EntitySprite_pale-wolf.png?85237)](/w/%E7%8B%BC "狼")[狼](/w/%E7%8B%BC "狼")

![字节型](/images/Data_node_byte.svg?eb0e0)CollarColor

[wolf/sound\_variant](#wolf/sound_variant)

![字符串](/images/Data_node_string.svg?42545)sound\_variant

[wolf/variant](#wolf/variant)

![字符串](/images/Data_node_string.svg?42545)variant

## 数据组件类型

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=6&veaction=edit "编辑章节：数据组件类型") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=6 "编辑章节的源代码： 数据组件类型")\]

游戏总共定义了下列数据组件。此处仅简要介绍，完整格式和作用见下文。

-   [![](/images/ItemSprite_wooden-spear.png?e2d0e)](#attack_range)[attack\_range](#attack_range)（[攻击距离](/w/%E6%94%BB%E5%87%BB%E8%B7%9D%E7%A6%BB "攻击距离")）
-   [![](/images/BlockSprite_chain-command-block.png?82b71)](#attribute_modifiers)[attribute\_modifiers](#attribute_modifiers)（[属性修饰符](/w/%E5%B1%9E%E6%80%A7%E4%BF%AE%E9%A5%B0%E7%AC%A6 "属性修饰符")）
-   [![](/images/BlockSprite_white-banner.png?ab75c)](#banner_patterns)[banner\_patterns](#banner_patterns)（[旗帜图案](/w/%E6%97%97%E5%B8%9C%E5%9B%BE%E6%A1%88 "旗帜图案")）
-   [![](/images/ItemSprite_shield.png?64fc2)](#base_color)[base\_color](#base_color)（[盾牌](/w/%E7%9B%BE%E7%89%8C "盾牌")基色）
-   [![](/images/BlockSprite_bee-nest.png?70047)](#bees)[bees](#bees)（[蜜蜂](/w/%E8%9C%9C%E8%9C%82 "蜜蜂")数据）
-   [![](/images/BlockSprite_monster-spawner.png?65148)](#block_entity_data)[block\_entity\_data](#block_entity_data)（[方块实体数据](/w/%E6%96%B9%E5%9D%97%E5%AE%9E%E4%BD%93%E6%95%B0%E6%8D%AE%E6%A0%BC%E5%BC%8F "方块实体数据格式")）
-   [![](/images/BlockSprite_oak-stairs.png?c204b)](#block_state)[block\_state](#block_state)（[方块状态](/w/%E6%96%B9%E5%9D%97%E7%8A%B6%E6%80%81 "方块状态")）
-   [![](/images/ItemSprite_shield.png?64fc2)](#blocks_attacks)[blocks\_attacks](#blocks_attacks)（[格挡](/w/%E6%A0%BC%E6%8C%A1 "格挡")攻击）
-   [![](/images/ItemSprite_wooden-hoe.png?d882a)](#break_sound)[break\_sound](#break_sound)（物品耐久耗尽音效）
-   [![](/images/ItemSprite_bucket-of-tropical-fish.png?d668e)](#bucket_entity_data)[bucket\_entity\_data](#bucket_entity_data)（[生物桶](/w/%E7%94%9F%E7%89%A9%E6%A1%B6 "生物桶")所装实体数据）
-   [![](/images/ItemSprite_bundle.png?cb35d)](#bundle_contents)[bundle\_contents](#bundle_contents)（[收纳袋](/w/%E6%94%B6%E7%BA%B3%E8%A2%8B "收纳袋")内物品）
-   [![](/images/ItemSprite_stone-pickaxe.png?e9b36)](#can_break)[can\_break](#can_break)（冒险模式下该物品可破坏的方块）
-   [![](/images/BlockSprite_cobble.png?226e6)](#can_place_on)[can\_place\_on](#can_place_on)（冒险模式下该物品可放置于的方块）
-   [![](/images/ItemSprite_crossbow.png?36e66)](#charged_projectiles)[charged\_projectiles](#charged_projectiles)（所装载的弹射物）
-   [![](/images/ItemSprite_golden-apple.png?846ed)](#consumable)[consumable](#consumable)（可消耗性）
-   [![](/images/BlockSprite_shulker-box.png?d84be)](#container)[container](#container)（[容器](/w/%E5%AE%B9%E5%99%A8 "容器")内物品）
-   [![](/images/BlockSprite_chest.png?05052)](#container_loot)[container\_loot](#container_loot)（容器[战利品表](/w/%E6%88%98%E5%88%A9%E5%93%81%E8%A1%A8 "战利品表")）
-   [![](/images/BlockSprite_barrier.png?7d049)](#custom_data)[custom\_data](#custom_data)（自定义数据）
-   [![](/images/ItemSprite_diamond.png?071fc)](#custom_model_data)[custom\_model\_data](#custom_model_data)（自定义模型数据）
-   [![](/images/ItemSprite_name-tag.png?9cd8e)](#custom_name)[custom\_name](#custom_name)（自定义名称）
-   [![](/images/ItemSprite_axe.png?8b022)](#damage)[damage](#damage)（物品损坏值）
-   [![](/images/ItemSprite_netherite-ingot.png?e3701)](#damage_resistant)[damage\_resistant](#damage_resistant)（不被指定伤害类型摧毁）
-   [![](/images/ItemSprite_iron-sword.png?687f2)](#damage_type)[damage\_type](#damage_type)（攻击造成的[伤害类型](/w/%E4%BC%A4%E5%AE%B3%E7%B1%BB%E5%9E%8B "伤害类型")）
-   [![](/images/ItemSprite_debug-stick.png?e37ad)](#debug_stick_state)[debug\_stick\_state](#debug_stick_state)（调试棒状态）
-   [![](/images/ItemSprite_totem-of-undying.png?1a460)](#death_protection)[death\_protection](#death_protection)（死亡保护）
-   [![](/images/ItemSprite_all-dyes.png?9adba)](#dye)[dye](#dye)（染料颜色）\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]
-   [![](/images/ItemSprite_leather-chestplate.png?14a04)](#dyed_color)[dyed\_color](#dyed_color)（所染颜色）
-   [![](/images/ItemSprite_diamond-boots.png?f5bed)](#enchantable)[enchantable](#enchantable)（在附魔台上的[附魔能力](/w/%E9%99%84%E9%AD%94%E8%83%BD%E5%8A%9B "附魔能力")）
-   [![](/images/EntitySprite_bottle-o%27-enchanting.png?5ee9b)](#enchantment_glint_override)[enchantment\_glint\_override](#enchantment_glint_override)（附魔[光效](/w/%E5%85%89%E6%95%88 "光效")）
-   [![](/images/ItemSprite_book.png?520c4)](#enchantments)[enchantments](#enchantments)（[魔咒](/w/%E9%AD%94%E5%92%92 "魔咒")）
-   [![](/images/ItemSprite_armor-stand.png?ea570)](#entity_data)[entity\_data](#entity_data)（实体数据）
-   [![](/images/ItemSprite_saddle.png?d10c2)](#equippable)[equippable](#equippable)（可穿戴性）
-   [![](/images/ItemSprite_firework-star.png?531eb)](#firework_explosion)[firework\_explosion](#firework_explosion)（[烟火之星](/w/%E7%83%9F%E7%81%AB%E4%B9%8B%E6%98%9F "烟火之星")爆裂数据）
-   [![](/images/ItemSprite_firework-rocket.png?2d031)](#fireworks)[fireworks](#fireworks)（[烟花火箭](/w/%E7%83%9F%E8%8A%B1%E7%81%AB%E7%AE%AD "烟花火箭")爆裂和飞行数据）
-   [![](/images/ItemSprite_cooked-beef.png?088c9)](#food)[food](#food)（[食物](/w/%E9%A3%9F%E7%89%A9 "食物")）
-   [![](/images/ItemSprite_elytra.png?2388a)](#glider)[glider](#glider)（穿戴后可[滑翔](/w/%E6%BB%91%E7%BF%94 "滑翔")）
-   [![](/images/ItemSprite_goat-horn.png?75c0e)](#instrument)[instrument](#instrument)（[乐器](/w/%E5%B1%B1%E7%BE%8A%E8%A7%92%E4%B9%90%E5%99%A8%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "山羊角乐器定义格式")）
-   [![](/images/ItemSprite_arrow.png?2eb1c)](#intangible_projectile)[intangible\_projectile](#intangible_projectile)（只能被创造模式玩家捡起的弹射物）
-   [![](/images/ItemSprite_emerald.png?5336a)](#item_model)[item\_model](#item_model)（[物品模型](/w/%E7%89%A9%E5%93%81%E6%A8%A1%E5%9E%8B%E6%98%A0%E5%B0%84 "物品模型映射")）
-   [![](/images/ItemSprite_name-tag.png?9cd8e)](#item_name)[item\_name](#item_name)（物品名称）
-   [![](/images/ItemSprite_music-disc-5.png?151f2)](#jukebox_playable)[jukebox\_playable](#jukebox_playable)（插入[唱片机](/w/%E5%94%B1%E7%89%87%E6%9C%BA "唱片机")并播放音乐）
-   [![](/images/ItemSprite_iron-spear.png?9354e)](#kinetic_weapon)[kinetic\_weapon](#kinetic_weapon)（设置冲锋攻击）
-   [![](/images/BlockSprite_tripwire-hook.png?714c8)](#lock)[lock](#lock)（锁）
-   [![](/images/ItemSprite_lodestone-compass.png?81b14)](#lodestone_tracker)[lodestone\_tracker](#lodestone_tracker)（所追踪的[磁石](/w/%E7%A3%81%E7%9F%B3 "磁石")位置）
-   [![](/images/ItemSprite_paper.png?cefb9)](#lore)[lore](#lore)（物品提示框中的描述信息）
-   [![](/images/ItemSprite_buried-treasure-map.png?c2380)](#map_color)[map\_color](#map_color)（[地图](/w/%E5%9C%B0%E5%9B%BE "地图")在物品栏内的纹理颜色）
-   [![](/images/ItemSprite_ocean-explorer-map.png?66136)](#map_decorations)[map\_decorations](#map_decorations)（地图图标）
-   [![](/images/ItemSprite_buried-treasure-map.png?c2380)](#map_id)[map\_id](#map_id)（地图编号）
-   [![](/images/ItemSprite_diamond-axe.png?fd475)](#max_damage)[max\_damage](#max_damage)（最大耐久度）
-   [![](/images/ItemSprite_egg.png?7cdef)](#max_stack_size)[max\_stack\_size](#max_stack_size)（最大堆叠数）
-   [![](/images/ItemSprite_iron-sword.png?687f2)](#minimum_attack_charge)[minimum\_attack\_charge](#minimum_attack_charge)（进行近战或穿刺攻击所需的[冷却进度](/w/%E8%BF%91%E6%88%98%E6%94%BB%E5%87%BB#攻击冷却 "近战攻击")最小值）
-   [![](/images/BlockSprite_jukebox-side.png?3ad1e)](#note_block_sound)[note\_block\_sound](#note_block_sound)（放有玩家的头的音符盒音效）
-   [![](/images/ItemSprite_ominous-bottle.png?7744e)](#ominous_bottle_amplifier)[ominous\_bottle\_amplifier](#ominous_bottle_amplifier)（物品的[不祥之兆](/w/%E4%B8%8D%E7%A5%A5%E4%B9%8B%E5%85%86 "不祥之兆")状态效果倍率）
-   [![](/images/ItemSprite_iron-spear.png?9354e)](#piercing_weapon)[piercing\_weapon](#piercing_weapon)（设置戳刺攻击）
-   [![](/images/ItemSprite_danger-pottery-shard.png?6496a)](#pot_decorations)[pot\_decorations](#pot_decorations)（[饰纹陶罐](/w/%E9%A5%B0%E7%BA%B9%E9%99%B6%E7%BD%90 "饰纹陶罐")陶片装饰）
-   [![](/images/ItemSprite_awkward-potion.png?398bd)](#potion_contents)[potion\_contents](#potion_contents)（药水效果与状态效果信息）
-   [![](/images/ItemSprite_awkward-lingering-potion.png?a36cf)](#potion_duration_scale)[potion\_duration\_scale](#potion_duration_scale)（状态效果时长缩放倍率）
-   [![](/images/BlockSprite_player-head.png?5ddf1)](#profile)[profile](#profile)（玩家游戏档案信息）
-   [![](/images/ItemSprite_creeper-charge-banner-pattern.png?d451b)](#provides_banner_patterns)[provides\_banner\_patterns](#provides_banner_patterns)（是否可放进织布机的旗帜图案槽位）
-   [![](/images/ItemSprite_raiser-armor-trim.png?5cb22)](#provides_trim_material)[provides\_trim\_material](#provides_trim_material)（为盔甲纹饰配方提供的纹饰材料）
-   [![](/images/ItemSprite_nether-star.png?9b2bd)](#rarity)[rarity](#rarity)（[稀有度](/w/%E7%A8%80%E6%9C%89%E5%BA%A6 "稀有度")）
-   [![](/images/ItemSprite_knowledge-book.png?793c1)](#recipes)[recipes](#recipes)（知识之书配方信息）
-   [![](/images/BlockSprite_anvil.png?a1169)](#repairable)[repairable](#repairable)（可在铁砧上被修复）
-   [![](/images/EntitySprite_experience-orb.png?d5ead)](#repair_cost)[repair\_cost](#repair_cost)（在铁砧上的[累计惩罚值](/w/%E9%93%81%E7%A0%A7%E6%9C%BA%E5%88%B6#累计惩罚 "铁砧机制")）
-   [![](/images/ItemSprite_enchanted-book.png?28fff)](#stored_enchantments)[stored\_enchantments](#stored_enchantments)（所存储的“无活性”魔咒）
-   [![](/images/ItemSprite_suspicious-stew.png?762dc)](#suspicious_stew_effects)[suspicious\_stew\_effects](#suspicious_stew_effects)（[迷之炖菜](/w/%E8%BF%B7%E4%B9%8B%E7%82%96%E8%8F%9C "迷之炖菜")效果）
-   [![](/images/ItemSprite_iron-sword.png?687f2)](#swing_animation)[swing\_animation](#swing_animation)（攻击动画）
-   [![](/images/ItemSprite_diamond-shovel.png?4600f)](#tool)[tool](#tool)（成为挖掘某方块的工具）
-   [![](/images/ItemSprite_item-frame.png?a2327)](#tooltip_display)[tooltip\_display](#tooltip_display)（物品提示框及附加信息的显示）
-   [![](/images/ItemSprite_painting.png?88c09)](#tooltip_style)[tooltip\_style](#tooltip_style)（物品提示框背景和边框样式）
-   [![](/images/ItemSprite_spire-armor-trim.png?c5f27)](#trim)[trim](#trim)（盔甲纹饰）
-   [![](/images/BlockSprite_bedrock.png?c6a65)](#unbreakable)[unbreakable](#unbreakable)（无法破坏）
-   [![](/images/ItemSprite_ender-pearl.png?72ab2)](#use_cooldown)[use\_cooldown](#use_cooldown)（使用后冷却）
-   [![](/images/ItemSprite_diamond-spear.png?f2e6b)](#use_effects)[use\_effects](#use_effects)（玩家使用物品时的行为）
-   [![](/images/ItemSprite_milk.png?8fc08)](#use_remainder)[use\_remainder](#use_remainder)（使用后返还物品）
-   [![](/images/ItemSprite_diamond-sword.png?434c7)](#weapon)[weapon](#weapon)（作为武器时的行为）
-   [![](/images/ItemSprite_book-and-quill.png?a97f6)](#writable_book_content)[writable\_book\_content](#writable_book_content)（书与笔内容）
-   [![](/images/ItemSprite_written-book.png?6794b)](#written_book_content)[written\_book\_content](#written_book_content)（成书内容）

[实体变种组件](#实体变种组件)

-   [![](/images/EntitySprite_axolotl.png?2597f)](#axolotl/variant)[axolotl/variant](#axolotl/variant)（美西螈变种）
-   [![](/images/EntitySprite_cat.png?4d91c)](#cat/collar)[cat/collar](#cat/collar)（猫项圈颜色）
-   [![](/images/EntitySprite_cat.png?4d91c)](#cat/sound_variant)[cat/sound\_variant](#cat/sound_variant)\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]（猫音效变种）
-   [![](/images/EntitySprite_cat.png?4d91c)](#cat/variant)[cat/variant](#cat/variant)（猫变种）
-   [![](/images/EntitySprite_chicken.png?ed28c)](#chicken/sound_variant)[chicken/sound\_variant](#chicken/sound_variant)\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]（鸡音效变种）
-   [![](/images/EntitySprite_chicken.png?ed28c)](#chicken/variant)[chicken/variant](#chicken/variant)（鸡变种）
-   [![](/images/EntitySprite_cow.png?54276)](#cow/sound_variant)[cow/sound\_variant](#cow/sound_variant)\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]（牛音效变种）
-   [![](/images/EntitySprite_cow.png?54276)](#cow/variant)[cow/variant](#cow/variant)（牛变种）
-   [![](/images/EntitySprite_fox.png?c5f58)](#fox/variant)[fox/variant](#fox/variant)（狐狸变种）
-   [![](/images/EntitySprite_frog.png?39d41)](#frog/variant)[frog/variant](#frog/variant)（青蛙变种）
-   [![](/images/EntitySprite_creamy-horse.png?80b78)](#horse/variant)[horse/variant](#horse/variant)（马变种）
-   [![](/images/EntitySprite_creamy-llama.png?b264d)](#llama/variant)[llama/variant](#llama/variant)（羊驼变种）
-   [![](/images/EntitySprite_mooshroom.png?837a6)](#mooshroom/variant)[mooshroom/variant](#mooshroom/variant)（哞菇变种）
-   [![](/images/EntitySprite_parrot.png?ed0cc)](#parrot/variant)[parrot/variant](#parrot/variant)（鹦鹉变种）
-   [![](/images/EntitySprite_kebab.png?15a6b)](#painting/variant)[painting/variant](#painting/variant)（画变种）
-   [![](/images/EntitySprite_pig.png?c2459)](#pig/sound_variant)[pig/sound\_variant](#pig/sound_variant)\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]（猪音效变种）
-   [![](/images/EntitySprite_pig.png?c2459)](#pig/variant)[pig/variant](#pig/variant)（猪变种）
-   [![](/images/EntitySprite_brown-rabbit.png?af02d)](#rabbit/variant)[rabbit/variant](#rabbit/variant)（兔子变种）
-   [![](/images/EntitySprite_salmon.png?913ea)](#salmon/size)[salmon/size](#salmon/size)（鲑鱼体型尺寸）
-   [![](/images/EntitySprite_white-sheep.png?07a4a)](#sheep/color)[sheep/color](#sheep/color)（绵羊变种）
-   [![](/images/EntitySprite_shulker.png?a8d4f)](#shulker/color)[shulker/color](#shulker/color)（潜影贝颜色）
-   [![](/images/EntitySprite_tropical-fish.png?20777)](#tropical_fish/base_color)[tropical\_fish/base\_color](#tropical_fish/base_color)（热带鱼体色）
-   [![](/images/EntitySprite_tropical-fish.png?20777)](#tropical_fish/pattern)[tropical\_fish/pattern](#tropical_fish/pattern)（热带鱼花纹）
-   [![](/images/EntitySprite_tropical-fish.png?20777)](#tropical_fish/pattern_color)[tropical\_fish/pattern\_color](#tropical_fish/pattern_color)（热带鱼花纹颜色）
-   [![](/images/EntitySprite_leatherworker.png?43c3d)](#villager/variant)[villager/variant](#villager/variant)（村民变种）
-   [![](/images/EntitySprite_pale-wolf.png?85237)](#wolf/collar)[wolf/collar](#wolf/collar)（狼项圈颜色）
-   [![](/images/EntitySprite_pale-wolf.png?85237)](#wolf/sound_variant)[wolf/sound\_variant](#wolf/sound_variant)（狼音效变种）
-   [![](/images/EntitySprite_pale-wolf.png?85237)](#wolf/variant)[wolf/variant](#wolf/variant)（狼变种）

[临时组件](#临时组件)

-   [additional\_trade\_cost](#additional_trade_cost)\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]
-   [creative\_slot\_lock](#creative_slot_lock)
-   [map\_post\_processing](#map_post_processing)

## 命令格式

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=7&veaction=edit "编辑章节：命令格式") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=7 "编辑章节的源代码： 命令格式")\]

`item_stack`和`item_predicate`[参数类型](/w/%E5%8F%82%E6%95%B0%E7%B1%BB%E5%9E%8B "参数类型")支持物品堆叠组件。

`item_stack`参数类型可以加载物品组件，也可通过在组件名前添加`!`来移除该物品的默认组件，在`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give")`等命令中使用。所提供的组件都会被设置，而未提供的组件会被设为默认值。格式参见[参数类型 § item\_stack](/w/%E5%8F%82%E6%95%B0%E7%B1%BB%E5%9E%8B#item_stack "参数类型")。

`item_predicate`参数类型可以检测物品组件，在`/[clear](/w/%E5%91%BD%E4%BB%A4/clear "命令/clear")`等命令中使用。另外，该参数类型还可以直接使用[数据组件谓词](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6%E8%B0%93%E8%AF%8D "数据组件谓词")检测物品堆叠组件。格式参见[参数类型 § item\_predicate](/w/%E5%8F%82%E6%95%B0%E7%B1%BB%E5%9E%8B#item_predicate "参数类型")。

## 数据格式

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=8&veaction=edit "编辑章节：数据格式") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=8 "编辑章节的源代码： 数据格式")\]

### attack\_range

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=9&veaction=edit "编辑章节：attack_range") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=9 "编辑章节的源代码： attack_range")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:attack\_range：生物持有此物品时的攻击距离，会覆写玩家的[实体交互距离](/w/%E5%B1%9E%E6%80%A7/%E5%AE%9E%E4%BD%93%E4%BA%A4%E4%BA%92%E8%B7%9D%E7%A6%BB "属性/实体交互距离")属性。
        -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)min\_reach：（0≤值≤64，默认为0）攻击者到目标的最小有效距离。以攻击者眼睛位置、沿视角方向到被攻击者攻击判定箱的最小距离计算。
        -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)max\_reach：（0≤值≤64，默认为3）攻击者到目标的最大有效距离。以攻击者眼睛位置、沿视角方向到被攻击者攻击判定箱的最小距离计算。
        -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)min\_creative\_reach：（0≤值≤64，默认为0）创造模式玩家到目标的最小有效距离，计算方式同上。
        -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)max\_creative\_reach：（0≤值≤64，默认为5）创造模式玩家到目标的最大有效距离，计算方式同上。
        -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)hitbox\_margin：（0≤值≤1，默认为0.3）决定攻击判定箱的大小。游戏将实体的碰撞箱向各个方向扩展此距离得到攻击判定箱。
        -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)mob\_factor：（0≤值≤2，默认为1.0）对于非玩家生物，其使用的最小有效距离和最大有效距离的缩放乘数。

### attribute\_modifiers

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=10&veaction=edit "编辑章节：attribute_modifiers") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=10 "编辑章节的源代码： attribute_modifiers")\]

存储修饰生物属性的[属性修饰符](/w/%E5%B1%9E%E6%80%A7%E4%BF%AE%E9%A5%B0%E7%AC%A6 "属性修饰符")，当物品在生物的指定槽位上时可以修改其所在生物的属性。物品存储的属性修饰符信息会在[物品提示框](/w/%E7%89%A9%E5%93%81%E6%8F%90%E7%A4%BA%E6%A1%86 "物品提示框")中显示。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)minecraft:attribute\_modifiers：物品为持有者提供的[属性修饰符](/w/%E5%B1%9E%E6%80%A7%E4%BF%AE%E9%A5%B0%E7%AC%A6 "属性修饰符")。
        -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)：一个修饰符。
            -   ![双精度浮点数](/images/Data_node_double.svg?14320)\*  
                \*amount：计算中修饰符调整基础值的数值。
            -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)display：属性修饰符在提示框的显示方式。
                -   ![字符串](/images/Data_node_string.svg?42545)\*  
                    \*type：显示类型，枚举值见下。
                
                -   当![字符串](/images/Data_node_string.svg?42545)type为`default`时，显示此项计算后的属性修饰符值。此项也为默认值。
                
                -   当![字符串](/images/Data_node_string.svg?42545)type为`hidden`时，不显示此项属性修饰符值。
                
                -   当![字符串](/images/Data_node_string.svg?42545)type为`override`时，替换所显示的属性修饰符文本，附加字段如下：
                -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)\*  
                    \*value：（[文本组件](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6 "文本组件")）替换后的文本。
            -   ![字符串](/images/Data_node_string.svg?42545)\*  
                \*id：（命名空间ID）当前[属性修饰符](/w/%E5%B1%9E%E6%80%A7%E4%BF%AE%E9%A5%B0%E7%AC%A6 "属性修饰符")的ID。
            -   ![字符串](/images/Data_node_string.svg?42545)\*  
                \*operation：定义修饰符对属性的基础值的[运算方法](/w/%E5%B1%9E%E6%80%A7#运算模式 "属性")。可以为`add_value`（Op0）、`add_multiplied_base`（Op1）、`add_multiplied_total`（Op2）。
            -   ![字符串](/images/Data_node_string.svg?42545)slot：（默认为`any`）一个[装备槽位组](/w/%E8%A3%85%E5%A4%87%E6%A7%BD%E4%BD%8D%E7%BB%84 "装备槽位组")，指定修饰符的有效槽位。
            -   ![字符串](/images/Data_node_string.svg?42545)\*  
                \*type：（命名空间ID）一个[属性](/w/%E5%B1%9E%E6%80%A7 "属性")的ID，表示当前属性修饰符要修饰的属性。

示例

给予一个木棍，玩家手持该木棍时，玩家的[尺寸](/w/%E5%B1%9E%E6%80%A7/%E5%B0%BA%E5%AF%B8 "属性/尺寸")会增加4倍（即原来的5倍大小）。物品提示框中将显示被`example:grow`属性修饰符修饰后的`scale`属性值。

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s stick[attribute_modifiers=[{type:"scale",slot:"hand",id:"example:grow",amount:4,operation:"add_multiplied_base"}]]`

### banner\_patterns

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=11&veaction=edit "编辑章节：banner_patterns") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=11 "编辑章节的源代码： banner_patterns")\]

存储[旗帜](/w/%E6%97%97%E5%B8%9C "旗帜")和[盾牌](/w/%E7%9B%BE%E7%89%8C "盾牌")上的[旗帜图案](/w/%E6%97%97%E5%B8%9C#方块实体 "旗帜")。旗帜图案信息会在提示框中显示。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)minecraft:banner\_patterns：旗帜图案的列表。
        -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)：一层图案。
            -   ![字符串](/images/Data_node_string.svg?42545)\*  
                \*color：这一层图案的颜色。取值为[染料颜色](/w/%E6%9F%93%E6%96%99%E9%A2%9C%E8%89%B2 "染料颜色")，即`white`、`orange`、`magenta`、`light_blue`、`yellow`、`lime`、`pink`、`gray`、`light_gray`、`cyan`、`purple`、`blue`、`brown`、`green`、`red`或`black`。
            -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)\*  
                \*pattern：这一层图案的样式。可以为旗帜图案的ID，也可以是旗帜图案的内联格式，见[旗帜 § 方块实体](/w/%E6%97%97%E5%B8%9C#方块实体 "旗帜")。
                
                -   旗帜图案，见[Template:Nbt inherit/banner pattern/source](/w/Template:Nbt_inherit/banner_pattern/source "Template:Nbt inherit/banner pattern/source")

### base\_color

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=12&veaction=edit "编辑章节：base_color") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=12 "编辑章节的源代码： base_color")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:base\_color：[盾牌](/w/%E7%9B%BE%E7%89%8C "盾牌")的基础颜色，同时影响盾牌的名称。取值为[染料颜色](/w/%E6%9F%93%E6%96%99%E9%A2%9C%E8%89%B2 "染料颜色")，即`white`、`orange`、`magenta`、`light_blue`、`yellow`、`lime`、`pink`、`gray`、`light_gray`、`cyan`、`purple`、`blue`、`brown`、`green`、`red`或`black`。

示例

给予玩家一个盾牌基色为黄绿色的盾牌：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s shield[base_color=lime]`

### bees

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=13&veaction=edit "编辑章节：bees") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=13 "编辑章节的源代码： bees")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)minecraft:bees：[蜂巢（方块）](/w/%E8%9C%82%E5%B7%A2%EF%BC%88%E6%96%B9%E5%9D%97%EF%BC%89 "蜂巢（方块）")和[蜂箱](/w/%E8%9C%82%E7%AE%B1 "蜂箱")的蜜蜂数据。
        -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)：一只蜜蜂的数据。
            -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)entity\_data：蜜蜂的部分实体数据。如果采用字符串格式进行定义，则游戏会将字符串的内容视为[SNBT](/w/SNBT "SNBT")加载，游戏只保存为复合标签格式。
                -   见[实体数据格式](/w/%E5%AE%9E%E4%BD%93%E6%95%B0%E6%8D%AE%E6%A0%BC%E5%BC%8F "实体数据格式")。下列标签不会被保存，也不会被加载：![短整型](/images/Data_node_short.svg?c1f72)Air、![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)drop\_chances、![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)equipment、![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)Brain、![布尔型](/images/Data_node_bool.svg?77754)CanPickUpLoot、![短整型](/images/Data_node_short.svg?c1f72)DeathTime、![单精度浮点数](/images/Data_node_float.svg?ae55e)fall\_distance、![布尔型](/images/Data_node_bool.svg?77754)FallFlying、![短整型](/images/Data_node_short.svg?c1f72)Fire、![整型](/images/Data_node_int.svg?8d24f)HurtByTimestamp、![短整型](/images/Data_node_short.svg?c1f72)HurtTime、![布尔型](/images/Data_node_bool.svg?77754)LeftHanded、![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)Motion、![布尔型](/images/Data_node_bool.svg?77754)NoGravity、![布尔型](/images/Data_node_bool.svg?77754)OnGround、![整型](/images/Data_node_int.svg?8d24f)PortalCooldown、![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)Pos、![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)Rotation、![整型数组](/images/Data_node_int-array.svg?546e8)sleeping\_pos、![整型](/images/Data_node_int.svg?8d24f)CannotEnterHiveTicks、![整型](/images/Data_node_int.svg?8d24f)TicksSincePollination、![整型](/images/Data_node_int.svg?8d24f)CropsGrownSincePollination、![整型数组](/images/Data_node_int-array.svg?546e8)hive\_pos、![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)Passengers、![整型数组](/images/Data_node_int-array.svg?546e8)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)leash、![整型数组](/images/Data_node_int-array.svg?546e8)UUID。
            -   ![整型](/images/Data_node_int.svg?8d24f)\*  
                \*min\_ticks\_in\_hive：蜜蜂会在巢内滞留的最短时间。
            -   ![整型](/images/Data_node_int.svg?8d24f)\*  
                \*ticks\_in\_hive：蜜蜂在巢内已滞留的时间。

### block\_entity\_data

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=14&veaction=edit "编辑章节：block_entity_data") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=14 "编辑章节的源代码： block_entity_data")\]

存储应用于方块实体的方块实体数据，放置诸如箱子或熔炉等具有对应方块实体的方块时将加载该组件的数据。

若放置的方块为指定了附加数据的任意[命令方块](/w/%E5%91%BD%E4%BB%A4%E6%96%B9%E5%9D%97 "命令方块")、[讲台](/w/%E8%AE%B2%E5%8F%B0 "讲台")、任意[告示牌](/w/%E5%91%8A%E7%A4%BA%E7%89%8C "告示牌")、任意[悬挂式告示牌](/w/%E6%82%AC%E6%8C%82%E5%BC%8F%E5%91%8A%E7%A4%BA%E7%89%8C "悬挂式告示牌")、[刷怪笼](/w/%E5%88%B7%E6%80%AA%E7%AC%BC "刷怪笼")或[试炼刷怪笼](/w/%E8%AF%95%E7%82%BC%E5%88%B7%E6%80%AA%E7%AC%BC "试炼刷怪笼")，则非管理员玩家使用这些物品时不会设置方块实体数据，且提示框中会显示安全警告。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:block\_entity\_data：物品放置方块时将套用到方块实体的数据。如果采用字符串格式进行定义，则游戏会将字符串的内容视为[SNBT](/w/SNBT "SNBT")加载，游戏只保存为复合标签格式。
        -   ![字符串](/images/Data_node_string.svg?42545)\*  
            \*id：（命名空间ID）方块实体。
        -   若干与该方块对应的方块实体数据标签，见[方块实体数据格式](/w/%E6%96%B9%E5%9D%97%E5%AE%9E%E4%BD%93%E6%95%B0%E6%8D%AE%E6%A0%BC%E5%BC%8F "方块实体数据格式")。

示例

给予一个蜘蛛刷怪笼。要放置该刷怪笼，玩家必须要有[管理员权限](/w/%E6%9D%83%E9%99%90%E7%AD%89%E7%BA%A7 "权限等级")：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s spawner[block_entity_data={id:"mob_spawner",SpawnData:{entity:{id:"spider"}}}]`

### block\_state

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=15&veaction=edit "编辑章节：block_state") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=15 "编辑章节的源代码： block_state")\]

存储方块物品被放置时将应用于方块的方块状态。未指定的属性依旧使用默认值，如果方块属性对于被放置的方块不存在或对应的方块属性值无效，则这项设置不起任何作用。指定属性为`honey_level`时提示框会显示蜂蜜等级信息。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:block\_state：物品放置方块时将要设置的[方块状态](/w/%E6%96%B9%E5%9D%97%E7%8A%B6%E6%80%81 "方块状态")。
        -   ![字符串](/images/Data_node_string.svg?42545)<*方块属性*\>：此项方块属性的值。

示例

给予一个被放置时总位于方块网格的上半部分的竹台阶：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s bamboo_slab[block_state={type:"top"}]`

### blocks\_attacks

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=16&veaction=edit "编辑章节：blocks_attacks") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=16 "编辑章节的源代码： blocks_attacks")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:blocks\_attacks：物品使用时的格挡行为。
        -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)block\_delay\_seconds：（值≥0，默认为`0`）成功阻挡攻击前需要按住右键的秒数。
        -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)block\_sound：成功阻挡攻击时播放的声音事件。
            
            -   声音事件，见[Template:Nbt inherit/sound event/source](/w/Template:Nbt_inherit/sound_event/source "Template:Nbt inherit/sound event/source")
        -   ![字符串](/images/Data_node_string.svg?42545)bypassed\_by\[失效：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]：（命名空间ID）可以无视此物品的阻挡而造成实际伤害的伤害类型。应为一个带`#`前缀的标签ID，游戏会将此值解析为[伤害类型标签](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E4%BC%A4%E5%AE%B3%E7%B1%BB%E5%9E%8B "Java版标签/伤害类型")，如果伤害类型标签不存在则可以阻挡任何伤害。
        -   ![字符串](/images/Data_node_string.svg?42545)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)bypassed\_by\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]：可以无视此物品的阻挡而造成实际伤害的伤害类型。可以为一个伤害类型ID、一个伤害类型标签ID，或一个伤害类型标签。
        -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)damage\_reductions：控制可阻挡多少伤害。未指定时，可阻挡一切伤害。
            -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)：控制可挡下的伤害量和伤害类型。阻挡成功时，伤害减少`clamp(base + factor * *所受攻击伤害*, 0, *所受攻击伤害*)`。
                -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)\*  
                    \*base：固定阻挡的伤害。
                -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)\*  
                    \*factor：应被阻挡的伤害比例。
                -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)horizontal\_blocking\_angle：（值>0，角度制，默认为`90`）在水平方向上，以当前玩家视角的水平分量向量为基准，如果受伤害方向与基准方向夹角小于此角度则伤害可被阻挡，否则不能阻挡。  
                    任何无来源伤害均被视为需要`180`度才能阻挡。
                -   ![字符串](/images/Data_node_string.svg?42545)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)type：可阻挡的伤害类型。允许单个ID、列表或标签。未指定则表示对所有伤害有效。
        -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)disable\_cooldown\_scale：（值≥0，默认为`1`）被可停用阻挡的攻击击中时，物品冷却时长的乘数。为`0`时，此物品不能被攻击停用。
        -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)disabled\_sound：此物品被攻击停用时播放的声音事件。
            
            -   声音事件，见[Template:Nbt inherit/sound event/source](/w/Template:Nbt_inherit/sound_event/source "Template:Nbt inherit/sound event/source")
        -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)item\_damage：控制攻击对物品造成的耐久损耗。未指定时，每次攻击损耗物品1点耐久。物品耐久最终损耗`floor(threshold, base + factor * *所受攻击伤害*)`。最终值可以为负数以使物品修复。
            -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)\*  
                \*base：损耗物品固定耐久度。
            -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)\*  
                \*factor：所受攻击伤害的乘数。
            -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)\*  
                \*threshold：（值≥0）攻击对此物品造成的最低耐久度损耗。

示例

给予一个[弓](/w/%E5%BC%93 "弓")，玩家使用此弓的同时会格挡前方所有类型所有伤害的攻击（因为使用物品和格挡的按键均为鼠标右键），且不会被[可停用阻挡的攻击](/w/%E7%9B%BE%E7%89%8C#停用 "盾牌")停用：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s bow[blocks_attacks={disable_cooldown_scale:0}]`

### break\_sound

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=17&veaction=edit "编辑章节：break_sound") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=17 "编辑章节的源代码： break_sound")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:break\_sound：物品耐久度耗尽时播放的声音事件。
        
        -   声音事件，见[Template:Nbt inherit/sound event/source](/w/Template:Nbt_inherit/sound_event/source "Template:Nbt inherit/sound event/source")

示例

给予一个物品破坏音效为紫水晶块被破坏的音效的金锹：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @a minecraft:golden_shovel[minecraft:break_sound={sound_id:block.amethyst_block.break}]`

### bucket\_entity\_data

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=18&veaction=edit "编辑章节：bucket_entity_data") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=18 "编辑章节的源代码： bucket_entity_data")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:bucket\_entity\_data：生物桶对桶中生物的部分实体数据。如果采用字符串格式进行定义，则游戏会将字符串的内容视为[SNBT](/w/SNBT "SNBT")加载，游戏只保存为复合标签格式。
        -   ![布尔型](/images/Data_node_bool.svg?77754)Glowing：表示桶中生物是否有发光的轮廓线。
        -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)Health：桶中生物的生命值。
        -   ![布尔型](/images/Data_node_bool.svg?77754)Invulnerable：表示桶中生物是否能抵抗绝大多数伤害。
        -   ![布尔型](/images/Data_node_bool.svg?77754)NoAI：表示桶中生物的AI是否被禁用。
        -   ![布尔型](/images/Data_node_bool.svg?77754)NoGravity：表示桶中生物是否不受重力影响。
        -   ![布尔型](/images/Data_node_bool.svg?77754)Silent：表示桶中生物是否不会发出任何声音。
        
        -   如果桶中生物是[蝌蚪](/w/%E8%9D%8C%E8%9A%AA "蝌蚪")，则有下列1\[失效：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]/2\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]个额外标签：
        -   ![整型](/images/Data_node_int.svg?8d24f)Age：桶中蝌蚪的年龄。大于等于24000时，蝌蚪会长大成青蛙。
        -   ![布尔型](/images/Data_node_bool.svg?77754)AgeLocked\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]：表示蝌蚪的年龄是否不会随时间自然增长。
        
        -   如果桶中生物是[美西螈](/w/%E7%BE%8E%E8%A5%BF%E8%9E%88 "美西螈")，则有下列2\[失效：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]/3\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]个额外标签：
        -   ![整型](/images/Data_node_int.svg?8d24f)Age：桶中美西螈的年龄。生物为幼体时为负值；生物为成体时为正值或0，如果为正值则表示距离生物能再次繁衍的时间。
        -   ![布尔型](/images/Data_node_bool.svg?77754)AgeLocked\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]：表示美西螈的年龄是否不会随时间自然增长或减少。
        -   ![长整型](/images/Data_node_long.svg?dde3f)HuntingCooldown：桶中美西螈[生物记忆](/w/%E7%94%9F%E7%89%A9%E8%AE%B0%E5%BF%86 "生物记忆")`has_hunting_cooldown`的过期倒计时。

### bundle\_contents

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=19&veaction=edit "编辑章节：bundle_contents") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=19 "编辑章节的源代码： bundle_contents")\]

存储[收纳袋](/w/%E6%94%B6%E7%BA%B3%E8%A2%8B "收纳袋")内部包含的物品。拥有此组件的物品实体被摧毁时会释放内容物，如果此组件不存在则收纳袋不会在提示框内显示容量条且不能保存物品。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)minecraft:bundle\_contents：[收纳袋](/w/%E6%94%B6%E7%BA%B3%E8%A2%8B "收纳袋")的内部物品栏。
        -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)\[失效：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]：一个物品。后加入的物品在列表前方，先加入的物品在列表后方。
            
            -   物品共通标签，见[Template:Nbt inherit/itemnoslot/source](/w/Template:Nbt_inherit/itemnoslot/source "Template:Nbt inherit/itemnoslot/source")
        -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]：一个物品。后加入的物品在列表前方，先加入的物品在列表后方。
            
            -   物品模板，见[Template:Nbt inherit/item template/source](/w/Template:Nbt_inherit/item_template/source "Template:Nbt inherit/item template/source")

示例

给予一个收纳袋，其中收纳袋内的物品从列表前方到后方分别为铜锭、铁锭、金锭：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s bundle[bundle_contents=[{id:"copper_ingot"},{id:"iron_ingot"},{id:"gold_ingot"}]]`

### can\_break和can\_place\_on

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=20&veaction=edit "编辑章节：can_break和can_place_on") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=20 "编辑章节的源代码： can_break和can_place_on")\]

控制[冒险模式](/w/%E5%86%92%E9%99%A9%E6%A8%A1%E5%BC%8F "冒险模式")玩家能否破坏指定方块或与指定方块交互，可互动方块会在提示框中提示。如果存在此组件但方块谓词未指定或不满足条件，则显示于提示框的方块为“未知”，且此物品可与任何方块互动。游戏不会测试方块实体组件。

`can_break`组件还可以触发[红石矿石](/w/%E7%BA%A2%E7%9F%B3%E7%9F%BF%E7%9F%B3 "红石矿石")、[龙蛋](/w/%E9%BE%99%E8%9B%8B "龙蛋")或[音符盒](/w/%E9%9F%B3%E7%AC%A6%E7%9B%92 "音符盒")的挖掘开始时效果。

指定列表时不能是空列表，且只有一个元素时游戏只保存为复合标签形式。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)minecraft:can\_break：检查被破坏的方块是否满足指定的方块谓词，作为列表时内部元素与此标签作为复合标签时相同。
        
        -   方块谓词，见[Template:Nbt inherit/block predicate/source](/w/Template:Nbt_inherit/block_predicate/source "Template:Nbt inherit/block predicate/source")
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)minecraft:can\_place\_on：检查被交互的方块是否满足指定的方块谓词，作为列表时内部元素与此标签作为复合标签时相同。
        
        -   方块谓词，见[Template:Nbt inherit/block predicate/source](/w/Template:Nbt_inherit/block_predicate/source "Template:Nbt inherit/block predicate/source")

示例

给予一把冒险模式下仅能挖掘一些矿石的金镐：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s golden_pickaxe[can_break={blocks:['copper_ore','coal_ore','iron_ore','gold_ore','diamond_ore','emerald_ore']}]`

给予一个冒险模式下仅能放置在砂岩上的石头： `/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s stone[can_place_on={blocks:'sandstone'}]`

### charged\_projectiles

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=21&veaction=edit "编辑章节：charged_projectiles") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=21 "编辑章节的源代码： charged_projectiles")\]

存储[弩](/w/%E5%BC%A9 "弩")装载的物品信息。此组件的所有物品将在提示框中显示，连续的相同内容物会合并显示。若物品列表存在烟花火箭则弩显示为“装填烟花火箭的弩”，否则为“装填箭的弩”。此组件不存在时代表弩没有装载任何物品。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)minecraft:charged\_projectiles：[弩](/w/%E5%BC%A9 "弩")的内部物品栏，表示弩的装填物。
        -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)\[失效：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]：一个物品。
            
            -   物品共通标签，见[Template:Nbt inherit/itemnoslot/source](/w/Template:Nbt_inherit/itemnoslot/source "Template:Nbt inherit/itemnoslot/source")
        -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]：一个物品。
            
            -   物品模板，见[Template:Nbt inherit/item template/source](/w/Template:Nbt_inherit/item_template/source "Template:Nbt inherit/item template/source")

### consumable

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=22&veaction=edit "编辑章节：consumable") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=22 "编辑章节的源代码： consumable")\]

控制物品是否具有消耗使用行为，以及使用后的效果。此处的消耗使用指物品数量在使用后会减少的操作，不包含方块物品的放置等对方块进行的有效操作。

这些物品即使使用此组件，也不能被消耗使用：[船](/w/%E8%88%B9 "船")、[运输船](/w/%E8%BF%90%E8%BE%93%E8%88%B9 "运输船")[\[3\]](#cite_note-3)；满足放置条件的[矿车](/w/%E7%9F%BF%E8%BD%A6 "矿车")、[漏斗矿车](/w/%E6%BC%8F%E6%96%97%E7%9F%BF%E8%BD%A6 "漏斗矿车")、[命令方块矿车](/w/%E5%91%BD%E4%BB%A4%E6%96%B9%E5%9D%97%E7%9F%BF%E8%BD%A6 "命令方块矿车")、[运输矿车](/w/%E8%BF%90%E8%BE%93%E7%9F%BF%E8%BD%A6 "运输矿车")、[动力矿车](/w/%E5%8A%A8%E5%8A%9B%E7%9F%BF%E8%BD%A6 "动力矿车")和[TNT矿车](/w/TNT%E7%9F%BF%E8%BD%A6 "TNT矿车")[\[4\]](#cite_note-4)；[成书](/w/%E6%88%90%E4%B9%A6 "成书")；[三叉戟](/w/%E4%B8%89%E5%8F%89%E6%88%9F "三叉戟")[\[5\]](#cite_note-5)；[刷子](/w/%E5%88%B7%E5%AD%90 "刷子")[\[6\]](#cite_note-6)；[铁桶](/w/%E9%93%81%E6%A1%B6 "铁桶")、[水桶](/w/%E6%B0%B4%E6%A1%B6 "水桶")、[熔岩桶](/w/%E7%86%94%E5%B2%A9%E6%A1%B6 "熔岩桶")、[鳕鱼桶](/w/%E9%B3%95%E9%B1%BC%E6%A1%B6 "鳕鱼桶")、[鲑鱼桶](/w/%E9%B2%91%E9%B1%BC%E6%A1%B6 "鲑鱼桶")、[河豚桶](/w/%E6%B2%B3%E8%B1%9A%E6%A1%B6 "河豚桶")、[热带鱼桶](/w/%E7%83%AD%E5%B8%A6%E9%B1%BC%E6%A1%B6 "热带鱼桶")、[美西螈桶](/w/%E7%BE%8E%E8%A5%BF%E8%9E%88%E6%A1%B6 "美西螈桶")、[蝌蚪桶](/w/%E8%9D%8C%E8%9A%AA%E6%A1%B6 "蝌蚪桶")[\[7\]](#cite_note-7)；[弓](/w/%E5%BC%93 "弓")、[弩](/w/%E5%BC%A9 "弩")[\[8\]](#cite_note-8)；[烟花火箭](/w/%E7%83%9F%E8%8A%B1%E7%81%AB%E7%AE%AD "烟花火箭")[\[9\]](#cite_note-9)和所有[刷怪蛋](/w/%E5%88%B7%E6%80%AA%E8%9B%8B "刷怪蛋")[\[10\]](#cite_note-10)。

若物品同时具有`food`、​`ominous_bottle_amplifier`、​`potion_contents`和​`suspicious_stew_effects`等组件，则这些组件的效果也会一并应用。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:consumable：物品的消耗使用行为。
        -   ![字符串](/images/Data_node_string.svg?42545)animation：（默认为`eat`）物品使用时的动画。可以为`none`（无动作）、`eat`（吃）、`drink`（饮用）、`block`（格挡）、`bow`（拉弓）、`brush`（清刷）、`crossbow`（弩上弦）、`spear`（矛蓄力）、`trident`（三叉戟投掷）、`spyglass`（看望远镜）、`toot_horn`（吹山羊角）和`bundle`（使用收纳袋）。
        -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)consume\_seconds：（值≥0，默认为1.6）物品使用的时间，单位为秒。当此值为0时，物品立刻使用，不会像拉弓等操作需要等待时间。
        -   ![布尔型](/images/Data_node_bool.svg?77754)has\_consume\_particles：（默认为`true`）物品在使用时是否产生物品破碎粒子。
        -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)on\_consume\_effects：当物品被使用后产生的效果列表。
            -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)：一项消耗使用效果。
                -   ![字符串](/images/Data_node_string.svg?42545)\*  
                    \*type：消耗使用效果类型。
                
                -   如果![字符串](/images/Data_node_string.svg?42545)type为`apply_effects`，则对使用此物品的生物添加状态效果：
                -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)\*  
                    \*effects：物品使用后添加的状态效果。
                    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)：一项状态效果。
                        
                        -   状态效果，见[Template:Nbt inherit/effect/source](/w/Template:Nbt_inherit/effect/source "Template:Nbt inherit/effect/source")
                -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)probability：（0≤值≤1，默认为1）食用后施加此状态效果的概率。
                
                -   如果![字符串](/images/Data_node_string.svg?42545)type为`clear_all_effects`，则对使用此物品的生物移除所有状态效果。
                
                -   如果![字符串](/images/Data_node_string.svg?42545)type为`play_sound`，则播放指定的声音：
                -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)\*  
                    \*sound：要播放的声音。
                    
                    -   声音事件，见[Template:Nbt inherit/sound event/source](/w/Template:Nbt_inherit/sound_event/source "Template:Nbt inherit/sound event/source")
                
                -   如果![字符串](/images/Data_node_string.svg?42545)type为`remove_effects`，则对使用此物品的生物移除指定状态效果：
                -   ![字符串](/images/Data_node_string.svg?42545)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)\*  
                    \*effects：物品使用后要移除的状态效果。可以为以`#`开头的[状态效果](/w/%E7%8A%B6%E6%80%81%E6%95%88%E6%9E%9C "状态效果")标签、一个状态效果ID、或以多个状态效果ID组成的列表。
                
                -   如果![字符串](/images/Data_node_string.svg?42545)type为`teleport_randomly`，则对使用此物品的生物进行随机传送：
                -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)diameter：（值>0，默认为16）随机传送的半径，以传送前的位置作为原点。
        -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)sound：（默认为`entity.generic.eat`）使用物品时产生的声音。
            
            -   声音事件，见[Template:Nbt inherit/sound event/source](/w/Template:Nbt_inherit/sound_event/source "Template:Nbt inherit/sound event/source")

示例

给予当前实体一个[铁镐](/w/%E9%93%81%E9%95%90 "铁镐")，按下右键会花费1.6秒（32[游戏刻](/w/%E6%B8%B8%E6%88%8F%E5%88%BB "游戏刻")）食用此铁镐，食用时播放声音事件“铁砧：着陆”，食用后获得6000[游戏刻](/w/%E6%B8%B8%E6%88%8F%E5%88%BB "游戏刻")（5分）的不显示效果粒子的16级[急迫](/w/%E6%80%A5%E8%BF%AB "急迫")效果：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s minecraft:iron_pickaxe[minecraft:consumable={animation:"eat",consume_seconds:1.6,on_consume_effects:[{type:"apply_effects",effects:[{amplifier:15,duration:6000,id:"haste",show_icon:true,show_particles:false}]}],sound:"block.anvil.land"}]`

### container

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=23&veaction=edit "编辑章节：container") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=23 "编辑章节的源代码： container")\]

存储容器方块内部物品栏的物品。当此方块的物品实体被摧毁时会释放内容物，提示框中会显示至多五项内容物信息。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)minecraft:container：方块物品的内部物品栏。
        -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)：一个槽位上的物品堆叠数据。
            -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)\*  
                \*item\[失效：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]：此槽位的物品堆叠数据。
                
                -   物品共通标签，见[Template:Nbt inherit/itemnoslot/source](/w/Template:Nbt_inherit/itemnoslot/source "Template:Nbt inherit/itemnoslot/source")
            -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)\*  
                \*item\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]：此槽位的物品堆叠数据。
                
                -   物品模板，见[Template:Nbt inherit/item template/source](/w/Template:Nbt_inherit/item_template/source "Template:Nbt inherit/item template/source")
            -   ![整型](/images/Data_node_int.svg?8d24f)\*  
                \*slot：（0≤值≤255）物品堆叠所在的槽位。

示例

给予一个木桶，其中的第一个槽位放了一个苹果：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s barrel[container=[{slot:0,item:{id:apple}}]]`

### container\_loot

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=24&veaction=edit "编辑章节：container_loot") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=24 "编辑章节的源代码： container_loot")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:container\_loot：战利品容器方块的战利品表数据。
        -   ![字符串](/images/Data_node_string.svg?42545)\*  
            \*loot\_table：（命名空间ID）生成战利品使用的[战利品表](/w/%E6%88%98%E5%88%A9%E5%93%81%E8%A1%A8 "战利品表")。
        -   ![长整型](/images/Data_node_long.svg?dde3f)seed：（默认为0）生成战利品使用的种子，0或不输入将使用[随机序列](/w/%E9%9A%8F%E6%9C%BA%E5%BA%8F%E5%88%97 "随机序列")。

### custom\_data

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=25&veaction=edit "编辑章节：custom_data") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=25 "编辑章节的源代码： custom_data")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:custom\_data：自定义的数据。如果采用字符串格式进行定义，则游戏会将字符串的内容视为[SNBT](/w/SNBT "SNBT")加载，游戏只保存为复合标签格式。
        -   ![任意类型](/images/Data_node_any.svg?d406c)<*自定义标签名*\>：一个可以为任意类型的自定义标签。

### custom\_model\_data

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=26&veaction=edit "编辑章节：custom_model_data") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=26 "编辑章节的源代码： custom_model_data")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:custom\_model\_data：自定义物品模型数据。
        -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)colors：定义物品模型映射中的[着色](/w/%E7%89%A9%E5%93%81%E6%A8%A1%E5%9E%8B%E6%98%A0%E5%B0%84#model "物品模型映射")列表。
            -   ![整型](/images/Data_node_int.svg?8d24f)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)：一个颜色。可以直接使用整数定义颜色，也可以使用RGB三个分量定义颜色，游戏只保存为整数形式。
                
                -   RGB颜色，见[Template:Nbt inherit/rgb color/source](/w/Template:Nbt_inherit/rgb_color/source "Template:Nbt inherit/rgb color/source")
        -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)flags：定义`[condition](/w/%E7%89%A9%E5%93%81%E6%A8%A1%E5%9E%8B%E6%98%A0%E5%B0%84#condition "物品模型映射")`物品模型映射类型的布尔值列表。
            -   ![布尔型](/images/Data_node_bool.svg?77754)：一个布尔值。
        -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)floats：定义`[range_dispatch](/w/%E7%89%A9%E5%93%81%E6%A8%A1%E5%9E%8B%E6%98%A0%E5%B0%84#range_dispatch "物品模型映射")`物品模型映射类型的浮点数列表。
            -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)：一个浮点数。
        -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)strings：定义`[select](/w/%E7%89%A9%E5%93%81%E6%A8%A1%E5%9E%8B%E6%98%A0%E5%B0%84#select "物品模型映射")`物品模型映射类型的字符串列表。
            -   ![字符串](/images/Data_node_string.svg?42545)：一个字符串。

正在加载互动小工具。如果加载失败，请您刷新本页面并检查JavaScript是否已启用。

### custom\_name

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=27&veaction=edit "编辑章节：custom_name") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=27 "编辑章节的源代码： custom_name")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)minecraft:custom\_name：（[文本组件](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6 "文本组件")）自定义名称。

示例

给予一个名为“Magic Wand”的木棍：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s stick[custom_name="Magic Wand"]`

### damage

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=28&veaction=edit "编辑章节：damage") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=28 "编辑章节的源代码： damage")\]

存储物品的[损坏值](/w/%E8%80%90%E4%B9%85%E5%BA%A6 "耐久度")，和`max_damage`一起控制物品能否被损坏。此组件不存在时代表物品处于最大耐久值。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![整型](/images/Data_node_int.svg?8d24f)minecraft:damage：（值≥0）物品的损坏值。

示例

给予一把缺少50点耐久的下界合金镐：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s netherite_pickaxe[damage=50]`

### damage\_resistant

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=29&veaction=edit "编辑章节：damage_resistant") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=29 "编辑章节的源代码： damage_resistant")\]

控制物品免疫的伤害类型。当物品实体受到此类伤害时不会被摧毁，且物品被装备时也不会因为受到此类伤害而消耗耐久度。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:damage\_resistant：物品免疫的伤害类型。
        -   ![字符串](/images/Data_node_string.svg?42545)\*  
            \*types\[失效：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]：（命名空间ID）物品免疫的伤害。应为一个带`#`前缀的标签ID，游戏会将此值解析为[伤害类型标签](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E4%BC%A4%E5%AE%B3%E7%B1%BB%E5%9E%8B "Java版标签/伤害类型")，如果伤害类型标签不存在则物品不会免疫任何伤害。
        -   ![字符串](/images/Data_node_string.svg?42545)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)\*  
            \*types\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]：物品免疫的伤害。可以为一个伤害类型ID、一个伤害类型标签ID，或一个伤害类型标签。

示例

给予玩家一个不会被火焰伤害损坏的铁胸甲：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s minecraft:iron_chestplate[minecraft:damage_resistant={types:"#is_fire"}]`

### damage\_type

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=30&veaction=edit "编辑章节：damage_type") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=30 "编辑章节的源代码： damage_type")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:damage\_type：（命名空间ID）使用此物品攻击时造成的伤害类型。

示例

给予玩家一个会造成箭伤害的钻石剑：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s diamond_sword[damage_type=arrow]`

### death\_protection

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=31&veaction=edit "编辑章节：death_protection") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=31 "编辑章节的源代码： death_protection")\]

控制物品是否具有类似[不死图腾](/w/%E4%B8%8D%E6%AD%BB%E5%9B%BE%E8%85%BE "不死图腾")的行为。当物品在手上时，如果生物受到伤害类型不为`#bypasses_invulnerability`的致死伤害，游戏会阻止生物死亡、将生命值设置为1，并消耗此物品。此消耗行为不属于消耗使用行为。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:death\_protection：持有者将要死亡时阻止生物死亡后的效果。
        -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)death\_effects：（默认为空）触发此物品后产生的效果。
            -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)：一项效果。
                -   ![字符串](/images/Data_node_string.svg?42545)\*  
                    \*type：消耗使用效果类型。
                
                -   如果![字符串](/images/Data_node_string.svg?42545)type为`apply_effects`，则对使用此物品的生物添加状态效果：
                -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)\*  
                    \*effects：物品使用后添加的状态效果。
                    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)：一项状态效果。
                        
                        -   状态效果，见[Template:Nbt inherit/effect/source](/w/Template:Nbt_inherit/effect/source "Template:Nbt inherit/effect/source")
                -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)probability：（0≤值≤1，默认为1）食用后施加此状态效果的概率。
                
                -   如果![字符串](/images/Data_node_string.svg?42545)type为`clear_all_effects`，则对使用此物品的生物移除所有状态效果。
                
                -   如果![字符串](/images/Data_node_string.svg?42545)type为`play_sound`，则播放指定的声音：
                -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)\*  
                    \*sound：要播放的声音。
                    
                    -   声音事件，见[Template:Nbt inherit/sound event/source](/w/Template:Nbt_inherit/sound_event/source "Template:Nbt inherit/sound event/source")
                
                -   如果![字符串](/images/Data_node_string.svg?42545)type为`remove_effects`，则对使用此物品的生物移除指定状态效果：
                -   ![字符串](/images/Data_node_string.svg?42545)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)\*  
                    \*effects：物品使用后要移除的状态效果。可以为以`#`开头的[状态效果](/w/%E7%8A%B6%E6%80%81%E6%95%88%E6%9E%9C "状态效果")标签、一个状态效果ID、或以多个状态效果ID组成的列表。
                
                -   如果![字符串](/images/Data_node_string.svg?42545)type为`teleport_randomly`，则对使用此物品的生物进行随机传送：
                -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)diameter：（值>0，默认为16）随机传送的半径，以传送前的位置作为原点。

### debug\_stick\_state

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=32&veaction=edit "编辑章节：debug_stick_state") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=32 "编辑章节的源代码： debug_stick_state")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:debug\_stick\_state：[调试棒](/w/%E8%B0%83%E8%AF%95%E6%A3%92 "调试棒")的调试数据。
        -   ![字符串](/images/Data_node_string.svg?42545)<*方块命名空间ID*\>：一个方块和此方块将要修改的方块属性的键值对。

### dye

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=33&veaction=edit "编辑章节：dye") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=33 "编辑章节的源代码： dye")\]

[![](/images/thumb/Grass_Block_JE7_BE6.png/16px-Grass_Block_JE7_BE6.png?d27c1)](/w/File:Grass_Block_JE7_BE6.png)

**本段落包含会在[下一次更新](/w/%E8%AE%A1%E5%88%92%E7%89%88%E6%9C%AC "计划版本")中出现的内容。**

这些特性在[Java版26.1](/w/Java%E7%89%8826.1 "Java版26.1")的开发版本中加入。

存储物品的染料颜色。具有`dye`组件是物品可作为染料使用的必要条件，在对应的游戏场景测试成功后将使用此组件的颜色进行计算：

-   在各种配方中染色时，必须是提供染料的原料。（依[配方](/w/%E9%85%8D%E6%96%B9 "配方")定义不同）
-   在[织布机](/w/%E7%BB%87%E5%B8%83%E6%9C%BA "织布机")中染色时，必须在标签`[#loom_dyes](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E7%89%A9%E5%93%81#loom_dyes "Java版标签/物品")`中。
-   给猫的项圈染色时，必须在标签`[#cat_collar_dyes](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E7%89%A9%E5%93%81#cat_collar_dyes "Java版标签/物品")`中。狼的项圈同理，为`[#wolf_collar_dyes](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E7%89%A9%E5%93%81#wolf_collar_dyes "Java版标签/物品")`。
-   给绵羊的羊毛染色时，必须是[染料物品](/w/%E6%9F%93%E6%96%99 "染料")。
-   给告示牌或悬挂式告示牌的文字染色时，必须是染料物品。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:dye：物品的[染料颜色](/w/%E6%9F%93%E6%96%99%E9%A2%9C%E8%89%B2 "染料颜色")数据。取值为`white`、`orange`、`magenta`、`light_blue`、`yellow`、`lime`、`pink`、`gray`、`light_gray`、`cyan`、`purple`、`blue`、`brown`、`green`、`red`或`black`。

### dyed\_color

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=34&veaction=edit "编辑章节：dyed_color") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=34 "编辑章节的源代码： dyed_color")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![整型](/images/Data_node_int.svg?8d24f)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)minecraft:dyed\_color：物品的颜色。只使用后24位，每个颜色通道占用8位，按RGB依次存储。
        
        -   RGB颜色，见[Template:Nbt inherit/rgb color/source](/w/Template:Nbt_inherit/rgb_color/source "Template:Nbt inherit/rgb color/source")

正在加载互动小工具。如果加载失败，请您刷新本页面并检查JavaScript是否已启用。

### enchantable

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=35&veaction=edit "编辑章节：enchantable") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=35 "编辑章节的源代码： enchantable")\]

存储物品的[附魔能力](/w/%E9%99%84%E9%AD%94%E8%83%BD%E5%8A%9B "附魔能力")。该组件不存在时，物品不可在附魔台中附魔。该组件存在时，若物品存在`enchantments`组件且为空，且存在可附加的魔咒，则该物品可在[附魔台](/w/%E9%99%84%E9%AD%94%E5%8F%B0 "附魔台")中附魔。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:enchantable
        -   ![整型](/images/Data_node_int.svg?8d24f)\*  
            \*value：（值≥1）物品的[附魔能力](/w/%E9%99%84%E9%AD%94%E8%83%BD%E5%8A%9B "附魔能力")。

示例

给予一把可以附魔且附魔能力为2的剪刀：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s shears[enchantable={value:2}]`

### enchantment\_glint\_override

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=36&veaction=edit "编辑章节：enchantment_glint_override") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=36 "编辑章节的源代码： enchantment_glint_override")\]

控制物品是否会显示[光效](/w/%E5%85%89%E6%95%88 "光效")，此组件的优先级高于其他任何影响光效的组件和物品自身属性。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![布尔型](/images/Data_node_bool.svg?77754)minecraft:enchantment\_glint\_override：是否显示[光效](/w/%E5%85%89%E6%95%88 "光效")。

### enchantments和stored\_enchantments

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=37&veaction=edit "编辑章节：enchantments和stored_enchantments") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=37 "编辑章节的源代码： enchantments和stored_enchantments")\]

存储物品的[魔咒](/w/%E9%AD%94%E5%92%92 "魔咒")信息。提示框中会显示魔咒及等级。

两个组件的区别在于：`enchantments`组件添加的是“带活性”的魔咒，其上的魔咒可以产生魔咒效果；而`stored_enchantments`组件添加的是“无活性”的魔咒，通常只用于[附魔书](/w/%E9%99%84%E9%AD%94%E4%B9%A6 "附魔书")存储魔咒，其上的魔咒不会产生效果。

需要注意的是，“无活性”魔咒在发挥实际附魔作用时会受到生存模式所能获取的对应附魔书的最大附魔等级的限制，例如给予自己一本[锋利](/w/%E9%94%8B%E5%88%A9 "锋利")VI的附魔书，而其发挥的附魔作用只能为锋利V。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:enchantments：物品的魔咒数据。
        -   ![整型](/images/Data_node_int.svg?8d24f)<*魔咒命名空间ID*\>：（1≤值≤255）一个[魔咒](/w/%E9%AD%94%E5%92%92 "魔咒")和对应魔咒的等级。
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:stored\_enchantments：[附魔书](/w/%E9%99%84%E9%AD%94%E4%B9%A6 "附魔书")保存的魔咒。
        -   ![整型](/images/Data_node_int.svg?8d24f)<*魔咒命名空间ID*\>：（1≤值≤255）一个[魔咒](/w/%E9%AD%94%E5%92%92 "魔咒")和对应魔咒的等级。

示例

给予一把带有锋利III和击退II的木剑：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s wooden_sword[enchantments={sharpness:3,knockback:2}]`

给予一本带有效率V和耐久III的附魔书： `/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s enchanted_book[stored_enchantments={efficiency:5,unbreaking:3}]`

### entity\_data

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=38&veaction=edit "编辑章节：entity_data") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=38 "编辑章节的源代码： entity_data")\]

存储物品生成对应实体时（如使用刷怪蛋或放置盔甲架）应用于所生成实体的数据。应用时采取合并的方式。若指定的实体类型无法在和平模式下存在，则提示框会提示“已在和平难度下禁用”。

若刷怪蛋将生成的实体为指定了附加数据的[下落的方块](/w/%E4%B8%8B%E8%90%BD%E7%9A%84%E6%96%B9%E5%9D%97 "下落的方块")、[命令方块矿车](/w/%E5%91%BD%E4%BB%A4%E6%96%B9%E5%9D%97%E7%9F%BF%E8%BD%A6 "命令方块矿车")或[刷怪笼矿车](/w/%E5%88%B7%E6%80%AA%E7%AC%BC%E7%9F%BF%E8%BD%A6 "刷怪笼矿车")，则非管理员玩家使用这些物品时不会设置实体数据，且提示框中会显示安全警告。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:entity\_data：物品放出实体时套用到实体上的数据。如果采用字符串格式进行定义，则游戏会将字符串的内容视为[SNBT](/w/SNBT "SNBT")加载，游戏只保存为复合标签格式。
        -   ![字符串](/images/Data_node_string.svg?42545)\*  
            \*id：（命名空间ID）实体类型。
        -   若干与该实体对应的实体数据标签，见[实体数据格式](/w/%E5%AE%9E%E4%BD%93%E6%95%B0%E6%8D%AE%E6%A0%BC%E5%BC%8F "实体数据格式")。

示例

给予一个在放置时成为小型盔甲架的盔甲架：`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s armor_stand[entity_data={id:"armor_stand",Small:1b}]`

### equippable

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=39&veaction=edit "编辑章节：equippable") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=39 "编辑章节的源代码： equippable")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:equippable：物品被穿戴的行为。
        -   ![字符串](/images/Data_node_string.svg?42545)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)allowed\_entities：（默认为全部生物）可以穿戴此物品的生物。可以为以`#`开头的实体类型标签、一个实体类型ID、或以多个实体类型ID组成的字符串列表。
        -   ![字符串](/images/Data_node_string.svg?42545)asset\_id：（命名空间ID）物品被穿戴时的[装备模型](/w/%E8%A3%85%E5%A4%87%E6%A8%A1%E5%9E%8B "装备模型")。此值不存在时，若装备在头部则根据物品模型渲染物品，否则什么也不会渲染。
        -   ![字符串](/images/Data_node_string.svg?42545)camera\_overlay：（命名空间ID）当此项存在且物品被玩家穿戴时，玩家第一人称视角将渲染指定的纹理遮罩。此遮罩可以使用多个设置此标签的物品互相叠加，每个物品指定的遮罩都会被渲染，且渲染顺序按照主手、副手、头盔、胸甲、护腿、靴子、身体、鞍的顺序依次叠加渲染。当遮罩纹理渲染时，遮罩纹理被视为**独立纹理**，即无法作为动态纹理或GUI纹理渲染，但可以指定纹理过滤方式。
        -   ![布尔型](/images/Data_node_bool.svg?77754)can\_be\_sheared：（默认为`false`）满足未被骑乘等其他条件时，玩家是否可以对装备此物品的生物进行修剪来卸下此物品。
        -   ![布尔型](/images/Data_node_bool.svg?77754)damage\_on\_hurt：（默认为`true`）生物在受到会影响损害盔甲的伤害时此物品是否会受损而减少耐久。
        -   ![布尔型](/images/Data_node_bool.svg?77754)equip\_on\_interact：（默认为`false`）对生物使用此物品时，是否可以让被交互的生物在允许的空槽位上穿戴此物品。
        -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)equip\_sound：（默认为`item.armor.equip_generic`，默认可装备鞍的生物的鞍除外）物品被穿戴时的声音。
            
            -   声音事件，见[Template:Nbt inherit/sound event/source](/w/Template:Nbt_inherit/sound_event/source "Template:Nbt inherit/sound event/source")
        -   ![布尔型](/images/Data_node_bool.svg?77754)dispensable：（默认为`true`）是否可以使用[发射器](/w/%E5%8F%91%E5%B0%84%E5%99%A8 "发射器")使生物穿戴此物品。如果物品本身有特殊的发射器行为则此项无效。
        -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)shearing\_sound：（默认为`item.shears.snip`）被玩家使用剪刀卸下此物品时播放的声音。
            
            -   声音事件，见[Template:Nbt inherit/sound event/source](/w/Template:Nbt_inherit/sound_event/source "Template:Nbt inherit/sound event/source")
        -   ![字符串](/images/Data_node_string.svg?42545)\*  
            \*slot：物品可被穿戴的[装备槽位](/w/%E8%A3%85%E5%A4%87%E6%A7%BD%E4%BD%8D "装备槽位")。
        -   ![布尔型](/images/Data_node_bool.svg?77754)swappable：（默认为`true`）物品是否可以直接使用穿戴。

### firework\_explosion

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=40&veaction=edit "编辑章节：firework_explosion") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=40 "编辑章节的源代码： firework_explosion")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:firework\_explosion：[烟火之星](/w/%E7%83%9F%E7%81%AB%E4%B9%8B%E6%98%9F "烟火之星")的数据。
        -   ![整型数组](/images/Data_node_int-array.svg?546e8)colors：（默认为空数组）表示爆裂时的粒子颜色，只使用后24位，每个颜色通道占用8位，按RGB依次存储。如果颜色没有对应的染料颜色，游戏将在[提示框](/w/%E6%8F%90%E7%A4%BA%E6%A1%86 "提示框")中显示为“自定义”，但爆裂时会产生正确的颜色。当存在多个值时，每个爆裂粒子在渲染时会随机选择一种颜色用于渲染。不存在或数组为空时被视为黑色。
        -   ![整型数组](/images/Data_node_int-array.svg?546e8)fade\_colors：（默认为空数组）表示爆裂后的淡化粒子颜色，只使用后24位，每个颜色通道占用8位，按RGB依次存储。当存在多个值时，每个爆裂粒子在渲染时会随机选择一种颜色用于渲染。
        -   ![布尔型](/images/Data_node_bool.svg?77754)has\_trail：（默认为`false`）表示烟火是否有拖曳痕迹（使用[钻石](/w/%E9%92%BB%E7%9F%B3 "钻石")合成时）。
        -   ![布尔型](/images/Data_node_bool.svg?77754)has\_twinkle：（默认为`false`）表示烟火是否出现闪烁效果（使用[荧石粉](/w/%E8%8D%A7%E7%9F%B3%E7%B2%89 "荧石粉")合成时）。
        -   ![字符串](/images/Data_node_string.svg?42545)\*  
            \*shape：爆裂时的形态。可以为`small_ball`（小型球状）、`large_ball`（大型球状）、`star`（星形）、`creeper`（苦力怕状）、`burst`（喷发状）。

### fireworks

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=41&veaction=edit "编辑章节：fireworks") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=41 "编辑章节的源代码： fireworks")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:fireworks：[烟花火箭](/w/%E7%83%9F%E8%8A%B1%E7%81%AB%E7%AE%AD "烟花火箭")的数据。
        -   ![字节型](/images/Data_node_byte.svg?eb0e0)flight\_duration：（无符号8位整数，默认为0）烟花火箭的飞行的时间，单位为“火药”（即表现为和在工作台上合成烟花火箭时所用的火药数相等）。
        -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)explosions：（最多256个元素）烟花火箭对应的烟火之星的数据，控制烟花火箭飞行结束时产生的爆裂烟花渲染。
            -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)：一个烟火之星的数据。
                -   ![整型数组](/images/Data_node_int-array.svg?546e8)colors：（默认为空数组）表示爆裂时的粒子颜色，只使用后24位，每个颜色通道占用8位，按RGB依次存储。如果颜色没有对应的染料颜色，游戏将在[提示框](/w/%E6%8F%90%E7%A4%BA%E6%A1%86 "提示框")中显示为“自定义”，但爆裂时会产生正确的颜色。当存在多个值时，每个爆裂粒子在渲染时会随机选择一种颜色用于渲染。不存在或数组为空时被视为黑色。
                -   ![整型数组](/images/Data_node_int-array.svg?546e8)fade\_colors：（默认为空数组）表示爆裂后的淡化粒子颜色，只使用后24位，每个颜色通道占用8位，按RGB依次存储。当存在多个值时，每个爆裂粒子在渲染时会随机选择一种颜色用于渲染。
                -   ![布尔型](/images/Data_node_bool.svg?77754)has\_trail：（默认为`false`）表示烟火是否有拖曳痕迹（使用[钻石](/w/%E9%92%BB%E7%9F%B3 "钻石")合成时）。
                -   ![布尔型](/images/Data_node_bool.svg?77754)has\_twinkle：（默认为`false`）表示烟火是否出现闪烁效果（使用[荧石粉](/w/%E8%8D%A7%E7%9F%B3%E7%B2%89 "荧石粉")合成时）。
                -   ![字符串](/images/Data_node_string.svg?42545)\*  
                    \*shape：爆裂时的形态。可以为`small_ball`（小型球状）、`large_ball`（大型球状）、`star`（星形）、`creeper`（苦力怕状）、`burst`（喷发状）。

### food

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=42&veaction=edit "编辑章节：food") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=42 "编辑章节的源代码： food")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:food：物品的食物属性。
        -   ![布尔型](/images/Data_node_bool.svg?77754)can\_always\_eat：（默认为`false`）表示物品是否可以无视当前饥饿值食用。
        -   ![整型](/images/Data_node_int.svg?8d24f)\*  
            \*nutrition：（值≥0）食用物品时增加的饥饿值。
        -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)\*  
            \*saturation：食用物品时增加的饱和度。

示例

给予一个海绵，该海绵可无视饥饿值食用，玩家食用后恢复玩家3点饥饿值和1点饱和度：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s sponge[food={can_always_eat:true,nutrition:3,saturation:1},consumable={}]`

### glider

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=43&veaction=edit "编辑章节：glider") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=43 "编辑章节的源代码： glider")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:glider：空标签，此组件存在时若被生物装备则可以[滑翔](/w/%E6%BB%91%E7%BF%94 "滑翔")，且滑翔时此物品每1秒消耗1耐久度。

示例

给予一个铁胸甲，穿戴后可以滑翔：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s iron_chestplate[glider={}]`

### instrument

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=44&veaction=edit "编辑章节：instrument") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=44 "编辑章节的源代码： instrument")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:instrument：玩家吹奏[山羊角](/w/%E5%B1%B1%E7%BE%8A%E8%A7%92 "山羊角")时使用的[山羊角乐器](/w/%E5%B1%B1%E7%BE%8A%E8%A7%92%E4%B9%90%E5%99%A8 "山羊角乐器")。命名空间ID或内联定义均可。
        -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)\*  
            \*description：（[文本组件](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6 "文本组件")）乐器的名称。
        -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)\*  
            \*sound\_event：吹奏时播放的声音。
            
            -   声音事件，见[Template:Nbt inherit/sound event/source](/w/Template:Nbt_inherit/sound_event/source "Template:Nbt inherit/sound event/source")
        -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)\*  
            \*range：（值>0）吹奏声音能传播的最远距离。
        -   ![整型](/images/Data_node_int.svg?8d24f)\*  
            \*use\_duration：（值>0）吹奏时间，影响物品冷却速度。

### intangible\_projectile

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=45&veaction=edit "编辑章节：intangible_projectile") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=45 "编辑章节的源代码： intangible_projectile")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:intangible\_projectile：空标签，此组件存在时若作为箭射出，则射出后只能被创造模式玩家捡起。

### item\_model

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=46&veaction=edit "编辑章节：item_model") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=46 "编辑章节的源代码： item_model")\]

控制物品的物品模型映射。物品模型映射会根据命名空间ID解析为`assets/<*命名空间*>/items/<*路径*>.json`。若对应的物品模型映射不存在或无法解析则使用[无效模型](/w/%E6%97%A0%E6%95%88%E6%A8%A1%E5%9E%8B "无效模型")。此组件不存在时什么也不会渲染。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:item\_model：（[命名空间ID](/w/%E5%91%BD%E5%90%8D%E7%A9%BA%E9%97%B4ID "命名空间ID")）为当前物品绑定一个[物品模型映射](/w/%E7%89%A9%E5%93%81%E6%A8%A1%E5%9E%8B%E6%98%A0%E5%B0%84 "物品模型映射")。

### item\_name

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=47&veaction=edit "编辑章节：item_name") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=47 "编辑章节的源代码： item_name")\]

控制物品的默认名称。该名称无法通过铁砧修改，不能在物品展示框中显示名称，带有该组件的旗帜在充当[地图标记](/w/%E5%9C%B0%E5%9B%BE#地图标记 "地图")时也不会显示名称。此组件对物品名称的控制等级永远最低，会被其他所有影响物品名称的组件覆盖。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)minecraft:item\_name：（[文本组件](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6 "文本组件")）物品的默认名称。

### jukebox\_playable

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=48&veaction=edit "编辑章节：jukebox_playable") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=48 "编辑章节的源代码： jukebox_playable")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:jukebox\_playable：（命名空间ID）[唱片机曲目](/w/%E5%94%B1%E7%89%87%E6%9C%BA%E6%9B%B2%E7%9B%AE "唱片机曲目")。此组件存在时物品可插进唱片机中播放。

示例

给予一个可以放进唱片机并播放cat的唱片残片：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @a minecraft:disc_fragment_5[minecraft:jukebox_playable=cat]`

### kinetic\_weapon

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=49&veaction=edit "编辑章节：kinetic_weapon") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=49 "编辑章节的源代码： kinetic_weapon")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:kinetic\_weapon：设置物品的冲锋攻击。
        -   ![整型](/images/Data_node_int.svg?8d24f)delay\_ticks：（值≥0，默认为0）武器生效前的时间，单位为游戏刻。
        -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)forward\_movement：（默认为0）动画期间脱离手的距离。
        -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)damage\_multiplier：（默认为1）攻击轴相对速度的最终伤害倍率。此处及下文的“攻击轴速度”定义为：上个游戏刻的位移向量（对于玩家）或速度改变量（对于非生物实体，等于![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)Motion）对攻击者视角向量的投影，是一个向量；“攻击轴相对速度”即攻击者攻击轴速度与被攻击者攻击轴速度之差，如果此差值小于0则游戏认为是0。
        -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)sound：使用此武器时播放的声音。
            
            -   声音事件，见[Template:Nbt inherit/sound event/source](/w/Template:Nbt_inherit/sound_event/source "Template:Nbt inherit/sound event/source")
        -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)hit\_sound：此武器攻击到生物时播放的声音。
            
            -   声音事件，见[Template:Nbt inherit/sound event/source](/w/Template:Nbt_inherit/sound_event/source "Template:Nbt inherit/sound event/source")
        -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)\*  
            \*dismount\_conditions：将目标强制脱离骑乘的条件。
            -   ![整型](/images/Data_node_int.svg?8d24f)\*  
                \*max\_duration\_ticks：不再检查条件的时间，单位为刻，从![整型](/images/Data_node_int.svg?8d24f)delay\_ticks开始计算。
            -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)min\_speed：（默认为0）攻击者的最低攻击轴速度。对于非玩家实体，实际最小速度为规定值的20%。
            -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)min\_relative\_speed：（默认为0）最小攻击轴相对速度。对于非玩家实体，实际最小速度为规定值的20%。
        -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)\*  
            \*knockback\_conditions：将目标击退的条件。
            -   格式同![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)dismount\_conditions。
        -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)\*  
            \*damage\_conditions：对目标造成伤害的条件。
            -   格式同![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)dismount\_conditions。
        -   ![整型](/images/Data_node_int.svg?8d24f)contact\_cooldown\_ticks：（值>0，默认为10）攻击的冷却时间，在此时间内无法与任何实体交互。

示例

给予一个可以蓄力攻击的下界合金剑，无法将目标强制脱离骑乘：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s minecraft:netherite_sword[minecraft:kinetic_weapon={dismount_conditions:{max_duration_ticks:0},knockback_conditions:{max_duration_ticks:2147483647},damage_conditions:{max_duration_ticks:2147483647}}]`

### lock

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=50&veaction=edit "编辑章节：lock") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=50 "编辑章节的源代码： lock")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:lock：容器方块的上锁数据。
        
        -   物品堆叠谓词，见[Template:Nbt inherit/item predicate/source](/w/Template:Nbt_inherit/item_predicate/source "Template:Nbt inherit/item predicate/source")

示例

给予一个箱子，此箱子放置后玩家仅能手持名为“密码”的物品来打开它：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s chest[lock={components:{custom_name:"密码"}}]`

### lodestone\_tracker

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=51&veaction=edit "编辑章节：lodestone_tracker") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=51 "编辑章节的源代码： lodestone_tracker")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:lodestone\_tracker：若指南针拥有此组件，则指南针将变为[磁石指针](/w/%E7%A3%81%E7%9F%B3%E6%8C%87%E9%92%88 "磁石指针")。
        -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)target：磁石指针指向的位置。
            -   ![字符串](/images/Data_node_string.svg?42545)\*  
                \*dimension：磁石指针指向位置的所在维度。
            -   ![整型数组](/images/Data_node_int-array.svg?546e8)\*  
                \*pos：磁石指针指向的坐标。内部的三个整数分别代表了位置的XYZ坐标值。
        -   ![布尔型](/images/Data_node_bool.svg?77754)tracked：（默认为`true`）表示磁石指针是否追踪绑定的磁石。为false时，当磁石被破坏后此组件不会被移除，磁石指针仍然指向对应位置。

示例

给予一个指南针，其始终指向主世界X=1、Y=2、Z=3处，无论磁石是否存在：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s compass[lodestone_tracker={target:{pos:[I;1,2,3],dimension:"overworld"},tracked:false}]`

### lore

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=52&veaction=edit "编辑章节：lore") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=52 "编辑章节的源代码： lore")\]

![](/images/Disambig_gray.svg?1bb41)“**Lore**”重定向至此。关于1.20.5前的Lore标签，请见“**[物品格式/Java版1.20.5前 § 通用标签](/w/%E7%89%A9%E5%93%81%E6%A0%BC%E5%BC%8F/Java%E7%89%881.20.5%E5%89%8D#通用标签 "物品格式/Java版1.20.5前")**”。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)minecraft:lore：物品的自定义描述信息，共计不允许超过256行。
        -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)：（[文本组件](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6 "文本组件")）一行描述信息。

示例

给予一根木棍，描述信息为第一行“Hello Minecraft”、第二行“Hello World”的木棍：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s stick[minecraft:lore=["Hello Minecraft", "Hello World"]]`

### map\_color

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=53&veaction=edit "编辑章节：map_color") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=53 "编辑章节的源代码： map_color")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![整型](/images/Data_node_int.svg?8d24f)minecraft:map\_color：（默认为4603950）物品栏内地图纹理上的颜色，在二进制形式下，只使用后24位，每个颜色通道占用8位，按RGB依次存储。

正在加载互动小工具。如果加载失败，请您刷新本页面并检查JavaScript是否已启用。

### map\_decorations

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=54&veaction=edit "编辑章节：map_decorations") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=54 "编辑章节的源代码： map_decorations")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:map\_decorations：地图图标数据。
        -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)<*图标名称*\>：一个图标的信息。
            -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)\*  
                \*rotation：图标的旋转角度，按顺时针角度计。游戏并不能真正显示所有角度，每经过22.5°，在地图上才会有区别。与图标纹理中的外观相比，旋转角度为0所显示的图标上下颠倒。
            -   ![字符串](/images/Data_node_string.svg?42545)\*  
                \*type：（命名空间ID）此图标显示的[地图图标类型](/w/%E5%9C%B0%E5%9B%BE#地图图标 "地图")。
            -   ![双精度浮点数](/images/Data_node_double.svg?14320)\*  
                \*x：图标在世界上所在的X坐标。如果超出地图所展示的范围且图标类型不是玩家，则图标无法添加到地图中。如果图标类型是玩家，位置超出显示范围但地图可以无限追踪玩家，那么图标类型会被修改为`player_off_limits`，且位置会显示在对应边；如果距离显示范围较近，则图标类型会被修改为`player_off_map`，且位置会显示在对应边；如果距离显示范围很远，则移除此图标。
            -   ![双精度浮点数](/images/Data_node_double.svg?14320)\*  
                \*z：图标在世界上所在的Z坐标。如果超出地图所展示的范围且图标类型不是玩家，则图标无法添加到地图中。如果图标类型是玩家，位置超出显示范围但地图可以无限追踪玩家，那么图标类型会被修改为`player_off_limits`，且位置会显示在对应边；如果距离显示范围较近，则图标类型会被修改为`player_off_map`，且位置会显示在对应边；如果距离显示范围很远，则移除此图标。

### map\_id

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=55&veaction=edit "编辑章节：map_id") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=55 "编辑章节的源代码： map_id")\]

存储地图编号。具有此组件的所有物品均会尝试读取相应编号的地图内容，且能被玩家展开，在物品展示框上铺开，作为地图被复制、锁定或扩展。提示框中会显示地图的缩放信息等数据。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![整型](/images/Data_node_int.svg?8d24f)minecraft:map\_id：地图编号。

### max\_damage

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=56&veaction=edit "编辑章节：max_damage") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=56 "编辑章节的源代码： max_damage")\]

存储物品的最大耐久度，和`damage`组件一起控制物品能否被损坏。此组件不存在时若物品被损坏则游戏将最大耐久度视为0。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![整型](/images/Data_node_int.svg?8d24f)minecraft:max\_damage：（值>0）物品的最大耐久度。

示例

给予一个耐久上限999点的金剑：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s golden_sword[max_damage=999]`

### max\_stack\_size

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=57&veaction=edit "编辑章节：max_stack_size") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=57 "编辑章节的源代码： max_stack_size")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![整型](/images/Data_node_int.svg?8d24f)minecraft:max\_stack\_size：（0≤值≤99）物品的最大堆叠数量。如果此组件不存在，则游戏默认为1。

示例

给予99个最大堆叠99个的雪球：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s snowball[max_stack_size=99] 99`

### minimum\_attack\_charge

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=58&veaction=edit "编辑章节：minimum_attack_charge") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=58 "编辑章节的源代码： minimum_attack_charge")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)minecraft:minimum\_attack\_charge：（0≤值≤1）玩家使用此物品进行近战攻击或穿刺攻击所需要[攻击冷却完成度](/w/%E8%BF%91%E6%88%98%E6%94%BB%E5%87%BB#攻击冷却 "近战攻击")的最小值。若添加了该组件，并且值大于0，则会影响[魔咒效果组件](/w/%E9%AD%94%E5%92%92%E6%95%88%E6%9E%9C%E7%BB%84%E4%BB%B6 "魔咒效果组件")`post_piercing_attack`的触发间隔。

### note\_block\_sound

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=59&veaction=edit "编辑章节：note_block_sound") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=59 "编辑章节的源代码： note_block_sound")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:note\_block\_sound：[玩家的头](/w/%E7%8E%A9%E5%AE%B6%E7%9A%84%E5%A4%B4 "玩家的头")被放置在音符盒上时播放的声音。应为一个来自资源包`[sounds.json](/w/Sounds.json "Sounds.json")`内定义的声音事件。

示例

给予一个被放在音符盒上时会播放声音“玩家：升级”的玩家的头：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s player_head[note_block_sound="entity.player.levelup"]`

### ominous\_bottle\_amplifier

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=60&veaction=edit "编辑章节：ominous_bottle_amplifier") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=60 "编辑章节的源代码： ominous_bottle_amplifier")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![整型](/images/Data_node_int.svg?8d24f)minecraft:ominous\_bottle\_amplifier：（0≤值≤4）玩家使用物品后获得的[不祥之兆](/w/%E4%B8%8D%E7%A5%A5%E4%B9%8B%E5%85%86 "不祥之兆")状态效果倍率。

### piercing\_weapon

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=61&veaction=edit "编辑章节：piercing_weapon") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=61 "编辑章节的源代码： piercing_weapon")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:piercing\_weapon：设置物品的戳刺攻击，跳过玩家对方块的点击和持续破坏行为。该组件也是触发[魔咒效果组件](/w/%E9%AD%94%E5%92%92%E6%95%88%E6%9E%9C%E7%BB%84%E4%BB%B6 "魔咒效果组件")`post_piercing_attack`的条件之一。
        -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)sound：使用此武器时播放的声音。
            
            -   声音事件，见[Template:Nbt inherit/sound event/source](/w/Template:Nbt_inherit/sound_event/source "Template:Nbt inherit/sound event/source")
        -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)hit\_sound：此武器攻击到生物时播放的声音。
            
            -   声音事件，见[Template:Nbt inherit/sound event/source](/w/Template:Nbt_inherit/sound_event/source "Template:Nbt inherit/sound event/source")
        -   ![布尔型](/images/Data_node_bool.svg?77754)deals\_knockback：（默认为`true`）攻击是否造成击退。
        -   ![布尔型](/images/Data_node_bool.svg?77754)dismounts：（默认为`false`）攻击是否将目标强制脱离骑乘。

### pot\_decorations

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=62&veaction=edit "编辑章节：pot_decorations") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=62 "编辑章节的源代码： pot_decorations")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)minecraft:pot\_decorations：[饰纹陶罐](/w/%E9%A5%B0%E7%BA%B9%E9%99%B6%E7%BD%90 "饰纹陶罐")的陶片数据。此列表应仅有四个元素，依次代表饰纹陶罐背面、左面、右面和前面的物品。默认每个面均为红砖。
        -   ![字符串](/images/Data_node_string.svg?42545)：（命名空间ID）饰纹陶罐这一个面的陶片物品。

### potion\_contents

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=63&veaction=edit "编辑章节：potion_contents") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=63 "编辑章节的源代码： potion_contents")\]

存储[药水效果](/w/%E8%8D%AF%E6%B0%B4%E6%95%88%E6%9E%9C "药水效果")和[状态效果](/w/%E7%8A%B6%E6%80%81%E6%95%88%E6%9E%9C "状态效果")信息。影响物品的名称和状态效果，提示框中会显示药水信息。下文的`<*药水物品类型*>`只对药水、喷溅药水、滞留药水和药箭有效。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:potion\_contents：物品的药水和自定义状态效果数据。如果设置此组件为字符串，则等价于只设置复合标签形式中的![字符串](/images/Data_node_string.svg?42545)potion，游戏在保存时只会保存为复合标签形式。
        -   ![整型](/images/Data_node_int.svg?8d24f)custom\_color：物品渲染中，药水部分使用的颜色。只使用后24位，每个颜色通道占用8位，按RGB依次存储。
        -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)custom\_effects：当前物品所含有的自定义[状态效果](/w/%E7%8A%B6%E6%80%81%E6%95%88%E6%9E%9C "状态效果")。
            -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)：一项状态效果。
                
                -   状态效果，见[Template:Nbt inherit/effect/source](/w/Template:Nbt_inherit/effect/source "Template:Nbt inherit/effect/source")
        -   ![字符串](/images/Data_node_string.svg?42545)custom\_name：覆盖物品的默认名称，游戏将以`<*药水物品名称翻译键*>.effect.<*此值*>`翻译键作为物品的名称，对于原版的药水物品而言就是`item.minecraft.<*药水物品类型*>.effect.<*此值*>`。
        -   ![字符串](/images/Data_node_string.svg?42545)potion：（命名空间ID）[药水效果](/w/%E8%8D%AF%E6%B0%B4%E6%95%88%E6%9E%9C "药水效果")，也会影响物品的名称和纹理。

正在加载互动小工具。如果加载失败，请您刷新本页面并检查JavaScript是否已启用。

示例

给予一瓶药水，药水颜色为紫色，药水的状态效果为持续1102游戏刻（55.1秒）的122倍率（即123级）的[发光](/w/%E5%8F%91%E5%85%89 "发光")效果：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s potion[potion_contents={custom_color:8388863,custom_effects:[{amplifier:122,duration:1102,id:"glowing"}]}]`

### potion\_duration\_scale

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=64&veaction=edit "编辑章节：potion_duration_scale") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=64 "编辑章节的源代码： potion_duration_scale")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)minecraft:potion\_duration\_scale：（值≥0）控制`potion_contents`组件存储的状态效果时长缩放倍率。此组件不存在时默认为1。

### profile

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=65&veaction=edit "编辑章节：profile") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=65 "编辑章节的源代码： profile")\]

存储玩家的头的对应玩家的游戏档案数据。

游戏会优先使用玩家档案获取皮肤等数据，然后再根据玩家皮肤指定的纹理模型等进行更改。例如游戏渲染玩家的头`player_head[profile={name:"jeb_", texture:"missingno"}]`时，会先获取玩家jeb\_的皮肤，再使用无效纹理进行覆盖。

如果设置了![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)properties，则游戏直接使用此游戏档案数据，不会因为对应玩家档案的更改而更改。

在未设置![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)properties的条件下，如果设置了![字符串](/images/Data_node_string.svg?42545)name，则游戏会将其视作玩家名称解析游戏档案。如果设置了![整型数组](/images/Data_node_int-array.svg?546e8)id，则游戏会将其视作UUID解析游戏档案，解析的优先级高于![字符串](/images/Data_node_string.svg?42545)name。游戏并不会将获取的游戏档案数据存储，而是实时获取，尽管需要客户端重新启动才能更改渲染效果。此时物品提示框也会显示“实时显示”，以与静态游戏档案相区分。

无论是静态档案还是动态档案，只有![字符串](/images/Data_node_string.svg?42545)name会影响玩家的头的物品名称。

如果设置此组件为字符串，则等价于只设置复合标签形式中的![字符串](/images/Data_node_string.svg?42545)name，游戏在保存时只会保存为复合标签形式。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:profile：玩家档案数据。
        
        -   游戏档案，见[Template:Nbt inherit/resolvable profile/source](/w/Template:Nbt_inherit/resolvable_profile/source "Template:Nbt inherit/resolvable profile/source")

正在加载互动小工具。如果加载失败，请您刷新本页面并检查JavaScript是否已启用。

游戏档案属性通常包括`textures`用于保存玩家的皮肤数据。在此属性的数据被Base64解码后具有如下结构：

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597) JSON数据根元素
    -   ![字符串](/images/Data_node_string.svg?42545)\*  
        \*profileId：游戏档案的UUID，不带连字符。
    -   ![字符串](/images/Data_node_string.svg?42545)\*  
        \*profileName：游戏档案名称。
    -   ![布尔型](/images/Data_node_bool.svg?77754)signatureRequired：代表此纹理属性是否已被签名。如果![字符串](/images/Data_node_string.svg?42545)signature存在，则此项也存在并为true。
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)\*  
        \*textures：纹理数据。
        -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)CAPE：[披风](/w/%E6%8A%AB%E9%A3%8E "披风")纹理。如果此游戏档案不包含披风，此项不存在。
            -   ![字符串](/images/Data_node_string.svg?42545)\*  
                \*url：披风纹理的URL链接。
        -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)SKIN：[皮肤](/w/%E7%9A%AE%E8%82%A4 "皮肤")纹理。如果此游戏档案不包含自定义皮肤，此项不存在。
            -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)metadata：皮肤的元数据。
                -   ![字符串](/images/Data_node_string.svg?42545)model：固定值`slim`。当皮肤模型手臂为3像素时存在，否则不存在。
            -   ![字符串](/images/Data_node_string.svg?42545)\*  
                \*url：皮肤纹理的URL链接。
    -   ![整型](/images/Data_node_int.svg?8d24f)\*  
        \*timestamp：[Unix时间戳](https://en.wikipedia.org/wiki/Unixtime "wikipedia:Unixtime")，以毫秒为单位，时间为请求玩家游戏档案数据的时间。

### provides\_banner\_patterns

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=66&veaction=edit "编辑章节：provides_banner_patterns") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=66 "编辑章节的源代码： provides_banner_patterns")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:provides\_banner\_patterns\[失效：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]：（命名空间ID）控制物品能否放进[织布机](/w/%E7%BB%87%E5%B8%83%E6%9C%BA "织布机")的旗帜图案槽位，以及可以制作的图案。应为一个带`#`前缀的标签ID，游戏会将此值解析为[旗帜图案标签](/w/Java%E7%89%88%E6%A0%87%E7%AD%BE/%E6%97%97%E5%B8%9C%E5%9B%BE%E6%A1%88 "Java版标签/旗帜图案")，如果旗帜图案标签不存在则织布机不会显示任何配方。
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:provides\_banner\_patterns\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]：控制物品能否放进[织布机](/w/%E7%BB%87%E5%B8%83%E6%9C%BA "织布机")的旗帜图案槽位，以及可以制作的图案。可以为一个旗帜图案ID、一个旗帜图案标签，或一个旗帜图案ID的列表。

### provides\_trim\_material

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=67&veaction=edit "编辑章节：provides_trim_material") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=67 "编辑章节的源代码： provides_trim_material")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:provides\_trim\_material：控制物品在锻造台上使用盔甲纹饰配方时为输出物品提供的盔甲纹饰材料。命名空间ID或内联定义均可。
        -   ![字符串](/images/Data_node_string.svg?42545)\*  
            \*asset\_name：一个字符串，实际盔甲纹饰纹理的后缀。
        -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)\*  
            \*description：（[文本组件](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6 "文本组件")）盔甲纹饰材料的名称。
        -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)override\_armor\_assets：对于指定的[装备模型](/w/%E8%A3%85%E5%A4%87%E6%A8%A1%E5%9E%8B "装备模型")，使用指定纹理覆盖而不使用![字符串](/images/Data_node_string.svg?42545)asset\_name。
            -   ![字符串](/images/Data_node_string.svg?42545)<*装备模型命名空间ID*\>：一个字符串，实际盔甲纹饰纹理的后缀。

### rarity

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=68&veaction=edit "编辑章节：rarity") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=68 "编辑章节的源代码： rarity")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:rarity：物品的基础[稀有度](/w/%E7%A8%80%E6%9C%89%E5%BA%A6 "稀有度")。可以为`common`（常见）、`uncommon`（少见）、`rare`（稀有）、`epic`（史诗）。

### recipes

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=69&veaction=edit "编辑章节：recipes") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=69 "编辑章节的源代码： recipes")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)minecraft:recipes：[知识之书](/w/%E7%9F%A5%E8%AF%86%E4%B9%8B%E4%B9%A6 "知识之书")保存的配方数据。
        -   ![字符串](/images/Data_node_string.svg?42545)：（命名空间ID）一个配方ID。

### repair\_cost

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=70&veaction=edit "编辑章节：repair_cost") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=70 "编辑章节的源代码： repair_cost")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![整型](/images/Data_node_int.svg?8d24f)minecraft:repair\_cost：（值≥0）物品在[铁砧](/w/%E9%93%81%E7%A0%A7 "铁砧")上修理、合并或重命名时在基础经验等级消耗之上额外增加的累积惩罚。

示例

给予当前实体一个累计惩罚值为30的下界合金剑：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s netherite_sword[repair_cost=30]`

### repairable

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=71&veaction=edit "编辑章节：repairable") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=71 "编辑章节的源代码： repairable")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:repairable：物品被铁砧进行原材料修复的有效物品。不论此值为何，物品永远可以被合并物品修复。
        -   ![字符串](/images/Data_node_string.svg?42545)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)\*  
            \*items：可用于修复的物品。可以为一个`#`开头的物品标签、一个物品ID、或一个物品ID的列表。

示例

给予当前实体一个只能被[橡木木板](/w/%E6%A9%A1%E6%9C%A8%E6%9C%A8%E6%9D%BF "橡木木板")或[下界合金剑](/w/%E4%B8%8B%E7%95%8C%E5%90%88%E9%87%91%E5%89%91 "下界合金剑")修复的下界合金剑：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s netherite_sword[repairable={items:"oak_planks"}]`

### suspicious\_stew\_effects

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=72&veaction=edit "编辑章节：suspicious_stew_effects") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=72 "编辑章节的源代码： suspicious_stew_effects")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)minecraft:suspicious\_stew\_effects：[谜之炖菜](/w/%E8%B0%9C%E4%B9%8B%E7%82%96%E8%8F%9C "谜之炖菜")的状态效果信息。
        -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)：一项状态效果信息。
            -   ![整型](/images/Data_node_int.svg?8d24f)duration：（默认为160）状态效果的时长，单位为[刻](/w/%E5%88%BB "刻")。
            -   ![字符串](/images/Data_node_string.svg?42545)\*  
                \*id：（命名空间ID）状态效果。

### swing\_animation

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=73&veaction=edit "编辑章节：swing_animation") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=73 "编辑章节的源代码： swing_animation")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:swing\_animation：使用此物品攻击时的动画。
        -   ![字符串](/images/Data_node_string.svg?42545)type：（默认为`whack`）摇摆动画类型。取值只能为`none`（轻微左右摇摆，第一人称下仅为物品上下略微移动）、`whack`（向前猛击，剑的默认攻击动画）、`stab`（矛戳刺攻击，被部分生物持有时还会有独特的手部动画）。
        -   ![整型](/images/Data_node_int.svg?8d24f)duration：（默认为6）动画播放的周期。

### tool

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=74&veaction=edit "编辑章节：tool") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=74 "编辑章节的源代码： tool")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:tool：物品的挖掘工具属性。
        -   ![布尔型](/images/Data_node_bool.svg?77754)can\_destroy\_blocks\_in\_creative：（默认为`true`）[创造模式](/w/%E5%88%9B%E9%80%A0%E6%A8%A1%E5%BC%8F "创造模式")玩家能否使用此物品破坏方块。
        -   ![整型](/images/Data_node_int.svg?8d24f)damage\_per\_block：（值≥0，默认为1）破坏硬度非0的方块时物品损失的耐久度。
        -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)default\_mining\_speed：（值≥0，默认为1）挖掘方块时的速度。
        -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)\*  
            \*rules：物品与对应可以挖掘的方块的映射列表。
            -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)：一项物品与方块列表的挖掘配置数据。
                -   ![字符串](/images/Data_node_string.svg?42545)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)\*  
                    \*blocks：此配置指定的有效方块。可以为一个`#`开头的[方块标签](/w/%E6%96%B9%E5%9D%97%E6%A0%87%E7%AD%BE "方块标签")、一个[方块ID](/w/Java%E7%89%88%E6%95%B0%E6%8D%AE%E5%80%BC "Java版数据值")、或一个方块ID的列表。
                -   ![布尔型](/images/Data_node_bool.svg?77754)correct\_for\_drops：此物品是否是所有上方指定方块的合适挖掘工具。
                -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)speed：覆盖所有上方指定方块的使用此物品挖掘时的挖掘速度。

示例

给予一把木锹，且这把木锹属于石头的适合挖掘工具：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s wooden_shovel[tool={rules:[{blocks:["stone"],correct_for_drops:True}]}]`

### tooltip\_display

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=75&veaction=edit "编辑章节：tooltip_display") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=75 "编辑章节的源代码： tooltip_display")\]

控制由组件产生的提示框文本可见性和提示框可见性。刷怪笼与试炼刷怪笼的介绍被视为由`block_entity_data`组件添加，由魔咒效果添加的属性修饰符文本被视为由`attribute_modifiers`组件添加。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:tooltip\_display：物品提示框的显示数据。
        -   ![布尔型](/images/Data_node_bool.svg?77754)hide\_tooltip：（默认为`false`）物品提示框是否总是隐藏。
        -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)hidden\_components：（默认为空列表）一个物品组件ID列表，列表内的所有组件提供的提示框文本都会被隐藏。如果组件不提供提示框文本，则对其没有效果。
            -   ![字符串](/images/Data_node_string.svg?42545)：（命名空间ID）一个物品组件。

示例

给予一个不显示提示框的[金斧](/w/%E9%87%91%E6%96%A7 "金斧")：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s golden_axe[tooltip_display={hide_tooltip:true}]`

### tooltip\_style

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=76&veaction=edit "编辑章节：tooltip_style") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=76 "编辑章节的源代码： tooltip_style")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:tooltip\_style：（命名空间ID）物品提示框外观。提示框外观分为两部分：背景由`<*命名空间*>:tooltip/<*路径*>_background`精灵图渲染，边框由`<*命名空间*>:tooltip/<*路径*>_frame`精灵图渲染。这两个精灵图都属于[GUI纹理](/w/%E7%BA%B9%E7%90%86#GUI纹理 "纹理")，默认会被解析为`assets/<*命名空间*>/textures/gui/sprites/tooltip/<*路径*>_background.png`和`assets/<*命名空间*>/textures/gui/sprites/tooltip/<*路径*>_frame.png`。

### trim

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=77&veaction=edit "编辑章节：trim") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=77 "编辑章节的源代码： trim")\]

该组件可以引用已存在的[盔甲纹饰定义文件](/w/%E7%9B%94%E7%94%B2%E7%BA%B9%E9%A5%B0%E5%AE%9A%E4%B9%89%E6%A0%BC%E5%BC%8F "盔甲纹饰定义格式")。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:trim：物品的[盔甲纹饰](/w/%E7%9B%94%E7%94%B2%E7%BA%B9%E9%A5%B0 "盔甲纹饰")信息。
        -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)\*  
            \*material：此盔甲纹饰的材料，命名空间ID或内联定义均可。
            -   ![字符串](/images/Data_node_string.svg?42545)\*  
                \*asset\_name：一个字符串，实际盔甲纹饰纹理的后缀。
            -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)\*  
                \*description：（[文本组件](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6 "文本组件")）盔甲纹饰材料的名称。
            -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)override\_armor\_assets：对于指定的[装备模型](/w/%E8%A3%85%E5%A4%87%E6%A8%A1%E5%9E%8B "装备模型")，使用指定纹理覆盖而不使用![字符串](/images/Data_node_string.svg?42545)asset\_name。
                -   ![字符串](/images/Data_node_string.svg?42545)<*装备模型命名空间ID*\>：一个字符串，实际盔甲纹饰纹理的后缀。
        -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)\*  
            \*pattern：此盔甲纹饰的图案，命名空间ID或内联定义均可。
            -   ![字符串](/images/Data_node_string.svg?42545)\*  
                \*asset\_id：（命名空间ID）用于推断盔甲纹饰纹理的位置。
            -   ![布尔型](/images/Data_node_bool.svg?77754)decal：（默认为`false`）是否使用贴花模式渲染盔甲纹饰（仅在非透明区域显示）。
            -   ![字符串](/images/Data_node_string.svg?42545)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)\*  
                \*description：（[文本组件](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6 "文本组件")）盔甲纹饰图案的名称。

### unbreakable

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=78&veaction=edit "编辑章节：unbreakable") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=78 "编辑章节的源代码： unbreakable")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:unbreakable：空标签，此组件存在时物品无法破坏，不存在耐久度。

示例

给予一把无法破坏的钻石镐：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s diamond_pickaxe[unbreakable={}]`

### use\_cooldown

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=79&veaction=edit "编辑章节：use_cooldown") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=79 "编辑章节的源代码： use_cooldown")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:use\_cooldown：设置物品的使用冷却行为。冷却时间会作用在一个“冷却组”上。
        -   ![字符串](/images/Data_node_string.svg?42545)cooldown\_group：（命名空间ID）设置物品冷却组。同冷却组的物品会同时受到同一个物品冷却影响，在冷却时间内所有同冷却组的物品都无法使用。如果此值不存在，游戏将以物品的命名空间ID作为冷却组ID使用。
        -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)\*  
            \*seconds：（值>0）物品使用后的冷却时间，单位为秒。

示例

给予16个每使用一次冷却一秒的[雪球](/w/%E9%9B%AA%E7%90%83 "雪球")：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s snowball[use_cooldown={seconds:1}] 16`

### use\_effects

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=80&veaction=edit "编辑章节：use_effects") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=80 "编辑章节的源代码： use_effects")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:use\_effects：设置物品被使用时的部分行为。
        -   ![布尔型](/images/Data_node_bool.svg?77754)can\_sprint：（默认为`false`）玩家在使用此物品时是否可以疾跑。
        -   ![布尔型](/images/Data_node_bool.svg?77754)interact\_vibrations：（默认为`true`）生物使用此物品时是否会发出`item_interact_finish`和`item_interact_start`[游戏事件](/w/%E6%B8%B8%E6%88%8F%E4%BA%8B%E4%BB%B6 "游戏事件")。此值为`false`或者此组件不存在时使用此物品不会发出这两个游戏事件。
        -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)speed\_multiplier：（0≤值≤1，默认为0.2）玩家使用此物品时的速度倍率。

### use\_remainder

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=81&veaction=edit "编辑章节：use_remainder") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=81 "编辑章节的源代码： use_remainder")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:use\_remainder\[失效：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]：控制物品在消耗使用且物品数量减少后游戏返还的物品。如果玩家物品栏在欲返还物品时已满，则掉落成为物品实体。
        
        -   物品共通标签，见[Template:Nbt inherit/itemnoslot/source](/w/Template:Nbt_inherit/itemnoslot/source "Template:Nbt inherit/itemnoslot/source")
    -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:use\_remainder\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]：控制物品在消耗使用且物品数量减少后游戏返还的物品。如果玩家物品栏在欲返还物品时已满，则掉落成为物品实体。
        
        -   物品模板，见[Template:Nbt inherit/item template/source](/w/Template:Nbt_inherit/item_template/source "Template:Nbt inherit/item template/source")

示例

给予16个能在使用后返回1个普通[雪球](/w/%E9%9B%AA%E7%90%83 "雪球")的雪球：`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s snowball[use_remainder={id:"snowball"}] 16`

### weapon

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=82&veaction=edit "编辑章节：weapon") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=82 "编辑章节的源代码： weapon")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:weapon：设置物品的武器数据。此组件存在时物品使用次数统计信息会在用此物品攻击时增加。
        -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)disable\_blocking\_for\_seconds：（值≥0，默认为0）攻击成功停用目标盾牌的秒数。
        -   ![整型](/images/Data_node_int.svg?8d24f)item\_damage\_per\_attack：（值≥0，默认为1）每次攻击对此物品造成的损伤值，即损耗的耐久度。

示例

给予一个下界合金剑，当攻击成功时停用目标盾牌60秒，且每次攻击损耗0耐久度：

`/[give](/w/%E5%91%BD%E4%BB%A4/give "命令/give") @s netherite_sword[weapon={disable_blocking_for_seconds:60,item_damage_per_attack:0}]`

### writable\_book\_content

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=83&veaction=edit "编辑章节：writable_book_content") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=83 "编辑章节的源代码： writable_book_content")\]

存储[书与笔](/w/%E4%B9%A6%E4%B8%8E%E7%AC%94 "书与笔")的数据。当玩家使用拥有此组件的书与笔或成书时会显示编辑界面，在讲台上打开拥有此组件的任意物品时游戏会显示每页的文本信息。此组件的优先级低于`written_book_content`组件。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:writable\_book\_content：[书与笔](/w/%E4%B9%A6%E4%B8%8E%E7%AC%94 "书与笔")的数据。
        -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)pages：（最多100个元素）书与笔内存储的页信息，必须为以下格式之一。
            -   ![字符串](/images/Data_node_string.svg?42545)：（长度不超过1024）书与笔内一页的文本信息。如果开启过滤，则代表文本信息与原信息一致。
            -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)：书与笔内一页的信息。
                -   ![字符串](/images/Data_node_string.svg?42545)filtered：（长度不超过1024）已过滤的文本信息。在开启过滤时，此字符串的优先级高于原始文本。被开启过滤的玩家更新时会删除原始文本并将此过滤文本作为原始文本，而被未开启过滤的玩家更新时会被移除。
                -   ![字符串](/images/Data_node_string.svg?42545)\*  
                    \*raw：（长度不超过1024）未过滤的文本原始信息。

### written\_book\_content

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=84&veaction=edit "编辑章节：written_book_content") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=84 "编辑章节的源代码： written_book_content")\]

存储[成书](/w/%E6%88%90%E4%B9%A6 "成书")的数据。当玩家使用拥有此组件的书与笔或成书，或在讲台上打开拥有此组件的任意物品时游戏会显示每页的文本信息。此组件的`title`字段也被视为物品的自定义名称。

所有拥有此组件的物品可以使用成书复制配方在工作台上进行复制。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:written\_book\_content：[成书](/w/%E6%88%90%E4%B9%A6 "成书")的数据。
        -   ![字符串](/images/Data_node_string.svg?42545)\*  
            \*author：成书的作者。
        -   ![整型](/images/Data_node_int.svg?8d24f)generation：（默认为`0`）决定成书的复制程度。可以为`0`（原稿），`1`（原稿的副本），`2`（副本的副本），`3`（破烂不堪）。
        -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)pages：成书内存储的页信息，必须使用以下格式之一。
            -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)：（[文本组件](/w/%E6%96%87%E6%9C%AC%E7%BB%84%E4%BB%B6 "文本组件")）成书内一页的信息。如果开启过滤，则代表过滤后文本信息与原信息一致。
            -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)：成书页内信息的另一种格式。如果采用复合标签定义文本组件，则只要存在`raw`字段就会以此格式解析。
                -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)filtered：（文本组件）已过滤的文本信息。在开启过滤时，此文本的优先级高于原始文本。被开启过滤的玩家更新时会删除原始文本并将此过滤文本作为原始文本，而被未开启过滤的玩家更新时会被移除。
                -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)\*  
                    \*raw：（文本组件）未过滤的原始信息。
        -   ![布尔型](/images/Data_node_bool.svg?77754)resolved：（默认为`false`）表示这本成书是否已经被解析，决定是否在打开成书时进行成书内文本的解析。
        -   ![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)\*  
            \*title：（长度不超过32）成书的标题信息。为![字符串](/images/Data_node_string.svg?42545)格式时，如果开启过滤，则代表此标题过滤后与原标题一致。
            -   ![字符串](/images/Data_node_string.svg?42545)filtered：（长度不超过32）已过滤的标题信息。在开启过滤时，此字符串优先级高于![字符串](/images/Data_node_string.svg?42545)raw。
            -   ![字符串](/images/Data_node_string.svg?42545)\*  
                \*raw：（长度不超过32）未过滤的标题原始信息。

### 实体变种组件

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=85&veaction=edit "编辑章节：实体变种组件") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=85 "编辑章节的源代码： 实体变种组件")\]

这些数据组件都可以作为实体组件，且专用于控制实体变种。若物品可以生成对应的实体（例如刷怪蛋、鸡蛋、画等）则使用这些物品生成实体时会生成指定的变种。例如，若棕色鸡蛋的`chicken/variant`组件值为`minecraft:cold`，则投掷后会生成寒带鸡而不是热带鸡。

#### axolotl/variant

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=86&veaction=edit "编辑章节：axolotl/variant") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=86 "编辑章节的源代码： axolotl/variant")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:axolotl/variant：[美西螈](/w/%E7%BE%8E%E8%A5%BF%E8%9E%88 "美西螈")的变种。取值只能为`lucy`（粉红色）、`wild`（棕色）、`gold`（金色）、`cyan`（青色）或`blue`（蓝色）。

#### cat/collar

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=87&veaction=edit "编辑章节：cat/collar") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=87 "编辑章节的源代码： cat/collar")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:cat/collar：猫的项圈颜色。取值为[染料颜色](/w/%E6%9F%93%E6%96%99%E9%A2%9C%E8%89%B2 "染料颜色")，即`white`、`orange`、`magenta`、`light_blue`、`yellow`、`lime`、`pink`、`gray`、`light_gray`、`cyan`、`purple`、`blue`、`brown`、`green`、`red`或`black`。

#### cat/sound\_variant

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=88&veaction=edit "编辑章节：cat/sound_variant") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=88 "编辑章节的源代码： cat/sound_variant")\]

[![](/images/thumb/Grass_Block_JE7_BE6.png/16px-Grass_Block_JE7_BE6.png?d27c1)](/w/File:Grass_Block_JE7_BE6.png)

**本段落包含会在[下一次更新](/w/%E8%AE%A1%E5%88%92%E7%89%88%E6%9C%AC "计划版本")中出现的内容。**

这些特性在[Java版26.1](/w/Java%E7%89%8826.1 "Java版26.1")的开发版本中加入。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:cat/sound\_variant：（命名空间ID）猫的音效变种。

#### cat/variant

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=89&veaction=edit "编辑章节：cat/variant") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=89 "编辑章节的源代码： cat/variant")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:cat/variant：（命名空间ID）[猫](/w/%E7%8C%AB "猫")的变种。

#### chicken/sound\_variant

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=90&veaction=edit "编辑章节：chicken/sound_variant") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=90 "编辑章节的源代码： chicken/sound_variant")\]

[![](/images/thumb/Grass_Block_JE7_BE6.png/16px-Grass_Block_JE7_BE6.png?d27c1)](/w/File:Grass_Block_JE7_BE6.png)

**本段落包含会在[下一次更新](/w/%E8%AE%A1%E5%88%92%E7%89%88%E6%9C%AC "计划版本")中出现的内容。**

这些特性在[Java版26.1](/w/Java%E7%89%8826.1 "Java版26.1")的开发版本中加入。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:chicken/sound\_variant：（命名空间ID）鸡的音效变种。

#### chicken/variant

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=91&veaction=edit "编辑章节：chicken/variant") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=91 "编辑章节的源代码： chicken/variant")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:chicken/variant：（命名空间ID）[鸡](/w/%E9%B8%A1 "鸡")的变种。

#### cow/sound\_variant

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=92&veaction=edit "编辑章节：cow/sound_variant") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=92 "编辑章节的源代码： cow/sound_variant")\]

[![](/images/thumb/Grass_Block_JE7_BE6.png/16px-Grass_Block_JE7_BE6.png?d27c1)](/w/File:Grass_Block_JE7_BE6.png)

**本段落包含会在[下一次更新](/w/%E8%AE%A1%E5%88%92%E7%89%88%E6%9C%AC "计划版本")中出现的内容。**

这些特性在[Java版26.1](/w/Java%E7%89%8826.1 "Java版26.1")的开发版本中加入。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:cow/sound\_variant：（命名空间ID）牛的音效变种。

#### cow/variant

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=93&veaction=edit "编辑章节：cow/variant") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=93 "编辑章节的源代码： cow/variant")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:cow/variant：（命名空间ID）[牛](/w/%E7%89%9B "牛")的变种。

#### fox/variant

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=94&veaction=edit "编辑章节：fox/variant") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=94 "编辑章节的源代码： fox/variant")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:fox/variant：[狐狸](/w/%E7%8B%90%E7%8B%B8 "狐狸")的变种。取值只能为`red`（红色）、`snow`（白色）。

#### frog/variant

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=95&veaction=edit "编辑章节：frog/variant") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=95 "编辑章节的源代码： frog/variant")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:frog/variant：（命名空间ID）[青蛙](/w/%E9%9D%92%E8%9B%99 "青蛙")的变种。

#### horse/variant

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=96&veaction=edit "编辑章节：horse/variant") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=96 "编辑章节的源代码： horse/variant")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:horse/variant：[马](/w/%E9%A9%AC "马")的基础毛色。取值只能为`white`（白色）、`creamy`（奶油色）、`chestnut`（栗色）、`brown`（褐色）、`black`（黑色）、`gray`（灰色）或`dark_brown`（深褐色）。

#### llama/variant

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=97&veaction=edit "编辑章节：llama/variant") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=97 "编辑章节的源代码： llama/variant")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:llama/variant：[羊驼](/w/%E7%BE%8A%E9%A9%BC "羊驼")的变种。取值只能为`creamy`（沙褐色）、`white`（奶油色）、`brown`（棕色）或`gray`（灰色）。[行商羊驼](/w/%E8%A1%8C%E5%95%86%E7%BE%8A%E9%A9%BC "行商羊驼")也使用此组件。

#### mooshroom/variant

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=98&veaction=edit "编辑章节：mooshroom/variant") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=98 "编辑章节的源代码： mooshroom/variant")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:mooshroom/variant：[哞菇](/w/%E5%93%9E%E8%8F%87 "哞菇")的变种。取值只能为`red`（红色）或`brown`（棕色）。

#### painting/variant

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=99&veaction=edit "编辑章节：painting/variant") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=99 "编辑章节的源代码： painting/variant")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:painting/variant：（命名空间ID）[画](/w/%E7%94%BB "画")的变种。

#### parrot/variant

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=100&veaction=edit "编辑章节：parrot/variant") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=100 "编辑章节的源代码： parrot/variant")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:parrot/variant：[鹦鹉](/w/%E9%B9%A6%E9%B9%89 "鹦鹉")的变种。取值只能为`red_blue`（红色）、`blue`（蓝色）、`green`（绿色）、`yellow_blue`（青色）或`gray`（灰色）。

#### pig/sound\_variant

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=101&veaction=edit "编辑章节：pig/sound_variant") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=101 "编辑章节的源代码： pig/sound_variant")\]

[![](/images/thumb/Grass_Block_JE7_BE6.png/16px-Grass_Block_JE7_BE6.png?d27c1)](/w/File:Grass_Block_JE7_BE6.png)

**本段落包含会在[下一次更新](/w/%E8%AE%A1%E5%88%92%E7%89%88%E6%9C%AC "计划版本")中出现的内容。**

这些特性在[Java版26.1](/w/Java%E7%89%8826.1 "Java版26.1")的开发版本中加入。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:pig/sound\_variant：（命名空间ID）猪的音效变种。

#### pig/variant

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=102&veaction=edit "编辑章节：pig/variant") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=102 "编辑章节的源代码： pig/variant")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:pig/variant：（命名空间ID）[猪](/w/%E7%8C%AA "猪")的变种。

#### rabbit/variant

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=103&veaction=edit "编辑章节：rabbit/variant") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=103 "编辑章节的源代码： rabbit/variant")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:rabbit/variant：[兔子](/w/%E5%85%94%E5%AD%90 "兔子")的变种。取值只能为`brown`（褐色）、`white`（白色）、`black`（黑色）、`white_splotched`（黑白相间）、`gold`（金色）、`salt`（胡椒盐色）或`evil`（杀手兔）。

#### salmon/size

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=104&veaction=edit "编辑章节：salmon/size") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=104 "编辑章节的源代码： salmon/size")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:salmon/size：鲑鱼的体型尺寸。取值只能为`small`（小型）、`medium`（中型）或`large`（大型）。

#### sheep/color

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=105&veaction=edit "编辑章节：sheep/color") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=105 "编辑章节的源代码： sheep/color")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:sheep/color：绵羊的毛色。取值为[染料颜色](/w/%E6%9F%93%E6%96%99%E9%A2%9C%E8%89%B2 "染料颜色")，即`white`、`orange`、`magenta`、`light_blue`、`yellow`、`lime`、`pink`、`gray`、`light_gray`、`cyan`、`purple`、`blue`、`brown`、`green`、`red`或`black`。

#### shulker/color

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=106&veaction=edit "编辑章节：shulker/color") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=106 "编辑章节的源代码： shulker/color")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:shulker/color：潜影贝的外壳颜色，如果此组件不存在，则潜影贝使用默认的颜色。取值为[染料颜色](/w/%E6%9F%93%E6%96%99%E9%A2%9C%E8%89%B2 "染料颜色")，即`white`、`orange`、`magenta`、`light_blue`、`yellow`、`lime`、`pink`、`gray`、`light_gray`、`cyan`、`purple`、`blue`、`brown`、`green`、`red`或`black`。

#### tropical\_fish/base\_color

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=107&veaction=edit "编辑章节：tropical_fish/base_color") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=107 "编辑章节的源代码： tropical_fish/base_color")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:tropical\_fish/base\_color：热带鱼的基础颜色。取值为[染料颜色](/w/%E6%9F%93%E6%96%99%E9%A2%9C%E8%89%B2 "染料颜色")，即`white`、`orange`、`magenta`、`light_blue`、`yellow`、`lime`、`pink`、`gray`、`light_gray`、`cyan`、`purple`、`blue`、`brown`、`green`、`red`或`black`。

#### tropical\_fish/pattern

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=108&veaction=edit "编辑章节：tropical_fish/pattern") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=108 "编辑章节的源代码： tropical_fish/pattern")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:tropical\_fish/pattern：热带鱼的花纹类型。取值只能为`kob`（石首类）、`sunstreak`（日纹类）、`snooper`（窥伺类）、`dasher`（速跃类）、`brinely`（咸水类）、`spotty`（多斑类）、`flopper`（飞翼类）、`stripey`（条纹类）、`glitter`（闪鳞类）、`blockfish`（方身类）、`betty`（背蒂类）或`clayfish`（陶鱼类）。

#### tropical\_fish/pattern\_color

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=109&veaction=edit "编辑章节：tropical_fish/pattern_color") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=109 "编辑章节的源代码： tropical_fish/pattern_color")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:tropical\_fish/pattern\_color：热带鱼的花纹颜色。取值为[染料颜色](/w/%E6%9F%93%E6%96%99%E9%A2%9C%E8%89%B2 "染料颜色")，即`white`、`orange`、`magenta`、`light_blue`、`yellow`、`lime`、`pink`、`gray`、`light_gray`、`cyan`、`purple`、`blue`、`brown`、`green`、`red`或`black`。

#### villager/variant

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=110&veaction=edit "编辑章节：villager/variant") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=110 "编辑章节的源代码： villager/variant")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:villager/variant：（命名空间ID）[村民](/w/%E6%9D%91%E6%B0%91 "村民")类型，取值可以为`desert`（沙漠）、`jungle`（丛林）、`plains`（默认）、`savanna`（热带草原）、`snow`（雪原）、`swamp`（沼泽）和`taiga`（针叶林）。

#### wolf/collar

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=111&veaction=edit "编辑章节：wolf/collar") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=111 "编辑章节的源代码： wolf/collar")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:wolf/collar：狼的项圈颜色。取值为[染料颜色](/w/%E6%9F%93%E6%96%99%E9%A2%9C%E8%89%B2 "染料颜色")，即`white`、`orange`、`magenta`、`light_blue`、`yellow`、`lime`、`pink`、`gray`、`light_gray`、`cyan`、`purple`、`blue`、`brown`、`green`、`red`或`black`。

#### wolf/sound\_variant

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=112&veaction=edit "编辑章节：wolf/sound_variant") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=112 "编辑章节的源代码： wolf/sound_variant")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:wolf/sound\_variant：（命名空间ID）狼的音效变种。

#### wolf/variant

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=113&veaction=edit "编辑章节：wolf/variant") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=113 "编辑章节的源代码： wolf/variant")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:wolf/variant：（命名空间ID）[狼](/w/%E7%8B%BC "狼")的变种。

#### zombie\_nautilus/variant

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=114&veaction=edit "编辑章节：zombie_nautilus/variant") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=114 "编辑章节的源代码： zombie_nautilus/variant")\]

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![字符串](/images/Data_node_string.svg?42545)minecraft:zombie\_nautilus/variant：（命名空间ID）[僵尸鹦鹉螺](/w/%E5%83%B5%E5%B0%B8%E9%B9%A6%E9%B9%89%E8%9E%BA "僵尸鹦鹉螺")的变种。

### 临时组件

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=115&veaction=edit "编辑章节：临时组件") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=115 "编辑章节的源代码： 临时组件")\]

这些数据组件仅网络同步，不会被保存也不会被加载。玩家只能判断组件是否存在，而不能主动设置或获取其组件信息。

#### additional\_trade\_cost

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=116&veaction=edit "编辑章节：additional_trade_cost") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=116 "编辑章节的源代码： additional_trade_cost")\]

[![](/images/thumb/Grass_Block_JE7_BE6.png/16px-Grass_Block_JE7_BE6.png?d27c1)](/w/File:Grass_Block_JE7_BE6.png)

**本段落包含会在[下一次更新](/w/%E8%AE%A1%E5%88%92%E7%89%88%E6%9C%AC "计划版本")中出现的内容。**

这些特性在[Java版26.1](/w/Java%E7%89%8826.1 "Java版26.1")的开发版本中加入。

村民生成交易时收购物品的基础增加量，当交易选项生成时会立刻被移除。

同步格式

-   -   ![整型](/images/Data_node_int.svg?8d24f)minecraft:additional\_trade\_cost：村民收购物品的增加量。

#### creative\_slot\_lock

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=117&veaction=edit "编辑章节：creative_slot_lock") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=117 "编辑章节的源代码： creative_slot_lock")\]

阻止玩家在物品栏内与此物品交互。默认附加到创造模式物品栏“已保存的快捷栏”中表示此快捷栏未保存的纸上。

同步格式

-   -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:creative\_slot\_lock：空标签，此组件存在时此物品若在创造模式物品栏内则玩家无法与其交互。

#### map\_post\_processing

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=118&veaction=edit "编辑章节：map_post_processing") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=118 "编辑章节的源代码： map_post_processing")\]

同步地图的缩放等级和锁定信息。默认附加到进行地图缩小或地图锁定操作的制图台的输出槽位或进行地图缩小配方的工作台的输出槽位的物品上，暂时提供输出物品的地图缩放信息或锁定信息。当物品从输出槽位取下时其组件会被立刻移除。

同步格式

-   -   ![整型](/images/Data_node_int.svg?8d24f)minecraft:map\_post\_processing：为0时使`map_id`组件额外增加“已锁定”行；为1时使`map_id`组件使用“![字节型](/images/Data_node_byte.svg?eb0e0)scale+1”而不是![字节型](/images/Data_node_byte.svg?eb0e0)scale显示地图比例缩放信息。

## 历史

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=119&veaction=edit "编辑章节：历史") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=119 "编辑章节的源代码： 历史")\]

[Java版](/w/Java%E7%89%88%E7%89%88%E6%9C%AC%E8%AE%B0%E5%BD%95 "Java版版本记录")

[1.20.5](/w/Java%E7%89%881.20.5 "Java版1.20.5")

[24w09a](/w/Java%E7%89%8824w09a "Java版24w09a")

加入了数据组件。

[24w10a](/w/Java%E7%89%8824w10a "Java版24w10a")

为`attribute_modifiers`、​`dyed_color`、​`enchantments`、​`potion_contents`、​`profile`和​`stored_enchantments`组件加入了简化定义。

将`lodestone_target`组件重命名为`lodestone_tracker`。

现在`container`组件适用于所有的[容器](/w/%E5%AE%B9%E5%99%A8 "容器")而不只有[潜影盒](/w/%E6%BD%9C%E5%BD%B1%E7%9B%92 "潜影盒")。

[24w12a](/w/Java%E7%89%8824w12a "Java版24w12a")

现在`unbreakable`组件与[耐久](/w/%E8%80%90%E4%B9%85 "耐久")魔咒兼容。[\[11\]](#cite_note-11)

加入了`fire_resistant`、​`food`、​`hide_tooltip`、​`max_damage`、​`max_stack_size`、​`rarity`和​`tool`组件。

[24w13a](/w/Java%E7%89%8824w13a "Java版24w13a")

加入了`item_name`和​`ominous_bottle_amplifier`组件。

现在方块实体会保存方块物品的全部组件，而非只能由继承序列化处理的组件。

[24w14a](/w/Java%E7%89%8824w14a "Java版24w14a")

将`writable_book_content`和​`written_book_content`组件中的未过滤的文本组件原始信息`text`被重命名为`raw`以避免歧义。[\[12\]](#cite_note-12)

现在`profile`组件中不存在`name`而存在`id`时，可由UUID解析玩家档案数据。[\[13\]](#cite_note-13)

[pre1](/w/Java%E7%89%881.20.5-pre1 "Java版1.20.5-pre1")

移除了`food`组件的营养价值字段![单精度浮点数](/images/Data_node_float.svg?ae55e)saturation\_modifier，现在由新字段![单精度浮点数](/images/Data_node_float.svg?ae55e)saturation直接指定饱和度。

[1.21](/w/Java%E7%89%881.21 "Java版1.21")

[24w19a](/w/Java%E7%89%8824w19a "Java版24w19a")

`custom_data`组件现在可以使用SNBT字符串定义。

为`food`组件加入了![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)using\_converts\_to字段。

[24w21a](/w/Java%E7%89%8824w21a "Java版24w21a")

加入了`jukebox_playable`组件。

向`attribute_modifiers`组件加入了![字符串](/images/Data_node_string.svg?42545)id字段，取代![字符串](/images/Data_node_string.svg?42545)name和![整型数组](/images/Data_node_int-array.svg?546e8)uuid字段。

[1.21.2](/w/Java%E7%89%881.21.2 "Java版1.21.2")

[24w33a](/w/Java%E7%89%8824w33a "Java版24w33a")

加入了`enchantable`和​`repairable`组件。

将`enchantments`和`stored_enchantments`组件中的魔咒等级下限由0提升至1。

现在`written_book_content`组件中的`title`具有比`custom_name`和`item_name`更高的优先级。

现在`instrument`组件的内联定义需要`description`字段。

[24w34a](/w/Java%E7%89%8824w34a "Java版24w34a")

加入了`consumable`、​`use_cooldown`和​`use_remainder`组件。

移除了`food`组件中![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)effects、![单精度浮点数](/images/Data_node_float.svg?ae55e)eat\_seconds和![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)using\_converts\_to字段，现在此组件仅作为食物数据容器。

此前字段格式：

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)food 物品堆叠组件
    -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)eat\_seconds：（默认为1.6）物品食用所需要的时间，以秒为单位。小于0.05秒（1游戏刻）时无法被食用。下列物品使用此标签无法修改食用时间：[蜂蜜瓶](/w/%E8%9C%82%E8%9C%9C%E7%93%B6 "蜂蜜瓶")（40游戏刻（2秒））、[药水](/w/%E8%8D%AF%E6%B0%B4 "药水")（32游戏刻（1.6秒））、[不祥之瓶](/w/%E4%B8%8D%E7%A5%A5%E4%B9%8B%E7%93%B6 "不祥之瓶")（32游戏刻（1.6秒））。
    -   ![NBT列表/JSON数组](/images/Data_node_list.svg?d6aa9)effects：物品被食用后可能被赋予的状态效果。
        -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)：一项可能施加的状态效果。
            -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)\*  
                \*effect：一项状态效果。
                -   ![字符串](/images/Data_node_string.svg?42545)\*  
                    \*id：状态效果的[命名空间ID](/w/%E5%91%BD%E5%90%8D%E7%A9%BA%E9%97%B4ID "命名空间ID")。
                -   ![布尔型](/images/Data_node_bool.svg?77754)ambient：表示状态效果是否是被信标添加的。如果不存在则为`false`。
                -   ![字节型](/images/Data_node_byte.svg?eb0e0)amplifier：（无符号8位整数，不小于0且不大于255）状态效果的等级。0表示等级1，以此类推。由于此数字为无符号整数，当值超过127时显示为负数但实际为正数，保存数字*s*和实际代表数字*a*的关系为a\=256+s。如果不存在则为0。
                -   ![整型](/images/Data_node_int.svg?8d24f)duration：距离状态效果失效的时间刻数。如果此值为-1，则此状态效果不会失效。如果不存在则为0。
                -   ![布尔型](/images/Data_node_bool.svg?77754)show\_icon：表示是否显示状态效果的图标。如果不存在则与![布尔型](/images/Data_node_bool.svg?77754)show\_particles值相同。
                -   ![布尔型](/images/Data_node_bool.svg?77754)show\_particles：表示是否显示粒子效果。如果不存在则为`true`。
            -   ![单精度浮点数](/images/Data_node_float.svg?ae55e)probability：（0≤值≤1，默认为1）食用后施加此状态效果的概率。
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)using\_converts\_to：物品被食用后返还的物品。
        -   ![字符串](/images/Data_node_string.svg?42545)\*  
            \*id：（[命名空间ID](/w/%E5%91%BD%E5%90%8D%E7%A9%BA%E9%97%B4ID "命名空间ID")）表示某种类的物品堆叠。不能为[空气](/w/%E7%A9%BA%E6%B0%94 "空气")（`air`）。
        -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components：关于当前物品的额外信息。此标签对于多数物品来说都是非必需项。
            -   ![任意类型](/images/Data_node_any.svg?d406c)<*物品堆叠组件*\>：一个物品堆叠组件，标签名称为物品组件的命名空间ID，类型则与对应的物品组件有关。在设置此数据时可以不写命名空间，但在导出时命名空间ID会自行加上`minecraft:`前缀。

现在`custom_name`组件具有比`written_book_content`中`title`字段更高的优先级。

[24w35a](/w/Java%E7%89%8824w35a "Java版24w35a")

现在`use_cooldown`组件使用![字符串](/images/Data_node_string.svg?42545)cooldown\_group字段而非![字符串](/images/Data_node_string.svg?42545)cooldownGroup字段。

现在`bucket_entity_data`组件使用![字符串](/images/Data_node_string.svg?42545)type字段用来保存鲑鱼的尺寸。

[24w36a](/w/Java%E7%89%8824w36a "Java版24w36a")

加入了`equippable`、​`glider`、​`item_model`和​`tooltip_style`组件。

[24w37a](/w/Java%E7%89%8824w37a "Java版24w37a")

加入了`damage_resistant`和​`death_protection`组件。

移除了`fire_resistant`组件。

为`potion_contents`组件加入了![字符串](/images/Data_node_string.svg?42545)custom\_name字段。

为`equippable`组件加入了![布尔型](/images/Data_node_bool.svg?77754)swappable和![布尔型](/images/Data_node_bool.svg?77754)damage\_on\_hurt字段。

`item_name`组件的优先级现在永远最低，因此它能被`potion_contents`之类影响物品名称的组件覆盖。

现在所有具有`map_id`组件的物品能像[地图](/w/%E5%9C%B0%E5%9B%BE "地图")一样被展开，在[物品展示框](/w/%E7%89%A9%E5%93%81%E5%B1%95%E7%A4%BA%E6%A1%86 "物品展示框")上直角旋转，作为地图被复制、锁定或拓展，且在未启用高级提示框时的物品提示框中显示地图编号。

[24w39a](/w/Java%E7%89%8824w39a "Java版24w39a")

`lock`组件现在是一个物品谓词，此前为任意字符串。

[pre1](/w/Java%E7%89%881.21.2-pre1 "Java版1.21.2-pre1")

为`equippable`组件加入了![字符串](/images/Data_node_string.svg?42545)camera\_overlay字段。

[1.21.4](/w/Java%E7%89%881.21.4 "Java版1.21.4")

[24w44a](/w/Java%E7%89%8824w44a "Java版24w44a")

现在`consumable`组件的第一人称格挡动画可以正常显示。[\[14\]](#cite_note-14)

[24w45a](/w/Java%E7%89%8824w45a "Java版24w45a")

修改了`custom_model_data`的格式。

现在`item_model`组件不再设置物品模型，而是设置物品模型映射。

将`equippable`组件的![字符串](/images/Data_node_string.svg?42545)model字段重命名为![字符串](/images/Data_node_string.svg?42545)asset\_id。

[24w46a](/w/Java%E7%89%8824w46a "Java版24w46a")

为`consumable`组件的![字符串](/images/Data_node_string.svg?42545)animation字段添加了新的可选值`bundle`。

[1.21.5](/w/Java%E7%89%881.21.5 "Java版1.21.5")

[25w02a](/w/Java%E7%89%8825w02a "Java版25w02a")

向`tool`组件加入新可选字段![布尔型](/images/Data_node_bool.svg?77754)can\_destroy\_blocks\_in\_creative。

加入了`potion_duration_scale`和​`weapon`组件。

[25w03a](/w/Java%E7%89%8825w03a "Java版25w03a")

将`weapon`组件的![整型](/images/Data_node_int.svg?8d24f)damage\_per\_attack字段重命名为![整型](/images/Data_node_int.svg?8d24f)item\_damage\_per\_attack。

向`equippable`组件加入新字段![布尔型](/images/Data_node_bool.svg?77754)equip\_on\_interaction。

加入了`axolotl/variant`、​`cat/collar`、​`cat/variant`、​`fox/variant`、​`frog/variant`、​`horse/variant`、​`llama/variant`、​`mooshroom/variant`、​`parrot/variant`、​`painting/variant`、​`pig/variant`、​`rabbit/variant`、​`salmon/size`、​`sheep/color`、​`shulker/color`、​`tropical_fish/pattern`、​`tropical_fish/base_color`、​`tropical_fish/pattern_color`、​`villager/variant`、​`wolf/collar`和​`wolf/variant`组件。

移除了`bucket_entity_data`的![整型](/images/Data_node_int.svg?8d24f)BucketVariantTag和![字符串](/images/Data_node_string.svg?42545)type字段，由对应的变种组件保存。

[25w04a](/w/Java%E7%89%8825w04a "Java版25w04a")

将`weapon`组件的![布尔型](/images/Data_node_bool.svg?77754)can\_disable\_blocking字段重做为![单精度浮点数](/images/Data_node_float.svg?ae55e)disable\_blocking\_for\_seconds。

加入了`blocks_attacks`、​`break_sound`、​`provides_banner_patterns`、​`provides_trim_material`和​`tooltip_display`组件。

移除了`hide_additional_tooltip`和​`hide_tooltip`组件。

移除了`attribute_modifiers`、​`can_break`、​`can_place_on`、​`dyed_color`、​`enchantments`、​`stored_enchantments`、​`jukebox_playable`、​`unbreakable`和​`trim`组件的![布尔型](/images/Data_node_bool.svg?77754)show\_in\_tooltip字段。

现在具有简化定义的组件始终以简化形式存储，`jukebox_playable`组件直接指定命名空间ID，`dyed_color`组件可接受RGB数组。

现在`villager/variant`组件对僵尸村民可用，`potion_contents`和​`potion_duration_scale`组件对区域效果云可用。

[25w05a](/w/Java%E7%89%8825w05a "Java版25w05a")

加入了`cow/variant`组件。

为`blocks_attacks`组件加入了![字符串](/images/Data_node_string.svg?42545)bypassed\_by和![单精度浮点数](/images/Data_node_float.svg?ae55e)horizontal\_blocking\_angle字段。

现在`provides_banner_patterns`接受的旗帜图案标签需带有`#`前缀。

[25w06a](/w/Java%E7%89%8825w06a "Java版25w06a")

加入了`chicken/variant`组件。

[25w08a](/w/Java%E7%89%8825w08a "Java版25w08a")

加入了`wolf/sound_variant`组件。

[1.21.6](/w/Java%E7%89%881.21.6 "Java版1.21.6")

[25w15a](/w/Java%E7%89%8825w15a "Java版25w15a")

为`attribute_modifiers`组件的属性修饰符加入了![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)display字段。

[25w16a](/w/Java%E7%89%8825w16a "Java版25w16a")

现在`painting/variant`组件只接受命名空间ID，不再接受内联定义。

[25w20a](/w/Java%E7%89%8825w20a "Java版25w20a")

向`equippable`组件加入![布尔型](/images/Data_node_bool.svg?77754)can\_be\_sheared和![字符串](/images/Data_node_string.svg?42545)![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)shearing\_sound字段。

[1.21.9](/w/Java%E7%89%881.21.9 "Java版1.21.9")

[25w31a](/w/Java%E7%89%8825w31a "Java版25w31a")

`entity_data`和`block_entity_data`组件不再支持字符串格式定义。

[25w34a](/w/Java%E7%89%8825w34a "Java版25w34a")

修改了`profile`组件的解析行为，不再自动存储解析的玩家游戏档案数据，现在若未指定游戏档案则动态获取档案数据。

[pre1](/w/Java%E7%89%881.21.9-pre1 "Java版1.21.9-pre1")

现在`profile`组件可指定本地纹理以进行覆盖，并使其对玩家模型生效。

[1.21.11](/w/Java%E7%89%881.21.11 "Java版1.21.11")

[25w41a](/w/Java%E7%89%8825w41a "Java版25w41a")

加入了`damage_type`、​`kinetic_weapon`、​`minimum_attack_charge`、​`piercing_weapon`、​`swing_animation`和​`use_effects`组件。

将`consumable`组件的动画`spear`重命名为`trident`，被新的`spear`取代。

[25w45a](/w/Java%E7%89%8825w45a "Java版25w45a")

`entity_data`和`block_entity_data`组件再次支持字符串格式定义。

[25w46a](/w/Java%E7%89%8825w46a "Java版25w46a")

`use_effects`组件加入了新字段![布尔型](/images/Data_node_bool.svg?77754)interact\_vibrations：使用该物品时是否触发`minecraft:item_interact_start`和`minecraft:item_interact_finish`游戏事件。默认为`true`。

[pre1](/w/Java%E7%89%881.21.11-pre1 "Java版1.21.11-pre1")

加入了`attack_range`组件。

[Java版（即将到来）](/w/%E8%AE%A1%E5%88%92%E7%89%88%E6%9C%AC#Java版 "计划版本")

[26.1](/w/Java%E7%89%8826.1 "Java版26.1")

[snapshot-1](/w/Java%E7%89%8826.1-snapshot-1 "Java版26.1-snapshot-1")

加入了`additional_trade_cost`组件。

[snapshot-5](/w/Java%E7%89%8826.1-snapshot-5 "Java版26.1-snapshot-5")

加入了`dye`组件。

[snapshot-7](/w/Java%E7%89%8826.1-snapshot-7 "Java版26.1-snapshot-7")

加入了`cat/sound_variant`、`chicken/sound_variant`、`cow/sound_variant`和`pig/sound_variant`组件。

### 已移除的组件

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=120&veaction=edit "编辑章节：已移除的组件") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=120 "编辑章节的源代码： 已移除的组件")\]

#### fire\_resistant

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=121&veaction=edit "编辑章节：fire_resistant") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=121 "编辑章节的源代码： fire_resistant")\]

控制物品是否无法在熔岩或火中燃烧，且装备时是否不会因为受到火焰伤害而消耗耐久度。

此组件已被`[damage_resistant](#damage_resistant)`组件替代。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:fire\_resistant：空标签，此组件存在时物品不受火焰伤害影响。

#### hide\_additional\_tooltip

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=122&veaction=edit "编辑章节：hide_additional_tooltip") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=122 "编辑章节的源代码： hide_additional_tooltip")\]

隐藏物品的提示框文本信息。部分由组件产生的提示框文本由对应组件的![布尔型](/images/Data_node_bool.svg?77754)show\_in\_tooltip决定。

此组件已被`[tooltip_display](#tooltip_display)`组件替代。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:hide\_additional\_tooltip：空标签，此组件存在时提示框不会显示附加信息。

#### hide\_tooltip

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=123&veaction=edit "编辑章节：hide_tooltip") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=123 "编辑章节的源代码： hide_tooltip")\]

隐藏物品的提示框。

此组件已被`[tooltip_display](#tooltip_display)`组件替代。

-   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)components
    -   ![NBT复合标签/JSON对象](/images/Data_node_structure.svg?3a597)minecraft:hide\_tooltip：空标签，此组件存在时不会渲染提示框。

## 参考

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=124&veaction=edit "编辑章节：参考") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=124 "编辑章节的源代码： 参考")\]

1.  [↑](#cite_ref-1) [Minecraft Snapshot 24w39a](https://www.minecraft.net/article/minecraft-snapshot-24w39a) — [Minecraft.net](/w/Minecraft.net "Minecraft.net")。
2.  ↑ [2.0](#cite_ref-24w09a_2-0) [2.1](#cite_ref-24w09a_2-1) [Minecraft Snapshot 24w09a](https://www.minecraft.net/article/minecraft-snapshot-24w09a) — [Minecraft.net](/w/Minecraft.net "Minecraft.net")。
3.  [↑](#cite_ref-3) [MC-269629](https://bugs.mojang.com/browse/MC-269629 "mojira:MC-269629") — 漏洞状态为“不予修复”。
4.  [↑](#cite_ref-4) [MC-269631](https://bugs.mojang.com/browse/MC-269631 "mojira:MC-269631") — 漏洞状态为“不予修复”。
5.  [↑](#cite_ref-5) [MC-269655](https://bugs.mojang.com/browse/MC-269655 "mojira:MC-269655") — 漏洞状态为“不予修复”。
6.  [↑](#cite_ref-6) [MC-269640](https://bugs.mojang.com/browse/MC-269640 "mojira:MC-269640") — 漏洞状态为“不予修复”。
7.  [↑](#cite_ref-7) [MC-269648](https://bugs.mojang.com/browse/MC-269648 "mojira:MC-269648") — 漏洞状态为“不予修复”。
8.  [↑](#cite_ref-8) [MC-269658](https://bugs.mojang.com/browse/MC-269658 "mojira:MC-269658") — 漏洞状态为“不予修复”。
9.  [↑](#cite_ref-9) [MC-269722](https://bugs.mojang.com/browse/MC-269722 "mojira:MC-269722") — 漏洞状态为“不予修复”。
10.  [↑](#cite_ref-10) [MC-269723](https://bugs.mojang.com/browse/MC-269723 "mojira:MC-269723") — 漏洞状态为“不予修复”。
11.  [↑](#cite_ref-11) [MC-268510](https://bugs.mojang.com/browse/MC-268510 "mojira:MC-268510") — 漏洞状态为“已修复”。
12.  [↑](#cite_ref-12) [MC-269677](https://bugs.mojang.com/browse/MC-269677 "mojira:MC-269677") — 漏洞状态为“已修复”。
13.  [↑](#cite_ref-13) [MC-269983](https://bugs.mojang.com/browse/MC-269983 "mojira:MC-269983") — 漏洞状态为“已修复”。
14.  [↑](#cite_ref-14) [MC-275917](https://bugs.mojang.com/browse/MC-275917 "mojira:MC-275917") — 漏洞状态为“已修复”。

## 导航

\[[编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?section=125&veaction=edit "编辑章节：导航") | [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit&section=125 "编辑章节的源代码： 导航")\]

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
-   数据组件
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

-   [查](/w/Template:Navbox_Java_files "Template:Navbox Java files")
-   [论](/w/Special:TalkPage/Template:Navbox_Java_files "Special:TalkPage/Template:Navbox Java files")
-   [编](/w/Special:EditPage/Template:Navbox_Java_files "Special:EditPage/Template:Navbox Java files")

[Java版](/w/Java%E7%89%88 "Java版")游戏文件

通用文件

-   [版本信息文件格式](/w/%E7%89%88%E6%9C%AC%E4%BF%A1%E6%81%AF%E6%96%87%E4%BB%B6%E6%A0%BC%E5%BC%8F "版本信息文件格式")
-   [信任符号链接列表文件格式](/w/%E4%BF%A1%E4%BB%BB%E7%AC%A6%E5%8F%B7%E9%93%BE%E6%8E%A5%E5%88%97%E8%A1%A8%E6%96%87%E4%BB%B6%E6%A0%BC%E5%BC%8F "信任符号链接列表文件格式")
-   [玩家档案缓存存储格式](/w/%E7%8E%A9%E5%AE%B6%E6%A1%A3%E6%A1%88%E7%BC%93%E5%AD%98%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "玩家档案缓存存储格式")
-   [性能分析报告文件](/w/%E6%80%A7%E8%83%BD%E5%88%86%E6%9E%90%E6%8A%A5%E5%91%8A%E6%96%87%E4%BB%B6 "性能分析报告文件")
-   [崩溃报告文件](/w/%E5%B4%A9%E6%BA%83%E6%8A%A5%E5%91%8A%E6%96%87%E4%BB%B6 "崩溃报告文件")

客户端文件

-   [散列资源文件](/w/%E6%95%A3%E5%88%97%E8%B5%84%E6%BA%90%E6%96%87%E4%BB%B6 "散列资源文件")
-   [客户端核心文件](/w/%E5%AE%A2%E6%88%B7%E7%AB%AF%E6%A0%B8%E5%BF%83%E6%96%87%E4%BB%B6 "客户端核心文件")
    -   [客户端数据生成器](/w/%E6%95%B0%E6%8D%AE%E7%94%9F%E6%88%90%E5%99%A8 "数据生成器")
-   [客户端选项文件格式](/w/%E5%AE%A2%E6%88%B7%E7%AB%AF%E9%80%89%E9%A1%B9%E6%96%87%E4%BB%B6%E6%A0%BC%E5%BC%8F "客户端选项文件格式")
-   [调试选项文件格式](/w/%E8%B0%83%E8%AF%95%E9%80%89%E9%A1%B9%E6%96%87%E4%BB%B6%E6%A0%BC%E5%BC%8F "调试选项文件格式")
-   [下载缓存目录](/w/%E4%B8%8B%E8%BD%BD%E7%BC%93%E5%AD%98%E7%9B%AE%E5%BD%95 "下载缓存目录")
-   [命令历史文件格式](/w/%E5%91%BD%E4%BB%A4%E5%8E%86%E5%8F%B2%E6%96%87%E4%BB%B6%E6%A0%BC%E5%BC%8F "命令历史文件格式")
-   [快捷栏存储格式](/w/%E5%BF%AB%E6%8D%B7%E6%A0%8F%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "快捷栏存储格式")
-   [服务器列表存储格式](/w/%E6%9C%8D%E5%8A%A1%E5%99%A8%E5%88%97%E8%A1%A8%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "服务器列表存储格式")
-   [Realms持久化数据存储格式](/w/Realms%E6%8C%81%E4%B9%85%E5%8C%96%E6%95%B0%E6%8D%AE%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "Realms持久化数据存储格式")

[服务端文件](/w/%E6%9C%8D%E5%8A%A1%E7%AB%AF%E6%96%87%E4%BB%B6%E7%9B%AE%E5%BD%95 "服务端文件目录")

-   [服务端核心文件](/w/%E6%9C%8D%E5%8A%A1%E7%AB%AF%E6%A0%B8%E5%BF%83%E6%96%87%E4%BB%B6 "服务端核心文件")
    -   [服务端数据生成器](/w/%E6%95%B0%E6%8D%AE%E7%94%9F%E6%88%90%E5%99%A8 "数据生成器")
-   [服务端配置文件格式](/w/%E6%9C%8D%E5%8A%A1%E7%AB%AF%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6%E6%A0%BC%E5%BC%8F "服务端配置文件格式")
-   [白名单存储格式](/w/%E7%99%BD%E5%90%8D%E5%8D%95%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "白名单存储格式")
-   [封禁列表存储格式](/w/%E5%B0%81%E7%A6%81%E5%88%97%E8%A1%A8%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "封禁列表存储格式")
-   [管理员列表存储格式](/w/%E7%AE%A1%E7%90%86%E5%91%98%E5%88%97%E8%A1%A8%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "管理员列表存储格式")

[存档文件](/w/Java%E7%89%88%E5%AD%98%E6%A1%A3%E6%A0%BC%E5%BC%8F "Java版存档格式")

-   [区域文件格式](/w/%E5%8C%BA%E5%9F%9F%E6%96%87%E4%BB%B6%E6%A0%BC%E5%BC%8F "区域文件格式")
-   [结构存储格式](/w/%E7%BB%93%E6%9E%84%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "结构存储格式")

存档数据

-   [存档基础数据存储格式](/w/%E5%AD%98%E6%A1%A3%E5%9F%BA%E7%A1%80%E6%95%B0%E6%8D%AE%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "存档基础数据存储格式")
-   [存档会话锁文件格式](/w/%E5%AD%98%E6%A1%A3%E4%BC%9A%E8%AF%9D%E9%94%81%E6%96%87%E4%BB%B6%E6%A0%BC%E5%BC%8F "存档会话锁文件格式")
-   [玩家数据格式](/w/%E7%8E%A9%E5%AE%B6%E6%95%B0%E6%8D%AE%E6%A0%BC%E5%BC%8F "玩家数据格式")
-   [统计存储格式](/w/%E7%BB%9F%E8%AE%A1%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "统计存储格式")
-   [进度存储格式](/w/%E8%BF%9B%E5%BA%A6%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "进度存储格式")
-   [记分板存储格式](/w/%E8%AE%B0%E5%88%86%E6%9D%BF%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "记分板存储格式")
-   [地图存储格式](/w/%E5%9C%B0%E5%9B%BE%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "地图存储格式")
-   [命令存储存储格式](/w/%E5%91%BD%E4%BB%A4%E5%AD%98%E5%82%A8%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "命令存储存储格式")
-   [随机序列存储格式](/w/%E9%9A%8F%E6%9C%BA%E5%BA%8F%E5%88%97%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "随机序列存储格式")\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]
-   [秒表时间存储格式](/w/%E7%A7%92%E8%A1%A8%E6%97%B6%E9%97%B4%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "秒表时间存储格式")
-   [自定义Boss栏存储格式](/w/%E8%87%AA%E5%AE%9A%E4%B9%89Boss%E6%A0%8F%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "自定义Boss栏存储格式")\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]
-   [游戏规则存储格式](/w/%E6%B8%B8%E6%88%8F%E8%A7%84%E5%88%99%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "游戏规则存储格式")\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]
-   [计划事件存储格式](/w/%E8%AE%A1%E5%88%92%E4%BA%8B%E4%BB%B6%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "计划事件存储格式")\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]
-   [流浪商人数据存储格式](/w/%E6%B5%81%E6%B5%AA%E5%95%86%E4%BA%BA%E6%95%B0%E6%8D%AE%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "流浪商人数据存储格式")\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]
-   [天气数据存储格式](/w/%E5%A4%A9%E6%B0%94%E6%95%B0%E6%8D%AE%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "天气数据存储格式")\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]
-   [世界生成设置存储格式](/w/%E4%B8%96%E7%95%8C%E7%94%9F%E6%88%90%E8%AE%BE%E7%BD%AE%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "世界生成设置存储格式")\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]
-   [世界时钟存储格式](/w/%E4%B8%96%E7%95%8C%E6%97%B6%E9%92%9F%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "世界时钟存储格式")\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]

维度数据

-   [区块存储格式](/w/%E5%8C%BA%E5%9D%97%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "区块存储格式")
    -   [方块实体数据格式](/w/%E6%96%B9%E5%9D%97%E5%AE%9E%E4%BD%93%E6%95%B0%E6%8D%AE%E6%A0%BC%E5%BC%8F "方块实体数据格式")
    -   [结构片段存储格式](/w/%E7%BB%93%E6%9E%84%E7%89%87%E6%AE%B5%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "结构片段存储格式")
    -   [物品格式](/w/%E7%89%A9%E5%93%81%E6%A0%BC%E5%BC%8F "物品格式")
    -   数据组件
-   [实体数据格式](/w/%E5%AE%9E%E4%BD%93%E6%95%B0%E6%8D%AE%E6%A0%BC%E5%BC%8F "实体数据格式")
    -   [生物记忆](/w/%E7%94%9F%E7%89%A9%E8%AE%B0%E5%BF%86 "生物记忆")
-   [兴趣点存储格式](/w/%E5%85%B4%E8%B6%A3%E7%82%B9%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "兴趣点存储格式")
-   [袭击存储格式](/w/%E8%A2%AD%E5%87%BB%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "袭击存储格式")
-   [随机序列存储格式](/w/%E9%9A%8F%E6%9C%BA%E5%BA%8F%E5%88%97%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "随机序列存储格式")\[失效：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]
-   [区块标签存储格式](/w/%E5%8C%BA%E5%9D%97%E6%A0%87%E7%AD%BE%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "区块标签存储格式")
-   [世界边界存储格式](/w/%E4%B8%96%E7%95%8C%E8%BE%B9%E7%95%8C%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "世界边界存储格式")
-   [末影龙战斗存储格式](/w/%E6%9C%AB%E5%BD%B1%E9%BE%99%E6%88%98%E6%96%97%E5%AD%98%E5%82%A8%E6%A0%BC%E5%BC%8F "末影龙战斗存储格式")\[新增：[JE 26.1](/w/Java%E7%89%8826.1 "Java版26.1")\]

启动器文件

-   [客户端清单文件格式](/w/%E5%AE%A2%E6%88%B7%E7%AB%AF%E6%B8%85%E5%8D%95%E6%96%87%E4%BB%B6%E6%A0%BC%E5%BC%8F "客户端清单文件格式")
-   [启动器档案文件格式](/w/%E5%90%AF%E5%8A%A8%E5%99%A8%E6%A1%A3%E6%A1%88%E6%96%87%E4%BB%B6%E6%A0%BC%E5%BC%8F "启动器档案文件格式")

已过时文件

-   [Classic世界格式](/w/Classic%E4%B8%96%E7%95%8C%E6%A0%BC%E5%BC%8F "Classic世界格式")
-   [Indev世界格式](/w/Indev%E4%B8%96%E7%95%8C%E6%A0%BC%E5%BC%8F "Indev世界格式")
-   [Alpha世界格式](/w/Alpha%E4%B8%96%E7%95%8C%E6%A0%BC%E5%BC%8F "Alpha世界格式")
-   [server\_level.dat](/w/Server_level.dat "Server level.dat")
-   [结构生成数据文件格式](/w/%E7%BB%93%E6%9E%84%E7%94%9F%E6%88%90%E6%95%B0%E6%8D%AE%E6%96%87%E4%BB%B6%E6%A0%BC%E5%BC%8F "结构生成数据文件格式")
-   [villages.dat格式](/w/Villages.dat%E6%A0%BC%E5%BC%8F "Villages.dat格式")
-   [物品格式（旧版）](/w/%E7%89%A9%E5%93%81%E6%A0%BC%E5%BC%8F/Java%E7%89%881.20.5%E5%89%8D "物品格式/Java版1.20.5前")

检索自“[https://zh.minecraft.wiki/w/数据组件?oldid=1305139](https://zh.minecraft.wiki/w/数据组件?oldid=1305139)”

[分类](/w/Special:Categories "Special:Categories")：​

-   [Java版独有特性](/w/Category:Java%E7%89%88%E7%8B%AC%E6%9C%89%E7%89%B9%E6%80%A7 "Category:Java版独有特性")

隐藏分类：​

-   [Java版即将到来/26.1](/w/Category:Java%E7%89%88%E5%8D%B3%E5%B0%86%E5%88%B0%E6%9D%A5/26.1 "Category:Java版即将到来/26.1")
-   [Java版即将移除/26.1](/w/Category:Java%E7%89%88%E5%8D%B3%E5%B0%86%E7%A7%BB%E9%99%A4/26.1 "Category:Java版即将移除/26.1")
-   [使用计算器的页面](/w/Category:%E4%BD%BF%E7%94%A8%E8%AE%A1%E7%AE%97%E5%99%A8%E7%9A%84%E9%A1%B5%E9%9D%A2 "Category:使用计算器的页面")

## 导航菜单

### 个人工具

-   [创建账号](/w/Special:CreateAccount?returnto=%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6&returntoquery=variant%3Dzh "我们推荐您创建账号并登录，但这不是强制性的")
-   [登录](/w/Special:UserLogin?returnto=%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6&returntoquery=variant%3Dzh "我们推荐您登录，但这不是强制性的​[o]")

### 命名空间

-   [页面](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6 "查看内容页面​[c]")
-   [讨论](/w/Talk:%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6 "有关内容页面的讨论​[t]")

 不转换

-   [不转换](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?variant=zh)
-   [简体](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?variant=zh-hans)
-   [繁體](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?variant=zh-hant)
-   [大陆简体](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?variant=zh-cn)
-   [香港繁體](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?variant=zh-hk)
-   [臺灣正體](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?variant=zh-tw)

### 查看

-   [阅读](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6)
-   [编辑](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?veaction=edit "编辑该页面​[v]")
-   [编辑源代码](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=edit "编辑该页面的源代码​[e]")
-   [查看历史](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=history "此页面过去的修订​[h]")

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
-   [开发版：26.1-pre-2](/w/Java%E7%89%8826.1-pre-2)
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

-   [链入页面](/w/Special:WhatLinksHere/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6 "所有链接至本页面的wiki页面列表​[j]")
-   [相关更改](/w/Special:RecentChangesLinked/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6 "链自本页的页面的最近更改​[k]")
-   [可打印版](javascript:print\(\); "本页面的可打印版本​[p]")
-   [固定链接](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?oldid=1305139 "此页面该修订版本的固定链接")
-   [页面信息](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=info "关于此页面的更多信息")
-   [特殊页面](/w/Special:SpecialPages)
-   [查看存储桶](/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?action=bucket "Bucket")

[](/hp/1773736664)

### 其他语言

-   [Deutsch](https://de.minecraft.wiki/w/Gegenstandsdaten "Gegenstandsdaten – Deutsch")
-   [English](https://minecraft.wiki/w/Data_component_format "Data component format – English")
-   [Français](https://fr.minecraft.wiki/w/Format_de_composant_de_donn%C3%A9es "Format de composant de données – français")
-   [日本語](https://ja.minecraft.wiki/w/%E3%83%87%E3%83%BC%E3%82%BF%E3%82%B3%E3%83%B3%E3%83%9D%E3%83%BC%E3%83%8D%E3%83%B3%E3%83%88 "データコンポーネント – 日本語")
-   [한국어](https://ko.minecraft.wiki/w/%EB%8D%B0%EC%9D%B4%ED%84%B0_%EA%B5%AC%EC%84%B1_%EC%9A%94%EC%86%8C_%ED%98%95%EC%8B%9D "데이터 구성 요소 형식 – 한국어")
-   [Português](https://pt.minecraft.wiki/w/Formato_de_componente_de_dado "Formato de componente de dado – português")
-   [Українська](https://uk.minecraft.wiki/w/%D0%A4%D0%BE%D1%80%D0%BC%D0%B0%D1%82_%D0%BA%D0%BE%D0%BC%D0%BF%D0%BE%D0%BD%D0%B5%D0%BD%D1%82%D0%B0_%D0%B4%D0%B0%D0%BD%D0%B8%D1%85 "Формат компонента даних – українська")

-   此页面最后编辑于2026年3月1日 (星期日) 05:42。
-   本网站内容采用[CC BY-NC-SA 3.0](https://creativecommons.org/licenses/by-nc-sa/3.0/)授权，[附加条款亦可能应用](https://meta.weirdgloop.org/w/Licensing "wgmeta:Licensing")。  
    本站并非Minecraft官方网站，与Mojang和微软亦无从属关系。

-   [隐私政策](https://weirdgloop.org/privacy)
-   [关于Minecraft Wiki](/w/Minecraft_Wiki:%E5%85%B3%E4%BA%8E)
-   [免责声明](https://meta.minecraft.wiki/w/General_disclaimer/zh)
-   [使用条款](https://weirdgloop.org/terms)
-   [联系Weird Gloop](/w/Special:Contact)
-   [移动版视图](https://zh.minecraft.wiki/w/%E6%95%B0%E6%8D%AE%E7%BB%84%E4%BB%B6?mobileaction=toggle_view_mobile&variant=zh)

-   [![CC BY-NC-SA 3.0](https://meta.weirdgloop.org/images/Creative_Commons_footer.png)](https://creativecommons.org/licenses/by-nc-sa/3.0/)
-   [![Hosted by Weird Gloop](https://meta.weirdgloop.org/images/Weird_Gloop_footer_hosted.png)](https://weirdgloop.org)