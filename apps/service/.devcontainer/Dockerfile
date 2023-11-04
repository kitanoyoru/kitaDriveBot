FROM ubuntu:22.10

ARG NODE_VERSION="18.12.1"
ARG GO_VERSION="1.21.0"

ARG TARGETARCH
ARG USERNAME=kitanoyoru
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN groupadd --gid $USER_GID $USERNAME && \
    useradd --uid $USER_UID --gid $USER_GID -m $USERNAME && \
    usermod -s /bin/bash $USERNAME

RUN apt update && \
    apt install -y build-essential git protobuf-compiler curl wget sudo
RUN wget https://go.dev/dl/go$GO_VERSION.linux-$TARGETARCH.tar.gz && \
    tar -C /usr/local -xzf go$GO_VERSION.linux-$TARGETARCH.tar.gz && \
    rm go$GO_VERSION.linux-$TARGETARCH.tar.gz

RUN echo "${USERNAME} ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers

RUN echo "${USERNAME}:kitanoyoru" | chpasswd
USER ${USERNAME}

ENV PATH="$PATH:/usr/local/go/bin:$PATH:/home/$USERNAME/go/bin"
