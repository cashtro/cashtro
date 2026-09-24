package agents

// NicheDeskSpec is one desk of a program's own brain.
// The shape matches the main brain. The work does not.
type NicheDeskSpec struct {
	Name  string
	Motto string
	Fn    string
}

// NicheBrain is a program's own brain. Epicenter Einstein stays the main brain.
// This one has the same five-desk shape and only does this program's functions.
type NicheBrain struct {
	ID                 string   `json:"id"`
	Niche              string   `json:"niche"`
	MainBrain          string   `json:"mainBrain"`
	OwnBrain           bool     `json:"ownBrain"`
	SameShape          bool     `json:"sameShape"`
	SameBrain          bool     `json:"sameBrain"`
	FunctionsOnly      bool     `json:"functionsOnly"`
	LoadsInstinctBrain bool     `json:"loadsInstinctBrain"`
	DecidedBy          string   `json:"decidedBy"`
	Desks              []Desk   `json:"desks"`
	Workers            int      `json:"workers"`
	Closed             []string `json:"closed"`
}

func buildNiche(id, niche string, specs []NicheDeskSpec, closed []string) NicheBrain {
	keys := []string{"pense", "voit", "fait", "tient", "borne"}
	counts := []int{200, 150, 200, 200, 150}
	directors := []int{1, 4, 2, 2, 4}
	parts := []string{"regard", "tri", "geste", "preuve", "pont"}
	desks := make([]Desk, len(specs))
	workers := 0
	for i, spec := range specs {
		key := keys[i]
		leads := make([]string, directors[i])
		for n := range leads {
			leads[n] = id + "-" + key + "-" + itoa(n+1)
		}
		subs := make([]SubBrain, len(parts))
		for n, part := range parts {
			role := spec.Fn + " Rien hors de ce créneau."
			if part == "pont" {
				role = spec.Fn + " Remet au cerveau principal. Epicenter Einstein décide."
			}
			subs[n] = SubBrain{ID: id + "-" + key + "-" + part, Name: part, Role: role, Bridges: part == "pont"}
		}
		desks[i] = Desk{
			ID: id + "-" + key, Name: spec.Name, Motto: spec.Motto,
			Directors: leads, Workers: counts[i], SubBrains: subs,
		}
		workers += counts[i]
	}
	return NicheBrain{
		ID: id, Niche: niche, MainBrain: "instinct",
		OwnBrain: true, SameShape: true, SameBrain: false, FunctionsOnly: true,
		LoadsInstinctBrain: false, DecidedBy: "instinct",
		Desks: desks, Workers: workers, Closed: closed,
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [4]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}

func d(name, motto, fn string) NicheDeskSpec {
	return NicheDeskSpec{Name: name, Motto: motto, Fn: fn}
}

// NicheBrains is every program's own brain. The control line is the main brain and is not in this list.
func NicheBrains() []NicheBrain {
	return []NicheBrain{
		buildNiche("wordpress", "thèmes WordPress Proximity", []NicheDeskSpec{
			d("Le Thème", "Pense le PHP et l'ACF.", "Penser le thème PHP et ACF Pro seulement."),
			d("La Liste", "Voit ces dépôts-là.", "Voir seulement les dépôts WordPress de ce département."),
			d("L'Atelier", "Fait le thème en local.", "Confectionner en local sur XAMPP."),
			d("La Relecture", "Tient le thème hors ligne.", "Relire le thème sans le mettre en ligne."),
			d("La Borne", "Refuse Lovable.", "Refuser Lovable et remettre le commit au cerveau principal."),
		}, []string{"lovable", "mint", "envoi"}),
		buildNiche("proximity", "comptes hors WordPress", []NicheDeskSpec{
			d("Le Compte", "Pense hors PHP.", "Penser les comptes qui ne sont pas des sites WordPress."),
			d("La Séparation", "Ne mélange pas.", "Voir ce département sans le mélanger au PHP."),
			d("L'Empilement", "Empile ces comptes.", "Empiler ces comptes, sans toucher ACF ni XAMPP."),
			d("La Revue", "Relit ces comptes.", "Tenir la revue de ces comptes seulement."),
			d("La Borne", "Refuse le commit PHP.", "Refuser un commit WordPress et remettre au cerveau principal."),
		}, []string{"acf", "xampp", "php"}),
		buildNiche("scanapp", "stock du scanner", []NicheDeskSpec{
			d("Le Stock", "Pense le scanner.", "Penser le stock du scanner interne."),
			d("Le Relevé", "Voit le stock interne.", "Voir le stock, pas la marketplace publique."),
			d("La Fiche", "Écrit le relevé.", "Tenir le relevé du scan."),
			d("L'Offre", "Mesure sans envoyer.", "Préparer l'offre mesurée du scan."),
			d("La Borne", "Remet l'envoi.", "Remettre l'envoi au cerveau principal."),
		}, []string{"marketplace", "encaissement"}),
		buildNiche("marketplace", "vente du Scan App", []NicheDeskSpec{
			d("La Fiche", "Pense l'item.", "Penser la marketplace du Scan App seulement."),
			d("Le Rayon", "Voit les items.", "Voir les items à vendre, pas le stock interne."),
			d("La Préparation", "Prépare la fiche.", "Préparer la fiche item."),
			d("La Mesure", "Mesure sans encaisser.", "Mesurer l'offre sans encaisser."),
			d("La Borne", "Remet la mise en vente.", "Remettre la mise en vente au cerveau principal."),
		}, []string{"stock interne", "encaissement"}),
		buildNiche("panda", "service white-glove", []NicheDeskSpec{
			d("Le Service", "Pense le white-glove.", "Penser le service white-glove seulement."),
			d("Le Besoin", "Voit le client Panda.", "Voir le besoin du client Panda."),
			d("Le Brouillon", "Prépare le geste.", "Préparer installation, cours ou IA en brouillon."),
			d("L'Offre", "Ne déploie pas.", "Tenir l'offre sans déployer chez le client."),
			d("La Borne", "Remet le déploiement.", "Remettre le déploiement au cerveau principal."),
		}, []string{"déploiement", "envoi"}),
		buildNiche("nft-giant", "art NFT", []NicheDeskSpec{
			d("La Pièce", "Pense l'art utile.", "Penser l'art NFT utile seulement."),
			d("Le Regard", "Voit la pièce.", "Voir la pièce, pas un ordre de marché."),
			d("Le Brouillon", "Écrit la pièce.", "Rédiger le brouillon de la pièce."),
			d("Le Coffre", "Ne minte pas.", "Tenir le brouillon sans minter."),
			d("La Borne", "Remet le mint.", "Remettre le mint au cerveau principal."),
		}, []string{"mint", "ordre"}),
		buildNiche("ecole", "cours tech", []NicheDeskSpec{
			d("Le Cours", "Pense la leçon.", "Penser le cours tech seulement."),
			d("Le Programme", "Voit le programme.", "Voir le programme, pas une promesse de gain."),
			d("La Leçon", "Écrit la leçon.", "Rédiger la leçon."),
			d("La Tenue", "Décrit le cours.", "Tenir le cours décrit."),
			d("La Borne", "Refuse le gain promis.", "Remettre une promesse de trading au cerveau principal."),
		}, []string{"trading", "gain"}),
		buildNiche("marketing", "copie de l'agence", []NicheDeskSpec{
			d("La Copie", "Pense l'agence.", "Penser la copie de l'agence seulement."),
			d("Le Carnet", "Voit les agents.", "Voir le carnet des agents, pas les clients Fix Tout."),
			d("Le Guide", "Tient les leçons.", "Tenir le guide de copie et ses leçons."),
			d("La Mesure", "Mesure sans envoyer.", "Mesurer l'offre sans l'envoyer."),
			d("La Borne", "Remet l'envoi.", "Remettre l'envoi au cerveau principal."),
		}, []string{"clients Fix Tout", "envoi"}),
		buildNiche("empire", "studio et live", []NicheDeskSpec{
			d("Le Studio", "Pense le live.", "Penser le studio et le live seulement."),
			d("Les Caméras", "Voit le plateau.", "Voir les caméras, le podcast et l'UGC."),
			d("La Vente", "Prépare le live.", "Préparer la vente live depuis le stock du scan."),
			d("La Sortie", "Ne encaisse pas ici.", "Tenir la sortie sans encaisser ici."),
			d("La Borne", "Remet l'encaissement.", "Remettre l'encaissement au cerveau principal."),
		}, []string{"encaissement"}),
		buildNiche("propres", "projets maison", []NicheDeskSpec{
			d("Le Projet", "Pense la maison.", "Penser les projets maison seulement."),
			d("La Carte", "Voit ces projets.", "Voir ces projets, pas les clients."),
			d("Le Geste", "Prépare le prochain pas.", "Préparer le prochain geste du projet."),
			d("Le Local", "Tient le brouillon.", "Tenir le brouillon local."),
			d("La Borne", "Remet la publication.", "Remettre la publication au cerveau principal."),
		}, []string{"publication", "clients"}),
		buildNiche("trading", "modèles de bots", []NicheDeskSpec{
			d("Le Modèle", "Pense le bot.", "Penser le modèle de bot seulement."),
			d("La Lecture", "Lit en public.", "Lire le marché en public, sans ordre."),
			d("La Note", "Note le modèle.", "Noter le modèle."),
			d("La Tenue", "Ne trade pas.", "Tenir la note sans trader."),
			d("La Borne", "Remet l'ordre.", "Remettre l'ordre au cerveau principal."),
		}, []string{"ordre", "mint"}),
		buildNiche("fix2", "chantiers Fix Tout", []NicheDeskSpec{
			d("Le Chantier", "Pense Fix Tout.", "Penser le chantier Fix Tout seulement."),
			d("Le Carnet", "Voit ce carnet-là.", "Voir le carnet Evolu-Jeunes/Fix2, pas l'autre CRM."),
			d("La Soumission", "Compte la taxe.", "Préparer la soumission et la taxe du Québec."),
			d("La Voix", "N'appelle pas.", "Tenir la voix Vapi sans appeler."),
			d("La Borne", "Remet l'envoi.", "Remettre l'envoi et l'encaissement au cerveau principal."),
		}, []string{"envoi", "encaissement", "autre CRM"}),
		buildNiche("fonds", "notes de fonds", []NicheDeskSpec{
			d("La Note", "Pense le fonds.", "Penser le fonds en note seulement."),
			d("Le Chiffre", "Voit sans virer.", "Voir le chiffre, sans mouvement bancaire."),
			d("Le Brouillon", "Prépare le contrôle.", "Préparer le brouillon du contrôle."),
			d("La Tenue", "Ne vire pas.", "Tenir la note sans virer."),
			d("La Borne", "Remet le virement.", "Remettre le virement au cerveau principal."),
		}, []string{"virement", "banque"}),
		buildNiche("ops", "opérations en brouillon", []NicheDeskSpec{
			d("L'Accueil", "Qualifie sans appeler.", "Qualifier le lead sans appeler."),
			d("Le Terrain", "Voit le chantier.", "Voir le chantier et la tournée en brouillon."),
			d("La Propriété", "Aide nos dépôts.", "Aider un dépôt cashtro ou Evolu-Jeunes seulement."),
			d("Le Devis", "N'envoie pas.", "Préparer la soumission sans l'envoyer."),
			d("La Borne", "Refuse le site.", "Refuser le site, l'envoi et le paiement, et remettre au cerveau principal."),
		}, []string{"site", "envoi", "paiement"}),
		buildNiche("hustle", "métaux légaux", []NicheDeskSpec{
			d("L'Offre", "Pense offre et mesure.", "Penser l'offre et la mesure seulement."),
			d("La Lecture", "Ne copie pas le livre.", "Lire la leçon sans copier le livre."),
			d("La Rue", "Juge, ne frappe pas.", "Tenir le métal de la rue comme jugement, pas comme coup."),
			d("La Société", "Garde le coup légal.", "Garder seulement le coup de la société, légal."),
			d("La Borne", "Refuse la fraude.", "Refuser la fraude et remettre le lancement au cerveau principal."),
		}, []string{"fraude", "lancement"}),
		buildNiche("eye", "veille défensive", []NicheDeskSpec{
			d("La Veille", "Liste nos dépôts.", "Lister nos dépôts seulement."),
			d("La Menace", "Nomme sans mode d'emploi.", "Nommer la faiblesse sans dire comment s'en servir."),
			d("La Faille", "Tient le brouillon.", "Tenir le correctif en brouillon."),
			d("La Garde", "Refuse l'attaque.", "Refuser l'attaque, le flood et la copie d'un secret."),
			d("Le Sceau", "Remet sans appliquer.", "Remettre le changement au cerveau principal sans l'appliquer."),
		}, []string{"attaque", "flood", "secret"}),
		buildNiche("forge", "développement", []NicheDeskSpec{
			d("L'Intake", "Prend le dépôt.", "Prendre chaque dépôt que Scrum amène."),
			d("Le Chantier", "Voit le développement.", "Voir le chantier de développement seulement."),
			d("Le Brouillon", "N'écrit pas une arme.", "Écrire le brouillon, pas une arme."),
			d("La Revue", "Ne déploie pas.", "Relire sans déployer."),
			d("La Borne", "Remet la décision.", "Remettre la décision au cerveau principal."),
		}, []string{"arme", "déploiement"}),
		buildNiche("giant", "crypto en formation", []NicheDeskSpec{
			d("Le Syllabus", "Tient la formation.", "Tenir le syllabus crypto."),
			d("La Lecture", "Lit sans ordre.", "Lire le marché public, sans ordre."),
			d("Le Brouillon", "Écrit le token en note.", "Rédiger le brouillon du token et de l'app."),
			d("Le Coffre", "Le mint n'est pas une option.", "Tenir la formation. Le mint n'est pas une option."),
			d("La Borne", "Remet le mint.", "Remettre mint, ordre, pont et déploiement au cerveau principal."),
		}, []string{"mint", "ordre", "pont", "déploiement"}),
		buildNiche("fusion", "application des lois", []NicheDeskSpec{
			d("La Règle", "Lit le checklist.", "Lire la loi sur Evolu-Jeunes et cashtro ensemble."),
			d("Les Deux", "Ne coupe pas le nom.", "Voir les deux peuples dans le même geste."),
			d("La Porte", "Arrête le raté.", "Arrêter le geste qui rate la règle, sans attaquer."),
			d("Le Registre", "Note sans secret.", "Tenir la preuve locale sans copier un secret."),
			d("La Borne", "Remet sans déployer.", "Remettre la décision au cerveau principal. Ne pas importer un cerveau TypeScript."),
		}, []string{"attaque", "secret", "surveillance", "import"}),
		buildNiche("pandora", "atelier Pandora", []NicheDeskSpec{
			d("L'Atelier", "Pense Pandora.", "Penser l'atelier Pandora seulement."),
			d("La Copie", "Tient le guide.", "Tenir le guide de copie en entier."),
			d("La Voix", "Ne réécrit pas Vapi.", "Décrire la voix Vapi sans la réécrire."),
			d("La Base", "Reste sur db:pandora.", "Rester sur db:pandora."),
			d("La Borne", "Remet la réécriture.", "Remettre l'appel, l'envoi et la réécriture au cerveau principal."),
		}, []string{"appel", "envoi", "réécriture"}),
	}
}

// NicheByID returns one program's own brain.
func NicheByID(id string) (NicheBrain, bool) {
	for _, brain := range NicheBrains() {
		if brain.ID == id {
			return brain, true
		}
	}
	return NicheBrain{}, false
}
