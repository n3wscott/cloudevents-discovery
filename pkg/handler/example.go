package handler

const exampleServices = `[{
  "id": "3db60532-e839-417e-8644-e255f338776a",
  "url": "https://storage.example.com/service/storage",
  "name": "storage",
  "description": "Blob storage in the cloud",
  "protocols": ["HTTP"],
  "subscriptionurl": "https://cloud.example.com/docs/storage",
  "types": [{
      "type": "com.example.storage.object.create",
      "specversions": [ "1.x-wip" ],
      "datacontenttype": "application/json",
      "dataschema": "http://schemas.example.com/download/com.example.storage.object.create.json",
      "sourcetemplate": "https://storage.example.com/service/storage/{objectID}"
    }]
},{
  "id": "cbdd62e8-c095-11ea-b3de-0242ac130004",
  "url": "https://example.com/services/widgetService",
  "name": "widgets",
  "specversions": [ "1.0" ],
  "subscriptionurl": "https://events.example.com",
  "protocols": [ "HTTP" ],
  "types": [{
      "type": "com.example.widget.create"
    }, {
      "type": "com.example.widget.delete"
    }]
}]`
