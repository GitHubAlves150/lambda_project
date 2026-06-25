// 2026-06-08
package main
import(
	"lambda_tracker/internal/repository"
)

func main() {

	repository.ConectaDB()
	repository.GPSData()
}
