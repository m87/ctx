# syntax=docker/dockerfile:1.6

########################
# BUILD STAGE
########################
FROM --platform=$BUILDPLATFORM node:22-alpine AS build

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .

RUN npm run build -- --configuration production

########################
# RUNTIME STAGE
########################
FROM --platform=$TARGETPLATFORM docker.io/library/nginx:1.29.4-alpine

RUN rm /etc/nginx/conf.d/default.conf

COPY nginx.conf /etc/nginx/conf.d/default.conf

COPY --from=build /app/dist/ctx/browser /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
