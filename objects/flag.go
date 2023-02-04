package objects

import (
	"os"
	"sort"
)

type CommandFlag struct {
	FlagExecutor func(args []string) []string
	Description  string
	Usage        string
	IsAliase     bool
	AliaseOf     string
}

type FlagHelpContent struct {
	Principal   string
	Description string
	Usage       string
	Aliases     []string
}

var flagExecutors = map[string](CommandFlag){}

func BuildFlagHelpMenu() []string {
	flags := map[string](FlagHelpContent){}
	for k, v := range flagExecutors {
		if !v.IsAliase {
			flagHelpContent := FlagHelpContent{k, v.Description, v.Usage, []string{k}}
			flagHelpContent.Description = v.Description
			flags[flagHelpContent.Principal] = flagHelpContent
		}
	}
	for k, v := range flagExecutors {
		if v.IsAliase {
			if val, ok := flags[v.AliaseOf]; ok {
				val.Aliases = append(val.Aliases, k)
				flags[val.Principal] = val
			}
		}
	}
	flagList := []FlagHelpContent{}
	for _, v := range flags {
		flagList = append(flagList, v)
	}
	sort.SliceStable(flagList, func(i, j int) bool {
		return flagList[i].Principal < flagList[j].Principal
	})
	helpBlock := []string{}
	for _, v := range flagList {
		helpBlock = append(helpBlock, "\t"+v.Principal+"\n")
		helpBlock = append(helpBlock, "\t\tDescription: "+v.Description+"\n")
		helpBlock = append(helpBlock, "\t\tUsage: "+v.Usage+"\n")
		helpBlock = append(helpBlock, "\t\tAliases: ")
		argAliasesStr := ""
		sort.SliceStable(v.Aliases, func(i, j int) bool {
			return v.Aliases[i] < v.Aliases[j]
		})
		for i, a := range v.Aliases {
			argAliasesStr += a
			if len(v.Aliases)-1 > i {
				argAliasesStr += ", "
			}
		}
		helpBlock = append(helpBlock, argAliasesStr+"\n\n")
	}
	return helpBlock
}

func ProcessArguments(args []string) {
	for len(args) > 0 {
		arg := args[0]
		if arg[0] == "-"[0] {
			if arg == "--help" || arg == "-h" {
				flagShowHelpMessage(args)
				return
			}
			if val, ok := flagExecutors[arg]; ok {
				cmdFlag := val
				if cmdFlag.IsAliase {
					cmdFlag = flagExecutors[cmdFlag.AliaseOf]
				}
				args = cmdFlag.FlagExecutor(args)
			} else {
				args = append(args[:0], args[1:]...)
				println("Can't find argument " + arg)
			}
		}
	}
}

func flagShowHelpMessage(args []string) []string {
	for _, l := range BuildFlagHelpMenu() {
		print(l)
	}
	os.Exit(0)
	return args
}
