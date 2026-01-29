# Guide : Ajouter un nouveau type de validation

Grâce à la centralisation de la logique de validation dans `validator.ApplyFilters()`, l'ajout d'un nouveau type ne nécessite des modifications que dans **UN SEUL fichier**.

## Exemple : Ajouter un type `email`

### Étape 1 : Ajouter la fonction de validation de base

Dans `internal/validator/validator.go`, ajoutez la fonction de validation :

```go
// ValidateEmail checks if the value is a valid email address.
func ValidateEmail(value string) error {
	if !strings.Contains(value, "@") || !strings.Contains(value, ".") {
		return fmt.Errorf("value %q is not a valid email address", value)
	}
	return nil
}
```

### Étape 2 : Ajouter le cas dans ApplyFilters

Dans la même fichier `internal/validator/validator.go`, ajoutez le cas dans la fonction `ApplyFilters` :

```go
func ApplyFilters(value string, filters []models.TemplateFilter) error {
	// ... code existant ...

	for _, f := range filters {
		switch f.Name {
		case "int":
			// ... code existant ...
		case "type":
			switch f.Arg {
			case "date":
				// ... code existant ...
			case "filepath":
				// ... code existant ...
			case "email":  // ← AJOUT ICI
				if err := ValidateEmail(value); err != nil {
					return fmt.Errorf("must be a valid email")
				}
			}
		}
	}
	return nil
}
```

### Étape 3 (Optionnel) : Ajouter un hint dans le TUI

Dans `internal/tui/form.go`, dans la fonction `getHint`, ajoutez :

```go
func getHint(filters []models.TemplateFilter) string {
	for _, f := range filters {
		if f.Name == "type" {
			switch f.Arg {
			case "date":
				return "DD-MM-YYYY"
			case "filepath":
				return "/path/to/file"
			case "email":  // ← AJOUT ICI
				return "user@example.com"
			}
		}
	}
	return ""
}
```

### C'est tout ! 🎉

Le nouveau type fonctionne automatiquement dans :
- ✅ Le formulaire TUI interactif
- ✅ Le mode CLI avec `--kv`
- ✅ Les deux utilisent la même validation

## Utilisation dans un template

```html
---
subject: Confirmation pour {{ Email | type:'email' }}
---
<html>
<body>
    <p>Email: {{ Email | type:'email' }}</p>
</body>
</html>
```

## Utilisation en CLI

```bash
./mailmate --template templates/contact.html --kv "Email='john@example.com'"
```

## Avant la centralisation ❌

Il fallait modifier **3 endroits** :
1. `internal/validator/validator.go` - Créer ValidateEmail()
2. `internal/tui/form.go` - Ajouter le switch case
3. `internal/kv/validator.go` - Ajouter le switch case (duplication!)

## Après la centralisation ✅

Il suffit de modifier **1 endroit** :
1. `internal/validator/validator.go` - Créer ValidateEmail() + ajouter dans ApplyFilters()

La fonction `ApplyFilters()` est automatiquement utilisée partout !
