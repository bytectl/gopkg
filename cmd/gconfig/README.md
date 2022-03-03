# gconfig

## usage

encrypts value

```bash
gconfig -e -k key -v value
# example:
gconfig  -e  -k 123456 -v "foo" 
```

decrypts value

```bash
gconfig -k key -v encrypted_value
# example:
gconfig -k 123456 -v "enc(xxx)" 
```
