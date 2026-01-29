# 📧 MailMate

> **Générez des emails standardisés en un éclair via Outlook, directement depuis votre terminal.**

MailMate simplifie la création d'emails répétitifs. Sélectionnez un modèle, remplissez les informations demandées dans une interface interactive, et laissez l'outil préparer votre brouillon dans Outlook.

---

## ✨ Fonctionnalités

*   🚀 **Rapide & Interactif** : Interface en ligne de commande (TUI) fluide pour saisir les données.
*   🎨 **Templates Flexibles** : Créez des modèles HTML avec variables dynamiques.
*   📫 **Outlook Natif** : Ouvre une fenêtre de rédaction Outlook locale (pas besoin d'accès admin ou API Graph).
*   🛡️ **Sûr** : Vous relisez et envoyez le mail vous-même, rien ne part sans votre validation.

## 🚀 Utilisation Rapide

### Mode Interactif (TUI)

1.  **Lancez l'application** :
    ```powershell
    go run ./cmd/mailmate/main.go
    ```
2.  **Sélectionnez un template** dans la liste.
3.  **Remplissez le formulaire** qui s'affiche.
4.  **Validez** : Outlook s'ouvre avec votre email prêt à partir !

### Mode Ligne de Commande (--kv)

Pour automatiser ou scripter la génération, utilisez le flag `--kv` :

```powershell
# Exemple avec plusieurs variables
./mailmate --template templates/invitation.html --kv "Name='John Doe';Date='25-01-2026';Count=5"

# Avec destinataires
./mailmate --template templates/relance.html --kv "ContactName='Marie';InvoiceNumber=12345;Date='20-01-2026'" --to "marie@example.com"
```

**Format** : `key1='value';key2='value2';key3=0`
- Séparez les paires par des points-virgules (`;`)
- Utilisez des guillemets simples ou doubles pour les valeurs avec espaces
- Les validations de type (int, date, filepath) s'appliquent automatiquement

## 📝 Créer vos Templates

C'est le cœur de l'outil ! Ajoutez vos fichiers `.html` dans le dossier `templates/`.

Un template ressemble à ça :

```html
---
subject: Relance facture {{ InvoiceNumber }}
---
<html>
<body>
    <p>Bonjour {{ ContactName }},</p>
    <p>Sauf erreur de notre part, la facture {{ InvoiceNumber }} du {{ Date | type:'date' }} est en attente.</p>
</body>
</html>
```

👉 **[Guide complet pour créer des templates](./templates/README.md)** (Syntaxe, variables, filtres...)

## 🔧 Configuration

### Dossier des Templates

Par défaut, MailMate cherche les templates dans le dossier `templates/` du répertoire courant.

Pour définir un emplacement personnalisé **de façon permanente**, utilisez la variable d'environnement `MAILMATE_TEMPLATES_DIR` :

**Windows PowerShell (permanent pour l'utilisateur)** :
```powershell
[System.Environment]::SetEnvironmentVariable('MAILMATE_TEMPLATES_DIR', 'C:\MesTemplates', 'User')
```

**Windows PowerShell (session courante uniquement)** :
```powershell
$env:MAILMATE_TEMPLATES_DIR = "C:\MesTemplates"
```

**Linux/macOS (permanent dans ~/.bashrc ou ~/.zshrc)** :
```bash
export MAILMATE_TEMPLATES_DIR="/home/user/mes-templates"
```

**Linux/macOS (session courante uniquement)** :
```bash
export MAILMATE_TEMPLATES_DIR="/home/user/mes-templates"
```

Une fois configurée, l'application utilisera automatiquement ce dossier sans configuration supplémentaire.

## ⚙️ Installation & Pré-requis

*   **OS** : Windows uniquement (dépendance à Outlook Desktop).
*   **Logiciel** : Microsoft Outlook installé et configuré.
*   **Go** : Go 1.21+ pour compiler.

```bash
# Cloner le repo
git clone https://github.com/votre-repo/mailmate.git

# Installer les dépendances
go mod download
```

---
*Note: Ancienne documentation technique disponible dans [README.md.old](./README.md.old)*
