# JFrog Applications Config

The JFrog Applications Config schema is used to define the configuration schema used by some for the JFrog applications, such as JFrog CLI, JFrog Frogbot and the JFrog IDE integrations. This configuration schema is used to define the rules and settings for the JFrog source code scanning tools.

By consolidating the relevant settings, rules and policies into a single file, developers and security teams can easily manage and update scanning configurations, ensuring consistent and effective code analysis.

## Project status

[![Scanned by Frogbot](https://raw.github.com/jfrog/frogbot/master/images/frogbot-badge.svg)](https://github.com/jfrog/frogbot#readme)
[![Test](https://github.com/jfrog/jfrog-apps-config/actions/workflows/test.yml/badge.svg)](https://github.com/jfrog/jfrog-apps-config/actions/workflows/test.yml)
[![Static Analysis](https://github.com/jfrog/jfrog-apps-config/actions/workflows/analysis.yml/badge.svg)](https://github.com/jfrog/jfrog-apps-config/actions/workflows/analysis.yml)

## Schema:

```yaml
# [Required] JFrog Applications Config version
version: "1.0"

modules:
  # [Required] Module name
  - name: FrogLeapApp
    # [Optional, default: "."] Application's root directory
    source_root: "src"
    # [Optional] Directories to exclude from scanning across all scanners
    exclude_patterns:
      - "docs/"
    # [Optional] Scanners to exclude from JFrog Advanced Security (Options: "secrets", "sast", "iac")
    exclude_scanners:
      - secrets
    # [Optional] Customize scanner configurations
    scanners:
      # [Optional] Configuration for Static Application Security Testing (SAST)
      sast:
        # [Optional] Specify the programming language for SAST
        language: java
        # [Optional] Working directories specific to SAST (Relative to source_root)
        working_dirs:
          - "dir1"
          - "dir2"
        # [Optional] Additional exclude patterns for this scanner
        exclude_patterns:
          - "dir1/test/**"
        # [Optional] List of specific scan rules to exclude from the scan
        excluded_rules:
          - xss-injection

      # [Optional] Configuration for secrets scan
      secrets:
        # [Optional] Working directories specific to the secret scanner (Relative to source_root)
        working_dirs:
          - "dir1"
          - "dir2"
        # [Optional] Additional exclude patterns for this scanner
        exclude_patterns:
          - "dir1/test/**"

      # [Optional] Configuration for Infrastructure as Code scan (IaC)
      iac:
        # [Optional] Working directories specific to IaC (Relative to source_root)
        working_dirs:
          - "dir1"
          - "dir2"
        # [Optional] Additional exclude patterns for this Scanner
        exclude_patterns:
          - "dir1/test/**"
```
