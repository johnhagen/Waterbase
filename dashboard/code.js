
// Maybe unsafe in some curcumstances?
let serverURL = location.origin;

function EmptyContents() {
    document.getElementsByClassName("docContentsInsides")[0].innerText = "";
}

function RefreshColumns() {

    EmptyContents()
    const docContainer = document.getElementById('column-documents');
    const colContainer = document.getElementById('column-collections');
    const container = document.getElementById('column-services');

    
    while (docContainer.firstChild) {
        docContainer.removeChild(docContainer.firstChild);
    }

    while (colContainer.firstChild) {
        colContainer.removeChild(colContainer.firstChild);
    }

    while (container.firstChild) {
        container.removeChild(container.firstChild);
    }

}

function RefreshCollections() {

    EmptyContents()
    const docContainer = document.getElementById('column-documents');
    const colContainer = document.getElementById('column-collections');
    
    while (docContainer.firstChild) {
        docContainer.removeChild(docContainer.firstChild);
    }

    while (colContainer.firstChild) {
        colContainer.removeChild(colContainer.firstChild);
    }

}

function RefreshDocuments() {

    EmptyContents()
    const docContainer = document.getElementById('column-documents');

    while (docContainer.firstChild) {
        docContainer.removeChild(docContainer.firstChild);
    }

}



function DeleteService(Name) {

    document.getElementById("statusDiv").innerText = "";
    const url = serverURL + "/waterbase/remove?type=service";

    if (document.getElementById("adminKeyInput").value.length === 0) {
        document.getElementById("statusDiv").innerText = "No Admin Key Specified";
        return;
    };

    if (Name.id) {
        Name = Name.id;
    };

    const data = {
        adminkey: document.getElementById("adminKeyInput").value.toString(),
        servicename: Name.toString()
    };

    console.log(data);

    //console.log(data);

    fetch(url, {
        method: 'DELETE',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('You fucked up fam: ' + response.statusText);
        }
        return response;
    })
    .then(data => {
        console.log("Deleting service: " + Name.toString());
        document.getElementById("column-services").removeChild(document.getElementById(Name.toString()));
        UpdateServices();
        EmptyContents();
    })
    .catch(error => {
        console.error('GET request is fucked mate', error);
    });
}

function DeleteCollection(Name, ServiceName) {

    document.getElementById("statusDiv").innerText = "";
    const url = serverURL + "/waterbase/remove?type=collection";

    if (document.getElementById("adminKeyInput").value.length === 0) {
        document.getElementById("statusDiv").innerText = "No Admin Key Specified";
        return;
    }

    if (Name.id) {
        Name = Name.id;
    }

    if (ServiceName.id) {
        ServiceName = ServiceName.id;
    }

    const data = {
        adminkey: document.getElementById("adminKeyInput").value.toString(),
        servicename: ServiceName.toString(),
        collectionname: Name.toString()
    };


    //console.log("Deleting collection: " + Name + " from service: " + ServiceName)
    fetch(url, {
        method: 'DELETE',
        headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('You fucked up fam: ' + response.statusText);
        }
        return response;
    })
    .then(data => {
        console.log("Deleting collection: " + Name.toString());
        document.getElementById("column-collections").removeChild(document.getElementById(Name.toString()));
        RefreshDocuments();
        EmptyContents();
    })
    .catch(error => {
        console.error('GET request is fucked mate', error);
    });
}

function DeleteDocument(Name, ServiceName, CollectionName) {
    document.getElementById("statusDiv").innerText = "";
    const url = serverURL + "/waterbase/remove?type=document";

    if (document.getElementById("adminKeyInput").value.length === 0) {
        document.getElementById("statusDiv").innerText = "No Admin Key Specified";
        return;
    }

    if (Name.id) {
        Name = Name.id;
    }

    if (ServiceName.id) {
        ServiceName = ServiceName.id;
    }

    if (CollectionName.id) {
        CollectionName = CollectionName.id;
    }

    const data = {
        adminkey: document.getElementById("adminKeyInput").value.toString(),
        servicename: ServiceName.toString(),
        collectionname: CollectionName.toString(),
        documentname: Name.toString()
    };

    //console.log("Document: " + Name + " from service: " + ServiceName)
    fetch(url, {
        method: 'DELETE',
        headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('You fucked up fam: ' + response.statusText);
        }
        return response;
    })
    .then(data => {
        console.log("Deleting document: " + Name.toString());
        document.getElementById("column-documents").removeChild(document.getElementById(Name.toString()));
        EmptyContents();
    })
    .catch(error => {
        console.error('GET request is fucked mate', error);
    });
}

function ListDocContent(ServiceName, CollectionName, DocumentName) {
    document.getElementById("statusDiv").innerText = "";
    const url = serverURL + "/waterbase/retrieve?type=document";

    if (document.getElementById("adminKeyInput").value.length === 0) {
        document.getElementById("statusDiv").innerText = "No Admin Key Specified";
        return;
    }

    if (DocumentName.id) {
        DocumentName = DocumentName.id;
    }

    if (ServiceName.id) {
        ServiceName = ServiceName.id;
    }

    if (CollectionName.id) {
        CollectionName = CollectionName.id;
    }

    //console.log("Document: " + Name + " from service: " + ServiceName)
    fetch(url, {
        method: 'GET',
        headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Adminkey': document.getElementById("adminKeyInput").value.toString(),
        'Servicename': ServiceName.toString(),
        'Collectionname': CollectionName.toString(),
        'Documentname': DocumentName.toString()
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('You fucked up fam: ' + response.statusText);
        }
        return response.json();
    })
    .then(data => {
        document.getElementsByClassName("docContentsInsides")[0].innerText = JSON.stringify(data.content);
    })
    .catch(error => {
        console.error('GET request is fucked mate', error);
    });
}

function CreateButton(Text, Function) {
    let button = document.createElement('button');
    button.textContent = Text;
    button.setAttribute("class", "pageButton");
    button.setAttribute("onclick", Function);
    return button;
}

function CreateCard(Name) {
    let card = document.createElement('div');
    card.className = "card";
    card.id = Name;
    return card;
}

function UpdateServices() {


    RefreshColumns();
    EmptyContents();
    
    const container = document.getElementById('column-services');

    const URL = serverURL + "/waterbase/transmitt?type=services";

    const rawResponse = fetch(URL, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          'AdminKey': document.getElementById("adminKeyInput").value
        },
    }).then(response => {
        if (!response.ok) {
            throw new Error('You fucked up fam' + response.statusText);
        }

        return response.json();
    }).then(data => {
        console.log(data);
        let content = data;

        content.forEach(element => {
        let newCard = CreateCard(element);
        let cardName = document.createElement('h3');
        cardName.innerText = element;
        newCard.setAttribute("key", "Keks")
        newCard.append(cardName);
        newCard.append(CreateButton("Delete", 'DeleteService("' + `${element.toString()}` + '")'));
        newCard.append(CreateButton("Collections", 'UpdateCollections("' + `${element.toString()}` + '")'));

        container.appendChild(newCard);
        });
    })
    .catch(error => {
        console.error('GET request is fucked mate', error);
    });
}

function UpdateCollections(serviceName) {

    const service = document.getElementById(serviceName.toString());

    if (document.getElementById("adminKeyInput").value.length === 0) {
        document.getElementById("statusDiv").innerText = "No Admin Key Specified";
        return;
    }

    RefreshCollections();



    const container = document.getElementById('column-collections');

    const URL = serverURL + "/waterbase/transmitt?type=collections";

    fetch(URL, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          'AdminKey': document.getElementById("adminKeyInput").value,
          'Servicename': service.id.toString()
        },
      })
      .then(response => {
        if (!response.ok) {
            throw new Error('You fucked up fam: ' + response.statusText);
        }

        return response.json();
    }).then(data => {
        console.log(data);

        data.forEach(element => {

        const newCard = document.createElement('div');
        newCard.className = 'card';
        newCard.id = element.toString();
        //newCard.textContent = element;
        let cardName = document.createElement('h3');
        cardName.innerText = element.toString();
        newCard.append(cardName);
        newCard.append(CreateButton("Delete", 'DeleteCollection("' + `${element.toString()}` + '", "' +  `${service.id.toString()}` + '")'));
        newCard.append(CreateButton("Documents", 'UpdateDocuments("' + `${service.id.toString()}` + '", "' + `${element.toString()}` + '")'));
        
        container.appendChild(newCard);
        });
        

        //container.appendChild(newCard);
    })
    .catch(error => {
        console.error('GET request is fucked mate', error);
    });
}

function UpdateDocuments(serviceName, collectionName) {

    const service = document.getElementById(serviceName.toString());
    const collection = document.getElementById(collectionName.toString());

    if (document.getElementById("adminKeyInput").value.length === 0) {
        document.getElementById("statusDiv").innerText = "No Admin Key Specified";
        return;
    }

    RefreshDocuments();

    const container = document.getElementById('column-documents');

    const URL = serverURL + "/waterbase/transmitt?type=documents";

    fetch(URL, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          'Adminkey': document.getElementById("adminKeyInput").value,
          'Servicename': service.id.toString(),
          'Collectionname': collection.id.toString()
        },
      })
      .then(response => {
        if (!response.ok) {
            throw new Error('You fucked up fam: ' + response.statusText);
        }

        return response.json();
    }).then(data => {
        console.log(data);

        data.forEach(element => {

        const newCard = document.createElement('div');
        newCard.className = 'card';
        newCard.id = element.toString();
        let cardName = document.createElement('h3');
        cardName.innerText = element.toString();
        newCard.append(cardName);
        newCard.append(CreateButton("Delete", 'DeleteDocument("' + `${element.toString()}` + '", "' +  `${serviceName.toString()}` + '", "' + `${collectionName.toString()}` + '")'));
        newCard.append(CreateButton("Contents", 'ListDocContent("' + `${service.id.toString()}` + '", "' + `${collection.id.toString()}` + '", "' + `${element.toString()}` + '")'));
        
        container.appendChild(newCard);
        });
        //container.appendChild(newCard);
    })
    .catch(error => {
        console.error('GET request is fucked mate', error);
    });

}
