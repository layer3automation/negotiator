package utils

import "github.com/c-robinson/iplib"

func NetsToStrings(nets []iplib.Net4) []string {
	strings := make([]string, len(nets))
	for i, n := range nets {
		strings[i] = n.String()
	}

	return strings
}

func NetToString(n iplib.Net4) string {
	return n.String()
}

func StringsToNets(strings []string) []iplib.Net4 {
	result := make([]iplib.Net4, len(strings))

	for i, s := range strings {
		result[i] = iplib.Net4FromStr(s)
	}

	return result
}

func StringToNet(s string) iplib.Net4 {
	return iplib.Net4FromStr(s)
}
