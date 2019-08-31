# Moat

Make your files safe before saving them to the cloud

## In Development - DO NOT USE

```
$ moat -service="fixtures"
Moat path is: C:\Users\selfup\Moat
Service path is: fixtures\Moat

Label Key written to: C:\Users\selfup\Moat\moatprivate
Private Key written to: C:\Users\selfup\Moat\moatprivate
Public Key written to: fixtures\Moat\moatpublic
Encrypted AES Key written to: fixtures\Moat\moatkey

$ echo "wow this is going to be encrypted and saved to a cloud service directory" >> ~/Moat/wow.txt

$ moat -service="fixtures" -cmd=push
Moat path is: C:\Users\selfup\Moat
Service path is: fixtures\Moat

Encrypted: C:\Users\selfup\Moat\wow.txt - to: fixtures\Moat\wow.txt

$ moat -service="fixtures" -cmd=pull
Moat path is: C:\Users\selfup\Moat
Service path is: fixtures\Moat

Decrypted: fixtures\Moat\wow.txt - to: C:\Users\selfup\Moat\wow.txt
```

## What does this do?

1. Encrypts files to other directories
1. Decrypts files from other directories back into your vault

By default `$HOME/Moat` or `%USERPROFILE%\Moat` are used for your Vault. This ensures any file written here will be encrypted to any service directory you want (Dropbox/Google Drive/A USB device/etc..).

Typically you would say something like: `-service=$HOME/Dropbox`

And when you run `moat` it will encrypt files from your Vault to `$HOME/Dropbox/Moat`.

You can also define your own Vault path, but we can talk about that later.

## How is this safe?

1. Keys are randomly generated (aesKey, label)
1. RSA 4096
1. AES 256

#### In your Vault:

Private RSA key is stored. As well as a random 32 byte label.

#### In your service directory:

Public RSA key as well as an encrypted (using said public RSA key) randomly generated 32 byte passphrase.

#### Decryption/Encryption of files

All files are encrypted using a 32byte passphrase paired with a 32 byte label.

All files are decrypted using the decrypted 32byte passphrase in your service directory paired with the local 32 byte label.

## Custom Paths (Vaults)

```
$ moat -home="archive" -service="fixtures"
Moat path is: archive\Moat
Service path is: fixtures\Moat

Label Key written to: archive\Moat\moatprivate
Private Key written to: archive\Moat\moatprivate
Public Key written to: fixtures\Moat\moatpublic
Encrypted AES Key written to: fixtures\Moat\moatkey

$ echo "wow this is going to be encrypted and saved to a cloud service directory" >> archive/Moat/wow.txt

$ moat -home="archive" -service="fixtures" -cmd=push
Moat path is: archive\Moat
Service path is: fixtures\Moat

Encrypted: archive\Moat\wow.txt - to: fixtures\Moat\wow.txt

$ moat -home="archive" -service="fixtures" -cmd=pull
Moat path is: archive\Moat
Service path is: fixtures\Moat

Decrypted: fixtures\Moat\wow.txt - to: archive\Moat\wow.txt
```

## Help

```
$ moat -h
Usage of moat:
  -cmd string
        REQUIRED
                main command
                push will encrypt Moat/filename.ext to Service/Moat/filename.ext
                pull will decrypt from Service/Moat/filename.ext to Moat/filename.ext
  -home string
        OPTIONAL
                Home dir (here you want Moat to be created at) - defaults to $HOME or USERPROFILE
  -moat string
        OPTIONAL
                What you want Moat to be called - essentially Vault names
  -service string
        REQUIRED
                Directory of cloud service that will sync on update
```
