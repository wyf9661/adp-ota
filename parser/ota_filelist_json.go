package parser

type OutputItem struct {
	Type         string `json:"type"`
	Architecture string `json:"architecture"`
	Name         string `json:"name"`
	Version      string `json:"version"`
	Md5          string `json:"md5"`
}

type ConfigStruct struct {
	Outputs []OutputItem `json:"outputs"`
}

var GlobalPackageInfoStruct ConfigStruct
