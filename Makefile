COMPOSE ?= docker compose
SERVICE ?=

.PHONY: help up up-d build down logs ps clean

help:
	@printf '%s\n' \
	  'Comandos disponíveis:' \
	  '  make up       Build e inicia os serviços no terminal' \
	  '  make up-d     Build e inicia os serviços em segundo plano' \
	  '  make build    Build das imagens' \
	  '  make down     Para e remove os containers' \
	  '  make logs     Acompanha os logs (SERVICE=nome opcional)' \
	  '  make ps       Mostra o status dos serviços' \
	  '  make clean    Para os serviços e remove os volumes'

up:
	$(COMPOSE) up --build

up-d:
	$(COMPOSE) up -d --build

build:
	$(COMPOSE) build

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f $(SERVICE)

ps:
	$(COMPOSE) ps

clean:
	$(COMPOSE) down -v