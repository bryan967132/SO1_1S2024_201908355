FROM node:alpine
WORKDIR /tarea1Front
COPY package*.json .
RUN npm i
COPY . .
EXPOSE 3000
CMD ["npm", "start"]