# Docker Setup - Backend Camisaria Store

Este guia explica como executar o backend da loja de camisas usando Docker.

## Pré-requisitos

- Docker instalado
- Docker Compose instalado (opcional, mas recomendado)

## Configuração

1. **Copie o arquivo de exemplo de variáveis de ambiente:**
   ```bash
   cp env.example .env
   ```

2. **Edite o arquivo `.env` com suas configurações:**
   - Ajuste a string de conexão do banco de dados
   - Configure a chave JWT
   - Ajuste outras configurações conforme necessário

## Execução

### Usando Docker Compose (Recomendado)

Para executar tanto a aplicação quanto o banco MySQL:

```bash
docker-compose up --build
```

Para executar em background:
```bash
docker-compose up -d --build
```

Para parar os serviços:
```bash
docker-compose down
```

### Usando apenas Docker

Se você já tem um banco MySQL rodando separadamente:

```bash
# Build da imagem
docker build -t camisaria-backend .

# Executar o container
docker run -p 4041:4041 --env-file .env camisaria-backend
```

## Configurações Importantes

### Banco de Dados
- A aplicação espera um MySQL 8.0+
- As migrações são executadas automaticamente na inicialização
- Configure a variável `DB` no arquivo `.env` com a string de conexão correta

### Porta
- A aplicação roda na porta 4041 por padrão
- Configure a variável `PORT` no `.env` se quiser alterar

### JWT
- Configure uma chave secreta forte na variável `JWT_SECRET`
- Nunca use a chave padrão em produção

## Estrutura dos Arquivos

- `Dockerfile` - Arquivo de build multi-stage otimizado
- `docker-compose.yml` - Configuração completa com banco de dados
- `.dockerignore` - Arquivos ignorados no build
- `env.example` - Exemplo de configurações de ambiente

## Desenvolvimento

Para desenvolvimento local com hot-reload, continue usando:
```bash
go run main.go
```

O Docker é recomendado para produção e testes de integração.

