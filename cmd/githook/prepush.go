package main

import (
	"errors"
	"fmt"
)

// prePush menjalankan seluruh test sebelum kode dikirim ke remote.
func prePush() error {
	fmt.Println("-> pre-push: menjalankan test...")
	out, err := run("go", "test", "./...")
	fmt.Print(out)
	if err != nil {
		return errors.New("test gagal, push dibatalkan")
	}
	fmt.Println("[ok] pre-push lolos")
	return nil
}
