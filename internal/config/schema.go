package config

// Config represents the main configuration structure
type Config struct {
	Daemons []Daemon `mapstructure:"daemons" yaml:"daemons" json:"daemons"`
}

// Daemon represents a daemon configuration
type Daemon struct {
	// Basic Information
	Name        string `mapstructure:"name" yaml:"name" json:"name"`
	Label       string `mapstructure:"label" yaml:"label" json:"label"`
	Description string `mapstructure:"description,omitempty" yaml:"description,omitempty" json:"description,omitempty"`

	// Program Information
	Program          string   `mapstructure:"program" yaml:"program" json:"program"`
	ProgramArguments []string `mapstructure:"program_arguments,omitempty" yaml:"program_arguments,omitempty" json:"program_arguments,omitempty"`
	WorkingDirectory string   `mapstructure:"working_directory,omitempty" yaml:"working_directory,omitempty" json:"working_directory,omitempty"`

	// Environment Variables
	EnvironmentVariables map[string]string `mapstructure:"environment_variables,omitempty" yaml:"environment_variables,omitempty" json:"environment_variables,omitempty"`

	// Logging
	StandardOutPath   string `mapstructure:"standard_out_path,omitempty" yaml:"standard_out_path,omitempty" json:"standard_out_path,omitempty"`
	StandardErrorPath string `mapstructure:"standard_error_path,omitempty" yaml:"standard_error_path,omitempty" json:"standard_error_path,omitempty"`

	// Launch Behavior
	RunAtLoad     bool `mapstructure:"run_at_load,omitempty" yaml:"run_at_load,omitempty" json:"run_at_load,omitempty"`
	StartInterval int  `mapstructure:"start_interval,omitempty" yaml:"start_interval,omitempty" json:"start_interval,omitempty"` // seconds

	// Keep Alive Settings
	KeepAlive        *KeepAlive `mapstructure:"keep_alive,omitempty" yaml:"keep_alive,omitempty" json:"keep_alive,omitempty"`
	ThrottleInterval int        `mapstructure:"throttle_interval,omitempty" yaml:"throttle_interval,omitempty" json:"throttle_interval,omitempty"` // seconds

	// Resource Limits
	ResourceLimits *ResourceLimits `mapstructure:"resource_limits,omitempty" yaml:"resource_limits,omitempty" json:"resource_limits,omitempty"`

	// Process Settings
	ProcessType   string `mapstructure:"process_type,omitempty" yaml:"process_type,omitempty" json:"process_type,omitempty"` // Background, Standard, Adaptive, Interactive
	Nice          *int   `mapstructure:"nice,omitempty" yaml:"nice,omitempty" json:"nice,omitempty"`
	InitGroups    bool   `mapstructure:"init_groups,omitempty" yaml:"init_groups,omitempty" json:"init_groups,omitempty"`
	UserName      string `mapstructure:"user_name,omitempty" yaml:"user_name,omitempty" json:"user_name,omitempty"`
	GroupName     string `mapstructure:"group_name,omitempty" yaml:"group_name,omitempty" json:"group_name,omitempty"`
	RootDirectory string `mapstructure:"root_directory,omitempty" yaml:"root_directory,omitempty" json:"root_directory,omitempty"`

	// Socket Activation
	Sockets map[string]Socket `mapstructure:"sockets,omitempty" yaml:"sockets,omitempty" json:"sockets,omitempty"`

	// Calendar/Timing
	StartCalendarInterval []CalendarInterval `mapstructure:"start_calendar_interval,omitempty" yaml:"start_calendar_interval,omitempty" json:"start_calendar_interval,omitempty"`

	// Watch Paths
	WatchPaths []string `mapstructure:"watch_paths,omitempty" yaml:"watch_paths,omitempty" json:"watch_paths,omitempty"`
	QueuePaths []string `mapstructure:"queue_paths,omitempty" yaml:"queue_paths,omitempty" json:"queue_paths,omitempty"`

	// Other Settings
	EnableGlobbing      bool `mapstructure:"enable_globbing,omitempty" yaml:"enable_globbing,omitempty" json:"enable_globbing,omitempty"`
	EnableTransactions  bool `mapstructure:"enable_transactions,omitempty" yaml:"enable_transactions,omitempty" json:"enable_transactions,omitempty"`
	EnablePressuredExit bool `mapstructure:"enable_pressured_exit,omitempty" yaml:"enable_pressured_exit,omitempty" json:"enable_pressured_exit,omitempty"`
	ExitTimeOut         int  `mapstructure:"exit_timeout,omitempty" yaml:"exit_timeout,omitempty" json:"exit_timeout,omitempty"` // seconds
}

// KeepAlive represents keep-alive settings
type KeepAlive struct {
	SuccessfulExit     *bool           `mapstructure:"successful_exit,omitempty" yaml:"successful_exit,omitempty" json:"successful_exit,omitempty"`
	NetworkState       *bool           `mapstructure:"network_state,omitempty" yaml:"network_state,omitempty" json:"network_state,omitempty"`
	PathState          map[string]bool `mapstructure:"path_state,omitempty" yaml:"path_state,omitempty" json:"path_state,omitempty"`
	OtherJobEnabled    map[string]bool `mapstructure:"other_job_enabled,omitempty" yaml:"other_job_enabled,omitempty" json:"other_job_enabled,omitempty"`
	Crashed            *bool           `mapstructure:"crashed,omitempty" yaml:"crashed,omitempty" json:"crashed,omitempty"`
	AfterInitialDemand *bool           `mapstructure:"after_initial_demand,omitempty" yaml:"after_initial_demand,omitempty" json:"after_initial_demand,omitempty"`
}

// ResourceLimits represents resource limitations
type ResourceLimits struct {
	CPU               *int `mapstructure:"cpu,omitempty" yaml:"cpu,omitempty" json:"cpu,omitempty"`
	FileSize          *int `mapstructure:"file_size,omitempty" yaml:"file_size,omitempty" json:"file_size,omitempty"`
	NumberOfFiles     *int `mapstructure:"number_of_files,omitempty" yaml:"number_of_files,omitempty" json:"number_of_files,omitempty"`
	Core              *int `mapstructure:"core,omitempty" yaml:"core,omitempty" json:"core,omitempty"`
	Data              *int `mapstructure:"data,omitempty" yaml:"data,omitempty" json:"data,omitempty"`
	MemoryLock        *int `mapstructure:"memory_lock,omitempty" yaml:"memory_lock,omitempty" json:"memory_lock,omitempty"`
	NumberOfProcesses *int `mapstructure:"number_of_processes,omitempty" yaml:"number_of_processes,omitempty" json:"number_of_processes,omitempty"`
	ResidentSetSize   *int `mapstructure:"resident_set_size,omitempty" yaml:"resident_set_size,omitempty" json:"resident_set_size,omitempty"`
	Stack             *int `mapstructure:"stack,omitempty" yaml:"stack,omitempty" json:"stack,omitempty"`
}

// Socket represents socket activation settings
type Socket struct {
	SockType        string   `mapstructure:"sock_type,omitempty" yaml:"sock_type,omitempty" json:"sock_type,omitempty"` // stream, dgram, seqpacket
	SockPassive     *bool    `mapstructure:"sock_passive,omitempty" yaml:"sock_passive,omitempty" json:"sock_passive,omitempty"`
	SockNodeName    string   `mapstructure:"sock_node_name,omitempty" yaml:"sock_node_name,omitempty" json:"sock_node_name,omitempty"`
	SockServiceName string   `mapstructure:"sock_service_name,omitempty" yaml:"sock_service_name,omitempty" json:"sock_service_name,omitempty"`
	SockFamily      string   `mapstructure:"sock_family,omitempty" yaml:"sock_family,omitempty" json:"sock_family,omitempty"`       // IPv4, IPv6
	SockProtocol    string   `mapstructure:"sock_protocol,omitempty" yaml:"sock_protocol,omitempty" json:"sock_protocol,omitempty"` // TCP, UDP
	SockPathName    string   `mapstructure:"sock_path_name,omitempty" yaml:"sock_path_name,omitempty" json:"sock_path_name,omitempty"`
	SockPathMode    *int     `mapstructure:"sock_path_mode,omitempty" yaml:"sock_path_mode,omitempty" json:"sock_path_mode,omitempty"`
	Bonjour         *bool    `mapstructure:"bonjour,omitempty" yaml:"bonjour,omitempty" json:"bonjour,omitempty"`
	BonjourMultiple []string `mapstructure:"bonjour_multiple,omitempty" yaml:"bonjour_multiple,omitempty" json:"bonjour_multiple,omitempty"`
}

// CalendarInterval represents calendar-based scheduling
type CalendarInterval struct {
	Minute  *int `mapstructure:"minute,omitempty" yaml:"minute,omitempty" json:"minute,omitempty"`
	Hour    *int `mapstructure:"hour,omitempty" yaml:"hour,omitempty" json:"hour,omitempty"`
	Day     *int `mapstructure:"day,omitempty" yaml:"day,omitempty" json:"day,omitempty"`
	Weekday *int `mapstructure:"weekday,omitempty" yaml:"weekday,omitempty" json:"weekday,omitempty"` // 0-7 (0 and 7 are Sunday)
	Month   *int `mapstructure:"month,omitempty" yaml:"month,omitempty" json:"month,omitempty"`
}
