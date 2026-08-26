# 1. Builder stage
FROM --platform=$BUILDPLATFORM node:22-alpine AS builder

WORKDIR /app

RUN npm install -g pnpm@9

COPY package.json pnpm-lock.yaml* ./
RUN pnpm install

COPY . .
RUN pnpm build

# 2. Production stage
FROM node:22-alpine

WORKDIR /app

# Install pnpm in the production stage as well
RUN npm install -g pnpm@9

COPY --from=builder /app/package.json ./
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/dist ./dist

EXPOSE 4173

# Run pnpm directly (binding to 0.0.0.0)
CMD ["pnpm", "preview", "--host", "0.0.0.0", "--port", "4173"]