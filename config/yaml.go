package config

import (
	"fmt"
	"io/ioutil"
	"github.com/go-yaml/yaml"
	"sort"
)

type config struct {
	Certs []Cert `yaml:"certs"`
}

type Yaml struct {
	Path string
}

type certsBySigner map[string][]Cert

func (y *Yaml) readConfig() ([]Cert, error) {
	b, err := ioutil.ReadFile(y.Path)
	if err != nil {
		return nil, fmt.Errorf("failed reading configuration file %v: %v", y.Path, err)
	}
	conf := &config{}
	if err := yaml.Unmarshal(b, conf); err != nil {
		return nil, fmt.Errorf("failed umarshaling yaml config (%v) %v: %v", y.Path, string(b), err)
	}

	return conf.Certs, nil
}

func (y *Yaml) Add(cert Cert) error {
	return fmt.Errorf("not implmented")
}

func (y *Yaml) Get(name string) (*Cert, error) {
	certs, err := y.readConfig()
	if err != nil {
		return nil, err
	}

	sort.Slice(certs, func(i, j int) bool {
		return certs[i].Name < certs[j].Name
	})

	i := sort.Search(len(certs), func(i int) bool {
		return certs[i].Name >= name
	})

	if i < len(certs) && certs[i].Name == name {
		return &certs[i], nil
	}

	return nil, nil
}

func (y *Yaml) Tree() ([]TreeNode, error) {
	certs, err := y.readConfig()
	if err != nil {
		return nil, err
	}

	flat := certsBySigner{}
	for _, cert := range certs {
		flat[cert.Signer] = append(flat[cert.Signer], cert)
	}

	for key, value := range flat {
		t := key + " - "
		for _, v := range value {
			t += v.Name + ", "
		}
	}

	var roots []TreeNode
	for _, cert := range certs {
		// Roots are all CA certificates that are self signed, or externally signed
		if cert.IsCA && (cert.Signer == "" || cert.Signer == cert.Name || flat[cert.Signer] == nil) {
			roots = append(roots, &YamlCertNode{self: cert, certsBySigner: &flat})
		}
	}

	return roots, nil
}

type YamlCertNode struct {
	self          Cert
	certsBySigner *certsBySigner
}

func (n *YamlCertNode) Self() Cert {
	return n.self
}

func (n *YamlCertNode) Children() []TreeNode {
	children := (*n.certsBySigner)[n.self.Name]
	if children == nil {
		return []TreeNode{}
	}

	var nodes []TreeNode
	for _, cert := range children {
		nodes = append(nodes, &YamlCertNode{self: cert, certsBySigner: n.certsBySigner})
	}

	return nodes
}
