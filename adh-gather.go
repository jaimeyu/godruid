package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	"github.com/go-openapi/loads"
	flags "github.com/jessevdk/go-flags"
	"github.com/spf13/viper"

	"github.com/accedian/adh-gather/config"
	"github.com/accedian/adh-gather/gather"
	adhh "github.com/accedian/adh-gather/handlers"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/restapi"
	"github.com/accedian/adh-gather/restapi/operations"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"

	pb "github.com/accedian/adh-gather/gathergrpc"
)

const (
	defaultIngestionDictionaryPath = "files/defaultIngestionDictionary.json"
	defaultSwaggerFile             = "files/swagger.yml"
)

var (
	configFilePath  string
	enableTLS       bool
	tlsKeyFile      string
	tlsCertFile     string
	ingDictFilePath string
	swaggerFilePath string

	enableChangeNotifications bool
	enableAuthorizationAAA    bool
)

func init() {
	pflag.StringVar(&configFilePath, "config", "config/adh-gather.yml", "Specify a configuration file to use")
	pflag.StringVar(&tlsKeyFile, "tlskey", "/run/secrets/tls_key", "Specify a TLS Key file")
	pflag.StringVar(&tlsCertFile, "tlscert", "/run/secrets/tls_crt", "Specify a TLS Cert file")
	pflag.BoolVar(&enableTLS, "tls", true, "Specify if TLS should be enabled")
	pflag.StringVar(&ingDictFilePath, "ingDict", defaultIngestionDictionaryPath, "Specify file path of default Ingestion Dictionary")
	pflag.StringVar(&swaggerFilePath, "swag", defaultSwaggerFile, "Specify file path of the Swagger documentation")

	pflag.BoolVar(&enableChangeNotifications, "changeNotifications", true, "Specify if Change Notifications should be enabled")

	pflag.BoolVar(&enableAuthorizationAAA, "enableAuthorizationAAA", true, "Specify if checking for Skylight AAA authorization is enabled")
}

func gRPCHandlerStart(cfg config.Provider) {
	gRPCAddress := fmt.Sprintf("%s:%d", cfg.GetString(gather.CK_server_grpc_ip.String()), cfg.GetInt(gather.CK_server_grpc_port.String()))

	lis, err := net.Listen("tcp", gRPCAddress)
	if err != nil {
		logger.Log.Fatalf("failed to start gRPC Service: %s", err.Error())
	}
	var opts []grpc.ServerOption

	gsh := adhh.CreateCoordinator()

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAdminProvisioningServiceServer(grpcServer, gsh)
	pb.RegisterTenantProvisioningServiceServer(grpcServer, gsh)

	logger.Log.Infof("gRPC service intiated on: %s", gRPCAddress)
	grpcServer.Serve(lis)
}

func startChangeNotificationHandler() {
	// Start monitoring changes and sending notifications
	hdlr := adhh.CreateChangeNotificationHandler()
	hdlr.SendChangeNotifications()
}

func modifySwagger(cfg config.Provider) {
	apiPort := cfg.GetInt(gather.CK_server_rest_port.String())

	var hostLine string
	var schemeLine string
	hostFromEnv := os.Getenv("API_TARGET")
	if len(hostFromEnv) == 0 {
		hostLine = fmt.Sprintf(`host: 'localhost:%d'`, apiPort)
		schemeLine = "  - http"
	} else {
		hostLine = fmt.Sprintf(`host: '%s'`, hostFromEnv)
		schemeLine = "  - https"
	}

	// Update the generated swagger file to contain the correct host
	input, err := ioutil.ReadFile(swaggerFilePath)
	if err != nil {
		logger.Log.Fatalf("Unable to locate swagger definition: %s", err.Error())
	}

	// Replace the host line
	containsHost := false
	lines := strings.Split(string(input), "\n")
	for i, line := range lines {
		if strings.Contains(line, `host:`) {
			containsHost = true
			lines[i] = hostLine
			break
		}
	}
	if !containsHost {
		// Insert the host into the swager file
		lines = append(lines[:2], append([]string{hostLine}, lines[2:]...)...)
	}

	// Append the appropriate scheme line:
	var index int
	for i, line := range lines {
		if strings.Contains(line, `schemes:`) {
			index = i + 1
			break
		}
	}
	lines[index] = schemeLine

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(swaggerFilePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func startRESTHandler(cfg config.Provider) {
	swaggerSpec, err := loads.Spec(swaggerFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewGatherAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	logger.Log.Info("Configuring REST Handler")
	apiHost := cfg.GetString(gather.CK_server_rest_ip.String())
	apiPort := cfg.GetInt(gather.CK_server_rest_port.String())

	server.Host = apiHost
	server.Port = apiPort

	if enableTLS {
		if _, err := os.Stat(tlsCertFile); os.IsNotExist(err) {
			// No TLS cert file
			logger.Log.Fatalf("Failed to start REST Handler: TLS cert %s does not exist", tlsCertFile)
		}
		if _, err := os.Stat(tlsKeyFile); os.IsNotExist(err) {
			// No TLS key file
			logger.Log.Fatalf("Failed to start REST Handler: TLS key %s does not exist", tlsKeyFile)
		}

		server.TLSCertificate = flags.Filename(tlsCertFile)
		server.TLSCertificateKey = flags.Filename(tlsKeyFile)
		server.TLSPort = apiPort
		server.EnabledListeners = []string{"https"}
	} else {
		server.EnabledListeners = []string{"http"}
	}

	server.ConfigureAPI()

	logger.Log.Info("Starting REST Handler")
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}

}

func main() {
	pflag.Parse()
	v := viper.New()

	v.BindPFlags(pflag.CommandLine)

	configFilePath = v.GetString("config")
	enableTLS = v.GetBool("tls")
	tlsCertFile = v.GetString("tlscert")
	tlsKeyFile = v.GetString("tlskey")
	ingDictFilePath = v.GetString("ingDict")
	swaggerFilePath = v.GetString("swag")
	enableChangeNotifications = v.GetBool("changeNotifications")

	// Load Configuration
	cfg := gather.LoadConfig(configFilePath, v)
	enableAuthorizationAAA = v.GetBool("enableAuthorizationAAA")
	cfg.Set(gather.CK_args_authorizationAAA.String(), enableAuthorizationAAA)

	debug := cfg.GetBool(gather.CK_args_debug.String())
	if debug {
		logger.SetDebugLevel(true)
	} else {
		logger.SetDebugLevel(false)
	}

	logger.Log.Infof("Starting adh-gather broker with config '%s'", configFilePath)

	if enableChangeNotifications {
		// Start monitoring changes and sending notifications
		go startChangeNotificationHandler()
	}

	// modify the swagger for this deployment
	modifySwagger(cfg)

	go gRPCHandlerStart(cfg)
	startRESTHandler(cfg)

}
