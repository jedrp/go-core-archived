package apicore

import (
	"log"
	"os"
	"strconv"

	"github.com/jedrp/go-core/until"
)

// long:"host" description:"the IP to listen on" default:"localhost" env:"HOST"`
// long:"port" description:"the port to listen on for insecure connections, defaults to a random value" env:"PORT"`
// long:"tls-host" description:"the IP to listen on for tls, when not specified it's the same as --host" env:"TLS_HOST"
// long:"tls-port" description:"the port to listen on for secure connections, defaults to a random value" env:"TLS_PORT"
// long:"tls-certificate" description:"the certificate to use for secure connections" env:"TLS_CERTIFICATE"
// long:"tls-key" description:"the private key to use for secure connections" env:"TLS_PRIVATE_KEY"
// long:"tls-ca" description:"the certificate authority file to be used with mutual tls auth" env:"TLS_CA_CERTIFICATE"
// format: "disable-grpc=true; disable-rest=true; grpc-port=80; rest-port=81; host=0.0.0.0|localhost; tls-certificate=; tls-key=; tls-ca=;"
// analyze app setting string to ENV that correct with go swagger
func SetupEnvVars() {
	setting, err := until.GetSettings(API_SETTING_STR)
	if err != nil {
		panic(err)
	}
	if !setting.IsConfigured() {
		return
	}

	tlsKey := setting.GetStringValue(API_TLS_KEY_NAME, "")
	tlsCert := setting.GetStringValue(API_TLS_CERT_NAME, "")
	tlsCa := setting.GetStringValue(API_TLS_CA_NAME, "")
	host := setting.GetStringValue(API_HOST, "localhost")
	sharePort := setting.GetIntValue(API_PORT, 0)
	grpcPort := setting.GetIntValue(API_GRPC_PORT, 0)
	restPort := setting.GetIntValue(API_REST_PORT, 0)
	disableRest := setting.GetBoolValue(API_DISABLE_REST, false)
	disableGrpc := setting.GetBoolValue(API_DISABLE_GRPC, false)

	if tlsCert == "" || tlsKey == "" || tlsCa == "" {
		log.Println("TLS Key and Cert was not configured, TLS skip (listen in http scheme) ")
		os.Setenv("SCHEME", "http")
	} else {
		os.Setenv("SCHEME", "https")
		os.Setenv("TLS_PRIVATE_KEY", tlsKey)
		os.Setenv("TLS_CERTIFICATE", tlsCert)
		os.Setenv("TLS_CA_CERTIFICATE", tlsCa)
		os.Setenv("TLS_HOST", host)
		os.Setenv("TLS_PORT", strconv.Itoa(sharePort))
	}

	if sharePort > 0 {
		os.Setenv("PORT", strconv.Itoa(sharePort))
		os.Setenv("GRPC_PORT", strconv.Itoa(sharePort))
	}

	//if rest port configured then override
	if restPort > 0 {
		os.Setenv("PORT", strconv.Itoa(restPort))
	}
	if grpcPort > 0 {
		os.Setenv("GRPC_PORT", strconv.Itoa(grpcPort))
	}

	if disableGrpc {
		os.Setenv("DISABLE_REST", "true")
	}

	if disableRest {
		os.Setenv("DISABLE_GRPC", "true")
	}
}
