# Use Node.js v14 as base image
FROM node:19.9.0

# Set working directory to /app
WORKDIR /app

# Copy package.json and package-lock.json to the container
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy remaining source code to the container
COPY . .

# Compile TypeScript
RUN npm run build

# Expose port 3000
EXPOSE 3000

# Set the command to start the bot
CMD [ "npm", "start" ]