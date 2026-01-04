# bashlog

A lightweight, powerful logging library for Bash scripts with support for multiple log levels, formatted output, and flexible configuration.

## Features

- 🎯 **Multiple Log Levels**: DEBUG, INFO, WARN, ERROR, and FATAL levels
- 📝 **Formatted Output**: Color-coded log messages with timestamps
- ⚙️ **Configurable**: Easy to customize log levels, formats, and output destinations
- 🎨 **Color Support**: Beautiful colored output for better readability
- 📤 **Multiple Outputs**: Log to console, files, or both simultaneously
- 🔍 **Debug Mode**: Enable/disable debug logging with a simple flag
- 📊 **Structured Logging**: Consistent log format across your scripts
- ✨ **Lightweight**: Minimal dependencies, pure Bash implementation
- 🛡️ **Error Handling**: Robust error handling and logging capabilities

## Installation

### Method 1: Clone the Repository

```bash
git clone https://github.com/interhack86/bashlog.git
cd bashlog
```

### Method 2: Download as a Source File

```bash
# Download the bashlog script
curl -O https://raw.githubusercontent.com/interhack86/bashlog/main/bashlog.sh
chmod +x bashlog.sh
```

### Method 3: Installation to System Path

```bash
git clone https://github.com/interhack86/bashlog.git
cd bashlog
sudo make install
```

## Quick Start

### Basic Usage

```bash
#!/bin/bash

# Source the bashlog library
source ./bashlog.sh

# Initialize logging (optional - sets log file)
log_init "/var/log/myapp.log"

# Use different log levels
log_debug "This is a debug message"
log_info "Application started successfully"
log_warn "This is a warning message"
log_error "An error occurred"
log_fatal "Fatal error - exiting"
```

## Usage Examples

### Simple Logging

```bash
#!/bin/bash
source ./bashlog.sh

log_info "Starting backup process..."
log_debug "Backup directory: /home/user/data"
log_warn "Some files may be skipped"
log_error "Failed to connect to backup server"
```

### With File Logging

```bash
#!/bin/bash
source ./bashlog.sh

# Initialize with log file
log_init "/tmp/myapp.log"
log_enable_file_logging

log_info "Application initialized"
log_error "Database connection failed"
```

### Debug Mode

```bash
#!/bin/bash
source ./bashlog.sh

# Enable debug mode
log_set_level DEBUG

log_debug "Detailed debugging information"
log_info "Normal information message"

# Disable debug mode
log_set_level INFO
log_debug "This won't be displayed"
```

### In Functions

```bash
#!/bin/bash
source ./bashlog.sh

backup_database() {
    local db_name="$1"
    
    log_info "Starting backup of database: $db_name"
    
    if ! mysqldump "$db_name" > "/tmp/$db_name.sql" 2>/dev/null; then
        log_error "Failed to backup database: $db_name"
        return 1
    fi
    
    log_info "Database backup completed successfully"
    return 0
}

backup_database "myapp_db"
```

### Error Handling and Logging

```bash
#!/bin/bash
source ./bashlog.sh

log_init "/var/log/deploy.log"

deploy_app() {
    log_info "Deployment started"
    
    if [[ ! -d "./app" ]]; then
        log_error "Application directory not found"
        log_fatal "Deployment failed - critical error"
        exit 1
    fi
    
    log_info "Deploying application..."
    cp -r ./app /opt/myapp || {
        log_error "Failed to copy application files"
        return 1
    }
    
    log_info "Deployment completed successfully"
}

deploy_app
```

### Conditional Logging

```bash
#!/bin/bash
source ./bashlog.sh

CONFIG_DEBUG=${CONFIG_DEBUG:-false}

if $CONFIG_DEBUG; then
    log_set_level DEBUG
else
    log_set_level INFO
fi

log_debug "Debug mode: $CONFIG_DEBUG"
log_info "Application configuration loaded"
```

### Multiple Log Destinations

```bash
#!/bin/bash
source ./bashlog.sh

# Log to file
log_init "/var/log/myapp.log"
log_enable_file_logging

# Log to console (default)
log_enable_console_logging

log_info "This message goes to both console and file"
log_error "Error logged to all destinations"
```

## Configuration

### Log Levels

Configure the minimum log level to display:

```bash
# Set log level (options: DEBUG, INFO, WARN, ERROR, FATAL)
log_set_level DEBUG
log_set_level INFO
log_set_level ERROR
```

### Log File

Initialize logging with a log file path:

```bash
log_init "/path/to/logfile.log"
```

### Output Control

```bash
# Enable/disable file logging
log_enable_file_logging
log_disable_file_logging

# Enable/disable console output
log_enable_console_logging
log_disable_console_logging

# Set custom date format (default: %Y-%m-%d %H:%M:%S)
log_set_date_format "%Y-%m-%d %T"
```

### Custom Log Prefix

```bash
# Set a custom prefix for all log messages
log_set_prefix "[MYAPP]"
```

## Log Levels Explained

| Level | Severity | Use Case |
|-------|----------|----------|
| DEBUG | Low | Detailed debugging information, variable values |
| INFO | Low | General information, application flow |
| WARN | Medium | Warning messages, deprecated features |
| ERROR | High | Error conditions, recoverable failures |
| FATAL | Critical | Fatal errors, non-recoverable failures |

## API Reference

### Core Functions

| Function | Description |
|----------|-------------|
| `log_debug(message)` | Log a debug message |
| `log_info(message)` | Log an info message |
| `log_warn(message)` | Log a warning message |
| `log_error(message)` | Log an error message |
| `log_fatal(message)` | Log a fatal message |

### Configuration Functions

| Function | Description |
|----------|-------------|
| `log_init(filepath)` | Initialize logging with a file path |
| `log_set_level(level)` | Set the minimum log level |
| `log_set_prefix(prefix)` | Set a custom log prefix |
| `log_set_date_format(format)` | Set custom date format |
| `log_enable_file_logging()` | Enable file logging |
| `log_disable_file_logging()` | Disable file logging |
| `log_enable_console_logging()` | Enable console output |
| `log_disable_console_logging()` | Disable console output |

## Advanced Examples

### Creating a Production-Ready Script

```bash
#!/bin/bash
set -euo pipefail

source ./bashlog.sh

# Configuration
SCRIPT_NAME="$(basename "$0")"
LOG_FILE="/var/log/${SCRIPT_NAME%.*}.log"
DEBUG=${DEBUG:-false}

# Initialize logging
log_init "$LOG_FILE"
log_enable_file_logging
log_enable_console_logging
log_set_prefix "[$SCRIPT_NAME]"

if $DEBUG; then
    log_set_level DEBUG
else
    log_set_level INFO
fi

log_info "Script started (PID: $$)"

# Your script logic here
main() {
    log_info "Executing main function"
    
    # Do something
    log_debug "Processing complete"
    
    log_info "Script completed successfully"
}

# Error handler
trap 'log_error "Script failed at line $LINENO"; exit 1' ERR

main "$@"
```

### Integration with Cron Jobs

```bash
#!/bin/bash
# /usr/local/bin/backup.sh
source /opt/bashlog/bashlog.sh

LOG_FILE="/var/log/backup.log"
log_init "$LOG_FILE"

# Suppress console output for cron jobs
log_disable_console_logging
log_enable_file_logging

log_info "Backup job started"

# Backup operations
# ...

log_info "Backup job completed"
```

## Performance

- **Minimal Overhead**: Pure Bash implementation with no external dependencies
- **Efficient Output**: Buffered logging for better performance
- **Lightweight**: < 5KB file size
- **No Subshells**: Avoids unnecessary subshell spawning

## Compatibility

- **Bash Version**: 4.0+ recommended
- **Systems**: Linux, macOS, BSD, WSL
- **Shell**: Compatible with bash and sh (for basic functions)

## Troubleshooting

### Colors Not Showing

Ensure your terminal supports ANSI colors:

```bash
# Check if colors are supported
echo $TERM
```

To disable colors:

```bash
export NO_COLOR=1
```

### Permission Denied

Make sure the log file has proper permissions:

```bash
touch /var/log/myapp.log
sudo chown $USER:$USER /var/log/myapp.log
chmod 644 /var/log/myapp.log
```

### Log File Not Created

Ensure the directory exists:

```bash
mkdir -p /var/log
sudo touch /var/log/myapp.log
```

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Changelog

### Version 1.0.0
- Initial release
- Basic logging functionality
- Multiple log levels
- File and console output
- Color support

## Support

For issues, questions, or suggestions, please open an [issue](https://github.com/interhack86/bashlog/issues) on GitHub.

## Author

**interhack86** - [GitHub Profile](https://github.com/interhack86)

---

**Happy Logging! 🎉**
