// Package evolujeunes is the delivery-org side of the Graphify tree.
// Private GitHub trees are not in this Cursor installation. Ships are
// catalog names. Live client production stays do-not-touch.
package evolujeunes

// Org is github.com/Evolu-Jeunes. Zero public repos on this token.
func Org() string { return "Evolu-Jeunes" }

// Member is the public org member github.com/evoluJeunes.
func Member() string { return "evoluJeunes" }

// Site is https://evolujeunes.ca. Live. Do not touch.
func Site() string { return "evolujeunes.ca" }

// PrivateRepos is the ~48-repo wall. Not scanned.
func PrivateRepos() string { return "private-repos-access-limited" }

// Fleet is every catalog ship under the org.
func Fleet() []string {
	return []string{
		BTKAvocats(),
		MDClinic(),
		SolutionHypothequeQC(),
		Educonnexion(),
		Proximity(),
		ScanApp(),
		CashtroDeliveryCatalog(),
	}
}

// BTKAvocats is a live law-firm platform. Do not touch.
func BTKAvocats() string { return "btk-avocats" }

// MDClinic is a live medical clinic. Do not touch.
func MDClinic() string { return "md-clinic" }

// SolutionHypothequeQC is a live mortgage service. Do not touch.
func SolutionHypothequeQC() string { return "solution-hypotheque-qc" }

// Educonnexion is a live education platform. Do not touch.
func Educonnexion() string { return "educonnexion" }

// Proximity is a live agency platform. Do not touch.
func Proximity() string { return "proximity" }

// ScanApp is concept. Safe to onboard after the private trees are visible.
func ScanApp() string { return "scanapp" }

// CashtroDeliveryCatalog is the idea-stage internal board.
func CashtroDeliveryCatalog() string { return "cashtro-delivery-catalog" }
