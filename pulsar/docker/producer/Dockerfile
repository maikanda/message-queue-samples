FROM ubuntu:noble-20240605

# OSパッケージインストール
RUN \
  rm -f /etc/apt/apt.conf.d/docker-clean; \
  echo 'Binary::apt::APT::Keep-Downloaded-Packages "true";' > /etc/apt/apt.conf.d/keep-cache
RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
  export DEBIAN_FRONTEND=noninteractive && \
  apt-get update && \
  apt-get upgrade -y && \
  apt-get install -y \
    curl \
    sudo \
    ca-certificates \
    # apt-getでcorrettoをinstallする際、httpsのsourceを利用するために必要
    apt-transport-https \
    # add-apt-repositoryを利用するために必要
    software-properties-common \
    wget \
    libxrender1 libxtst6 libxi6 libx264-dev && \
  wget -O- https://apt.corretto.aws/corretto.key | apt-key add - && \
  add-apt-repository 'deb https://apt.corretto.aws stable main' && \
  apt-get update && \
  apt-get install -y java-11-amazon-corretto-jdk


# mavenインストール
RUN \
  curl -sSL https://archive.apache.org/dist/maven/maven-3/3.8.3/binaries/apache-maven-3.8.3-bin.tar.gz \
    | tar xz -C /opt && \
  ln -s /opt/apache-maven-3.8.3/bin/mvn /usr/local/bin/mvn
