// Package main initializes and runs the DNA analyzer application, setting up routes,
// controllers, and database connections necessary for processing and analyzing DNA sequences.
package main

import (
	_ "golang/docs"
	"math/rand"
	_ "net/http/pprof"
	"time"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /api/v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func generateRandomDNA(length int) string {
	bases := []rune{'A', 'C', 'G', 'T'}
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // optional seeding
	sequence := make([]rune, length)
	for i := 0; i < length; i++ {
		sequence[i] = bases[r.Intn(len(bases))]
	}
	return string(sequence)
}

func main() {
	minMutation("AACCGGTT", "AAACGGTA", []string{"AACCGGTA", "AACCGCTA", "AAACGGTA"})

}

func minMutation(startGene string, endGene string, bank []string) int {
	numMut := 0
	charDiff := make(ma)
	for i := 0; i < 8; i++ {
		if startGene[i] != endGene[i] {
			for _, seq := range bank {
				// startGeneBytes := []byte(startGene)
				startGene[i] = endGene[i]

				if string(startGeneBytes) == seq {
					numMut++
					continue;
				} 
			}
		}
	}

	return numMut
}

// func main() {
// 	rand.Seed(time.Now().UnixNano())
//
// 	sequence := generateRandomDNA(1_000_000) // 1 million bases
// 	start := time.Now()
// 	result := findRepeatedDnaSequences(sequence)
// 	duration := time.Since(start)
//
// 	fmt.Println("Found", len(result), "repeated 10-letter sequences.")
// 	fmt.Println("Execution Time:", duration)
// 	//	if err := godotenv.Load(); err != nil {
// 	//		log.Fatal(".env file not found, skipping loading", err)
// 	//	}
// 	//
// 	// postgres, err := db.New()
// 	//
// 	//	if err != nil {
// 	//		log.Fatal(err)
// 	//	}
// 	//
// 	// defer postgres.Close()
// 	//
// 	// repositories := repository.NewRepositories(postgres.GetDB())
// 	// services := services.NewServices(repositories)
// 	// controllers := controllers.NewControllers(services)
// 	//
// 	// server := server.NewServer(
// 	//
// 	//	server.WithPort(os.Getenv("PORT")),
// 	//	server.WithControllers(controllers),
// 	//
// 	// )
// 	//
// 	// server.StartServer()
// }
