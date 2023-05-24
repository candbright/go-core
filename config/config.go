package config

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"os"
	"strconv"
	"strings"
)

const YAML ParseType = "YAML"
const XML ParseType = "XML"
const JSON ParseType = "JSON"

type ParseType string

type Config struct {
	data map[interface{}]interface{}
}

func Parse(data []byte, parseType ParseType) (*Config, error) {
	config := &Config{}
	err := config.parse(string(data), parseType)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) Parse(fileName string, parseType ParseType) error {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	err = c.parse(string(file), parseType)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) parse(data string, parseType ParseType) error {
	expandData := os.ExpandEnv(data)
	var err error
	switch parseType {
	case YAML:
		err = yaml.Unmarshal([]byte(expandData), &c.data)
		if err != nil {
			return err
		}
		break
	case XML:
		err = xml.Unmarshal([]byte(expandData), &c.data)
		if err != nil {
			return err
		}
		break
	case JSON:
		err = json.Unmarshal([]byte(expandData), &c.data)
		if err != nil {
			return err
		}
		break
	default:
		return errors.New("unsupported parse type")
	}
	return nil
}

func (c *Config) Get(key string) string {
	if c.data == nil {
		fmt.Println("config data is nil, please parse first")
		return ""
	}
	split := strings.Split(key, ".")
	value := get(split, c.data)
	if value == nil {
		fmt.Println("value is nil, please input the right key")
		return ""
	}
	if val, ok := value.(string); ok {
		return val
	}
	if val, ok := value.(int); ok {
		return strconv.Itoa(val)
	}
	if val, ok := value.(int64); ok {
		return strconv.FormatInt(val, 10)
	}
	if valMap, ok := value.(map[interface{}]interface{}); ok {
		env := envParser(valMap)
		if env != nil {
			return env.Get()
		}
		archV := archParser(valMap)
		if archV != nil {
			return archV.Get()
		}
		osV := osParser(valMap)
		if osV != nil {
			return osV.Get()
		}
	}
	fmt.Println("-value type is unsupported")
	return ""
}

func get(keys []string, tree map[interface{}]interface{}) interface{} {
	if keys == nil {
		return nil
	}
	if len(keys) == 1 {
		value, ok := tree[keys[0]]
		if ok {
			return value
		} else {
			return nil
		}
	}
	value, ok := tree[keys[0]]
	if valueMap, valueOk := value.(map[interface{}]interface{}); valueOk && ok {
		return get(keys[1:], valueMap)
	} else {
		return nil
	}
}
