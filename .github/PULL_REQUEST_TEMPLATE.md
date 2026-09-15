
## Contexto e Demanda
<!-- Link para a task/história e um breve resumo do motivo dessa mudança de negócio ou técnica. -->
- **Ticket/Issue:** [PROJ-123](https://jira.suaempresa.com/browse/PROJ-123)
- **Motivação:** 

## Arquitetura e Solução Técnica
<!-- Descreva como você resolveu o problema. Se houver impacto no fluxo do CQRS, roteamento de eventos ou novos padrões aplicados, detalhe aqui. -->
- 

## Impacto e Breaking Changes
<!-- Avalie o impacto no ecossistema da aplicação e infraestrutura. -->
- [ ] Alteração de contrato de API (Swagger/OpenAPI)
- [ ] Mudança estrutural em Banco de Dados (Migrations DDL)
- [ ] Alteração em payloads ou tópicos de mensageria (ex: Apache Kafka)
- [ ] Adição ou remoção de Variáveis de Ambiente (ENVs)

## Observabilidade e Rollback
<!-- Como garantimos que isso está funcionando em produção e como voltamos atrás se der erro? -->
- **Monitoramento:** Foi adicionada alguma nova métrica de negócio (Prometheus), trace distribuído (OpenTelemetry/Jaeger) ou log estruturado (Zap)?
- **Rollback:** Qual o procedimento se a feature falhar em PRD? (ex: Apenas reverter o PR, executar *down migration*, etc).

## Definition of Done (Checklist de Qualidade)
- [ ] O pipeline de CI (Dry-Run / Lint / Testes) executou com sucesso.
- [ ] Testes unitários/integração foram adicionados ou atualizados.
- [ ] Não há vazamento de credenciais, senhas ou tokens (Hardcoded Secrets).
- [ ] A documentação técnica (README ou docs internos) foi atualizada.