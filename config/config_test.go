package config

import "testing"

func TestReadJson(t *testing.T) {
	invalidInput := "bla"
	validInput := `{ "server": { "port": 80 }}`

	cfg, err := ReadJson([]byte(invalidInput))
	if err == nil {
		t.Log("Error should be present on invalid input.")
		t.Fail()
	}
	if cfg != nil {
		t.Log("An error should have occured and no configuration returned.")
		t.Fail()
	}

	cfg, err = ReadJson([]byte(validInput))
	if err != nil {
		t.Log("No error should be present.")
		t.Fail()
	}
	if cfg == nil {
		t.Log("There should be configuration parsed.")
		t.FailNow()
	}
	if cfg.Server.Port != 80 {
		t.Logf("The server.port should be set to 80 but it is %d", cfg.Server.Port)
		t.Fail()
	}
}

func TestReadYaml(t *testing.T) {
	invalidInput := "bla"
	validInput := `{ "server": { "port": 80 }}` // don't know how this works but cool I guess
	anotherValidInput := `
server:
  port: 80`

	cfg, err := ReadYaml([]byte(invalidInput))
	if err == nil {
		t.Log("Error should be present on invalid input.")
		t.Fail()
	}
	if cfg != nil {
		t.Log("An error should have occured and no configuration returned.")
		t.Fail()
	}

	cfg, err = ReadYaml([]byte(validInput))
	if err != nil {
		t.Logf("No error should be present but got %v\n", err)
		t.Fail()
	}
	if cfg == nil {
		t.Log("There should be configuration parsed.")
		t.Fail()
	} else if cfg.Server.Port != 80 {
		t.Logf("The server.port should be set to 80 but it is %d", cfg.Server.Port)
		t.Fail()
	}

	cfg, err = ReadYaml([]byte(anotherValidInput))
	if err != nil {
		t.Logf("No error should be present but got %v\n", err)
		t.Fail()
	}
	if cfg == nil {
		t.Log("There should be configuration parsed.")
		t.Fail()
	} else if cfg.Server.Port != 80 {
		t.Logf("The server.port should be set to 80 but it is %d\n", cfg.Server.Port)
		t.Fail()
	}
}

func TestAsJson(t *testing.T) {
	cfg := &Config{
		Server: &ServerConfig{
			Port: 80,
		},
		Checks: []CheckConfig{
			{
				Name:    "wow",
				Enabled: false,
			},
		},
	}
	resShort := `{"server":{"port":80},"checks":[{"name":"wow","enabled":false}]}`
	resLong := `{
  "server": {
    "port": 80
  },
  "checks": [
    {
      "name": "wow",
      "enabled": false
    }
  ]
}`

	bytes, err := cfg.AsJson(false)
	if err != nil {
		t.Logf("No error should be present but got %v\n", err)
		t.Fail()
	}
	if string(bytes) != resShort {
		t.Logf("Result should be:\n%s\nBut got:\n%s\n", resShort, string(bytes))
		t.Fail()
	}

	bytes, err = cfg.AsJson(true)
	if err != nil {
		t.Logf("No error should be present but got %v\n", err)
		t.Fail()
	}
	if string(bytes) != resLong {
		t.Logf("Result should be:\n%s\nBut got:\n%s\n", resLong, string(bytes))
		t.Fail()
	}
}

func TestAsYaml(t *testing.T) {
	cfg := &Config{
		Server: &ServerConfig{
			Port: 80,
		},
		Checks: []CheckConfig{
			{
				Name:    "wow",
				Enabled: false,
			},
		},
	}
	res := `server:
    port: 80
checks:
    - name: wow
      enabled: false
`

	bytes, err := cfg.AsYaml()
	if err != nil {
		t.Logf("No error should be present but got %v\n", err)
		t.Fail()
	}
	if string(bytes) != res {
		t.Logf("Result should be:\n%s\nBut got:\n%s\n", res, string(bytes))
		t.Fail()
	}
}
