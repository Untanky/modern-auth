FROM node:20-alpine as build

RUN npm install -g pnpm

WORKDIR /build
COPY ./package.json .
COPY ./pnpm-lock.yaml .

RUN pnpm install

COPY . .

RUN pnpm build

FROM node:20-alpine as server

RUN apk update && apk upgrade

WORKDIR /app
COPY --from=build /build/.svelte-kit/output .

EXPOSE 5173

ENTRYPOINT [ "node server/index.js" ]
