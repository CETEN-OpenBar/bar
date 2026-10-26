FROM node:24-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM gcr.io/distroless/nodejs24-debian13
WORKDIR /app
COPY package.json package-lock.json ./
COPY --from=builder /app/build .
CMD ["/app/index.js"]