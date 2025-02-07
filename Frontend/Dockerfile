FROM node:18

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

#ARG REACT_APP_BACKEND_URL
#ENV REACT_APP_BACKEND_URL=${REACT_APP_BACKEND_URL}
RUN echo ${REACT_APP_BACKEND_URL}

RUN npm run build

FROM nginx:alpine
RUN apk add --no-cache gettext
COPY --from=0 /app/build /usr/share/nginx/html
COPY nginx.conf.template /etc/nginx/templates/nginx.conf.template
CMD envsubst '${REACT_APP_BACKEND_URL}' < /etc/nginx/templates/nginx.conf.template > /etc/nginx/conf.d/default.conf && nginx -g 'daemon off;'

EXPOSE 80