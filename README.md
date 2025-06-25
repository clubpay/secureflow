# Automated Security Scanning Workflows

This repository contains **reusable GitHub workflows** for automated security scanning, including:

- **Static Application Security Testing (SAST)**  
- **Software Composition Analysis (SCA)**  
- **Infrastructure as Code (IaC) Scanning**  
- **Secret Detection**  

These workflows are designed to streamline security checks in your CI/CD pipelines. Remarkably, they complete scans in under a minute, providing swift and efficient protection for your codebase.
![secureflow_arch_720](https://github.com/user-attachments/assets/3e6152c0-4463-4b03-bbf8-43cc8deea490)

## Overview of Scanning Types

Below is a detailed explanation of each scan type:

### 1. **SAST (Static Application Security Testing)**

SAST is a white-box testing method that analyzes your source code or binaries for vulnerabilities without executing the application. It helps identify issues such as:

- **Injection vulnerabilities** (e.g., SQL Injection, Command Injection)
- **Cross-site Scripting (XSS)**
- **Insecure deserialization**

SAST scans are typically run early in the development lifecycle, ensuring security flaws are caught before deployment.

---

### 2. **SCA (Software Composition Analysis)**

SCA focuses on identifying vulnerabilities in third-party libraries, frameworks, and dependencies used in your project. It also helps with:

- Checking for outdated or vulnerable libraries.
- Licensing issues to ensure compliance.
- Understanding the risk profile of your software supply chain.

SCA is critical for managing risks from open-source components and external libraries.  

---

### 3. **IaC Scanning (Infrastructure as Code Scanning)**

IaC Scanning analyzes your configuration files for cloud infrastructure, such as:

- Terraform
- Kubernetes manifests
- Dockerfiles

It helps detect security risks, such as:

- Misconfigured access controls
- Open ports or insecure network configurations
- Lack of encryption or improper key management

By scanning IaC files, you can prevent deploying vulnerable infrastructure to production.  

---

### 4. **Secret Detection**

Secret Detection is the process of scanning for hardcoded sensitive information in the codebase, such as:

- API keys
- Authentication tokens
- Passwords
- Certificates or private keys

Exposing such secrets can lead to unauthorized access to sensitive systems or data breaches.  

---

Each of these scanning methods targets specific aspects of your application’s security, and together they provide a comprehensive security posture for your project.


## Available Workflows

The repository provides the following reusable workflows for automated security scanning:

| **Workflow**       | **Purpose**                        | **Tool(s) Used**       |
|---------------------|------------------------------------|------------------------|
| `sast.yml`          | Static Application Security Testing (SAST) | [Semgrep](https://github.com/semgrep/semgrep) |
| `sca.yml`           | Software Composition Analysis (SCA)       | [Grype](https://github.com/anchore/grype) |
| `iac-scanning.yml`           | Infrastructure as Code (IaC) Scanning    | [KICS](https://github.com/Checkmarx/kics) |
| `secret-detection.yml`       | Secret Detection                   | [Gitleaks](https://github.com/gitleaks/gitleaks) |

## How to Use the Workflows

These workflows can be reused in your GitHub repositories to automate security scanning and you can easily integrate them into your GitHub jobs by referencing their name and path in your repository. Below is an example of how to integrate the **Secure Workflows**:

```yaml
name: Trigger SecureFlow
on:
  push: # This will trigger the workflow on every push to the repository

jobs:
  sast:
    uses: clubpay/secureflow/.github/workflows/sast.yml@main
    secrets:
      GLOBAL_REPO_TOKEN: ${{ secrets.GLOBAL_REPO_TOKEN }}
      DEFECTDOJO_TOKEN: ${{ secrets.DEFECTDOJO_TOKEN }}
    permissions:
      id-token: write #This is essential for authentication to Teleport

  sca:
    uses: clubpay/secureflow/.github/workflows/sca.yml@main
    secrets:
      GLOBAL_REPO_TOKEN: ${{ secrets.GLOBAL_REPO_TOKEN }}
      DEFECTDOJO_TOKEN: ${{ secrets.DEFECTDOJO_TOKEN }}
    permissions:
      id-token: write

  iac-scanning:
    uses: clubpay/secureflow/.github/workflows/iac-scanning.yml@main
    secrets:
      GLOBAL_REPO_TOKEN: ${{ secrets.GLOBAL_REPO_TOKEN }}
      DEFECTDOJO_TOKEN: ${{ secrets.DEFECTDOJO_TOKEN }}
    permissions:
      id-token: write

  secret-detection:
    uses: clubpay/secureflow/.github/workflows/secret-detection.yml@main
    secrets:
      GLOBAL_REPO_TOKEN: ${{ secrets.GLOBAL_REPO_TOKEN }}
      DEFECTDOJO_TOKEN: ${{ secrets.DEFECTDOJO_TOKEN }}
    permissions:
      id-token: write

```
## Skipping Files or Lines from Scanning

You can skip certain files, folders, or lines from being scanned by SecureFlow.

Important Considerations:

- Documentation: Always add a comment explaining why you're suppressing the rule. This helps other developers understand your decision and avoids accidentally re-introducing the vulnerability.
- Review: Regularly review your skipped codes to ensure they still need to be skipped. Security landscapes change, and previously acceptable suppressions might become risky.
- Alternative: If possible, prefer fixing the underlying vulnerability rather than excluding the code. Skipping should be a last resort.


### 1. Skip Specific Lines

If you want to ignore specific lines of code, you can add in-line comments at the end of the target line of code. Use **nosemgrep** comment to skip SAST scans, and **gitleaks:allow** comment to skip secret detection.

Skipping SAST Scans (Example 1):
```python
import os

def get_user_input():
  user_input = input("Enter something: ") #nosemgrep
  print(f"You entered: {user_input}")
  return user_input

if __name__ == "__main__":
  get_user_input()
```
Skipping SAST Scans (Example 2):
```go
package main

import (
    "fmt"
    "os"
)

func main() {
    command := os.Getenv("MY_COMMAND") //nosemgrep
    fmt.Println("Executing:", command)
}
```
Skipping Secrets (Example 1):
```go
package main
import "fmt"

func main() {
    apiKey := "THIS_IS_A_SAMPLE_API_KEY" //gitleaks:allow 
    fmt.Println("API Key:", apiKey)
}
```
Skipping Secrets (Example 2):
```python
DB_PASSWORD = "TestDBPassword"  #gitleaks:allow

def my_function():
#sample funcion to read from DB

my_function()
```
### 2. Skip Files and Directories

You can exclude specific directories or files from SAST and secret detection scans by using the **SAST_EXCLUDE_LIST** and **SECRET_DETECTION_EXCLUDE_LIST** variables. These variables accept a space-separated list of file and directory names that should be ignored during scanning. To apply these exclusions, you can reuse the workflows **with** these variables.
```yaml
jobs:
  sast:
    uses: clubpay/secureflow/.github/workflows/sast.yml@main
    with:
      SAST_EXCLUDE_LIST: "community/certs/server.key README.md MyAwsomeDirectory stage.env"
    secrets:
      GLOBAL_REPO_TOKEN: ${{ secrets.GLOBAL_REPO_TOKEN }}
      DEFECTDOJO_TOKEN: ${{ secrets.DEFECTDOJO_TOKEN }}
    permissions:
      id-token: write

  secret-detection:
    uses: clubpay/secureflow/.github/workflows/secret-detection.yml@main
    with:
      SECRET_DETECTION_EXCLUDE_LIST: "services/community.go mock.go README.md AwsomeDirectory"
    secrets:
      GLOBAL_REPO_TOKEN: ${{ secrets.GLOBAL_REPO_TOKEN }}
      DEFECTDOJO_TOKEN: ${{ secrets.DEFECTDOJO_TOKEN }}
    permissions:
      id-token: write
```
## Note for Private Repositories
For private repositories, it's essential to configure an action secret named **GLOBAL_REPO_TOKEN** with the appropriate permissions. Ensure that the token has both repository (repo) and workflow (workflow) access, allowing GitHub Actions to authenticate and execute workflows smoothly. Without this, attempts to access private repositories during checkout or workflow execution will fail due to insufficient authorization.
