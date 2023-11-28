FROM node:18-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM gcr.io/distroless/nodejs20-debian12
WORKDIR /app
COPY package.json package-lock.json ./
COPY --from=builder /app/build .
CMD ["/app/index.js"]