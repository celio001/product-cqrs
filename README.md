# Product CQRS

Sistema de produtos baseado em CQRS (Command Query Responsibility Segregation), eventos Kafka e leitura otimizada com Redis.

## Arquitetura

O projeto separa escrita e leitura em serviços independentes, unificados pelo Kong API Gateway:

![Arquitetura da aplicacao](docs/architecture.png)

O Kong encaminha `POST` e `DELETE` de produtos para o `product-command` e `GET` para o `product-query`. O `worker` sincroniza eventos de criação e exclusão de produtos com o MongoDB. O fluxo é assíncrono: logo após criar um produto, pode haver um intervalo até que ele esteja disponível para consulta.

## Modelo de dados

O diagrama abaixo apresenta as tabelas do banco relacional e seus relacionamentos, incluindo as categorias, marcas, produtos, estoque e dados fiscais.

![Tabelas e relacionamentos do banco de dados](docs/tables.png)

```mermaid
flowchart LR
		Client[Cliente] --> Gateway[API Gateway\nKong]
		Gateway -- POST, PUT, DELETE --> Command[product-command]
		Gateway -- GET --> Query[product-query]
		Command --> Postgres[(PostgreSQL)]
		Command --> Kafka[(Kafka)]
		Kafka -- product.created\nproduct.deleted\nbrand.created --> Worker[worker]
		Worker --> Mongo[(MongoDB)]
		Worker -- Ao criar produto, grava no cache --> Redis[(Redis com TTL)]
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
* `product-command`: API de escrita. Persiste produtos, marcas e categorias no PostgreSQL e publica eventos de domínio no Kafka.
* `worker`: consome `product.created`, `product.deleted` e `brand.created`. Sincroniza produtos e marcas no MongoDB; ao criar um produto, também grava o produto no Redis. Ao excluir, altera o status no MongoDB e remove a chave do produto no Redis. O worker não consome eventos de atualização nem de categoria.
* `product-query`: API de leitura. Consulta o Redis primeiro e usa o MongoDB como fallback. Em caso de cache miss, grava o produto no Redis por 5 minutos.
* PostgreSQL: banco relacional otimizado para a parte de comandos (escrita).
* MongoDB: banco de documentos utilizado como modelo de leitura.
* Redis: cache de produtos com TTL de 5 minutos.
* Kafka: transporte dos eventos de domínio entre o `product-command` e o `worker`.
* Prometheus, Grafana, Loki e Jaeger: métricas, visualização, logs e traces. O Grafana provisiona os datasources diretamente; command e worker enviam traces ao Jaeger.

## Requisitos

* Docker Desktop com Docker Compose
* Go `1.26.5` apenas para desenvolvimento local

## Executando com Docker

Na raiz do projeto, use o Makefile:

```bash
make up
```

Para executar em segundo plano, acompanhar os logs ou consultar o estado:

```bash
make up-d
make logs
make logs SERVICE=product-query
make ps
```

Para parar os serviços ou remover também os volumes persistentes:

```bash
make down
make clean
```

## Endpoints

A comunicação externa é feita pelo Kong em `http://localhost:8000`. Atualmente, o gateway expõe apenas rotas de produtos:

| Método | Rota | Descrição |
| --- | --- | --- |
| `POST` | `/api/v1/product/` | Cria um produto |
| `GET` | `/api/v1/product/:id` | Consulta um produto por UUID |
| `DELETE` | `/api/v1/product/:id` | Desativa um produto |

O Kong também encaminha `PUT` para o serviço de comandos, mas o handler de produto ainda não registra essa rota. Marcas e categorias têm handlers no `product-command`, porém não estão expostas pelo Kong.

O produto precisa referenciar `brand_id` e `category_id` já existentes no PostgreSQL. Exemplo de criação:

```bash
curl -X POST http://localhost:8000/api/v1/product/ \
	-H 'Content-Type: application/json' \
	-d '{
		"brand_id": "b52694f9-7e1a-49dc-9931-8c0645ed9076",
		"category_id": "61053615-c6ce-4b63-b139-adceef262a14",
		"name": "Café Especial",
		"sku": "CAF-ESP-504G",
		"barcode_ean13": "7891020304050",
		"short_description": "Café 100% arábica com notas de chocolate.",
		"detailed_description": "Produzido em grandes altitudes, este café passa por um processo rigoroso de seleção de grãos para garantir a melhor experiência na sua xícara.",
		"unit_of_measure": "UN",
		"cost_price": 12.50,
		"sale_price": 24.90,
		"promotional_price": 21.90,
		"gross_weight": 0.52,
		"net_weight": 0.50,
		"height": 18.5,
		"width": 9.0,
		"length": 6.0,
		"status": "ACTIVE",
		"stock": {
			"location_aisle": "Corredor B - Prateleira 4",
			"quantity_available": 150,
			"minimum_stock": 20,
			"maximum_stock": 500
		},
		"fiscal": {
			"ncm_code": "09012100",
			"cest_code": "1709500",
			"origin_code": 0,
			"icms_rate": 18.0,
			"pis_rate": 1.65,
			"cofins_rate": 7.6,
			"ipi_rate": 0.0
		}
	}'
```

O GET usa o UUID retornado pelo POST e responde com um envelope `status`, `message` e `data`:

```bash
curl http://localhost:8000/api/v1/product/0823bdda-852a-4362-b6f1-b84d461ec8e1
```

Exemplo abreviado da resposta:

```json
{
  "status": 200,
  "message": "product retrieved successfully",
  "data": {
    "id": "0823bdda-852a-4362-b6f1-b84d461ec8e1",
    "name": "Café Especial",
    "sku": "CAF-ESP-504G",
    "stock": {
      "quantity_available": 150,
      "minimum_stock": 20,
      "maximum_stock": 500
    },
    "fiscal": {
      "ncm_code": "09012100",
      "cest_code": "1709500",
      "icms_rate": 18,
      "pis_rate": 1.65,
      "cofins_rate": 7.6,
      "ipi_rate": 0
    }
  }
}
```

O GET ainda não retorna todos os campos aceitos no POST: o modelo de leitura não inclui `detailed_description` nem `fiscal.origin_code`. O worker e o modelo de leitura também representam `width`, `length`, quantidades de estoque, `icms_rate` e `ipi_rate` como inteiros; valores fracionários nesses campos podem impedir o processamento do evento pelo worker.

O worker processa eventos de forma assíncrona, então o GET pode retornar `404` logo após o POST até a sincronização terminar. Na exclusão, após processar o evento, o worker remove a chave do Redis antes de confirmar o evento Kafka; a próxima consulta busca o produto no MongoDB, onde produtos com status `DISABLED` não são retornados.

## Kafka

O Compose cria os tópicos `product.created`, `product.update`, `product.deleted`, `brand.created` e `category.created`. Atualmente, o worker consome `product.created`, `product.deleted` e `brand.created`; não consome atualizações nem categorias. Além disso, a configuração padrão do consumidor de atualização espera `product.updated` (diferente do `product.update` criado pelo Compose), e esse consumidor não é iniciado pelo processo principal. Em caso de falha de processamento, o worker publica mensagens nas DLQs `product.dlq` e `brand.dlq`.

## Portas e ferramentas

| Serviço | URL |
| --- | --- |
| API Gateway | `http://localhost:8000` |
| Product Command | `8081` (interno no Compose; sem porta publicada no host) |
| Product Query | `8082` (interno no Compose; sem porta publicada no host) |
| Kafka UI | `http://localhost:8080` |
| Grafana | `http://localhost:3000` |
| Prometheus | `http://localhost:9090` |
| Redis Insight | `http://localhost:5540` |
| Jaeger | `http://localhost:16686` |
| Loki | `http://localhost:3100` |

Credenciais padrão do Grafana:

* Usuário: `admin`
* Senha: `admin`

O Grafana provisiona automaticamente os datasources Prometheus, Loki e Jaeger. O Compose não possui um serviço OpenTelemetry Collector.

## Variáveis principais

Os valores abaixo já possuem defaults para o ambiente Docker:

| Variável | Exemplo |
| --- | --- |
| `POSTGRES_DB_DSN` | `postgres://postgres:postgres@postgres-main:5432/product?sslmode=disable` |
| `KAFKA_BROKERS` | `kafka1:9092` |
| `REDIS_HOST` | `redis:6379` |
| `MONGO_DB_DSN` | `mongodb://root:MongoDB2019!@mongo:27017/?authSource=admin` |
| `JAEGER_URL` | `jaeger:4317` |

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
* A sincronização PostgreSQL → Kafka → MongoDB/Redis é assíncrona. O worker invalida o Redis ao excluir produtos, mas ainda não processa atualizações de produtos nem eventos de categorias.
* Command e worker exportam traces para o Jaeger usando `JAEGER_URL`; o Grafana consulta Prometheus, Loki e Jaeger diretamente.

```

```
