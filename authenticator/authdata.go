package main

import pb "github.com/approzium/approzium/authenticator/protos"

const (
	KeyAuthType    = "auth_type"
	ValAuthTypeAWS = "aws"

	KeyClientLang       = "client_language"
	ValClientLangGo     = "go"
	ValClientLangPython = "python"

	// AWS only.
	KeySignedGetCallerIdentity = "signed_get_caller_identity"
	KeyClaimedIamArn           = "claimed_iam_arn"
)

func toMap(authData []*pb.AuthData) map[string]string {
	m := make(map[string]string, len(authData))
	for _, kv := range authData {
		m[kv.Key] = kv.Value
	}
	return m
}
