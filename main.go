package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Mantém o mapeamento do JSON que vem do ESP32 (tudo minúsculo)
type TelemetryData struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Speed     float64 `json:"speed"`
}

// CORREÇÃO: Alinhando com o comportamento padrão do seu código antigo
type DBItem struct {
	VehicleID string  `json:"vehicle_id" dynamodbav:"vehicle_id"`
	Timestamp string  `json:"timestamp"  dynamodbav:"timestamp"`
	Latitude  float64 `json:"latitude"   dynamodbav:"latitude"`
	Longitude float64 `json:"longitude"  dynamodbav:"longitude"`
	Speed     float64 `json:"speed"      dynamodbav:"speed"`
}


var db *dynamodb.DynamoDB
var tableName = "VehiclesPositions"

// Cabeçalhos CORS necessários para o HTML/Navegador funcionar
var corsHeaders = map[string]string{
	"Content-Type":                 "application/json",
	"Access-Control-Allow-Origin":  "*", // Permite qualquer origem (HTML local ou servidor)
	"Access-Control-Allow-Headers": "Content-Type",
	"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
}

func init() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
	}))
	db = dynamodb.New(sess)
	log.Println("Lambda inicializada! Região: us-east-2")
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Roteamento baseado no método HTTP
	switch request.HTTPMethod {
	case "GET":
		return handleGet(ctx)
	case "POST":
		return handlePost(ctx, request)
	case "OPTIONS": // Necessário para o "pre-flight" do CORS no navegador
		return events.APIGatewayProxyResponse{StatusCode: 200, Headers: corsHeaders}, nil
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 405,
			Headers:    corsHeaders,
			Body:       `{"error": "Method Not Allowed"}`,
		}, nil
	}
}

// 🔍 NOVA ROTA GET: Busca os dados no DynamoDB
func handleGet(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	// Cria o input para escanear a tabela (limita a 50 itens para teste)
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
		Limit:     aws.Int64(50), 
	}

	result, err := db.Scan(input)
	if err != nil {
		log.Printf("Erro Scan DynamoDB: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers:    corsHeaders,
			Body:       `{"error": "Failed to fetch data"}`,
		}, nil
	}

	// Converte os itens do DynamoDB de volta para a estrutura Go
	var items []DBItem
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		log.Printf("Erro Unmarshal List: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers:    corsHeaders,
			Body:       `{"error": "Failed to parse data"}`,
		}, nil
	}

	// Transforma a lista em JSON
	jsonBody, _ := json.Marshal(items)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    corsHeaders,
		Body:       string(jsonBody),
	}, nil
}

// 📥 ROTA POST ANTIGA: Mantém o salvamento do ESP32
func handlePost(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var data TelemetryData
	err := json.Unmarshal([]byte(request.Body), &data)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Headers: corsHeaders, Body: `{"error": "Invalid JSON"}`}, nil
	}

	if data.VehicleID == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Headers: corsHeaders, Body: `{"error": "vehicle_id is required"}`}, nil
	}

	item := DBItem{
    VehicleID: data.VehicleID,
    Timestamp: time.Now().UTC().Format(time.RFC3339),
    Latitude:  data.Latitude,
    Longitude: data.Longitude,
    Speed:     data.Speed,
}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Headers: corsHeaders, Body: `{"error": "Internal error"}`}, nil
	}

	_, err = db.PutItem(&dynamodb.PutItemInput{Item: av, TableName: aws.String(tableName)})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Headers: corsHeaders, Body: `{"error": "Database error"}`}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    corsHeaders,
		Body:       `{"message": "OK"}`,
	}, nil
}

func main() {
	lambda.Start(handler)
}