# Moat

Make your files safe before saving them to the cloud

#### TODO

`StartPrompt`

- once a command is given then run the scanner
- if this is the first time service boots up ask for a password
- RSA priv/pub will be created
- RSA pub will be stored in Service/Moat/rsa_id.pub
- RSA priv will be stored in Moat/rsa_id
- Machine is assumed to be encrypted and safe
- if the command says to encrypt, scan all files, encrypt in memory, and then write to Service/Moat
- if the command says to decrypt, read all files in moat, decrypt what is in memory, write to Moat decrypted
- open -> read contents -> keep contents in memory -> check RSA -> find aes256 key (64 characters) -> encrypt or decrypt
