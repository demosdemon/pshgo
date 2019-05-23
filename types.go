package pshgo

type (
	JSONArray  = []interface{}
	JSONObject = map[string]interface{}
	StringMap  = map[string]string
)

type (
	Access           map[AccessType]AccessLevel
	Caches           map[string]CacheConfiguration
	Crons            map[string]Cron
	Mounts           map[string]Mount
	RedirectPaths    map[string]RedirectPath
	SourceOperations map[string]SourceOperation
	Variables        JSONObject
	WebLocations     map[string]WebLocation
	WebRules         map[string]WebRule
	Workers          map[string]Worker
)

type (
	Application struct {
		ApplicationCore
		Web     Web     `json:"web"`
		Hooks   Hooks   `json:"hooks"`
		Crons   Crons   `json:"crons"`
		Workers Workers `json:"workers"`
		TreeID  string  `json:"tree_id"`
		SlugID  string  `json:"slug_id"`
		AppDir  string  `json:"app_dir"`
	}

	ApplicationBase struct {
		Size          ServiceSize `json:"size"`
		Disk          uint32      `json:"disk"`
		Access        Access      `json:"access"`
		Relationships StringMap   `json:"relationships"`
		Mounts        Mounts      `json:"mounts"`
		Timezone      string      `json:"timezone"` // TODO: replace with serializable time.Location
		Variables     Variables   `json:"variables"`
	}

	ApplicationBuilder struct {
		ApplicationCore
		Dependencies JSONObject `json:"dependencies"`
		Build        Build      `json:"build"`
		Source       Source     `json:"source"`
	}

	ApplicationCore struct {
		ApplicationBase
		Name      string      `json:"name"`
		Type      string      `json:"type"`
		Runtime   interface{} `json:"runtime"`
		Preflight Preflight   `json:"preflight"`
	}

	Build struct {
		Flavor string `json:"flavor"`
		Caches Caches `json:"caches"`
	}

	Cache struct {
		Enabled    bool     `json:"enabled"`
		DefaultTTL int      `json:"default_ttl"`
		Cookies    []string `json:"cookies"`
		Headers    []string `json:"headers"`
	}

	CacheConfiguration struct {
		Directory        string   `json:"directory"`
		Watch            []string `json:"watch"`
		AllowStale       bool     `json:"allow_stale"`
		ShareBetweenApps bool     `json:"share_between_apps"`
	}

	Commands struct {
		Start string `json:"start"`
		Stop  string `json:"stop,omitempty"`
	}

	Cron struct {
		Spec string `json:"spec"`
		Cmd  string `json:"cmd"`
	}

	HTTPAccess struct {
		Addresses []string          `json:"addresses"`
		BasicAuth map[string]string `json:"basic_auth"`
	}

	Hooks struct {
		Build      string `json:"build"`
		Deploy     string `json:"deploy"`
		PostDeploy string `json:"post_deploy"`
	}

	Mount struct {
		Source     ApplicationMount `json:"source"`
		SourcePath string           `json:"path"`
		Service    string           `json:"service,omitempty"`
	}

	Preflight struct {
		Enabled      bool     `json:"enabled"`
		IgnoredRules []string `json:"ignored_rules"`
	}

	Redirects struct {
		Expires Duration      `json:"expires"`
		Paths   RedirectPaths `json:"paths"`
	}

	RedirectPath struct {
		Regexp       bool     `json:"regexp"`
		To           string   `json:"to"`
		Prefix       bool     `json:"prefix"`
		AppendSuffix bool     `json:"append_suffix"`
		Code         int      `json:"code"`
		Expires      Duration `json:"expires"`
	}

	Route struct {
		Primary        bool              `json:"primary"`
		ID             *string           `json:"id"`
		OriginalURL    string            `json:"original_url"`
		Attributes     map[string]string `json:"attributes"`
		Type           string            `json:"type"`
		Redirects      Redirects         `json:"redirects"`
		TLS            TLSSettings       `json:"tls"`
		HTTPAccess     HTTPAccess        `json:"http_access"`
		RestrictRobots bool              `json:"restrict_robots"`

		// Upstream Routes
		Cache    Cache  `json:"cache"`
		SSI      SSI    `json:"ssi"`
		Upstream string `json:"upstream"`

		// Redirect Routes
		To string `json:"to"`
	}

	RouteIdentification struct {
		Scheme string `json:"scheme"`
		Host   string `json:"host"`
		Path   string `json:"path"`
	}

	RouteRepresentation struct {
		Project     string              `json:"project"`
		Environment string              `json:"environment"`
		Route       RouteIdentification `json:"route"`
	}

	SSI struct {
		Enabled bool `json:"enabled"`
	}

	Source struct {
		Operations SourceOperations `json:"operations"`
	}

	SourceOperation struct {
		Command string `json:"command"`
	}

	TLSSettings struct {
		StrictTransportSecurity      TLSSTS        `json:"strict_transport_security"`
		MinVersion                   *TLSVersion   `json:"min_version"`
		ClientAuthentication         string        `json:"client_authentication"`
		ClientCertificateAuthorities []Certificate `json:"client_certificate_authorities"`
	}

	TLSSTS struct {
		Enabled           bool `json:"enabled"`
		IncludeSubdomains bool `json:"include_subdomains"`
		Preload           bool `json:"preload"`
	}

	Upstream struct {
		SocketFamily SocketFamily   `json:"socket_family"`
		Protocol     SocketProtocol `json:"socket_protocol"`
	}

	Web struct {
		// ApplicationBase
		Locations    WebLocations `json:"locations"`
		Commands     Commands     `json:"commands"`
		Upstream     Upstream     `json:"upstream"`
		DocumentRoot *string      `json:"document_root,omitempty"` // deprecated
		Passthru     *string      `json:"passthru,omitempty"`      // deprecated
		IndexFiles   []string     `json:"index_files,omitempty"`   // deprecated
		Whitelist    []string     `json:"whitelist,omitempty"`     // deprecated
		Blacklist    []string     `json:"blacklist,omitempty"`     // deprecated
		Expires      *Duration    `json:"expires,omitempty"`       // deprecated
		MoveToRoot   *bool        `json:"move_to_root,omitempty"`  // deprecated
	}

	WebLocation struct {
		Root     string    `json:"root"`
		Expires  Duration  `json:"expires"`
		Passthru Passthru  `json:"passthru"`
		Scripts  bool      `json:"scripts"`
		Index    []string  `json:"index"`
		Allow    bool      `json:"allow"`
		Headers  StringMap `json:"headers"`
		Rules    WebRules  `json:"rules"`
	}

	WebRule struct {
		Expires  Duration  `json:"expires"`
		Passthru Passthru  `json:"passthru"`
		Scripts  bool      `json:"scripts"`
		Allow    bool      `json:"allow"`
		Headers  StringMap `json:"headers"`
	}

	Worker struct {
		// ApplicationBase
		Commands Commands `json:"commands"`
	}
)
