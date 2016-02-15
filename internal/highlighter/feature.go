package highlighter

import "strings"

func Feature(def string) string {
	def = strings.Replace(def, "Feature: ", red+"Feature: "+reset, -1)
	def = strings.Replace(def, "Scenario: ", red+"Scenario: "+reset, -1)
	def = strings.Replace(def, " Given ", green+" Given "+reset, -1)
	def = strings.Replace(def, " And ", green+" And "+reset, -1)
	def = strings.Replace(def, " When ", blue+" When "+reset, -1)
	def = strings.Replace(def, " Then ", yellow+" Then "+reset, -1)

	return def
}
