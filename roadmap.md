L'objectif initial devrait être : **prendre un relevé bancaire et produire un fichier propre, catégorisé, avec le minimum d'intervention humaine.**

## 🗺️ Roadmap

### Phase 0 — Comprendre le relevé bancaire

**Objectif :** connaître exactement le format des données entrantes.

**Livrable :** un petit fichier d'exemple représentatif.

---

# Phase 1 — Modèle de données

Créer les structures Go.

```go
type Transaction struct {
    Date        time.Time
    Description string
    Amount      float64

    Vendor   string
    Category string
    Status   Status
}

type Status string

const (
    Matched Status = "MATCHED"
    Unknown Status = "UNKNOWN"
)
```

À ce stade, aucune logique complexe.

👉 **Livrable :** package `model`.

---

# Phase 2 — Parser

```go
transactions, err := parser.Parse("statement.csv")
```

Tests importants :

```text
✓ date valide
✓ montant positif
✓ montant négatif
✓ description vide
✓ ligne invalide
✓ caractères accentués
✓ mauvais format de fichier
```

**Livrable :** parser fiable et testé.

---

# Phase 3 — Configuration des vendeurs

Créer ton fichier de configuration.

Par exemple :

```yaml
vendors:
  - name: Tim Hortons
    category: Meals
    patterns:
      - "TIM HORTONS"

  - name: Uber
    category: Transportation
    patterns:
      - "UBER *TRIP"
      - "UBER CANADA"

  - name: Uber Eats
    category: Meals
    patterns:
      - "UBER EATS"

  - name: Netflix
    category: Subscriptions
    patterns:
      - "NETFLIX"
```

Puis :

```go
config, err := config.Load("vendors.yaml")
```

**Livrable :** les vendeurs sont configurables **sans modifier le code Go**.

---

# Phase 4 — Normalisation

Avant de chercher un vendeur :

```text
" Tim Hortons #1234 "
          ↓
"TIM HORTONS #1234"
```

Créer quelque chose comme :

```go
func Normalize(description string) string
```

Tu peux gérer :

* majuscules/minuscules
* espaces multiples
* espaces en début/fin
* caractères spéciaux
* éventuellement accents

C'est une petite étape mais elle va énormément améliorer le matching.

---

# Phase 5 — Matching

C'est le cœur de ton application.

```text
Transaction
     │
     ▼
Normalize()
     │
     ▼
Matcher
     │
 ┌───┴────┐
 ▼        ▼
MATCH    UNKNOWN
```

Exemple :

```text
"TIM HORTONS #12345"
        ↓
pattern "TIM HORTONS"
        ↓
Vendor = Tim Hortons
Category = Meals
```

Résultat :

```go
Transaction{
    Description: "TIM HORTONS #12345",
    Amount:      -8.45,
    Vendor:      "Tim Hortons",
    Category:    "Meals",
    Status:      Matched,
}
```

Sinon :

```go
Status: Unknown
```

👉 **Livrable :** premier système automatique de catégorisation.

---

# Phase 6 — CLI

Maintenant tu rends ton outil utilisable.

Par exemple :

```bash
bookkeeping import statement.csv
```

Output :

```text
Reading statement.csv...

153 transactions found.

141 matched
12 unknown

Output:
  output/bookkeeping.csv
```

Puis :

```bash
bookkeeping unknown output/bookkeeping.csv
```

```text
Unknown transactions:

1. AMAZON MKTPLACE       -52.30
2. SQ *RESTAURANT        -34.20
3. DOLLARAMA              -18.75
```

👉 **Livrable :** outil utilisable depuis le terminal.

Et puisque tu aimes travailler avec le terminal, je pense que c'est un très bon choix pour la V1.

---

# Phase 7 — Export

Produire un CSV exploitable par ta sœur.

Par exemple :

```text
Date,Description,Vendor,Category,Amount,Status
2026-08-01,TIM HORTONS #123,Tim Hortons,Meals,-8.45,MATCHED
2026-08-02,UBER *TRIP,Uber,Transportation,-24.50,MATCHED
2026-08-03,AMAZON MKTPLACE,,, -52.30,UNKNOWN
```

Tu pourrais également produire un JSON pour debug :

```bash
bookkeeping import statement.csv --format json
```

---

# Phase 8 — Feedback loop

C'est ici que ton outil commence à devenir **vraiment pratique**.

Supposons :

```text
UNKNOWN

AMAZON MKTPLACE -52.30
```

Ta sœur choisit :

```text
Amazon → Office Supplies
```

Ton programme ajoute la règle :

```yaml
- name: Amazon
  category: Office Supplies
  patterns:
    - "AMAZON MKTPLACE"
```

Au prochain relevé :

```text
AMAZON MKTPLACE
       ↓
Amazon
       ↓
Office Supplies
```

Plus elle l'utilise, plus le nombre de `UNKNOWN` diminue.

---

# Phase 9 — Gestion des conflits

Il faut prévoir ce cas :

```text
UBER EATS
```

correspond potentiellement à :

```text
Uber
Uber Eats
```

Il faut donc introduire une priorité.

Par exemple :

```yaml
- name: Uber Eats
  priority: 100
  patterns:
    - "UBER EATS"

- name: Uber
  priority: 50
  patterns:
    - "UBER"
```

Le matcher prend la règle la plus spécifique/prioritaire.

Tu peux même commencer sans `priority` et utiliser simplement **le pattern le plus long**.

```text
UBER EATS     ← 9 caractères
UBER          ← 4 caractères

→ Uber Eats
```

---

# Phase 10 — Validation / audit

Pour de la comptabilité, je mettrais beaucoup d'efforts ici.

Chaque transaction devrait être traçable :

```text
Transaction
    │
    ├── Original description
    ├── Matched vendor
    ├── Category
    ├── Matching rule
    └── Status
```

Par exemple :

```json
{
  "description": "TIM HORTONS #12345",
  "vendor": "Tim Hortons",
  "category": "Meals",
  "status": "MATCHED",
  "matchedPattern": "TIM HORTONS"
}
```

Ça permet de répondre à :

> Pourquoi le programme a mis cette dépense dans cette catégorie ?

---

# Phase 11 — Interface graphique

**Seulement après** avoir validé le workflow CLI.


## 🎯 MVP

Je limiterais volontairement la **V1** à seulement :

```text
1. Lire CSV
2. Parser transactions
3. Normaliser descriptions
4. Charger vendors.yaml
5. Matcher vendeurs
6. Catégoriser
7. Identifier UNKNOWN
8. Exporter CSV
```

Donc :

```text
CSV
 ↓
Parser
 ↓
Normalize
 ↓
Matcher
 ↓
Category
 ↓
CSV
```

**Pas de DB, pas d'IA, pas de web app, pas de Docker, pas d'authentification.**

Une fois que ta sœur l'aura utilisé sur **2–3 vrais relevés**, tu verras précisément où sont les problèmes réels. C'est à ce moment-là que ça vaudra la peine d'ajouter les fonctionnalités plus avancées.
