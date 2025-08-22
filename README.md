# Stashly

Stashly é uma ferramenta CLI escrita em Go para facilitar a revisão, staging, commit e push de arquivos em repositórios Git.
Ela combina a simplicidade do Git com uma interface interativa no terminal, mostrando os arquivos modificados com cores, seleção por setas e filtros avançados.

# Funcionalidades

- Listagem interativa de arquivos modificados (git status) com cores:
  
  - Verde → arquivos novos
  
  - Amarelo → modificados
  
  - Vermelho → deletados

- Seleção de arquivos para staging com setas + espaço

- Commit interativo com mensagem personalizada

- Push opcional após commit

- Filtros por:
  
  - Substring (--include, --exclude)
  
  - Regex (--includer, --excluder)

- Suporte SSH com agent ou chave privada (~/.ssh/id_rsa)

- Cross-platform: Windows, Linux e macOS



## Instalação

1. Clone o repositório:

```bash
git clone https://github.com/seu-usuario/stashly.git
cd stashly
```

2. Build da CLI:
   **Windows (PowerShell):**
   `make build`
   **Linux/macOS:**
   `make build-linux`

3. O binário será gerado em `./bin/stashly` (ou `stashly.exe` no Windows)



## Uso

Basta chamar a CLI no diretório do repositório Git:

```bash
./bin/stashly --review
```

## Flags de filtragem

| Flag         | Descrição                                 | Tipo      |
| ------------ | ----------------------------------------- | --------- |
| `--include`  | Inclui arquivos que contenham a substring | Substring |
| `--exclude`  | Exclui arquivos que contenham a substring | Substring |
| `--includer` | Inclui arquivos que combinem com regex    | Regex     |
| `--excluder` | Exclui arquivos que combinem com regex    | Regex     |

**Exemplos**

```bash
# Mostrar só arquivos .go
./bin/stashly --includer ".*\.go$"

# Excluir arquivos de teste
./bin/stashly --excluder "_test\.go$"

# Combinar substring + regex
./bin/stashly --include cmd --excluder "_test\.go$"

```

## Cores

- **Verde** → arquivos novos (`Untracked`)

- **Amarelo** → modificados (`Modified`)

- **Vermelho** → deletados (`Deleted`)

## Requisitos

- [Go 1.20+](https://golang.org/dl/)

- [Git](https://git-scm.com/downloads)

- golangci-lint (opcional para lint)

## Desenvolvimento

Clone o projeto e rode:

```bash
make test      # rodar testes
make lint      # rodar lint
make run       # rodar a CLI diretamente

```

## SSH

Stashly suporta autenticação via:

1. **SSH Agent** (recomendado)

2. **Chave privada padrão** (`~/.ssh/id_rsa`)



## Estrutura do Projeto

stashly/
├── cmd/
│   └── stashly/       # main.go
├── internal/
│   ├── gitx/          # helpers Git: auth, status, stage, commit, push
│   └── ui/            # interação com usuário, cores e prompts
├── bin/               # binários gerados
└── Makefile           # build, run, lint, test