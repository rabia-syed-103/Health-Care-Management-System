# ── Stage 1: BUILD ───────────────────────────────────────────
# Use Node to install dependencies and build the React app
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files first (layer caching — faster rebuilds)
COPY package.json package-lock.json* ./

# Install all dependencies
RUN npm install

# Copy all frontend source code
COPY . .

# Build the production bundle (outputs to /app/dist)
RUN npm run build

# ── Stage 2: SERVE ───────────────────────────────────────────
# Use Nginx to serve the built React app — very fast and lightweight
FROM nginx:stable-alpine

# Clear default nginx page
RUN rm -rf /usr/share/nginx/html/*

# Copy built files from stage 1 into nginx's serve folder
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy our custom nginx config (handles routing + API proxy)
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]