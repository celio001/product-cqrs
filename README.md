# Product CQRS

Sistema de produtos baseado em CQRS (Command Query Responsibility Segregation), eventos Kafka e leitura otimizada com Redis.

## Arquitetura

O projeto separa escrita e leitura em serviços independentes, unificados por um API Gateway:

![Arquitetura da aplicacao](docs/architecture.png)

Visão geral da arquitetura, destacando o Kong API Gateway como ponto de entrada. O gateway roteia as requisições de mutação (POST, PUT, DELETE) para o `product-command` e as de leitura (GET) para o `product-query`. A sincronização de dados é garantida pelo `worker`, que consome os tópicos do Kafka, atualiza o MongoDB e gerencia a invalidação do cache no Redis.

```mermaid
flowchart LR
		Client[Cliente] --> Gateway[API Gateway\nKong]
		Gateway -- POST, PUT, DELETE --> Command[product-command]
		Gateway -- GET --> Query[product-query]
		Command --> Postgres[(PostgreSQL)]
		Command --> Kafka[(Kafka)]
		Kafka -- product.topics\ncategory.topics\nbrand.topics --> Worker[worker]
		Worker --> Mongo[(MongoDB)]
		Worker -- Invalida cache --> Redis[(Redis com TTL)]
		Query --> Redis
		Query --> Mongo
		Redis --> Exporter[Redis Exporter]
		Exporter --> Prometheus[Prometheus]
		Prometheus --> Grafana[Grafana]
		Command -. traces/logs .-> OTel[OpenTelemetry]
		Query -. traces/logs .-> OTel
		Worker -. traces/logs .-> OTel
		OTel --> Jaeger[Jaeger]
		OTel --> Loki[(Loki)]

```

### Componentes

* API Gateway (Kong): Ponto de entrada único. Roteia tráfego para os serviços de comandos ou consultas dependendo do método HTTP.
* `product-command`: API de escrita. Persiste produtos, marcas e categorias no PostgreSQL e publica eventos no Kafka.
* `worker`: consome eventos Kafka, sincroniza a base de leitura no MongoDB e **invalida o cache** no Redis quando há alterações.
* `product-query`: API de leitura. Consulta Redis primeiro (cache-aside) e usa MongoDB como fallback, populando o cache com TTL em caso de miss.
* PostgreSQL: banco relacional otimizado para a parte de comandos (escrita).
* MongoDB: banco de documentos utilizado como modelo de leitura.
* Redis: cache de produtos com TTL de 5 minutos.
* Kafka: transporte dos eventos de domínio agrupados por entidades.
* OpenTelemetry, Prometheus, Grafana, Loki e Jaeger: stack unificada de observabilidade para métricas, logs estruturados e distributed tracing.

## Requisitos

* Docker Desktop com Docker Compose
* Go `1.26.5` apenas para desenvolvimento local

## Executando com Docker

Na raiz do projeto:

```bash
docker compose up --build

```

Para executar em segundo plano:

```bash
docker compose up -d --build

```

Para acompanhar os logs:

```bash
docker compose logs -f product-command product-query worker gateway

```

Para encerrar os containers:

```bash
docker compose down

```

Para remover também os volumes persistentes:

```bash
docker compose down -v

```

## Endpoints

Toda a comunicação externa deve ser feita através do API Gateway.

### API de comandos (via Gateway)

| Método | Rota | Descrição |
| --- | --- | --- |
| `POST` | `/product/` | Cria um produto |
| `PUT` | `/product/:id` | Atualiza um produto |
| `DELETE` | `/product/:id` | Desativa um produto |
| `POST` | `/brands/` | Cria uma marca |
| `PUT` | `/brands/:id` | Atualiza uma marca |
| `DELETE` | `/brands/:id` | Desativa uma marca |
| `POST` | `/categories/` | Cria uma categoria |
| `PUT` | `/categories/:id` | Atualiza uma categoria |
| `DELETE` | `/categories/:id` | Desativa uma categoria |

Exemplo de criação de produto via Gateway:

```bash
curl -X POST http://localhost:8000/v1/product/ \
	-H 'Content-Type: application/json' \
	-d '{
		"name": "Notebook",
		"sku": "NOTEBOOK-001",
		"unit_of_measure": "UN",
		"cost_price": 2500,
		"sale_price": 3200,
		"stock": {
			"quantity_available": 10,
			"minimum_stock": 2
		},
		"fiscal": {}
	}'

```

### API de consultas (via Gateway)

| Método | Rota | Descrição |
| --- | --- | --- |
| `GET` | `/product/:id` | Consulta um produto por UUID |

Exemplo:

```bash
curl http://localhost:8000/v1/product/<product-uuid>

```

O fluxo de leitura implementa cache-aside: em requisições GET, o `product-query` consulta o Redis. Em caso de *cache miss*, consulta o MongoDB, grava o resultado no Redis (com TTL de 5 minutos) e retorna. Modificações (PUT/DELETE/POST) disparam eventos que fazem o `worker` **invalidar** ativamente esse cache.

## Kafka

O Compose cria as seguintes famílias de tópicos, que centralizam os eventos do domínio:

* `product.topics` (incluindo created, updated, deleted)
* `brand.topics`
* `category.topics`

O worker também utiliza tópicos `*.dlq` (Dead Letter Queue) quando configurados pelas variáveis de ambiente para lidar com falhas de processamento.

## Portas e ferramentas

| Serviço | URL |
| --- | --- |
| API Gateway | `http://localhost:8000` |
| Product Command | `http://localhost:8081` (Interno) |
| Product Query | `http://localhost:8082` (Interno) |
| Kafka UI | `http://localhost:8080` |
| Grafana | `http://localhost:3000` |
| Prometheus | `http://localhost:9090` |
| Redis Insight | `http://localhost:5540` |
| Jaeger | `http://localhost:16686` |
| Loki | `http://localhost:3100` |

Credenciais padrão do Grafana:

* Usuário: `admin`
* Senha: `admin`

O Grafana provisiona automaticamente os datasources Prometheus, Loki e Jaeger via OpenTelemetry Collector.

## Variáveis principais

Os valores abaixo já possuem defaults para o ambiente Docker:

| Variável | Exemplo |
| --- | --- |
| `POSTGRES_DB_DSN` | `postgres://postgres:postgres@postgres-main:5432/product?sslmode=disable` |
| `KAFKA_BROKERS` | `kafka1:9092` |
| `REDIS_HOST` | `redis:6379` |
| `MONGO_DB_DSN` | `mongodb://root:MongoDB2019!@mongo:27017/?authSource=admin` |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | `otel-collector:4317` |

## Desenvolvimento local

Cada serviço Go é um módulo independente:

```bash
cd product-command && go test ./...
cd ../product-query && go test ./...
cd ../worker && go test ./...

```

Para formatar um módulo:

```bash
go fmt ./...

```

O `product-command` possui comandos de migration Atlas no Makefile:

```bash
cd product-command
make migrate.status
make migrate
make migrate.diff

```

## Observações

* Os serviços devem usar os nomes dos containers como hostnames quando executados no Compose (ex: `redis:6379`, `kafka1:9092`).
* A consistência do cache Redis é mantida primariamente pelo `worker`, que realiza a **invalidação da chave** no momento em que um evento de atualização ou exclusão é processado, além da segurança extra do TTL de 5 minutos.
* Toda a emissão de traces e logs estruturados é roteada padronizadamente pelo OpenTelemetry.

```

```
