package main

import (
	"fmt"
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/orbs-network/orbs-client-sdk-go/gammacli/jsoncodec"
	"github.com/orbs-network/orbs-client-sdk-go/orbsclient"
	"github.com/pkg/errors"
	"io/ioutil"
	"net"
	"net/url"
	"path"
	"strings"
	"syscall"
)

const DEPLOY_SYSTEM_CONTRACT_NAME = "_Deployments"
const DEPLOY_SYSTEM_METHOD_NAME = "deployService"
const PROCESSOR_TYPE_NATIVE = uint32(1)
const PROCESSOR_TYPE_JAVASCRIPT = uint32(2)

func commandDeploy(requiredOptions []string) {
	codeFile := requiredOptions[0]

	if *flagContractName == "" {
		*flagContractName = getFilenameWithoutExtension(codeFile)
	}

	processorType := getProcessorTypeFromFilename(codeFile)
	code, err := ioutil.ReadFile(codeFile)
	if err != nil {
		die("Could not open code file.\n\n%s", err.Error())
	}

	signer := getTestKeyFromFile(*flagSigner)

	client := createOrbsClient()
	payload, txId, err := client.CreateSendTransactionPayload(signer.PublicKey, signer.PrivateKey, DEPLOY_SYSTEM_CONTRACT_NAME, DEPLOY_SYSTEM_METHOD_NAME, string(*flagContractName), uint32(processorType), []byte(code))
	if err != nil {
		die("Could not encode payload of the message about to be sent to server.\n\n%s", err.Error())
	}

	response, clientErr := client.SendTransaction(payload)
	handleNoConnectionGracefully(clientErr, client)
	if response != nil {
		output, err := jsoncodec.MarshalSendTxResponse(response, txId)
		if err != nil {
			die("Could not encode send-tx response to json.\n\n%s", err.Error())
		}

		log("%s\n", string(output))
	}

	if clientErr != nil {
		die("Request transaction failed on server.\n\n%s", clientErr.Error())
	}
}

func commandSendTx(requiredOptions []string) {
	inputFile := requiredOptions[0]

	signer := getTestKeyFromFile(*flagSigner)

	bytes, err := ioutil.ReadFile(inputFile)
	if err != nil {
		die("Could not open input file.\n\n%s", err.Error())
	}

	sendTx, err := jsoncodec.UnmarshalSendTx(bytes)
	if err != nil {
		die("Failed parsing input json file '%s'.\n\n%s", inputFile, err.Error())
	}

	// override contract name
	if *flagContractName != "" {
		sendTx.ContractName = *flagContractName
	}

	overrideArgsWithFlags(sendTx.Arguments)
	inputArgs, err := jsoncodec.UnmarshalArgs(sendTx.Arguments, getTestKeyFromFile)
	if err != nil {
		die(err.Error())
	}

	client := createOrbsClient()
	payload, txId, err := client.CreateSendTransactionPayload(signer.PublicKey, signer.PrivateKey, sendTx.ContractName, sendTx.MethodName, inputArgs...)
	if err != nil {
		die("Could not encode payload of the message about to be sent to server.\n\n%s", err.Error())
	}

	response, clientErr := client.SendTransaction(payload)
	handleNoConnectionGracefully(clientErr, client)
	if response != nil {
		output, err := jsoncodec.MarshalSendTxResponse(response, txId)
		if err != nil {
			die("Could not encode send-tx response to json.\n\n%s", err.Error())
		}

		log("%s\n", string(output))
	}

	if clientErr != nil {
		die("Request send-tx failed on server.\n\n%s", clientErr.Error())
	}
}

func commandRunQuery(requiredOptions []string) {
	inputFile := requiredOptions[0]

	signer := getTestKeyFromFile(*flagSigner)

	bytes, err := ioutil.ReadFile(inputFile)
	if err != nil {
		die("Could not open input file.\n\n%s", err.Error())
	}

	runQuery, err := jsoncodec.UnmarshalRead(bytes)
	if err != nil {
		die("Failed parsing input json file '%s'.\n\n%s", inputFile, err.Error())
	}

	// override contract name
	if *flagContractName != "" {
		runQuery.ContractName = *flagContractName
	}

	overrideArgsWithFlags(runQuery.Arguments)
	inputArgs, err := jsoncodec.UnmarshalArgs(runQuery.Arguments, getTestKeyFromFile)
	if err != nil {
		die(err.Error())
	}

	client := createOrbsClient()
	payload, err := client.CreateCallMethodPayload(signer.PublicKey, runQuery.ContractName, runQuery.MethodName, inputArgs...)
	if err != nil {
		die("Could not encode payload of the message about to be sent to server.\n\n%s", err.Error())
	}

	response, clientErr := client.CallMethod(payload)
	handleNoConnectionGracefully(clientErr, client)
	if response != nil {
		output, err := jsoncodec.MarshalReadResponse(response)
		if err != nil {
			die("Could not encode run-query response to json.\n\n%s", err.Error())
		}

		log("%s\n", string(output))
	}

	if clientErr != nil {
		die("Request run-query failed on server.\n\n%s", clientErr.Error())
	}
}

func commandTxStatus(requiredOptions []string) {
	txId := requiredOptions[0]

	client := createOrbsClient()
	payload, err := client.CreateGetTransactionStatusPayload(txId)
	if err != nil {
		die("Could not encode payload of the message about to be sent to server.\n\n%s", err.Error())
	}

	response, clientErr := client.GetTransactionStatus(payload)
	handleNoConnectionGracefully(clientErr, client)
	if response != nil {
		output, err := jsoncodec.MarshalTxStatusResponse(response)
		if err != nil {
			die("Could not encode status response to json.\n\n%s", err.Error())
		}

		log("%s\n", string(output))
	}

	if clientErr != nil {
		die("Request status failed on server.\n\n%s", clientErr.Error())
	}
}

func commandTxProof(requiredOptions []string) {
	txId := requiredOptions[0]

	client := createOrbsClient()
	payload, err := client.CreateGetTransactionReceiptProofPayload(txId)
	if err != nil {
		die("Could not encode payload of the message about to be sent to server.\n\n%s", err.Error())
	}

	response, clientErr := client.GetTransactionReceiptProof(payload)
	handleNoConnectionGracefully(clientErr, client)
	if response != nil {
		output, err := jsoncodec.MarshalTxProofResponse(response)
		if err != nil {
			die("Could not encode tx proof response to json.\n\n%s", err.Error())
		}

		log("%s\n", string(output))
	}

	if clientErr != nil {
		die("Request status failed on server.\n\n%s", clientErr.Error())
	}
}

func createOrbsClient() *orbsclient.OrbsClient {
	env := getEnvironmentFromConfigFile(*flagEnv)
	if len(env.Endpoints) == 0 {
		die("Environment Endpoints key does not contain any endpoints.")
	}

	endpoint := env.Endpoints[0]
	if endpoint == "localhost" {
		if !isDockerGammaRunning() && !isPortListening(*flagPort) {
			die("Local Gamma server is not running, use 'gamma-cli start-local' to start it.")
		}
		endpoint = fmt.Sprintf("http://localhost:%d", *flagPort)
	}

	return orbsclient.NewOrbsClient(endpoint, env.VirtualChain, codec.NETWORK_TYPE_TEST_NET)
}

func getProcessorTypeFromFilename(filename string) uint32 {
	if strings.HasSuffix(filename, ".go") {
		return PROCESSOR_TYPE_NATIVE
	}
	if strings.HasSuffix(filename, ".js") {
		return PROCESSOR_TYPE_JAVASCRIPT
	}
	die("Unsupported code file type.\n\nSupported code file extensions are: .go .js")
	return 0
}

// TODO: this needs to be simplified
func handleNoConnectionGracefully(err error, client *orbsclient.OrbsClient) {
	msg := fmt.Sprintf("Cannot connect to server at endpoint %s\n\nPlease check that:\n - The server is started and running.\n - The server is accessible over the network.\n - The endpoint is properly configured if a config file is used.", client.Endpoint)
	switch err := errors.Cause(err).(type) {
	case *url.Error:
		die(msg)
	case *net.OpError:
		if err.Op == "dial" || err.Op == "read" {
			die(msg)
		}
	case net.Error:
		if err.Timeout() {
			die(msg)
		}
	case syscall.Errno:
		if err == syscall.ECONNREFUSED {
			die(msg)
		}
	default:
		if err == orbsclient.NoConnectionError {
			die(msg)
		}
		return
	}
}

func getFilenameWithoutExtension(filename string) string {
	return strings.Split(path.Base(filename), ".")[0]
}

func overrideArgsWithFlags(args []*jsoncodec.Arg) {
	if *flagArg1 != "" {
		args[0].Value = *flagArg1
	}
	if *flagArg2 != "" {
		args[1].Value = *flagArg2
	}
	if *flagArg3 != "" {
		args[2].Value = *flagArg3
	}
	if *flagArg4 != "" {
		args[3].Value = *flagArg4
	}
	if *flagArg5 != "" {
		args[4].Value = *flagArg5
	}
	if *flagArg6 != "" {
		args[5].Value = *flagArg6
	}
}
