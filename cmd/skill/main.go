// пакеты исполняемых приложений должны называться main
package main

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/DmitriiSvarovskii/alice-skill/internal/logger"
)

// функция main вызывается автоматически при запуске приложения
func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

// функция run будет полезна при инициализации зависимостей сервера перед запуском
func run() error {
	if err := logger.Initialize(flagLogLevel); err != nil {
		return err
	}

	logger.Log.Info("Running server", zap.String("address", flagRunAddr))
	// оборачиваем хендлер webhook в middleware с логированием
	return http.ListenAndServe(flagRunAddr, logger.RequestLogger(webhook))
}

func webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Log.Debug("got request with bad method", zap.String("method", r.Method))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`
          {
            "response": {
              "text": "Извините, я пока ничего не умею"
            },
            "version": "1.0"
          }
        `))
	logger.Log.Debug("sending HTTP 200 response")
}
