/*
Copyright 2022 The Knative Authors

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
package dispatcher

import (
	"testing"

	"knative.dev/eventing-natss/pkg/apis/messaging/v1alpha1"
	"knative.dev/eventing-natss/pkg/channel/jetstream/utils"
	reconciletesting "knative.dev/eventing-natss/pkg/reconciler/testing"
	"knative.dev/eventing/pkg/channel"
)

func TestDispatcher_RegisterChannelHost(t *testing.T) {
	nc := reconciletesting.NewNatsJetStreamChannel(testNS, ncName)
	config := createChannelConfig(nc)

	d := &Dispatcher{}

	err := d.RegisterChannelHost(*config)
	if err != nil {
		t.Fatal(err)
	}
}

func createChannelConfig(nc *v1alpha1.NatsJetStreamChannel, subs ...Subscription) *ChannelConfig {
	if subs == nil {
		subs = []Subscription{}
	}
	return &ChannelConfig{
		ChannelReference: channel.ChannelReference{
			Namespace: nc.Namespace,
			Name:      nc.Name,
		},
		StreamName:             utils.StreamName(nc),
		HostName:               "a.b.c.d",
		ConsumerConfigTemplate: nc.Spec.ConsumerConfigTemplate,
		Subscriptions:          subs,
	}
}
