#!/bin/bash

# Godot MCP Daemon Control Script

PLIST_NAME="com.mantle.godot-mcp"
PLIST_FILE="$(pwd)/com.mantle.godot-mcp.plist"
LAUNCH_AGENTS_DIR="$HOME/Library/LaunchAgents"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_usage() {
    echo "Usage: $0 {install|uninstall|start|stop|restart|status|logs|tail}"
    echo ""
    echo "Commands:"
    echo "  install    - Install the daemon (copies plist to LaunchAgents)"
    echo "  uninstall  - Uninstall the daemon"
    echo "  start      - Start the daemon"
    echo "  stop       - Stop the daemon"
    echo "  restart    - Restart the daemon"
    echo "  status     - Check daemon status"
    echo "  logs       - Show recent logs"
    echo "  tail       - Tail logs in real-time"
}

check_plist_exists() {
    if [ ! -f "$PLIST_FILE" ]; then
        echo -e "${RED}Error: Plist file not found at $PLIST_FILE${NC}"
        exit 1
    fi
}

is_installed() {
    [ -f "$LAUNCH_AGENTS_DIR/$PLIST_NAME.plist" ]
}

is_running() {
    launchctl list | grep -q "$PLIST_NAME"
}

install_daemon() {
    check_plist_exists
    
    if is_installed; then
        echo -e "${YELLOW}Daemon already installed${NC}"
        return
    fi
    
    echo "Installing daemon..."
    mkdir -p "$LAUNCH_AGENTS_DIR"
    cp "$PLIST_FILE" "$LAUNCH_AGENTS_DIR/"
    
    # Load the daemon
    launchctl load "$LAUNCH_AGENTS_DIR/$PLIST_NAME.plist"
    
    echo -e "${GREEN}Daemon installed successfully${NC}"
}

uninstall_daemon() {
    if ! is_installed; then
        echo -e "${YELLOW}Daemon not installed${NC}"
        return
    fi
    
    echo "Uninstalling daemon..."
    
    # Stop if running
    if is_running; then
        launchctl unload "$LAUNCH_AGENTS_DIR/$PLIST_NAME.plist"
    fi
    
    # Remove plist
    rm -f "$LAUNCH_AGENTS_DIR/$PLIST_NAME.plist"
    
    echo -e "${GREEN}Daemon uninstalled successfully${NC}"
}

start_daemon() {
    if ! is_installed; then
        echo -e "${RED}Error: Daemon not installed. Run '$0 install' first${NC}"
        exit 1
    fi
    
    if is_running; then
        echo -e "${YELLOW}Daemon already running${NC}"
        return
    fi
    
    echo "Starting daemon..."
    launchctl start "$PLIST_NAME"
    
    # Wait a moment and check status
    sleep 2
    if is_running; then
        echo -e "${GREEN}Daemon started successfully${NC}"
    else
        echo -e "${RED}Failed to start daemon. Check logs with '$0 logs'${NC}"
    fi
}

stop_daemon() {
    if ! is_running; then
        echo -e "${YELLOW}Daemon not running${NC}"
        return
    fi
    
    echo "Stopping daemon..."
    launchctl stop "$PLIST_NAME"
    
    echo -e "${GREEN}Daemon stopped${NC}"
}

restart_daemon() {
    stop_daemon
    sleep 2
    start_daemon
}

check_status() {
    echo "Daemon status:"
    echo "-------------"
    
    if is_installed; then
        echo -e "Installed: ${GREEN}YES${NC}"
    else
        echo -e "Installed: ${RED}NO${NC}"
    fi
    
    if is_running; then
        echo -e "Running: ${GREEN}YES${NC}"
        echo ""
        echo "Process info:"
        launchctl list | grep "$PLIST_NAME"
    else
        echo -e "Running: ${RED}NO${NC}"
    fi
}

show_logs() {
    LOG_DIR="$(pwd)/logs"
    
    if [ ! -d "$LOG_DIR" ]; then
        echo -e "${RED}Log directory not found${NC}"
        return
    fi
    
    echo "=== STDOUT ==="
    if [ -f "$LOG_DIR/stdout.log" ]; then
        tail -n 50 "$LOG_DIR/stdout.log"
    else
        echo "(empty)"
    fi
    
    echo ""
    echo "=== STDERR ==="
    if [ -f "$LOG_DIR/stderr.log" ]; then
        tail -n 50 "$LOG_DIR/stderr.log"
    else
        echo "(empty)"
    fi
}

tail_logs() {
    LOG_DIR="$(pwd)/logs"
    
    if [ ! -d "$LOG_DIR" ]; then
        echo -e "${RED}Log directory not found${NC}"
        return
    fi
    
    echo "Tailing logs (Ctrl+C to stop)..."
    tail -f "$LOG_DIR/stdout.log" "$LOG_DIR/stderr.log"
}

# Main script logic
case "$1" in
    install)
        install_daemon
        ;;
    uninstall)
        uninstall_daemon
        ;;
    start)
        start_daemon
        ;;
    stop)
        stop_daemon
        ;;
    restart)
        restart_daemon
        ;;
    status)
        check_status
        ;;
    logs)
        show_logs
        ;;
    tail)
        tail_logs
        ;;
    *)
        print_usage
        exit 1
        ;;
esac