FROM node:18-alpine

# Create app directory
WORKDIR /app

# Install app dependencies
COPY package*.json ./
RUN npm install --only=production

# Bundle app source
COPY . .

# Expose port
EXPOSE 3000

# Start the application
CMD ["node", "src/index.js"]