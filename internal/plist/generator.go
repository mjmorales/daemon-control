package plist

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/mjmorales/daemon-control/internal/config"
)

// Generator creates plist files from daemon configurations
type Generator struct {
	outputDir string
}

// NewGenerator creates a new plist generator
func NewGenerator(outputDir string) *Generator {
	return &Generator{
		outputDir: outputDir,
	}
}

// GenerateAll generates plist files for all daemons
func (g *Generator) GenerateAll(daemons []config.Daemon) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(g.outputDir, 0750); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	for _, daemon := range daemons {
		if err := g.Generate(&daemon); err != nil {
			return fmt.Errorf("failed to generate plist for %s: %w", daemon.Name, err)
		}
	}

	return nil
}

// Generate creates a plist file for a single daemon
func (g *Generator) Generate(daemon *config.Daemon) error {
	plist := g.daemonToPlist(daemon)

	// Marshal to XML
	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	buf.WriteString(`<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">` + "\n")

	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", "    ")

	if err := encoder.Encode(plist); err != nil {
		return fmt.Errorf("failed to encode plist: %w", err)
	}

	// Write to file
	outputPath := filepath.Join(g.outputDir, daemon.Name+".plist")
	if err := os.WriteFile(outputPath, buf.Bytes(), 0600); err != nil {
		return fmt.Errorf("failed to write plist file: %w", err)
	}

	log.Info().
		Str("daemon", daemon.Name).
		Str("output", outputPath).
		Msg("Generated plist file")

	return nil
}

// daemonToPlist converts a daemon config to plist structure
func (g *Generator) daemonToPlist(daemon *config.Daemon) *Plist {
	dict := &Dict{}

	// Label (required)
	dict.AddString("Label", daemon.Label)

	// Program or ProgramArguments
	if len(daemon.ProgramArguments) > 0 {
		dict.AddStringArray("ProgramArguments", daemon.ProgramArguments)
	} else if daemon.Program != "" {
		dict.AddStringArray("ProgramArguments", []string{daemon.Program})
	}

	// Working Directory
	if daemon.WorkingDirectory != "" {
		dict.AddString("WorkingDirectory", daemon.WorkingDirectory)
	}

	// Environment Variables
	if len(daemon.EnvironmentVariables) > 0 {
		envDict := &Dict{}
		for k, v := range daemon.EnvironmentVariables {
			envDict.AddString(k, v)
		}
		dict.AddDict("EnvironmentVariables", envDict)
	}

	// Logging
	if daemon.StandardOutPath != "" {
		dict.AddString("StandardOutPath", daemon.StandardOutPath)
	}
	if daemon.StandardErrorPath != "" {
		dict.AddString("StandardErrorPath", daemon.StandardErrorPath)
	}

	// Launch behavior
	if daemon.RunAtLoad {
		dict.AddBool("RunAtLoad", daemon.RunAtLoad)
	}

	if daemon.StartInterval > 0 {
		dict.AddInteger("StartInterval", daemon.StartInterval)
	}

	// Keep Alive
	if daemon.KeepAlive != nil {
		keepAliveDict := &Dict{}

		if daemon.KeepAlive.SuccessfulExit != nil {
			keepAliveDict.AddBool("SuccessfulExit", *daemon.KeepAlive.SuccessfulExit)
		}
		if daemon.KeepAlive.NetworkState != nil {
			keepAliveDict.AddBool("NetworkState", *daemon.KeepAlive.NetworkState)
		}
		if daemon.KeepAlive.Crashed != nil {
			keepAliveDict.AddBool("Crashed", *daemon.KeepAlive.Crashed)
		}
		if daemon.KeepAlive.AfterInitialDemand != nil {
			keepAliveDict.AddBool("AfterInitialDemand", *daemon.KeepAlive.AfterInitialDemand)
		}

		// PathState
		if len(daemon.KeepAlive.PathState) > 0 {
			pathDict := &Dict{}
			for path, state := range daemon.KeepAlive.PathState {
				pathDict.AddBool(path, state)
			}
			keepAliveDict.AddDict("PathState", pathDict)
		}

		// OtherJobEnabled
		if len(daemon.KeepAlive.OtherJobEnabled) > 0 {
			jobDict := &Dict{}
			for job, enabled := range daemon.KeepAlive.OtherJobEnabled {
				jobDict.AddBool(job, enabled)
			}
			keepAliveDict.AddDict("OtherJobEnabled", jobDict)
		}

		if len(keepAliveDict.Items) > 0 {
			dict.AddDict("KeepAlive", keepAliveDict)
		}
	}

	// Throttle Interval
	if daemon.ThrottleInterval > 0 {
		dict.AddInteger("ThrottleInterval", daemon.ThrottleInterval)
	}

	// Process Settings
	if daemon.ProcessType != "" {
		dict.AddString("ProcessType", daemon.ProcessType)
	}
	if daemon.Nice != nil {
		dict.AddInteger("Nice", *daemon.Nice)
	}
	if daemon.InitGroups {
		dict.AddBool("InitGroups", daemon.InitGroups)
	}
	if daemon.UserName != "" {
		dict.AddString("UserName", daemon.UserName)
	}
	if daemon.GroupName != "" {
		dict.AddString("GroupName", daemon.GroupName)
	}
	if daemon.RootDirectory != "" {
		dict.AddString("RootDirectory", daemon.RootDirectory)
	}

	// Resource Limits
	if daemon.ResourceLimits != nil {
		limitsDict := &Dict{}

		if daemon.ResourceLimits.CPU != nil {
			limitsDict.AddInteger("CPU", *daemon.ResourceLimits.CPU)
		}
		if daemon.ResourceLimits.FileSize != nil {
			limitsDict.AddInteger("FileSize", *daemon.ResourceLimits.FileSize)
		}
		if daemon.ResourceLimits.NumberOfFiles != nil {
			limitsDict.AddInteger("NumberOfFiles", *daemon.ResourceLimits.NumberOfFiles)
		}
		if daemon.ResourceLimits.Core != nil {
			limitsDict.AddInteger("Core", *daemon.ResourceLimits.Core)
		}
		if daemon.ResourceLimits.Data != nil {
			limitsDict.AddInteger("Data", *daemon.ResourceLimits.Data)
		}
		if daemon.ResourceLimits.MemoryLock != nil {
			limitsDict.AddInteger("MemoryLock", *daemon.ResourceLimits.MemoryLock)
		}
		if daemon.ResourceLimits.NumberOfProcesses != nil {
			limitsDict.AddInteger("NumberOfProcesses", *daemon.ResourceLimits.NumberOfProcesses)
		}
		if daemon.ResourceLimits.ResidentSetSize != nil {
			limitsDict.AddInteger("ResidentSetSize", *daemon.ResourceLimits.ResidentSetSize)
		}
		if daemon.ResourceLimits.Stack != nil {
			limitsDict.AddInteger("Stack", *daemon.ResourceLimits.Stack)
		}

		if len(limitsDict.Items) > 0 {
			dict.AddDict("SoftResourceLimits", limitsDict)
			dict.AddDict("HardResourceLimits", limitsDict) // Same as soft for now
		}
	}

	// Socket Activation
	if len(daemon.Sockets) > 0 {
		socketsDict := &Dict{}
		for name, socket := range daemon.Sockets {
			socketDict := &Dict{}

			if socket.SockType != "" {
				socketDict.AddString("SockType", socket.SockType)
			}
			if socket.SockPassive != nil {
				socketDict.AddBool("SockPassive", *socket.SockPassive)
			}
			if socket.SockNodeName != "" {
				socketDict.AddString("SockNodeName", socket.SockNodeName)
			}
			if socket.SockServiceName != "" {
				socketDict.AddString("SockServiceName", socket.SockServiceName)
			}
			if socket.SockFamily != "" {
				socketDict.AddString("SockFamily", socket.SockFamily)
			}
			if socket.SockProtocol != "" {
				socketDict.AddString("SockProtocol", socket.SockProtocol)
			}
			if socket.SockPathName != "" {
				socketDict.AddString("SockPathName", socket.SockPathName)
			}
			if socket.SockPathMode != nil {
				socketDict.AddInteger("SockPathMode", *socket.SockPathMode)
			}
			if socket.Bonjour != nil {
				socketDict.AddBool("Bonjour", *socket.Bonjour)
			}
			if len(socket.BonjourMultiple) > 0 {
				socketDict.AddStringArray("Bonjour", socket.BonjourMultiple)
			}

			socketsDict.AddDict(name, socketDict)
		}
		dict.AddDict("Sockets", socketsDict)
	}

	// Calendar Intervals
	if len(daemon.StartCalendarInterval) > 0 {
		if len(daemon.StartCalendarInterval) == 1 {
			// Single interval
			calDict := g.calendarIntervalToDict(daemon.StartCalendarInterval[0])
			if len(calDict.Items) > 0 {
				dict.AddDict("StartCalendarInterval", calDict)
			}
		} else {
			// Multiple intervals
			intervals := make([]*Dict, 0, len(daemon.StartCalendarInterval))
			for _, interval := range daemon.StartCalendarInterval {
				calDict := g.calendarIntervalToDict(interval)
				if len(calDict.Items) > 0 {
					intervals = append(intervals, calDict)
				}
			}
			if len(intervals) > 0 {
				dict.AddDictArray("StartCalendarInterval", intervals)
			}
		}
	}

	// Watch Paths
	if len(daemon.WatchPaths) > 0 {
		dict.AddStringArray("WatchPaths", daemon.WatchPaths)
	}
	if len(daemon.QueuePaths) > 0 {
		dict.AddStringArray("QueueDirectories", daemon.QueuePaths)
	}

	// Other Settings
	if daemon.EnableGlobbing {
		dict.AddBool("EnableGlobbing", daemon.EnableGlobbing)
	}
	if daemon.EnableTransactions {
		dict.AddBool("EnableTransactions", daemon.EnableTransactions)
	}
	if daemon.EnablePressuredExit {
		dict.AddBool("EnablePressuredExit", daemon.EnablePressuredExit)
	}
	if daemon.ExitTimeOut > 0 {
		dict.AddInteger("ExitTimeOut", daemon.ExitTimeOut)
	}

	return &Plist{
		Version: "1.0",
		Dict:    dict,
	}
}

// calendarIntervalToDict converts a calendar interval to a dict
func (g *Generator) calendarIntervalToDict(interval config.CalendarInterval) *Dict {
	dict := &Dict{}

	if interval.Minute != nil {
		dict.AddInteger("Minute", *interval.Minute)
	}
	if interval.Hour != nil {
		dict.AddInteger("Hour", *interval.Hour)
	}
	if interval.Day != nil {
		dict.AddInteger("Day", *interval.Day)
	}
	if interval.Weekday != nil {
		dict.AddInteger("Weekday", *interval.Weekday)
	}
	if interval.Month != nil {
		dict.AddInteger("Month", *interval.Month)
	}

	return dict
}
