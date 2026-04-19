# rabi-food-core — Contexto para Claude Code

## Visão Geral
Monolito Go para gestão de pedidos e delivery multi-tenant.
Arquitetura: Feature-based Modular Monolith.
Princípio central: dependências entre features devem ser **unidirecionais e intencionais**.

## Stack
- **Runtime**: Go 1.26
- **HTTP**: Fiber v2
- **ORM**: GORM + pgx/v5 (Postgres)
- **DI**: samber/do
- **Logs**: zerolog + Wide Event pattern
- **Validação**: go-playground/validator
- **Testes**: testify suite + httpexpect + fixtures próprias

## Estrutura de Pacotes
```
app/
├── main.go
├── config/             — config.go, environment.go
├── app_context/        — session.go (UserSession)
├── domain/             — tipos compartilhados (auth/role, payment_status)
├── features/           — módulos de negócio
│   ├── category/
│   ├── order/
│   ├── product/
│   ├── subscription/
│   ├── tenant/
│   └── user/
│   └── [feature]/
│       ├── controller/     — handlers HTTP (fiber_*.go)
│       ├── gateway/        — interface.go + gorm_*.go
│       ├── model/          — struct GORM
│       ├── routes/         — registro de rotas
│       └── usecases/       — lógica de negócio + integration_test.go
├── fixtures/           — helpers para testes de integração
└── libs/
    ├── database/       — GormAdapter, interface Database, paginate
    ├── di/             — wiring do DI (new.go)
    ├── errs/           — AppError + erros de domínio por feature
    ├── http/           — FiberAdapter + middlewares
    ├── logger/         — WideEvent pattern
    └── validator/      — wrapper do validator
```

## Idioma
Todo código, comentário, documentação e mensagem de commit DEVE ser em **inglês**. Isso inclui comentários inline, `COMMENT ON` em migrations, nomes de variáveis, mensagens de erro e docstrings.

## Desvios de Padrão

Os padrões definidos neste documento refletem as decisões de design atuais — não são imutáveis. Se durante o desenvolvimento você identificar que um padrão não se aplica bem a um caso concreto, **não o ignore silenciosamente nem o contorne com gambiarras**. Em vez disso:

1. Sinalize o desvio de forma explícita no PR/commit.
2. Descreva objetivamente o tradeoff: o que o padrão atual sacrifica nesse cenário e o que a alternativa ganha (e perde).
3. Se o desvio for pontual, documente-o no próprio código com um comentário direto. Se indicar uma falha no padrão em si, proponha a atualização deste documento.

Exemplos de justificativas aceitáveis: "o padrão X adiciona indireção desnecessária aqui porque Y nunca será reutilizado"; "seguir o padrão Z neste caso obriga a duplicar lógica em N lugares". Exemplos inaceitáveis: preferência estética, pressa, ou "é só uma exceção".

## Padrões Obrigatórios

### 1. Criar nova feature
Cada feature DEVE ter exatamente essa estrutura de arquivos:
- `gateway/interface.go` — interface do gateway
- `gateway/gorm_struct.go` — struct do adapter GORM
- `gateway/gorm_create.go`, `gorm_get_by_id.go`, etc. — um arquivo por operação
- `usecases/case_struct.go` — struct do use case + construtor New()
- `usecases/create.go`, `get_by_id.go`, etc. — um arquivo por use case
- `controller/controller_struct.go` — struct do controller + construtor New()
- `controller/fiber_create.go`, etc. — um handler por arquivo
- `model/[feature]_model.go` — struct GORM
- `routes/[feature]_routes.go` — registro de rotas

### 2. TenantID e UserID
SEMPRE vêm do session JWT — NUNCA do body da requisição:
```go
session := app_context.GetSession(ctx)
session.UserID   // string
session.TenantID // string
session.Role     // auth.Role
```

A extração da sessão (`app_context.GetSession`) é responsabilidade da **controller**. Use cases recebem os dados já extraídos como parâmetros explícitos, mantendo-os agnósticos ao transporte HTTP e reutilizáveis por outros use cases.

```go
// controller — extrai da sessão e repassa explicitamente
func (c *ProductController) Delete(ctx *fiber.Ctx) error {
    session := app_context.GetSession(ctx.UserContext())
    filter := gateway.DeleteFilter{ID: ctx.Params("id"), TenantID: session.TenantID}
    ...
}

// use case — recebe os dados, não conhece a origem deles
func (c *ProductCase) Delete(ctx context.Context, filter g.DeleteFilter) (bool, error) {
    ...
}
```

**Exceção aceitável:** use cases que precisam de rastreabilidade (ex: `logger.GetWideEvent`) ou propagação de contexto de cancelamento continuam usando `ctx context.Context` — isso é infraestrutura, não acoplamento ao transporte. O que deve ser evitado é usar o contexto para extrair dados de negócio (TenantID, Role) que deveriam ser parâmetros explícitos.

### 3. Erros de domínio
SEMPRE em `libs/errs/`, nunca inline:
```go
// libs/errs/order.go
var ErrOrderNotFound = newErr("ORDER_NOT_FOUND", http.StatusNotFound)

// uso no usecase:
return "", errs.ErrOrderNotFound
```

### 4. Wide Event Logger
Todo handler define o evento antes de chamar o use case:
```go
uctx := ctx.UserContext()
logger.GetWideEvent(uctx).Event = "create-order"
id, err := c.usecase.Create(uctx, data)
```

### 5. Handler padrão
```go
func (c *OrderController) Create(ctx *fiber.Ctx) error {
    data := usecases.CreateInput{}
    if err := ctx.BodyParser(&data); err != nil {
        return ctx.JSON(err)
    }
    if err := validator.V.Struct(data); err != nil {
        return err
    }
    uctx := ctx.UserContext()
    logger.GetWideEvent(uctx).Event = "create-order"
    id, err := c.usecase.Create(uctx, data)
    if err != nil {
        return err
    }
    return ctx.Status(http.StatusCreated).SendString(id)
}
```

### 6. Gateway interface
Gateway SEMPRE recebe `context.Context` como primeiro parâmetro:
```go
type OrderGateway interface {
    Create(ctx context.Context, input CreateInput) (string, error)
    GetByID(ctx context.Context, filter GetByIDFilter) (*GetByIDOutput, error)
    Patch(ctx context.Context, filter PatchFilter, values PatchValues) (bool, error)
    Paginate(ctx context.Context, filter PaginateFilter, paginate database.PaginateInput) (PaginateOutput, error)
    Delete(ctx context.Context, filter DeleteFilter) (bool, error)
}
```

### 7. Context nas queries do gateway
**Leituras (SELECT):** passam `ctx` diretamente — canceláveis pelo client.
**Escritas (INSERT/UPDATE/DELETE):** usam `context.WithoutCancel(ctx)` — preserva valores do contexto (trace, request ID) mas não propaga cancel nem deadline do request, evitando que um disconnect do client interrompa uma operação de escrita em andamento.

```go
// leitura — cancelável
rows, err := g.DB.Pool.Query(ctx, sql, args...)

// escrita — não cancelável por disconnect
_, err = g.DB.Pool.Exec(context.WithoutCancel(ctx), sql, args...)
```

**Nota:** `context.WithoutCancel` remove tanto o cancel quanto o deadline do contexto pai. A proteção contra queries travadas indefinidamente fica por conta do `statement_timeout` configurado no PostgreSQL.

### 7. Dependências entre use cases
Features podem importar outras features diretamente quando a dependência é **unidirecional e intencional** (nunca cíclica). Use interface local apenas quando testabilidade isolada for um requisito explícito.

```go
// PERMITIDO — dependência unidirecional intencional
import tenant_usecases "github.com/.../features/tenant/usecases"

type AuthCase struct {
    gateway    gateway.AuthGateway
    tenantCase *tenant_usecases.TenantCase
}

// Use interface local quando precisar mockar em testes unitários
type productLister interface {
    List(ctx context.Context, filter product_gateway.ListFilter) ([]product_gateway.ListOutput, error)
}
type OrderCase struct {
    gateway     OrderGateway
    productCase productLister
}
```

Nunca crie uma interface local apenas para disfarçar um import cíclico — isso é sinal de que a direção da dependência está errada.

### 8. Registrar no DI
Todo novo serviço deve ser registrado em `libs/di/new.go`:
```go
do.Provide(injector, func(i *do.Injector) (*feature_case.FeatureCase, error) {
    gw := do.MustInvoke[feature_gateway.FeatureGateway](i)
    return feature_case.New(gw), nil
})
```

### 9. Testes de integração
```go
package usecases_test

type TestSuite struct {
    suite.Suite
    app *fixtures.App
}

func (t *TestSuite) SetupSuite()    { t.app = fixtures.NewApp(); t.app.Start(t.T()) }
func (t *TestSuite) SetupSubTest()  { fixtures.CleanDatabase(t.T()) }
func (t *TestSuite) TearDownSuite() { t.app.Stop(t.T()) }
func TestMySuite(t *testing.T)      { suite.Run(t, new(TestSuite)) }
```

## Issues Conhecidos (não replicar)
1. `gateway.Create` em Order não recebe `context.Context` — inconsistente, novos gateways devem receber
3. Módulo de pagamento externo ainda não implementado
4. Scheduling de entregas de Subscription ainda não implementado

## Comandos úteis
```bash
# Rodar testes de integração
task test

# Subir infra local
task infra:up

# Build
task build
```
