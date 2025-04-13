package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"tasks/internal/storage"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rabbitmq/amqp091-go"
	"github.com/spf13/cobra"
)

type reportMessage struct {
	S3Bucket       string `json:"s3_bucket"`
	S3Key          string `json:"s3_key"`
	RecipientEmail string `json:"recipient_email"`
}

var recipientEmail string

var sendReportCmd = &cobra.Command{
	Use:   "send-report",
	Short: "Send tasks report to processing",
	Long:  `Upload tasks.csv to S3 and send processing message to RabbitMQ.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Se o email não foi fornecido via flag, tenta obter da variável de ambiente
		if recipientEmail == "" {
			recipientEmail = os.Getenv("DEFAULT_RECIPIENT_EMAIL")
			if recipientEmail == "" {
				fmt.Fprintln(os.Stderr, "Recipient email not provided and DEFAULT_RECIPIENT_EMAIL not set")
				os.Exit(1)
			}
		}

		// Carrega o arquivo CSV
		f, err := storage.LoadFile(storage.DataFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading tasks file: %v\n", err)
			os.Exit(1)
		}
		defer storage.CloseFile(f)

		// Configura o cliente AWS S3
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to load AWS SDK config: %v\n", err)
			os.Exit(1)
		}

		// Obtém o nome do bucket da variável de ambiente
		bucket := os.Getenv("AWS_S3_BUCKET")
		if bucket == "" {
			fmt.Fprintln(os.Stderr, "AWS_S3_BUCKET environment variable not set")
			os.Exit(1)
		}

		// Gera o nome do arquivo com timestamp
		timestamp := time.Now().Format("20060102-150405")
		s3Key := fmt.Sprintf("reports/tasks-%s.csv", timestamp)

		// Faz upload do arquivo para S3
		client := s3.NewFromConfig(cfg)
		_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: &bucket,
			Key:    &s3Key,
			Body:   f,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error uploading file to S3: %v\n", err)
			os.Exit(1)
		}

		// Conecta ao RabbitMQ
		rabbitmqURL := os.Getenv("RABBITMQ_URL")
		if rabbitmqURL == "" {
			fmt.Fprintln(os.Stderr, "RABBITMQ_URL environment variable not set")
			os.Exit(1)
		}

		conn, err := amqp091.Dial(rabbitmqURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error connecting to RabbitMQ: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close()

		ch, err := conn.Channel()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating RabbitMQ channel: %v\n", err)
			os.Exit(1)
		}
		defer ch.Close()

		// Prepara a mensagem para o RabbitMQ
		msg := reportMessage{
			S3Bucket:       bucket,
			S3Key:          s3Key,
			RecipientEmail: recipientEmail,
		}

		msgBody, err := json.Marshal(msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling message: %v\n", err)
			os.Exit(1)
		}

		// Publica a mensagem no RabbitMQ
		exchange := os.Getenv("RABBITMQ_EXCHANGE")
		routingKey := os.Getenv("RABBITMQ_ROUTING_KEY")
		if exchange == "" || routingKey == "" {
			fmt.Fprintln(os.Stderr, "RABBITMQ_EXCHANGE or RABBITMQ_ROUTING_KEY environment variable not set")
			os.Exit(1)
		}

		err = ch.PublishWithContext(
			context.TODO(),
			exchange,
			routingKey,
			false,
			false,
			amqp091.Publishing{
				ContentType: "application/json",
				Body:        msgBody,
			},
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error publishing message to RabbitMQ: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Report successfully uploaded to S3 and message sent to RabbitMQ\n")
		fmt.Printf("S3 Location: s3://%s/%s\n", bucket, s3Key)
		fmt.Printf("Recipient Email: %s\n", recipientEmail)
	},
}

func init() {
	sendReportCmd.Flags().StringVarP(&recipientEmail, "recipient", "r", "", "Email of the report recipient")
}
