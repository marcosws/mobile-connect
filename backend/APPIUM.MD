# Integração Appium

O Mobile Connect atua como um gateway para o Appium. Ao iniciar, ele verifica se
o Appium já responde no endereço configurado; caso contrário, inicia o binário
`appium`, aguarda o health check e o encerra quando o Mobile Connect recebe um
sinal de parada. Em seguida recebe as requisições WebDriver, seleciona o
dispositivo Android pelo ID do Mobile Connect e as encaminha ao Appium.

O principal benefício é evitar que um teste ou uma sessão do Inspector escolha
acidentalmente outro aparelho conectado por ADB. Na criação da sessão, o gateway
sempre substitui `appium:udid` pelo `{id}` presente na URL.

## Pré-requisitos

1. Tenha `adb` disponível e um dispositivo em `device` em `GET /devices`.
2. Instale o Appium e o driver Android necessário. Um exemplo para UiAutomator2
   é:

   ```bash
   npm install -g appium
   appium driver install uiautomator2
   ```

3. Inicie o Mobile Connect. O Appium será iniciado automaticamente em
   `127.0.0.1:4723`:

   ```bash
cd backend
APPIUM_URL=http://127.0.0.1:4723 go run ./cmd/api
```

   Esse é o valor padrão, portanto a variável pode ser omitida nesse caso. Para
   usar outra porta ou base, configure `APPIUM_URL`; o Mobile Connect repassa
   automaticamente endereço, porta e `--base-path` ao Appium:

   ```bash
   cd backend
   APPIUM_URL=http://127.0.0.1:4723/wd/hub go run ./cmd/api
   ```

### Configuração do processo Appium

| Variável | Padrão | Uso |
| --- | --- | --- |
| `APPIUM_URL` | `http://127.0.0.1:4723` | Endereço do Appium iniciado ou utilizado pelo gateway. |
| `APPIUM_COMMAND` | `appium` | Caminho/nome do executável Appium. |
| `APPIUM_AUTO_START` | `true` | Use `false` para conectar a um Appium já gerenciado externamente. |

O auto-start é intencionalmente limitado a `localhost`, `127.0.0.1` ou `::1`.
Para um Appium remoto, use `APPIUM_AUTO_START=false`; o gateway não deve iniciar
processos em outra máquina. Se o binário não estiver no `PATH`, defina, por
exemplo, `APPIUM_COMMAND=/opt/node/bin/appium`.

Confirme a comunicação entre o gateway e o Appium:

```bash
curl http://localhost:8080/appium/status
```

O retorno é a resposta original de `GET /status` do Appium. Um `502 Bad Gateway`
significa que o Mobile Connect não conseguiu alcançá-lo.

## Rotas

| Rota | Uso |
| --- | --- |
| `GET /appium/status` | Health check do Appium configurado. |
| `POST /devices/{id}/appium/sessions` | Criação explícita de sessão pelo gateway. |
| `POST /devices/{id}/appium/session` | Criação W3C compatível com clientes WebDriver. |
| `/devices/{id}/appium/session/{sessionId}/...` | Comandos da sessão para Selenium e Appium Inspector. |
| `/appium/session/{sessionId}/...` | Encaminhamento de comandos de uma sessão existente. |

Prefira a base `http://localhost:8080/devices/{id}/appium` nos clientes. Ela é
compatível com o protocolo WebDriver: o cliente acrescenta `/session` ao criar a
sessão e reutiliza a mesma base nos comandos seguintes.

Não use `POST /appium/session` diretamente. Essa rota é rejeitada de propósito,
pois não contém o ID do dispositivo que deve ser fixado.

## Capabilities Android

Appium 2 usa o formato W3C. `platformName` é uma capability padrão; capabilities
do Appium devem levar o prefixo `appium:`, como `appium:automationName` e
`appium:appPackage`.

Exemplo para abrir um app já instalado:

```json
{
  "capabilities": {
    "alwaysMatch": {
      "platformName": "Android",
      "appium:automationName": "UiAutomator2",
      "appium:deviceName": "Android",
      "appium:appPackage": "com.exemplo.app",
      "appium:appActivity": ".MainActivity",
      "appium:noReset": true
    },
    "firstMatch": [{}]
  }
}
```

`appium:udid` não precisa ser informado: se estiver presente, será sobrescrito
pelo ID da URL. Para automatizar Chrome no Android, substitua `appPackage` e
`appActivity` por `browserName: "Chrome"`.

Teste a criação da sessão sem uma biblioteca cliente:

```bash
curl -X POST http://localhost:8080/devices/emulator-5554/appium/session \
  -H 'Content-Type: application/json' \
  -d '{
    "capabilities": {
      "alwaysMatch": {
        "platformName": "Android",
        "appium:automationName": "UiAutomator2",
        "appium:appPackage": "com.exemplo.app",
        "appium:appActivity": ".MainActivity"
      },
      "firstMatch": [{}]
    }
  }'
```

Guarde o `value.sessionId` da resposta e encerre a sessão ao final:

```bash
curl -X DELETE http://localhost:8080/devices/emulator-5554/appium/session/SEU_SESSION_ID
```

## Uso com Appium Java Client (Selenium)

O Appium Java Client é a opção adequada para testes mobile em Java e usa o
protocolo Selenium/WebDriver por baixo. Dependências e versões devem ser
compatíveis com a versão do Selenium usada pelo projeto.

```java
import io.appium.java_client.android.AndroidDriver;
import io.appium.java_client.android.options.UiAutomator2Options;
import java.net.URL;

UiAutomator2Options options = new UiAutomator2Options()
    .setDeviceName("Android")
    .setAppPackage("com.exemplo.app")
    .setAppActivity(".MainActivity")
    .setNoReset(true);

AndroidDriver driver = new AndroidDriver(
    new URL("http://localhost:8080/devices/emulator-5554/appium"),
    options
);

// Localização e interação usam a API Selenium normalmente.
driver.findElement(AppiumBy.accessibilityId("Entrar")).click();
driver.quit();
```

Não defina `udid` no teste; a URL é a fonte de verdade para o aparelho. Para
usar um aparelho físico, substitua `emulator-5554` pelo campo `id` retornado por
`GET /devices`.

### Selenium puro para Chrome Android

Quando o teste só precisa da API WebDriver de navegador, um `RemoteWebDriver`
também pode ser usado:

```java
import org.openqa.selenium.MutableCapabilities;
import org.openqa.selenium.remote.RemoteWebDriver;
import java.net.URL;

MutableCapabilities caps = new MutableCapabilities();
caps.setCapability("platformName", "Android");
caps.setCapability("browserName", "Chrome");
caps.setCapability("appium:automationName", "UiAutomator2");

RemoteWebDriver driver = new RemoteWebDriver(
    new URL("http://localhost:8080/devices/emulator-5554/appium"), caps
);
driver.get("https://example.com");
driver.quit();
```

Para gestos e comandos específicos de mobile, prefira o Appium Java Client.

## Uso com Appium Python Client

```python
from appium import webdriver
from appium.options.android import UiAutomator2Options

options = UiAutomator2Options()
options.platform_name = "Android"
options.automation_name = "UiAutomator2"
options.app_package = "com.exemplo.app"
options.app_activity = ".MainActivity"
options.no_reset = True

driver = webdriver.Remote(
    "http://localhost:8080/devices/emulator-5554/appium",
    options=options,
)
driver.find_element("accessibility id", "Entrar").click()
driver.quit()
```

## Uso com Appium Inspector

Use o aplicativo **Appium Inspector** (não o Appium Desktop, que foi
descontinuado para Appium 2). Crie uma sessão em **Start Session** com os valores
abaixo:

| Campo do Inspector | Valor |
| --- | --- |
| Remote Host | `127.0.0.1` ou o host do Mobile Connect |
| Remote Port | `8080` |
| Remote Path | `/devices/emulator-5554/appium` |
| SSL | desativado, salvo se o Mobile Connect estiver atrás de HTTPS |

No editor de capabilities, informe somente as capabilities da aplicação. Por
exemplo:

```json
{
  "platformName": "Android",
  "appium:automationName": "UiAutomator2",
  "appium:deviceName": "Android",
  "appium:appPackage": "com.exemplo.app",
  "appium:appActivity": ".MainActivity",
  "appium:noReset": true
}
```

Clique em **Start Session**. O Inspector fará `POST` para
`/devices/emulator-5554/appium/session`; o gateway adicionará o `udid` correto e
encaminhará a sessão ao Appium. Depois disso, use as abas **Source**, **Commands**
e os seletores gerados pelo Inspector normalmente. Ao terminar, encerre a sessão
no Inspector para liberar recursos no dispositivo e no Appium.

Se o Inspector estiver em outra máquina, o Mobile Connect precisa estar acessível
por essa rede e o CORS/firewall deve permitir a conexão. Não exponha este gateway
diretamente à internet: nesta implementação ainda não há autenticação e ele pode
controlar dispositivos conectados.

## Diagnóstico rápido

| Sintoma | Verificação |
| --- | --- |
| `502 Bad Gateway` | Confirme `APPIUM_URL` e consulte `/appium/status`; veja os logs do Appium no mesmo terminal do Mobile Connect. |
| `NoSuchDriver` ou driver ausente | Execute `appium driver install uiautomator2` e reinicie o Appium. |
| Dispositivo não encontrado | Confirme `adb devices` e `GET /devices`; use exatamente o `id` retornado. |
| Erro de capability | Use W3C e prefixe capabilities não padronizadas com `appium:`. |
| `404` com `/wd/hub` | Alinhe `APPIUM_URL` ao `--base-path` real do Appium; Appium 2 normalmente usa a raiz `/`. |

## Referências

- [Capabilities do Appium](https://appium.io/docs/en/latest/guides/caps/)
- [Protocolo WebDriver suportado pelo Appium](https://appium.io/docs/en/3.1/reference/api/webdriver/)
- [Migração e capabilities no Appium 2](https://appium.io/docs/en/2.19/guides/migrating-1-to-2/)
- [Documentação do Appium Inspector](https://github.com/appium/appium-inspector/tree/main/docs)
