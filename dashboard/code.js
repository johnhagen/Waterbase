function DeleteService(Name) {

    const url = "http://localhost:8080/waterbase/remove?type=service";

    const data = {
        auth: document.getElementById("serviceKeyInput").value,
        servicename: Name.toString()
    };

    console.log(data);

    console.log(document.getElementById("serviceKeyInput").value);

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
            throw new Error('You fucked up fam' + response.statusText);
        }
        return response;
    })
    .then(data => {
        console.log("Deleting service: " + Name);
        document.getElementById("column-services").removeChild(document.getElementById(Name));
    })
    .catch(error => {
        console.error('GET request is fucked mate', error);
    });
}

function CreateButton(Text, Function) {
    let button = document.createElement('button');
    button.textContent = Text;
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
    const container = document.getElementById('column-services');

    const URL = 'http://localhost:8080/waterbase/transmitt?type=services';

    while (container.firstChild) {
        container.removeChild(container.firstChild)
    }

    /*(async () => {
        const rawResponse = await fetch(URL, {
            method: 'GET',
            headers: {
              'Accept': 'application/json',
              'Content-Type': 'application/json',
              'AdminKey': document.getElementById("adminKeyInput").value
            },
          });

        const content = await rawResponse.json();
        console.log(document.getElementById("adminKeyInput").value)

        console.log(content)

        content.forEach(element => {
        var newCard = document.createElement('div');
        var button = document.createElement('button');
        button.textContent = 'Get Collections';
        button.onclick = `UpdateCollections(${element}, ${123})`;
        newCard.appendChild(button);
        newCard.className = 'card';
        newCard.textContent = element;

        container.appendChild(newCard);
            
        });
    })(); */

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
        newCard.append(CreateButton("Delete Service", `DeleteService(${element})`));
        newCard.append(CreateButton("Get Collections", `UpdateCollections(${element}, ${123})`));

        container.appendChild(newCard);
        });
    })
    .catch(error => {
        console.error('GET request is fucked mate', error);
    });
}

function UpdateCollections(service ,key) {
    const container = document.getElementById('column-collections');

    const URL = 'http://localhost:8080/waterbase/transmitt?type=collections';

    while (container.firstChild) {
        container.removeChild(container.firstChild)
    }

    (async () => {
        const rawResponse = await fetch(URL, {
            method: 'GET',
            headers: {
              'Accept': 'application/json',
              'Content-Type': 'application/json',
              'Auth': key,
              'Servicename': service
            },
          });

        const content = await rawResponse.json();

        console.log(content)

        content.forEach(element => {
        const newCard = document.createElement('div');
        newCard.className = 'card';
        newCard.textContent = element;

        container.appendChild(newCard);
            
        });
    })();
}
