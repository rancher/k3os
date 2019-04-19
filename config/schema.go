package config

var schema = `{
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "hostname": {
      "type": "string"
    },
    "k3s": {
      "$ref": "#/definitions/k3s_config"
    },
    "k3os": {
      "$ref": "#/definitions/k3os_config"
    },
    "runcmd": {
      "type": "array"
    },
    "write_files": {
      "items": {
        "$ref": "#/definitions/file_config"
      }
    }
  },
  "definitions": {
    "k3s_config": {
      "id": "#/definitions/k3s_config",
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "role": {
          "type": "string"
        },
        "extra_args": {
          "$ref": "#/definitions/list_of_strings"
        }
      }
    },
    "k3os_config": {
      "id": "#/definitions/k3os_config",
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "defaults": {
          "$ref": "#/definitions/defaults_config"
        },
        "environment": {
          "type": "object"
        },
        "modules": {
          "$ref": "#/definitions/list_of_strings"
        },
        "ssh": {
          "$ref": "#/definitions/ssh_config"
        },
        "sysctl": {
          "type": "object"
        },
        "upgrade": {
          "$ref": "#/definitions/upgrade_config"
        },
        "network": {
          "$ref": "#/definitions/network_config"
        },
        "password": {
          "type": "string"
        }
      }
    },
    "defaults_config": {
      "id": "#/definitions/defaults_config",
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "modules": {
          "$ref": "#/definitions/list_of_strings"
        }
      }
    },
    "file_config": {
      "id": "#/definitions/file_config",
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "content": {
          "type": "string"
        },
        "encoding": {
          "type": "string"
        },
        "owner": {
          "type": "string"
        },
        "path": {
          "type": "string"
        },
        "permissions": {
          "type": "string"
        }
      }
    },
    "ssh_config": {
      "id": "#/definitions/ssh_config",
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "address": {
          "type": "string"
        },
        "authorized_keys": {
          "$ref": "#/definitions/list_of_strings"
        },
        "daemon": {
          "type": "boolean"
        },
        "host_keys": {
          "type": "object"
        },
        "port": {
          "type": "integer"
        }
      }
    },
    "upgrade_config": {
      "id": "#/definitions/upgrade_config",
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "url": {
          "type": "string"
        },
        "rollback": {
          "type": "string"
        },
        "policy": {
          "type": "string"
        }
      }
    },
    "network_config": {
      "id": "#/definitions/network_config",
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "dns": {
          "$ref": "#/definitions/dns_config"
        },
        "interfaces": {
          "$ref": "#/definitions/interface_additional"
        },
        "proxy": {
          "$ref": "#/definitions/proxy_config"
        }
      }
    },
    "dns_config": {
      "id": "#/definitions/dns_config",
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "searches": {
          "$ref": "#/definitions/list_of_strings"
        },
        "nameservers": {
          "$ref": "#/definitions/list_of_strings"
        }
      }
    },
    "interface_additional": {
      "id": "#/definitions/interface_additional",
      "type": "object",
      "additionalProperties": {
        "$ref": "#/definitions/interface_config"
      }
    },
    "interface_config": {
      "id": "#/definitions/interface_config",
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "addresses": {
          "$ref": "#/definitions/list_of_strings"
        },
        "gateway": {
          "type": "string"
        },
        "ipv4ll": {
          "type": "boolean"
        },
        "metric": {
          "type": "integer"
        }
      }
    },
    "proxy_config": {
      "id": "#/definitions/proxy_config",
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "http_proxy": {
          "type": "string"
        },
        "https_proxy": {
          "type": "string"
        },
        "no_proxy": {
          "type": "string"
        }
      }
    },
    "list_of_strings": {
      "type": "array",
      "items": {
        "type": "string"
      },
      "uniqueItems": true
    }
  }
}
`
