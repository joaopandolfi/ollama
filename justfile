# Lints Go code via golangci-lint within Docker
lint:
  docker run \
    -t \
    --rm \
    -v "$(pwd)/:/app" \
    -w /app \
    golangci/golangci-lint:v1.60 \
    golangci-lint run -v
