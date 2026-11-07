package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Format Conventional Commits: type(scope)!: deskripsi
var conventional = regexp.MustCompile(
	`^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)` +
		`(\([a-z0-9._-]+\))?(!)?: \S.{0,71}$`)

// commitMsg membaca file pesan commit (path diberikan Git) dan memeriksa baris pertamanya.
func commitMsg(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	first := strings.SplitN(strings.TrimSpace(string(data)), "\n", 2)[0]
	// Pesan merge dan revert otomatis dari Git tidak perlu diperiksa.
	if strings.HasPrefix(first, "Merge ") || strings.HasPrefix(first, "Revert ") {
		return nil
	}
	if !conventional.MatchString(first) {
		return fmt.Errorf("pesan commit tidak sesuai Conventional Commits:\n  %q\n\n"+
			"format : <type>(scope opsional): <deskripsi>\n"+
			"contoh : feat(auth): tambah login dengan Google\n"+
			"type   : feat fix docs style refactor perf test build ci chore revert", first)
	}
	fmt.Println("[ok] commit-msg lolos")
	return nil
}
