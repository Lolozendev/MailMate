# 🎨 Guide de Création de Templates

Créez vos propres modèles d'emails pour standardiser vos communications. Les templates utilisent une syntaxe simple basée sur HTML et le moteur de template Pongo2 (similaire à Jinja2/Django).

## Structure d'un Template

Un fichier template (`.html`) se compose de deux parties :
1. **L'En-tête (Frontmatter)** : Pour définir le sujet de l'email.
2. **Le Corps** : Le contenu HTML de l'email.

### Exemple Complet

```html
---
subject: Invitation pour {{ RecipientName }}
---
<html>
<body>
    <h1>Bonjour {{ RecipientName }},</h1>
    
    <p>Nous sommes ravis de vous inviter à {{ EventName | default:'notre événement exclusif' }}.</p>
    
    <p>Date : {{ EventDate | type:'date' }}</p>
    <p>Lieu : {{ Location }}</p>
    
    <p>Cordialement,<br>
    {{ SenderName }}</p>
</body>
</html>
```

## 📝 Syntaxe des Variables

Utilisez les doubles accolades `{{ }}` pour insérer des variables. Ces variables généreront automatiquement un formulaire interactif lors de l'exécution du programme.

- **Texte simple** : `{{ ClientName }}` -> Crée un champ texte.
- **Valeurs par défaut** : `{{ Company | default:'Ma Société' }}` -> Pré-remplit le champ. //NE MARCHE PAS POUR L'INSTANT

## 🛠️ Filtres Spéciaux

Nous avons ajouté des filtres spécifiques pour améliorer les formulaires de saisie :

| Filtre | Usage | Description |
|--------|-------|-------------|
| `type:'date'` | `{{ MyDate \| type:'date' }}` | Demande une date valide. |
| `type:'filepath'` | `{{ Report \| type:'filepath' }}` | Demande un chemin de fichier (utile pour validation). |
| `int` | `{{ Count \| int }}` | Assure que la valeur saisie est un nombre entier. |

## 💡 Astuces

- **Sujet Dynamique** : Vous pouvez utiliser des variables dans le sujet (voir l'exemple ci-dessus).
- **Nommage** : Donnez à vos fichiers des noms clairs (ex: `relance_client.html`, `compte_rendu.html`) car c'est ce qui apparaîtra dans le menu de sélection.

## 🧹 Astuce : Créer un Template depuis Outlook

Si vous avez déjà un email bien formaté dans Outlook, vous pouvez l'utiliser comme base pour votre template. Cependant, le code HTML généré par Outlook est souvent très verbeux et complexe.

Voici comment le nettoyer pour l'utiliser facilement :

1. Dans Outlook, ouvrez votre email et faites **Enregistrer sous** -> choisissez le format **HTML**.
2. Placez le fichier `.htm` ou `.html` dans le dossier où vous allez exécuter la commande.
3. Utilisez **Pandoc** via Docker pour nettoyer le fichier :

```bash
docker run --rm --volume ".:/data" pandoc/core "votre_export_outlook.htm" --from=html --to=html --strip-comments --syntax-highlighting=none -o "template_propre.html"
```

> **Note** : Remplacez `"votre_export_outlook.htm"` par le nom de votre fichier exporté.

> ⚠️ **Conseil d'expert** :
> Le moteur d'Outlook utilise des techniques anciennes (tableaux imbriqués) pour la mise en page. Le nettoyage par Pandoc modernise ce code.
> - **Si votre mail est simple** (texte, images, gras/italique) : Le résultat sera parfait.
> - **Si votre mise en page est complexe** (colonnes multiples, design précis) : Le nettoyage pourrait altérer l'alignement. Utilisez le résultat comme une **base propre** à retravailler, plutôt qu'un résultat final immédiat.
