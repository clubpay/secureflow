package main

import "fmt"

func main() {
    // AWS Access Key ID
    awsAccessKeyID := "AKIAIOSFODNN7EXAMPLE"
    fmt.Println("AWS Access Key ID:", awsAccessKeyID)

    // AWS Secret Access Key
    awsSecretAccessKey := "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
    fmt.Println("AWS Secret Access Key:", awsSecretAccessKey)

    // Database Connection String (with password)
    dbConnectionString := "postgres://user:password123@localhost:5432/mydb"
    fmt.Println("Database Connection String:", dbConnectionString)

    // GitHub Personal Access Token
    githubToken := "ghp_aBcDeFgHiJkLmNoPqRsTuVwXyZ1234567890"
    fmt.Println("GitHub Token:", githubToken)

    // Basic Auth Password
    basicAuthPassword := "supersecretpassword!"
    fmt.Println("Basic Auth Password:", basicAuthPassword)

    // Slack Webhook URL
    slackWebhookURL := "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"
    fmt.Println("Slack Webhook URL:", slackWebhookURL)
}
