package engine

import (
	"fmt"
	"path/filepath"
)

// Engine is a code execution engine, supporting a set of languages for
// sandboxed code execution.
type Engine struct {
	languages []*Language
}

// Languages returns the list of supported execution languages.
func (eng *Engine) Languages() []*Language {
	return eng.languages
}

// Run executes the given source code with the given options, using the
// specified language.
func (eng *Engine) Run(languageName string, opts *RunOptions) (*RunResult, error) {
	lang := eng.findLanguage(languageName)
	if lang == nil {
		return nil, fmt.Errorf("Language not found: '%s'", languageName)
	}

	return lang.Run(opts)
}

// New creates a new execution engine using the yaml config files in the
// specified directory. Only disableApparmor for development, docker alone isn't
// sufficient to fully sandbox the code.
func New(confPath string, disableApparmor bool) (result *Engine, err error) {
	result = &Engine{}
	configs, err := filepath.Glob(confPath + "/lang-*.yml")
	if err != nil {
		return
	}

	if len(configs) < 1 {
		err = fmt.Errorf("no languages found at path '%s'", confPath)
		return
	}

	for _, config := range configs {
		var lang *Language
		lang, err = loadLanguage(config)
		if err != nil {
			err = fmt.Errorf("Error loading language '%s': %s", config, err)
			// fail everything if one language fails to load
			return
		}
		if disableApparmor {
			lang.disableAppArmor()
		}
		result.languages = append(result.languages, lang)
	}

	return
}

func (eng *Engine) findLanguage(name string) *Language {
	for _, lang := range eng.Languages() {
		if lang.Name == name {
			return lang
		}
	}
	return nil
}
