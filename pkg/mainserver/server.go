package mainserver

import (
	"fmt"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"

	"github.com/gigamono/gigamono/pkg/inits"
)

// MainServer represents the main server.
type MainServer struct {
	inits.App
	GinEngine *gin.Engine
}

// NewMainServer creates a new server.
func NewMainServer(app inits.App) (MainServer, error) {
	return MainServer{
		App:       app,
		GinEngine: gin.Default(),
	}, nil
}

// Listen makes the server listen on specified port.
func (server *MainServer) Listen() error {
	// Listener on TCP port.
	listener, err := net.Listen("tcp", fmt.Sprint(":", server.Config.Services.Types.WorkflowEngine.Ports.MainServer))
	if err != nil {
		return err
	}

	// Create multiplexer and delegate content-types.
	multiplexer := cmux.New(listener)
	grpcListener := multiplexer.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpListener := multiplexer.Match(cmux.HTTP1Fast())

	// Run servers concurrently and sync errors.
	grp := new(errgroup.Group)
	grp.Go(func() error { return server.grpcServe(grpcListener) })
	grp.Go(func() error { return server.httpServe(httpListener) })
	grp.Go(func() error { return multiplexer.Serve() })
	return grp.Wait()
}
