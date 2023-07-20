FROM node:18-alpine AS external-website
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build
RUN rm -rf src/ static/ emailTemplates/ docker-compose.yml
USER node:node
CMD ["node","build/index.js"]