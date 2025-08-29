package handler

import (
	_ "github.com/external-secrets-inc/reloader/internal/handler/deployment"
	_ "github.com/external-secrets-inc/reloader/internal/handler/externalsecret"
	_ "github.com/external-secrets-inc/reloader/internal/handler/pushsecret"
	_ "github.com/external-secrets-inc/reloader/internal/handler/workflow"
)
