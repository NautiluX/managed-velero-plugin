package k8s

import "encoding/json"

type GenericObject struct {
	Kind     string            `json:"kind"`
	Metadata GenericObjectMeta `json:"metadata"`
}
type GenericObjectMeta struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

func GetGenericObject(cr string) (GenericObject, error) {
	genericObject := GenericObject{}
	err := json.Unmarshal([]byte(cr), &genericObject)
	if err != nil {
		return GenericObject{}, err
	}
	return genericObject, nil
}
