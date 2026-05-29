/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BlackboxModuleSpec defines the desired state of BlackboxModule
type BlackboxModuleSpec struct {
	// +kubebuilder:validation:Enum=http;tcp;icmp;dns;grpc
	Prober string `json:"prober" yaml:"prober"`
	// +kubebuilder:default:="5s"
	// +optional
	Timeout string `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	// +optional
	HTTP *HTTPProbe `json:"http,omitempty" yaml:"http,omitempty"`
	// +optional
	TCP *TCPProbe `json:"tcp,omitempty" yaml:"tcp,omitempty"`
	// +optional
	ICMP *ICMPProbe `json:"icmp,omitempty" yaml:"icmp,omitempty"`
	// +optional
	DNS *DNSProbe `json:"dns,omitempty" yaml:"dns,omitempty"`
	// +optional
	GRPC *GRPCProbe `json:"grpc,omitempty" yaml:"grpc,omitempty"`
}

type HTTPProbe struct {
	// +kubebuilder:default={200}
	// +optional
	ValidStatusCodes []int `json:"validStatusCodes,omitempty" yaml:"valid_status_codes,omitempty"`
	// +kubebuilder:default={HTTP/1.1,HTTP/2}
	// +optional
	ValidHTTPVersions []string `json:"validHttpVersions,omitempty" yaml:"valid_http_versions,omitempty"`
	// +kubebuilder:default=GET
	// +optional
	Method string `json:"method,omitempty" yaml:"method,omitempty"`
	// +optional
	Headers map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
	// +kubebuilder:default=1048576
	// +optional
	BodySizeLimit int `json:"bodySizeLimit,omitempty" yaml:"body_size_limit,omitempty"`
	// +kubebuilder:default=none
	// +kubebuilder:validation:Enum=none;gzip;deflate
	// +optional
	Compression string `json:"compression,omitempty" yaml:"compression,omitempty"`
	// +kubebuilder:default=true
	// +optional
	FollowRedirects *bool `json:"followRedirects,omitempty" yaml:"follow_redirects,omitempty"`
	// +optional
	FailIfSSL *bool `json:"failIfSsl,omitempty" yaml:"fail_if_ssl,omitempty"`
	// +optional
	FailIfNotSSL *bool `json:"failIfNotSsl,omitempty" yaml:"fail_if_not_ssl,omitempty"`
	// +optional
	FailIfBodyJSONMatchesCEL string `json:"failIfBodyJsonMatchesCel,omitempty" yaml:"fail_if_body_matches_cel,omitempty"`
	// +optional
	FailIfBodyJSONNotMatchesCEL string `json:"failIfBodyJsonNotMatchesCel,omitempty" yaml:"fail_if_body_not_matches_cel,omitempty"`
	// +optional
	FailIfBodyMatchesRegexp []string `json:"failIfBodyMatchesRegexp,omitempty" yaml:"fail_if_body_matches_regexp,omitempty"`
	// +optional
	FailIfBodyNotMatchesRegexp []string `json:"failIfBodyNotMatchesRegexp,omitempty" yaml:"fail_if_body_not_matches_regexp,omitempty"`
	// +optional
	FailIfHeaderMatchesRegexp []HTTPHeaderMatch `json:"failIfHeaderMatchesRegexp,omitempty" yaml:"fail_if_header_matches_regexp,omitempty"`
	// +optional
	FailIfHeaderNotMatchesRegexp []HTTPHeaderMatch `json:"failIfHeaderNotMatchesRegexp,omitempty" yaml:"fail_if_header_not_matches_regexp,omitempty"`
	// +optional
	TLSConfig *TLSConfig `json:"tlsConfig,omitempty" yaml:"tls_config,omitempty"`
	// +optional
	BasicAuth *BasicAuth `json:"basicAuth,omitempty" yaml:"basic_auth,omitempty"`
	// +optional
	Authorization *Authorization `json:"authorization,omitempty" yaml:"authorization,omitempty"`
	// +optional
	ProxyURL string `json:"proxyUrl,omitempty" yaml:"proxy_url,omitempty"`
	// +optional
	NoProxy string `json:"noProxy,omitempty" yaml:"no_proxy,omitempty"`
	// +kubebuilder:default=false
	// +optional
	ProxyFromEnvironment *bool `json:"proxyFromEnvironment,omitempty" yaml:"proxy_from_environment,omitempty"`
	// +optional
	ProxyConnectHeaders map[string]string `json:"proxyConnectHeaders,omitempty" yaml:"proxy_connect_headers,omitempty"`
	// +kubebuilder:default=false
	// +optional
	SkipResolvePhaseWithProxy *bool `json:"skipResolvePhaseWithProxy,omitempty" yaml:"skip_resolve_phase_with_proxy,omitempty"`
	// +optional
	OAuth2 *OAuth2Config `json:"oauth2,omitempty" yaml:"oauth2,omitempty"`
	// +kubebuilder:default=true
	// +optional
	IPProtocolFallback *bool `json:"ipProtocolFallback,omitempty" yaml:"ip_protocol_fallback,omitempty"`
	// +optional
	Body string `json:"body,omitempty" yaml:"body,omitempty"`
}

type BasicAuth struct {
	// +optional
	Username string `json:"username,omitempty" yaml:"username,omitempty"`
	// +optional
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
	// +optional
	PasswordFile string `json:"passwordFile,omitempty" yaml:"password_file,omitempty"`
}

type Authorization struct {
	// +kubebuilder:validation:Enum=Bearer;Basic
	// +optional
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
	// +optional
	Credentials string `json:"credentials,omitempty" yaml:"credentials,omitempty"`
	// +optional
	CredentialsFile string `json:"credentialsFile,omitempty" yaml:"credentials_file,omitempty"`
}

type OAuth2Config struct {
	// +optional
	ClientID string `json:"clientId,omitempty" yaml:"client_id,omitempty"`
	// +optional
	ClientSecret string `json:"clientSecret,omitempty" yaml:"client_secret,omitempty"`
	// +optional
	ClientSecretFile string `json:"clientSecretFile,omitempty" yaml:"client_secret_file,omitempty"`
	// +optional
	ClientSecretRef string `json:"clientSecretRef,omitempty" yaml:"client_secret_ref,omitempty"`
	// +optional
	Scopes []string `json:"scopes,omitempty" yaml:"scopes,omitempty"`
	// +optional
	TokenURL string `json:"tokenUrl,omitempty" yaml:"token_url,omitempty"`
	// +optional
	EndpointParams map[string]string `json:"endpointParams,omitempty" yaml:"endpoint_params,omitempty"`
	// +optional
	TLSConfig *TLSConfig `json:"tlsConfig,omitempty" yaml:"tls_config,omitempty"`
	// +optional
	ProxyURL string `json:"proxyUrl,omitempty" yaml:"proxy_url,omitempty"`
	// +optional
	NoProxy string `json:"noProxy,omitempty" yaml:"no_proxy,omitempty"`
}

type HTTPHeaderMatch struct {
	Header string `json:"header" yaml:"header"`
	Regexp string `json:"regexp" yaml:"regexp"`
	// +kubebuilder:default=false
	// +optional
	AllowMissing bool `json:"allowMissing,omitempty" yaml:"allow_missing,omitempty"`
}

type TLSConfig struct {
	// +kubebuilder:default=false
	// +optional
	InsecureSkipVerify *bool `json:"insecureSkipVerify,omitempty" yaml:"insecure_skip_verify,omitempty"`
	// +optional
	CAFile string `json:"caFile,omitempty" yaml:"ca_file,omitempty"`
	// +optional
	CertFile string `json:"certFile,omitempty" yaml:"cert_file,omitempty"`
	// +optional
	KeyFile string `json:"keyFile,omitempty" yaml:"key_file,omitempty"`
	// +optional
	ServerName string `json:"serverName,omitempty" yaml:"server_name,omitempty"`
}

type TCPProbe struct {
	// +kubebuilder:validation:Enum=ip4;ip6
	// +kubebuilder:default=ip4
	// +optional
	PreferredIPProtocol string `json:"preferredIpProtocol,omitempty" yaml:"preferred_ip_protocol,omitempty"`
	// +optional
	SourceIPAddress string `json:"sourceIpAddress,omitempty" yaml:"source_ip_address,omitempty"`
	// +optional
	QueryResponse []QueryResponseEntry `json:"queryResponse,omitempty" yaml:"query_response,omitempty"`
	// +kubebuilder:default=false
	// +optional
	TLS *bool `json:"tls,omitempty" yaml:"tls,omitempty"`
	// +optional
	TLSConfig *TLSConfig `json:"tlsConfig,omitempty" yaml:"tls_config,omitempty"`
}

type QueryResponseEntry struct {
	// +optional
	Expect string `json:"expect,omitempty" yaml:"expect,omitempty"`
	// +optional
	Send string `json:"send,omitempty" yaml:"send,omitempty"`
	// +kubebuilder:default=false
	// +optional
	StartTLS *bool `json:"startTls,omitempty" yaml:"starttls,omitempty"`
}

type DNSProbe struct {
	// +kubebuilder:validation:Enum=ip4;ip6
	// +kubebuilder:default=ip4
	// +optional
	PreferredIPProtocol string `json:"preferredIpProtocol,omitempty" yaml:"preferred_ip_protocol,omitempty"`
	// +optional
	SourceIPAddress string `json:"sourceIpAddress,omitempty" yaml:"source_ip_address,omitempty"`
	// +kubebuilder:validation:Enum=udp;tcp
	// +kubebuilder:default=udp
	// +optional
	TransportProtocol string `json:"transportProtocol,omitempty" yaml:"transport_protocol,omitempty"`
	// Required
	QueryName string `json:"queryName" yaml:"query_name"`
	// +kubebuilder:default=A
	// +optional
	QueryType string `json:"queryType,omitempty" yaml:"query_type,omitempty"`
	// +optional
	ValidRCodes []string `json:"validRCodes,omitempty" yaml:"valid_rcodes,omitempty"`
	// +optional
	ValidateAnswerRRs *DNSValidateRR `json:"validateAnswerRRs,omitempty" yaml:"validate_answer_rrs,omitempty"`
	// +optional
	ValidateAuthorityRRs *DNSValidateRR `json:"validateAuthorityRRs,omitempty" yaml:"validate_authority_rrs,omitempty"`
	// +optional
	ValidateAdditionalRRs *DNSValidateRR `json:"validateAdditionalRRs,omitempty" yaml:"validate_additional_rrs,omitempty"`
}

type DNSValidateRR struct {
	// +optional
	FailIfMatchesRegexp []string `json:"failIfMatchesRegexp,omitempty" yaml:"fail_if_matches_regexp,omitempty"`
	// +optional
	FailIfNotMatchesRegexp []string `json:"failIfNotMatchesRegexp,omitempty" yaml:"fail_if_not_matches_regexp,omitempty"`
	// +optional
	FailIfAllMatchRegexp []string `json:"failIfAllMatchRegexp,omitempty" yaml:"fail_if_all_match_regexp,omitempty"`
	// +optional
	FailIfNoneMatchRegexp []string `json:"failIfNoneMatchRegexp,omitempty" yaml:"fail_if_none_match_regexp,omitempty"`
}

type ICMPProbe struct {
	// +kubebuilder:validation:Enum=ip4;ip6
	// +kubebuilder:default=ip4
	// +optional
	PreferredIPProtocol string `json:"preferredIpProtocol,omitempty" yaml:"preferred_ip_protocol,omitempty"`
	// +optional
	SourceIPAddress string `json:"sourceIpAddress,omitempty" yaml:"source_ip_address,omitempty"`
	// +kubebuilder:default=24
	// +optional
	PayloadSize int `json:"payloadSize,omitempty" yaml:"payload_size,omitempty"`
	// +kubebuilder:default=false
	// +optional
	DontFragment *bool `json:"dontFragment,omitempty" yaml:"dont_fragment,omitempty"`
}

type GRPCProbe struct {
	// +optional
	Service string `json:"service,omitempty" yaml:"service,omitempty"`
	// +kubebuilder:validation:Enum=ip4;ip6
	// +kubebuilder:default=ip4
	// +optional
	PreferredIPProtocol string `json:"preferredIpProtocol,omitempty" yaml:"preferred_ip_protocol,omitempty"`
	// +kubebuilder:default=true
	// +optional
	IPProtocolFallback *bool `json:"ipProtocolFallback,omitempty" yaml:"ip_protocol_fallback,omitempty"`
	// +kubebuilder:default=false
	// +optional
	UseTLS *bool `json:"useTls,omitempty" yaml:"tls,omitempty"`
	// +optional
	TLSConfig *TLSConfig `json:"tlsConfig,omitempty" yaml:"tls_config,omitempty"`
}

// BlackboxModuleStatus defines the observed state of BlackboxModule
type BlackboxModuleStatus struct {
	//+operator-sdk:csv:customresourcedefinitions:type=status
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// BlackboxModule is the Schema for the blackboxmodules API
type BlackboxModule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BlackboxModuleSpec   `json:"spec,omitempty"`
	Status BlackboxModuleStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BlackboxModuleList contains a list of BlackboxModule
type BlackboxModuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BlackboxModule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BlackboxModule{}, &BlackboxModuleList{})
}
