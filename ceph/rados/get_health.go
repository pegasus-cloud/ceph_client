package rados

import (
	"encoding/json"
	"errors"
)

//GetHealth ....
func (c *RadosCeph) GetHealth() (string, error) {
	resp, _, err := c.conn.MonCommand([]byte(`{"prefix": "status", "format": "json-pretty"}`))
	if err != nil {
		return "Health_Err", err
	}

	var clusterStatus interface{}
	err = json.Unmarshal(resp, &clusterStatus)
	if err != nil {
		return "Health_Err", err
	}
	m := clusterStatus.(map[string]interface{})
	health := m["health"]
	m = health.(map[string]interface{})
	status := m["overall_status"].(string)
	summary := m["summary"].([]interface{})
	if len(summary) == 0 {
		return status, nil
	} else {
		b, _ := json.Marshal(summary)
		return status, errors.New(string(b))
	}
}
