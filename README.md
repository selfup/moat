# Moat

Make your files safe before saving them to the cloud

## In Development - DO NOT USE - IN DEBUG MODE

```
selfup@win42 MINGW64 ~/go/src/github.com/selfup/moat (master)
$ moat -cmd=push -service=fixtures
Moat path is: C:\Users\selfup\Moat
Service path is: fixtures\Moat

Encrypted: C:\Users\selfup\Moat\omg.txt - to: fixtures\Moat\omg.txt
Encrypted: C:\Users\selfup\Moat\wow.txt - to: fixtures\Moat\wow.txt

selfup@win42 MINGW64 ~/go/src/github.com/selfup/moat (master)
$ moat -cmd=pull -service=fixtures
Moat path is: C:\Users\selfup\Moat
Service path is: fixtures\Moat

Decrypted: fixtures\Moat\omg.txt - to: C:\Users\selfup\Moat\omg.txt
Decrypted: fixtures\Moat\wow.txt - to: C:\Users\selfup\Moat\wow.txt
```

## Help

```
$ moat -h
Usage of moat.exe:
  -cmd string
        REQUIRED
                main command
                push will encrypt Moat/filename.ext to Service/Moat/filename.ext
                pull will decrypt from Service/Moat/filename.ext to Moat/filename.ext
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
