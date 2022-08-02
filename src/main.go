package main

import (
	"context"
)


func main() {
	client := authenticate()
	client.CurrentUser(context.Background())
}