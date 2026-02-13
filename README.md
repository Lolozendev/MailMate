# MailMate — Générateur de brouillons Outlook à partir de templates HTML

MailMate permet de **générer rapidement des emails récurrents** (relances, invitations, confirmations…) à partir de **modèles**.  
Vous choisissez un template, vous remplissez quelques champs, et MailMate ouvre **un brouillon dans Outlook**.

- Aucun envoi automatique : **vous relisez puis cliquez sur “Envoyer” dans Outlook**
- Ciblé **Windows + Outlook Desktop** (usage perso / petite équipe)

---

## Installation

### Option 1 — Télécharger une release (recommandé)

1. Allez sur la page **Releases** du dépôt GitHub
2. Téléchargez `mailmate.exe` (ou l’archive) depuis les assets de la dernière version
3. Placez l’exécutable où vous voulez (ex: `C:\Tools\MailMate\`)

> Téléchargez l’exécutable depuis la page **GitHub Releases** (assets de la dernière version).

### Option 2 — Compiler depuis les sources

- Go **1.21+**

```bash
git clone https://github.com/Lolozendev/MailMate.git
cd MailMate
go mod download
go build ./cmd/mailmate
```

---

## Démarrage rapide

### Pré-requis

- **Windows**
- **Outlook Desktop** installé et configuré

### Lancer (mode interactif)

```powershell
.\mailmate.exe
```

Résultat : Outlook s’ouvre avec
- le sujet rempli
- le corps HTML rendu
- (optionnel) les destinataires (voir templates et flags CLI)

---

## À quoi ça sert ?

- Relances (facture / paiement / devis)
- Invitations / convocations
- Confirmations / suivi
- Messages “standard” qui ne changent que sur quelques variables (nom, date, référence…)

---

## Utilisation

### Mode interactif (recommandé)

```powershell
.\mailmate.exe
```

Sélection du template → formulaire → validation → brouillon Outlook.

### Mode “scriptable” (optionnel)

```powershell
# Exemple
.\mailmate.exe --template templates/invitation.html --kv "Name='John Doe';Date='25-01-2026';Count=5"

# Avec destinataires
.\mailmate.exe --template templates/relance.html --kv "ContactName='Marie';InvoiceNumber=12345;Date='20-01-2026'" --to "marie@example.com"
```

Format des variables (`--kv`) :

- `key1='value';key2='value2';key3=0`
- séparateur : `;`
- guillemets simples ou doubles si espaces

---

## Templates (créer / modifier)

Les templates sont des fichiers `.html` dans `templates/`.

### Exemple minimal

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

### Destinataires par défaut (dans le template)

```html
---
subject: "Relance facture {{ InvoiceNumber }}"
to: "{{ ContactEmail }}"
cc: "comptabilite@example.com"
bcc: "archive@example.com"
---
<html>
<body>
  <p>Bonjour {{ ContactName }},</p>
  <p>Votre facture {{ InvoiceNumber }} est en attente...</p>
</body>
</html>
```

- Les destinataires sont pré-remplis
- Variables dynamiques supportées (ex: `{{ ContactEmail }}`)
- Les flags CLI `--to`, `--cc`, `--bcc` restent prioritaires

Guide complet : **[templates/README.md](./templates/README.md)**

---

## Où mettre mes templates ?

Par défaut, MailMate cherche les templates dans `templates/` (dans le répertoire courant).

Vous pouvez définir un emplacement permanent via la variable d'environnement `MAILMATE_TEMPLATES_DIR`.

**PowerShell (permanent pour l’utilisateur)**

```powershell
[System.Environment]::SetEnvironmentVariable('MAILMATE_TEMPLATES_DIR', 'C:\\MesTemplates', 'User')
```

**PowerShell (session courante uniquement)**

```powershell
$env:MAILMATE_TEMPLATES_DIR = "C:\\MesTemplates"
```

---

## Limites / Notes

- Outlook Web non supporté (Outlook Desktop uniquement)
- L’ouverture du brouillon repose sur l’intégration Outlook locale : Outlook doit être configuré
