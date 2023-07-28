# Package `util` :

The `util` package provides utility functions for working with big integers and string-to-big.Int conversions. The `ParseBigInt` function parses a string and returns a pointer to a big.Int if successful, while the `StringToUniqueBigInt` function uniquely converts a string to a big integer using ASCII character values. These utility functions are helpful for handling big integer arithmetic and unique representations of strings as big integers.

1. **ParseBigInt Function:**
   - `ParseBigInt` function is used to parse a string and return a pointer to a `big.Int` if successful.
   - It takes two parameters: `str` (the string to parse) and `param` (a string identifying the parameter being parsed).
   - Inside the function, a new big.Int variable `bigInt` is created.
   - `bigInt.SetString()` method is used to attempt parsing the input string `str` as a base-10 integer.
   - If the parsing is successful (i.e., the string can be converted to a big integer), it returns the pointer to the big integer.
   - If the parsing fails (e.g., the string contains non-numeric characters), it returns an error indicating the failure.

2. **StringToUniqueBigInt Function:**
   - `StringToUniqueBigInt` function converts a given string uniquely to a big.Int.
   - The function assumes that the given string consists of only ASCII characters.
   - Inside the function, a new big.Int variable `result` is created to store the final result.
   - A base of 256 is set (`base := big.NewInt(256)`), assuming an ASCII character set (8-bit).
   - The function iterates over each character `ch` in the input string `input`.
   - It performs the following steps for each character:
     - Shift left the current `result` by the base value (effectively multiplying it by 256).
     - Add the ASCII value of the character `ch` to the `result`.
   - The final result is a unique big integer representation of the input string.


