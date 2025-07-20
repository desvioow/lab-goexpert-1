# API de Consulta de Temperatura por CEP

## Descrição
Este projeto é uma API REST desenvolvida em Go que permite consultar a temperatura atual de uma cidade brasileira a partir de um CEP (Código de Endereçamento Postal). A aplicação integra-se com a API ViaCEP para obter informações de localização a partir do CEP fornecido e, em seguida, utiliza a API WeatherAPI para obter dados meteorológicos da cidade correspondente.

## Funcionalidades
- Consulta de temperatura a partir de um CEP brasileiro
- Validação de formato de CEP (deve conter 8 dígitos numéricos)
- Retorno de temperatura em três escalas: Celsius (°C), Fahrenheit (°F) e Kelvin (K)
- Tratamento de erros para CEPs inválidos, inexistentes ou cidades não reconhecidas

## Tecnologias Utilizadas
- Go 1.22
- Chi Router para gerenciamento de rotas HTTP
- Docker para containerização
- Google Cloud Run para deploy

## Endpoints

### GET /{cep}
Retorna a temperatura atual da cidade correspondente ao CEP informado.

**Parâmetros:**
- `cep`: CEP brasileiro com 8 dígitos numéricos (sem hífen)

**Respostas:**
- `200 OK`: Retorna os dados de temperatura no formato JSON
- `404 Not Found`: CEP não encontrado ou cidade não reconhecida pela API de clima
- `422 Unprocessable Entity`: Formato de CEP inválido
- `500 Internal Server Error`: Erro interno do servidor

**Exemplo de resposta bem-sucedida:**
```json
{
  "temp_C": "25.0",
  "temp_F": "77.0",
  "temp_K": "298.0"
}
```

## Como testar a aplicação no Google CLoud Run
# Buscar CEP válido
```json
curl -X GET \
  -H "Content-Type: application/json" \
  https://lab-goexpert-1-392400719582.southamerica-east1.run.app/69402289
```
# Buscar CEP inválido
```json
curl -X GET \
  https://lab-goexpert-1-392400719582.southamerica-east1.run.app/1234567890
```
# Buscar CEP válido mas inexistente
```json
curl -X GET \
  -H "Content-Type: application/json" \
  https://lab-goexpert-1-392400719582.southamerica-east1.run.app/08556111
```
# Buscar CEP válido mas cidade(Poá) não considerada válida na API da weatherapi.com
```json
curl -X GET \
  -H "Content-Type: application/json" \
  https://lab-goexpert-1-392400719582.southamerica-east1.run.app/08556450
```

## Deploy
A aplicação está configurada para ser implantada no Google Cloud Run. O arquivo Dockerfile fornece as instruções necessárias para a criação da imagem de container.

## APIs Externas Utilizadas
- [ViaCEP](https://viacep.com.br/) - API para consulta de CEPs brasileiros
- [WeatherAPI](https://www.weatherapi.com/) - API para consulta de dados meteorológicos

## Exemplos de Uso
Consulte o arquivo `requests.http` para exemplos de requisições à API, tanto localmente quanto na versão implantada no Google Cloud Run.