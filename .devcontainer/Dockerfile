FROM docker.io/library/golang:1.19-bullseye

ENV TZ=Europe \
    EDITOR=vi \
    LANG=en_US.UTF-8 \
    GO111MODULE=on

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone && \
    apt update && \
    apt install -y sudo git make curl gnupg iputils-ping dnsutils && \
    wget -qO - https://www.mongodb.org/static/pgp/server-6.0.asc | sudo apt-key add - && \
    echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu focal/mongodb-org/6.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-6.0.list && \
    apt update && \
    apt install -y mongodb-mongosh &&\
    echo 'alias ll="ls -lsah"' >> ~/.bashrc && \
    go install golang.org/x/tools/gopls@latest

RUN \
    # Base Go tools needed for VS code Go extension
    go install golang.org/x/tools/cmd/guru@latest && \
    go install golang.org/x/tools/cmd/gorename@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    # Extra tools integrating with VS code
    go install github.com/fatih/gomodifytags@latest && \
    go install github.com/cweill/gotests/...@latest &&\
    go install github.com/davidrjenni/reftools/cmd/fillstruct@latest && \
    # Terminal tools
    go install github.com/golang/mock/mockgen@latest && \
    go install github.com/vektra/mockery/v2/...@latest && \
    2>&1 && \
    rm -rf $GOPATH/pkg/* $GOPATH/src/* /root/.cache/go-build && \
    export PATH="$PATH:$(go env GOPATH)/bin"