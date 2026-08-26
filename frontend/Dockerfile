# Use Node.js 22 LTS (Alpine or Slim)
FROM node:22-alpine

WORKDIR /app

# Enable corepack (built-in pnpm manager) or install pnpm globally
RUN npm install -g pnpm

# Copy package manifests first for efficient Docker layer caching
COPY package.json pnpm-lock.yaml* ./

RUN pnpm install

COPY . .

RUN pnpm build
EXPOSE 4173

CMD ["pnpm", "preview"]