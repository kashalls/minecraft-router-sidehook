package unifi

type UnifiNetworkObject struct {
	Meta struct {
		RC string `json:"rc"`
	} `json:"meta"`
	Data []NetworkGroup `json:"data"`
}

type NetworkGroup struct {
	GroupMembers []string `json:"group_members"`
	Name         string   `json:"name"`
	SiteID       string   `json:"site_id,omitempty"`
	ID           string   `json:"_id,omitempty"`
	GroupType    string   `json:"group_type"`
}
