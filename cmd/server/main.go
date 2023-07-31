package main

import (
	"errors"
	"fmt"
	accountRPCImpl "myschool/internal/controller/rpc/account"
	studentRPCImpl "myschool/internal/controller/rpc/student"
	"myschool/internal/repositories"
	accountSvcImpl "myschool/internal/services/account"
	"myschool/pkg/controller/rpc"
	accountRPCIface "myschool/pkg/controller/rpc/account"
	studentRPCIface "myschool/pkg/controller/rpc/student"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-chi/chi/v5"

	"github.com/dewanggasurya/logger"
	"github.com/dewanggasurya/logger/log"
)

func main() {

	loggerEngine := logger.Default().SetLevel(logger.Level(logger.DebugLevel)).SetTemplate(logger.DefaultTemplate())
	logger.SetLogger(loggerEngine)

	wg := &sync.WaitGroup{}
	repo := repositories.Setup(wg)

	studentServer := &studentRPCImpl.StudentRPCServer{Repo: repo}
	studentTwirpHandler := studentRPCIface.NewStudentServer(studentServer,
		rpc.TwirpHookOption(studentRPCIface.StudentPathPrefix))
	studentTwirpHandler2 := rpc.WithEvaluateHeaders(studentTwirpHandler)

	accountService := accountSvcImpl.AccountService{Repo: repo}
	accountRPCServer := &accountRPCImpl.AccountRPCServer{AccountSvc: &accountService}
	accountTwirpHandler := accountRPCIface.NewAccountServer(accountRPCServer,
		rpc.TwirpHookOption(accountRPCIface.AccountPathPrefix))
	accountTwirpHandler2 := rpc.WithEvaluateHeaders(accountTwirpHandler)

	router := chi.NewRouter()
	router.Mount(accountRPCIface.AccountPathPrefix, accountTwirpHandler2)
	router.Mount(studentRPCIface.StudentPathPrefix, studentTwirpHandler2)

	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", 8080),
		Handler: router,
	}

	go func(srv *http.Server) {
		log.Infof("listen on port : %d", 8080)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Infof("listen: %s", err.Error())
		}
	}(server)

	defer func() {
		if panicErr := recover(); panicErr != nil {
			log.Debugf("panic: %s", panicErr)
			//send to i.e. sentry
			//sentry.CurrentHub().Recover(panicErr)
			//sentry.Flush(time.Second * 5)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	wg.Add(1)
	go repo.Disconnect()
	wg.Wait()
	log.Infoln("Bye!")

}
