# base image
FROM golang:1-alpine as stage1
# copy code
COPY backend /codebase/backend
# RUN ls /codebase/src/main.go
# build binary
RUN cd /codebase/backend && go build -v -o bin/dgframework main.go

# decrease size
FROM alpine as stage2
# copy final binary
COPY --from=stage1 /codebase/backend/bin/dgframework /dgframework
# ENV PORT=8080
# run command for the binary
ENTRYPOINT [ "/dgframework" ]