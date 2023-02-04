package objects

type ConfigValue struct {
	Value   any `json:"value"`
	Default any `json:"-"`
}

type ConfigKey string

const (
	Static_FlagConfig_Verbose = false
)

const (
	ServerPort      ConfigKey = "serverport"
	DownloadPicture ConfigKey = "downloadpictures"
)

func InitConfig() {
	configMap[ServerPort] = ConfigValue{"8080", "8080"}
	configMap[DownloadPicture] = ConfigValue{true, true}
}

var configMap = make(map[ConfigKey]ConfigValue)

func GetConfigValue[T any](key ConfigKey) T {
	return configMap[key].Value.(T)
}
