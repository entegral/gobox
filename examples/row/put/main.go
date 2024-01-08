package main

import (
	"context"
	"fmt"

	"github.com/entegral/gobox/examples/exampleLib"
)

func main() {
	ctx := context.Background()
	user := exampleLib.PutUser(ctx, "")
	fmt.Println("User:", user)
}
