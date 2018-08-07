package api

import (
	"log"
	"github.com/google/easypki/pkg/easypki"
	"crypto/x509"
	"github.com/google/easypki/pkg/certificate"
	"os"
	"encoding/pem"

	"github.com/gorilla/mux"
	"easypki-ui/config"
	"net/http"
	"time"
	"encoding/json"
	"fmt"
	"crypto/x509/pkix"
	"net/url"
)

type API struct {
	cfg *config.Config
	r   *mux.Router
}

type Routes string

const (
	ListHandler Routes = "CertificateList"

	CAInfo   Routes = "CAInfo"
	CertInfo Routes = "CertInfo"

	CaCertFile    Routes = "CACertFile"
	CertFile      Routes = "CertFile"
	CertChainFile Routes = "CertChainFile"
)

func (a *API) Setup(cfg *config.Config, r *mux.Router) *mux.Router {
	a.cfg = cfg
	a.r = r

	r.HandleFunc("/", a.CertificateListHandler).
		Methods("GET").
		Name(string(ListHandler))
	r.HandleFunc("/{issuer}", a.CertificateHandler).
		Methods("GET").
		Name(string(CAInfo))
	r.HandleFunc("/{issuer}/{name}", a.CertificateHandler).
		Methods("GET").
		Name(string(CertInfo))
	r.HandleFunc("/{issuer}/file/cert", a.CertificateBundleHandler).
		Methods("GET").
		Name(string(CaCertFile))
	r.HandleFunc("/{issuer}/{name}/file/cert", a.CertificateBundleHandler).
		Methods("GET").
		Name(string(CertFile))
	r.HandleFunc("/{issuer}/{name}/file/cert", a.CertificateBundleHandler).
		Methods("GET").
		Queries("chain", "{chain}").
		Name(string(CertChainFile))

	return r
}

type ErrorResp struct {
	Error string `json:"error"`
}

type LightWeightCertificate struct {
	Name       string    `json:"name"`
	CommonName string    `json:"commonName"`
	NotAfter   time.Time `json:"notAfter"`
	NotBefore  time.Time `json:"notBefore"`
	Issuer     string    `json:"issuer"`

	Href string `json:"href"`

	Children []LightWeightCertificate `json:"children,omitempty"`
}

type Certificate struct {
	Name           string                 `json:"name"`
	Subject        pkix.Name              `json:"subject"`
	DNSNames       []string               `json:"dnsNames"`
	EmailAddresses []string               `json:"emailAddresses"`
	NotBefore      time.Time              `json:"notBefore"`
	NotAfter       time.Time              `json:"notAfter"`
	Issuer         LightWeightCertificate `json:"issuer"`

	Href string `json:"href"`
}

func (a *API) CertificateListHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	tree := []LightWeightCertificate{}

	roots, err := a.cfg.Store.Tree()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(ErrorResp{Error: fmt.Sprintf("%s", err)}); err != nil {
			panic(err)
		}
		return
	}

	for _, root := range roots {
		tree = append(tree, a.walk(root, req))
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tree); err != nil {
		panic(err)
	}
}

func (a *API) get(issuer string, name string) (*certificate.Bundle, *url.URL, error) {
	var err error
	var bundle *certificate.Bundle
	var href *url.URL

	if name == "" {
		bundle, err = a.cfg.EasyPKI.GetCA(issuer)
		href, _ = a.r.Get(string(CAInfo)).URL("issuer", issuer)
	} else {
		bundle, err = a.cfg.EasyPKI.GetBundle(issuer, name)
		href, _ = a.r.Get(string(CertInfo)).URL("issuer", issuer, "name", name)
	}

	return bundle, href, err
}

func (a *API) CertificateHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(req)

	issuerName := vars["issuer"]
	name := vars["name"]

	bundle, href, err := a.get(issuerName, name)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(ErrorResp{Error: fmt.Sprintf("%s", err)}); err != nil {
			panic(err)
		}
		return
	}
	if bundle == nil {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(ErrorResp{Error: fmt.Sprintf("not found")}); err != nil {
			panic(err)
		}
		return
	}

	issuer, issuerHref, err := a.get(bundle.Cert.Issuer.CommonName, "")

	cert := Certificate{
		Name:           name,
		Subject:        bundle.Cert.Subject,
		DNSNames:       bundle.Cert.DNSNames,
		EmailAddresses: bundle.Cert.EmailAddresses,
		NotAfter:       bundle.Cert.NotAfter,
		NotBefore:      bundle.Cert.NotBefore,
		Issuer:         LightWeightCertificate{
			Name:       issuer.Name,
			CommonName: issuer.Cert.Subject.CommonName,
			NotAfter:   issuer.Cert.NotAfter,
			NotBefore:  issuer.Cert.NotBefore,
			Issuer:     issuer.Cert.Issuer.CommonName,

			Href: decorateUrl(issuerHref, req).String(),
		},

		Href: decorateUrl(href, req).String(),
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(cert); err != nil {
		panic(err)
	}
}

func (a *API) CertificateBundleHandler(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	name := vars["name"]
	fullChain := vars["chain"] == "full"

	var err error
	conf, err := a.cfg.Store.Get(name)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(ErrorResp{Error: fmt.Sprintf("%s", err)}); err != nil {
			panic(err)
		}
		return
	}

	if conf == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(ErrorResp{Error: fmt.Sprintf("not found")}); err != nil {
			panic(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/x-pem-file; charset=UTF-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.crt", conf.Name))

	var bundle *certificate.Bundle
	if conf.IsCA {
		bundle, err = a.cfg.EasyPKI.GetCA(conf.Name)
	} else {
		bundle, err = a.cfg.EasyPKI.GetBundle(conf.Signer, conf.Name)
	}

	leaf := bundle
	chain := []*certificate.Bundle{bundle}
	if fullChain {
		for {
			if leaf.Cert.Issuer.CommonName == leaf.Cert.Subject.CommonName {
				break
			}
			ca, err := a.cfg.EasyPKI.GetCA(leaf.Cert.Issuer.CommonName)
			if err != nil {
				log.Fatalf("Failed getting signing CA %v: %v", leaf.Cert.Issuer.CommonName, err)
			}
			chain = append(chain, ca)
			leaf = ca
		}
	}
	for _, c := range chain {
		if err := pem.Encode(w, &pem.Block{
			Bytes: c.Cert.Raw,
			Type:  "CERTIFICATE",
		}); err != nil {
			log.Fatalf("Failed ecoding %v certificate: %v", c.Name, err)
		}
	}
}

func (a *API) walk(node config.TreeNode, req *http.Request) LightWeightCertificate {
	conf := node.Self()

	var err error
	var bundle *certificate.Bundle
	var href *url.URL

	if conf.IsCA {
		bundle, err = a.cfg.EasyPKI.GetCA(conf.Name)
		href, _ = a.r.Get(string(CAInfo)).URL("issuer", conf.Name)
	} else {
		bundle, err = a.cfg.EasyPKI.GetBundle(conf.Signer, conf.Name)
		href, _ = a.r.Get(string(CertInfo)).URL("issuer", conf.Signer, "name", conf.Name)
	}
	if err != nil {
		return LightWeightCertificate{
			Name: conf.Name,
		}
	}

	lw := LightWeightCertificate{
		Name:       bundle.Name,
		CommonName: bundle.Cert.Subject.CommonName,
		NotAfter:   bundle.Cert.NotAfter,
		NotBefore:  bundle.Cert.NotBefore,
		Issuer:     bundle.Cert.Issuer.CommonName,

		Href: decorateUrl(href, req).String(),

		Children: []LightWeightCertificate{},
	}

	for _, child := range node.Children() {
		lw.Children = append(lw.Children, a.walk(child, req))
	}

	return lw
}

func decorateUrl(u *url.URL, req *http.Request) *url.URL {
	if u.Scheme == "" {
		if req.TLS == nil {
			u.Scheme = "http"
		} else {
			u.Scheme = "https"
		}
	}

	if u.Host == "" {
		u.Host = req.Host
	}

	return u
}

// get retrieves a bundle from the bolt database. If fullChain is true, the
// certificate will be the chain of trust from the primary tup to root CA.
func get(pki *easypki.EasyPKI, caName, bundleName string, fullChain bool) {
	var bundle *certificate.Bundle
	if caName == "" {
		caName = bundleName
	}
	bundle, err := pki.GetBundle(caName, bundleName)
	if err != nil {
		log.Fatalf("Failed getting bundle %v within CA %v: %v", bundleName, caName, err)
	}
	leaf := bundle
	chain := []*certificate.Bundle{bundle}
	if fullChain {
		for {
			if leaf.Cert.Issuer.CommonName == leaf.Cert.Subject.CommonName {
				break
			}
			ca, err := pki.GetCA(leaf.Cert.Issuer.CommonName)
			if err != nil {
				log.Fatalf("Failed getting signing CA %v: %v", leaf.Cert.Issuer.CommonName, err)
			}
			chain = append(chain, ca)
			leaf = ca
		}
	}
	key, err := os.Create(bundleName + ".key")
	if err != nil {
		log.Fatalf("Failed creating key output file: %v", err)
	}
	if err := pem.Encode(key, &pem.Block{
		Bytes: x509.MarshalPKCS1PrivateKey(bundle.Key),
		Type:  "RSA PRIVATE KEY",
	}); err != nil {
		log.Fatalf("Failed ecoding private key: %v", err)
	}
	crtName := bundleName + ".crt"
	if fullChain {
		crtName = bundleName + "+chain.crt"
	}
	cert, err := os.Create(crtName)
	if err != nil {
		log.Fatalf("Failed creating chain output file: %v", err)
	}
	for _, c := range chain {
		if err := pem.Encode(cert, &pem.Block{
			Bytes: c.Cert.Raw,
			Type:  "CERTIFICATE",
		}); err != nil {
			log.Fatalf("Failed ecoding %v certificate: %v", c.Name, err)
		}
	}
}
