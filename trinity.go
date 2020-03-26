package trinity

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/bluele/gcache"
	"github.com/gin-gonic/gin"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	// GlobalTrinity global instance
	GlobalTrinity  *Trinity
	runMode        = "Local"
	rootPath       string
	configFilePath string
)

func init() {
	rootPath, _ = os.Getwd()
	configFilePath = filepath.Join(rootPath, "config", "config.yml")
}

// Trinity struct for app subconfig
type Trinity struct {
	mu             sync.RWMutex
	runMode        string
	rootpath       string
	configFilePath string

	// COMMON
	ctx              context.Context
	setting          *Setting
	customizeSetting CustomizeSetting
	db               *gorm.DB
	vCfg             *ViewSetCfg
	logger           Logger
	cache            gcache.Cache
	serviceMesh      ServiceMesh

	// GRPC
	gServer *grpc.Server
	context *Context

	// HTTP
	router *gin.Engine
}

func (t *Trinity) initDefaultValue() {
	GlobalTrinity = t
}

// GetRootPath  get rootpath
func GetRootPath() string {
	return rootPath
}

// SetRootPath  get rootpath
func SetRootPath(path string) {
	rootPath = path
}

// GetConfigFilePath  get rootpath
func GetConfigFilePath() string {
	return configFilePath
}

// SetConfigFilePath  get rootpath
func SetConfigFilePath(path string) {
	configFilePath = path
}

// GetRunMode  get RunMode
func GetRunMode() string {
	return runMode
}

// SetRunMode  set RunMode
func SetRunMode(runmode string) {
	runMode = runmode
}

// New app
// initial global trinity object
func New(ctx context.Context, customizeSetting ...CustomizeSetting) *Trinity {
	t := &Trinity{
		runMode:        runMode,
		rootpath:       rootPath,
		configFilePath: configFilePath,
		ctx:            ctx,
		context:        &Context{root: true},
	}
	t.mu.Lock()
	t.setting = newSetting(t.runMode, t.configFilePath).GetSetting()
	t.loadCustomizeSetting(customizeSetting...)
	t.logger = initLogger(t.setting)
	t.context.logger = initLogger(t.setting)
	t.context.setting = t.setting
	t.InitDatabase()
	t.db.SetLogger(t.logger)
	t.context.db.SetLogger(t.logger)

	switch t.setting.Webapp.Type {
	case "HTTP":
		t.initRouter()
		t.initViewSetCfg()
		break
	case "GRPC":
		t.initGRPCServer()
		break
	default:
		log.Fatal("wrong app type")
		break
	}

	if t.setting.ServiceMesh.AutoRegister {
		switch t.setting.GetServiceMeshType() {
		case "consul":
			c, err := NewConsulRegister(
				t.setting.GetServiceMeshAddress(),
				t.setting.GetServiceMeshPort(),
			)
			if err != nil {
				log.Fatal("get service mesh client err")
			}
			t.serviceMesh = c
			break
		case "etcd":
			c, err := NewEtcdRegister(
				t.setting.GetServiceMeshAddress(),
				t.setting.GetServiceMeshPort(),
			)
			if err != nil {
				log.Fatal("get service mesh client err")
			}
			t.serviceMesh = c
			break
		default:
			log.Fatal("wrong service mash type")
		}

	}

	t.initCache()
	t.initDefaultValue()
	t.mu.Unlock()
	return t
}

// Reload  reload trinity
func (t *Trinity) Reload(runMode string) {
	t.mu.RLock()
	t.runMode = runMode
	t.setting = newSetting(t.runMode, t.configFilePath).GetSetting()
	t.logger = initLogger(t.setting)
	t.context.logger = initLogger(t.setting)
	t.InitDatabase()
	t.initRouter()
	t.initViewSetCfg()
	t.initDefaultValue()
	t.mu.RUnlock()
}

// reloadTrinity for reload some config
func (t *Trinity) reloadTrinity() {
	t.setting = newSetting(t.runMode, t.configFilePath).GetSetting()
	t.logger = initLogger(t.setting)
	t.context.logger = initLogger(t.setting)
	t.InitDatabase()
	t.initRouter()
	t.initViewSetCfg()
	t.initDefaultValue()
}

// GetVCfg  get vcfg
func (t *Trinity) GetVCfg() *ViewSetCfg {
	t.mu.RLock()
	v := t.vCfg
	t.mu.RUnlock()
	return v
}

// NewContext  get context
func (t *Trinity) NewContext() *Context {
	t.mu.RLock()
	v := t.context
	t.mu.RUnlock()
	return v
}

// SetVCfg  get vcfg
func (t *Trinity) SetVCfg(newVCfg *ViewSetCfg) *Trinity {
	t.mu.Lock()
	t.vCfg = newVCfg
	t.reloadTrinity()
	t.mu.Unlock()
	return t
}

// GetLogger  get vcfg
func (t *Trinity) GetLogger() Logger {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.logger

}

func (t *Trinity) initGRPCServer() {

	if t.setting.GetTLSEnabled() {
		cert, err := tls.LoadX509KeyPair(t.setting.GetServerPemFile(), t.setting.GetServerKeyFile())
		if err != nil {
			log.Fatalf("tls.LoadX509KeyPair err: %v", err)
		}
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(t.setting.GetCAPemFile())
		if err != nil {
			log.Fatalf("ioutil.ReadFile err: %v", err)
		}
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			log.Fatalf("certPool.AppendCertsFromPEM err")
		}
		c := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    certPool,
		})
		opts := []grpc.ServerOption{
			grpc.Creds(c),
			grpc_middleware.WithUnaryServerChain(
				RecoveryInterceptor(t.logger),
				LoggingInterceptor(t.logger),
				UserAuthInterceptor(t.logger),
			),
		}
		t.gServer = grpc.NewServer(opts...)
	} else {
		opts := []grpc.ServerOption{
			grpc_middleware.WithUnaryServerChain(
				RecoveryInterceptor(t.logger),
				LoggingInterceptor(t.logger),
				UserAuthInterceptor(t.logger),
			),
		}
		t.gServer = grpc.NewServer(opts...)
	}

}

// GetGRPCServer get grpc server instance
func (t *Trinity) GetGRPCServer() *grpc.Server {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.gServer
}

// ServeGRPC serve GRPC
func (t *Trinity) ServeGRPC() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", t.setting.GetWebAppPort()))
	if err != nil {
		log.Fatalf("tcp port : %v  listen err: %v", t.setting.GetWebAppPort(), err)
	}
	errors := make(chan error)
	go func() {

		if err := t.serviceMesh.RegService(
			t.setting.GetProjectName(),
			t.setting.GetProjectVersion(),
			t.setting.GetWebAppAddress(),
			t.setting.GetWebAppPort(),
			t.setting.GetTags(),
			t.setting.GetDeregisterAfterCritical(),
			t.setting.GetHealthCheckInterval(),
			t.setting.GetTLSEnabled(),
		); err != nil {
			errors <- err
		}
		// logger.Logger.Log("transport", "GRPC", "address", port, "msg", "listening")
		errors <- t.gServer.Serve(lis)
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errors <- fmt.Errorf("%s", <-c)
	}()

	t.serviceMesh.DeRegService(
		t.setting.GetProjectName(),
		t.setting.GetProjectVersion(),
		t.setting.GetWebAppAddress(),
		t.setting.GetWebAppPort(),
	)

}

// ServeHTTP serve HTTP
func (t *Trinity) ServeHTTP() {
	defer t.Close()
	s := &http.Server{
		Addr:              fmt.Sprintf(":%v", t.setting.Webapp.Port),
		Handler:           t.router,
		ReadTimeout:       time.Duration(t.setting.Webapp.ReadTimeoutSecond) * time.Second,
		ReadHeaderTimeout: time.Duration(t.setting.Webapp.ReadHeaderTimeoutSecond) * time.Second,
		WriteTimeout:      time.Duration(t.setting.Webapp.WriteTimeoutSecond) * time.Second,
		IdleTimeout:       time.Duration(t.setting.Webapp.IdleTimeoutSecond) * time.Second,
		MaxHeaderBytes:    t.setting.Webapp.MaxHeaderBytes,
	}
	t.logger.Print(fmt.Sprintf("[info] %v  start http server listening : %v, version : %v", GetCurrentTimeString(time.RFC3339), t.setting.Webapp.Port, t.setting.Version))
	s.ListenAndServe()
	return
}

// Serve http
func (t *Trinity) Serve() {
	switch t.setting.GetWebAppType() {
	case "HTTP":
		t.ServeHTTP()
		break
	case "GRPC":
		t.ServeGRPC()
		break
	default:
		log.Fatal("Unsupported Web method")
		break
	}
}

// Close http
func (t *Trinity) Close() {
	t.db.Close()
	t.logger.Print(fmt.Sprintf("[info] %v  end http server listening : %v, version : %v", GetCurrentTimeString(time.RFC3339), t.setting.Webapp.Port, t.setting.Version))
}
