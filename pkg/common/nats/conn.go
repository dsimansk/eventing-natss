/*
Copyright 2021 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package nats

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientsetcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/system"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"

	commonconfig "knative.dev/eventing-natss/pkg/common/config"
	"knative.dev/eventing-natss/pkg/common/constants"
)

var (
	ErrBadCredentialFileOption = errors.New("bad auth.credentialFile option")
	ErrBadMTLSOption           = errors.New("bad auth.tls option")
	ErrBadTLSOption            = errors.New("bad tls option")
)

func NewNatsConn(ctx context.Context, config commonconfig.EventingNatsConfig) (*nats.Conn, error) {
	url := config.URL
	if url == "" {
		url = constants.DefaultNatsURL
	}

	coreV1Client, err := clientsetcorev1.NewForConfig(injection.GetConfig(ctx))
	if err != nil {
		return nil, err
	}

	secrets := coreV1Client.Secrets(getNamespace(ctx))

	var opts []nats.Option

	if config.Auth != nil {
		o, err := buildAuthOption(ctx, *config.Auth, secrets)
		if err != nil {
			return nil, err
		}

		opts = append(opts, o...)
	}

	if config.RootCA != nil {
		o, err := buildRootCAOption(ctx, *config.RootCA, secrets)
		if err != nil {
			return nil, err
		}

		opts = append(opts, o)
	}

	return nats.Connect(url, opts...)
}

func buildAuthOption(ctx context.Context, config commonconfig.ENConfigAuth, secrets clientsetcorev1.SecretInterface) ([]nats.Option, error) {
	opts := make([]nats.Option, 0, 2)
	if config.CredentialFile != nil {
		o, err := buildCredentialFileOption(ctx, *config.CredentialFile, secrets)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrBadCredentialFileOption, err.Error())
		}

		opts = append(opts, o)
	}

	if config.TLS != nil {
		o, err := buildClientTLSOption(ctx, *config.TLS, secrets)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrBadMTLSOption, err.Error())
		}

		opts = append(opts, o)
	}

	return opts, nil
}

func buildCredentialFileOption(ctx context.Context, config commonconfig.ENConfigAuthCredentialFile, secrets clientsetcorev1.SecretInterface) (nats.Option, error) {
	if config.Secret == nil {
		return nil, nil
	}

	secret, err := secrets.Get(ctx, config.Secret.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	contents, err := loadCredentialFileSecret(*config.Secret, secret)
	if err != nil {
		return nil, err
	}

	return credentialFileOption(contents), nil
}

func buildClientTLSOption(ctx context.Context, config commonconfig.ENConfigAuthTLS, secrets clientsetcorev1.SecretInterface) (nats.Option, error) {
	if config.Secret == nil {
		return nil, nil
	}

	secret, err := secrets.Get(ctx, config.Secret.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return ClientCert(secret), nil
}

func buildRootCAOption(ctx context.Context, config commonconfig.ENConfigRootCA, secrets clientsetcorev1.SecretInterface) (nats.Option, error) {
	var (
		decoded []byte
		err     error
	)

	if config.CABundle != "" {
		decoded, err = base64.StdEncoding.DecodeString(config.CABundle)
		if err != nil {
			return nil, err
		}
	} else if config.Secret != nil {
		secret, err := secrets.Get(ctx, config.Secret.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		var ok bool
		if decoded, ok = secret.Data[TLSCaCertKey]; !ok {
			return nil, ErrTLSCaCertMissing
		}
	}

	return RootCA(decoded), nil
}

func loadCredentialFileSecret(config v1.SecretKeySelector, secret *v1.Secret) ([]byte, error) {
	key := constants.DefaultCredentialFileSecretKey
	if config.Key != "" {
		key = config.Key
	}

	contents, ok := secret.Data[key]
	if !ok {
		return nil, fmt.Errorf("failed to load secret, key does not exist: %s", key)
	}

	return contents, nil
}

// credentialFileOption processes the raw credential file contents and returns the nats.Option. This logic has been
// derived from the nats.UserCredentials() function but modified for when the file has already been parsed.
func credentialFileOption(contents []byte) nats.Option {
	userCB := func() (string, error) {
		return nkeys.ParseDecoratedJWT(contents)
	}

	sigCB := func(nonce []byte) ([]byte, error) {
		// nkeys.KeyPair, error
		kp, err := nkeys.ParseDecoratedNKey(contents)
		if err != nil {
			return nil, err
		}
		// Wipe our key on exit.
		defer kp.Wipe()

		return kp.Sign(nonce)
	}

	return nats.UserJWT(userCB, sigCB)
}

func getNamespace(ctx context.Context) string {
	if injection.HasNamespaceScope(ctx) {
		return injection.GetNamespaceScope(ctx)
	}

	return system.Namespace()
}
