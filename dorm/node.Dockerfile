FROM node:14.1.0

RUN apt-get update && apt-get install -y git 

WORKDIR /usr/src/app

ENV PATH /usr/src/app/node_modules/.bin:$PATH

RUN git clone --branch develop --single-branch --depth 1 https://github.com/pavel232/dorm.git \
  && chown -R www-data:www-data dorm

WORKDIR /usr/src/app/dorm 

RUN npm install 
RUN npm install -g @angular/cli@7.3.9

COPY . /usr/src/app/dorm

EXPOSE 4200

CMD ng serve --host 0.0.0.0 