# Multi-stage Dockerfile for FOSSBilling Documentation Site (VitePress)
FROM node:20-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app
COPY package*.json ./
RUN npm install

COPY . .
RUN npm run docs:build

# Production Web Server
FROM nginx:alpine
COPY --from=builder /app/.vitepress/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
