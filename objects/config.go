package objects

import "errors"

type ConfigItemBoolean struct {
	Value   bool
	Default bool
}
type ConfigItemInt struct {
	Value   int
	Default int
}

type ConfigItemString struct {
	Value   string
	Default string
}

type WebConfig struct {
	BoolItems   map[ConfigKey]ConfigItemBoolean
	IntItems    map[ConfigKey]ConfigItemInt
	StringItems map[ConfigKey]ConfigItemString
}

type ConfigKey string

const (
	ServerPort ConfigKey = "serverport"
)

func (gc *WebConfig) InitConfig() {
	gc.IntItems = map[ConfigKey]ConfigItemInt{
		ServerPort: {8080, 8080},
	}

	gc.BoolItems = map[ConfigKey]ConfigItemBoolean{}

	gc.StringItems = map[ConfigKey]ConfigItemString{}
}

func (wc *WebConfig) GetConfigItem(key ConfigKey) (int, bool, string) {
	if k, ok := wc.BoolItems[key]; ok {
		return 0, k.Value, ""
	} else if k, ok := wc.IntItems[key]; ok {
		return k.Value, false, ""
	} else if k, ok := wc.StringItems[key]; ok {
		return 0, false, k.Value
	}
	panic(errors.New("Unable to find " + string(key) + " config key !"))
}

func (wc *WebConfig) SetConfigItemValue(key ConfigKey, keyValue interface{}) {
	if k, ok := wc.BoolItems[key]; ok {
		k.Value = keyValue.(bool)
		wc.BoolItems[key] = k
		return
	} else if k, ok := wc.IntItems[key]; ok {
		k.Value = keyValue.(int)
		wc.IntItems[key] = k
		return
	} else if k, ok := wc.StringItems[key]; ok {
		k.Value = keyValue.(string)
		wc.StringItems[key] = k
		return
	}
	panic(errors.New("Unable to find " + string(key) + " config key !"))
}
