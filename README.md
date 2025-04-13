# Task Reporter Microservice 📝📧

![Badge](https://img.shields.io/badge/Status-Em%20Desenvolvimento-yellow)
![Badge](https://img.shields.io/badge/Go-1.x-blue)
![Badge](https://img.shields.io/badge/Node.js-18.x-green)
![Badge](https://img.shields.io/badge/License-MIT-brightgreen)

## 📖 Descrição

Este projeto implementa um sistema de microserviços para gerenciar tarefas via CLI e enviar relatórios dessas tarefas (em formato CSV) por email. Ele demonstra a comunicação assíncrona usando RabbitMQ e a integração com serviços AWS (S3 para armazenamento e SES para envio de emails), com infraestrutura gerenciada por Terraform.

## ✨ Funcionalidades

* **Gerenciamento de Tarefas (CLI Go):**
    * Adicionar novas tarefas.
    * Listar tarefas pendentes ou todas as tarefas.
    * Marcar tarefas como concluídas (TODO).
    * Deletar tarefas (TODO).
* **Envio de Relatório (Fluxo Completo):**
    * Comando no CLI Go (`send-report`) para iniciar o processo.
    * Upload do arquivo `tasks.csv` para o AWS S3.
    * Publicação de mensagem no RabbitMQ com detalhes do arquivo no S3.
    * Consumo da mensagem pela API Node.js.
    * Download do CSV do S3 pela API Node.js.
    * Envio do CSV como anexo de email usando AWS SES.

## 🛠️ Tecnologias Utilizadas

* **Linguagens:** Go, Node.js
* **Frameworks/Bibliotecas:** Cobra (Go), Express (Node.js), Nodemailer (Node.js), AMQPLib (Node.js), AWS SDK (Go & Node.js)
* **Mensageria:** RabbitMQ
* **Cloud:** AWS (S3, SES)
* **Infraestrutura como Código:** Terraform
* **Containerização:** Docker, Docker Compose

## 🚀 Como Executar o Projeto (Localmente com Docker)

**Pré-requisitos:**

* Docker e Docker Compose instalados.
* Conta AWS configurada.
* AWS CLI configurado localmente (opcional, mas útil para testes e configuração inicial de credenciais).
* Terraform CLI instalado (para provisionar a infraestrutura).

**Passos:**

1.  **Clonar o Repositório:**
    ```bash
    git clone <url-do-seu-repositorio>
    cd <diretorio-do-repositorio>
    ```

2.  **Provisionar Infraestrutura AWS (Terraform):**
    * Navegue até o diretório `terraform/`.
    * **(Opcional)** Crie um arquivo `terraform.tfvars` para definir variáveis como `aws_region`, `s3_report_bucket_name` (use um nome globalmente único), `ses_sender_email`.
    * Execute `terraform init`.
    * Execute `terraform plan`.
    * Execute `terraform apply` e confirme com `yes`.
    * 🚨 **IMPORTANTE:** Após o `apply`, acesse seu email (`ses_sender_email`) e clique no link de verificação enviado pela AWS para autorizar o envio pelo SES.
    * Obtenha o nome do bucket S3 criado:
        ```bash
        terraform output s3_report_bucket_name
        ```
    * **(Se aplicável e com cuidado!)** Obtenha as credenciais IAM se geradas pelo Terraform:
        ```bash
        terraform output aws_access_key_id
        terraform output -raw aws_secret_access_key
        ```
        **Aviso:** Não comite chaves secretas ou o arquivo `.tfstate` se ele contiver segredos.

3.  **Configurar Variáveis de Ambiente (Local):**
    * Crie um arquivo `.env` dentro do diretório `go-cli/` e outro dentro de `node-api/`. Use os arquivos `.env.example` (se existirem) como base.
    * **Variáveis necessárias (exemplo):**
        ```ini
        # Comum para ambos ou específicos
        AWS_REGION=us-east-1
        AWS_ACCESS_KEY_ID=AKIAXXXXXXXXXXXXXX # Obtido do Terraform ou use credenciais existentes
        AWS_SECRET_ACCESS_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxx # Obtido do Terraform ou use credenciais existentes

        # Específico para Go CLI (go-cli/.env)
        AWS_S3_BUCKET=seu-nome-unico-de-bucket-aqui # Obtido do terraform output
        RABBITMQ_URL=amqp://user:password@rabbitmq_broker:5672 # Usuário/senha do docker-compose
        RABBITMQ_EXCHANGE=tasks_exchange
        RABBITMQ_ROUTING_KEY=report.todo
        DEFAULT_RECIPIENT_EMAIL=seu_email_padrao@example.com # Email para onde enviar por padrão

        # Específico para Node API (node-api/.env)
        RABBITMQ_URL=amqp://user:password@rabbitmq_broker:5672 # Usuário/senha do docker-compose
        RABBITMQ_QUEUE=todo_report_queue
        RABBITMQ_EXCHANGE=tasks_exchange
        RABBITMQ_ROUTING_KEY=report.todo
        SES_FROM_EMAIL=seu_email_verificado@example.com # O email verificado no SES
        PORT=3000 # Porta para a API Node.js
        ```
    * *Nota:* O `docker-compose.yml` pode ser configurado para ler esses arquivos `.env`.

4.  **Construir e Iniciar os Contêineres:**
    * Na raiz do projeto, execute:
        ```bash
        docker-compose up --build -d
        ```
    * Isso iniciará o RabbitMQ e a API Node.js.

5.  **Utilizar a Aplicação:**
    * **Adicionar/Listar Tarefas (via Go CLI):**
        * Você pode executar o CLI dentro do contêiner Docker (se configurado no compose) ou construir e rodar localmente (após instalar Go e configurar as vars de ambiente no seu terminal).
        * Exemplo via Docker Compose (se o serviço `go-cli` existir):
            ```bash
            docker-compose run --rm go-cli add "Minha tarefa"
            docker-compose run --rm go-cli list
            ```
    * **Enviar Relatório:**
        ```bash
        # Via Docker Compose (se o serviço go-cli existir)
        docker-compose run --rm go-cli send-report --recipient email_destino@example.com
        # Ou usando o email padrão definido em DEFAULT_RECIPIENT_EMAIL
        docker-compose run --rm go-cli send-report
        ```
    * Verifique a caixa de entrada do `email_destino@example.com`.
    * **Acessar RabbitMQ UI:** Abra `http://localhost:15672` no seu navegador (login: `user`/`password` ou o que foi definido no `docker-compose.yml`).

    Isso permitiria que você execute comandos como ./tasks.sh add "Nova tarefa" em vez do comando mais longo.
    `docker-compose run --rm go-cli "$@"`

## 🔄 Fluxo da Aplicação

1.  **Usuário (via CLI Go):** Executa `./tasks send-report [--recipient <email>]`.
2.  **Go CLI:** Lê `tasks.csv`.
3.  **Go CLI:** Faz upload do `tasks.csv` para o **AWS S3**.
4.  **Go CLI:** Publica uma mensagem no **RabbitMQ** contendo a localização do arquivo no S3 (bucket/key) e o email do destinatário.
5.  **RabbitMQ:** Roteia a mensagem para a fila da API Node.js.
6.  **Node API:** Consome a mensagem da fila.
7.  **Node API:** Baixa o arquivo CSV do **AWS S3** usando as informações da mensagem.
8.  **Node API:** Envia o email usando **AWS SES**, com o CSV como anexo, para o destinatário.
9.  **Usuário:** Recebe o email com o relatório.

## 👤 Autor

Gabriel Morais

[![LinkedIn](https://img.shields.io/badge/LinkedIn-blue)](https://www.linkedin.com/in/gabrielmaroco/)
[![GitHub](https://img.shields.io/badge/GitHub-grey)](https://github.com/GabrielMM182)

---