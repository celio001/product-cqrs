# Product CQRS

Sistema de produtos baseado em CQRS (Command Query Responsibility Segregation), eventos Kafka e leitura otimizada com Redis.

## Arquitetura

O projeto separa escrita e leitura em servicos independentes:

```mermaid
flowchart LR
		Client[Cliente] --> Command[product-command\nHTTP :8081]
		Command --> Postgres[(PostgreSQL)]
		Command --> Kafka[(Kafka)]
		Kafka --> Worker[worker]
		Worker --> Mongo[(MongoDB)]
		Worker --> Redis[(Redis)]
		Client --> Query[product-query\nHTTP :8082]
		Query --> Redis
		Query --> Mongo
		Redis --> Exporter[Redis Exporter]
		Exporter --> Prometheus[Prometheus]
		Prometheus --> Grafana[Grafana]
		Command -. traces .-> Jaeger[Jaeger]
		Query -. traces/logs .-> Observability[Observabilidade]
		Worker -. traces/logs .-> Observability
```

### Componentes

- `product-command`: API de escrita. Persiste produtos, marcas e categorias no PostgreSQL e publica eventos no Kafka.
- `worker`: consome eventos Kafka, sincroniza produtos no MongoDB e popula o cache Redis.
- `product-query`: API de leitura. Consulta Redis primeiro e usa MongoDB como fallback.
- PostgreSQL: banco da parte de comandos.
- MongoDB: modelo de leitura.
- Redis: cache de produtos com TTL de 5 minutos.
- Kafka: transporte dos eventos de domínio.
- Prometheus, Grafana, Loki, Promtail e Jaeger: monitoramento, logs e tracing.

## Requisitos

- Docker Desktop com Docker Compose
- Go `1.26.5` apenas para desenvolvimento local

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
docker compose logs -f product-command product-query worker
```

Para encerrar os containers:

```bash
docker compose down
```

Para remover tambem os volumes persistentes:

```bash
docker compose down -v
```

## Endpoints

### API de comandos

Base URL: `http://localhost:8081/v1`

| Metodo | Rota | Descricao |
| --- | --- | --- |
| `POST` | `/product/` | Cria um produto |
| `DELETE` | `/product/:id` | Desativa um produto |
| `POST` | `/brands/` | Cria uma marca |
| `DELETE` | `/brands/:id` | Desativa uma marca |
| `POST` | `/categories/` | Cria uma categoria |
| `DELETE` | `/categories/:id` | Desativa uma categoria |

Exemplo de criacao de produto:

```bash
curl -X POST http://localhost:8081/v1/product/ \
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

### API de consultas

Base URL: `http://localhost:8082`

| Metodo | Rota | Descricao |
| --- | --- | --- |
| `GET` | `/product/:id` | Consulta um produto por UUID |

Exemplo:

```bash
curl http://localhost:8082/product/<product-uuid>
```

O fluxo de leitura e cache-aside: primeiro consulta o Redis; em caso de cache miss, consulta o MongoDB, grava o resultado no Redis por 5 minutos e retorna o produto.

## Kafka

O Compose cria os seguintes topicos:

- `product.created`
- `product.deleted`
- `brand.created`
- `category.created`

O worker tambem utiliza `product.dlq` e os DLQs de marcas/categorias quando configurados pelas variaveis de ambiente.

## Portas e ferramentas

| Servico | URL |
| --- | --- |
| Product Command | `http://localhost:8081` |
| Product Query | `http://localhost:8082` |
| Kafka UI | `http://localhost:8080` |
| Grafana | `http://localhost:3000` |
| Prometheus | `http://localhost:9090` |
| Redis Insight | `http://localhost:5540` |
| Redis Exporter | `http://localhost:9121/metrics` |
| Jaeger | `http://localhost:16686` |
| Loki | `http://localhost:3100` |

Credenciais padrão do Grafana:

- Usuário: `admin`
- Senha: `admin`

O Grafana provisiona automaticamente os datasources Prometheus, Loki e Jaeger. O Prometheus coleta o proprio servidor e o Redis Exporter.

## Variaveis principais

Os valores abaixo ja possuem defaults para o ambiente Docker:

| Variavel | Exemplo |
| --- | --- |
| `POSTGRES_DB_DSN` | `postgres://postgres:postgres@postgres-main:5432/product?sslmode=disable` |
| `KAFKA_BROKERS` | `kafka1:9092` |
| `REDIS_HOST` | `redis:6379` |
| `MONGO_DB_DSN` | `mongodb://root:MongoDB2019!@mongo:27017/?authSource=admin` |
| `JAEGER_URL` | `jaeger:4317` |

## Desenvolvimento local

Cada servico Go e um modulo independente:

```bash
cd product-command && go test ./...
cd ../product-query && go test ./...
cd ../worker && go test ./...
```

Para formatar um modulo:

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

## Observacoes

- Os servicos devem usar os nomes dos containers como hostnames quando executados no Compose, por exemplo `redis:6379` e `kafka1:9092`.
- O cache Redis e temporario. Alteracoes e exclusoes de produtos devem considerar a invalidacao da chave ou o TTL de 5 minutos.
- O projeto possui configuracao de tracing via OpenTelemetry/Jaeger e logging estruturado nos servicos.