package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"

	"github.com/golang/glog"
	"github.com/sbezverk/lsprefix-edge/pkg/arangodb"
	"github.com/sbezverk/lsprefix-edge/pkg/kafkamessenger"

	_ "net/http/pprof"
)

var (
	msgSrvAddr     string
	dbSrvAddr      string
	dbName         string
	dbUser         string
	dbPass         string
	vertexPrefix   string
	vertexNode     string
	edgeCollection string
)

func init() {
	runtime.GOMAXPROCS(1)
	flag.StringVar(&msgSrvAddr, "message-server", "", "URL to the messages supplying server")
	flag.StringVar(&dbSrvAddr, "database-server", "", "{dns name}:port or X.X.X.X:port of the graph database")
	flag.StringVar(&dbName, "database-name", "", "DB name")
	flag.StringVar(&dbUser, "database-user", "", "DB User name")
	flag.StringVar(&dbPass, "database-pass", "", "DB User's password")
	flag.StringVar(&vertexPrefix, "vertex-prefix", "LSPrefix_Test", "Vertex LSPrefix Collection name, default: \"LSPrefix_Test\"")
	flag.StringVar(&vertexNode, "vertex-node", "LSNode_Test", "Vertex LSNode Collection name, default: \"LSNode_Test\"")
	flag.StringVar(&edgeCollection, "edge-name", "LSPrefix_Edge_Test", "Edge Collection name, default \"LSLink_Test\"")
}

var (
	onlyOneSignalHandler = make(chan struct{})
	shutdownSignals      = []os.Signal{os.Interrupt}
)

func setupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler) // panics when called twice

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}

func main() {
	flag.Parse()
	_ = flag.Set("logtostderr", "true")
	dbSrv, err := arangodb.NewDBSrvClient(dbSrvAddr, dbUser, dbPass, dbName, vertexPrefix, vertexNode, edgeCollection)
	if err != nil {
		glog.Errorf("failed to initialize databse client with error: %+v", err)
		os.Exit(1)
	}

	if err := dbSrv.Start(); err != nil {
		if err != nil {
			glog.Errorf("failed to connect to database with error: %+v", err)
			os.Exit(1)
		}
	}

	// Initializing messenger process
	msgSrv, err := kafkamessenger.NewKafkaMessenger(msgSrvAddr, dbSrv.GetInterface())
	if err != nil {
		glog.Errorf("failed to initialize message server with error: %+v", err)
		os.Exit(1)
	}

	msgSrv.Start()

	stopCh := setupSignalHandler()
	<-stopCh

	msgSrv.Stop()
	dbSrv.Stop()

	os.Exit(0)
}
