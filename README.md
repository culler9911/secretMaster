# 序列战争--诡秘之主粉丝序列QQ群插件

欢迎大家给我打赏和发电~~ https://afdian.net/@molin

欢迎大家使用文字游戏编辑器，编写自己的文字副本： https://mission-editor.now.sh/

## 更新记录

2020年3月2 全新梳理和设计的重要更新版本，5大系统，新的功能等待你探索  v2.0.0

2020年2月21 增加了红剧场门票，增加了货币对接其它插件ini功能 v1.0.7

2020年2月16 增加分群开关，指令为`序列战争开`或`序列战争关` v1.0.3

2020年2月16 感谢flak.更新好多探险文案 v1.0.2

2020年2月14 更新部分序列名称，正式更名为序列战争

2020年2月11 增加查询排行榜指定人物属性功能，增加了大量探险文案。优化了文案添加方法。

2020年2月1 增加私聊查询功能，只需指令@群号私聊即可；在热心书友flak.的建议下，新增部分探险文案。

2020年1月16 增加命运途径默认幸运加25%特质，增加灵性属性，灵性枯竭则不加经验和金镑，修复删除人物的bug

2020年1月16 整理代码分布，修复水群bug

2020年1月15 尊名系统上线，现在序列3可以设置个性化的尊名了

2020年1月14 增加购买探险卷轴功能和删除人物功能

2020年1月13 基本功能完成

2020年1月10 开始开发并完成第一版

## 介绍

序列战争是一款在QQ群中酷Q机器人使用的聊天游戏插件。故事背景是《诡秘之主》的世界观。

策划群：1028799086

游玩群：1030551041，466238445

《诡秘之主》粉丝群：731419992

开发日志（三）

提前恭喜本群第一位真神即将诞生。愿众生在祂的指引下前进。

之前的几次大版本中我们逐步添加了金镑系统，尊名设置和探险系统，甜姜茶的滋味想必大家一定很喜欢。

在现行的0.4.5版本之中，我们发现由于金镑用途单一，导致大家花费大量的金镑在探险卷轴上，所以我们会增加商店系统，增添一些骗取金镑的新物品，海量超值货物，正在东拜朗中转。

随着真神的诞生，需要有新的教会引导群里迷茫的人们，所以我们准备开放组织建设功能，此功能会逐步完善，尽情期待，成为神灵的眷者或者是主教牧首，全看你的选择。

很多非凡者反应需要一些能够体现非凡能力的设定，于是在上次更新中增加了命运途径（怪物途径）的幸运，那么接下来我们会逐步完善各途径的技能，而技能的开发也是战斗系统的先行版。

希望大家能够健康聊天，合理水群，毕竟这只是个拥有统计功能的工具，我们只不过尽量让这个工具更有趣。

待更新内容

1）商店系统

售卖商品包括：

a. 探险卷轴

b. 灵性药剂

c. 权杖

d. 仪式材料

2）教会组织系统

必须序列3以上才能创建教会，通过购买权杖，获得创建教会资格。

眷者可以加入教会享受被动增益

3）技能系统

分为主动技能和被动技能

主动技能PK时使用

被动技能可由序列获得或由教会赐予。

## 规则说明
### 1. 使用`@Yami 帮助`指令可以查看基础用法

```
帮助：回复 帮助 可显示帮助信息。
属性：回复 属性 可查询当前人物的属性信息。
途径：回复 途径 可查询途径列表。
更换：回复 更换+途径序号 可更改当前人物的非凡途径。
排行：回复 排行 可查询当前群内的非凡者排行榜。
探险：回复 探险 可主动触发每日1次的奇遇探险经历。
```

### 2. 使用`@Yami 属性`可查看自己的人物属性
```
尊名：小鱼
途径：黑夜
序列：序列9：不眠者
经验：186
金镑：198
修炼时间：5小时
战力评价：毫不足虑
```
其中尊名为QQ名字，后续可能开放自定义尊名。

序列相当于等级。50点经验一下为普通人，普通人最低，序列0最高。

经验和金镑随着群内发言条数，每条加1点。

战力评价每100点经验变更一次，达到上限后会出现`转生`字样。

查询属性和排行会消耗2点金镑。

### 3. 使用`@Yami 途径`可查看当前22条成神途径
```
1: 黑夜
2: 死神
3: 战神
4: 暗
5: 心灵
6: 风暴
7: 智慧
8: 太阳
9: 异种
10: 深渊
11: 秘密
12: 工匠
13: 愚者
14: 门
15: 时间
16: 大地
17: 月亮
18: 命运
19: 审判
20: 黑皇帝
21: 战争
22: 魔女
```

### 4. 使用`@Yami 更换16`命令可以更换到途径16大地

使用更换命令有所限制，在序列9~序列5区间无法更换。

在序列4以上只可在相近途径更换。

如果某一序列存在序列0真神，则其余该序列人员最高只能升到序列2.

普通人阶段可以任意转换序列。

### 5. 使用`@Yami 排行`命令可以查看当前群内排行
排行默认显示前31名。
### 6. 使用`@Yami 探险`命令可以执行随机奇遇探险
探险每天可以执行3次。

奖励最高100点，最低-50点。

## 待更新内容

### 1）商店系统
售卖商品包括：

a. 探险卷轴

b. 灵性药剂

c. 权杖

d. 仪式材料

### 2) 教会组织系统
必须序列3以上才能创建教会，通过购买权杖，获得创建教会资格。

眷者可以加入教会享受被动增益

### 3) 技能系统

分为主动技能和被动技能

主动技能PK时使用

被动技能可由序列获得或由教会赐予。


