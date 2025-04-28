FROM node:22.15-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm install
COPY . .
RUN npm run build

FROM node:22.15-alpine
WORKDIR /app
# Update packages to fix vulnerabilities
RUN apk update && apk upgrade
# Create a non-root user
RUN addgroup -g 1001 nodejs && adduser -u 1001 -G nodejs -s /bin/sh -D nextjs
COPY --from=builder /app/.next .next
COPY --from=builder /app/package.json package.json
COPY --from=builder /app/public public
COPY --from=builder /app/node_modules node_modules
# Set proper ownership
RUN chown -R nextjs:nodejs /app
USER nextjs
EXPOSE 3000
CMD ["npm", "start"]