# MailMate (Mail templating pour Outlook, sans prise de tête)

MailMate aide à **rédiger plus vite des emails récurrents** (relances, invitations, confirmations…), à partir de **modèles**. 
Vous choisissez un modèle, vous remplissez 2–3 champs, et MailMate ouvre **un brouillon dans Outlook**.

- Rien n’est envoyé automatiquement : **vous relisez et cliquez sur “Envoyer” dans Outlook**.
- Pensé pour un usage perso / équipe, sur **Windows + Outlook Desktop**.

---

## À quoi ça sert ?

- Relance de facture / paiement
- Invitation / convocation
- Compte-rendu / suivi
- Messages “standard” qui changent juste sur quelques variables (nom, date, référence…)

---

## Fonctionnement (simple)

1. Vous lancez MailMate
2. Vous choisissez un modèle
3. Vous remplissez les champs demandés
4. Outlook s’ouvre avec un **brouillon prêt**

---

## Démarrage rapide

### 1) Pré-requis

- **Windows**
- **Outlook Desktop** installé et configuré

### 2) Lancer

```powershell
mailmate.exe
```

### 3) Résultat

Une fenêtre Outlook s’ouvre avec :
- le sujet déjà rempli
- le corps HTML rendu
- (optionnel) les destinataires si vous les fournissez en CLI

---

## Utilisation

### Mode interactif (recommandé)

Lancez simplement l’exécutable :

```powershell
mailmate.exe
```

Puis : sélection du template → formulaire → validation → brouillon Outlook.

### Mode “scriptable” (optionnel)

Si vous voulez automatiser depuis un script (CI perso, raccourci, PowerShell…), vous pouvez passer :
- le template
- les variables
- (optionnel) les destinataires

```powershell
# Exemple
./mailmate.exe --template templates/invitation.html --kv \"Name='John Doe';Date='25-01-2026';Count=5\"

# Avec destinataires
./mailmate.exe --template templates/relance.html --kv \"ContactName='Marie';InvoiceNumber=12345;Date='20-01-2026'\" --to \"marie@example.com\"
```

Format des variables (`--kv`) :

- `key1='value';key2='value2';key3=0`
- séparateur : `;`
- guillemets simples ou doubles si espaces

---

## Créer vos modèles (templates)

Les templates sont des fichiers `.html` dans `templates/`.

### Exemple basique :

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

### Nouveau : Destinataires par défaut

Vous pouvez désormais définir des destinataires par défaut directement dans vos templates :

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

**Avantages :**
- ✅ Les destinataires sont automatiquement pré-remplis
- ✅ Support des variables dynamiques (ex: `{{ ContactEmail }}`)
- ✅ Les flags CLI `--to`, `--cc`, `--bcc` peuvent toujours les remplacer

**Priorité :** Template (défaut) < Flags CLI (override)

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
$env:MAILMATE_TEMPLATES_DIR = \"C:\\MesTemplates\"
```

---

## Limites / Notes

- Projet orienté **Outlook Desktop** : pas de support “Outlook Web”.
- L’ouverture du brouillon repose sur l’intégration Outlook locale : si Outlook n’est pas configuré, ça ne marchera pas.

---

## Développeurs (build)

Si vous voulez compiler vous-même :

- Go **1.21+**

```bash
git clone https://github.com/Lolozendev/MailMate.git
go mod download
```


> **Générez des emails standardisés en un éclair via Outlook, directement depuis votre terminal.**

MailMate simplifie la création d'emails répétitifs. Sélectionnez un modèle, remplissez les informations demandées dans une interface interactive, et laissez l'outil préparer votre brouillon dans Outlook.
