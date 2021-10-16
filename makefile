NAME = raycaster

all:
	go build -o ${NAME} .

run:
	./${NAME}
