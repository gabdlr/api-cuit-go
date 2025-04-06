package cuit

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gabdlr/api-cuit-go/utils"
)

const CUIT_REGEX = `^([\d]{2}-[\d]{8}-[\d]{1}|[\d]{11})$`

var CUIT_TYPES = map[uint8]bool{
	30: true,
	33: true,
	34: true,
}

func IsValid(cuit string) (isValidCuit bool) {
	if validateFormat(cuit) {
		cuit = standardizeCuit(cuit)
		if validateCuitType(cuit) {
			isValidCuit = validateWithVerifierDigit(cuit)
		}
	}
	return isValidCuit
}

func standardizeCuit(cuit string) string {
	if len(cuit) > 11 {
		cuit = strings.ReplaceAll(cuit, "-", "")
	}
	return cuit
}

func validateWithVerifierDigit(cuit string) bool {
	verificationResult := false
	toVerify := utils.ReverseStringWithBuffer(cuit[:10])

	weightUpResult := 0
	weightUpFactorCounter := -1
	weightUpCheckFactor := []int{2, 3, 4, 5, 6, 7}

	verifierDigit, err := strconv.Atoi(cuit[len(cuit)-1:])
	if err != nil {
		return false
	}

	for i := range 10 {
		if i%6 == 0 {
			weightUpFactorCounter += 1
		}
		weightUp, err := strconv.Atoi(string(toVerify[i]))
		if err != nil {
			return verificationResult
		}
		weightUpResult += weightUp * weightUpCheckFactor[i-6*weightUpFactorCounter]
	}
	mod11WeightupResult := weightUpResult % 11

	switch mod11WeightupResult {
	case 11:
		verificationResult = verifierDigit == 0
	case 10:
		verificationResult = verifierDigit == 9
	default:
		verificationResult = verifierDigit == 11-mod11WeightupResult
	}

	return verificationResult
}

func validateCuitType(cuit string) bool {
	validationResult := false
	cuitType, err := strconv.Atoi(cuit[:2])
	if err == nil {
		validationResult = CUIT_TYPES[uint8(cuitType)]
	}
	return validationResult
}

func validateFormat(cuit string) bool {
	regexExp, err := regexp.Compile(CUIT_REGEX)
	if err == nil {
		return regexExp.MatchString(cuit)
	}
	return false
}
