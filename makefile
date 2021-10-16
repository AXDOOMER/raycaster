NAME = raycaster
GO = go

all: tidy build run

build:
	$(GO) build -o $(NAME) .

tidy:
	$(GO) mod tidy

run:
	./$(NAME)
  