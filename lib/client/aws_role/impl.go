package aws_role

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Cloud-Foundations/keymaster/lib/paths"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

const rsaKeySize = 2048

func parseArn(arnString string) (*arn.ARN, error) {
	parsedArn, err := arn.Parse(arnString)
	if err != nil {
		return nil, err
	}
	switch parsedArn.Service {
	case "iam", "sts":
	default:
		return nil, fmt.Errorf("unsupported service: %s", parsedArn.Service)
	}
	splitResource := strings.Split(parsedArn.Resource, "/")
	if len(splitResource) < 2 || splitResource[0] != "assumed-role" {
		return nil, fmt.Errorf("invalid resource: %s", parsedArn.Resource)
	}
	// Normalise to the actual role ARN, rather than an ARN showing how the
	// credentials were obtained. This mirrors the way AWS policy documents are
	// written.
	parsedArn.Region = ""
	parsedArn.Service = "iam"
	parsedArn.Resource = "role/" + splitResource[1]
	return &parsedArn, nil
}

func newManager(p Params) (*Manager, error) {
	cert, err := p.getRoleCertificateTLS()
	if err != nil {
		return nil, err
	}
	p.Logger.Printf("got AWS Role certificate for: %s\n", p.roleArn)
	manager := &Manager{
		Params:  p,
		tlsCert: cert,
	}
	go manager.refreshLoop()
	return manager, nil
}

func (m *Manager) getClientCertificate(cri *tls.CertificateRequestInfo) (
	*tls.Certificate, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.tlsCert, m.tlsError
}

func (m *Manager) refreshLoop() {
	for ; ; time.Sleep(time.Minute) {
		m.refreshOnce()
	}
}

func (m *Manager) refreshOnce() {
	if m.tlsCert != nil {
		refreshTime := m.tlsCert.Leaf.NotBefore.Add(
			m.tlsCert.Leaf.NotAfter.Sub(m.tlsCert.Leaf.NotBefore) * 3 / 4)
		time.Sleep(time.Until(refreshTime))
	}
	if cert, err := m.getRoleCertificateTLS(); err != nil {
		m.Logger.Println(err)
		if m.tlsCert == nil {
			m.mutex.Lock()
			m.tlsError = err
			m.mutex.Unlock()
		}
	} else {
		m.mutex.Lock()
		m.tlsCert = cert
		m.tlsError = nil
		m.mutex.Unlock()
		m.Logger.Printf("refreshed AWS Role certificate for: %s\n", m.roleArn)
	}
}

// Returns certificate PEM block.
func (p *Params) getRoleCertificate() ([]byte, error) {
	if err := p.setupVerify(); err != nil {
		return nil, err
	}
	presignedReq, err := p.stsPresignClient.PresignGetCallerIdentity(p.Context,
		&sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}
	p.Logger.Debugf(1, "presigned URL: %v\n", presignedReq.URL)
	hostPath := p.KeymasterServer + paths.RequestAwsRoleCertificatePath
	body := &bytes.Buffer{}
	body.Write(p.pemPubKey)
	req, err := http.NewRequestWithContext(p.Context, "POST", hostPath, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("claimed-arn", p.roleArn)
	req.Header.Add("presigned-method", presignedReq.Method)
	req.Header.Add("presigned-url", presignedReq.URL)
	resp, err := p.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("got error from call %s, url='%s'\n",
			resp.Status, hostPath)
	}
	return ioutil.ReadAll(resp.Body)
}

func (p *Params) getRoleCertificateTLS() (*tls.Certificate, error) {
	certPEM, err := p.getRoleCertificate()
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return nil, fmt.Errorf("unable to decode certificate PEM block")
	}
	if block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("invalid certificate type: %s", block.Type)
	}
	x509Cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return &tls.Certificate{
		Certificate: [][]byte{block.Bytes},
		PrivateKey:  p.Signer,
		Leaf:        x509Cert,
	}, nil
}

func (p *Params) setupVerify() error {
	if p.isSetup {
		return nil
	}
	if p.KeymasterServer == "" {
		return fmt.Errorf("no keymaster server specified")
	}
	if p.Logger == nil {
		return fmt.Errorf("no logger specified")
	}
	if p.Context == nil {
		p.Context = context.TODO()
	}
	if p.HttpClient == nil {
		p.HttpClient = http.DefaultClient
	}
	if p.Signer == nil {
		signer, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
		if err != nil {
			return err
		}
		p.Signer = signer
	}
	derPubKey, err := x509.MarshalPKIXPublicKey(p.Signer.Public())
	if err != nil {
		return err
	}
	p.derPubKey = derPubKey
	p.pemPubKey = pem.EncodeToMemory(&pem.Block{
		Bytes: p.derPubKey,
		Type:  "PUBLIC KEY",
	})
	awsConfig, err := config.LoadDefaultConfig(p.Context,
		config.WithEC2IMDSRegion())
	if err != nil {
		return err
	}
	p.awsConfig = awsConfig
	p.stsClient = sts.NewFromConfig(awsConfig)
	p.stsPresignClient = sts.NewPresignClient(p.stsClient)
	idOutput, err := p.stsClient.GetCallerIdentity(p.Context,
		&sts.GetCallerIdentityInput{})
	if err != nil {
		return err
	}
	p.Logger.Debugf(0, "Account: %s, ARN: %s, UserId: %s\n",
		*idOutput.Account, *idOutput.Arn, *idOutput.UserId)
	parsedArn, err := parseArn(*idOutput.Arn)
	if err != nil {
		return err
	}
	p.roleArn = parsedArn.String()
	p.isSetup = true
	return nil
}