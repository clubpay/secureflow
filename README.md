# Automated Security Scanning Workflows

This repository contains **reusable GitHub workflows** for automated security scanning, including:

- **Static Application Security Testing (SAST)**  
- **Software Composition Analysis (SCA)**  
- **Infrastructure as Code (IaC) Scanning**  
- **Secret Detection**  

These workflows are designed to streamline security checks in your CI/CD pipelines. Remarkably, they complete scans in under a minute, providing swift and efficient protection for your codebase.
![SecureFlow Arch](https://github.com/user-attachments/assets/062ad6e9-0f62-43a6-ab2d-e4139875fec0)

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
  push:  # This will trigger the workflow on every push to the repository

jobs:
  sast:
    uses: clubpay/secureflow/.github/workflows/sast.yml@main
    secrets: inherit
    
  sca:
    uses: clubpay/secureflow/.github/workflows/sca.yml@main
    secrets: inherit
    
  iac-scanning:
    uses: clubpay/secureflow/.github/workflows/iac-scanning.yml@main
    secrets: inherit
    
  secret-detection:
    uses: clubpay/secureflow/.github/workflows/secret-detection.yml@main
    secrets: inherit

```
