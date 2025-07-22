ARG REGISTRY=docker.io
FROM --platform=$BUILDPLATFORM ${REGISTRY}/node:18 AS builder

WORKDIR /web
COPY ./VERSION .
COPY ./web .

# 清理npm缓存并安装依赖
RUN npm cache clean --force && \
    cd /web/default && rm -rf node_modules package-lock.json && npm install --legacy-peer-deps && \
    cd /web/berry && rm -rf node_modules package-lock.json && npm install --legacy-peer-deps && \
    cd /web/air && rm -rf node_modules package-lock.json && npm install --legacy-peer-deps

# 创建build目录
RUN mkdir -p /web/build

# 构建前端项目
RUN cd /web/default && npm audit fix --force || true && DISABLE_ESLINT_PLUGIN='true' REACT_APP_VERSION=$(cat ../VERSION) npm run build
RUN cd /web/berry && npm audit fix --force || true && DISABLE_ESLINT_PLUGIN='true' REACT_APP_VERSION=$(cat ../VERSION) npm run build
RUN cd /web/air && npm audit fix --force || true && DISABLE_ESLINT_PLUGIN='true' REACT_APP_VERSION=$(cat ../VERSION) npm run build

# 验证构建结果
RUN ls -la /web/build/

FROM ${REGISTRY}/golang:alpine AS builder2

RUN apk add --no-cache \
    gcc \
    musl-dev \
    sqlite-dev \
    build-base

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux

WORKDIR /build

ADD go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=builder /web/build ./web/build

RUN go build -trimpath -ldflags "-s -w -X 'github.com/songquanpeng/one-api/common.Version=$(cat VERSION)' -linkmode external -extldflags '-static'" -o one-api

FROM ${REGISTRY}/alpine:latest

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder2 /build/one-api /

EXPOSE 3000
WORKDIR /data
ENTRYPOINT ["/one-api"]