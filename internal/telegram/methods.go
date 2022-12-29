package bot

import (
	"context"
	"errors"
	"log"
	"os/exec"
	"strings"
	"time"
)

func getPrompt(rawPrompt string) string {
	return strings.Join(strings.Split(rawPrompt, " ")[1:], "")
}

func getPromptURL(rawPrompt string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Second)
	defer cancel()

	imgChn := make(chan string)

	prompt := getPrompt(rawPrompt)
	go useImagine(imgChn, prompt)

	select {
	case URL := <-imgChn:
		defer close(imgChn)

		return URL, nil
	case <-ctx.Done():
		return "", errors.New(errImagineTimeOut)
	}
}

func useImagine(imgChn chan string, prompt string) {
	script := "from app.model import imagine; print(imagine('" + prompt + "'))"
	cmd := exec.Command("python", "-c", script)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	imgChn <- string(out)
}
