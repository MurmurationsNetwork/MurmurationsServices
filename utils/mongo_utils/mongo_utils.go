package mongo_utils

import (
	"strings"

	"github.com/MurmurationsNetwork/MurmurationsServices/utils/rest_errors"
)

func ParseError(err error) rest_errors.RestErr {
	errString := err.Error()
	if strings.Contains(errString, "document is nil") {
		return rest_errors.NewBadRequestError("Cannot find the document.")
	}
	return rest_errors.NewBadRequestError(err.Error())
}
