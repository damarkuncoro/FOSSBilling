# Multi-stage Dockerfile for FOSSBilling Admin Portal
# --------------------------------------------------------
# Build Stage
# --------------------------------------------------------
FROM node:22-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY package*.json ./
RUN npm ci

# Copy source and build
COPY . .
RUN npm run build

# --------------------------------------------------------
# Runtime Stage (Nginx)
# --------------------------------------------------------
FROM nginx:1.27-alpine

# Security: Remove default Nginx content and use non-root if needed
# (Note: Standard nginx image runs as root to bind to port 80)

# Copy built assets
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy custom nginx config for React routing support
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
