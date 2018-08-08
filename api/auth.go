package api

import (
	"net/url"
	"net/http"
	"fmt"
	"encoding/json"
	"time"

	"github.com/gorilla/mux"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"strings"
	"encoding/base64"
	"sort"
	"github.com/lestrrat-go/jwx/jwa"
		"context"
)

type DiscoveryMetadata struct {
	// REQUIRED . URL using the https scheme with no query or fragment component that the OP asserts as its Issuer Identifier. If Issuer discovery is supported (see Section 2), this value MUST be identical to the issuer value returned by WebFinger. This also MUST be identical to the iss Claim value in ID Tokens issued from this Issuer.
	Issuer string `json:"issuer"`
	// REQUIRED . URL of the OP's OAuth 2.0 Authorization Endpoint [OpenID.Core].
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	// REQUIRED . URL of the OP's OAuth 2.0 Token Endpoint [OpenID.Core]. This is REQUIRED unless only the Implicit Flow is used.
	TokenEndpoint string `json:"token_endpoint"`
	// RECOMMENDED . URL of the OP's UserInfo Endpoint [OpenID.Core]. This URL MUST use the https scheme and MAY contain port, path, and query parameter components.
	UserinfoEndpoint string `json:"userinfo_endpoint"`
	// REQUIRED . URL of the OP's JSON Web Key Set [JWK] document. This contains the signing key(s) the RP uses to validate signatures from the OP. The JWK Set MAY also contain the Server's encryption key(s), which are used by RPs to encrypt requests to the Server. When both signing and encryption keys are made available, a use (Key Use) parameter value is REQUIRED for all keys in the referenced JWK Set to indicate each key's intended usage. Although some algorithms allow the same key to be used for both signatures and encryption, doing so is NOT RECOMMENDED, as it is less secure. The JWK x5c parameter MAY be used to provide X.509 representations of keys provided. When used, the bare key values MUST still be present and MUST match those in the certificate.
	JwksUri string `json:"jwks_uri"`
	// RECOMMENDED . URL of the OP's Dynamic Client Registration Endpoint [OpenID.Registration].
	RegistrationEndpoint string `json:"registration_endpoint"`
	// RECOMMENDED . JSON array containing a list of the OAuth 2.0 [RFC6749] scope values that this server supports. The server MUST support the openid scope value. Servers MAY choose not to advertise some supported scope values even when this parameter is used, although those defined in [OpenID.Core] SHOULD be listed, if supported.
	ScopesSupported string `json:"scopes_supported"`
	// REQUIRED . JSON array containing a list of the OAuth 2.0 response_type values that this OP supports. Dynamic OpenID Providers MUST support the code, id_token, and the token id_token Response Type values.
	ResponseTypesSupported string `json:"response_types_supported"`
	// OPTIONAL . JSON array containing a list of the OAuth 2.0 response_mode values that this OP supports, as specified in OAuth 2.0 Multiple Response Type Encoding Practices [OAuth.Responses]. If omitted, the default for Dynamic OpenID Providers is ["query", "fragment"].
	ResponseModesSupported string `json:"response_modes_supported,omitempty"`
	// OPTIONAL . JSON array containing a list of the OAuth 2.0 Grant Type values that this OP supports. Dynamic OpenID Providers MUST support the authorization_code and implicit Grant Type values and MAY support other Grant Types. If omitted, the default value is ["authorization_code", "implicit"].
	GrantTypesSupported string `json:"grant_types_supported,omitempty"`
	// OPTIONAL . JSON array containing a list of the Authentication Context Class References that this OP supports.
	AcrValuesSupported string `json:"acr_values_supported,omitempty"`
	// REQUIRED . JSON array containing a list of the Subject Identifier types that this OP supports. Valid types include pairwise and public.
	SubjectTypesSupported string `json:"subject_types_supported"`
	// REQUIRED . JSON array containing a list of the JWS signing algorithms (alg values) supported by the OP for the ID Token to encode the Claims in a JWT [JWT]. The algorithm RS256 MUST be included. The value none MAY be supported, but MUST NOT be used unless the Response Type used returns no ID Token from the Authorization Endpoint (such as when using the Authorization Code Flow).
	IdTokenSigningAlgValuesSupported string `json:"id_token_signing_alg_values_supported"`
	// OPTIONAL . JSON array containing a list of the JWE encryption algorithms (alg values) supported by the OP for the ID Token to encode the Claims in a JWT [JWT].
	IdTokenEncryptionAlgValuesSupported string `json:"id_token_encryption_alg_values_supported,omitempty"`
	// OPTIONAL . JSON array containing a list of the JWE encryption algorithms (enc values) supported by the OP for the ID Token to encode the Claims in a JWT [JWT].
	IdTokenEncryptionEncValuesSupported string `json:"id_token_encryption_enc_values_supported,omitempty"`
	// OPTIONAL . JSON array containing a list of the JWS [JWS] signing algorithms (alg values) [JWA] supported by the UserInfo Endpoint to encode the Claims in a JWT [JWT]. The value none MAY be included.
	UserinfoSigningAlgValuesSupported string `json:"userinfo_signing_alg_values_supported,omitempty"`
	// OPTIONAL . JSON array containing a list of the JWE [JWE] encryption algorithms (alg values) [JWA] supported by the UserInfo Endpoint to encode the Claims in a JWT [JWT].
	UserinfoEncryptionAlgValuesSupported string `json:"userinfo_encryption_alg_values_supported,omitempty"`
	// OPTIONAL . JSON array containing a list of the JWE encryption algorithms (enc values) [JWA] supported by the UserInfo Endpoint to encode the Claims in a JWT [JWT].
	UserinfoEncryptionEncValuesSupported string `json:"userinfo_encryption_enc_values_supported,omitempty"`
	// OPTIONAL . JSON array containing a list of the JWS signing algorithms (alg values) supported by the OP for Request Objects, which are described in Section 6.1 of OpenID Connect Core 1.0 [OpenID.Core]. These algorithms are used both when the Request Object is passed by value (using the request parameter) and when it is passed by reference (using the request_uri parameter). Servers SHOULD support none and RS256.
	RequestObjectSigningAlgValuesSupported string `json:"request_object_signing_alg_values_supported,omitempty"`
	// OPTIONAL . JSON array containing a list of the JWE encryption algorithms (alg values) supported by the OP for Request Objects. These algorithms are used both when the Request Object is passed by value and when it is passed by reference.
	RequestObjectEncryptionAlgValuesSupported string `json:"request_object_encryption_alg_values_supported,omitempty"`
	// OPTIONAL . JSON array containing a list of the JWE encryption algorithms (enc values) supported by the OP for Request Objects. These algorithms are used both when the Request Object is passed by value and when it is passed by reference.
	RequestObjectEncryptionEncValuesSupported string `json:"request_object_encryption_enc_values_supported,omitempty"`
	// OPTIONAL . JSON array containing a list of Client Authentication methods supported by this Token Endpoint. The options are client_secret_post, client_secret_basic, client_secret_jwt, and private_key_jwt, as described in Section 9 of OpenID Connect Core 1.0 [OpenID.Core]. Other authentication methods MAY be defined by extensions. If omitted, the default is client_secret_basic -- the HTTP Basic Authentication Scheme specified in Section 2.3.1 of OAuth 2.0 [RFC6749].
	TokenEndpointAuthMethodsSupported string `json:"token_endpoint_auth_methods_supported,omitempty"`
	// OPTIONAL . JSON array containing a list of the JWS signing algorithms (alg values) supported by the Token Endpoint for the signature on the JWT [JWT] used to authenticate the Client at the Token Endpoint for the private_key_jwt and client_secret_jwt authentication methods. Servers SHOULD support RS256. The value none MUST NOT be used.
	TokenEndpointAuthSigningAlgValuesSupported string `json:"token_endpoint_auth_signing_alg_values_supported,omitempty"`
	// OPTIONAL . JSON array containing a list of the display parameter values that the OpenID Provider supports. These values are described in Section 3.1.2.1 of OpenID Connect Core 1.0 [OpenID.Core].
	DisplayValuesSupported string `json:"display_values_supported,omitempty"`
	// OPTIONAL . JSON array containing a list of the Claim Types that the OpenID Provider supports. These Claim Types are described in Section 5.6 of OpenID Connect Core 1.0 [OpenID.Core]. Values defined by this specification are normal, aggregated, and distributed. If omitted, the implementation supports only normal Claims.
	ClaimTypesSupported string `json:"claim_types_supported,omitempty"`
	// RECOMMENDED . JSON array containing a list of the Claim Names of the Claims that the OpenID Provider MAY be able to supply values for. Note that for privacy or other reasons, this might not be an exhaustive list.
	ClaimsSupported string `json:"claims_supported"`
	// OPTIONAL . URL of a page containing human-readable information that developers might want or need to know when using the OpenID Provider. In particular, if the OpenID Provider does not support Dynamic Client Registration, then information on how to register Clients needs to be provided in this documentation.
	ServiceDocumentation string `json:"service_documentation,omitempty"`
	// OPTIONAL . Languages and scripts supported for values in Claims being returned, represented as a JSON array of BCP47 [RFC5646] language tag values. Not all languages and scripts are necessarily supported for all Claim values.
	ClaimsLocalesSupported string `json:"claims_locales_supported,omitempty"`
	// OPTIONAL . Languages and scripts supported for the user interface, represented as a JSON array of BCP47 [RFC5646] language tag values.
	UiLocalesSupported string `json:"ui_locales_supported,omitempty"`
	// OPTIONAL . Boolean value specifying whether the OP supports use of the claims parameter, with true indicating support. If omitted, the default value is false.
	ClaimsParameterSupported string `json:"claims_parameter_supported,omitempty"`
	// OPTIONAL . Boolean value specifying whether the OP supports use of the request parameter, with true indicating support. If omitted, the default value is false.
	RequestParameterSupported string `json:"request_parameter_supported,omitempty"`
	// OPTIONAL . Boolean value specifying whether the OP supports use of the request_uri parameter, with true indicating support. If omitted, the default value is true.
	RequestUriParameterSupported string `json:"request_uri_parameter_supported,omitempty"`
	// OPTIONAL . Boolean value specifying whether the OP requires any request_uri values used to be pre-registered using the request_uris registration parameter. Pre-registration is REQUIRED when the value is true. If omitted, the default value is false.
	RequireRequestUriRegistration string `json:"require_request_uri_registration,omitempty"`
	// OPTIONAL . URL that the OpenID Provider provides to the person registering the Client to read about the OP's requirements on how the Relying Party can use the data provided by the OP. The registration process SHOULD display this URL to the person registering the Client if it is given.
	OpPolicyUri string `json:"op_policy_uri,omitempty"`
	// OPTIONAL . URL that the OpenID Provider provides to the person registering the Client to read about OpenID Provider's terms of service. The registration process SHOULD display this URL to the person registering the Client if it is given.
	OpTosUri string `json:"op_tos_uri,omitempty"`
}

type JwtHeader struct {
	Type       string `json:"typ"`
	Algorithm  string `json:"alg"`
	Thumbprint string `json:"x5t"`
	KeyId      string `json:"kid"`
}

type Provider struct {
	discovery url.URL
	metadata  *DiscoveryMetadata
	jwks      *jwk.Set
}

func (p *Provider) getMetadata(c *http.Client) error {
	r, err := c.Get(p.discovery.String())
	if err != nil {
		return err
	}

	if r.StatusCode != 200 {
		return fmt.Errorf("error retreiving discovery metadata %d: %s", r.StatusCode, r.Status)
	}

	var b []byte
	r.Body.Read(b)

	err = json.Unmarshal(b, p.metadata)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) getCertificates(c *http.Client) error {
	r, err := c.Get(p.metadata.JwksUri)
	if err != nil {
		return err
	}

	if r.StatusCode != 200 {
		return fmt.Errorf("error retreiving jwks %d: %s", r.StatusCode, r.Status)
	}

	var b []byte
	r.Body.Read(b)

	p.jwks, err = jwk.Parse(b)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) Setup(c *http.Client) {
	p.getMetadata(c)
	p.getCertificates(c)
}

type Providers struct {
	providers []*Provider
}
func (p *Providers) find(issuer string) *Provider {
	i := sort.Search(len(p.providers), func(i int) bool {
		return p.providers[i].metadata.Issuer >= issuer
	})
	if i < len(p.providers) && p.providers[i].metadata.Issuer == issuer {
		return p.providers[i]
	}

	return nil
}

func (p *Providers) setup(c *http.Client) {
	for _, provider := range p.providers {
		provider.Setup(c)
	}
	sort.Slice(p.providers, func(i, j int) bool {
		return p.providers[i].metadata.Issuer < p.providers[i].metadata.Issuer
	})
}

func (p *Providers) Add(metadataUrl url.URL) {
	p.providers = append(p.providers, &Provider{discovery: metadataUrl})
}

func Verify(providers Providers, raw string) (bool, string) {
	sectionsB64 := strings.Split(raw, ".")
	// Verify we have the expected number of JWT bits
	if len(sectionsB64) != 3 {
		return false, ""
	}

	// Decode the JWT header so we can get the Alg and Key ID / thumbprint needed for verification
	headerJson, err := base64.RawURLEncoding.DecodeString(sectionsB64[0])
	if err != nil {
		return false, ""
	}

	var header JwtHeader
	err = json.Unmarshal(headerJson, header)
	if err != nil {
		return false, ""
	}

	// Parse the token without verification so that we can find the provider that matches the issuer
	parsedToken, err := jwt.ParseString(raw, nil)
	provider := providers.find(parsedToken.Issuer())
	if provider != nil {
		// JWT signing algorithm
		var algorithm jwa.SignatureAlgorithm
		algorithm.Accept(header.Algorithm)

		// Find the key by id
		k := provider.jwks.LookupKeyID(header.KeyId)
		// or if that failed, try and find it by thumbprint instead
		if k == nil {
			for _, value := range provider.jwks.Keys {
				if value.X509CertThumbprint() == header.Thumbprint {
					k = []jwk.Key{value}
				}
			}
		}

		_, err := jwt.ParseVerify(strings.NewReader(raw), algorithm, k[0])
		if err != nil {
			return false, ""
		}

		bodyJson, _ := base64.RawURLEncoding.DecodeString(sectionsB64[1])
		return true, string(bodyJson)
	}


	return false, ""
}

func AuthMiddleware(providers Providers) mux.MiddlewareFunc {
	c := &http.Client{Timeout: time.Second * 5}
	providers.setup(c)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			authHeader := req.Header.Get("Authorization")
			headerBits := strings.Split(authHeader, " ")
			if len(headerBits) == 2 && strings.ToLower(headerBits[0]) == "bearer" {
				success, claims := Verify(providers, headerBits[1])
				if success {
					*req = *req.WithContext(context.WithValue(req.Context(), "user", claims))
				}
			}

			next.ServeHTTP(w, req)
		})
	}
}
