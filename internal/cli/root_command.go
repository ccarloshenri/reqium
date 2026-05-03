package cli

import (
	"github.com/spf13/cobra"

	"reqium/internal/app"
	"reqium/internal/implementations/filesystem"
	"reqium/internal/implementations/formatter"
	httpinfra "reqium/internal/implementations/http"
	"reqium/internal/implementations/storage"
	"reqium/internal/implementations/variables"
	"reqium/internal/tui"
)

func NewRootCommand() *cobra.Command {
	store, err := storage.NewDefaultJSONStore()
	if err != nil {
		panic(err)
	}

	reader := filesystem.NewOSFileReader()
	client := httpinfra.NewNetHTTPClient()
	historyRepo := storage.NewHistoryRepository(store)
	envRepo := storage.NewEnvironmentRepository(store)
	collectionRepo := storage.NewCollectionRepository(store)
	resolver := variables.NewTemplateVariableResolver()
	envService := app.NewEnvironmentService(envRepo)
	requestService := app.NewRequestServiceWithHistory(client, formatter.NewResponseFormatter(true), historyRepo)
	deps := dependencies{
		reader:            reader,
		client:            client,
		historyRepo:       historyRepo,
		envRepo:           envRepo,
		collectionRepo:    collectionRepo,
		resolver:          resolver,
		envService:        envService,
		requestService:    requestService,
		collectionService: app.NewCollectionService(collectionRepo),
		historyService:    app.NewHistoryService(historyRepo, requestService),
	}

	root := &cobra.Command{
		Use:          "reqium",
		Short:        "A fast, minimal terminal API client",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return tui.Run(cmd.Context(), tui.Services{
				History:      historyRepo,
				Environments: envRepo,
				Collections:  collectionRepo,
			})
		},
	}

	root.AddCommand(newMethodCommand("get", "GET", deps))
	root.AddCommand(newMethodCommand("post", "POST", deps))
	root.AddCommand(newMethodCommand("put", "PUT", deps))
	root.AddCommand(newMethodCommand("patch", "PATCH", deps))
	root.AddCommand(newMethodCommand("delete", "DELETE", deps))
	root.AddCommand(newHistoryCommand(deps))
	root.AddCommand(newEnvironmentCommand(deps))
	root.AddCommand(newCollectionCommand(deps))
	root.AddCommand(newRunCommand(deps))

	return root
}

type dependencies struct {
	reader            *filesystem.OSFileReader
	client            *httpinfra.NetHTTPClient
	historyRepo       *storage.HistoryRepository
	envRepo           *storage.EnvironmentRepository
	collectionRepo    *storage.CollectionRepository
	resolver          *variables.TemplateVariableResolver
	envService        *app.EnvironmentService
	requestService    *app.RequestService
	collectionService *app.CollectionService
	historyService    *app.HistoryService
}
