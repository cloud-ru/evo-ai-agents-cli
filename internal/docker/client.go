package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/log"
)

// Client предоставляет функционал для работы с Docker
type Client struct {
	registryURL string
}

// NewClient создает новый Docker клиент
func NewClient(registryURL string) *Client {
	return &Client{
		registryURL: registryURL,
	}
}

// Image представляет Docker образ
type Image struct {
	Tag        string // Полный тег включая registry URL
	LocalName  string // Локальное имя образа
	Registry   string // Registry URL
	Repository string // Имя репозитория
	Version    string // Версия/тег
}

// BuildImage собирает Docker образ
func (c *Client) BuildImage(ctx context.Context, dockerfilePath, contextPath string, imageName string) error {
	log.Info("Building Docker image", "name", imageName, "context", contextPath)

	cmd := exec.CommandContext(ctx, "docker", "build",
		"-t", imageName,
		"-f", dockerfilePath,
		contextPath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build image: %w", err)
	}

	log.Info("Docker image built successfully", "name", imageName)
	return nil
}

// PushImage загружает Docker образ в registry
func (c *Client) PushImage(ctx context.Context, imageName string) error {
	log.Info("Pushing Docker image", "name", imageName)

	cmd := exec.CommandContext(ctx, "docker", "push", imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push image: %w", err)
	}

	log.Info("Docker image pushed successfully", "name", imageName)
	return nil
}

// TagImage создает новый тег для образа
func (c *Client) TagImage(ctx context.Context, sourceImage, targetImage string) error {
	log.Info("Tagging Docker image", "source", sourceImage, "target", targetImage)

	cmd := exec.CommandContext(ctx, "docker", "tag", sourceImage, targetImage)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to tag image: %w", err)
	}

	log.Info("Docker image tagged successfully", "source", sourceImage, "target", targetImage)
	return nil
}

// ImageExists проверяет существование локального образа
func (c *Client) ImageExists(ctx context.Context, imageName string) (bool, error) {
	cmd := exec.CommandContext(ctx, "docker", "image", "inspect", imageName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return false, nil
	}

	var images []interface{}
	if err := json.Unmarshal(output, &images); err != nil {
		return false, nil
	}

	return len(images) > 0, nil
}

// BuildAndPush собирает и загружает Docker образ
func (c *Client) BuildAndPush(ctx context.Context, dockerfilePath, contextPath, imageName, registryURL string) error {
	// 1. Собираем образ
	localImageName := imageName
	if err := c.BuildImage(ctx, dockerfilePath, contextPath, localImageName); err != nil {
		return fmt.Errorf("failed to build image: %w", err)
	}

	// 2. Создаем тег для registry
	fullImageName := fmt.Sprintf("%s/%s", registryURL, imageName)
	if err := c.TagImage(ctx, localImageName, fullImageName); err != nil {
		return fmt.Errorf("failed to tag image: %w", err)
	}

	// 3. Загружаем в registry
	if err := c.PushImage(ctx, fullImageName); err != nil {
		return fmt.Errorf("failed to push image: %w", err)
	}

	log.Info("Successfully built and pushed image", "image", fullImageName)
	return nil
}

// FindDockerfile находит Dockerfile в директории проекта
func FindDockerfile(projectDir string) (string, error) {
	dockerfile := filepath.Join(projectDir, "Dockerfile")

	// Проверяем существование Dockerfile
	if _, err := os.Stat(dockerfile); err != nil {
		return "", fmt.Errorf("Dockerfile not found in %s", projectDir)
	}

	return dockerfile, nil
}
