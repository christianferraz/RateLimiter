# Rate Limiter em Go com Redis

Este projeto implementa um rate limiter em Go, utilizando Redis como mecanismo de armazenamento para controlar o número de requisições a um serviço web.

## Funcionalidades

- Limitação de requisições por IP ou token.
- Configuração de limites de requisições e tempo de bloqueio.
- Armazenamento e consulta de dados no Redis.

## Como Usar

### Pré-requisitos

- Go instalado na sua máquina.
- Instância do Redis rodando localmente ou em um servidor remoto.

### Configuração

1. Clone o repositório:
2. Edite o arquivo cmd/.env para limitar por token ou IP de acordo com os exemplos:
   TOKEN_1=token_1:30:100
   TOKEN_2=token_2:110:0
   TOKEN_3=token_3:120:20
   IP_LIMIT_1=192.168.1.1:10:0
   IP_LIMIT_2=10.0.0.4:5:10
   1. No final da chave deve-se seguir a sequência numérica;
   2. No valor deve-se atribuir o valor do token seguido de ":" Nr de Requisições por segundo seguido de ":" tempo de bloqueio em segundos
3. Para realizar o teste, a máquina deve estar instalado o docker com o plugin de docker-compose
4. Executar o comando

```bash
   RateLimiter/cmd$ docker compose up -d
```

5. Execute o programa em go

   ```bash
   RateLimiter$ go mod tidy
   RateLimiter$ cd cmd
   RateLimiter/cmd$ go run main.go
   ```

# Testes

Testes poderão ser realizados na pasta cmd com os comandos

```
RateLimiter/cmd$ go test . -v -run=TestRateLimiterUnderLimit
```
