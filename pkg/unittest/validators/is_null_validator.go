package validators

import (
	"github.com/helm-unittest/helm-unittest/internal/common"
	"github.com/helm-unittest/helm-unittest/pkg/unittest/valueutils"
	log "github.com/sirupsen/logrus"
)

// IsNullValidator validate value of Path id kind
type IsNullValidator struct {
	Path string
}

func (v IsNullValidator) failInfo(actual interface{}, index int, not bool) []string {
	actualYAML := common.TrustedMarshalYAML(actual)

	log.WithField("validator", "is_null").Debugln("actual content:", actualYAML)

	return splitInfof(
		setFailFormat(not, true, false, false, " to be null, got"),
		index,
		v.Path,
		actualYAML,
	)
}

// Validate implement Validatable
func (v IsNullValidator) Validate(context *ValidateContext) (bool, []string) {
	manifests, err := context.getManifests()
	if err != nil {
		return false, splitInfof(errorFormat, -1, err.Error())
	}

	validateSuccess := false
	validateErrors := make([]string, 0)

	for idx, manifest := range manifests {
		actual, err := valueutils.GetValueOfSetPath(manifest, v.Path)
		if err != nil {
			validateSuccess = false
			errorMessage := splitInfof(errorFormat, idx, err.Error())
			validateErrors = append(validateErrors, errorMessage...)
			continue
		}

		var singleActual interface{}
		if len(actual) > 0 {
			singleActual = actual[0]
		} else {
			singleActual = nil
		}

		if singleActual == nil == context.Negative {
			validateSuccess = false
			errorMessage := v.failInfo(singleActual, idx, context.Negative)
			validateErrors = append(validateErrors, errorMessage...)
			continue
		}

		validateSuccess = determineSuccess(idx, validateSuccess, true)
	}

	return validateSuccess, validateErrors
}
