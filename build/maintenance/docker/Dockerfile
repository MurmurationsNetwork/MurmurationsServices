# Use the Nginx image from Docker Hub
FROM nginx:alpine

# Remove the default Nginx configuration file
RUN rm /etc/nginx/conf.d/default.conf

# Copy the custom Nginx configuration file into the container
COPY default.conf /etc/nginx/conf.d/

# Copy the JSON file into Nginx's serve directory
COPY message.json /usr/share/nginx/html/message.json

# Expose port 8080 to the Docker host
EXPOSE 8080
