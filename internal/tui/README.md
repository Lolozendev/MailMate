# Guide d'extension du TUI (Terminal User Interface)

Ce document explique comment ajouter de nouveaux types de variables ou filtres pour qu'ils soient correctement gérés par l'interface en ligne de commande (TUI) lors de la saisie des variables du modèle.

## Vue d'ensemble

L'application utilise deux étapes principales pour gérer les variables :
1.  **Saisie et Validation (TUI)** : Géré dans `internal/tui/form.go`. C'est ici que l'utilisateur entre les valeurs et que la validation interactive se produit.
2.  **Rendu (Templates)** : Géré dans `internal/templates/render.go`. C'est ici que `pongo2` traite les filtres lors de la génération du HTML final.

Pour ajouter un nouveau type de validation (par exemple `email` ou `color`), vous devez modifier ces deux fichiers.

## Étape 1 : Enregistrement du filtre (Templates)

Fichier : `internal/templates/render.go`

Même si la validation principale se fait dans le TUI, le moteur de template doit connaître le filtre pour ne pas planter lors du parsing.

### Cas A : Ajouter un nouveau "type" (ex: `| type:"email"`)

Modifiez la fonction `filterType` dans `internal/templates/render.go` :

```go
func filterType(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
    // ...
    switch typ {
    case "date":
        // ...
    case "email": // <--- AJOUTER CECI
        // Validation simple ou pass-through
        return pongo2.AsValue(val), nil
    // ...
    }
}
```

### Cas B : Ajouter un nouveau filtre nommé (ex: `| uppercase`)

1.  Enregistrez le filtre dans la fonction `init()` :
    ```go
    pongo2.RegisterFilter("uppercase", filterUppercase)
    ```
2.  Implémentez la fonction du filtre :
    ```go
    func filterUppercase(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
        return pongo2.AsValue(strings.ToUpper(in.String())), nil
    }
    ```

## Étape 2 : Validation dans le formulaire (TUI)

Fichier : `internal/tui/form.go`

C'est ici que l'expérience utilisateur est définie.

### 1. Ajouter la validation

Modifiez la fonction `createValidator` pour inclure la logique de validation de votre nouveau type.

```go
func createValidator(filters []app.TemplateFilter) func(string) error {
    return func(str string) error {
        // ...
        for _, f := range filters {
            switch f.Name {
            case "type":
                if f.Arg == "email" { // <--- AJOUTER CECI
                    if !strings.Contains(str, "@") {
                        return fmt.Errorf("doit être une adresse email valide")
                    }
                }
            case "uppercase": // Si vous avez ajouté un filtre nommé
                 // Validation spécifique si nécessaire
            }
        }
        return nil
    }
}
```

### 2. Ajouter un indice visuel (Optionnel)

Si vous souhaitez afficher un texte de remplacement (placeholder) dans le champ de saisie, modifiez `getHint` :

```go
func getHint(filters []app.TemplateFilter) string {
    for _, f := range filters {
        if f.Name == "type" && f.Arg == "email" {
            return "exemple@domaine.com"
        }
    }
    return ""
}
```

## Résumé

Pour ajouter `{{ Variable | type:"monType" }}` :

1.  **`internal/templates/render.go`** : Ajoutez `case "monType":` dans `filterType`.
2.  **`internal/tui/form.go`** : Ajoutez la logique de validation dans `createValidator`.
3.  **`internal/tui/form.go`** : (Optionnel) Ajoutez un placeholder dans `getHint`.
