package plan

import (
	"time"

	"github.com/vatebur/dbinstall/internal/host"
)

type Risk string

const (
	RiskInfo   Risk = "info"
	RiskLow    Risk = "low"
	RiskMedium Risk = "medium"
	RiskHigh   Risk = "high"
)

type Privilege string

const (
	PrivilegeUser Privilege = "user"
	PrivilegeRoot Privilege = "root"
)

type Warning struct {
	Code    string `json:"code"`
	Risk    Risk   `json:"risk"`
	Message string `json:"message"`
}

type Resource struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type Step struct {
	ID              string     `json:"id"`
	Description     string     `json:"description"`
	DependsOn       []string   `json:"dependsOn,omitempty"`
	Risk            Risk       `json:"risk"`
	Privilege       Privilege  `json:"privilege"`
	Implemented     bool       `json:"implemented"`
	RequiresRestart bool       `json:"requiresRestart,omitempty"`
	OwnedResources  []Resource `json:"ownedResources,omitempty"`
	Timeout         string     `json:"timeout"`
}

type Document struct {
	APIVersion   string    `json:"apiVersion"`
	OperationID  string    `json:"operationID"`
	GeneratedAt  time.Time `json:"generatedAt"`
	InstanceName string    `json:"instanceName"`
	SpecDigest   string    `json:"specDigest"`
	Provider     string    `json:"provider"`
	Version      string    `json:"version"`
	Host         host.Info `json:"host"`
	Executable   bool      `json:"executable"`
	Warnings     []Warning `json:"warnings,omitempty"`
	Steps        []Step    `json:"steps"`
}
