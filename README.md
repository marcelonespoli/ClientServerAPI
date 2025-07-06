# Client Server API desenvolvida em GO
Solicita a cotação do dolar do dia.

## Pré requisitos

- Instalar o banco de dados SQLite3
- Criar o database antes de executar a aplicação.

Comandos abaixo podem ser utilizados para criar esta tabela:

```sql
sqlite3 cotacoes.db;

create table cotacoes (
  id INTEGER PRIMARY KEY,
  bid string
);

select * from cotacoes;
```


