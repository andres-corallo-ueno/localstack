package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	url, err := GeneratePresignedPutURL("my-bucket", "test.txt", 15*time.Minute)
	if err != nil {
		panic(err)
	}
	fmt.Println("URL prefirmada para subir:", url)
}

func GeneratePresignedPutURL(bucket, key string, expiresIn time.Duration) (string, error) {
	ctx := context.Background()
	// Cargar la configuración predeterminada con credenciales locales
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			"test", // Access Key ID
			"test", // Secret Access Key
			"",     // Session Token (optional)
		)),
	)
	if err != nil {
		return "", fmt.Errorf("error cargando config de AWS: %w", err)
	}
	// Crear el cliente de S3 con un EndpointResolverV2 personalizado
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("http://localhost:4566")
		o.UsePathStyle = true
	})
	// Crear el presignador
	presigner := s3.NewPresignClient(client)
	// Generar la URL prefirmada para subir (PUT)
	req, err := presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expiresIn))
	if err != nil {
		return "", fmt.Errorf("error generando URL prefirmada: %w", err)
	}
	return req.URL, nil
}
