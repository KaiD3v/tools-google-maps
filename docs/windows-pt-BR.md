# Executar no Windows

Este fork inclui uma cópia fixada de `scrapemate v1.3.0` com as correções de
ciclo de vida das páginas e inicialização do Chromium. O build normal já usa
essa cópia por meio do `go.mod`; não é necessário editar o cache do Go.

## Compilar

Instale o Git e o Go 1.27.0 ou uma versão estável posterior. Abra novamente
o PowerShell após instalar o Go para atualizar o PATH.

```powershell
git clone https://github.com/KaiD3v/tools-google-maps.git
Set-Location tools-google-maps
go version
go mod download
go mod verify
go build -o .\bin\gms.exe .
```

Interrompa o procedimento se qualquer comando falhar.

## Instalar o navegador

Execute uma vez, na pasta do projeto. Não é necessário fazer login no Google.

```powershell
$env:DISABLE_TELEMETRY = '1'
$env:PLAYWRIGHT_DRIVER_PATH = Join-Path $PWD 'cache\playwright-driver'
$env:PLAYWRIGHT_BROWSERS_PATH = Join-Path $PWD 'cache\playwright-browsers'
$env:PLAYWRIGHT_INSTALL_ONLY = '1'
try {
    & .\bin\gms.exe
    if ($LASTEXITCODE -ne 0) { throw 'Falha ao instalar o Playwright.' }
} finally {
    Remove-Item Env:PLAYWRIGHT_INSTALL_ONLY -ErrorAction SilentlyContinue
}
```

## Iniciar e parar

Sempre execute na pasta do projeto. Repita as variáveis ao abrir outro terminal.

```powershell
$env:DISABLE_TELEMETRY = '1'
$env:PLAYWRIGHT_DRIVER_PATH = Join-Path $PWD 'cache\playwright-driver'
$env:PLAYWRIGHT_BROWSERS_PATH = Join-Path $PWD 'cache\playwright-browsers'
Remove-Item Env:PLAYWRIGHT_INSTALL_ONLY -ErrorAction SilentlyContinue
& .\bin\gms.exe -web -addr 127.0.0.1:8080 -data-folder .\webdata -c 1
```

Abra <http://127.0.0.1:8080/app> e mantenha o terminal aberto. Para parar,
pressione **Ctrl+C no terminal**. Fechar a aba do navegador não encerra o
servidor. Os resultados ficam em `webdata/`.

## Primeira busca

- Keywords: uma busca por linha, por exemplo `oficinas mecânicas em Jundiaí SP`.
- Language: `pt`.
- Latitude e longitude: deixe ambas vazias para usar a cidade do texto;
  não mantenha `0,0` se não deseja buscar nessa coordenada.
- Depth: `1`; Max job time: `3m` ou mais.
- Fast Mode e Fetch Emails: desmarcados; Proxies: vazio.

Clique em **Start Scraping** e depois em **Download**. Confira se o CSV tem
registros: a interface ainda pode mostrar `ok` quando erros internos não são
propagados. Em caso de arquivo vazio, examine as mensagens no terminal.

## Testes e limites

```powershell
go test github.com/gosom/scrapemate/adapters/fetchers/jshttp
go test ./gmaps ./runner/... ./web
```

Os testes de ciclo de vida não consultam o Google. A validação manual no
Windows retornou 20 oficinas e 20 clínicas em duas buscas pequenas, com nome,
telefone e endereço. Isso não garante a quantidade ou o sucesso de outras
consultas. Nenhum CSV dessa validação está incluído no repositório.

A interface não possui autenticação: mantenha o endereço em `127.0.0.1` e
não publique a porta na internet. `DISABLE_TELEMETRY=1` desativa o envio de
eventos ao PostHog, mas o código original ainda pode consultar serviços de IP
público ao construir eventos. As correções deste fork são funcionais; não
resolvem todos os riscos de segurança e dependências do projeto original.
