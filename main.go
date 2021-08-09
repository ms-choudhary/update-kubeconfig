package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"k8s.io/client-go/tools/clientcmd"
)

func kubeconfigPath() (string, error) {
	env := os.Getenv("KUBECONFIG")
	if env != "" {
		return env, nil
	}

	path, err := homedir.Expand("~/.kube/config")
	if err != nil {
		return "", err
	}

	return path, nil
}

func updateKubecfg(token string) {
	path, err := kubeconfigPath()
	if err != nil {
		log.Fatalf("could not find kubeconfig path: %v", err)
	}

	config, err := clientcmd.LoadFromFile(path)
	if err != nil {
		log.Fatalf("could not load kubeconfig: %v", err)
	}

	config.AuthInfos[config.CurrentContext].Token = token

	if err = clientcmd.WriteToFile(*config, path); err != nil {
		log.Fatalf("could not write to file: %v", err)
	}
}

func main() {

	// take care of subcommands
	for len(os.Args) > 0 {
		if strings.HasPrefix(os.Args[1], "--token") {
			break
		}
		os.Args = os.Args[1:]
	}

	var token string
	flag.StringVar(&token, "token", "", "token")
	flag.String("server", "", "server")
	flag.Parse()

	if token == "" {
		log.Fatalf("token empty")
	}

	updateKubecfg(token)
}
