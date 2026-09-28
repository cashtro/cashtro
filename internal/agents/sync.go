package agents

// SyncChain connects two repos without putting them in the same list.
// Giant stays an environment. The trading bot stays its own repo.
// Each one keeps its own agentic chain. They move together from here.
type SyncChain struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Giant      string `json:"giant"`
	GiantRepo  string `json:"giantRepo"`
	GiantChain string `json:"giantChain"`
	Bot        string `json:"bot"`
	BotRepo    string `json:"botRepo"`
	BotChain   string `json:"botChain"`
	Mixed      bool   `json:"mixed"`
	Together   bool   `json:"together"`
	Trade      string `json:"trade"`
	LiveOrder  bool   `json:"liveOrder"`
}

// GiantSync is the chain between Giant and the trading bot.
// The two GitHub repos stay apart. AMF approval is handled outside this kernel.
func GiantSync() SyncChain {
	return SyncChain{
		ID:         "giant-trading",
		Name:       "Synchronie Giant",
		Giant:      "écosystème et environnement",
		GiantRepo:  "Evolu-Jeunes/Giant",
		GiantChain: "nft-giant",
		Bot:        "bot de trading",
		BotRepo:    "cashtro/trading_bot-main",
		BotChain:   "trading",
		Mixed:      false,
		Together:   true,
		Trade:      "blockchain",
		LiveOrder:  false,
	}
}
