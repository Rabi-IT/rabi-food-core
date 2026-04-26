# rabi-food-core/web — Padrões de código frontend

## Regras que nunca quebrar

### Render — renderização condicional
- **Sem ternário no render.** Nunca `condition ? <A /> : <B />` nem `condition ? x : y` dentro do JSX.
- **Bloco único opcional** → short-circuit: `{condition && <Component />}`
- **Dois blocos mutuamente exclusivos** → early return antes do JSX principal. Nunca `else` — o `if` retorna e o fluxo continua.
- **If aninhado no render** → mova para uma função auxiliar com early return interno.

### Código
- **Sem `else`** em nenhum lugar. Prefira early return.
- **If aninhado** → extraia para função com early return interno.

### Tipos
- **`readonly` em todos os tipos e arrays.** `readonly string[]`, `Readonly<Props>`, etc.
- **Sem tipagem inline** no JSX ou em funções — defina um tipo nomeado.
- **Sem números mágicos** — extraia como constante nomeada com intenção clara.
- **Listas e valores reutilizáveis** declarados com `as const` e tipados explicitamente.

### Internacionalização
- **Todo texto visível ao usuário via i18n** — sem string literal no JSX.
- **Placeholders de formato** (ex: `(11) 99999-9999`, `00000-000`) são constantes do componente, não chaves de tradução — o formato não muda entre idiomas.

## Exemplos

```tsx
// ERRADO — ternário no render
{isNew ? <FieldA /> : <FieldB />}

// CERTO — early return para blocos mutuamente exclusivos
if (!isNew) return <FieldB />
return <FieldA />
```

```tsx
// CERTO — short-circuit para bloco único opcional
{isNew && <NameField />}
{!isNew && <EmailOnlyHint />}
```

```tsx
// ERRADO
const items = ["a", "b", "c"]

// CERTO
const ITEMS = ["a", "b", "c"] as const
type Item = typeof ITEMS[number]
```

```tsx
// ERRADO
if (code.length < 6) return

// CERTO — nomeie a intenção
const OTP_CODE_LENGTH = 6
if (code.length < OTP_CODE_LENGTH) return
```
