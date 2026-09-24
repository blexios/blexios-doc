# Blexios DOC Generator

## Présentation

Blexios DOC est le premier composant du moteur de génération de documents de Blexios. Il permet de transformer des données JSON brutes en documents PDF structurés, avec une architecture modulaire prête à évoluer vers HTML, DOCX, XLSX et d'autres formats.

## Installation

```bash
go mod tidy
```

Ensuite, construire le binaire :

```bash
go build -o blexios .
```

## Utilisation

```bash
./blexios doc:generate --pdf --path=data.json
./blexios doc:generate --pdf --path=data.json --output=rapport.pdf
./blexios doc:generate --pdf --path=data.json --title="Rapport de présence"
./blexios doc:generate --pdf --path=data.json --orientation=landscape
```

## Format JSON supporté

Le moteur accepte surtout des tableaux d'objets JSON brut :

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

## Options CLI

- `--pdf` : active la génération PDF
- `--path` : chemin du fichier JSON source
- `--output` : chemin du fichier PDF cible
- `--title` : titre du document
- `--template` : template à appliquer (préparé pour l'avenir)
- `--orientation` : `auto`, `portrait` ou `landscape`
- `--verbose` : mode détaillé

## Exemples

### Exemple de présence

```bash
./blexios doc:generate --pdf --path=examples/doc/presences.json --title="Rapport de présence"
```

### Exemple de produits

```bash
./blexios doc:generate --pdf --path=examples/doc/products.json --title="Catalogue produits"
```

## Templates

La version actuelle fonctionne en mode automatique. Une logique de templates est prévue pour les futures versions : `default`, `report`, `invoice`, `attendance`, etc.

## Architecture

```text
cmd/
  doc/
    generate.go
internal/
  doc/
    parser/
    model/
    formatter/
    renderer/
    service/
```

Le flux logique est :

```text
JSON -> Parser -> Data analysis -> Document model -> Renderer -> PDF
```

## Développement futur

- support HTML, DOCX, XLSX
- templates avancés
- styles et thèmes
- i18n/locale
- logos et en-têtes personnalisés
- streaming pour gros volumes
