package repository

import (
	"fmt"
	"log"
	"time"
)


type Posicao struct {
	VeiculoID  string
	Latitude   float64
	Longitude  float64
	Velocidade float64
	TimeStamp  time.Time
}

func GPSData() {
	// ─────────────────────────────────────
	// DADOS DO GPS
	// ─────────────────────────────────────

	db, err := ConectaDB()
	if err != nil {
		log.Println("Erro ao conectar com o banco")
	}
	//Cria uma instancia da struct Posicao com dados reais do veículo
	//usar o envio dos dados do GPS da lilygo
	posicoes := Posicao{
		VeiculoID:  "Car-005",
		Latitude:   -34.544,
		Longitude:  -33.898,
		Velocidade: 90.0,
		TimeStamp:  time.Now(),
	}

	// ─────────────────────────────────────
	// 3. INSERT NO BANCO
	// ─────────────────────────────────────
	// 🔥 APENAS SALVA! (INSERT)
	query := `INSERT INTO posicoes (veiculo_id, latitude, longitude, velocidade, timestamp) 
              VALUES ($1, $2, $3, $4, $5)`

	db.SetConnMaxIdleTime(25)
	db.SetMaxOpenConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	// db.Exec executa a query sem retornar linhas (ideal para INSERT, UPDATE, DELETE)
	// os valores substituem os placeholders $1, $2... na mesma ordem
	// $1 = posicao.VeiculoID
	// $2 = posicao.Latitude
	// $3 = posicao.Longitude
	// $4 = posicao.Velocidade
	// $5 = posicao.Timestamp
	// o _ ignora o resultado (sql.Result) pois não precisamos do ID gerado aqui

	_, err = db.Exec(query,
		posicoes.VeiculoID,
		posicoes.Latitude,
		posicoes.Longitude,
		posicoes.Velocidade,
		posicoes.TimeStamp,
	)
	// Se o INSERT falhou (ex: tabela não existe, violação de constraint)
	// encerra o programa com o erro
	if err != nil {
		log.Fatal("Erro ao salvar: ", err)
	}

	//Se chegou até aqui é por que salvou no banco e o INSERT funcionou
	fmt.Println("✅ Posição slva com sucesso")
	//Fecha o banco ao finalizar a função
	defer db.Close()
}


//Futuramente adicionar funcionalidade como pesquisar posições antigas/ especificas...
