FROM golang:1.12-alpine

# Add Maintainer Info
LABEL maintainer="Nico Schaefer <nschaefer@cs.uni-kl.de>"

########## Build Container ##########

# Set the Current Working Directory inside the container
WORKDIR /app

# Add git
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY cmd/ cmd/
COPY internal/ internal/
COPY web/ web/

# Build the Go app
RUN go build cmd/joda-web/joda-web.go && mv joda-web cmd/joda-web/joda-web
WORKDIR /app/cmd/joda-web/


########## Application Container ##########
FROM alpine:latest  
COPY --from=0 /app/cmd/ /app/cmd
COPY --from=0 /app/web/ /app/web
WORKDIR /app/cmd/joda-web/

# Expose port 8080 to the outside world
EXPOSE 8080
# Run the executable
ENTRYPOINT ["./joda-web", "-host"]
