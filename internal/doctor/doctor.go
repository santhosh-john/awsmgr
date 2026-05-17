package doctor

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/santhoshjohn/awsmgr/internal/aws"
	"github.com/santhoshjohn/awsmgr/internal/kube"
	"github.com/santhoshjohn/awsmgr/internal/network"
	"github.com/santhoshjohn/awsmgr/internal/terraform"
)

const stepTimeout = 8 * time.Second

type Status string

const (
	StatusOK   Status = "OK"
	StatusFail Status = "FAIL"
	StatusWarn Status = "WARN"
)

var nonAlphaNumeric = regexp.MustCompile(`[^a-z0-9]+`)

var (
	validateAWSProfile = aws.ValidateProfile
	currentKubeContext = kube.CurrentContext
	terraformAvailable = terraform.Available
	networkReachable   = network.Reachable
)

type Check struct {
	Name    string
	Status  Status
	OK      bool
	Message string
	Err     error
}

func (c Check) ResultStatus() Status {
	if c.Status != "" {
		return c.Status
	}
	if c.OK {
		return StatusOK
	}
	return StatusFail
}

type Report struct {
	Checks []Check
}

func (r Report) OK() bool {
	for _, check := range r.Checks {
		if check.ResultStatus() == StatusFail {
			return false
		}
	}
	return true
}

func Run(ctx context.Context, profile string) Report {
	awsResult := validateAWSProfile(ctx, profile, stepTimeout)
	kubeResult := currentKubeContext(ctx, stepTimeout)

	checks := []Check{
		checkAWS(profile, awsResult),
		checkKubernetes(kubeResult),
		checkTerraform(ctx),
		checkNetwork(ctx),
	}
	if alignmentCheck, ok := checkAlignment(profile, awsResult, kubeResult); ok {
		checks = append(checks, alignmentCheck)
	}
	return Report{Checks: checks}
}

func checkAWS(profile string, result aws.ValidationResult) Check {
	if result.Valid {
		return Check{
			Name:    "AWS",
			Status:  StatusOK,
			OK:      true,
			Message: fmt.Sprintf("profile=%s account=%s", profile, result.AccountID),
		}
	}
	message := result.Message
	if message == "" {
		message = fmt.Sprintf("profile=%s invalid", profile)
	}
	return Check{Name: "AWS", Status: StatusFail, Message: message, Err: result.Err}
}

func checkKubernetes(result kube.ContextResult) Check {
	if result.OK {
		return Check{Name: "Kubernetes", Status: StatusOK, OK: true, Message: fmt.Sprintf("context=%s", result.Context)}
	}
	return Check{Name: "Kubernetes", Status: StatusFail, Message: result.Message, Err: result.Err}
}

func checkTerraform(ctx context.Context) Check {
	result := terraformAvailable(ctx, stepTimeout)
	if result.OK {
		return Check{Name: "Terraform", Status: StatusOK, OK: true, Message: "terraform available"}
	}
	return Check{Name: "Terraform", Status: StatusFail, Message: result.Message, Err: result.Err}
}

func checkNetwork(ctx context.Context) Check {
	result := networkReachable(ctx, "aws.amazon.com:443", stepTimeout)
	if result.OK {
		return Check{Name: "Network", Status: StatusOK, OK: true, Message: "network reachable"}
	}
	return Check{Name: "Network", Status: StatusFail, Message: "network unreachable", Err: result.Err}
}

func checkAlignment(profile string, awsResult aws.ValidationResult, kubeResult kube.ContextResult) (Check, bool) {
	if !awsResult.Valid || !kubeResult.OK {
		return Check{}, false
	}

	contextName := kubeResult.Context
	if strings.Contains(strings.ToLower(contextName), strings.ToLower(awsResult.AccountID)) {
		return Check{
			Name:    "Alignment",
			Status:  StatusOK,
			OK:      true,
			Message: fmt.Sprintf("kube context matches AWS account %s", awsResult.AccountID),
		}, true
	}

	profileToken := normalizeToken(profile)
	contextToken := normalizeToken(contextName)
	if profileToken != "" && strings.Contains(contextToken, profileToken) {
		return Check{
			Name:    "Alignment",
			Status:  StatusOK,
			OK:      true,
			Message: fmt.Sprintf("kube context name appears aligned with profile %s", profile),
		}, true
	}

	return Check{
		Name:    "Alignment",
		Status:  StatusWarn,
		OK:      true,
		Message: fmt.Sprintf("context=%s does not appear to match profile=%s or account=%s", contextName, profile, awsResult.AccountID),
	}, true
}

func normalizeToken(value string) string {
	normalized := strings.ToLower(value)
	normalized = nonAlphaNumeric.ReplaceAllString(normalized, "")
	return normalized
}
