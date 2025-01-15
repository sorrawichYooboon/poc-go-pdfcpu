FROM golang:1.23

# Set the working directory
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

RUN apt-get -y update && apt-get -y upgrade && apt-get -y install chromium

# Expose port 8080 to the outside world
EXPOSE 8080

CMD ["./app"]

