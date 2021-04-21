# 04. 文本格式

## 基本结构

- wasm 文本格式使用 `S-` 表达式描述模块
  - `S-` 表达式源自 Lisp 语言，使用大量圆括号，特别适合描述类似抽象语法树的树形结构
- 结构上，文本格式和二进制格式有以下较大的不同之处
  - 二进制格式以段（section）为单位组织数据，文本格式则以域（field）为单位组织内容。域相当于二进制段的项目，但不一定连续出现，WAT 编译器会把同类型的域收集起来，合并成二进制段
  - 二进制格式中，除自定义段以外，其他段必须按照 ID 递增的顺序排列，文本格式的域则没有这么严格的限制，但是导入域必须出现在函数域、表域、内存域和全局域之前
  - 文本格式没有单独的代码域，只有函数域。WAT 编译器会将函数域收集起来，分别生成函数段和代码段。文本格式没有自定义域
  - 文本格式提供多种内联写法

### 类型域

```
(module
  (type (func (param i32) (param i32) (result i32)))
)
```

- 圆括号是 WAT 语言的主要分隔符
- `module`、`type`、`func`、`param`、`result` 等是 WAT 的关键字
- 可为函数类型分配以 `$` 开头的**标识符**，使得其他地方可以通过标识符引用函数类型
- 多个参数可以简写在同一个 `param` 块
- 多个返回值可以简写在同一个 `result` 块

### 导入和导出域

```
(module
  (type $ft1 (func (param i32 i32) (result i32)))

  (import "env" "f1" (func $f1 (type $ft1)))
  (import "env" "t1" (table $t 1 8 funcref))
  (import "env" "m1" (memory $m 4 16))
  (import "env" "g1" (global $g1 i32))
  (import "env" "g2" (global $g2 (mut i32)))
  (import "env" )
)
```

- 导入域中，需要指明模块名、元素名和元素具体类型。其中，模块名和元素名用**字符串**指定，以双引号分隔。导入域也可以附带标识符便于后续引用
- WAT 支持两种类型的注释
  - 单行注释以 `;;` 开始
  - 多行注释以 `(;;` 开始，`;;)` 结束
- 上例的类型域是单独出现的，并在导入函数通过名字引用。当多个函数有相同类型时，这种写法可避免代码重复
- 如果某个函数类型只被使用一次，可以将其内联进导入域如下
  ```
  (module
    (import "env" "f1"
      (func $f1
        (param i32 i32) (result i32)
      )
    )
  )
  ```
- 导出域只需指定导出名和元素索引。建议通过标识符引用元素，实际索引交给 WAT 编译器计算即可。导出名在整个模块内必须唯一。样例如下
  ```
  (module
    (export "f1" (func $f1))
    (export "f2" (func $f2))
    (export "t1" (table $t))
    (export "m1" (memory $m))
    (export "g1" (global $g1))
    (export "g2" (global $g2))
  )
  ```
- 导入导出域可以内联在函数、表、内存和全局域
  - 导入域内联例如下
    ```
    (module
      (type $ft1 (func (param i32 i32) (result i32)))

      (func $f1 (import "env" "f1") (type $ft1))
      (table $t1 (import "env" "t") 1 8 funcref)
      (memory $m1 (import "env" "m") 4 16)
      (global $g1 (import "env" "g1") i32)
      (global $g2 (import "env" "g2") (mut i32))
    )
    ```
  - 导出域内联样例如下
    ```
    (module
      (func   $f (export "f1") ...)
      (table  $t (export "t" ) ...)
      (memory $m (export "m" ) ...)
      (global $g (export "g1") ...)
    )
    ```

### 函数域
- 函数域定义函数的类型和局部变量，并给出函数的指令
- WAT 编译器会把函数拆开，把类型索引放在函数段，局部变量信息和字节码放在代码段
- 样例代码
  ```
  (module
    (type $ft (func (param i32 i32) (result i32)))

    (func $add (type $ft)
      (local i64 i64)

      (i64.add (local.get 2) (local.get 3)) (drop)
      (i32.add (local.get 0) (local.get 1))
    )
  )
  ```
- 函数参数本质上也是局部变量，与函数的局部变量一起构成函数的局部变量空间，索引从 0 开始
- 可将函数类型内联进函数域，并拆分 `param` 块后为参数命名，拆分 `result` 块为结果命名，样例代码如下
  ```
  (module
    (func $add (param $a i32) (param $b i32) (result i32)
      (local $c i64)
      (local $d i64)

      (i64.add (local.get $c) (local.get $d)) (drop)
      (i32.add (local.get $a) (local.get $b))
    )
  )
  ```

### 表域和元素域

- 目前，表域最多只能出现一次，元素域可以出现多次
- 表域需要描述表的类型，包括限制和元素类型（目前只能是 `funcref`）
- 元素域可以指定若干个函数索引，以及第一个索引的表内偏移量
- 样例代码
  ```
  (module
    (func $f1)
    (func $f2)
    (func $f3)

    (table 10 20 funcref)

    (elem (offset (i32.const 5)) $f1 $f2 $f3)
  )
  ```
- 表和内存偏移量以及全局变量的初始值通过常量指定
- 表域也可以内联一个元素域，但使用这种方式无法指定表的限制和元素的表内偏移量
  ```
  (module
    (func $f1)
    (func $f2)
    (func $f3)

    (table funcref
      (elem $f1 $f2 $f3)
    )
  )
  ```

### 内存域和数据域
- 目前，内存域最多出现一次，需要描述内存的类型（即页数的上下限）
- 数据域可以出现多次，需要指定内存的偏移量和初始数据
- 样例代码
  ```
  (module
    (memory 4 16)

    (data (offset (i32.const 100)) "Hello, ")
    (data (offset (i32.const 108)) "World!\n")
  )
  ```
- 内存的数据以字符串形式指定。除普通字符外，还可以使用
  - 转义字符
  - 十六进制编码的任意字节
  - Unicode 代码点
- 表可以内联一个数据域，但内联后无法指定内存页数和偏移量。样例代码如下
  ```
  (module
    (memory                             ;; min=1, max=1
      (data "Hello, " "World!\n")       ;; inline data, offset=0
    )
  )
  ```

### 全局域
- 全局域定义全局变量，需要描述变量的类型和可变性，并给定初始值
- 全局变量也可以通过标识符引用
- 样例代码
  ```
  (module
    (global $g1 (mut i32) (i32.const 100))
    (global $g2 (mut i32) (i32.const 200))
    (global $g3 f32 (f32.const 3.14))
    (global $g4 f64 (f64.const 2.71))

    (func
      (global.get $g1)
      (global.set $g2)
    )
  )
  ```

### 起始域
- 起始域只需指定一个起始函数或索引即可
- 样例代码
  ```
  (module
    (func $main ... )

    (start $main)
  )
  ```

## 指令
- 指令写法有两种
  - 普通形式：与二进制编码格式基本一致
  - 折叠形式：完全是语法糖，便于人类编写，会被 WAT 编译器展开

### 普通形式
- 大部分指令都是操作码后跟立即数，立即数往往不能省略，必须以数值或名称的形式跟在操作码后面
  - 内存读写指令例外，可选的 `offset` 和 `align` 立即数需要显式指定（名称和数值用等号分开）
  - 结构化控制指令 `block`/`loop`/`if` 可以指定可选的参数和结果类型，必须以 `end` 结尾，其中 `if` 指令还可以用 `else` 分割成两条分支
- 样例代码
  ```
  (module
    (memory 1 2)

    (global $g1 (mut i32) (i32.const 0))

    (func $f1)
    (func $f2 (param $a i32)
      i32.const 123
      i32.load offset=100 align=4
      i32.const 456
      i32.store offset=200
      global.get $g1
      local.get $a
      i32.add
      call $f1
      drop
    )

    (func $foo
      block $l1 (result i32)
        i32.const 123
        br $l1
        loop $l2
          i32.const 123
          br_if $l2
        end
      end
      drop
    )

    (func $max (param $a i32) (param $b i32) (result i32)
      local.get $a
      local.get $b
      i32.gt_s
      if (result i32)
        local.get $a
      else
        local.get $b
      end
    )
  )
  ```

### 折叠形式
- 三步将普通形式调整为折叠形式
  1. 用圆括号把指令包起来
  2. 如果是结构化指令，把 `end` 去掉
  3. 这一步可选。如果某条指令（无论普通还是已折叠）和它前面的几条指令从逻辑上可以看成一组操作，则把几条指令折叠进该指令
- 前一小节样例代码的折叠形式如下

```
(module
  (func $foo
    (block $l1 (result i32)
      (i32.const 123)
      (br $l1)
      (loop $l2
        (br_if $l2 (i32.const 123))
      )
    )
    (drop)
  )

  (func $max (param $a i32) (param $b i32) (result i32)
    (if (result i32)
      (i32.get_s (local.get $a) (local.get $b))
      (then (local.get $a))
      (else (local.get $b))
    )
  )
)
```