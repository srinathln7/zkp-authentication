package config

// ZKP- System parameters
// `p` and `q` are set to length of 164 bits
const (
	CPZKP_PARAM_P string = "42765216643065397982265462252423826320512529931694366715111734768493812630447"
	CPZKP_PARAM_Q string = "21382608321532698991132731126211913160256264965847183357555867384246906315223"
	CPZKP_PARAM_G string = "4"
	CPZKP_PARAM_H string = "9"
)

// Only for testing purposes
var (
	CPZKP_TEST_X_CORRECT   = "7"
	CPZKP_TEST_X_INCORRECT = "666"
)
