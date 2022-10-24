package pkg1

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// формулируем ожидания: анализатор должен находить ошибку,
	// описанную в комментарии want
	fmt.Println("1")
	os.Exit(0)  // want "using os Exit!"
	log.Fatal() // want "using os Exit!"
}

func failmain() {
	fmt.Println("1")
	os.Exit(0)
}
