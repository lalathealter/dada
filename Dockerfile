# The base go-image
FROM golang:latest

# Set working directory
WORKDIR /
 
# Copy all files from the current directory to the app directory
COPY . .
 

# Run command as described:
# go build will build an executable file named server in the current directory
RUN go build -o server . 
 
EXPOSE 8080
# Run the server executable
CMD [ "./server" ]
