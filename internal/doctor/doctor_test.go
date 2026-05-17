package doctor

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/santhosh-john/awsmgr/internal/aws"
	"github.com/santhosh-john/awsmgr/internal/kube"
	"github.com/santhosh-john/awsmgr/internal/network"
	"github.com/santhosh-john/awsmgr/internal/terraform"
)

func TestReportOK(t *testing.T) {
	report := Report{Checks: []Check{{Name: "AWS", OK: true}, {Name: "Network", OK: true}}}
	if !report.OK() {
		t.Fatal("Report.OK() = false, want true")
	}

	report.Checks[1].OK = false
	if report.OK() {
		t.Fatal("Report.OK() = true, want false")
	}
}

func TestRunAggregatesDependencyResults(t *testing.T) {
	restore := stubDoctorDependencies(
		func(context.Context, string, time.Duration) aws.ValidationResult {
			return aws.ValidationResult{Valid: true, AccountID: "123456789012"}
		},
		func(context.Context, time.Duration) kube.ContextResult {
			return kube.ContextResult{OK: true, Context: "prod-cluster"}
		},
		func(context.Context, time.Duration) terraform.AvailabilityResult {
			return terraform.AvailabilityResult{OK: true}
		},
		func(context.Context, string, time.Duration) network.ReachabilityResult {
			return network.ReachabilityResult{OK: true}
		},
	)
	defer restore()

	report := Run(context.Background(), "prod")

	got := report.Checks
	want := []Check{
		{Name: "AWS", Status: StatusOK, OK: true, Message: "profile=prod account=123456789012"},
		{Name: "Kubernetes", Status: StatusOK, OK: true, Message: "context=prod-cluster"},
		{Name: "Terraform", Status: StatusOK, OK: true, Message: "terraform available"},
		{Name: "Network", Status: StatusOK, OK: true, Message: "network reachable"},
		{Name: "Alignment", Status: StatusOK, OK: true, Message: "kube context name appears aligned with profile prod"},
	}

	if len(got) != len(want) {
		t.Fatalf("len(Run().Checks) = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i].Name != want[i].Name || got[i].OK != want[i].OK || got[i].Message != want[i].Message || got[i].Err != nil {
			t.Fatalf("Run().Checks[%d] = %#v, want %#v", i, got[i], want[i])
		}
	}
	if !report.OK() {
		t.Fatal("Run().OK() = false, want true")
	}
}

func TestRunMarksFailures(t *testing.T) {
	terraformErr := errors.New("terraform missing")
	restore := stubDoctorDependencies(
		func(context.Context, string, time.Duration) aws.ValidationResult {
			return aws.ValidationResult{Valid: true, AccountID: "123456789012"}
		},
		func(context.Context, time.Duration) kube.ContextResult {
			return kube.ContextResult{OK: true, Context: "prod-cluster"}
		},
		func(context.Context, time.Duration) terraform.AvailabilityResult {
			return terraform.AvailabilityResult{Message: "terraform unavailable", Err: terraformErr}
		},
		func(context.Context, string, time.Duration) network.ReachabilityResult {
			return network.ReachabilityResult{OK: true}
		},
	)
	defer restore()

	report := Run(context.Background(), "prod")

	if report.OK() {
		t.Fatal("Run().OK() = true, want false")
	}
	if got := report.Checks[2]; got.OK || got.Message != "terraform unavailable" || !errors.Is(got.Err, terraformErr) {
		t.Fatalf("Run() terraform check = %#v, want failed terraform check", got)
	}
}

func TestRunUsesAWSFailureMessage(t *testing.T) {
	awsErr := errors.New("sso expired")
	restore := stubDoctorDependencies(
		func(context.Context, string, time.Duration) aws.ValidationResult {
			return aws.ValidationResult{Message: "profile=prod AWS SSO session expired or invalid", Err: awsErr}
		},
		func(context.Context, time.Duration) kube.ContextResult {
			return kube.ContextResult{OK: true, Context: "prod-cluster"}
		},
		func(context.Context, time.Duration) terraform.AvailabilityResult {
			return terraform.AvailabilityResult{OK: true}
		},
		func(context.Context, string, time.Duration) network.ReachabilityResult {
			return network.ReachabilityResult{OK: true}
		},
	)
	defer restore()

	report := Run(context.Background(), "prod")

	if got := report.Checks[0]; got.OK || got.Message != "profile=prod AWS SSO session expired or invalid" || !errors.Is(got.Err, awsErr) {
		t.Fatalf("Run() AWS check = %#v, want actionable AWS failure", got)
	}
}

func TestRunAddsAlignmentWarningForMismatch(t *testing.T) {
	restore := stubDoctorDependencies(
		func(context.Context, string, time.Duration) aws.ValidationResult {
			return aws.ValidationResult{Valid: true, AccountID: "123456789012"}
		},
		func(context.Context, time.Duration) kube.ContextResult {
			return kube.ContextResult{OK: true, Context: "dev-cluster"}
		},
		func(context.Context, time.Duration) terraform.AvailabilityResult {
			return terraform.AvailabilityResult{OK: true}
		},
		func(context.Context, string, time.Duration) network.ReachabilityResult {
			return network.ReachabilityResult{OK: true}
		},
	)
	defer restore()

	report := Run(context.Background(), "prod")

	if !report.OK() {
		t.Fatal("Run().OK() = false, want true for warning-only mismatch")
	}
	if got := report.Checks[4]; got.ResultStatus() != StatusWarn || got.Message != "context=dev-cluster does not appear to match profile=prod or account=123456789012" {
		t.Fatalf("Run() alignment check = %#v, want warning mismatch", got)
	}
}

func stubDoctorDependencies(
	awsFn func(context.Context, string, time.Duration) aws.ValidationResult,
	kubeFn func(context.Context, time.Duration) kube.ContextResult,
	terraformFn func(context.Context, time.Duration) terraform.AvailabilityResult,
	networkFn func(context.Context, string, time.Duration) network.ReachabilityResult,
) func() {
	prevAWS := validateAWSProfile
	prevKube := currentKubeContext
	prevTerraform := terraformAvailable
	prevNetwork := networkReachable

	validateAWSProfile = awsFn
	currentKubeContext = kubeFn
	terraformAvailable = terraformFn
	networkReachable = networkFn

	return func() {
		validateAWSProfile = prevAWS
		currentKubeContext = prevKube
		terraformAvailable = prevTerraform
		networkReachable = prevNetwork
	}
}
