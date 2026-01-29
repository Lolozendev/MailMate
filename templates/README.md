# 🎨 Guide de Création de Templates

Créez vos propres modèles d'emails pour standardiser vos communications. Les templates utilisent une syntaxe simple basée sur HTML et le moteur de template Pongo2 (similaire à Jinja2/Django).

## Structure d'un Template

Un fichier template (`.html`) se compose de deux parties :
1. **L'En-tête (Frontmatter)** : Pour définir le sujet, les destinataires par défaut, etc.
2. **Le Corps** : Le contenu HTML de l'email.

### Exemple Complet

```html
---
subject: Invitation pour {{ RecipientName }}
to: "{{ RecipientEmail }}"
cc: "manager@example.com"
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

## 📧 Destinataires par Défaut (Nouveau !)

Vous pouvez maintenant définir des destinataires par défaut directement dans le frontmatter du template :

```yaml
---
subject: "Relance facture {{ InvoiceNumber }}"
to: "{{ ContactEmail }}"
cc: "comptabilite@example.com"
bcc: "archive@example.com"
---
```

### Champs Disponibles

- **`to`** : Destinataire principal
- **`cc`** : Copie (Carbon Copy)
- **`bcc`** : Copie cachée (Blind Carbon Copy)

### Avantages

- ✅ **Automatisation** : Les emails récurrents ont leurs destinataires pré-remplis
- ✅ **Variables dynamiques** : Vous pouvez utiliser des variables (ex: `{{ ContactEmail }}`)
- ✅ **Texte statique** : Ou définir des emails fixes (ex: `comptabilite@example.com`)
- ✅ **Flexible** : Les flags CLI `--to`, `--cc`, `--bcc` peuvent toujours remplacer ces valeurs

### Ordre de Priorité

```
Template (par défaut) < Flags CLI (override)
```

Si vous définissez `to: "client@example.com"` dans le template mais utilisez `--to "autre@example.com"` en CLI, c'est la valeur CLI qui sera utilisée.

### Exemple Complet

```html
---
subject: "Relance facture {{ InvoiceNumber }}"
to: "{{ ContactEmail }}"
cc: "comptabilite@example.com"
---
<html>
<body>
    <p>Bonjour {{ ContactName }},</p>
    <p>Votre facture {{ InvoiceNumber }} est en attente...</p>
</body>
</html>
```

Lors de l'utilisation :
- Le formulaire demandera `ContactEmail`, `ContactName`, `InvoiceNumber`
- Le destinataire principal sera automatiquement `{{ ContactEmail }}`
- Une copie sera toujours envoyée à `comptabilite@example.com`

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
