package config

import (
	"time"
	"crypto/x509/pkix"
	"github.com/google/easypki/pkg/easypki"
	"crypto/x509"
	"github.com/google/easypki/pkg/certificate"
	"fmt"
)

type Store interface {
	Add(cert Cert) error
	Get(name string) (*Cert, error)
	Tree() ([]TreeNode, error)
}

type TreeNode interface {
	Children() []TreeNode
	Self() Cert
}

type Cert struct {
	Name           string    `yaml:"name"`
	Subject        pkix.Name `yaml:"subject"`
	CommonName     string    `yaml:"commonName"`
	DNSNames       []string  `yaml:"dnsNames"`
	EmailAddresses []string  `yaml:"emailAddresses"`

	Signer string        `yaml:"signer"`
	Expire time.Duration `yaml:"expire"`

	IsCA     bool `yaml:"isCA"`
	IsClient bool `yaml:"isClient"`
}

type Config struct {
	Store   Store
	EasyPKI *easypki.EasyPKI
}

func (c *Config) Init() error {
	tree, err := c.Store.Tree()
	if err != nil {
		return err
	}

	for _, node := range tree {
		err := c.walk(node)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) walk(node TreeNode) error {
	c.makeCert(node.Self())
	for _, node := range node.Children() {
		err := c.walk(node)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) makeCert(cert Cert) error {
	req := &easypki.Request{
		Name: cert.Name,
		Template: &x509.Certificate{
			Subject:        cert.Subject,
			NotAfter:       time.Now().Add(cert.Expire),
			IsCA:           cert.IsCA,
			DNSNames:       cert.DNSNames,
			EmailAddresses: cert.EmailAddresses,
		},
		IsClientCertificate: cert.IsClient,
	}
	if cert.IsCA {
		req.Template.MaxPathLen = -1
	}
	req.Template.Subject.CommonName = cert.CommonName

	var signer *certificate.Bundle
	var err error
	if cert.Signer != "" {
		signer, err = c.EasyPKI.GetCA(cert.Signer)
		if err != nil {
			return fmt.Errorf("cannot sign %v because cannot get CA %v: %v", cert.Name, cert.Signer, err)
		}
	}

	if err := c.EasyPKI.Sign(signer, req); err != nil {
		return fmt.Errorf("cannot create bundle for %v: %v", cert.Name, err)
	}

	return nil
}
