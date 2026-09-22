/**
 * Visual plan — corporation as-run after step-up.
 * Five brains + ten lines + L'Inquisiteur. Divorced agent rules,
 * per-department self-improve, contradict/optimize department.
 */
export const title = "Agentic OS · corporation step-up";

export const mermaid = `
flowchart TB
  subgraph corp["Cashtro OS 0.5.0 · Epicenter"]
    manager["Manager"]
    inquisitor["L'Inquisiteur · live<br/>contredit · teste · bloque · optimise"]
    manager --> inquisitor
    architecte["Architecte"]
    cartographe["Cartographe"]
    forgeron["Forgeron"]
    orfevre["Orfèvre"]
    hustler["Hustler"]
    manager --> architecte
    manager --> cartographe
    manager --> forgeron
    manager --> orfevre
    manager --> hustler
    inquisitor -.->|audit / smarter / block| architecte
    inquisitor -.->|audit / smarter / block| cartographe
    inquisitor -.->|audit / smarter / block| forgeron
    inquisitor -.->|audit / smarter / block| orfevre
    inquisitor -.->|audit / smarter / block| hustler
  end
  subgraph lines["Lignes · chacune a un smarter loop"]
    proximity["Proximity"]
    scanapp["Scan App"]
    panda["Panda"]
    nft["NFT + Giant"]
    ecole["École"]
    marketing["Marketing"]
    empire["Empire"]
    trading["Trading"]
    pandora["Pandora"]
    teal["Teal / Vapi"]
  end
  forgeron --> proximity
  forgeron --> empire
  hustler --> scanapp
  hustler --> panda
  hustler --> nft
  hustler --> marketing
  hustler --> pandora
  hustler --> teal
  cartographe --> ecole
  architecte --> trading
  inquisitor --> proximity
  inquisitor --> scanapp
  inquisitor --> panda
  inquisitor --> nft
  inquisitor --> ecole
  inquisitor --> marketing
  inquisitor --> empire
  inquisitor --> trading
  inquisitor --> pandora
  inquisitor --> teal
  subgraph rules["Divorced rules + ask-gate"]
    spec["Spec.Rules[] par agentic"]
    ask["Ask queue · questions à Castro avant mutation"]
    block["vapi.call refuse si la ligne est bloquée"]
  end
  inquisitor --> spec
  inquisitor --> ask
  teal --> block
`;

export default function AgenticOSPlan() {
  return (
    <section data-canvas="agentic-os-plan">
      <h1>Corporation · step-up</h1>
      <p>
        Cinq cerveaux, dix lignes, L&apos;Inquisiteur live. Chaque
        agentic porte ses propres rules. Chaque département a un
        smarter loop. Le contradicteur teste, bloque jusqu&apos;à
        auto-amélioration, et pousse l&apos;option la plus optimisée.
        Hustler ne dump plus six lignes sans boucle spécialisée.
      </p>
      <pre>{mermaid}</pre>
    </section>
  );
}
