package visual

import (
	"encoding/json"
	"fmt"
	"strings"

	"moonbridge/internal/extension/plugin"
	"moonbridge/internal/config"
	"moonbridge/internal/format"
)

const PluginName = "visual"

type EnabledFunc func(modelAlias string) bool

type Config struct {
	Provider  string `json:"provider,omitempty" yaml:"provider"`
	Model     string `json:"model,omitempty" yaml:"model"`
	MaxRounds int    `json:"max_rounds,omitempty" yaml:"max_rounds"`
	MaxTokens int    `json:"max_tokens,omitempty" yaml:"max_tokens"`
}

// Plugin injects the Visual tools for models that opt in.
type Plugin struct {
	plugin.BasePlugin
	isEnabled EnabledFunc
	pluginCfg config.PluginConfig
}

func NewPlugin(isEnabled ...EnabledFunc) *Plugin {
	var enabled EnabledFunc
	if len(isEnabled) > 0 {
		enabled = isEnabled[0]
	}
	return &Plugin{isEnabled: enabled}
}

func (p *Plugin) Name() string { return PluginName }

func (p *Plugin) ConfigSpecs() []config.ExtensionConfigSpec { return ConfigSpecs() }

func (p *Plugin) Init(ctx plugin.PluginContext) error {
	p.pluginCfg = config.PluginFromGlobalConfig(&ctx.AppConfig)
	return nil
}

func (p *Plugin) EnabledForModel(model string) bool {
	if p.isEnabled != nil {
		return p.isEnabled(model)
	}
	if setting, ok := p.pluginCfg.Extensions[PluginName]; ok && setting.Enabled != nil {
		return *setting.Enabled
	}
	return false
}

func (p *Plugin) InjectTools(_ *plugin.RequestContext) []format.CoreTool {
	return CoreTools()
}

// CoreInjectTools returns Core-format tool definitions for visual analysis.
func (p *Plugin) CoreInjectTools(_ *plugin.RequestContext) []format.CoreTool {
	return CoreTools()
}

func ConfigSpecs() []config.ExtensionConfigSpec {
	return []config.ExtensionConfigSpec{{
		Name: PluginName,
		Scopes: []config.ExtensionScope{
			config.ExtensionScopeGlobal,
			config.ExtensionScopeProvider,
			config.ExtensionScopeModel,
			config.ExtensionScopeRoute,
		},
		Factory: func() any { return &Config{} },
		Validate: func(cfg config.Config) error {
			return ValidateConfig(
				config.PluginFromGlobalConfig(&cfg),
				config.ProviderFromGlobalConfig(&cfg),
			)
		},
	}}
}

func ConfigForModel(pluginCfg config.PluginConfig, modelAlias string) (Config, bool) {
	// Check if visual extension is enabled globally.
	if !pluginExtensionEnabled(pluginCfg, PluginName) {
		return Config{}, false
	}
	// Decode typed config from global RawConfig.
	var cfg *Config
	if setting, ok := pluginCfg.Extensions[PluginName]; ok && len(setting.RawConfig) > 0 {
		data, err := json.Marshal(setting.RawConfig)
		if err == nil {
			_ = json.Unmarshal(data, &cfg)
		}
	}
	if cfg == nil {
		return Config{}, true
	}
	return cfg.Normalized(), true
}

func pluginExtensionEnabled(pluginCfg config.PluginConfig, name string) bool {
	if setting, ok := pluginCfg.Extensions[name]; ok && setting.Enabled != nil {
		return *setting.Enabled
	}
	return false
}

func (cfg Config) Normalized() Config {
	cfg.Provider = strings.TrimSpace(cfg.Provider)
	cfg.Model = strings.TrimSpace(cfg.Model)
	if cfg.MaxRounds <= 0 {
		cfg.MaxRounds = 4
	}
	if cfg.MaxTokens <= 0 {
		cfg.MaxTokens = 2048
	}
	return cfg
}

func ValidateConfig(pluginCfg config.PluginConfig, providerCfg config.ProviderConfig) error {
	for alias := range providerCfg.Routes {
		if err := validateModelConfig(pluginCfg, providerCfg, alias); err != nil {
			return err
		}
	}
	for providerKey, def := range providerCfg.Providers {
		for modelName := range def.Models {
			if err := validateModelConfig(pluginCfg, providerCfg, providerKey+"/"+modelName); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateModelConfig(pluginCfg config.PluginConfig, providerCfg config.ProviderConfig, modelAlias string) error {
	cfg, ok := ConfigForModel(pluginCfg, modelAlias)
	if !ok {
		return nil
	}
	if cfg.Provider == "" {
		return fmt.Errorf("extensions.%s.config.provider is required when visual is enabled for %s", PluginName, modelAlias)
	}
	if cfg.Model == "" {
		return fmt.Errorf("extensions.%s.config.model is required when visual is enabled for %s", PluginName, modelAlias)
	}
	def, ok := providerCfg.Providers[cfg.Provider]
	if !ok {
		return fmt.Errorf("extensions.%s.config.provider references unknown provider %q", PluginName, cfg.Provider)
	}
	// Protocol constraint removed — visual extension operates on Core format.
	_ = def
	return nil
}

var (
	_ plugin.Plugin             = (*Plugin)(nil)
	_ plugin.ConfigSpecProvider = (*Plugin)(nil)
)

var (
	_ plugin.Plugin             = (*Plugin)(nil)
	_ plugin.ConfigSpecProvider = (*Plugin)(nil)
	_ plugin.ToolInjector       = (*Plugin)(nil)
)
