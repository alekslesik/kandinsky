# Config Package

## Overview
The `config` package is a comprehensive and user-friendly configuration management solution for Go applications. It offers seamless integration of environment variables and `.env` file support, along with planned features for hot-reloading, multiple format support, and external system integration.

## Features
- **Environment Variables**: Leverage environment variables for configuration, with built-in default values.
- **`.env` File Support**: Manage your configurations conveniently through a `.env` file.
- **Planned Features**:
  - Hot-Reloading: Dynamically reload configurations without restarting the application.
  - Multiple Format Support: Extendable to support YAML, JSON, TOML formats.
  - External System Integration: Facilitate integration with systems like Consul, etcd for distributed configurations.
  - Validation: Validate configuration values to ensure reliability and correctness.

## Installation
To install the package, execute the following command:
```
go get github.com/alekslesik/config
```
## Usage
To use the package in your Go application, follow these steps:

1. Import the package:
   ```
   import "github.com/alekslesik/config"
   ```
2. Define a configuration structure:
   ```
   type AppConfig struct {
    Port int `env:"PORT" env-default:"8080"`
    // ... other configuration fields
    }
3. Load the configuration:

    ```
    var cfg AppConfig
    err := config.Load(&cfg)
    if err != nil {
        log.Fatalf("Error loading config: %v", err)
    }
    ```

## Documentation

Refer to the config documentation for detailed information on all functionalities of the package.
Contributing

Contributions are welcome! If you wish to contribute to the project, please fork the repository and submit a pull request. For substantial changes, please open an issue first to discuss what you would like to change.
License

This project is licensed under the MIT License, which permits both personal and commercial use, modification, distribution, and sublicensing.

## Future Enhancements

 Hot-reloading of configurations.
 Support for additional formats like YAML, JSON, TOML.
 Integration with external configuration management systems.
 Configuration validation.

