package config

var schema = `{
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "hostname": {
      "type": "string"
    },
    "k3os": {
      "$ref": "#/definitions/k3os_config"
    }
  },
  "definitions": {
    "k3os_config": {
      "id": "#/definitions/k3os_config",
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "defaults": {
          "$ref": "#/definitions/defaults_config"
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
    "ssh_config": {
      "id": "#/definitions/ssh_config",
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "address": {
          "type": "string"
        },
        "daemon": {
          "type": "boolean"
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
