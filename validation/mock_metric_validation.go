package validation

import (
	"fmt"
	"reflect"

	"google.golang.org/genproto/googleapis/monitoring/v3"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

const (
	missingFieldMsg = "%v must contain requried %v field"
)

func IsValidRequest(req interface{}) error {
	reqReflect := reflect.ValueOf(req)
	requiredFields := []string{"Name"}
	requestName := ""

	switch req.(type) {
	case *monitoring.CreateMetricDescriptorRequest:
		requiredFields = append(requiredFields, "MetricDescriptor")
		requestName = "CreateMetricDescriptor"
	case *monitoring.GetMetricDescriptorRequest:
		requestName = "GetMetricDescriptorRequest"
	case *monitoring.DeleteMetricDescriptorRequest:
		requestName = "DeleteMetricDescriptorRequest"
	case *monitoring.GetMonitoredResourceDescriptorRequest:
		requestName = "GetMonitoredResourceDescriptorRequest"
	case *monitoring.ListMetricDescriptorsRequest:
		requestName = "ListMetricDescriptorsRequest"
	}

	return CheckForRequiredFields(requiredFields, reqReflect, requestName)
}

func CheckForRequiredFields(requiredFields []string, reqReflect reflect.Value, requestName string) error {
	br := &errdetails.BadRequest{}

	for _, field := range requiredFields {
		if isZero := reflect.Indirect(reqReflect).FieldByName(field).IsZero(); isZero {
			v := &errdetails.BadRequest_FieldViolation{
				Field:       field,
				Description: fmt.Sprintf(missingFieldMsg, requestName, field),
			}
			br.FieldViolations = append(br.FieldViolations, v)
		}
	}

	if len(br.FieldViolations) > 0 {
		st, err := MissingFieldError.WithDetails(br)

		if err != nil {
			panic(fmt.Sprintf("unexpected error attaching metadata: %v", err))
		}

		return st.Err()
	}
	
	return nil
}