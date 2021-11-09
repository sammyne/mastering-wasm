# 03. 指令集

## 指令集介绍

wasm 指令包含两部分信息：
- 操作码（Opcode）：指令的 ID，决定指令将执行的操作
- 操作数（Operands）：指令的参数，决定指令的执行结果

### 操作码

wasm 指令的操作码固定为一个字节（**字节码**之称由此而来），因此指令集最多只能包含 256 条指令。

wasm 规范共定义 178 条指令，按功能分为 5 大类
- 控制指令
- 参数指令
- 变量指令
- 内存指令
- 数值指令
  - 常量指令
  - 测试指令
  - 比较指令
  - 算数运算指令
  - 类型转换指令

### 助记符
- 为了方便开发者书写和理解，wasm 规范为每个操作码定义了助记符（Mnemonic）
- 两条命名规则
  - 类型前缀：数值指令的前缀通常为 `i32`、`i64`、`f32`、`f64`，例如，`i32.load`
  - 符号后缀
    - 如果整数指令的结果不受符号影响，则操作码助记符无特别后缀，例如，`i32.add`
    - 否则指令决定将整数解释为有符号（操作码助记符以 `_s` 结尾）还是无符号（操作码助记符以 `_u` 结尾），两者一般成对出现，例如 `i64.div_s`<->`i64.div_u`

### 立即数

- 操作数分为两种
  - **静态操作数**
    - 直接编码在指令里，跟在操作码后面
    - 又称为指令的**静态立即参数**（Static Immediate Arguments），简称**立即数**
  - **动态操作数**：在运行时从操作数栈获取。后续部分如无特别说明，操作数特指动态操作数
- 立即数大致分为
  - 数值（包括常量和索引）
  - 内存指令参数
  - 控制指令参数
- 内存指令：内存加载/存储系列指令需要指定内存偏移量和对齐提示
- `block` 和 `loop` 指令
  - wasm 使用 `block`、`loop` 和 `if` 这三种指令定义顺序、循环和分支结构的起点，均以 `end` 指令为终点，形成内部是嵌套的指令序列
  - `br` 系列指令可跳出 `block` 和 `if` 块，或者重新开始 `loop` 块
  - **多返回值提案接受前**，块最多只有一个结果，其类型用一个字节表示（0x7F->`i32`，0x7E->`i64`，0x7D->`f32`，0x7C->`f64`，`0x40`->`void`）
  - **多返回值提案接受后**，块类型被重新解释为 LEB128 有符号整数
    - 负数：-1、-2、-3、-4 和 -64 分别对应限制放开前的 5 种结果
    - 非负数：必须是有效的类型索引（**块类型也存在类型段**）
- `if` 指令类似 `block` 指令，只是需要额外考虑 `else` 分支
- `br_table` 指令
  - `br` 系列指令包括 4 条：`br`、`br_if`、`br_table` 和 `return`
  - `return` 没有立即数
  - `br` 和 `br_if` 指令的立即数是索引类型
  - `br_table` 的立即数包括一个跳转表和默认跳转标签，类似 Go 语言的 `switch`

### 操作数
- wasm 规范实际上定义了一台概念上的**栈式虚拟机**。绝大多数得 wasm 指令都是基于这个栈式虚拟机工作：从栈顶弹出若干个数，进行计算，然后把结果压栈
- 运行时位于栈顶并被指令操纵的数叫做指令的**动态操作数**，简称**操作数**。相应地，这个栈称为**操作数栈**。
- 由于采用了栈式虚拟机，大部分 wasm 指令（特别是数值指令）都很短，只有一个操作码，因为操作数都已隐含在栈上。

> Python 和 Ruby 等语言也用栈式虚拟机。Lua 和 Android 早期的 Dalvik 虚拟机采用的是寄存器虚拟机，其指令需要包含寄存器索引，所以寄存器虚拟机的指令一般较长。

## 指令分析
### 数值指令
- 4 条常量指令和饱和截断指令
  - `i32.const`/`i64.const` 带 `s32`/`s64` 类型的立即数，使用 LEB128 有符号编码
  - `f32.const`/`f64.const` 带 `f32`/`f64` 类型的立即数，固定占用 4/8 字节
  - `trunc_sat`（操作码 0xFC）：格式为`前缀操作码（0xFC）+ 子操作码`，带一个单字节的子操作码作为立即数
      > 书中前缀操作码为 0x0F 的说法个人觉得是不对的
- 编码格式
  ```
  i32.const: 0x41|s32
  i64.const: 0x42|s64
  f32.const: 0x43|f32
  f64.const: 0x44|f64
  trunc_sat: 0xfc|byte
  num_instr: opcode
  ```
- 样例代码参见 [01-numeric.wat](./code/01-numeric.wat)，编译并解码可得如下输出

  ```bash
  wat2wasm 01-numeric.wat 
  wasm-objdump -d 01-numeric.wasm 

  01-numeric.wasm:	file format wasm 0x1

  Code Disassembly:

  000016 func[0]:
  000017: 43 cd cc 44 41             | f32.const 0x1.89999ap+3
  00001c: 43 66 66 36 42             | f32.const 0x1.6cccccp+5
  000021: 92                         | f32.add
  000022: fc 00                      | i32.trunc_sat_f32_s
  000024: 1a                         | drop
  000025: 0b                         | end
  ```

### 变量指令

- 变量指令共 5 条
  - 3 条用于读写局部变量，立即数是局部变量索引
  - 2 条用于读写全局变量，立即数是全局变量索引
- 编码格式
  ```
   local.get:  0x20|local_idx
   local.set:  0x21|local_idx
   local.tee:  0x22|local_idx
  global.get: 0x23|global_idx
  global.set: 0x24|global_idx
  ```
- 样例代码参见 [02-variable.wat](./code/02-variable.wat)，编译并执行可得如下输出

  ```bash
  wat2wasm 02-variable.wat 
  wasm-objdump -d 02-variable.wasm 

  02-variable.wasm:	file format wasm 0x1

  Code Disassembly:

  000025 func[0]:
  000026: 23 00                      | global.get 0
  000028: 24 01                      | global.set 1
  00002a: 20 00                      | local.get 0
  00002c: 21 01                      | local.set 1
  00002e: 0b                         | end
  ```

### 内存指令
- 内存指令共 25 条
  - 14 条加载指令，用于将内存数据加载到操作数栈，有两个立即数：对齐提示和内存偏移量
  - 9 条存储指令，用于将操作数栈顶数据写回内存，有两个立即数：对齐提示和内存偏移量
  - 2 条指令用于获取和拓展内存页数，立即数是内存索引。wasm 规范目前规定模块只能导入或定义一块内存，所以内存索引只起到占位作用，必须为 0
- 编码格式

  ```
   load_instr: opcode|align|offset # align: u32, offset: u32
  store_instr: opcode|align|offset
  memory.size: 0x3f|0x00
  memory.grow: 0x40|0x00
  ```
- 样例代码参见 [03-memory.wat](./code/03-memory.wat)，编译并执行可得如下输出

  ```bash
  wat2wasm 03-memory.wat 
  wasm-objdump -d 03-memory.wasm 

  03-memory.wasm:	file format wasm 0x1

  Code Disassembly:

  00001c func[0]:
  00001d: 41 01                      | i32.const 1
  00001f: 41 02                      | i32.const 2
  000021: 28 02 64                   | i32.load 2 100
  000024: 36 02 64                   | i32.store 2 100
  000027: 3f 00                      | memory.size 0
  000029: 1a                         | drop
  00002a: 41 04                      | i32.const 4
  00002c: 40 00                      | memory.grow 0
  00002e: 1a                         | drop
  00002f: 0b                         | end
  ```

### 结构化控制指令

- 控制指令共 13 条，包括结构化控制指令、跳转指令、函数调用指令等
- 结构化控制指令有 3 条
  - `block`、`loop` 和 `if`
  - 必须和 `end` （操作码 0x0B）指令搭配，成对出现
  - 如果 `if` 指令有两条分支，则中间由 `else` 指令（操作码 0x05）分隔
  - `end` 和 `else` 也称为**伪指令**
- 编码格式

  ```
  block_instr: 0x02|block_type|instr*|0x0b
   loop_instr: 0x03|block_type|instr*|0x0b
     if_instr: 0x04|block_type|instr*|(0x05|instr*)?|0x0b
   block_type: s32
  ```
- 样例代码参见 [04-block.wat](./code/04-block.wat)，编译并执行可得如下输出

  ```bash
  wat2wasm 04-block.wat 
  wasm-objdump -d 04-block.wasm 

  04-block.wasm:	file format wasm 0x1

  Code Disassembly:

  000017 func[0]:
  000018: 02 7f                      | block i32
  00001a: 41 01                      |   i32.const 1
  00001c: 03 7f                      |   loop i32
  00001e: 41 02                      |     i32.const 2
  000020: 04 7f                      |     if i32
  000022: 41 03                      |       i32.const 3
  000024: 05                         |     else
  000025: 41 04                      |       i32.const 4
  000027: 0b                         |     end
  000028: 0b                         |   end
  000029: 1a                         |   drop
  00002a: 0b                         | end
  00002b: 0b                         | end
  ```

### 跳转指令
- 共 4 条指令如下

  指令 | 操作码| 作用 | 立即数
  ----|----|----|-------
  `br` | 0x0C | 无条件跳转 | 目标标签索引
  `br_if` | 0x0D | 有条件跳转 | 目标标签索引
  `br_table` | 0x0E | 查表跳转 | 目标标签索引和默认标签索引
  `return` | 0x0F | 直接跳出最外层循环并导致整个函数返回 | N/A

- 样例代码参见 [05-break.wat](./code/05-break.wat)，编译并运行可得如下输出

  ```bash
  wat2wasm 05-break.wat 
  wasm-objdump -d 05-break.wasm 

  05-break.wasm:	file format wasm 0x1

  Code Disassembly:

  000016 func[0]:
  000017: 02 40                      | block
  000019: 02 40                      |   block
  00001b: 02 40                      |     block
  00001d: 0c 01                      |       br 1
  00001f: 41 e4 00                   |       i32.const 100
  000022: 0d 02                      |       br_if 2
  000024: 0e 03 00 01 02 03          |       br_table 0 1 2 3
  00002a: 0f                         |       return
  00002b: 0b                         |     end
  00002c: 0b                         |   end
  00002d: 0b                         | end
  00002e: 0b                         | end
  ```

### 函数调用指令
- 两种函数调用方式
  - `call`（操作码 0x10）进行直接函数调用，函数索引由立即数指定
  - `call_indirect`（操作码 0x11）进行间接函数调用，函数签名的签名由立即数指定，到运行时才能知道具体调用的函数
- 编码格式

  ```
     call_instr: 0x10|func_idx
  call_indirect: 0x11|type_idx|0x00
  ```
- 间接函数调用指令需要查表才能完成，由第 2 个立即数指定查哪张表。目前由于模块最多只能导入或定义一张表，所以这个立即数只起到占位作用，必须为 0
- 样例代码参见 [06-call.wat](./code/06-call.wat)，编译并解码可得如下输出

  ```bash
  wat2wasm 06-call.wat 
  wasm-objdump -d 06-call.wasm 

  06-call.wasm:	file format wasm 0x1

  Code Disassembly:

  00002b func[0]:
  00002c: 10 00                      | call 0
  00002e: 41 02                      | i32.const 2
  000030: 11 01 00                   | call_indirect 1 0
  000033: 0b                         | end
  ```

## 指令解码

前面章节没有介绍指令和表达式解码逻辑，包括
- 全局项的初始化表达式
- 元素和数据项的偏移量表达式
- 代码项的字节码

其编码格式如下

```
global: global_type|init_expr
  elem: table_idx|offset_expr|vec<func_idx>
  data: mem_idx|offset_expr|vec<byte>
  code: byte_count|vec<locals>|expr
  expr: instr*|0x0b
```
