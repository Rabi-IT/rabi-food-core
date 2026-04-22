# rabi-food-core — Contexto para Claude Code

## Agentes especializados

Antes de implementar qualquer coisa neste projeto, invoque o agente correspondente via Skill:

| Tarefa | Agente |
|--------|--------|
| Implementar features, handlers, gateways, use cases, migrations, DI | `/dev` |
| Escrever ou revisar testes de integração | `/testar` |
| Decisões de arquitetura ou dependências entre features | `/arquiteto` |
| Revisar código antes de commitar | `/revisar` |

**Regra:** se a tarefa envolve código Go, invoque `/dev` antes de escrever. O agente conhece todos os padrões obrigatórios do projeto (estrutura de features, gateway, filter structs, Wide Event, DI, etc.) e vai guiar a implementação correta.

## Visão Geral
Monolito Go para gestão de pedidos e delivery multi-tenant.
Arquitetura: Feature-based Modular Monolith.
Princípio central: dependências entre features devem ser **unidirecionais e intencionais**.

## Stack
- **Runtime**: Go 1.26
- **HTTP**: Fiber v2
- **DB**: pgx/v5 (Postgres)
- **DI**: samber/do
- **Logs**: zerolog + Wide Event pattern
- **Validação**: go-playground/validator
- **Testes**: testify suite + httpexpect + fixtures próprias

## Estrutura de Pacotes
```
api/
├── config/
├── app_context/        — session.go (UserSession)
├── domain/             — tipos compartilhados (auth/role, payment_status)
├── features/           — auth, category, order, product, subscription, tenant
│   └── [feature]/
│       ├── controller/ — fiber_*.go (handlers HTTP)
│       ├── gateway/    — interface.go + pgx_*.go
│       ├── model/
│       ├── routes/
│       └── usecases/   — lógica de negócio + integration_test.go
├── fixtures/           — helpers para testes de integração
└── libs/
    ├── database/       — PgxAdapter, PaginateInput, DefaultPage/DefaultPageSize
    │   └── filter/     — filter.Bool (optional boolean filter type)
    ├── di/             — new.go (wiring DI)
    ├── errs/           — erros de domínio por feature
    ├── http/           — FiberAdapter + middlewares
    ├── logger/         — WideEvent pattern
    └── validator/
```

## Idioma
Todo código, comentário, documentação e mensagem de commit DEVE ser em **inglês**.

## Desvios de Padrão
Se um padrão não se aplica a um caso concreto: sinalize explicitamente no PR, descreva o tradeoff e, se pontual, documente no código com um comentário direto. Não ignore silenciosamente.

## Migrations
Migrations rodam externamente — **não são executadas pelo binário da aplicação**. Para criar uma nova migration, adicione um arquivo `.sql` em `libs/database/migrations/` com anotações goose e rode `task migrate`.

## Issues Conhecidos (não replicar)
1. `gateway.Create` em Order não recebe `context.Context` — inconsistente, novos gateways devem receber
2. Módulo de pagamento externo ainda não implementado
3. Scheduling de entregas de Subscription ainda não implementado

## Comandos úteis
```bash
task test          # rodar testes de integração
task db            # subir Postgres local
task migrate       # rodar migrations pendentes
task build         # build
```
