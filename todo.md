# Stashly TODO

## Funcionalidades já implementadas

* [x] Listar arquivos modificados, novos e deletados
* [x] Selecionar arquivos para staging (toggle)
* [x] Commit de arquivos selecionados com input modal para mensagem
* [x] Checkout de branch existente
* [x] Criar nova branch (`git checkout -b`)
* [x] Push/Pull de branch atual
* [x] Highlight da branch atual (em cor diferente)
* [x] Navegação de foco entre painéis com Tab / Shift+Tab
* [x] Mostrar diff colorido de arquivo selecionado
* [x] Mostrar helperText com atalhos dependendo do foco

## Funcionalidades sugeridas

### Staging / Reset / Undo

* [ ] Desfazer staging de arquivos (`git reset HEAD <file>`)
* [ ] Descartar alterações locais de arquivos (`git checkout -- <file>`)
* [ ] Atalhos sugeridos: `u` para unstaging, `r` para reset

### Merge / Rebase

* [ ] Merge interativo de branch (`git merge <branch>`)
* [ ] Rebase interativo (`git rebase <branch>`)
* [ ] Destacar conflitos no painel de diff

### Histórico de commits

* [ ] Selecionar commit para ver diff detalhado
* [ ] Checkout de commit específico (modo detached HEAD)
* [ ] Busca de commits por mensagem ou autor

### Filtros e buscas

* [ ] Filtrar arquivos por status (`modified`, `new`, `deleted`)
* [ ] Busca por nome de arquivo ou commit
* [ ] Modal de input temporário para busca

### Branch management

* [ ] Separar visual de branches locais e remotas (`git branch -a`)
* [ ] Indicar upstream e divergência no nome da branch
* [ ] Ícone ou cor para branch atual já implementado

### Stash

* [ ] Criar stash (`git stash`)
* [ ] Aplicar stash (`git stash apply`)
* [ ] Remover stash (`git stash drop`)
* [ ] Painel lateral de stashes

### Notificações / Logs

* [ ] Mostrar mensagens de sucesso ou erro no rodapé

  * Commit feito
  * Push ou Pull falhou
  * Checkout realizado

### Customizações visuais

* [ ] Mais cores e destaques no diff (`+` verde, `-` vermelho, conflitos amarelos)
* [ ] Suporte a temas (dark/light) ou cores personalizadas
* [ ] Configuração de cor do highlight do item ativo (background transparente, texto rosa)

## Possível roadmap de implementação

1. Staging / Reset / Undo (mais útil no dia a dia)
2. Histórico de commits detalhado
3. Filtros e buscas de arquivos / commits
4. Merge / Rebase interativo
5. Stash management
6. Notificações e logs
7. Customizações visuais e temas
