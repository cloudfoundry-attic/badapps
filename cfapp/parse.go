package cfapp

import "encoding/json"

type Info struct {
	ApplicationID   string         `json:"application_id"`
	ApplicationName string         `json:"application_name"`
	ApplicationUris []string       `json:"application_uris"`
	InstanceID      string         `json:"instance_id"`
	Limits          map[string]int `json:"limits"`
	SpaceId         string         `json:"space_id"`
	SpaceName       string         `json:"space_name"`
}

func Parse(vcapApplication string) (Info, error) {
	var info Info
	if err := json.Unmarshal([]byte(vcapApplication), &info); err != nil {
		return Info{}, err
	}
	return info, nil
}
