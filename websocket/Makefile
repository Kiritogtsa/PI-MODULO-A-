
# Nome da imagem e do container
IMAGE_NAME=minha-imagem
CONTAINER_NAME=meu-servidor
PORT=8080

# Build da imagem Docker
build:
	docker build -t $(IMAGE_NAME) .

# Rodar o container com nome fixo
run:
	docker run -d -p $(PORT):$(PORT) --name $(CONTAINER_NAME) $(IMAGE_NAME)

# Parar o container
start:
	docker start $(CONTAINER_NAME) || true


# Parar o container
stop:
	docker stop $(CONTAINER_NAME) || true

# Remover o container (caso precise recriar)
remove:
	docker rm $(CONTAINER_NAME) || true

# Executar um shell dentro do container
shell:
	docker exec -it $(CONTAINER_NAME) sh

# Exibir logs do container
logs:
	docker logs -f $(CONTAINER_NAME)

# Parar, remover e rodar novamente
restart: stop remove run

# Limpar imagens antigas
clean:
	docker system prune -f

.PHONY: build run stop remove shell logs restart clean
