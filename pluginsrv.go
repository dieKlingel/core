package main

import (
	"plugin"
)

type PluginService struct {
	storageService *StorageService
	plugins        []*plugin.Plugin
}

type PluginProviderFunc func(...interface{}) any

func NewPluginService(storageService *StorageService) *PluginService {
	return &PluginService{
		storageService: storageService,
		plugins:        make([]*plugin.Plugin, 0),
	}
}

func (pluginService *PluginService) LoadPluginProviderFunctions() []any {
	providerFuncs := make([]any, 0)

	plugin, err := plugin.Open("plugins/ruler/ruler.so")
	if err != nil {
		panic(err.Error())
	}

	newServiceFunc, err := plugin.Lookup("NewPlugin")
	if err != nil {
		panic(err.Error())
	}
	pluginService.plugins = append(pluginService.plugins, plugin)
	providerFuncs = append(providerFuncs, newServiceFunc)

	return providerFuncs
}

func (pluginService *PluginService) LoadPluginInitalizer() []any {
	initializerFuncs := make([]any, 0)

	for _, plugin := range pluginService.plugins {
		initFunc, err := plugin.Lookup("Run")
		if err != nil {
			continue
		}
		initializerFuncs = append(initializerFuncs, initFunc)
	}

	return initializerFuncs
}

func (pluginService *PluginService) Run() {

	println("plugin service")
}
