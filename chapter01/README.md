# 01. Wasm 介绍

## 简介
- 高级编程语言只能通过以下两种形式在只懂机器语言的机器上运行
  - 编译器**预先编译**（Ahead-of-Time Compilation，简写 AOT）成机器码，然后执行
  - 解析器**即时编译**（Just-in-Time Compilation，简写 JIT）执行
- 现代编译器工作原理
    ```
    +-----+                      +--------+       +--------+
    |C/C++|---------+    +------>+ASM(x86)+------>+EXE(x86)|
    +-----|         |    |       +--------+       +--------+
    +-----+         v----+       +--------+       +--------+
    |Rust +-------->+ IR +------>+ASM(ARM)+------>+EXE(ARM)|
    +-----+         +----+       +--------+       +--------+
    +-----+         ^    |       +--------+       +--------+
    |...  +---------+    +------>+ASM(...)+------>+EXE(...)|
    +-----+                      +--------+       +--------+
    ```
    - 现代编译器一般按模块化的方式设计与实现，具体分为以下模块
        - 前端：负责预处理、词法分析、语法分析、语义分析，生成便于后续处理的中间表示（Intermediate Representation，IR）
        - 中端：分析并优化 IR，例如常量折叠、死代码消除、函数内联
        - 后端：生成目标代码，把 IR 转换成平台相关的汇编语言，最终由汇编器编译为机器码

### 特点  

1. 规范  
    - [WebAssembly Core Specification]：描述平台无关的 Wasm 模块的结构和语义，任何 Wasm 实现都必须满足这些语义
    - [WebAssembly Web API]
    - [WebAssembly JavaScript Interface]
2. 模块  
    - 模块是 Wasm 程序编译、传输和加载的单位，主要有两种格式
      - 二进制格式：Wasm 模块的主要编码格式，文件后缀一般为 `.wasm`
      - 文本格式：为了方便开发者理解 Wasm 模块，或者编写小型测试代码，文件后缀一般为 `.wat`
3. 指令集
    - 栈式虚拟机和字节码
4. 验证
    - Wasm 模块包含大量类型信息，使得静态分析在代码执行前能够发现大多数问题，只有少数问题需要推迟到运行时检查

### 语义阶段

- 3 种模块表现形式
    ```
    +-------+  Emscripten
    |.c/.cpp+--------------+
    +-------+              |         wat2wasm
    +-------+  rustc      +v----+<---------------+----+
    |.rs    +------------>+.wasm|                |.wat|
    +-------+             ++-+--+--------------->+----+
    +-------+              ^ ^  |    wat2wasm
    |...    +--------------+ |  |
    +-------+                |  |decode
                            |  |
                    encode|  v
                        +-------+-+
                        |in-memory|
                        +---------+
    ```

- 语义阶段
  ```
             decode         instantiate
             +-----+        +---------+
             |     v        |         v 
  +----+  +--+--+  +--------++  +-----+--+
  |.wat|  |.wasm|  |in|memory|  |instance|
  +-+--+  +--+--+  +-+-----+-+  +--------+
    |        ^       |     ^
    +--------+       +-----+
      compile         validate
  ```

## 开发环境搭建

本项目的所有操作环境均基于这份 [Dockerfile](../docker/Dockerfile) 构建的 docker 镜像。

其中 WABT 是 Wasm 二进制工具箱，提供诸多有用工具
- WAT 汇编器 wat2wasm
- 反汇编器 wasm2wat
- 二进制格式查看工具 wasm-objdump
- 二进制格式验证工具 wasm-validate
- ...

## hello world 示例

1. 切换到 code/hello-world 目录
2. 编译 生成 wasm 二进制文件
    ```bash
    outDir=$PWD/webapp
    targetDir=$PWD/target/wasm32-unknown-unknown/release

    cargo build --release

    cp $targetDir/*.wasm $outDir
    ```
3. 编写演示用的 [index.html](./code/hello-world/webapp/index.html)
4. 在支持 docker 的宿主机运行部署这个应用
    ```bash
    workdir=$PWD/code/hello-world/webapp

    docker run --rm \
    -v $workdir:/usr/share/nginx/html \
    --name html -p 8090:80 \
    nginx:1.19.6-alpine
    ```
5. 再浏览器访问 http://localhost:8090 即可看到页面弹窗

[WebAssembly Core Specification]: https://www.w3.org/TR/2019/PR-wasm-core-1-20191001/
[WebAssembly JavaScript Interface]: https://www.w3.org/TR/wasm-js-api-1/
[WebAssembly Web API]: https://www.w3.org/TR/wasm-web-api-1/
