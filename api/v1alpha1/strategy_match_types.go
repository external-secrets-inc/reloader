package v1alpha1

type MatchStrategy struct {
	Path       string      `json:"path"`
	Conditions []Condition `json:"conditions"`
}

type Condition struct {
	Value     string             `json:"value"`
	Operation ConditionOperation `json:"operation"`
}

type ConditionOperation string

const (
	ConditionOperationEqual       ConditionOperation = "Equal"
	ConditionOperationNotEqual    ConditionOperation = "NotEqual"
	ConditionOperationContains    ConditionOperation = "Contains"
	ConditionOperationNotContains ConditionOperation = "NotContains"
	ConditionOperationIn          ConditionOperation = "RegularExpression"
)
