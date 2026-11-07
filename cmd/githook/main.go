// Command githook menjalankan semua Git hooks kita.
package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const hooksDir = ".githooks"

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "pemakaian: githook <install|nama-hook>")
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "install":
		err = install()
	case "pre-commit":
		err = preCommit()
	case "commit-msg":
		if len(os.Args) < 3 {
			err = fmt.Errorf("commit-msg butuh path file pesan commit")
		} else {
			err = commitMsg(os.Args[2])
		}
	case "pre-push":
		err = prePush()
	default:
		err = fmt.Errorf("hook tidak dikenal: %q", os.Args[1])
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nGAGAL: %v\n", err)
		os.Exit(1) // kode keluar bukan 0 = Git membatalkan aksinya
	}
}

// install membuat folder .githooks berisi shim untuk tiap hook,
// lalu mengarahkan Git ke folder itu lewat core.hooksPath.
func install() error {
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		return err
	}
	for _, name := range []string{"pre-commit", "commit-msg", "pre-push"} {
		// "$@" meneruskan argumen dari Git (mis. path file pesan commit).
		script := "#!/bin/sh\nexec go run ./cmd/githook " + name + " \"$@\"\n"
		path := filepath.Join(hooksDir, name)
		if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
			return err
		}
		fmt.Println("[ok] dibuat", path)
	}
	out, err := exec.Command("git", "config", "core.hooksPath", hooksDir).CombinedOutput()
	if err != nil {
		return fmt.Errorf("git config gagal: %v\n%s", err, out)
	}
	fmt.Println("[ok] core.hooksPath =", hooksDir)
	return nil
}

// run menjalankan perintah dan mengembalikan stdout+stderr gabungan.
func run(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var buf bytes.Buffer
	cmd.Stdout, cmd.Stderr = &buf, &buf
	if err := cmd.Run(); err != nil {
		cmdline := name + " " + strings.Join(args, " ")
		return buf.String(), fmt.Errorf("%s gagal:\n%s", cmdline, buf.String())
	}
	return buf.String(), nil
}
