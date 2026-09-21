package agents

import "github.com/cashtro/cashtro/internal/kernel"

// seedChooser parks the 2026-09-21 Samsung/Gmail scan and the jobs Castro
// can pick himself. Bodies, passwords, and student records stay off the desk.
func seedChooser(k *kernel.Kernel) {
	k.RecordScan(kernel.Scan{
		Account: "alejandro@proximityagency.ca",
		Device:  "Samsung / Gmail",
		Inbox:   620,
		Unread:  163,
		Starred: 42,
		Note:    "Connected Gmail is the mailbox the Samsung Gmail app shows for this account. Other Samsung mail apps and brandingeniousceo@gmail.com are not bound here. Linear, Todoist, and Asana still need Connect. Outbound stays gated.",
	})

	jobs := []kernel.Choice{
		{Key: "scan-gmail-2026-09-21", Kind: kernel.KindScan, Urgency: "info", Title: "Gmail scan · Samsung", Client: "Cashtro", Sector: "ops", Source: "gmail",
			Summary: "620 inbox · 163 unread · 42 starred on alejandro@proximityagency.ca. Newsletters (DataCamp, OpenAI, Shakepay) dominate unread. Take this to park a triage mandate."},
		{Key: "mail-clic-v2", Kind: kernel.KindMail, Urgency: "now", Title: "Clic Inspection v2", Client: "Proximity", Sector: "web", Source: "gmail", Stack: []string{"WordPress"},
			Summary: "Xavier sent the v2 zip. Fouad could not open Drive, then reported the staging admin up. Pick this to QA the deploy and close the loop."},
		{Key: "mail-alaska-v2", Kind: kernel.KindMail, Urgency: "now", Title: "Alaska v2 — Excel blocker", Client: "Alaska", Sector: "ecommerce", Source: "gmail", Stack: []string{"WordPress", "WooCommerce"},
			Summary: "Xavier finished items 1, 3, 4. Item 2 is blocked on Excel access. Souad's product-display notes still sit on the thread."},
		{Key: "mail-astro-poulet", Kind: kernel.KindMail, Urgency: "now", Title: "Astro Poulet — franchise + EN menu", Client: "Astro Poulet", Sector: "food", Source: "gmail", Stack: []string{"WordPress"},
			Summary: "Franchise zones for the form, Locations page, menu rebuild. FR menu is in. Mohamed approved the English pass."},
		{Key: "mail-proximity-catalog", Kind: kernel.KindMail, Urgency: "soon", Title: "Proximity service catalog / admin", Client: "Proximity", Sector: "agency", Source: "gmail", Stack: []string{"WordPress", "ACF"},
			Summary: "Priority admin adjustments after the full service catalogue was integrated."},
		{Key: "mail-proxisoins", Kind: kernel.KindMail, Urgency: "soon", Title: "ProxiSoins — confirm then English", Client: "ProxiSoins à Domicile", Sector: "health", Source: "gmail", Stack: []string{"WordPress"},
			Summary: "FR copy changes are done. Waiting on client confirm before the English version."},
		{Key: "mail-asselin-backup", Kind: kernel.KindMail, Urgency: "soon", Title: "Asselin & Greene backup", Client: "Asselin & Greene", Sector: "finance", Source: "gmail", Stack: []string{"WordPress"},
			Summary: "Fouad sent the latest live-site backup. Review and keep it off the public journal."},
		{Key: "mail-sos-besoin", Kind: kernel.KindMail, Urgency: "soon", Title: "SOS-Besoin Android test", Client: "SOS-Besoin", Sector: "mobile", Source: "gmail", Stack: []string{"Android"},
			Summary: "Tester invite is unread. Run it from the Samsung, not from this VM."},
		{Key: "mail-supabase-paused", Kind: kernel.KindMail, Urgency: "now", Title: "Supabase ProximityAPI paused", Client: "Proximity", Sector: "infra", Source: "gmail", Stack: []string{"Supabase"},
			Summary: "Free-tier project Api was paused after inactivity. Restore it or let it sleep — you pick."},
		{Key: "mail-azure-monitor", Kind: kernel.KindMail, Urgency: "soon", Title: "Azure Monitor — verify email", Client: "Proximity", Sector: "infra", Source: "gmail", Stack: []string{"Azure"},
			Summary: "Action-required verify for an Azure Monitor action group."},
		{Key: "mail-reco-letter", Kind: kernel.KindMail, Urgency: "soon", Title: "Recommendation letter follow-up", Client: "Evolu-Jeunes", Sector: "ops", Source: "gmail",
			Summary: "Starred forward of a recommendation letter. Decide whether to answer or park."},
		{Key: "mail-flag-football", Kind: kernel.KindMail, Urgency: "later", Title: "Flag football assistant coach", Client: "CSSPI", Sector: "ops", Source: "gmail",
			Summary: "Starred follow-up on the adjoint entraîneur post. Personal, not a ship unless you take it."},
		{Key: "mail-unread-noise", Kind: kernel.KindMail, Urgency: "later", Title: "Unread triage · 163", Client: "Inbox", Sector: "ops", Source: "gmail",
			Summary: "Unread is mostly DataCamp, OpenAI product mail, and Shakepay. Take this to park a cleanup mandate; I will not mass-archive without you."},

		{Key: "verb-scan-mail", Kind: kernel.KindVerb, Urgency: "info", Title: "Scan Gmail again", Source: "chooser",
			Summary: "Re-read the connected inbox. I cannot open Samsung device mail apps directly."},
		{Key: "verb-draft-reply", Kind: kernel.KindVerb, Urgency: "info", Title: "Draft a reply (gated)", Source: "comms",
			Summary: "comms.send parks a confirm. Nothing leaves until you Allow at the human gate."},
		{Key: "verb-close-desk", Kind: kernel.KindVerb, Urgency: "info", Title: "Close desk / pulse / open", Source: "watch",
			Summary: "Closed hours keep research, delivery, memory, and planner moving. The banner must say CLOSED HOURS when closed and hide when open."},
		{Key: "verb-park-mandate", Kind: kernel.KindVerb, Urgency: "info", Title: "Park a mandate on idea", Source: "delivery",
			Summary: "delivery.create / planner.backlog. You still advance idea → concept → production."},
		{Key: "verb-ingest-research", Kind: kernel.KindVerb, Urgency: "info", Title: "Ingest a research note", Source: "research",
			Summary: "Sourced findings live in the library, not in chat."},
		{Key: "verb-remember", Kind: kernel.KindVerb, Urgency: "info", Title: "Remember / recall", Source: "memory",
			Summary: "Episodic facts without a model key."},
		{Key: "verb-search-desk", Kind: kernel.KindVerb, Urgency: "info", Title: "Search the desk", Source: "explorer",
			Summary: "explorer.search across processes, ships, and notes."},
		{Key: "verb-review-video", Kind: kernel.KindVerb, Urgency: "info", Title: "Review a walkthrough video", Source: "reviewer",
			Summary: "Resident QA. The last closed-hours demo video had a broken banner; this desk now sets CLOSED HOURS text from watch.closed."},
		{Key: "verb-operator", Kind: kernel.KindVerb, Urgency: "info", Title: "Drive the browser (resident)", Source: "operator",
			Summary: "operator.browse is on the desk but waits for a worker bind."},
		{Key: "verb-deploy", Kind: kernel.KindVerb, Urgency: "info", Title: "Deploy / promote (resident)", Source: "deploy",
			Summary: "CI, preview, production. Gated during closed hours."},
		{Key: "verb-security", Kind: kernel.KindVerb, Urgency: "info", Title: "Triage CVE / SAST (resident)", Source: "security",
			Summary: "security.triage before a ship moves right."},
		{Key: "verb-linear", Kind: kernel.KindVerb, Urgency: "info", Title: "Connect Linear / Todoist / Asana", Source: "chooser",
			Summary: "Task apps are installed but not authenticated. Connect one if you want this desk to pull those boards."},
	}
	for _, job := range jobs {
		k.Offer(job)
	}
}
