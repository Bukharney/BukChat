# Stage 1: Build the static assets
FROM --platform=$BUILDPLATFORM node:22-alpine AS builder

WORKDIR /app

RUN npm install -g pnpm@9

COPY package.json pnpm-lock.yaml* ./
RUN pnpm install

COPY . .
RUN pnpm build

# Stage 2: Serve static files with Nginx (No Node/Rollup needed)
FROM nginx:alpine

# Copy the built dist directory
COPY --from=builder /app/dist /usr/share/nginx/html

# Configure Nginx to listen on port 4173 and handle SPA routing
RUN printf 'server {\n\
    listen 4173;\n\
    server_name localhost;\n\
    location / {\n\
    root /usr/share/nginx/html;\n\
    index index.html index.htm;\n\
    try_files $uri $uri/ /index.html;\n\
    }\n\
    }\n' > /etc/nginx/conf.d/default.conf

EXPOSE 4173

CMD ["nginx", "-g", "daemon off;"]