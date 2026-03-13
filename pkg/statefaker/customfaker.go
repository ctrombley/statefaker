package statefaker

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"reflect"
)

func GenerateAttributes(resourceType string) (json.RawMessage, error) {
	// Generate realistic resource attributes based on resource type
	var generator func() map[string]any

	switch resourceType {
	case "aws_s3_bucket":
		generator = generateS3BucketAttributes
	case "aws_iam_user":
		generator = generateIAMUserAttributes
	case "aws_instance":
		generator = generateEC2InstanceAttributes
	case "aws_lambda_function":
		generator = generateLambdaFunctionAttributes
	case "aws_db_instance":
		generator = generateRDSInstanceAttributes
	case "aws_api_gateway_rest_api":
		generator = generateAPIGatewayRestAPIAttributes
	default:
		// Fallback for types without specific generators
		generator = func() map[string]any {
			// Generate minimal valid attributes based on ID only
			// Most resources have an ID, and extra unknown attributes can cause errors
			name := generateResourceName()
			return map[string]any{
				"id": name,
			}
		}
	}

	resourceAttributes := generator()

	b, err := json.Marshal(resourceAttributes)
	if err != nil {
		return nil, err
	}

	return json.RawMessage(b), nil
}

func tfidentityProvider(v reflect.Value) (any, error) {
	// Always return nil for identity as most resources don't support it
	return nil, nil
}

func tfprivateProvider(v reflect.Value) (any, error) {
	// Occasionally generate a non-empty private field
	if rand.IntN(5) == 0 {
		// Generate some random bytes and encode as a base64 string
		bytes := make([]byte, 16)
		for i := range bytes {
			bytes[i] = byte(rand.IntN(256))
		}
		return base64.StdEncoding.EncodeToString(bytes), nil
	}
	return "", nil
}

func tfdependenciesProvider(v reflect.Value) (any, error) {
	// Generate a list of dependencies (0-3) for the resource
	numDeps := rand.IntN(4)
	dependencies := make([]string, numDeps)

	for i := 0; i < numDeps; i++ {
		resourceType := generateResourceType()
		resourceName := generateResourceName()
		var moduleAddress string
		if rand.IntN(10) < 3 {
			moduleAddress = generateModuleAddress()
		}
		dep := fmt.Sprintf("%s.%s", resourceType, resourceName)
		if moduleAddress != "" {
			dep = fmt.Sprintf("%s.%s", moduleAddress, dep)
		}
		dependencies[i] = dep
	}

	return dependencies, nil
}

func tfemptystringsliceProvider(v reflect.Value) (any, error) {
	// Always return an empty string slice
	return []string{}, nil
}
