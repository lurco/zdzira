FROM node:22-alpine AS frontend
WORKDIR /app
COPY frontend/package*.json frontend/
RUN cd frontend && npm ci
COPY frontend/ frontend/
RUN cd frontend && npm run build
# vite outDir is ../dist relative to frontend/, so output lands at /app/dist

FROM nginx:alpine
COPY --from=frontend /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
