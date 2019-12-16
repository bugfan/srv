package main

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
)

type hs struct {
}

func (*hs) Process(r *http.Request) {

}

func (*hs) Config() ActionConfig {
	return nil
}

func testHttpServer() *hs {
	return &hs{}
}

func main() {
	mConfig := &URLMatcherConfig{
		Method: "ALL",
		Path:   "/api",
	}
	aConfig := &HandlerActionConfig{Handler: testHttpServer()} // 测试http服务
	rule, err := MakeRule(mConfig, aConfig)
	if err != nil {
		return //nil, err
	}
	wsmConfig := &URLMatcherConfig{
		Method: "ALL",
		Path:   "/ws",
	}
	wsaConfig := &HandlerActionConfig{Handler: testHttpServer()} //测试websocket服务
	wsRule, err := MakeRule(wsmConfig, wsaConfig)
	if err != nil {
		return //nil, err
	}
	ruleSet, err := NewRuleSet(rule, wsRule)
	if err != nil {
		return //nil, err
	}
	fmt.Println("ERR:", ruleSet)
}

var matcherConfigs = make(map[string]MatcherConfig)

func RegisterMatcherConfig(name string, config MatcherConfig) error {
	if _, ok := matcherConfigs[name]; ok {
		return fmt.Errorf("factory %s already exists", name)
	}
	matcherConfigs[name] = config
	return nil
}

type MatcherConfig interface {
	Name() string
	Factory(MatcherConfig) (Matcher, error)
}

type Matcher interface {
	Match(*http.Request) bool
	Config() MatcherConfig
}

func GetMatcher(config MatcherConfig) (Matcher, error) {
	if factory, ok := matcherConfigs[config.Name()]; ok {
		return factory.Factory(config)
	}
	return nil, fmt.Errorf("factory %s not exists", config.Name())
}

// type Handler interface {
// 	ServeHTTP(http.ResponseWriter, *http.Request)
// }

var actionConfigs = make(map[string]ActionConfig)

type Action interface {
	Process(*http.Request)
	Config() ActionConfig
}

type ActionConfig interface {
	Name() string
	Factory(ActionConfig) (Action, error)
}

func RegisterActionConfig(name string, config ActionConfig) error {
	if _, ok := actionConfigs[name]; ok {
		return fmt.Errorf("factory %s already exists", name)
	}
	actionConfigs[name] = config
	return nil
}

func GetAction(config ActionConfig) (Action, error) {
	if factory, ok := actionConfigs[config.Name()]; ok {
		return factory.Factory(config)
	}
	return nil, fmt.Errorf("factory %s not exists", config.Name())
}

/// impl
func init() {
	RegisterMatcherConfig("url_matcher", new(URLMatcherConfig))
}

type URLMatcherConfig struct {
	Method string
	Path   string
	Query  map[string]string
}

func (URLMatcherConfig) Name() string {
	return "url_matcher"
}

type URLMatcher struct {
	Method     string
	pathRegex  *regexp.Regexp
	queryRegex map[string]*regexp.Regexp
	config     *URLMatcherConfig
}

func (*URLMatcherConfig) Factory(mc MatcherConfig) (Matcher, error) {
	if config, ok := mc.(*URLMatcherConfig); ok && mc != nil {
		m := new(URLMatcher)
		m.Method = config.Method
		m.config = config
		if config.Path != "" {
			pathRegex, err := RegexpCompileStart(config.Path)
			if err != nil {
				return nil, err
			}
			m.pathRegex = pathRegex
		}

		if config.Query != nil && len(config.Query) > 0 {
			m.queryRegex = make(map[string]*regexp.Regexp)
			for key, val := range config.Query {
				re, err := RegexpCompileStart(val)
				if err != nil {
					return nil, err
				}
				m.queryRegex[key] = re
			}
		}
		return m, nil
	}
	return nil, errors.New("matcher config error")
}

func NewURLMatcher(method, path string) Matcher {
	config := &URLMatcherConfig{
		Method: method,
		Path:   path,
		Query:  make(map[string]string),
	}

	pathRegex, err := RegexpCompileStart(config.Path)
	fmt.Printf("NewURLMatcher err:%v\n", err)
	// if err != nil {
	// 	return nil
	// }
	queryRegex := make(map[string]*regexp.Regexp)
	for key, val := range config.Query {
		re, err := RegexpCompileStart(val)
		fmt.Printf("RegexpCompileStart err:", err, key, val)
		// if err != nil {
		// 	return nil, err
		// }
		queryRegex[key] = re
	}

	urlMatcher := &URLMatcher{
		Method:     method,
		pathRegex:  pathRegex,
		queryRegex: queryRegex,
		config:     config,
	}
	return urlMatcher
}

func (m *URLMatcher) Match(r *http.Request) bool {
	if m.Method != "" && m.Method != "ALL" && m.Method != r.Method {
		return false
	}
	if m.pathRegex != nil {
		if !m.pathRegex.MatchString(r.URL.Path) {
			return false
		}
	}
	if m.queryRegex != nil && len(m.queryRegex) > 0 {
		for key, val := range m.queryRegex {
			if !val.MatchString(r.URL.Query().Get(key)) {
				return false
			}
		}
	}
	return true
}

func (m *URLMatcher) Config() MatcherConfig {
	return m.config
}

func RegexpCompileStart(s string) (*regexp.Regexp, error) {
	s = "^" + s
	return regexp.Compile(s)
}

// impl action

func init() {
	RegisterActionConfig("handler", new(HandlerActionConfig))
}

type HandlerActionConfig struct {
	Handler
}

func (HandlerActionConfig) Name() string {
	return "handler"
}

func (*HandlerActionConfig) Factory(ac ActionConfig) (Action, error) {
	if config, ok := ac.(*HandlerActionConfig); ok && ac != nil {
		return &Handler{
			config,
		}, nil
	}
	return nil, errors.New("action config error")
}

type Handler struct {
	*HandlerActionConfig
}

func (a *Handler) Process(r *http.Request) {
	// i.Handle = a.Handler
}

func (a *Handler) Config() ActionConfig {
	return a.HandlerActionConfig
}

type Rule struct {
	NotMatcher bool
	Matcher
	Action
}

func MakeRule(mc MatcherConfig, ac ActionConfig) (*Rule, error) {
	m, err := GetMatcher(mc)
	if err != nil {
		return nil, err
	}

	a, err := GetAction(ac)
	if err != nil {
		return nil, err
	}
	return &Rule{
		Matcher: m,
		Action:  a,
	}, nil
}

func NewRuleSet(rules ...*Rule) (*RuleSet, error) {
	return &RuleSet{
		Rules: rules,
	}, nil
}

type RuleSet struct {
	Rules []*Rule
}
