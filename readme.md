To deploy:

1. Build the executable:

```bash
make build-linux
```

2. Shut down the existing site:

```bash
ssh-pierre
```

```bash
screen -r
```

```bash
ctrl+C
```

3. Back on the host machine, transfer binary

```bash
scp-pierre
```

4. On the remote, make the binary executable, run and exit screen

```bash
chmod +x app_prod
```

```bash
./app_prod
```

```bash
ctrl+A
```

```bash
ctrl+D
```
