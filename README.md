# Todo App

Uma simples aplicação para gerenciar suas tarefas diárias.

## 🚀 Funcionalidades

- ✨ Interface limpa e fácil de usar
- 🔐 Sistema de autenticação (registro, login, logout)
- 📝 Criação e gerenciamento de listas de tarefas
- ✅ Adição e acompanhamento de tarefas
- 🎯 Organização de tarefas por listas
- 🔄 Atualizações em tempo real

## 🛠️ Tecnologias Utilizadas

- [Go](https://golang.org/) - Linguagem de programação
- [Fiber](https://gofiber.io/) - Framework web rápido
- [GORM](https://gorm.io/) - ORM para Go
- [SQLite](https://www.sqlite.org/) - Banco de dados
- [Bulma](https://bulma.io/) - Framework CSS
- [HTMX](https://htmx.org/) - Interatividade moderna
- [Alpine.js](https://alpinejs.dev/) - Framework JavaScript minimalista

## 📋 Pré-requisitos

- Go 1.20 ou superior
- SQLite3
- [Air](https://github.com/cosmtrek/air) (opcional, para desenvolvimento)

## 🚀 Instalação

1. Clone o repositório:
```bash
git clone https://github.com/seu-usuario/todo-app.git
cd todo-app
```

2. Instale as dependências:
```bash
go mod download
```

3. Execute o projeto:
```bash
go run main.go
```

Para desenvolvimento, você pode usar o Air para live-reload:
```bash
# Instalar o Air (caso não tenha)
go install github.com/cosmtrek/air@latest

# Executar o projeto com Air
air
```

## 🗄️ Estrutura do Projeto

```
todo-app/
├── app/
│   ├── middleware/   # Middlewares da aplicação
│   ├── public/      # Arquivos estáticos (CSS, JS)
│   ├── routes/      # Rotas da aplicação
│   └── views/       # Templates HTML
├── config/
│   └── database/    # Configurações do banco de dados
├── tmp/            # Arquivos temporários (gerados pelo Air)
└── main.go         # Ponto de entrada da aplicação
```

## 🤝 Contribuindo

1. Faça um Fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.
