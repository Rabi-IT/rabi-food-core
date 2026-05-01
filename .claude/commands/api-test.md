Aja como o especialista em QA e Testes do rabi-food-core.

Gere um arquivo `integration_test.go` completo seguindo o padrão exato do projeto.

## Padrão obrigatório de arquivo

```go
package usecases_test

import (
    "net/http"
    "testing"

    "github.com/Rabi-IT/rabi-food-core/fixtures"
    "github.com/gavv/httpexpect/v2"
    "github.com/stretchr/testify/suite"
)

type TestSuite struct {
    suite.Suite
    app *fixtures.App
}

func (t *TestSuite) SetupSuite()    { t.app = fixtures.NewApp(); t.app.Start(t.T()) }
func (t *TestSuite) SetupSubTest()  { fixtures.CleanDatabase(t.T()) }
func (t *TestSuite) TearDownSuite() { t.app.Stop(t.T()) }
func TestMySuite(t *testing.T)      { suite.Run(t, new(TestSuite)) }
```

## Fixtures disponíveis
- `fixtures.Tenant.Create(t, nil)` → retorna struct com `UserID`, `ID`
- `fixtures.Auth.UserToken(t, userID)` → retorna JWT string
- `fixtures.Order.Create(t, tenantID, userID)` → cria e retorna order
- `fixtures.Product.Create(t, tenantID)` → cria e retorna product
- `fixtures.Category.Create(t, tenantID)` → cria e retorna category
- `fixtures.Subscription.Create(t, tenantID, userID)` → cria e retorna subscription
- `fixtures.CleanDatabase(t)` → já chamado no SetupSubTest

## Padrão de request
```go
httpexpect.Default(t.T(), fixtures.AppURL).
    Request(http.MethodPost, "/order/").
    WithHeader("Authorization", "Bearer "+token).
    WithJSON(body).
    Expect().
    Status(http.StatusCreated).
    Body().NotEmpty()
```

## Casos OBRIGATÓRIOS em todo teste de endpoint
1. Happy path (201/200/204)
2. Sem token de autenticação → 401
3. Token de tenant diferente tentando acessar recurso de outro tenant → 403 ou 404
4. Body inválido / campos obrigatórios ausentes → 400
5. Recurso inexistente (GET/PATCH/DELETE) → 404
6. Duplicidade quando aplicável → 409

## Feature/endpoint para testar
$ARGUMENTS
