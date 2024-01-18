document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('retrieveAuthorizedUsersButton').addEventListener('click', function() {
        fetch('https://admin.ryanschnabel.com/bff/getUsers', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        })
        .then(response => response.text())
        .then(data => {
            const usersArray = data.split(',');

            const usersDiv = document.getElementById('whitelistedUsers');
            usersDiv.innerHTML = '';

            usersArray.forEach(user => {
                const userSpan = document.createElement('span');
                userSpan.textContent = user.trim();
                usersDiv.appendChild(userSpan);

                if (user !== 'ryan.d.schnabel@gmail.com' && user !== 'ryan.schnabel@gmail.com') {
                    const removeButton = document.createElement('button');
                    removeButton.textContent = 'Remove User';
                    removeButton.onclick = function() { removeUser(user); };
                    usersDiv.appendChild(removeButton);
                }

                usersDiv.appendChild(document.createElement('br'));
            });
        })
        .catch(error => {
            console.error('Error:', error);
        });
    });
    document.getElementById('addUserButton').addEventListener('click', function() {
        const inputText = document.getElementById('inputText').value;
        const requestData = {
            endpointVar: 'addUser',
            inputText: inputText
        };

        fetch('https://admin.ryanschnabel.com/bff/', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(requestData)
        })
        .then(response => response.text())
        .then(data => {
            // Process the response data
            console.log(data);
            // You can also update the UI here based on the response
        })
        .catch(error => {
            console.error('Error:', error);
        });
    });
});

function removeUser(userEmail) {

    const requestData = {
        endpointVar: 'removeUser',
        inputText: userEmail
    };

    fetch('https://admin.ryanschnabel.com/bff/', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
    })
    .then(response => response.text())
    .then(data => {
        console.log(data);
        // Update UI or show confirmation message
    })
    .catch(error => {
        console.error('Error:', error);
    });
}