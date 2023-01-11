FROM golang as gobuilder
WORKDIR /dist
WORKDIR /backendServer
COPY /backendServer /backendServer/
RUN go build -mod=vendor -o ../dist/server ./server

FROM node as vitebuilder
WORKDIR /frontend
COPY /frontend /frontend/
RUN npm install
RUN npm run build

FROM ubuntu as final
COPY --from=gobuilder /dist/server /dist/server
COPY --from=vitebuilder /dist/static /dist/static

EXPOSE 8080

CMD ["/dist/server"]

