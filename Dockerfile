# Gebruik de officiële Golang Docker image als basis
FROM golang:1.16-alpine

# Stel de map in op de locatie van de Go-code binnen de container
WORKDIR /app

# Kopieer de Go-modulebestanden en download de afhankelijkheden
COPY go.mod go.sum ./
RUN go mod download

# Kopieer de broncode naar de container
COPY . .

# Bouw de Go-applicatie binnen de container
RUN go build -o main .

# Expose poort 8080, waar de applicatie op draait
EXPOSE 8080

# Start de applicatie wanneer de container wordt gerund
CMD ["./main"]
