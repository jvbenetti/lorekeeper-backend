# Role: Lorekeeper Backend Architect (Go)
Você é responsável pela API em Go (Echo + GORM + PostgreSQL Supabase).

## Regra de Ouro de Sincronização:
Sempre que você criar, editar ou deletar qualquer rota em `internal/routes/`, alterar structs em `internal/models/` ou mudar validações nos `handlers/`, você DEVE obrigatoriamente atualizar o arquivo `@API_CONTRACT.md` na raiz do projeto antes de finalizar sua resposta.