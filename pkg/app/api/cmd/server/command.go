package server

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
	"golang.org/x/sync/errgroup"

	"github.com/gihyodocker/taskapp/pkg/app/api/handler"
	"github.com/gihyodocker/taskapp/pkg/cli"
	"github.com/gihyodocker/taskapp/pkg/config"
	"github.com/gihyodocker/taskapp/pkg/db"
	"github.com/gihyodocker/taskapp/pkg/repository"
	"github.com/gihyodocker/taskapp/pkg/server"
)

type command struct {
	port        int
	gracePeriod time.Duration
	configFile  string
}

// 　CLIライブラリcobraのインスタンスを作成
func NewCommand() *cobra.Command {
	//　CLIのオプションとして取る値のデフォルト値を設定
	c := &command{
		port:        8180,
		gracePeriod: 5 * time.Second,
	}
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start up the api server",
		RunE:  cli.WithContext(c.execute),
	}
	//　APIサーバをListenするポート番号を定義
	cmd.Flags().IntVar(&c.port, "port", c.port, "The port number used to run HTTP api.")
	// graceful shutdownまでの待ち時間を定義
	cmd.Flags().DurationVar(&c.gracePeriod, "grace-period", c.gracePeriod, "How long to wait for graceful shutdown.")
	// APIサーバの設定ファイルのパスを定義
	cmd.Flags().StringVar(&c.configFile, "config-file", c.configFile, "The path to the config file.")
	// --config-fileオプションの指定を必須に設定
	cmd.MarkFlagRequired("config-file")
	return cmd
}

func (c *command) execute(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	// --config-fileで指定された設定ファイルの読み込み
	appConfig, err := config.LoadConfigFile(c.configFile)
	if err != nil {
		slog.Error("failed to load api configuration",
			slog.String("config-file", c.configFile),
			err,
		)
		return err
	}
	// Open MySQL connection
	dbConn, err := createMySQL(*appConfig.Database)
	if err != nil {
		slog.Error("failed to open MySQL connection", err)
		return err
	}

	// Initialize repositories
	taskRepo := repository.NewTask(dbConn)

	// Handlers
	taskHandler := handler.NewTask(taskRepo)

	options := []server.Option{
		server.WithGracePeriod(c.gracePeriod),
	}
	httpServer := server.NewHTTPServer(c.port, options...)

	httpServer.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	httpServer.Put("/api/tasks/{id}", taskHandler.Update)
	httpServer.Delete("/api/tasks/{id}", taskHandler.Delete)
	httpServer.Get("/api/tasks/{id}", taskHandler.Get)
	httpServer.Post("/api/tasks", taskHandler.Create)
	httpServer.Get("/api/tasks", taskHandler.List)

	group.Go(func() error {
		return httpServer.Serve(ctx)
	})

	if err := group.Wait(); err != nil {
		slog.Error("failed while running", err)
		return err
	}
	return nil
}

func createMySQL(conf config.Database) (*sql.DB, error) {
	options := []db.Option{
		db.WithMaxIdleConns(conf.MaxIdleConns),
		db.WithMaxOpenConns(conf.MaxOpenConns),
		db.WithConnMaxLifetime(conf.ConnMaxLifetime),
	}

	ds := db.NewMySQLDatasource(conf.Username, conf.Password, conf.Host, conf.DBName)
	return db.OpenDB(ds, options...)
}
