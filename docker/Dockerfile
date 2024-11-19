# 编译前端页面
# FROM nodejs as nodejs-builder


# =============================================================================
# 编译后端API

FROM golang:1.21.4-alpine as go-builder

ENV GOPROXY=https://goproxy.cn,direct

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update && apk upgrade && \
    apk add --no-cache ca-certificates git mercurial

# Build version and commit should be passed in when performing docker build
ARG PROJECT_NAME=myadmin
ARG VERSION=0.0.1
ARG GIT_SHA=0000000

WORKDIR /myadmin

# ARG http_proxy=http://proxy.cn:1080
# ARG https_proxy=http://proxy.cn:1080

COPY go.mod go.sum ./
RUN go mod download

COPY pkg ./ cmd ./ version ./
COPY --from=nodejs-builder /myadmin/dist /myadmin/web/dist/

RUN go mod download

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${GOBIN}/${PROJECT_NAME} \
    -ldflags "-X myadmin/internal/version.Version=${VERSION} -X myadmin/internal/version.GitSHA=${GIT_SHA}" \
    /myadmin

# =============================================================================
# 将编译文件打包

FROM alpine:3.18.4 AS final

ARG PROJECT_NAME=myadmin

COPY --from=go-builder ${GOBIN}/${PROJECT_NAME} /usr/local/bin/${PROJECT_NAME}

RUN adduser -D ${PROJECT_NAME}
USER ${PROJECT_NAME}

ENTRYPOINT ["/usr/local/bin/myadmin"]

