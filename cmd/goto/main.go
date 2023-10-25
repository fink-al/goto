package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/grafviktor/goto/internal/config"
	"github.com/grafviktor/goto/internal/logger"
	"github.com/grafviktor/goto/internal/storage"
	"github.com/grafviktor/goto/internal/ui"
	"github.com/grafviktor/goto/internal/version"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	// Set application version and build details
	version.Set(buildVersion, buildDate, buildCommit)

	lg, err := logger.New()
	if err != nil {
		log.Fatalf("Can't create log file %v", err)
	}

	lg.Debug("Starting application")
	lg.Debug("Version %s", version.BuildVersion())
	lg.Debug("Build date %s", version.BuildDate())
	lg.Debug("Commit %s", version.BuildCommit())

	ctx := context.Background()
	appConfig := config.New(ctx, &lg)

	st, err := storage.GetStorage(ctx, appConfig)
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	uiComponent := ui.NewMainModel(ctx, st)
	p := tea.NewProgram(uiComponent, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Println("Error running program:", err)

		os.Exit(1)
	}

	err = appConfig.Save()
	if err != nil {
		log.Fatalf("Can't save application config before closing %v", err)
	}
}
