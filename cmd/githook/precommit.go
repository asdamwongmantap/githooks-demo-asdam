package main

import (
	"fmt"
	"strings"
)

// preCommit memeriksa format (gofmt) file .go yang di-stage, lalu go vet.
func preCommit() error {
	fmt.Println("-> pre-commit: memeriksa format dan vet...")
	out, err := run("git", "diff", "--cached", "--name-only", "--diff-filter=ACM")
	if err != nil {
		return err
	}
	var goFiles []string
	for _, f := range strings.Fields(out) {
		if strings.HasSuffix(f, ".go") {
			goFiles = append(goFiles, f)
		}
	}
	if len(goFiles) == 0 {
		fmt.Println("   tidak ada file .go yang di-stage, dilewati")
		return nil
	}
	unformatted, err := run("gofmt", append([]string{"-l"}, goFiles...)...)
	if err != nil {
		return err
	}
	if strings.TrimSpace(unformatted) != "" {
		return fmt.Errorf("file belum rapi, jalankan `gofmt -w`:\n%s", unformatted)
	}
	if _, err := run("go", "vet", "./..."); err != nil {
		return err
	}
	fmt.Println("[ok] pre-commit lolos")
	return nil
}
