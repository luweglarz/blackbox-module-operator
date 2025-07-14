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
	Prober string `json:"prober"`
	// +kubebuilder:default:="5s"
	// +optional
	Timeout string `json:"timeout,omitempty"`
	// +optional
	HTTP *HTTPProbe `json:"http,omitempty"`
	// +optional
	TCP *TCPProbe `json:"tcp,omitempty"`
	// +optional
	ICMP *ICMPProbe `json:"icmp,omitempty"`
	// +optional
	DNS *DNSProbe `json:"dns,omitempty"`
	// +optional
	GRPC *GRPCProbe `json:"grpc,omitempty"`
}

type HTTPProbe struct {
	// +kubebuilder:default={200}
	// +optional
	ValidStatusCodes []int `json:"validStatusCodes,omitempty"`
	// +kubebuilder:default={HTTP/1.1,HTTP/2}
	// +optional
	ValidHTTPVersions []string `json:"validHttpVersions,omitempty"`
	// +kubebuilder:default=GET
	// +optional
	Method string `json:"method,omitempty"`
	// +optional
	Headers map[string]string `json:"headers,omitempty"`
	// +kubebuilder:default=1048576
	// +optional
	BodySizeLimit int `json:"bodySizeLimit,omitempty"`
	// +kubebuilder:default=none
	// +kubebuilder:validation:Enum=none;gzip;deflate
	// +optional
	Compression string `json:"compression,omitempty"`
	// +kubebuilder:default=true
	// +optional
	FollowRedirects *bool `json:"followRedirects,omitempty"`
	// +optional
	FailIfSSL *bool `json:"failIfSsl,omitempty"`
	// +optional
	FailIfNotSSL *bool `json:"failIfNotSsl,omitempty"`
	// +optional
	FailIfBodyJSONMatchesCEL string `json:"failIfBodyJsonMatchesCel,omitempty"`
	// +optional
	FailIfBodyJSONNotMatchesCEL string `json:"failIfBodyJsonNotMatchesCel,omitempty"`
	// +optional
	FailIfBodyMatchesRegexp []string `json:"failIfBodyMatchesRegexp,omitempty"`
	// +optional
	FailIfBodyNotMatchesRegexp []string `json:"failIfBodyNotMatchesRegexp,omitempty"`
	// +optional
	FailIfHeaderMatchesRegexp []HTTPHeaderMatch `json:"failIfHeaderMatchesRegexp,omitempty"`
	// +optional
	FailIfHeaderNotMatchesRegexp []HTTPHeaderMatch `json:"failIfHeaderNotMatchesRegexp,omitempty"`
	// +optional
	TLSConfig *TLSConfig `json:"tlsConfig,omitempty"`
	// +optional
	BasicAuth *BasicAuth `json:"basicAuth,omitempty"`
	// +optional
	Authorization *Authorization `json:"authorization,omitempty"`
	// +optional
	ProxyURL string `json:"proxyUrl,omitempty"`
	// +optional
	NoProxy string `json:"noProxy,omitempty"`
	// +kubebuilder:default=false
	// +optional
	ProxyFromEnvironment *bool `json:"proxyFromEnvironment,omitempty"`
	// +optional
	ProxyConnectHeaders map[string]string `json:"proxyConnectHeaders,omitempty"`
	// +kubebuilder:default=false
	// +optional
	SkipResolvePhaseWithProxy *bool `json:"skipResolvePhaseWithProxy,omitempty"`
	// +optional
	OAuth2 *OAuth2Config `json:"oauth2,omitempty"`
	// +kubebuilder:default=true
	// +optional
	IPProtocolFallback *bool `json:"ipProtocolFallback,omitempty"`
	// +optional
	Body string `json:"body,omitempty"`
}

type BasicAuth struct {
	// +optional
	Username string `json:"username,omitempty"`
	// +optional
	Password string `json:"password,omitempty"`
	// +optional
	PasswordFile string `json:"passwordFile,omitempty"`
}

type Authorization struct {
	// +kubebuilder:validation:Enum=Bearer;Basic
	// +optional
	Type string `json:"type,omitempty"`
	// +optional
	Credentials string `json:"credentials,omitempty"`
	// +optional
	CredentialsFile string `json:"credentialsFile,omitempty"`
}

type OAuth2Config struct {
	// +optional
	ClientID string `json:"clientId,omitempty"`
	// +optional
	ClientSecret string `json:"clientSecret,omitempty"`
	// +optional
	ClientSecretFile string `json:"clientSecretFile,omitempty"`
	// +optional
	ClientSecretRef string `json:"clientSecretRef,omitempty"`
	// +optional
	Scopes []string `json:"scopes,omitempty"`
	// +optional
	TokenURL string `json:"tokenUrl,omitempty"`
	// +optional
	EndpointParams map[string]string `json:"endpointParams,omitempty"`
	// +optional
	TLSConfig *TLSConfig `json:"tlsConfig,omitempty"`
	// +optional
	ProxyURL string `json:"proxyUrl,omitempty"`
	// +optional
	NoProxy string `json:"noProxy,omitempty"`
}

type HTTPHeaderMatch struct {
	Header string `json:"header"`
	Regexp string `json:"regexp"`
	// +kubebuilder:default=false
	// +optional
	AllowMissing bool `json:"allowMissing,omitempty"`
}

type TLSConfig struct {
	// +kubebuilder:default=false
	// +optional
	InsecureSkipVerify *bool `json:"insecureSkipVerify,omitempty"`
	// +optional
	CAFile string `json:"caFile,omitempty"`
	// +optional
	CertFile string `json:"certFile,omitempty"`
	// +optional
	KeyFile string `json:"keyFile,omitempty"`
	// +optional
	ServerName string `json:"serverName,omitempty"`
}

type TCPProbe struct {
	// +kubebuilder:validation:Enum=ip4;ip6
	// +kubebuilder:default=ip4
	// +optional
	PreferredIPProtocol string `json:"preferredIpProtocol,omitempty"`
	// +optional
	SourceIPAddress string `json:"sourceIpAddress,omitempty"`
	// +optional
	QueryResponse []QueryResponseEntry `json:"queryResponse,omitempty"`
	// +kubebuilder:default=false
	// +optional
	TLS *bool `json:"tls,omitempty"`
	// +optional
	TLSConfig *TLSConfig `json:"tlsConfig,omitempty"`
}

type QueryResponseEntry struct {
	// +optional
	Expect string `json:"expect,omitempty"`
	// +optional
	Send string `json:"send,omitempty"`
	// +kubebuilder:default=false
	// +optional
	StartTLS *bool `json:"startTls,omitempty"`
}

type DNSProbe struct {
	// +kubebuilder:validation:Enum=ip4;ip6
	// +kubebuilder:default=ip4
	// +optional
	PreferredIPProtocol string `json:"preferredIpProtocol,omitempty"`
	// +optional
	SourceIPAddress string `json:"sourceIpAddress,omitempty"`
	// +kubebuilder:validation:Enum=udp;tcp
	// +kubebuilder:default=udp
	// +optional
	TransportProtocol string `json:"transportProtocol,omitempty"`
	// Required
	QueryName string `json:"queryName"`
	// +kubebuilder:default=A
	// +optional
	QueryType string `json:"queryType,omitempty"`
	// +optional
	ValidRCodes []string `json:"validRCodes,omitempty"`
	// +optional
	ValidateAnswerRRs *DNSValidateRR `json:"validateAnswerRRs,omitempty"`
	// +optional
	ValidateAuthorityRRs *DNSValidateRR `json:"validateAuthorityRRs,omitempty"`
	// +optional
	ValidateAdditionalRRs *DNSValidateRR `json:"validateAdditionalRRs,omitempty"`
}

type DNSValidateRR struct {
	// +optional
	FailIfMatchesRegexp []string `json:"failIfMatchesRegexp,omitempty"`
	// +optional
	FailIfNotMatchesRegexp []string `json:"failIfNotMatchesRegexp,omitempty"`
	// +optional
	FailIfAllMatchRegexp []string `json:"failIfAllMatchRegexp,omitempty"`
	// +optional
	FailIfNoneMatchRegexp []string `json:"failIfNoneMatchRegexp,omitempty"`
}

type ICMPProbe struct {
	// +kubebuilder:validation:Enum=ip4;ip6
	// +kubebuilder:default=ip4
	// +optional
	PreferredIPProtocol string `json:"preferredIpProtocol,omitempty"`
	// +optional
	SourceIPAddress string `json:"sourceIpAddress,omitempty"`
	// +kubebuilder:default=24
	// +optional
	PayloadSize int `json:"payloadSize,omitempty"`
	// +kubebuilder:default=false
	// +optional
	DontFragment *bool `json:"dontFragment,omitempty"`
}

type GRPCProbe struct {
	// +optional
	Service string `json:"service,omitempty"`
	// +kubebuilder:validation:Enum=ip4;ip6
	// +kubebuilder:default=ip4
	// +optional
	PreferredIPProtocol string `json:"preferredIpProtocol,omitempty"`
	// +kubebuilder:default=true
	// +optional
	IPProtocolFallback *bool `json:"ipProtocolFallback,omitempty"`
	// +kubebuilder:default=false
	// +optional
	UseTLS *bool `json:"useTls,omitempty"`
	// +optional
	TLSConfig *TLSConfig `json:"tlsConfig,omitempty"`
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
