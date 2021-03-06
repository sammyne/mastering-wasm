# 07. 函数调用（上）

## 介绍
- 按照代码的位置，可以把 wasm 函数分为**内部函数**和**外部函数**
  - 内部函数完全在 wasm 模块内定义，其字节码写在代码段
  - 外部函数从宿主环境或其他模块导入，函数类型等信息由导入段指定
- 外部函数和内部函数按照导入或定义顺序依次排列，构成模块的函数索引空间
- 对于外部函数，按照实现方式可分为**普通函数**和**本地函数**
  - 如果外部函数是从普通 wasm 模块导入，它和内部函数没有本质区别，都是 wasm 字节码
  - 否则，外部函数就是本地语言（wasm 的实现语言）实现的，被称为**本地函数**
- 两种函数调用方式
  - **直接函数调用**通过指令立即数指定函数索引直接调用函数，称为**静态函数调用**
  - **间接函数调用**借助栈顶操作数和表间接调用函数，称为**动态函数调用**
- 参数和返回值
  - wasm 函数可以返回多个值，这些值按顺序出现在栈顶（*第一个返回值在最下面*）
- 局部变量
  - 操作数的生命周期最短，可能只持续几条指令
  - 局部变量的生命周期是整个函数的执行过程。函数参数实际上也是局部变量
  - 全局变量和内存数据的生命周期最长，存活于整个模块执行期间。如果全局变量或内存是从外部导入的，它的生命周期更长，很可能跨越多个模块实例
- 调用栈和调用帧
  - 函数指令依赖操作数栈，函数也需要为参数和局部变量分配空间
  - 实现函数调用往往还需要记录一些其他信息
  - 函数调用所需数据看成一个整体，称为函数的**调用帧**
  - 每调用一个函数就需要创建一个调用帧，函数执行完毕销毁调用帧
  - 函数调用过程创建的调用帧组成函数的调用栈
  - 任意时刻，只有位于栈顶的函数调用帧是活跃的，称为**当前帧**，与之关联的函数称为**当前函数**

## 实现
- 4 个步骤
  1. 增强操作数栈，给它添加一些实用方法
      - 主要是为参数和返回值传递以及局部变量空间做准备
  2. 给虚拟机添加函数调用栈
  3. 增强虚拟机，让新版操作数栈和函数调用栈发挥作用
  4. 实现函数调用指令

### 增强操作数栈
- 把参数、局部变量和操作数放在一起有两个好处
  - 实现简单，不需要再定义一个额外的局部变量表
  - 参数传递变成了无操作
- 具体改动
  - 添加 3 个方法，用于获取操作数数量以及按索引读写操作数
  - 添加两个方法用于批量压入和弹出操作数

### 添加调用栈
- 调用帧是控制帧的一种特殊情况
    - 后续实现控制指令时会用到控制帧
- 所有函数和控制帧共享操作数栈，所以需要记录与帧对应的操作数栈起始索引
- 程序计数器记录指令执行的位置

### 增强虚拟机
- 函数栈起始地址 BP 等于当前操作数栈高度减去参数数量

### call 指令
- 指令执行时，根据被调函数的类型从栈顶弹出参数
- 指令结束后，被调函数的返回值出现在栈顶

## 局部变量指令
- 功能：读写函数的参数和局部变量
- 共 3 条汇总如下
  指令 | 操作码 |说明
  ----|-------|----
  local.get | 0x20 | 获取并压栈局部变量的值。该指令带一个立即数指定局部变量索引。指令执行后，栈顶操作数和所访问局部变量类型相同
  local.set | 0x21 | 设置局部变量的值。局部变量的索引由立即数指定，新值从栈顶弹出（必须和待修改局部变量类型相同）
  local.tee | 0x22 | 使用栈顶操作数设置局部变量的值，并将操作数留在栈顶

## 全局变量指令

共 2 条汇总如下

指令 | 操作码 |说明
----|-------|----
global.get | 0x23 | 获取并压栈全局变量的值。该指令带一个立即数指定全局变量索引。指令执行后，栈顶操作数和所访问全局变量类型相同
global.set | 0x24 | 设置全局变量的值。全局变量的索引由立即数指定，新值从栈顶弹出（必须和待修改全局变量类型相同）
