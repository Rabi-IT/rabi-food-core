# rabi-food-core — Contexto para Claude Code

## Idioma
Todo código, comentário, documentação e mensagem de commit DEVE ser em **inglês**.

## Agentes especializados

Antes de implementar qualquer coisa neste projeto, invoque o agente correspondente via Skill:

| Tarefa | Agente |
|--------|--------|
| Decisões de UX, fluxo de navegação ou regras de produto | `/product` |
| Implementar features, handlers, gateways, use cases, migrations, DI (Go) | `/api-dev` |
| Escrever ou revisar testes de integração (Go) | `/api-test` |
| Decisões de arquitetura ou dependências entre features | `/api-architect` |
| Revisar código Go antes de commitar | `/api-review` |
| Implementar componentes, páginas e UI (web) | `/web-dev` |
| Revisar código web antes de commitar | `/web-review` |

**Regra:** se a tarefa envolve código Go, invoque `/api-dev` antes de escrever. Se envolve código web, invoque `/web-dev` antes de escrever. Se a tarefa envolve uma decisão de UX, fluxo ou produto — antes de qualquer implementação — invoque `/product`. Os agentes conhecem todos os padrões obrigatórios de cada subprojeto e vão guiar a implementação correta.
