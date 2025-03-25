package listener

import (
	_ "github.com/external-secrets-inc/reloader/internal/listener/eventgrid"
	_ "github.com/external-secrets-inc/reloader/internal/listener/hashivault"
	_ "github.com/external-secrets-inc/reloader/internal/listener/k8ssecret"
	_ "github.com/external-secrets-inc/reloader/internal/listener/mock"
	_ "github.com/external-secrets-inc/reloader/internal/listener/pubsub"
	_ "github.com/external-secrets-inc/reloader/internal/listener/sqs"
	_ "github.com/external-secrets-inc/reloader/internal/listener/tcp"
	_ "github.com/external-secrets-inc/reloader/internal/listener/webhook"
)
