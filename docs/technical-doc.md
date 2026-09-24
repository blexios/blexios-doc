# Documentation technique – Blexios DOC

## 1. Vue d’ensemble

Blexios DOC est un composant CLI écrit en Go qui transforme un fichier JSON brut en document professionnel au format PDF.

Le but est de fournir un moteur modulaire capable de :

- lire des données JSON structurées ;
- détecter automatiquement la structure et les colonnes ;
- transformer les clés techniques en libellés lisibles ;
- formater les valeurs (dates, heures, booléens, nombres, null) ;
- générer un document PDF exploitable et professionnel.

Le système a été conçu pour évoluer vers d’autres formats comme HTML, DOCX, XLSX, etc., sans coupler directement le parser JSON au rendu PDF.

---

## 2. Commande principale

La commande de base est :

```powershell
blexios doc:generate --pdf --path=D:\data\presences.json
```

### Commande minimale

```powershell
blexios doc:generate --pdf --path=D:\data\presences.json
```

### Avec sortie personnalisée

```powershell
blexios doc:generate --pdf --path=D:\data\presences.json --output=D:\data\rapport.pdf
```

### Avec titre

```powershell
blexios doc:generate --pdf --path=D:\data\presences.json --title="Rapport de présence"
```

### Avec orientation forcée

```powershell
blexios doc:generate --pdf --path=D:\data\presences.json --orientation=landscape
```

### Avec aide

```powershell
blexios doc:generate --help
```

---

## 3. Comment la commande est interprétée

La commande est gérée dans le fichier :

- [cmd/doc/generate.go](../cmd/doc/generate.go)

Le point d’entrée est :

- [main.go](../main.go)

### Flux de lecture

1. `main.go` passe les arguments du terminal à la commande `doc:generate`.
2. `cmd/doc/generate.go` parse les flags :
   - `--pdf`
   - `--path`
   - `--output`
   - `--title`
   - `--template`
   - `--orientation`
   - `--verbose`
3. La commande valide les paramètres.
4. La génération est déléguée au service `Generator`.
5. Le service lit le fichier JSON, valide sa structure et génère le PDF.

---

## 4. Structure technique du projet

```text
blexios-doc/
├── cmd/
│   └── doc/
│       └── generate.go
├── internal/
│   └── doc/
│       ├── parser/
│       │   ├── json.go
│       │   └── detector.go
│       │
│       ├── model/
│       │   ├── document.go
│       │   ├── column.go
│       │   └── row.go
│       │
│       ├── formatter/
│       │   ├── date.go
│       │   ├── time.go
│       │   └── value.go
│       │
│       ├── renderer/
│       │   ├── renderer.go
│       │   └── pdf/
│       │       └── renderer.go
│       │
│       └── service/
│           └── generator.go
├── docs/
│   ├── doc-generator.md
│   └── technical-doc.md
├── examples/
│   └── doc/
│       ├── presences.json
│       ├── products.json
│       └── clients.json
├── main.go
├── go.mod
├── README.md
└── .gitignore
```

---

## 5. Les modules et leur rôle

### 5.1 Parser JSON

Fichier :

- [internal/doc/parser/json.go](../internal/doc/parser/json.go)

Rôle :

- lire le contenu brut du fichier ;
- décoder le JSON en structure Go ;
- vérifier que la racine est un tableau d’objets ;
- renvoyer une erreur claire si le JSON est vide ou invalide.

#### Exemple de logique

```go
rows, err := parser.ParseJSON(raw)
```

Le parseur accepte par exemple :

```json
[
  { "id": 1, "name": "Alice" },
  { "id": 2, "name": "Bob" }
]
```

Et rejette les structures non supportées comme :

```json
{ "id": 1, "name": "Alice" }
```

---

### 5.2 Détection de colonnes

Fichier :

- [internal/doc/parser/detector.go](../internal/doc/parser/detector.go)

Rôle :

- parcourir toutes les lignes du tableau ;
- collecter les clés ;
- déterminer les colonnes disponibles ;
- inférer le type de chaque colonne : `string`, `date`, `time`, `number`, `bool`.

Exemple de clés détectées :

```text
employee_id
date_presence
heure_arrivee
heure_depart
statut
justification
```

---

### 5.3 Humanisation des libellés

Le système transforme automatiquement les noms techniques en libellés lisibles.

Exemples :

- `employee_id` → `Employee ID`
- `date_presence` → `Date Presence`
- `heure_arrivee` → `Heure Arrivee`
- `firstName` → `First Name`

La transformation est assurée par la logique de `HumanizeKey` dans :

- [internal/doc/parser/detector.go](../internal/doc/parser/detector.go)

---

### 5.4 Modèle de document

Fichier :

- [internal/doc/model/document.go](../internal/doc/model/document.go)
- [internal/doc/model/column.go](../internal/doc/model/column.go)
- [internal/doc/model/row.go](../internal/doc/model/row.go)

Le modèle interne est indépendant du rendu PDF.

```go
type Document struct {
    Title   string
    Columns []Column
    Rows    []Row
    Metadata Metadata
}
```

Le document ne dépend pas d’une bibliothèque PDF. Il est simplement une structure métier prête à être rendue dans n’importe quel format.

---

### 5.5 Formatters

Fichiers :

- [internal/doc/formatter/date.go](../internal/doc/formatter/date.go)
- [internal/doc/formatter/time.go](../internal/doc/formatter/time.go)
- [internal/doc/formatter/value.go](../internal/doc/formatter/value.go)

Ils servent à afficher les valeurs de manière lisible :

- `null` → `—`
- `2026-09-21` → `21/09/2026`
- `19:37:00` → `19:37`
- `true` → `Oui`
- `false` → `Non`

La logique est centralisée dans :

```go
formatter.FormatValue(value, valueType)
```

---

### 5.6 Service de génération

Fichier :

- [internal/doc/service/generator.go](../internal/doc/service/generator.go)

C’est l’orchestrateur principal.

Il réalise la chaîne complète :

```text
validation du fichier
  ↓
lecture JSON
  ↓
validation du tableau
  ↓
détection des colonnes
  ↓
création du modèle Document
  ↓
choix de l’orientation
  ↓
création du fichier de sortie
  ↓
appel du renderer PDF
```

Il décide aussi :

- si le chemin de sortie est fourni ;
- sinon, il crée automatiquement un fichier `.pdf` au même endroit ;
- s’il existe déjà, il ajoute un suffixe `_1`, `_2`, etc.

---

### 5.7 Renderer PDF

Fichier :

- [internal/doc/renderer/pdf/renderer.go](../internal/doc/renderer/pdf/renderer.go)

Rôle :

- créer le document PDF ;
- définir les marges et la page A4 ;
- choisir l’orientation auto selon le nombre de colonnes ;
- tracer le titre ;
- dessiner la table ;
- gérer les lignes et les sauts de page ;
- écrire le fichier final.

Le renderer ne lit pas le JSON directement : il reçoit le modèle interne déjà construit.

---

## 6. Flux complet d’exécution

Voici le flux exact quand on tape :

```powershell
blexios doc:generate --pdf --path=D:\data\presences.json --title="Rapport de présence"
```

### Étape 1 – démarrage

`main.go` reçoit les arguments système.

### Étape 2 – traitement CLI

`cmd/doc/generate.go` lit :

- `--pdf` → active le format PDF
- `--path` → fichier source
- `--title` → titre du document

### Étape 3 – validation

Le service vérifie :

- le fichier existe ;
- le fichier est lisible ;
- le contenu JSON est valide ;
- la racine est un tableau d’objets ;
- le tableau contient au moins un enregistrement.

### Étape 4 – détection des données

Le parseur JSON remonte une liste de lignes :

```go
[]map[string]any
```

Puis le détecteur analyse les clés pour construire les colonnes.

### Étape 5 – construction du document

Le service crée un objet `Document` avec :

- `Title`
- `Columns`
- `Rows`
- `Metadata`

### Étape 6 – rendu PDF

Le PDF renderer reçoit le document et produit le fichier.

### Étape 7 – sortie terminal

Le CLI affiche alors :

```text
Blexios DOC Generator

Input:
D:/data/presences.json
Records: 17
Columns: 7

Generating PDF...
████████████████████████ 100%

PDF generated successfully.

Output:
D:/data/presences.pdf
```

---

## 7. Gestion des chemins

Le projet supporte les chemins Windows, Linux et macOS grâce à l’utilisation des fonctions standards du package Go :

- `filepath`
- `os.Stat`
- `os.ReadFile`
- `filepath.Ext`

Cela permet de gérer correctement :

- `D:\data\file.json`
- `D:/data/file.json`
- `/tmp/data.json`

---

## 8. Gestion des erreurs

Les erreurs sont volontairement formulées pour l’utilisateur final.

Exemples :

- fichier introuvable :

```text
Error: input file not found:
D:/data/presences.json
```

- JSON invalide :

```text
Error: invalid JSON.
Unable to parse:
D:/data/presences.json
```

- tableau vide :

```text
Error: the JSON array contains no records.
```

- structure non supportée :

```text
Error: unsupported JSON structure.
Expected an array of objects.
```

- génération PDF impossible :

```text
Error: unable to generate PDF.
```

---

## 9. Exemple réel

Fichier d’entrée :

```json
[
  {
    "id": "2138",
    "employee_id": "53",
    "date_presence": "2026-09-21",
    "heure_arrivee": "19:37:00",
    "heure_depart": null,
    "statut": "présent",
    "justification": null
  }
]
```

Le système détecte :

- 1 enregistrement
- 7 colonnes
- `employee_id` comme identifiant
- `date_presence` comme date
- `heure_arrivee` comme heure
- `heure_depart` comme null

Ensuite le document est généré avec :

- titre configurable ;
- table avec en-têtes ;
- lignes de données ;
- valeurs null remplacées par `—` ;
- PDF en sortie.

---

## 10. Points d’évolution

Le projet est prêt pour la prochaine étape :

- HTML renderer
- DOCX renderer
- XLSX renderer
- templates
- styles et thème
- internationalisation
- logos dans l’entête
- gestion de gros volumes via streaming

---

## 11. Commandes utiles

### Lancer les tests

```powershell
go test ./...
```

### Compiler le binaire

```powershell
go build -o blexios.exe .
```

### Générer un PDF

```powershell
./blexios.exe doc:generate --pdf --path=examples/doc/presences.json --title="Rapport de présence"
```

### Afficher l’aide

```powershell
./blexios.exe doc:generate --help
```

---

## 12. Résumé

Le moteur fonctionne comme une chaîne de transformation :

```text
JSON brut
  → parsing
  → détection automatique
  → modélisation interne
  → formatage des valeurs
  → rendu PDF
  → fichier généré
```

C’est cette architecture qui permet de faire évoluer le projet vers d’autres types de documents sans casser la logique existante.
