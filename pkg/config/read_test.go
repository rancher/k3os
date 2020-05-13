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

func TestAuthorizedKeys(t *testing.T) {
	c1 := map[string]interface{}{
		"ssh_authorized_keys": []string{
			"one...",
		},
	}
	c2 := map[string]interface{}{
		"ssh_authorized_keys": []string{
			"two...",
		},
	}
	cc, err := readersToObject(
		func() (map[string]interface{}, error) {
			return c1, nil
		},
		func() (map[string]interface{}, error) {
			return c2, nil
		},
	)
	if len(cc.SSHAuthorizedKeys) != 1 {
		t.Fatal(err, "got %d keys, expected 2", len(cc.SSHAuthorizedKeys))
	}
	if err != nil {
		t.Fatal(err)
	}
}
