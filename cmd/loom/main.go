package main

import (
	"fmt"
	"os"

	"github.com/ekosachev/loom/internal/adapters/cli"
	"github.com/ekosachev/loom/internal/adapters/llm"
	"github.com/ekosachev/loom/internal/adapters/storage"
	"github.com/ekosachev/loom/internal/domain/models"
	"github.com/ekosachev/loom/internal/domain/services"
	"github.com/pterm/pterm"
	"github.com/spf13/viper"
)

func main() {
	home, err := os.UserHomeDir()

	if err != nil {
		pterm.Fatal.Println("failed to find user's home directory")
		return
	}

	loomDir := fmt.Sprintf("%s\\.loom", home)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(loomDir)

	if err := viper.ReadInConfig(); err != nil {
		pterm.Fatal.Printfln("failed to read config file: %w", err)
		return
	}

	var config models.Config

	if err := viper.Unmarshal(&config); err != nil {
		pterm.Fatal.Printfln("failed to decode config: %w", err)
		return
	}

	dbPath := fmt.Sprintf("%s/loom.db", loomDir)
	sqliteRepo, err := storage.NewSQLiteAdapter(dbPath)
	if err != nil {
		pterm.Fatal.Printfln("failed to connect to database: %w", err)
		return
	}

	defer sqliteRepo.Close()

	modelStorage, err := storage.NewModelStorage(loomDir + "\\models.yaml")
	if err != nil {
		pterm.Fatal.Printfln("failed to initiate model storage: %w", err)
		return
	}

	toolsPath := loomDir + "\\tools"
	toolStorage, err := storage.NewToolStorage(toolsPath)
	if err != nil {
		pterm.Fatal.Printfln("failed to initiate tool storage: %w", err)
		return
	}

	openRouterClient := llm.NewOpenRouterAdapter()
	chatService := services.NewChatService(sqliteRepo, openRouterClient)
	workspaceService := services.NewWorkspaceService(sqliteRepo, sqliteRepo, sqliteRepo)
	branchService := services.NewBranchServcie(sqliteRepo, sqliteRepo)
	messageService := services.NewMessageService(sqliteRepo)
	modelService := services.NewModelService(modelStorage, openRouterClient, sqliteRepo)
	toolService := services.NewToolService(toolStorage)

	cliApp := cli.NewCLIApp(chatService, workspaceService, branchService, messageService, modelService, toolService, &config)
	cliApp.Execute()
}
