//go:build integration

package integration

const (
	// projectName - имя проекта для Docker-контейнеров и сети
	projectName = "inventory-service"

	// inventoryCollectionName - имя коллекции MongoDB для частей
	inventoryCollectionName = "parts"

	// Параметры для IAM контейнера
	iamAppName    = "iam-app"
	iamDockerfile = "deploy/docker/iam/Dockerfile"
	iamGRPCPort   = "50053"
)
