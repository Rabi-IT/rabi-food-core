FROM node:20-alpine AS development-dependencies-env
COPY . /api
WORKDIR /api
RUN npm ci

FROM node:20-alpine AS production-dependencies-env
COPY ./package.json package-lock.json /api/
WORKDIR /api
RUN npm ci --omit=dev

FROM node:20-alpine AS build-env
COPY . /api/
COPY --from=development-dependencies-env /api/node_modules /api/node_modules
WORKDIR /api
RUN npm run build

FROM node:20-alpine
COPY ./package.json package-lock.json /api/
COPY --from=production-dependencies-env /api/node_modules /api/node_modules
COPY --from=build-env /api/build /api/build
WORKDIR /api
CMD ["npm", "run", "start"]