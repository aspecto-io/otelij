package config

// OTLP types
type Type string

const (
	TRACES  Type = "trace"
	METRICS Type = "metrics"
	LOGS    Type = "logs"
)

// Protocols
type Exporter string

const (
	OTLP   Exporter = "otlp"
	JAEGER Exporter = "jaeger"
	ZIPKIN Exporter = "zipkin"
	STDOUT Exporter = "stdout"
)

// Message format
type OtlpProtocol string

const (
	HttpJson     OtlpProtocol = "http/json"
	HttpProtobuf OtlpProtocol = "http/protobuf"
	GRPC         OtlpProtocol = "grpc"
)

type JaegerProtocol string

const (
	HttpThriftBinary JaegerProtocol = "http/thrift.binary"
	//JaegerGrpc       JaegerProtocol = "grpc" //Not supported in go jaeger exporter
	UdpThriftCompact JaegerProtocol = "udp/thrift.compact"
	//UdpThriftBinary  JaegerProtocol = "udp/thrift.binary"  //Not supported in go jaeger exporter
)

func isValidExporter(exporter string) bool {
	return isInList(getExporters(), Exporter(exporter))
}

func getExporters() []Exporter {
	return []Exporter{OTLP, JAEGER, ZIPKIN, STDOUT}
}

func isValidOtlpProtocol(protocol string) bool {
	return isInList(getOtlpProtocols(), OtlpProtocol(protocol))
}

func getOtlpProtocols() []OtlpProtocol {
	return []OtlpProtocol{HttpJson, HttpProtobuf, GRPC}
}

func isValidJaegerProtocol(protocol string) bool {
	return isInList(getJaegerProtocols(), JaegerProtocol(protocol))
}

func getJaegerProtocols() []JaegerProtocol {
	return []JaegerProtocol{HttpThriftBinary, UdpThriftCompact}
}

func isInList[T comparable](list []T, value T) bool {
	for _, val := range list {
		if val == value {
			return true
		}
	}
	return false
}
