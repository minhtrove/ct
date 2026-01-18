// Package email handles sending emails using AWS SES
package email

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/minhtranin/ct/internal/logger"
)

// SESClient is the email client for AWS SES
type SESClient struct {
	client    *sesv2.Client
	fromEmail string
	fromName  string
	region    string
}

// SendEmailInput wraps the SES SendEmailInput for use by handlers
type SendEmailInput struct {
	ToEmailAddress   string
	FromEmailAddress string
	Content          *EmailContent
}

// EmailContent wraps the SES EmailContent
type EmailContent struct {
	Simple *Message
}

// Message wraps the SES Message
type Message struct {
	Subject *Content
	Body    *Body
}

// Content wraps the SES Content
type Content struct {
	Data    string
	Charset string
}

// Body wraps the SES Body
type Body struct {
	Html *Content
	Text *Content
}

// NewSESClient creates a new AWS SES email client
func NewSESClient() (*SESClient, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1" // Default region
	}

	fromEmail := os.Getenv("SES_FROM_EMAIL")
	fromName := os.Getenv("SES_FROM_NAME")

	// Set defaults if env vars not provided
	if fromEmail == "" {
		fromEmail = "noreply@example.com"
	}
	if fromName == "" {
		appName := os.Getenv("APP_NAME")
		if appName != "" {
			fromName = appName
		} else {
			fromName = "CT App"
		}
	}

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		logger.Error("Email", "Failed to load AWS config: "+err.Error())
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &SESClient{
		client:    sesv2.NewFromConfig(cfg),
		fromEmail: fromEmail,
		fromName:  fromName,
		region:    region,
	}, nil
}

// SendEmail sends an email using AWS SES
func (c *SESClient) SendEmail(input *SendEmailInput) (*sesv2.SendEmailOutput, error) {
	// Convert to SES format
	sesInput := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(input.FromEmailAddress),
		Destination: &types.Destination{
			ToAddresses: []string{input.ToEmailAddress},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data:    aws.String(input.Content.Simple.Subject.Data),
					Charset: aws.String(input.Content.Simple.Subject.Charset),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data:    aws.String(input.Content.Simple.Body.Html.Data),
						Charset: aws.String(input.Content.Simple.Body.Html.Charset),
					},
					Text: &types.Content{
						Data:    aws.String(input.Content.Simple.Body.Text.Data),
						Charset: aws.String(input.Content.Simple.Body.Text.Charset),
					},
				},
			},
		},
	}

	return c.client.SendEmail(context.TODO(), sesInput)
}

// FormatFromAddress formats the from address with name
func (c *SESClient) FormatFromAddress() string {
	if c.fromName != "" {
		return fmt.Sprintf("%s <%s>", c.fromName, c.fromEmail)
	}
	return c.fromEmail
}

// GetBaseURL returns the base URL from environment or defaults to localhost
func (c *SESClient) GetBaseURL() string {
	baseURL := os.Getenv("APP_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}
	return baseURL
}

// GetClient returns the underlying SES client (for advanced usage)
func (c *SESClient) GetClient() *sesv2.Client {
	return c.client
}
