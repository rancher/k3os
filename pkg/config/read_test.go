package config

import "testing"

func TestDataSource(t *testing.T) {
	cc, err := readersToObject(func() (map[string]interface{}, error) {
		return map[string]interface{}{
			"k3os": map[string]interface{}{
				"datasource": "foo",
			},
		}, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(cc.K3OS.DataSources) != 1 {
		t.Fatal("no datasources")
	}
	if cc.K3OS.DataSources[0] != "foo" {
		t.Fatalf("%s != foo", cc.K3OS.DataSources[0])
	}
}
