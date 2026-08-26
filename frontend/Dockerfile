# 1. Build stage runs natively on the builder's host architecture (x86_64) - NO QEMU
FROM --platform=$BUILDPLATFORM node:22-alpine AS builder

WORKDIR /app

RUN npm install -g pnpm@9

COPY package.json pnpm-lock.yaml* ./
RUN pnpm install

COPY . .
RUN pnpm build

# 2. Production stage (Target architecture: linux/arm64)
FROM node:22-alpine

WORKDIR /app

# Copy built files from builder stage
COPY --from=builder /app/package.json ./
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/dist ./dist

EXPOSE 4173

# Adjust command to your preview/start script
CMD ["pnpm", "preview", "--host"]