# 05. 操作数栈

wasm 程序的执行环境时一台栈式虚拟机。绝大多数 wasm 指令都要借助操作数栈来工作：从上面弹出若干数值，对数值进行计算，然后再把计算结果压栈。

## 操作数栈
- 对于计算机来说，一切信息终究只是 0 和 1 组成的序列。如果把数据的单位放大一点，一切信息都只是字节序列而已。一串 0 和 1 代表的含义取决于我们如何解释它。

## 虚拟机
### 指令循环
- 虚拟机最核心的逻辑如下
  ```
  loop 还有更多指令需要执行 {
    取出一条指令
    执行这条指令
  }
  ```

### 指令分派
- 解释器执行指令的方式有两种
  - 借助 `switch-case` 语句分派指令执行逻辑
  - 查表分配指令执行逻辑
- 本项目采用查表法

## 参数指令
### drop
- `drop` 指令（操作码 0x1A）从栈顶弹出一个操作数并把它扔掉

### select
- `select` 指令（操作码 0x1B）从栈顶弹出 3 个操作数，然后根据最先弹出的操作数从其他两个操作数选择一个压栈
- 最先弹出的操作数必须是 `i32` 类型，其他两个操作数类型相同即可
- 若果最先弹出的操作数不为 0，则把最后弹出的操作数压栈，否则把中间的操作数压栈

## 数值指令
### 常量
- 常量指令共 4 条如下
  指令 | 操作码
  ----:|:-----
  `i32.const` | 0x41
  `i64.const` | 0x42
  `f32.const` | 0x43
  `f64.const` | 0x44
- 常量指令带一个相应类型的立即数，效果是将立即数压栈

> 除常量指令以外，其余数值指令都没有立即数

### 测试
- 测试指令从栈顶弹出一个操作数，先测试它是否为 0，然后把测试结果（`i32` 类型的布尔值）压栈
- 测试指令只有两条如下
  指令 | 操作码
  ----:|:----
  `i32.eqz` | 0x45
  `i64.eqz` | 0x50

### 比较
- 比较指令从栈顶弹出两个同类型的操作数进行比较，然后把比较结果（`i32` 类型的布尔值）压栈
- 32 条比较指令汇总如下

<table border="1px">
    <thead>
      <tr>
        <td>比较</td>
        <td>i32</td>
        <td>i64</td>
        <td>f32</td>
        <td>f64</td>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td>等于</td>
        <td>i32.eq(0x46)</td>
        <td>i64.eq(0x51)</td>
        <td>f32.eq(0x5B)</td>
        <td>f64.eq(0x61)</td>
      </tr>
      <tr>
        <td>不等于</td>
        <td>i32.ne(0x47)</td>
        <td>i64.ne(0x52)</td>
        <td>f32.ne(0x5C)</td>
        <td>f64.ne(0x62)</td>
      </tr>
      <tr>
        <td rowspan="2">小于</td>
        <td>i32.lt_s(0x48)</td>
        <td>i64.lt_s(0x53)</td>
        <td rowspan="2">f32.lt(0x5D)</td>
        <td rowspan="2">f64.lt(0x63)</td>
      </tr>
      <tr>
        <td>i32.lt_u(0x49)</td>
        <td>i64.lt_u(0x54)</td>
      </tr>
      <tr>
        <td rowspan="2">大于</td>
        <td>i32.gt_s(0x4A)</td>
        <td>i64.gt_s(0x55)</td>
        <td rowspan="2">f32.gt(0x5E)</td>
        <td rowspan="2">f64.gt(0x64)</td>
      </tr>
      <tr>
        <td>i32.gt_u(0x4B)</td>
        <td>i64.gt_u(0x56)</td>
      </tr>
      <tr>
        <td rowspan="2">小于等于</td>
        <td>i32.le_s(0x4C)</td>
        <td>i64.le_s(0x57)</td>
        <td rowspan="2">f32.le(0x5F)</td>
        <td rowspan="2">f64.le(0x65)</td>
      </tr>
      <tr>
        <td>i32.le_u(0x4D)</td>
        <td>i64.le_u(0x58)</td>
      </tr>
      <tr>
        <td rowspan="2">大于等于</td>
        <td>i32.ge_s(0x4E)</td>
        <td>i64.ge_s(0x59)</td>
        <td rowspan="2">f32.ge(0x60)</td>
        <td rowspan="2">f64.ge(0x66)</td>
      </tr>
      <tr>
        <td>i32.ge_u(0x4F)</td>
        <td>i64.ge_u(0x5A)</td>
      </tr>
    </tbody>
  </table>

### 一元算术
- 一元算数指令从栈顶弹出一个操作数进行计算，然后将同类型的结果压栈
- 整数一元算数指令汇总如下
  运算 | i32 | i64
  ----|------|-----
  统计前置 0 的比特数 | i32.clz(0x67) | i64.clz(0x79)
  统计后置 0 的比特数 | i32.ctz(0x68) | i64.ctz(0x7A)
  统计前置 1 的比特数 | i32.popcnt(0x69) | i64.ctz(0x7B)
- 浮点数一元算数指令汇总如下
  运算 | f32 | f64
  ----|------|-----
  绝对值 | f32.abs(0x8B) | f64.abs(0x99)
  取反 | f32.neg(0x8C) | f64.neg(0x9A)
  向上取整 | f32.ceil(0x8D) | f64.ceil(0x9B)
  向下取整 | f32.floor(0x8E) | f64.floor(0x9C)
  截断取整 | f32.trunc(0x8F) | f64.trunc(0x9D)
  就近取整 | f32.nearest(0x90) | f64.nearest(0x9E)
  平方根 | f32.sqrt(0x91) | f64.sqrt(0x9F)

### 二元算术
- 二元算术指令从栈顶弹出两个相同类型的操作数进行计算，然后将同类型结果压栈
- 整数二元运算指令汇总如下
  运算 | i32 | i64
  ----|------|-----
  加 | i32.add(0x6A) | i64.add(0x7C)
  减 | i32.sub(0x6B) | i64.sub(0x7D)
  乘 | i32.mul(0x6C) | i64.mul(0x7E)
  有符号除 | i32.div_s(0x6D) | i64.div_s(0x7F)
  有符号除 | i32.div_u(0x6E) | i64.div_u(0x80)
  有符号求余 | i32.rem_s(0x6F) | i64.rem_s(0x81)
  无符号求余 | i32.rem_u(0x70) | i64.rem_u(0x82)
  按位与 | i32.and(0x71) | i64.and(0x83)
  按位或 | i32.or(0x72) | i64.or(0x84)
  按位异或 | i32.xor(0x73) | i64.xor(0x85)
  左移 | i32.shl(0x74) | i64.shl(0x86)
  有符号右移 | i32.shr_s(0x75) | i64.shr_s(0x87)
  无符号右移 | i32.shr_u(0x76) | i64.shr_u(0x88)
  左旋转 | i32.rotl(0x77) | i64.rotl(0x89)
  右旋转 | i32.rotr(0x78) | i64.rotr(0x8A)

- 浮点数二元运算指令汇总如下
  运算 | f32 | f64
  ----|-----|----
  加 | f32.add(0x92) | f64.add(0xA0)
  减 | f32.sub(0x93) | f64.sub(0xA1)
  乘 | f32.mul(0x94) | f64.mul(0xA2)
  除 | f32.div(0x95) | f64.div(0xA3)
  取最小值 | f32.min(0x96) | f64.min(0xA4)
  取最大值 | f32.max(0x97) | f64.max(0xA5)
  拷贝符号位 | f32.copysign(0x98) | f64.copysign(0xA6)

  - `copysign(v1, v2)` 指令将 `v1` 的绝对值和 `v2` 的符号相结合
  - 整数 `div` 有可能导致除零错误

### 类型转换
- 类型转换指令从栈顶弹出一个操作数进行类型转换，然后把结果压栈
- 相关指令汇总如下
  <table border="1px">
    <thead>
      <tr>
        <td></td>
        <td>转换后</td>
        <td>i32</td>
        <td>i64</td>
        <td>f32</td>
        <td>f64</td>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td rowspan="12">转换前</td>
        <td rowspan="3">i32</td>
        <td>i32.extend8_s</td>
        <td>i64.extend_i32_s</td>
        <td>f32.convert_i32_s</td>
        <td>f64.convert_i32_s</td>
      </tr>
      <tr>
        <td>i32.extend16_s</td>
        <td>i64.extend_i32_u</td>
        <td>f32.convert_i32_u</td>
        <td>f64.convert_i32_u</td>
      </tr>
      <tr>
        <td></td>
        <td></td>
        <td>f32.reinterpret_i32</td>
        <td></td>
      </tr>
      <tr>
        <td rowspan="3">i64</td>
        <td>i32.wrap_i64</td>
        <td>i64.extend8_s</td>
        <td>f32.convert_i64_s</td>
        <td>f64.convert_i64_s</td>
      </tr>
      <tr>
        <td rowspan="2"></td>
        <td>i64.extend16_s</td>
        <td>f32.convert_i64_u</td>
        <td>f64.convert_i64_u</td>
      </tr>
      <tr>
        <td>i64.extend32_s</td>
        <td></td>
        <td>f64.reinterpret_i64</td>
      </tr>
      <tr>
        <td rowspan="3">f32</td>
        <td>i32.trunc_f32_s</td>
        <td>i64.trunc_f32_s</td>
        <td rowspan="3"></td>
        <td>f64.promote_f32</td>
      </tr>
      <tr>
        <td>i32.trunc_f32_u</td>
        <td>i64.trunc_f32_u</td>
        <td rowspan="2"></td>
      </tr>
      <tr>
        <td>i32.reinterpret_f32</td>
      </tr>
      <tr>
        <td rowspan="3">f64</td>
        <td>i32.trunc_f64_s</td>
        <td>i64.trunc_f64_s</td>
        <td>f32.demote_f64</td>
        <td rowspan="3"></td>
      </tr>
      <tr>
        <td>i32.trunc_f64_s</td>
        <td>i64.trunc_f64_s</td>
        <td rowspan="2"></td>
      </tr>
      <tr>
        <td></td>
        <td>i64.reinterpret_f64</td>
      </tr>
    </tbody>
  </table>

- 两个某种类型的数的计算结果可能会超出该类型的表达范围，触发上溢（overflow）或下溢（underflow）
- 溢出的处理方式有 3 种如下
  - 环绕，整数运算通常采用这种方式
  - 饱和，浮点数运算通常采用这种方式，超出范围的值会被表示为正或负“无穷”
  - 异常，例如整数除零异常
- 为了避免异常情况，某提案建议增加了 8 条饱和截断指令进行特殊处理，比如将 `NaN` 转换为 0，将正/负无穷转换为整数最大/小值。饱和截断指令通过操作码前缀 0xFC 引入如下
  指令 | 操作码
  ----:|:----
  i32.trunc_sat_f32_s  | 0xFC 0x00
  i32.trunc_sat_f32_ u | 0xFC 0x01
  i32.trunc_sat_f64_s  | 0xFC 0x02
  i32.trunc_sat_f64_u  | 0xFC 0x03
  i64.trunc_sat_f32_s  | 0xFC 0x04
  i64.trunc_sat_f32_u  | 0xFC 0x05
  i64.trunc_sat_f64_s  | 0xFC 0x06
  i64.trunc_sat_f64_u  | 0xFC 0x07