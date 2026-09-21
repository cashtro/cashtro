targetScope = 'subscription'

@description('Name of the azd environment. Used for resource group and unique hashes.')
param environmentName string

@description('Azure region for the resource group and App Service. Canada Central is the default for Castro / Proximity.')
param location string = 'canadacentral'

@secure()
@description('Optional OpenRouter key. Leave empty to boot the kernel unbound.')
param openRouterApiKey string = ''

@description('Optional OpenRouter model id.')
param openRouterModel string = 'openai/gpt-4o-mini'

var resourceToken = uniqueString(subscription().id, location, environmentName)
var tags = {
  'azd-env-name': environmentName
  project: 'cashtro'
}

resource rg 'Microsoft.Resources/resourceGroups@2024-03-01' = {
  name: 'rg-${environmentName}'
  location: location
  tags: tags
}

module resources './modules/resources.bicep' = {
  name: 'resources'
  scope: rg
  params: {
    name: 'cashtro'
    location: location
    tags: tags
    resourceToken: resourceToken
    openRouterApiKey: openRouterApiKey
    openRouterModel: openRouterModel
  }
}

output AZURE_RESOURCE_GROUP string = rg.name
output AZURE_LOCATION string = location
output AZURE_LOG_ANALYTICS_WORKSPACE_ID string = resources.outputs.logAnalyticsWorkspaceId
output WEB_URL string = resources.outputs.webUrl
output CASHTRO_URL string = resources.outputs.webUrl
