package build

import (
	"encoding/json"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func HelmValuesToJSON(helmValues any) (v1.JSON, error) {
	// marshal HelmValues struct to JSON
	helmValuesRaw, err := json.Marshal(helmValues)
	// Create the apiextensionsv1.JSON instance with the raw data
	helmValuesJSONObj := v1.JSON{Raw: helmValuesRaw}
	return helmValuesJSONObj, err
}
