# react
FROM node:20.11.0-alpine as build
WORKDIR /frontend
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# nginx server
FROM nginx
COPY --from=build /frontend/build /usr/share/nginx/html
COPY --from=build /frontend/default.conf /etc/nginx/conf.d/
EXPOSE 80:80
CMD ["nginx", "-g", "daemon off;"]