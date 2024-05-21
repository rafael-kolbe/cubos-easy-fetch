package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/chzyer/readline"
)

func runCommand(command string, dir string) error {
	parts := strings.Fields(command)
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = nil
	cmd.Stderr = os.Stderr
	if dir != "" {
		cmd.Dir = dir
	}
	return cmd.Run()
}

func input(msg string) string {
	rl, err := readline.New(msg)
	if err != nil {
		panic(fmt.Sprintf("Erro ao criar readline: %v", err))
	}
	defer rl.Close()

	input, err := rl.Readline()
	if err != nil {
		panic(fmt.Sprintf("Erro ao ler o input: %v", err))
	}

	return input
}

func main() {
	sshAddress := input("Digite o endereço SSH do repositório: ")
	sshAddress = strings.TrimSpace(sshAddress)

	prsInput := input("Digite os números dos PRs (separados por vírgula): ")

	prs := strings.Split(prsInput, ",")
	for i := range prs {
		prs[i] = strings.TrimSpace(prs[i])
	}

	repoNameParts := strings.Split(sshAddress, "/")
	repoName := strings.TrimSuffix(repoNameParts[len(repoNameParts)-1], ".git")

	fmt.Printf("Executando: git clone %s\n", sshAddress)
	if err := runCommand(fmt.Sprintf("git clone %s", sshAddress), ""); err != nil {
		fmt.Printf("Erro ao executar comando: %v\n", err)
		return
	}

	for _, pr := range prs {
		commands := []string{
			fmt.Sprintf("git remote add pr-%s %s", pr, sshAddress),
			fmt.Sprintf("git fetch pr-%s pull/%s/head:pr-%s", pr, pr, pr),
		}

		for _, command := range commands {
			fmt.Printf("Executando: %s\n", command)
			if err := runCommand(command, repoName); err != nil {
				fmt.Printf("Erro ao executar comando: %v\n", err)
			}
		}
	}
}
