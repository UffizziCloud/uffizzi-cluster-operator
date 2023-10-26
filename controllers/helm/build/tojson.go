package build

import (
	"encoding/json"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func HelmValuesToJSON(helmValues any) v1.JSON {
	// marshal HelmValues struct to JSON
	helmValuesRaw, _ := json.Marshal(helmValues)
	//if err != nil {
	//	return nil, errors.Wrap(err, "failed to marshal HelmValues struct to JSON")
	//}

	// Create the apiextensionsv1.JSON instance with the raw data
	helmValuesJSONObj := v1.JSON{Raw: helmValuesRaw}
	return helmValuesJSONObj
}
