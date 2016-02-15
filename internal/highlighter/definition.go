package highlighter

import "strings"

func Definition(def string) string {
	def = strings.Replace(def, "package main", yellow+"package "+blue+"main"+reset, -1)
	def = strings.Replace(def, "import ", yellow+"import"+reset+" ", -1)
	def = strings.Replace(def, "func ", blue+"func"+reset+" ", -1)
	def = strings.Replace(def, "func(", blue+"func"+reset+"(", -1)
	def = strings.Replace(def, "defer ", blue+"defer"+reset+" ", -1)
	def = strings.Replace(def, " error ", " "+red+"error"+reset+" ", -1)

	def = strings.Replace(def, "\nFeature(", red+"\nFeature"+reset+"(", -1)
	def = strings.Replace(def, "\nScenario(", red+"\nScenario"+reset+"(", -1)
	def = strings.Replace(def, "\nGiven(", green+"\nGiven"+reset+"(", -1)
	def = strings.Replace(def, "\nAnd(", green+"\nAnd"+reset+"(", -1)
	def = strings.Replace(def, "\nWhen(", blue+"\nWhen"+reset+"(", -1)
	def = strings.Replace(def, "\nThen(", yellow+"\nThen"+reset+"(", -1)

	def = strings.Replace(def, "main()", yellow+"main"+reset+"()", -1)
	def = strings.Replace(def, "setup()", yellow+"setup"+reset+"()", -1)

	def = strings.Replace(def, " Pending(", " "+yellow+"Pending"+reset+"(", -1)
	def = strings.Replace(def, "os.Open(", yellow+"os.Open"+reset+"(", -1)
	def = strings.Replace(def, "fd.Close(", yellow+"fd.Close"+reset+"(", -1)
	def = strings.Replace(def, "suite.Test(", yellow+"suite.Test"+reset+"(", -1)
	def = strings.Replace(def, "ParseBool(", yellow+"ParseBool"+reset+"(", -1)
	def = strings.Replace(def, "stdres.EnableColor(", yellow+"stdres.EnableColor"+reset+"(", -1)
	def = strings.Replace(def, "stdres.DisableColor(", yellow+"stdres.DisableColor"+reset+"(", -1)
	def = strings.Replace(def, "os.Args[", yellow+"os.Args"+reset+"[", -1)

	return def
}
