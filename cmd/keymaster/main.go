package main

import (
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"github.com/Symantec/Dominator/lib/log/cmdlogger"

	"github.com/Symantec/keymaster/lib/client/config"
	"github.com/Symantec/keymaster/lib/client/twofa"
	"github.com/Symantec/keymaster/lib/client/twofa/u2f"
	"github.com/Symantec/keymaster/lib/client/util"

	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

const DefaultKeysLocation = "/.ssh/"
const FilePrefix = "keymaster"

var (
	// Must be a global variable in the data segment so that the build
	// process can inject the version number on the fly when building the
	// binary. Use only from the Usage() function.
	Version = "No version provided"
)

var (
	configFilename = flag.String("config", filepath.Join(os.Getenv("HOME"), ".keymaster", "client_config.yml"), "The filename of the configuration")
	rootCAFilename = flag.String("rootCAFilename", "", "(optional) name for using non OS root CA to verify TLS connections")
	configHost     = flag.String("configHost", "", "Get a bootstrap config from this host")
	cliUsername    = flag.String("username", "", "username for keymaster")
	checkDevices   = flag.Bool("checkDevices", false, "CheckU2F devices in your system")
)

func Usage() {
	fmt.Fprintf(
		os.Stderr, "Usage of %s (version %s):\n", os.Args[0], Version)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = Usage
	flag.Parse()
	logger := cmdlogger.New()

	if *checkDevices {
		u2f.CheckU2FDevices(logger)
		return
	}

	var rootCAs *x509.CertPool
	if len(*rootCAFilename) > 1 {
		caData, err := ioutil.ReadFile(*rootCAFilename)
		if err != nil {
			logger.Printf("Failed to read caFilename")
			logger.Fatal(err)
		}
		rootCAs = x509.NewCertPool()
		if !rootCAs.AppendCertsFromPEM(caData) {
			logger.Fatal("cannot append file data")
		}

	}

	usr, err := user.Current()
	if err != nil {
		logger.Printf("cannot get current user info")
		logger.Fatal(err)
	}
	userName := usr.Username

	homeDir, err := util.GetUserHomeDir(usr)
	if err != nil {
		logger.Fatal(err)
	}

	configPath, _ := filepath.Split(*configFilename)

	err = os.MkdirAll(configPath, 0755)
	if err != nil {
		logger.Fatal(err)
	}

	if len(*configHost) > 1 {
		err = config.GetConfigFromHost(*configFilename, *configHost, rootCAs, logger)
		if err != nil {
			logger.Fatal(err)
		}
	} else if len(defaultConfigHost) > 1 { // if there is a configHost AND there is NO config file, create one
		if _, err := os.Stat(*configFilename); os.IsNotExist(err) {
			err = config.GetConfigFromHost(
				*configFilename, defaultConfigHost, rootCAs, logger)
			if err != nil {
				logger.Fatal(err)
			}
		}
	}

	config, err := config.LoadVerifyConfigFile(*configFilename)
	if err != nil {
		logger.Fatal(err)
	}

	if len(config.Base.Username) > 0 {
		userName = config.Base.Username
	}
	// command line always wins over pref or config
	if *cliUsername != "" {
		userName = *cliUsername
	}

	//sshPath := homeDir + "/.ssh/"
	privateKeyPath := filepath.Join(homeDir, DefaultKeysLocation, FilePrefix)
	sshConfigPath, _ := filepath.Split(privateKeyPath)
	err = os.MkdirAll(sshConfigPath, 0700)
	if err != nil {
		logger.Fatal(err)
	}

	tempPrivateKeyPath := filepath.Join(homeDir, DefaultKeysLocation, "keymaster-temp")
	signer, tempPublicKeyPath, err := util.GenKeyPair(
		tempPrivateKeyPath, userName+"@keymaster", logger)
	if err != nil {
		logger.Fatal(err)
	}
	defer os.Remove(tempPrivateKeyPath)
	defer os.Remove(tempPublicKeyPath)

	password, err := util.GetUserCreds(userName)
	if err != nil {
		logger.Fatal(err)
	}

	sshCert, x509Cert, err := twofa.GetCertFromTargetUrls(
		signer,
		userName,
		password,
		strings.Split(config.Base.Gen_Cert_URLS, ","),
		rootCAs,
		false,
		logger)
	if err != nil {
		logger.Fatal(err)
	}
	if sshCert == nil || x509Cert == nil {
		err := errors.New("Could not get cert from any url")
		logger.Fatal(err)
	}
	logger.Debugf(0, "Got Certs from server")
	//..
	if _, ok := os.LookupEnv("SSH_AUTH_SOCK"); ok {
		// TODO(rgooch): Parse certificate to get actual lifetime.
		cmd := exec.Command("ssh-add", "-d", privateKeyPath)
		cmd.Run()
	}

	//rename files to expected paths
	err = os.Rename(tempPrivateKeyPath, privateKeyPath)
	if err != nil {
		err := errors.New("Could not rename private Key")
		logger.Fatal(err)
	}

	err = os.Rename(tempPublicKeyPath, privateKeyPath+".pub")
	if err != nil {
		err := errors.New("Could not rename public Key")
		logger.Fatal(err)
	}

	// now we write the cert file...
	sshCertPath := privateKeyPath + "-cert.pub"
	err = ioutil.WriteFile(sshCertPath, sshCert, 0644)
	if err != nil {
		err := errors.New("Could not write ssh cert")
		logger.Fatal(err)
	}
	x509CertPath := privateKeyPath + "-x509Cert.pem"
	err = ioutil.WriteFile(x509CertPath, x509Cert, 0644)
	if err != nil {
		err := errors.New("Could not write ssh cert")
		logger.Fatal(err)
	}

	logger.Printf("Success")
	if _, ok := os.LookupEnv("SSH_AUTH_SOCK"); ok {
		// TODO(rgooch): Parse certificate to get actual lifetime.
		lifetime := fmt.Sprintf("%ds", uint64((*twofa.Duration).Seconds()))
		cmd := exec.Command("ssh-add", "-t", lifetime, privateKeyPath)
		cmd.Run()
	}
}
