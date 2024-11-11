# rate-limiter-redis
O objetivo desta aplicação é resolver o primeiro desafio de conclusão de curso da FullCycle na trilha de Go Expert.

## Desafio
Criar um sistema de rate limiter utilizando Golang e Redis. O serviço tem o intuito de salvar valores no Redis  durante
10 segundos como se fosse um cache. Outro ponto é que o serviço deve limitar a quantidade de requisições por segundo se
baseando no token ou no IP.

## Requisitos
- Docker
- Docker Compose
- Go

## Execução
- Criar um arquivo `.env` na raíz do projeto e definir as variáveis de ambiente:
```shell
RATE_LIMIT_TOKEN=5
RATE_LIMIT_TOKEN_TIME=30
RATE_LIMIT_IP=5
RATE_LIMIT_IP_TIME=10
JWT_SECRET=secret
JWT_EXPIRESIN=300
```

- Subir o Redis:
```shell
make up-redis
```

- Executar a aplicação:
```shell
make run
```

- Realizar uma requisição POST para buscar o token:
```curl
curl --location --request POST 'localhost:8080/generate_token' \
--data ''
```

- Realizar uma requisição POST para salvar o user no Redis:
```curl
curl --location --request POST 'localhost:8080/user/jhonathann10' \
--header 'Authorization: Bearer <TOKEN>' \
--data ''
```

- Realizar uma requisição GET para buscar o user no Redis:
```curl
curl --location 'localhost:8080/user' \
--header 'Authorization: Bearer <TOKEN>' \
--data ''
```

- Quando atingir o limite de requisições, a aplicação retornará o status 429 e informando as seguintes possibilidades de
mensagem:
```json
{
    "Message": "Token rate limit exceeded",
    "Status": 429
}
```
ou
```json
{
    "Message": "IP rate limit exceeded",
    "Status": 429
}
```

## Teste
O teste que apliquei foi com 10 milhões de requisições, onde o limite era de 100 por segundos. No teste, será possível
analisar a quantidade de requisições processas e quantos tiveram sucesso:

### Exemplo de mensagem
- `Processed 1000000 requests in 2.539150417s with 300 successful responses`

Para executar o teste, basta rodar o comando:
```shell
make test
```