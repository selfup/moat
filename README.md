# Moat

Make your files safe before saving them to the cloud

## In Development - DO NOT USE

```
$ moat -service="fixtures"
Moat path is: C:\Users\selfup\Moat
Service path is: fixtures\Moat

Private Key written to: C:\Users\selfup\Moat\privatemoatssh
Public Key written to: fixtures\Moat\publicmoatssh
Encrypted AES Key written to: fixtures\Moat\aesKey

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

## Custom Paths (Vaults)

```bash
$ moat -home="archive" -service="fixtures" -cmd=push
Moat path is: archive\Moat
Service path is: fixtures\Moat

Private Key written to: archive\Moat\privatemoatssh
Public Key written to: fixtures\Moat\publicmoatssh
Encrypted AES Key written to: fixtures\Moat\aesKey
Encrypted: archive\Moat\wow.txt - to: fixtures\Moat\wow.txt

$ moat -home="archive" -service="fixtures" -cmd=pull
Moat path is: archive\Moat
Service path is: fixtures\Moat

Decrypted: fixtures\Moat\wow.txt - to: archive\Moat\wow.txt

$ moat -home="archive" -service="fixtures" -cmd=push
Moat path is: archive\Moat
Service path is: fixtures\Moat

Encrypted: archive\Moat\wow.txt - to: fixtures\Moat\wow.txt
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

#### TODO

`StartPrompt`

- ~~once a command is given then run the scanner~~
- if this is the first time service boots up ask for a password
- RSA priv/pub will be created
- RSA pub will be stored in Service/Moat/rsa_id.pub
- RSA priv will be stored in Moat/rsa_id
- ~~if the command says to push, scan all files, encrypt in memory, and then write to Service/Moat~~
- ~~if the command says to pull, read all files in moat, decrypt what is in memory, write to Moat decrypted~~
- open -> read contents -> keep contents in memory -> check RSA -> find aes256 key (64 characters) -> encrypt or decrypt
