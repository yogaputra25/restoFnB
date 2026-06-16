FROM node:25-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
ARG NUXT_PUBLIC_API_BASE=http://backend:8080
ENV NUXT_PUBLIC_API_BASE=$NUXT_PUBLIC_API_BASE
RUN npm run build

FROM node:25-alpine AS runner
WORKDIR /app
COPY --from=builder /app/.output ./.output
EXPOSE 3000
CMD ["node", ".output/server/index.mjs"]
