FROM rust:1.51.0-buster 

ENV DEBIAN_FRONTEND=noninteractive \
  RUSTUP_UPDATE_ROOT=https://mirrors.ustc.edu.cn/rust-static/rustup \
  RUSTUP_DIST_SERVER=https://mirrors.ustc.edu.cn/rust-static

RUN export GO_VERSION=1.16.3                                          &&\
  cd /root                                                            &&\
  wget "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz"  &&\
  tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz           &&\
  rm -f go${GO_VERSION}.linux-amd64.tar.gz                            &&\
  ln -sf /usr/local/go/bin/* /usr/bin

RUN curl -L https://hub.fastgit.org/WebAssembly/wabt/releases/download/1.0.23/wabt-1.0.23-ubuntu.tar.gz |\
  tar -zxv -C /opt/                                           &&\
  echo 'export PATH=/opt/wabt-1.0.23/bin:$PATH' >> ~/.bashrc

# 这一步成功执行很艰难，建议注释掉，进入容器后手动执行
RUN rustup target add wasm32-unknown-unknown
